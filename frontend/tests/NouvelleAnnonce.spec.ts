import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { createI18n } from 'vue-i18n'
import NouvelleAnnonce from '../src/pages/client/NouvelleAnnonce.vue'
import { useClientStore } from '../src/stores/client'
import fr from '../src/i18n/locales/fr.json'
import en from '../src/i18n/locales/en.json'

const i18n = createI18n({ legacy: false, locale: 'fr', fallbackLocale: 'fr', messages: { fr, en } })

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/particulier/annonces/nouvelle', component: { template: '<div></div>' } },
        { path: '/particulier/annonces', component: { template: '<div></div>' } }
    ]
})

describe('NouvelleAnnonce Component', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('shows validation errors when submitting empty form', async () => {
        const wrapper = mount(NouvelleAnnonce, {
            global: { plugins: [router, i18n] }
        })

        await wrapper.find('form').trigger('submit.prevent')

        expect(wrapper.text()).toContain('Le titre est requis')
        expect(wrapper.text()).toContain('La description est requise')
        expect(wrapper.text()).toContain('Un prix valide est requis')
    })

    it('successfully submits listing form', async () => {
        const store = useClientStore()
        store.cities = [{ id: 1, name: 'Paris', zip_code: '75000' }]
        
        const createAnnonceSpy = vi.spyOn(store, 'createAnnonce').mockResolvedValue({ id: 100 })
        const fetchAnnoncesSpy = vi.spyOn(store, 'fetchAnnonces').mockResolvedValue()
        const pushSpy = vi.spyOn(router, 'push')

        const wrapper = mount(NouvelleAnnonce, {
            global: { plugins: [router, i18n] }
        })

        // Fill form
        await wrapper.find('input[type="text"]').setValue('Ma superbe chaise')
        await wrapper.find('textarea').setValue('Une description très longue et détaillée.')
        
        const selects = wrapper.findAll('select')
        await selects[0].setValue('Mobilier')
        await selects[1].setValue('1') // City ID
        
        await wrapper.find('input[type="number"]').setValue('45.00')

        await wrapper.find('form').trigger('submit.prevent')

        expect(createAnnonceSpy).toHaveBeenCalledWith(expect.objectContaining({
            name: 'Ma superbe chaise',
            price: 45
        }))
        expect(fetchAnnoncesSpy).toHaveBeenCalled()
        expect(pushSpy).toHaveBeenCalledWith('/particulier/annonces')
    })
})
