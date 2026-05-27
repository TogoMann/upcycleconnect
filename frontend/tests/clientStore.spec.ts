import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useClientStore } from '../src/stores/client'
import { useAuthStore } from '../src/stores/auth'

// Mock global fetch
global.fetch = vi.fn()

describe('Client Store (Cart Logic)', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        
        // Mock auth store
        const auth = useAuthStore()
        auth.token = 'test-token'
        auth.user = { id: 5, username: 'testuser' }
    })

    it('fetches cart successfully', async () => {
        const mockCart = [
            { id: 1, listing_id: { Int64: 10, Valid: true }, listing: { name: 'Object' } }
        ]
        
        ;(fetch as any).mockResolvedValue({
            ok: true,
            json: () => Promise.resolve(mockCart)
        })

        const store = useClientStore()
        await store.fetchCart()

        expect(store.cart).toEqual(mockCart)
    })

    it('adds item to cart and refreshes', async () => {
        ;(fetch as any).mockResolvedValue({ ok: true, json: () => Promise.resolve([]) })

        const store = useClientStore()
        await store.addToCart({ listingId: 10 })

        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart'), expect.objectContaining({
            method: 'POST',
            body: JSON.stringify({ listing_id: 10 })
        }))
        // Refresh call (fetchCart) doesn't specify method 'GET' explicitly in code
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart'), expect.objectContaining({
            headers: expect.objectContaining({ Authorization: 'Bearer test-token' })
        }))
    })

    it('removes item from cart and refreshes', async () => {
        ;(fetch as any).mockResolvedValue({ ok: true, json: () => Promise.resolve([]) })

        const store = useClientStore()
        await store.removeFromCart('listing', 10)

        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart/listing/10'), expect.objectContaining({
            method: 'DELETE'
        }))
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart'), expect.anything())
    })

    it('checkouts direct-pay items and refreshes', async () => {
        ;(fetch as any).mockResolvedValue({ 
            ok: true, 
            json: () => Promise.resolve({ message: 'Success' }) 
        })

        const store = useClientStore()
        await store.checkoutCart()

        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart/checkout'), expect.objectContaining({
            method: 'POST'
        }))
        // Parallel refresh calls
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart'), expect.anything())
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/event-participation'), expect.anything())
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/course-order/me'), expect.anything())
    })
})
