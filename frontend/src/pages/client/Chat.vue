<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useChatStore, type Conversation, type Message } from '@/stores/chat';
import { useAuthStore } from '@/stores/auth';

const chatStore = useChatStore();
const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const conversations = ref<Conversation[]>([]);
const selectedConv = ref<Partial<Conversation> | null>(null);
const messages = ref<Message[]>([]);
const newMessage = ref('');
const showProposalInput = ref(false);
const proposedPrice = ref<number | null>(null);
const editingMessage = ref<Message | null>(null);

async function loadConversations() {
    try {
        conversations.value = await chatStore.getConversations();
        
        // Check if we come from a listing
        const queryListingId = route.query.listingId;
        if (queryListingId) {
            const lid = parseInt(queryListingId as string);
            const existing = conversations.value.find(c => {
                const cLid = typeof c.listing_id === 'object' ? (c.listing_id as any).Int64 : c.listing_id;
                return cLid === lid;
            });
            if (existing) {
                selectConversation(existing);
            } else {
                // Prepare a "virtual" conversation
                selectedConv.value = {
                    listing_id: lid,
                    id: 0 // Indicates new
                };
                messages.value = [];
            }
        }
    } catch (e) {
        console.error(e);
    }
}

async function selectConversation(conv: Partial<Conversation>) {
    selectedConv.value = conv;
    if (conv.id && conv.id > 0) {
        try {
            messages.value = await chatStore.getMessages(conv.id);
        } catch (e) {
            console.error(e);
        }
    } else {
        messages.value = [];
    }
}

async function send() {
    if (!selectedConv.value || (!newMessage.value && !proposedPrice.value)) return;
    
    try {
        if (editingMessage.value) {
            await chatStore.editMessage(editingMessage.value.id, newMessage.value, proposedPrice.value ?? undefined);
            editingMessage.value = null;
        } else {
            const type = proposedPrice.value ? 'price_proposal' : 'text';
            const listingId = selectedConv.value.listing_id!;
            await chatStore.sendMessage(listingId, newMessage.value, type, proposedPrice.value ?? undefined);
        }
        
        newMessage.value = '';
        proposedPrice.value = null;
        showProposalInput.value = false;
        
        // Reload messages or conversations
        if (selectedConv.value.id && selectedConv.value.id > 0) {
            messages.value = await chatStore.getMessages(selectedConv.value.id);
        } else {
            await loadConversations();
        }
        
    } catch (e) {
        console.error(e);
    }
}

function startEdit(msg: Message) {
    editingMessage.value = msg;
    newMessage.value = msg.content;
    if (msg.message_type === 'price_proposal') {
        showProposalInput.value = true;
        proposedPrice.value = msg.proposed_price;
    } else {
        showProposalInput.value = false;
        proposedPrice.value = null;
    }
}

function cancelEdit() {
    editingMessage.value = null;
    newMessage.value = '';
    proposedPrice.value = null;
    showProposalInput.value = false;
}

async function respondToProposal(msgId: number, accept: boolean) {
    try {
        await chatStore.handleProposal(msgId, accept);
        if (selectedConv.value && selectedConv.value.id) {
            messages.value = await chatStore.getMessages(selectedConv.value.id);
        }
    } catch (e) {
        console.error(e);
    }
}

function goToPayment(price: number) {
    if (!selectedConv.value || !selectedConv.value.listing_id) return;
    router.push({
        path: '/particulier/paiement',
        query: {
            id: selectedConv.value.listing_id.toString(),
            price: price.toString(),
            type: 'listing'
        }
    });
}

onMounted(loadConversations);

// Watch for query changes if user clicks another listing while on chat page
watch(() => route.query.listingId, () => {
    loadConversations();
});
</script>

