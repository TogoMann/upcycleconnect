import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'

export const useCatalogueStore = defineStore('catalogue', () => {
    const events = ref<any[]>([])
    const courses = ref<any[]>([])
    const participations = ref<any[]>([])
    const courseOrders = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
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

    async function fetchParticipations() {
        try {
            const res = await fetch(`${API_BASE}/event-participation`, { headers: authHeaders() })
            if (res.ok) {
                const data = await res.json()
                if (!Array.isArray(data)) {
                    participations.value = []
                    return
                }
                const authStore = useAuthStore()
                participations.value = data.filter((p: any) => {
                    const pUid = p.user_id && typeof p.user_id === 'object' ? p.user_id.Int64 : p.user_id
                    return Number(pUid) === authStore.user?.id
                })
            }
        } catch (e) {
            console.error('Fetch Participations Error:', e)
        }
    }

    async function fetchCourseOrders() {
        try {
            const res = await fetch(`${API_BASE}/course-order/me`, { headers: authHeaders() })
            if (res.ok) {
                const data = await res.json()
                courseOrders.value = Array.isArray(data) ? data : []
            }
        } catch (e) {
            console.error('Fetch Course Orders Error:', e)
        }
    }

    async function createCourseOrder(courseId: number, price: number) {
        const res = await fetch(`${API_BASE}/course-order/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({
                course_id: courseId,
                price: price,
            }),
        })
        if (!res.ok) throw new Error('Erreur inscription atelier')
        return await res.json()
    }

    async function createEventParticipation(eventId: number) {
        const res = await fetch(`${API_BASE}/event-participation/`, {
            method: 'POST',
            headers: { ...authHeaders(), 'Content-Type': 'application/json' },
            body: JSON.stringify({
                event_id: eventId,
            }),
        })
        if (!res.ok) throw new Error('Erreur inscription événement')
        return await res.json()
    }

    return {
        events,
        courses,
        participations,
        courseOrders,
        isLoading,
        error,
        fetchCatalogue,
        fetchParticipations,
        fetchCourseOrders,
        createCourseOrder,
        createEventParticipation
    }
})