import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { useCatalogueStore } from './catalogue'
import { API_BASE } from '@/config'

export const useCartStore = defineStore('cart', () => {
    const cart = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
    }

    async function fetchCart() {
        isLoading.value = true
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/cart`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement panier')
            const data = await res.json()
            cart.value = Array.isArray(data) ? data : []
        } catch (e: any) {
            console.error('Fetch Cart Error:', e)
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function addToCart(ids: { listingId?: number; eventId?: number; courseId?: number }) {
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/cart`, {
                method: 'POST',
                headers: { ...authHeaders(), 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    listing_id: ids.listingId,
                    event_id: ids.eventId,
                    course_id: ids.courseId,
                }),
            })
            if (!res.ok) throw new Error('Erreur ajout au panier')
            await fetchCart()
        } catch (e: any) {
            console.error('Add To Cart Error:', e)
            error.value = e.message
            throw e
        }
    }

    async function removeFromCart(type: 'listing' | 'event' | 'course', id: number) {
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/cart/${type}/${id}`, {
                method: 'DELETE',
                headers: authHeaders(),
            })
            if (!res.ok) throw new Error('Erreur retrait du panier')
            await fetchCart()
        } catch (e: any) {
            console.error('Remove From Cart Error:', e)
            error.value = e.message
            throw e
        }
    }

    async function clearCart() {
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/cart`, {
                method: 'DELETE',
                headers: authHeaders(),
            })
            if (!res.ok) throw new Error('Erreur lors de la suppression du panier')
            await fetchCart()
        } catch (e: any) {
            console.error('Clear Cart Error:', e)
            error.value = e.message
            throw e
        }
    }

    async function checkoutCart() {
        const catalogueStore = useCatalogueStore()
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/cart/checkout`, {
                method: 'POST',
                headers: authHeaders(),
            })
            if (!res.ok) throw new Error('Erreur paiement panier')
            await Promise.all([
                fetchCart(),
                catalogueStore.fetchParticipations(),
                catalogueStore.fetchCourseOrders()
            ])
            return await res.json()
        } catch (e: any) {
            console.error('Checkout Error:', e)
            error.value = e.message
            throw e
        }
    }

    return {
        cart,
        isLoading,
        error,
        fetchCart,
        addToCart,
        removeFromCart,
        clearCart,
        checkoutCart
    }
})