<template>
    <div class="chat-container">
        <div class="conv-list">
            <h3>Mes Conversations</h3>
            <div v-for="conv in conversations" :key="conv.id" 
                 class="conv-item" :class="{ active: selectedConv?.id === conv.id }"
                 @click="selectConversation(conv)">
                <p>Listing #{{ conv.listing_id }}</p>
                <small>Mis à jour : {{ new Date(conv.updated_at.Time ?? conv.updated_at).toLocaleString() }}</small>
            </div>
        </div>

        <div class="chat-window" v-if="selectedConv">
            <div class="chat-header">
                <h3>{{ selectedConv.id === 0 ? 'Nouvelle conversation' : 'Conversation' }} - Objet #{{ selectedConv.listing_id }}</h3>
            </div>
            <div class="messages">
                <div v-for="msg in messages" :key="msg.id" 
                     class="message" :class="{ mine: msg.sender_id === authStore.user?.id }">
                    <div class="bubble">
                        <p v-if="msg.content">{{ msg.content }}</p>
                        
                        <div v-if="msg.message_type === 'price_proposal'" class="proposal-box">
                            <strong>Proposition de prix : {{ msg.proposed_price }}€</strong>
                            <div v-if="msg.proposal_status === 'pending'">
                                <span class="badge pending">En attente</span>
                                <div v-if="msg.sender_id !== authStore.user?.id" class="actions">
                                    <button @click="respondToProposal(msg.id, true)" class="btn-accept">Accepter</button>
                                    <button @click="respondToProposal(msg.id, false)" class="btn-decline">Refuser</button>
                                </div>
                            </div>
                            <div v-else-if="msg.proposal_status === 'accepted'">
                                <span class="badge accepted">Acceptée</span>
                                <div v-if="msg.sender_id === authStore.user?.id" class="buy-now">
                                    <button @click="goToPayment(msg.proposed_price!)" class="btn-pay">Payer {{ msg.proposed_price }}€</button>
                                </div>
                            </div>
                            <span v-else-if="msg.proposal_status === 'declined'" class="badge declined">Refusée</span>
                        </div>
                        <small class="time">
                            {{ new Date(msg.created_at.Time ?? msg.created_at).toLocaleTimeString() }}
                            <span v-if="msg.is_edited" class="edited-label">(modifié)</span>
                        </small>
                        
                        <div v-if="msg.sender_id === authStore.user?.id && !editingMessage" class="msg-actions">
                            <button @click="startEdit(msg)" class="btn-text">Modifier</button>
                        </div>
                    </div>
                </div>
            </div>

            <div class="input-area">
                <div v-if="editingMessage" class="edit-banner">
                    Mode édition <button @click="cancelEdit" class="btn-text">Annuler</button>
                </div>
                <textarea v-model="newMessage" :placeholder="editingMessage ? 'Modifier votre message...' : 'Votre message...'"></textarea>
                <div class="controls">
                    <button @click="showProposalInput = !showProposalInput" class="btn-alt" :disabled="editingMessage && editingMessage.message_type === 'text'">
                        {{ showProposalInput ? 'Annuler proposition' : 'Négocier prix' }}
                    </button>
                    <input v-if="showProposalInput" type="number" v-model="proposedPrice" placeholder="Prix (€)">
                    <button @click="send" class="btn-send">{{ editingMessage ? 'Mettre à jour' : 'Envoyer' }}</button>
                </div>
            </div>
        </div>
        <div class="no-chat" v-else>
            Sélectionnez une conversation pour commencer.
        </div>
    </div>
</template>

<style scoped>
.chat-container {
    display: grid;
    grid-template-columns: 300px 1fr;
    height: 80vh;
    border: 1px solid #ddd;
    border-radius: 8px;
    background: white;
}

.conv-list {
    border-right: 1px solid #ddd;
    overflow-y: auto;
    padding: 10px;
}

.conv-item {
    padding: 10px;
    border-bottom: 1px solid #eee;
    cursor: pointer;
}

.conv-item.active {
    background: #f0f7ff;
}

.chat-window {
    display: flex;
    flex-direction: column;
}

.chat-header {
    padding: 15px 20px;
    border-bottom: 1px solid #ddd;
    background: #fdfdfd;
}

.messages {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.message {
    display: flex;
    justify-content: flex-start;
}

.message.mine {
    justify-content: flex-end;
}

.bubble {
    max-width: 70%;
    padding: 10px 15px;
    border-radius: 18px;
    background: #f1f0f0;
    position: relative;
}

.message.mine .bubble {
    background: #0084ff;
    color: white;
}

.proposal-box {
    margin-top: 10px;
    padding: 10px;
    border: 1px solid rgba(0,0,0,0.1);
    border-radius: 8px;
    background: rgba(0,0,0,0.05);
}

.badge {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 0.8em;
    margin-top: 5px;
}

.pending { background: #ffeeba; color: #856404; }
.accepted { background: #d4edda; color: #155724; }
.declined { background: #f8d7da; color: #721c24; }

.input-area {
    padding: 20px;
    border-top: 1px solid #ddd;
}

textarea {
    width: 100%;
    height: 60px;
    margin-bottom: 10px;
    padding: 10px;
    border-radius: 8px;
    border: 1px solid #ddd;
    resize: none;
}

.controls {
    display: flex;
    gap: 10px;
    align-items: center;
}

.btn-send {
    background: #0084ff;
    color: white;
    border: none;
    padding: 8px 20px;
    border-radius: 20px;
    cursor: pointer;
    margin-left: auto;
}

.btn-alt {
    background: #eee;
    border: none;
    padding: 5px 15px;
    border-radius: 15px;
    cursor: pointer;
}

.btn-accept { background: #28a745; color: white; border: none; padding: 5px 10px; border-radius: 5px; margin-right: 5px; cursor: pointer; }
.btn-decline { background: #dc3545; color: white; border: none; padding: 5px 10px; border-radius: 5px; cursor: pointer; }

.btn-pay {
    background: #ffc107;
    color: #212529;
    border: none;
    padding: 8px 15px;
    border-radius: 5px;
    font-weight: bold;
    margin-top: 10px;
    cursor: pointer;
    width: 100%;
}
.btn-pay:hover {
    background: #e0a800;
}

.time {
    font-size: 0.7em;
    opacity: 0.7;
    display: block;
    margin-top: 5px;
}

.edited-label {
    font-style: italic;
    margin-left: 5px;
}

.msg-actions {
    margin-top: 5px;
    text-align: right;
}

.btn-text {
    background: none;
    border: none;
    color: inherit;
    font-size: 0.75em;
    cursor: pointer;
    text-decoration: underline;
    opacity: 0.8;
}

.edit-banner {
    background: #fff3cd;
    padding: 5px 10px;
    border-radius: 5px;
    font-size: 0.85em;
    margin-bottom: 10px;
    display: flex;
    justify-content: space-between;
}
</style>
