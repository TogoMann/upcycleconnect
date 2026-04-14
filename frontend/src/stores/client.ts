import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

const API_BASE = 'http://localhost:8081'

export const useClientStore = defineStore('client', () => {
    const annonces = ref<any[]>([])
    const depots = ref<any[]>([])
    const entries = ref<any[]>([])
    const score = ref<number>(0)
    const events = ref<any[]>([])
    const courses = ref<any[]>([])
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
            const res = await fetch(`${API_BASE}/listing`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement annonces')
            const all = await res.json()
            const authStore = useAuthStore()
            annonces.value = all.filter((a: any) => a.created_by?.Int64 === authStore.user?.id)
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    const allAnnonces = ref<any[]>([])
    async function fetchAllAnnonces() {
        isLoading.value = true
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/listing`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement annonces')
            allAnnonces.value = await res.json()
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function createAnnonce(data: { name: string; description: string; price: number }) {
        const authStore = useAuthStore()
        const res = await fetch(`${API_BASE}/listing/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({ ...data, created_by: { Int64: authStore.user?.id, Valid: true } }),
        })
        if (!res.ok) throw new Error('Erreur création annonce')
        return await res.json()
    }

    async function deleteAnnonce(id: number) {
        const res = await fetch(`${API_BASE}/listing/${id}`, {
            method: 'DELETE',
            headers: authHeaders(),
        })
        if (!res.ok) throw new Error('Erreur suppression annonce')
        annonces.value = annonces.value.filter((a: any) => a.id?.Int64 !== id)
    }

    async function fetchDepots() {
        isLoading.value = true
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/items`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement dépôts')
            const all = await res.json()
            const authStore = useAuthStore()
            depots.value = all.filter((i: any) => i.owner_id?.Int64 === authStore.user?.id)
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function fetchScore() {
        const authStore = useAuthStore()
        if (!authStore.user) return
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/users/${authStore.user.id}/score`, {
                headers: authHeaders(),
            })
            if (!res.ok) throw new Error('Erreur score')
            const data = await res.json()
            score.value = data.score
        } catch (e: any) {
            error.value = e.message
        }
    }

    async function fetchEntries() {
        isLoading.value = true
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/entry`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur planning')
            const all = await res.json()
            const authStore = useAuthStore()
            entries.value = all.filter((e: any) => e.created_by?.Int64 === authStore.user?.id)
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function createEntry(data: { schedule: string; start: string; ending: string }) {
        const authStore = useAuthStore()
        const res = await fetch(`${API_BASE}/entry/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({ ...data, created_by: { Int64: authStore.user?.id, Valid: true } }),
        })
        if (!res.ok) throw new Error('Erreur création créneau')
        return await res.json()
    }

    async function deleteEntry(id: number) {
        const res = await fetch(`${API_BASE}/entry/${id}`, {
            method: 'DELETE',
            headers: authHeaders(),
        })
        if (!res.ok) throw new Error('Erreur suppression créneau')
        entries.value = entries.value.filter((e: any) => e.id?.Int64 !== id)
    }

    async function fetchCatalogue() {
        isLoading.value = true
        error.value = null
        try {
            const [evtRes, crsRes] = await Promise.all([
                fetch(`${API_BASE}/event`, { headers: authHeaders() }),
                fetch(`${API_BASE}/course`, { headers: authHeaders() }),
            ])
            events.value = evtRes.ok ? await evtRes.json() : []
            courses.value = crsRes.ok ? await crsRes.json() : []
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function createOrder(listingId: number, price: number) {
        const authStore = useAuthStore()
        const res = await fetch(`${API_BASE}/listing-order/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({
                listing_id: { Int64: listingId, Valid: true },
                user_id: { Int64: authStore.user?.id, Valid: true },
                price: price,
                status: 'pending',
            }),
        })
        if (!res.ok) throw new Error('Erreur création commande')
        return await res.json()
    }

    async function markTutorialSeen() {
        const authStore = useAuthStore()
        if (!authStore.user) return
        const res = await fetch(`${API_BASE}/users/${authStore.user.id}/tutorial`, {
            method: 'PATCH',
            headers: authHeaders(),
        })
        if (res.ok && authStore.user) {
            authStore.user.has_seen_tutorial = true
        }
    }

    return {
        annonces,
        allAnnonces,
        depots,
        entries,
        score,
        events,
        courses,
        isLoading,
        error,
        fetchAnnonces,
        fetchAllAnnonces,
        createAnnonce,
        deleteAnnonce,
        fetchDepots,
        fetchScore,
        fetchEntries,
        createEntry,
        deleteEntry,
        fetchCatalogue,
        createOrder,
        markTutorialSeen,
    }
})
