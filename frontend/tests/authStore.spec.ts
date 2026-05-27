import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../src/stores/auth'

// Mock global fetch
global.fetch = vi.fn()

describe('Auth Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        localStorage.clear()
    })

    it('logs in successfully', async () => {
        const mockToken = 'fake.token.here'
        const mockUser = { id: 1, username: 'testuser' }

        ;(fetch as any).mockImplementation((url: string) => {
            if (url.includes('/login')) {
                return Promise.resolve({
                    ok: true,
                    json: () => Promise.resolve({ token: mockToken })
                })
            }
            if (url.includes('/users/me')) {
                return Promise.resolve({
                    ok: true,
                    json: () => Promise.resolve(mockUser)
                })
            }
            return Promise.reject(new Error('Unknown URL'))
        })

        const store = useAuthStore()
        await store.login('testuser', 'password')

        expect(store.token).toBe(mockToken)
        expect(store.user).toEqual(mockUser)
        expect(localStorage.getItem('auth_token')).toBe(mockToken)
    })

    it('fails login with invalid credentials', async () => {
        ;(fetch as any).mockResolvedValue({
            ok: false,
            status: 401
        })

        const store = useAuthStore()
        await expect(store.login('wrong', 'pass')).rejects.toThrow('Identifiants invalides')
        expect(store.token).toBeNull()
    })

    it('logs out and clears state', async () => {
        const store = useAuthStore()
        store.token = 'some-token'
        store.user = { id: 1 } as any
        localStorage.setItem('auth_token', 'some-token')

        store.logout()

        expect(store.token).toBeNull()
        expect(store.user).toBeNull()
        expect(localStorage.getItem('auth_token')).toBeNull()
    })

    it('registers successfully', async () => {
        const mockToken = 'new.token.here'
        ;(fetch as any).mockImplementation((url: string) => {
            if (url.includes('/register')) {
                return Promise.resolve({
                    ok: true,
                    json: () => Promise.resolve({ token: mockToken })
                })
            }
            return Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
        })

        const store = useAuthStore()
        await store.register('user', 'First', 'Last', 'email@test.com', 'pass')

        expect(store.token).toBe(mockToken)
    })
})
