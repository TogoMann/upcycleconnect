import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'
export interface Conversation {
  id: number
  listing_id: number
  buyer_id: number
  seller_id: number
  is_closed: boolean
  created_at: any
  updated_at: any
  listing_title?: string
  buyer_name?: string
  seller_name?: string
}

export interface Message {
  id: number
  conversation_id: number
  sender_id: number
  content: string
  message_type: 'text' | 'price_proposal'
  proposed_price: number | null
  proposal_status: 'pending' | 'accepted' | 'declined' | null
  created_at: string
  is_edited: boolean
}

export interface EditHistory {
  id: number
  message_id: number
  old_content: string
  old_proposed_price: number | null
  edited_at: string
}

export const useChatStore = defineStore('chat', () => {
  const authStore = useAuthStore()

  async function sendMessage(
    listingId: number,
    content: string,
    type: 'text' | 'price_proposal' = 'text',
    price?: number,
    conversationId?: number,
  ) {
    const res = await fetch(`${API_BASE}/chat/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${authStore.token}`,
      },
      body: JSON.stringify({
        listing_id: listingId,
        conversation_id: conversationId,
        content: content,
        message_type: type,
        proposed_price: price,
      }),
    })
    if (!res.ok) throw new Error('Failed to send message')
    return await res.json()
  }

  async function getConversations() {
    const res = await fetch(`${API_BASE}/chat/conversations`, {
      headers: { Authorization: `Bearer ${authStore.token}` },
    })
    if (!res.ok) throw new Error('Failed to fetch conversations')
    return await res.json()
  }

  async function getMessages(convId: number) {
    const res = await fetch(`${API_BASE}/chat/conversations/${convId}/messages`, {
      headers: { Authorization: `Bearer ${authStore.token}` },
    })
    if (!res.ok) throw new Error('Failed to fetch messages')
    return await res.json()
  }

  async function handleProposal(messageId: number, accept: boolean) {
    const res = await fetch(`${API_BASE}/chat/messages/${messageId}/proposal`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${authStore.token}`,
      },
      body: JSON.stringify({ accept }),
    })
    if (!res.ok) throw new Error('Failed to handle proposal')
  }

  async function editMessage(messageId: number, content: string, price?: number) {
    const res = await fetch(`${API_BASE}/chat/messages/${messageId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${authStore.token}`,
      },
      body: JSON.stringify({
        content: content,
        proposed_price: price,
      }),
    })
    if (!res.ok) throw new Error('Failed to edit message')
  }

  // Admin Methods
  async function adminGetConversations() {
    const res = await fetch(`${API_BASE}/admin/chat/conversations`, {
      headers: { Authorization: `Bearer ${authStore.token}` },
    })
    if (!res.ok) throw new Error('Failed to fetch admin conversations')
    return await res.json()
  }

  async function adminGetConversationDetails(convId: number) {
    const res = await fetch(`${API_BASE}/admin/chat/conversations/${convId}`, {
      headers: { Authorization: `Bearer ${authStore.token}` },
    })
    if (!res.ok) throw new Error('Failed to fetch conversation details')
    return await res.json()
  }

  async function adminCensorMessage(messageId: number) {
    const res = await fetch(`${API_BASE}/admin/chat/messages/${messageId}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${authStore.token}` },
    })
    if (!res.ok) throw new Error('Failed to censor message')
  }

  return {
    sendMessage,
    getConversations,
    getMessages,
    handleProposal,
    editMessage,
    adminGetConversations,
    adminGetConversationDetails,
    adminCensorMessage,
  }
})
