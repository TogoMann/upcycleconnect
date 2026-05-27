import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'

export const usePlanningStore = defineStore('planning', () => {
    const entries = ref<any[]>([])
    const planning = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
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
            headers: authHeaders(),
        })
        if (!res.ok) throw new Error('Erreur suppression créneau')
        entries.value = entries.value.filter((e: any) => e.id?.Int64 !== id)
    }

    return {
        entries,
        planning,
        isLoading,
        error,
        fetchEntries,
        fetchPlanning,
        createPersonalEvent,
        deletePersonalEvent,
        createEntry,
        deleteEntry
    }
})