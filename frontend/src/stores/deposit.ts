import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'

export const useDepositStore = defineStore('deposit', () => {
    const depots = ref<any[]>([])
    const sites = ref<any[]>([])
    const lockerAccesses = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
    }

    async function fetchDepots() {
        const authStore = useAuthStore()
        isLoading.value = true
        try {
            const res = await fetch(`${API_BASE}/items/me`, { headers: authHeaders() })
            if (!res.ok) throw new Error('Erreur chargement dépôts')
            depots.value = await res.json()
        } catch (e: any) {
            error.value = e.message
        } finally {
            isLoading.value = false
        }
    }

    async function fetchSites() {
        try {
            const res = await fetch(`${API_BASE}/sites`, { headers: authHeaders() })
            if (!res.ok) return
            sites.value = await res.json()
        } catch {
            sites.value = []
        }
    }

    async function fetchLockerAccesses() {
        try {
            const res = await fetch(`${API_BASE}/users/me/locker-access`, { headers: authHeaders() })
            if (res.ok) {
                const data = await res.json()
                lockerAccesses.value = Array.isArray(data) ? data : []
            }
        } catch (e) {
            console.error('Fetch Locker Accesses Error:', e)
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

    return {
        depots,
        sites,
        lockerAccesses,
        isLoading,
        error,
        fetchDepots,
        fetchSites,
        fetchLockerAccesses,
        createItem
    }
})