<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useChatStore, type Conversation, type Message } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const chatStore = useChatStore()
const authStore = useAuthStore()
const route = useRoute()

const conversations = ref<Conversation[]>([])
const selectedConv = ref<Partial<Conversation> | null>(null)
const messages = ref<Message[]>([])
const newMessage = ref('')

function getNumericId(val: any): number | undefined {
    if (val === null || val === undefined) return undefined
    if (typeof val === 'object') return val.Int64
    return val
}

function isCourseConversation(conv: Partial<Conversation>): boolean {
    return !!getNumericId(conv.course_id)
}

async function loadConversations() {
    try {
        const all = await chatStore.getConversations()
        conversations.value = all.filter((c: Conversation) => !!getNumericId(c.course_id))

        const queryCourseId = route.query.courseId
        if (queryCourseId) {
            const cid = parseInt(queryCourseId as string)
            const existing = conversations.value.find(c => getNumericId(c.course_id) === cid)
            if (existing) {
                selectConversation(existing)
            } else {
                selectedConv.value = { course_id: cid, id: 0 }
                messages.value = []
            }
        }
    } catch (e) {
        console.error(e)
    }
}

async function selectConversation(conv: Partial<Conversation>) {
    selectedConv.value = conv
    if (conv.id && conv.id > 0) {
        try {
            messages.value = await chatStore.getMessages(conv.id)
        } catch (e) {
            console.error(e)
        }
    } else {
        messages.value = []
    }
}

async function send() {
    if (!selectedConv.value || !newMessage.value) return

    try {
        const courseId = getNumericId(selectedConv.value.course_id)!
        const convId = getNumericId(selectedConv.value.id)
        await chatStore.sendMessage(undefined, newMessage.value, 'text', undefined, convId, courseId)

        newMessage.value = ''

        if (selectedConv.value.id && selectedConv.value.id > 0) {
            messages.value = await chatStore.getMessages(selectedConv.value.id)
        } else {
            await loadConversations()
        }
    } catch (e) {
        console.error(e)
    }
}

onMounted(loadConversations)

watch(() => route.query.courseId, () => {
    loadConversations()
})
</script>

