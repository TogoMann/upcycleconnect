import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { useCartStore } from './cart'
import { API_BASE } from '@/config'

export const useListingStore = defineStore('listing', () => {
    const annonces = ref<any[]>([])
    const allAnnonces = ref<any[]>([])
    const conversations = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
    }

    async function fetchAnnonces() {
        isLoading.value = true
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/listing/me`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement annonces')
            annonces.value = await res.json()
        } catch (e: any) {
            console.error('Fetch Annonces Error:', e)
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function fetchAllAnnonces() {
        isLoading.value = true
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/listing`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement annonces')
            allAnnonces.value = await res.json()
        } catch (e: any) {
            console.error('Fetch All Annonces Error:', e)
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function uploadImage(file: File): Promise<string> {
        const formData = new FormData()
        formData.append('image', file)

        const res = await fetch(`${API_BASE}/listing/upload`, {
            method: 'POST',
            headers: authHeaders(),
            body: formData,
        })
        if (!res.ok) {
            const errData = await res.text()
            throw new Error(errData || 'Erreur lors de l\'envoi de l\'image')
        }
        const data = await res.json()
        return data.url
    }

    async function fetchSitesWithLockers(cityId: number): Promise<any[]> {
        const res = await fetch(`${API_BASE}/sites/lockers?city_id=${cityId}`, { headers: authHeaders() })
        if (!res.ok) throw new Error('Erreur chargement des sites de dépôt')
        const data = await res.json()
        return Array.isArray(data) ? data : []
    }

    async function createAnnonce(data: { name: string; description: string; price: number; category?: string; city_id?: number; image_url?: string; handoff_mode?: string; address?: string; weight?: number; locker_id?: number; physical_state?: string; size?: string }) {
        const res = await fetch(`${API_BASE}/listing/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        })
        if (!res.ok) {
            const errText = await res.text()
            throw new Error(errText || 'Erreur création annonce')
        }
        return await res.json()
    }

    async function deleteAnnonce(id: number) {
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/listing/${id}`, {
                method: 'DELETE',
                headers: authHeaders(),
            })
            if (!res.ok) {
                const errText = await res.text()
                throw new Error(errText || 'Erreur suppression annonce')
            }
            annonces.value = annonces.value.filter((a: any) => (a.id?.Int64 ?? a.id) !== id)
        } catch (e: any) {
            console.error('Delete Annonce Error:', e)
            error.value = e.message
            throw e
        }
    }

    async function createOrderCheckout(listingId: number): Promise<{ free: boolean; url?: string; order_id?: number }> {
        const cartStore = useCartStore()
        const res = await fetch(`${API_BASE}/listing-order/checkout`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({ listing_id: listingId }),
        })
        if (!res.ok) {
            const errText = await res.text()
            throw new Error(errText || 'Erreur lors de la création de la commande')
        }

        const data = await res.json()

        try {
            await cartStore.removeFromCart('listing', listingId)
        } catch (e) {
            console.warn('Could not remove from cart after order:', e)
        }

        return data
    }

    return {
        annonces,
        allAnnonces,
        conversations,
        isLoading,
        error,
        fetchAnnonces,
        fetchAllAnnonces,
        uploadImage,
        fetchSitesWithLockers,
        createAnnonce,
        createOrderCheckout,
        deleteAnnonce
    }
})