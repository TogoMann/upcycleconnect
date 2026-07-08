import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useChatStore } from '../src/stores/chat'
import { useAuthStore } from '../src/stores/auth'


global.fetch = vi.fn()

describe('Chat Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        
        const auth = useAuthStore()
        auth.token = 'test-token'
    })

    it('sends message successfully', async () => {
        ;(fetch as any).mockResolvedValue({
            ok: true,
            json: () => Promise.resolve({ id: 100 })
        })

        const store = useChatStore()
        const result = await store.sendMessage(1, 'Hello', 'text')

        expect(result.id).toBe(100)
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/chat/messages'), expect.objectContaining({
            method: 'POST',
            body: expect.stringContaining('"content":"Hello"')
        }))
    })

    it('fetches conversations', async () => {
        const mockConvs = [{ id: 1, listing_title: 'Object' }]
        ;(fetch as any).mockResolvedValue({
            ok: true,
            json: () => Promise.resolve(mockConvs)
        })

        const store = useChatStore()
        const convs = await store.getConversations()

        expect(convs).toEqual(mockConvs)
        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/chat/conversations'), expect.anything())
    })

    it('handles price proposal', async () => {
        ;(fetch as any).mockResolvedValue({ ok: true })

        const store = useChatStore()
        await store.handleProposal(50, true)

        expect(fetch).toHaveBeenCalledWith(expect.stringContaining('/chat/messages/50/proposal'), expect.objectContaining({
            method: 'POST',
            body: JSON.stringify({ accept: true })
        }))
    })
})
