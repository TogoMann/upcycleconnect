import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import Catalogue from '../src/pages/client/Catalogue.vue'
import { useClientStore } from '../src/stores/client'

// Mock router
const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/particulier/catalogue', component: { template: '<div></div>' } }]
})

describe('Catalogue Component', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('renders empty state when no items available', async () => {
        const store = useClientStore()
        store.events = []
        store.courses = []
        store.allAnnonces = []
        store.cart = []

        const wrapper = mount(Catalogue, {
            global: {
                plugins: [router]
            }
        })

        expect(wrapper.text()).toContain('Aucun élément disponible')
    })

    it('filters out items already in cart', async () => {
        const store = useClientStore()
        
        // Item in catalogue
        store.allAnnonces = [{ id: 10, name: 'Annonce 1', _type: 'annonce' }]
        
        // Same item in cart
        store.cart = [{ id: 1, listing_id: { Int64: 10, Valid: true } }]

        const wrapper = mount(Catalogue, {
            global: {
                plugins: [router]
            }
        })

        // Should not find the card for the item in cart
        expect(wrapper.text()).not.toContain('Annonce 1')
        expect(wrapper.text()).toContain('Aucun élément disponible')
    })

    it('shows items not in cart', async () => {
        const store = useClientStore()
        
        store.allAnnonces = [
            { id: 10, name: 'Annonce 1', description: 'Desc 1', price: 10 },
            { id: 11, name: 'Annonce 2', description: 'Desc 2', price: 20 }
        ]
        
        // Only one item in cart
        store.cart = [{ id: 1, listing_id: { Int64: 10, Valid: true } }]

        const wrapper = mount(Catalogue, {
            global: {
                plugins: [router]
            }
        })

        expect(wrapper.text()).not.toContain('Annonce 1')
        expect(wrapper.text()).toContain('Annonce 2')
    })
})
