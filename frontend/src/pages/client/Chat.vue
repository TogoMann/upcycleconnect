<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useChatStore, type Conversation, type Message } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'

const chatStore = useChatStore()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const conversations = ref<Conversation[]>([])
const selectedConv = ref<Partial<Conversation> | null>(null)
const messages = ref<Message[]>([])
const newMessage = ref('')
const showProposalInput = ref(false)
const proposedPrice = ref<number | null>(null)
const editingMessage = ref<Message | null>(null)

async function loadConversations() {
    try {
        conversations.value = await chatStore.getConversations()

        const queryListingId = route.query.listingId
        if (queryListingId) {
            const lid = parseInt(queryListingId as string)
            const existing = conversations.value.find((c) => getNumericId(c.listing_id) === lid)
            if (existing) {
                selectConversation(existing)
            } else {
                selectedConv.value = {
                    listing_id: lid,
                    id: 0,
                }
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

function getNumericId(val: any): number | undefined {
    if (val === null || val === undefined) return undefined
    if (typeof val === 'object') return val.Int64
    return val
}

async function send() {
    if (!selectedConv.value || (!newMessage.value && !proposedPrice.value)) return

    try {
        if (editingMessage.value) {
            await chatStore.editMessage(
                editingMessage.value.id,
                newMessage.value,
                proposedPrice.value ?? undefined,
            )
            editingMessage.value = null
        } else {
            const type = proposedPrice.value ? 'price_proposal' : 'text'
            const listingId = getNumericId(selectedConv.value.listing_id)!
            const convId = getNumericId(selectedConv.value.id)
            await chatStore.sendMessage(
                listingId,
                newMessage.value,
                type,
                proposedPrice.value ?? undefined,
                convId,
            )
        }

        newMessage.value = ''
        proposedPrice.value = null
        showProposalInput.value = false

        if (selectedConv.value.id && selectedConv.value.id > 0) {
            messages.value = await chatStore.getMessages(selectedConv.value.id)
        } else {
            await loadConversations()
        }
    } catch (e) {
        console.error(e)
    }
}

function startEdit(msg: Message) {
    editingMessage.value = msg
    newMessage.value = msg.content
    if (msg.message_type === 'price_proposal') {
        showProposalInput.value = true
        proposedPrice.value = msg.proposed_price
    } else {
        showProposalInput.value = false
        proposedPrice.value = null
    }
}

function cancelEdit() {
    editingMessage.value = null
    newMessage.value = ''
    proposedPrice.value = null
    showProposalInput.value = false
}

async function respondToProposal(msgId: number, accept: boolean) {
    try {
        await chatStore.handleProposal(msgId, accept)
        if (selectedConv.value && selectedConv.value.id) {
            messages.value = await chatStore.getMessages(selectedConv.value.id)
        }
    } catch (e) {
        console.error(e)
    }
}

function goToPayment(price: number) {
    if (!selectedConv.value || !selectedConv.value.listing_id) return
    router.push({
        path: '/particulier/paiement',
        query: {
            id: selectedConv.value.listing_id.toString(),
            price: price.toString(),
            type: 'listing',
        },
    })
}

onMounted(loadConversations)

watch(
    () => route.query.listingId,
    () => {
        loadConversations()
    },
)
</script>

<template>
    <div class="chat-container">
        <div class="conv-list">
            <h3>Mes Conversations</h3>
            <div
                v-for="conv in conversations"
                :key="conv.id"
                class="conv-item"
                :class="{ active: selectedConv?.id === conv.id }"
                @click="selectConversation(conv)"
            >
                <p>Listing #{{ getNumericId(conv.listing_id) }}</p>
                <small
                    >Mis à jour :
                    {{ new Date(conv.updated_at.Time ?? conv.updated_at).toLocaleString() }}</small
                >
            </div>
        </div>

        <div class="chat-window" v-if="selectedConv">
            <div class="chat-header">
                <h3>
                    {{ selectedConv.id === 0 ? 'Nouvelle conversation' : 'Conversation' }} - Objet
                    #{{ getNumericId(selectedConv.listing_id) }}
                </h3>
                <span v-if="selectedConv.is_closed" class="status-closed"
                    >Discussion fermée (Objet vendu)</span
                >
            </div>
            <div class="messages">
                <div
                    v-for="msg in messages"
                    :key="msg.id"
                    class="message"
                    :class="{ mine: msg.sender_id === authStore.user?.id }"
                >
                    <div class="bubble">
                        <p v-if="msg.content">{{ msg.content }}</p>

                        <div v-if="msg.message_type === 'price_proposal'" class="proposal-box">
                            <strong>Proposition de prix : {{ msg.proposed_price }}€</strong>
                            <div v-if="msg.proposal_status === 'pending'">
                                <span class="badge pending">En attente</span>
                                <div
                                    v-if="msg.sender_id !== authStore.user?.id && !selectedConv.is_closed"
                                    class="actions"
                                >
                                    <button
                                        @click="respondToProposal(msg.id, true)"
                                        class="btn-accept"
                                    >
                                        Accepter
                                    </button>
                                    <button
                                        @click="respondToProposal(msg.id, false)"
                                        class="btn-decline"
                                    >
                                        Refuser
                                    </button>
                                </div>
                            </div>
                            <div v-else-if="msg.proposal_status === 'accepted'">
                                <span class="badge accepted">Acceptée</span>
                                <div
                                    v-if="msg.sender_id === authStore.user?.id && !selectedConv.is_closed"
                                    class="buy-now"
                                >
                                    <button
                                        @click="goToPayment(msg.proposed_price!)"
                                        class="btn-pay"
                                    >
                                        Payer {{ msg.proposed_price }}€
                                    </button>
                                </div>
                            </div>
                            <span
                                v-else-if="msg.proposal_status === 'declined'"
                                class="badge declined"
                                >Refusée</span
                            >
                        </div>
                        <small class="time">
                            {{
                                new Date(msg.created_at.Time ?? msg.created_at).toLocaleTimeString()
                            }}
                            <span v-if="msg.is_edited" class="edited-label">(modifié)</span>
                        </small>

                        <div
                            v-if="msg.sender_id === authStore.user?.id && !editingMessage && !selectedConv.is_closed"
                            class="msg-actions"
                        >
                            <button @click="startEdit(msg)" class="btn-text">Modifier</button>
                        </div>
                    </div>
                </div>
            </div>

            <div class="input-area" v-if="!selectedConv.is_closed">
                <div v-if="editingMessage" class="edit-banner">
                    Mode édition <button @click="cancelEdit" class="btn-text">Annuler</button>
                </div>
                <textarea
                    v-model="newMessage"
                    :placeholder="editingMessage ? 'Modifier votre message...' : 'Votre message...'"
                ></textarea>
                <div class="controls">
                    <button
                        @click="showProposalInput = !showProposalInput"
                        class="btn-alt"
                        :disabled="editingMessage && editingMessage.message_type === 'text'"
                    >
                        {{ showProposalInput ? 'Annuler proposition' : 'Négocier prix' }}
                    </button>
                    <input
                        v-if="showProposalInput"
                        type="number"
                        v-model="proposedPrice"
                        placeholder="Prix (€)"
                    />
                    <button @click="send" class="btn-send">
                        {{ editingMessage ? 'Mettre à jour' : 'Envoyer' }}
                    </button>
                </div>
            </div>
            <div class="input-area closed-notice" v-else>
                Cette discussion est terminée car l'objet a été vendu.
            </div>
        </div>
        <div class="no-chat" v-else>Sélectionnez une conversation pour commencer.</div>
    </div>
</template>

<style scoped>
.chat-container {
    display: grid;
    grid-template-columns: 320px 1fr;
    height: calc(100vh - 120px);
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    overflow: hidden;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.conv-list {
    border-right: 1.5px solid rgba(53, 53, 53, 0.1);
    display: flex;
    flex-direction: column;
    background: #fafafa;
}

.conv-list h3 {
    padding: 20px;
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    border-bottom: 1px solid rgba(53, 53, 53, 0.05);
}

.conv-item {
    padding: 16px 20px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    cursor: pointer;
    transition: background 0.2s;
    border-bottom: 1px solid rgba(53, 53, 53, 0.03);
}

.conv-item:hover {
    background: rgba(52, 137, 91, 0.05);
}

.conv-item.active {
    background: var(--green-pale);
    border-right: 3px solid var(--green-mid);
}

.conv-item p {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
}

.conv-item small {
    font-size: 0.75rem;
    color: rgba(53, 53, 53, 0.5);
}

.chat-window {
    display: flex;
    flex-direction: column;
    background: var(--white);
}

.chat-header {
    padding: 18px 24px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.status-closed {
    font-size: 0.8rem;
    font-weight: 700;
    color: #c0392b;
    background: #fdf2f2;
    padding: 4px 12px;
    border-radius: 20px;
}

.closed-notice {
    text-align: center;
    padding: 30px;
    color: rgba(53, 53, 53, 0.5);
    font-style: italic;
    font-size: 0.95rem;
    background: #f9f9f9;
}

.chat-header h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--charcoal);
}

.messages {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    background: #fcfcfc;
}

.message {
    display: flex;
    justify-content: flex-start;
}

.message.mine {
    justify-content: flex-end;
}

.bubble {
    max-width: 75%;
    padding: 12px 18px;
    border-radius: 16px;
    background: var(--white);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    border: 1px solid rgba(53, 53, 53, 0.05);
    position: relative;
    color: var(--charcoal);
}

.message.mine .bubble {
    background: var(--green-dark);
    color: var(--white);
    border: none;
}

.bubble p {
    margin: 0;
    font-size: 0.95rem;
    line-height: 1.5;
}

.proposal-box {
    margin-top: 12px;
    padding: 12px;
    border-radius: 8px;
    background: rgba(53, 53, 53, 0.03);
    border-left: 3px solid var(--green-mid);
}

.message.mine .proposal-box {
    background: rgba(255, 255, 255, 0.1);
    border-left-color: var(--green-pale);
}

.badge {
    display: inline-block;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 700;
    margin-top: 8px;
}

.pending {
    background: #ffeeba;
    color: #856404;
}
.accepted {
    background: #d4edda;
    color: #155724;
}
.declined {
    background: #f8d7da;
    color: #721c24;
}

.actions {
    margin-top: 10px;
    display: flex;
    gap: 8px;
}

.btn-accept {
    background: #28a745;
    color: white;
    border: none;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
}

.btn-decline {
    background: #dc3545;
    color: white;
    border: none;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
}

.btn-pay {
    background: #ffc107;
    color: #212529;
    border: none;
    padding: 10px 16px;
    border-radius: 8px;
    font-weight: 700;
    margin-top: 12px;
    cursor: pointer;
    width: 100%;
    transition: background 0.2s;
}

.btn-pay:hover {
    background: #e0a800;
}

.time {
    font-size: 0.7rem;
    opacity: 0.6;
    display: block;
    margin-top: 8px;
}

.message.mine .time {
    opacity: 0.8;
}

.edited-label {
    font-style: italic;
    margin-left: 5px;
}

.msg-actions {
    margin-top: 8px;
    display: flex;
    justify-content: flex-end;
}

.btn-text {
    background: none;
    border: none;
    color: inherit;
    font-size: 0.75rem;
    cursor: pointer;
    text-decoration: underline;
    opacity: 0.7;
}

.input-area {
    padding: 24px;
    border-top: 1px solid rgba(53, 53, 53, 0.08);
}

.edit-banner {
    background: #fff3cd;
    padding: 8px 16px;
    border-radius: 8px;
    font-size: 0.85rem;
    margin-bottom: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

textarea {
    width: 100%;
    height: 80px;
    margin-bottom: 12px;
    padding: 14px;
    border-radius: 10px;
    border: 1.5px solid rgba(53, 53, 53, 0.12);
    font-family: inherit;
    font-size: 0.95rem;
    resize: none;
    outline: none;
    transition: border-color 0.2s;
    box-sizing: border-box;
}

textarea:focus {
    border-color: var(--green-mid);
}

.controls {
    display: flex;
    gap: 12px;
    align-items: center;
}

.price-input {
    width: 100px;
    padding: 8px 12px;
    border-radius: 8px;
    border: 1.5px solid rgba(53, 53, 53, 0.12);
    outline: none;
}

.btn-send {
    background: var(--green-dark);
    color: var(--white);
    border: none;
    padding: 10px 24px;
    border-radius: 8px;
    font-weight: 700;
    cursor: pointer;
    margin-left: auto;
    transition: background 0.2s;
}

.btn-send:hover {
    background: var(--green-mid);
}

.btn-alt {
    background: var(--green-pale);
    color: var(--green-dark);
    border: none;
    padding: 8px 16px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
}

.no-chat {
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fafafa;
    text-align: center;
}

.no-chat p {
    font-size: 1rem;
    color: rgba(53, 53, 53, 0.4);
}
</style>
