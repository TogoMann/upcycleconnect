import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAdminStore } from '../src/stores/admin'

global.fetch = vi.fn()

describe('Admin Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('fetches users successfully', async () => {
        const mockUsers = [{ id: 1, username: 'admin', role: 'admin' }]
        const mockResponse = { data: mockUsers, total_pages: 2, page: 1 }
        
        ;(fetch as any).mockResolvedValueOnce({
            ok: true,
            json: () => Promise.resolve(mockResponse)
        })

        const store = useAdminStore()
        await store.fetchUsers()

        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/users'), expect.any(Object))
        expect(store.users).toEqual(mockUsers)
        expect(store.usersTotalPages).toBe(2)
        expect(store.usersCurrentPage).toBe(1)
        expect(store.isLoading).toBe(false)
    })

    it('handles fetch users error', async () => {
        ;(fetch as any).mockResolvedValueOnce({
            ok: false,
            text: () => Promise.resolve('Failed to fetch users')
        })

        const store = useAdminStore()
        await store.fetchUsers()

        expect(store.error).toBe('Failed to fetch users')
        expect(store.users).toEqual([])
    })

    it('approves an event', async () => {
        const store = useAdminStore()
        store.events = [{ id: 1, title: 'Event 1', approved: false }] as any[]
        
        ;(fetch as any).mockResolvedValueOnce({ ok: true })

        await store.approveEvent(1)

        expect(fetch).toHaveBeenCalledWith(
            expect.stringContaining('/event/1/approve'), 
            expect.objectContaining({ method: 'PATCH' })
        )
        
        expect(store.events[0].approved).toBe(true)
    })

    it('deletes a user and removes it from state', async () => {
        const store = useAdminStore()
        store.users = [{ id: 1, username: 'user1' }, { id: 2, username: 'user2' }] as any[]
        
        ;(fetch as any).mockResolvedValueOnce({ ok: true })

        await store.deleteUser(1)

        expect(fetch).toHaveBeenCalledWith(
            expect.stringContaining('/users/1'), 
            expect.objectContaining({ method: 'DELETE' })
        )
        
        expect(store.users.length).toBe(1)
        expect(store.users[0].id).toBe(2)
    })
})
