import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import DeposerObjet from '../src/pages/client/DeposerObjet.vue'
import { useClientStore } from '../src/stores/client'

const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/particulier/deposer', component: { template: '<div></div>' } }]
})

describe('DeposerObjet Component', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('shows validation errors when submitting empty form', async () => {
        const wrapper = mount(DeposerObjet, {
            global: { plugins: [router] }
        })

        await wrapper.find('form').trigger('submit.prevent')

        expect(wrapper.text()).toContain('Le type de matériau est requis')
        expect(wrapper.text()).toContain("L'état de l'objet est requis")
    })

    it('successfully submits form', async () => {
        const store = useClientStore()
        const createItemSpy = vi.spyOn(store, 'createItem').mockResolvedValue({ id: 1 })
        const createEntrySpy = vi.spyOn(store, 'createEntry').mockResolvedValue({ id: 2 })

        const wrapper = mount(DeposerObjet, {
            global: { plugins: [router] }
        })

        // Fill form using indices
        const selects = wrapper.findAll('select')
        await selects[0].setValue('Bois')
        await selects[1].setValue('Neuf')
        
        const inputs = wrapper.findAll('input')
        await inputs[0].setValue('2026-12-31') // Date
        await inputs[1].setValue('10:00')      // Start
        await inputs[2].setValue('11:00')      // End

        await wrapper.find('form').trigger('submit.prevent')

        expect(createItemSpy).toHaveBeenCalled()
        expect(createEntrySpy).toHaveBeenCalled()
        expect(wrapper.text()).toContain('Dépôt enregistré !')
    })
})