<template>
    <div class="salarie-chat">
        <div class="page-header">
            <h1 class="page-title">{{ t('salarie.chat.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('salarie.chat.subtitle') }}</p>
        </div>

        <div class="chat-container">
            <div class="conv-list">
                <div class="conv-list-header">
                    <h3>{{ t('salarie.chat.conversations') }}</h3>
                </div>
                <div v-if="conversations.length === 0" class="empty-conv">
                    {{ t('salarie.chat.noConversations') }}
                </div>
                <div v-for="conv in conversations" :key="conv.id"
                    class="conv-item" :class="{ active: selectedConv?.id === conv.id }"
                    @click="selectConversation(conv)">
                    <div class="conv-item-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                        </svg>
                    </div>
                    <div class="conv-item-info">
                        <p class="conv-item-title">{{ conv.course_title || t('salarie.chat.training', { id: getNumericId(conv.course_id) }) }}</p>
                        <small class="conv-item-sender">{{ conv.buyer_name || t('salarie.chat.participant') }}</small>
                        <small class="conv-item-date">{{ new Date(conv.updated_at.Time ?? conv.updated_at).toLocaleDateString() }}</small>
                    </div>
                </div>
            </div>

            <div class="chat-window" v-if="selectedConv">
                <div class="chat-header">
                    <h3>{{ selectedConv.course_title || t('salarie.chat.training', { id: getNumericId(selectedConv.course_id) }) }}</h3>
                    <span class="chat-header-sender">{{ t('salarie.chat.withParticipant', { name: selectedConv.buyer_name || t('salarie.chat.aParticipant') }) }}</span>
                </div>
                <div class="messages">
                    <div v-for="msg in messages" :key="msg.id"
                        class="message" :class="{ mine: msg.sender_id === authStore.user?.id }">
                        <div class="bubble">
                            <p v-if="msg.content">{{ msg.content }}</p>
                            <small class="time">
                                {{ new Date(msg.created_at.Time ?? msg.created_at).toLocaleTimeString() }}
                            </small>
                        </div>
                    </div>
                </div>

                <div class="input-area">
                    <textarea v-model="newMessage" :placeholder="t('salarie.chat.responsePlaceholder')"></textarea>
                    <div class="controls">
                        <button @click="send" class="btn-send">{{ t('salarie.chat.send') }}</button>
                    </div>
                </div>
            </div>
            <div class="no-chat" v-else>
                <div class="no-chat-content">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                    </svg>
                    <p>{{ t('salarie.chat.selectPrompt') }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.salarie-chat { display: flex; flex-direction: column; height: calc(100vh - 120px); }
.page-header { margin-bottom: 24px; }
.page-title { font-size: 2rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.02em; margin: 0 0 4px; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.chat-container { display: grid; grid-template-columns: 320px 1fr; flex: 1; background: var(--white); border: 1.5px solid rgba(53, 53, 53, 0.1); border-radius: 14px; overflow: hidden; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05); }
.conv-list { border-right: 1.5px solid rgba(53, 53, 53, 0.1); display: flex; flex-direction: column; background: #fafafa; }
.conv-list-header { padding: 20px; border-bottom: 1px solid rgba(53, 53, 53, 0.05); }
.conv-list-header h3 { margin: 0; font-size: 1rem; font-weight: 700; color: var(--charcoal); }
.empty-conv { padding: 40px 20px; text-align: center; color: rgba(53, 53, 53, 0.4); font-size: 0.85rem; }
.conv-item { padding: 16px 20px; display: flex; align-items: center; gap: 14px; cursor: pointer; transition: background 0.2s; border-bottom: 1px solid rgba(53, 53, 53, 0.03); }
.conv-item:hover { background: rgba(52, 137, 91, 0.05); }
.conv-item.active { background: var(--green-pale); border-right: 3px solid var(--green-mid); }
.conv-item-icon { width: 40px; height: 40px; background: var(--white); border: 1px solid rgba(53, 53, 53, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: var(--green-mid); flex-shrink: 0; }
.conv-item-icon svg { width: 20px; height: 20px; }
.conv-item-info { flex: 1; overflow: hidden; }
.conv-item-title { margin: 0 0 2px; font-size: 0.9rem; font-weight: 600; color: var(--charcoal); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.conv-item-sender { display: block; font-size: 0.78rem; color: var(--green-dark); font-weight: 600; margin-bottom: 2px; }
.conv-item-date { font-size: 0.75rem; color: rgba(53, 53, 53, 0.5); }
.chat-window { display: flex; flex-direction: column; background: var(--white); }
.chat-header { padding: 18px 24px; border-bottom: 1px solid rgba(53, 53, 53, 0.08); display: flex; flex-direction: column; gap: 2px; }
.chat-header h3 { margin: 0; font-size: 1.1rem; font-weight: 700; color: var(--charcoal); }
.chat-header-sender { font-size: 0.8rem; color: var(--charcoal); opacity: 0.6; }
.messages { flex: 1; overflow-y: auto; padding: 24px; display: flex; flex-direction: column; gap: 16px; background: #fcfcfc; }
.message { display: flex; justify-content: flex-start; }
.message.mine { justify-content: flex-end; }
.bubble { max-width: 75%; padding: 12px 18px; border-radius: 16px; background: var(--white); box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05); border: 1px solid rgba(53, 53, 53, 0.05); position: relative; color: var(--charcoal); }
.message.mine .bubble { background: var(--green-dark); color: var(--white); border: none; }
.bubble p { margin: 0; font-size: 0.95rem; line-height: 1.5; }
.time { font-size: 0.7rem; opacity: 0.6; display: block; margin-top: 8px; }
.message.mine .time { opacity: 0.8; }
.input-area { padding: 24px; border-top: 1px solid rgba(53, 53, 53, 0.08); }
textarea { width: 100%; height: 80px; margin-bottom: 12px; padding: 14px; border-radius: 10px; border: 1.5px solid rgba(53, 53, 53, 0.12); font-family: inherit; font-size: 0.95rem; resize: none; outline: none; transition: border-color 0.2s; box-sizing: border-box; }
textarea:focus { border-color: var(--green-mid); }
.controls { display: flex; gap: 12px; align-items: center; }
.btn-send { background: var(--green-dark); color: var(--white); border: none; padding: 10px 24px; border-radius: 8px; font-weight: 700; cursor: pointer; margin-left: auto; transition: background 0.2s; }
.btn-send:hover { background: var(--green-mid); }
.no-chat { display: flex; align-items: center; justify-content: center; background: #fafafa; text-align: center; padding: 40px; }
.no-chat-content svg { width: 64px; height: 64px; color: rgba(53, 53, 53, 0.1); margin-bottom: 16px; }
.no-chat-content p { font-size: 1rem; color: rgba(53, 53, 53, 0.4); max-width: 300px; line-height: 1.6; }
</style>
