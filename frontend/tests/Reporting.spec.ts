import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Reporting from '@/pages/admin/Reporting.vue'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import fr from '../src/i18n/locales/fr.json'
import en from '../src/i18n/locales/en.json'

const i18n = createI18n({ legacy: false, locale: 'fr', fallbackLocale: 'fr', messages: { fr, en } })

// Mock Chart.js to avoid issues in test environment
vi.mock('chart.js', () => ({
  Chart: {
    register: vi.fn(),
  },
  Title: vi.fn(),
  Tooltip: vi.fn(),
  Legend: vi.fn(),
  BarElement: vi.fn(),
  CategoryScale: vi.fn(),
  LinearScale: vi.fn(),
  ArcElement: vi.fn(),
  PointElement: vi.fn(),
  LineElement: vi.fn(),
}))

vi.mock('vue-chartjs', () => ({
  Bar: { template: '<div>Bar Chart</div>' },
  Pie: { template: '<div>Pie Chart</div>' },
  Doughnut: { template: '<div>Doughnut Chart</div>' },
}))

describe('Reporting.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    
    // Mock global fetch
    global.fetch = vi.fn((url: string) => {
      if (url.includes('/reporting/actors')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([
            { role: 'client', count: 10 },
            { role: 'pro', count: 5 }
          ])
        })
      }
      if (url.includes('/reporting/prestations')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([
            { type: 'event', count: 8 },
            { type: 'course', count: 3 }
          ])
        })
      }
      if (url.includes('/reporting/predictions/distribution')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({
            'event': 10,
            'course': 5
          })
        })
      }
      if (url.includes('/reporting/predictions')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({
            data: [
              { user_id: 1, username: 'testuser', email: 'test@test.com', predicted_service_type: 'event', probability: 0.85 }
            ],
            total_pages: 1
          })
        })
      }
      if (url.includes('/reporting/ml-status')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({
            last_run: '2026-06-25T00:00:00Z',
            total_predictions: 1
          })
        })
      }
      return Promise.reject(new Error('Unknown URL: ' + url))
    }) as any
  })

  it('renders loading state initially', () => {
    const wrapper = mount(Reporting, { global: { plugins: [i18n] } })
    expect(wrapper.text()).toContain('Chargement des données...')
  })

  it('renders charts and table after fetching data', async () => {
    const wrapper = mount(Reporting, { global: { plugins: [i18n] } })
    
    // Wait for async setup
    await new Promise(resolve => setTimeout(resolve, 0))
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Répartition des Acteurs')
    expect(wrapper.text()).toContain('Succès des Prestations')
    expect(wrapper.text()).toContain('Prédictions de Prochaines Prestations')
    
    // Check table content
    expect(wrapper.text()).toContain('testuser')
    expect(wrapper.text()).toContain('test@test.com')
    expect(wrapper.text()).toContain('event')
    expect(wrapper.text()).toContain('85.0%')
  })
})
