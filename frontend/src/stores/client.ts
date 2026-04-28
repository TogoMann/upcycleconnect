import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

const API_BASE = 'http://localhost:8081'

export const useClientStore = defineStore('client', () => {
    const annonces = ref<any[]>([])
    const depots = ref<any[]>([])
    const entries = ref<any[]>([])
    const planning = ref<any[]>([])
    const score = ref<number>(0)
    const scoreHistory = ref<any[]>([])
    const events = ref<any[]>([])
    const courses = ref<any[]>([])
    const cities = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
    }

    async function fetchCities() {
        try {
            const res = await fetch(`${API_BASE}/city`)
            if (!res.ok) return
            cities.value = await res.json()
        } catch {
            cities.value = []
        }
    }

    async function fetchAnnonces() {
        isLoading.value = true
        try {
            const res = await fetch(`${API_BASE}/listing/me`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement annonces')
            annonces.value = await res.json()
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    const allAnnonces = ref<any[]>([])
    async function fetchAllAnnonces() {
        isLoading.value = true
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

    async function createAnnonce(data: { name: string; description: string; price: number; category?: string; city_id?: number; image_url?: string }) {
        const res = await fetch(`${API_BASE}/listing/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        })
        if (!res.ok) throw new Error('Erreur création annonce')
        return await res.json()
    }

    async function fetchDepots() {
        const authStore = useAuthStore()
        isLoading.value = true
        try {
            const res = await fetch(`${API_BASE}/items/owner/${authStore.user?.id}`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement dépôts')
            depots.value = await res.json()
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function fetchScore() {
        const authStore = useAuthStore()
        if (!authStore.user) return
        try {
            const res = await fetch(`${API_BASE}/users/${authStore.user.id}/score`, { headers: authHeaders() })
            if (!res.ok) return
            const data = await res.json()
            score.value = data.score
        } catch {}
    }

    async function fetchScoreHistory() {
        const authStore = useAuthStore()
        if (!authStore.user) return
        isLoading.value = true
        try {
            const res = await fetch(`${API_BASE}/users/${authStore.user.id}/score/history`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement historique points')
            scoreHistory.value = await res.json()
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function fetchEntries() {
        const authStore = useAuthStore()
        isLoading.value = true
        try {
            const res = await fetch(`${API_BASE}/entry`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement créneaux')
            const all: any[] = await res.json()
            entries.value = all.filter((a: any) => a.created_by?.Int64 === authStore.user?.id)
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function fetchPlanning() {
        isLoading.value = true
        try {
            const res = await fetch(`${API_BASE}/planning/me`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement planning')
            planning.value = await res.json()
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function createPersonalEvent(data: { title: string; description: string; date: string; start_time: string; end_time: string }) {
        const res = await fetch(`${API_BASE}/planning/personal`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        })
        if (!res.ok) throw new Error('Erreur création événement')
        await fetchPlanning()
    }

    async function deletePersonalEvent(id: number) {
        const res = await fetch(`${API_BASE}/planning/personal/${id}`, {
            method: 'DELETE',
            headers: authHeaders(),
        })
        if (!res.ok) throw new Error('Erreur suppression événement')
        await fetchPlanning()
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
            headers: headers(),
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
                fetch(`${API_BASE}/course/catalogue`, { headers: authHeaders() }),
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

    async function createCourseOrder(courseId: number, price: number) {
        const authStore = useAuthStore()
        const res = await fetch(`${API_BASE}/course-order/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({
                course_id: courseId,
                buyer_id: authStore.user?.id,
                price: price,
            }),
        })
        if (!res.ok) throw new Error('Erreur inscription atelier')
        return await res.json()
    }

    async function createEventParticipation(eventId: number) {
        const authStore = useAuthStore()
        const res = await fetch(`${API_BASE}/event-participation/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({
                event_id: eventId,
                user_id: authStore.user?.id,
            }),
        })
        if (!res.ok) throw new Error('Erreur inscription événement')
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

    const sites = ref<any[]>([])
    async function fetchSites() {
        try {
            const res = await fetch(`${API_BASE}/sites`, { headers: authHeaders() })
            if (!res.ok) return
            sites.value = await res.json()
        } catch {
            sites.value = []
        }
    }

    async function createItem(data: {
        material_type: string
        physical_state: string
        weight?: number
        site_id?: number
    }) {
        const authStore = useAuthStore()
        const body: any = {
            owner_id: { Int64: authStore.user?.id, Valid: true },
            material_type: data.material_type,
            physical_state: data.physical_state,
            status: 'deposited',
        }
        if (data.site_id) body.site_id = { Int64: data.site_id, Valid: true }
        const res = await fetch(`${API_BASE}/items/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        })
        if (!res.ok) throw new Error('Erreur création objet')
        return await res.json()
    }

    async function updateProfile(data: { first_name: string; last_name: string; email: string }) {
        const authStore = useAuthStore()
        if (!authStore.user) return
        const res = await fetch(`${API_BASE}/users/${authStore.user.id}`, {
            method: 'PATCH',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        })
        if (!res.ok) throw new Error('Erreur mise à jour profil')
        await authStore.fetchCurrentUser()
    }

    return {
        annonces,
        allAnnonces,
        depots,
        entries,
        planning,
        score,
        scoreHistory,
        events,
        courses,
        cities,
        sites,
        isLoading,
        error,
        fetchAnnonces,
        fetchAllAnnonces,
        createAnnonce,
        fetchDepots,
        fetchScore,
        fetchScoreHistory,
        fetchEntries,
        fetchPlanning,
        createPersonalEvent,
        deletePersonalEvent,
        createEntry,
        deleteEntry,
        fetchCatalogue,
        createOrder,
        createCourseOrder,
        createEventParticipation,
        markTutorialSeen,
        fetchSites,
        fetchCities,
        createItem,
        updateProfile,
    }
})