import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCartStore } from '../src/stores/cart'

global.fetch = vi.fn()

describe('Cart Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('fetches cart successfully', async () => {
        const mockCartItems = [{ id: 1, listing_id: 10, price: 150 }]
        
        ;(fetch as any).mockResolvedValueOnce({
            ok: true,
            json: () => Promise.resolve(mockCartItems)
        })

        const store = useCartStore()
        await store.fetchCart()

        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/cart'), expect.any(Object))
        expect(store.cart).toEqual(mockCartItems)
        expect(store.isLoading).toBe(false)
    })

    it('handles fetch cart error', async () => {
        ;(fetch as any).mockResolvedValueOnce({
            ok: false,
            text: () => Promise.resolve('Erreur chargement panier')
        })

        const store = useCartStore()
        await store.fetchCart()

        expect(store.error).toBe('Erreur chargement panier')
        expect(store.cart).toEqual([])
    })

    it('adds to cart successfully', async () => {
        ;(fetch as any)
            .mockResolvedValueOnce({ ok: true }) 
            .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) }) 

        const store = useCartStore()

        await store.addToCart({ listingId: 10 })

        expect(fetch).toHaveBeenCalledWith(
            expect.stringContaining('/cart'), 
            expect.objectContaining({ method: 'POST' })
        )
    })

    it('removes from cart successfully', async () => {
        ;(fetch as any)
            .mockResolvedValueOnce({ ok: true }) 
            .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) }) 

        const store = useCartStore()

        await store.removeFromCart('listing', 10)

        expect(fetch).toHaveBeenCalledWith(
            expect.stringContaining('/cart/listing/10'), 
            expect.objectContaining({ method: 'DELETE' })
        )
    })
})
