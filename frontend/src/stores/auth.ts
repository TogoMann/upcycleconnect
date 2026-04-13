import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

const API_BASE = 'http://localhost:8081'

interface JwtPayload {
    username: string
    role: string
    exp: number
    iat: number
}

interface AuthUser {
    id: number
    username: string
    first_name: string
    last_name: string
    email: string
    role: string
    has_seen_tutorial: boolean
}

export const useAuthStore = defineStore('auth', () => {
    const token = ref<string | null>(localStorage.getItem('auth_token'))
    const user = ref<AuthUser | null>(null)

    const isAuthenticated = computed(() => !!token.value)
    const userRole = computed(() => {
        if (!token.value) return null
        const payload = parseJwt(token.value)
        return payload?.role ?? null
    })

    function parseJwt(t: string): JwtPayload | null {
        try {
            const base64 = t.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
            return JSON.parse(atob(base64))
        } catch {
            return null
        }
    }

    async function login(username: string, password: string) {
        const res = await fetch(`${API_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        })
        if (!res.ok) throw new Error('Identifiants invalides')
        const data = await res.json()
        token.value = data.token
        localStorage.setItem('auth_token', data.token)
        await fetchCurrentUser()
    }

    async function fetchCurrentUser() {
        if (!token.value) return
        const payload = parseJwt(token.value)
        if (!payload) return
        const res = await fetch(`${API_BASE}/users`, {
            headers: { Authorization: `Bearer ${token.value}` },
        })
        if (!res.ok) return
        const users: AuthUser[] = await res.json()
        user.value = users.find(u => u.username === payload.username) ?? null
    }

    function logout() {
        token.value = null
        user.value = null
        localStorage.removeItem('auth_token')
    }

    if (token.value) {
        fetchCurrentUser()
    }

    return { token, user, isAuthenticated, userRole, login, logout, fetchCurrentUser }
})
