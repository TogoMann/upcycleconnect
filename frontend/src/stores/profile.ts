import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'

export const useProfileStore = defineStore('profile', () => {
    const score = ref<number>(0)
    const scoreHistory = ref<any[]>([])
    const quests = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    function authHeaders(): Record<string, string> {
        const authStore = useAuthStore()
        return { Authorization: `Bearer ${authStore.token}` }
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

    async function fetchQuests() {
        try {
            const res = await fetch(`${API_BASE}/users/me/quests`, { headers: authHeaders() })
            if (!res.ok) return
            quests.value = await res.json()
        } catch {}
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
        score,
        scoreHistory,
        quests,
        isLoading,
        error,
        fetchScore,
        fetchScoreHistory,
        fetchQuests,
        markTutorialSeen,
        updateProfile
    }
})