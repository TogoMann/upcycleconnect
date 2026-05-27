import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import Chat from '../src/pages/client/Chat.vue'
import { useClientStore } from '../src/stores/client'
import { useChatStore } from '../src/stores/chat'
import { useAuthStore } from '../src/stores/auth'

// Mock router
const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/particulier/chat', component: { template: '<div></div>' } }]
})

// Mock global fetch to prevent actual API calls during mount
global.fetch = vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve([])
})

describe('Chat Component', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        // Override window.alert for tests
        vi.spyOn(window, 'alert').mockImplementation(() => {})
    })

    it('displays the "Retirer du panier" button if listing is in cart and removes it when clicked', async () => {
        const clientStore = useClientStore()
        const chatStore = useChatStore()
        const authStore = useAuthStore()

        authStore.user = { id: 1, username: 'test' } as any

        // Mock a conversation
        const mockConv = { id: 1, listing_id: 10, buyer_id: 1, seller_id: 2, listing_title: 'Object 10', updated_at: new Date() }
        vi.spyOn(chatStore, 'getConversations').mockResolvedValue([mockConv] as any)
        vi.spyOn(chatStore, 'getMessages').mockResolvedValue([])

        const fetchCartSpy = vi.spyOn(clientStore, 'fetchCart').mockImplementation(async () => {})
        
        // Put the listing in the cart
        clientStore.cart = [{ id: 5, listing_id: { Int64: 10, Valid: true } }]

        const removeFromCartSpy = vi.spyOn(clientStore, 'removeFromCart').mockResolvedValue()

        await router.push('/particulier/chat?listingId=10')

        const wrapper = mount(Chat, {
            global: {
                plugins: [router]
            }
        })

        // Wait for all promises (getConversations, selectConversation, getMessages)
        await flushPromises()
        await wrapper.vm.$nextTick() // Ensure DOM is updated

        // Check if the chat window is rendered
        if (!wrapper.text().includes('Object 10')) {
            console.log('HTML Output:', wrapper.html())
            throw new Error('Conversation was not selected')
        }

        // Check if the button is rendered
        const removeBtn = wrapper.find('.btn-remove-cart')
        if (!removeBtn.exists()) {
            console.log('Cart:', clientStore.cart)
            console.log('SelectedConv:', (wrapper.vm as any).selectedConv)
            console.log('isCurrentListingInCart:', (wrapper.vm as any).isCurrentListingInCart)
            console.log('HTML Output:', wrapper.html())
        }
        
        expect(removeBtn.exists()).toBe(true)
        expect(removeBtn.text()).toContain('Retirer du panier')

        // Click the button
        await removeBtn.trigger('click')
        await flushPromises()

        // Verify the store method was called with the correct arguments
        expect(removeFromCartSpy).toHaveBeenCalledWith('listing', 10)
        expect(window.alert).toHaveBeenCalledWith('Objet retiré du panier avec succès.')
    })

    it('does not display the button if listing is not in cart', async () => {
        const clientStore = useClientStore()
        const chatStore = useChatStore()

        const mockConv = { id: 1, listing_id: 20, buyer_id: 1, seller_id: 2, listing_title: 'Object 20', updated_at: new Date() }
        vi.spyOn(chatStore, 'getConversations').mockResolvedValue([mockConv] as any)

        // Cart is empty
        clientStore.cart = []

        await router.push('/particulier/chat?listingId=20')

        const wrapper = mount(Chat, {
            global: {
                plugins: [router]
            }
        })

        await flushPromises()
        await wrapper.vm.$nextTick()

        const removeBtn = wrapper.find('.btn-remove-cart')
        expect(removeBtn.exists()).toBe(false)
    })
})
