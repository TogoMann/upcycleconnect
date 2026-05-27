import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'

export const useLocationStore = defineStore('location', () => {
    const cities = ref<any[]>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    async function fetchCities() {
        error.value = null
        try {
            const res = await fetch(`${API_BASE}/city`)
            if (!res.ok) return
            cities.value = await res.json()
        } catch (e: any) {
            console.error('Fetch Cities Error:', e)
            cities.value = []
        }
    }

    return {
        cities,
        isLoading,
        error,
        fetchCities,
    }
})
