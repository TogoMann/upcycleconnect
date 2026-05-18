<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useChatStore, type Conversation, type Message, type EditHistory } from '@/stores/chat';

const chatStore = useChatStore();

const conversations = ref<Conversation[]>([]);
const selectedConv = ref<Conversation | null>(null);
const messages = ref<Message[]>([]);
const history = ref<EditHistory[]>([]);
const loading = ref(false);

async function loadConversations() {
    loading.value = true;
    try {
        conversations.value = await chatStore.adminGetConversations();
    } catch (e) {
        console.error('Failed to load conversations:', e);
    } finally {
        loading.value = false;
    }
}

async function viewConversation(conv: Conversation) {
    selectedConv.value = conv;
    loading.value = true;
    try {
        const res = await chatStore.adminGetConversationDetails(conv.id);
        messages.value = res.messages || [];
        history.value = res.history || [];
    } catch (e) {
        console.error('Failed to load conversation details:', e);
        messages.value = [];
        history.value = [];
    } finally {
        loading.value = false;
    }
}

function getMessageHistory(msgId: number) {
    return history.value.filter(h => h.message_id === msgId);
}

onMounted(loadConversations);
</script>

<template>
    <div class="admin-chat">
        <header class="page-header">
            <h1>Modération des Conversations</h1>
            <button @click="loadConversations" class="btn-refresh" :disabled="loading">Actualiser la liste</button>
        </header>

        <div class="admin-layout-grid">
            <!-- Sidebar: Conversation List -->
            <aside class="conversation-sidebar">
                <div v-if="loading && conversations.length === 0" class="status-msg">Chargement...</div>
                <div v-else-if="conversations.length === 0" class="status-msg">Aucune conversation trouvée.</div>
                
                <div v-for="conv in conversations" :key="conv.id" 
                     class="conv-card" :class="{ active: selectedConv?.id === conv.id }"
                     @click="viewConversation(conv)">
                    <div class="conv-card-header">
                        <span class="conv-id">#{{ conv.id }}</span>
                        <span class="conv-date">{{ new Date(conv.updated_at.Time ?? conv.updated_at).toLocaleDateString() }}</span>
                    </div>
                    <div class="conv-listing">{{ conv.listing_title }}</div>
                    <div class="conv-users">
                        <span>{{ conv.buyer_name }}</span>
                        <span class="arrow">→</span>
                        <span>{{ conv.seller_name }}</span>
                    </div>
                </div>
            </aside>

            <!-- Main: Transcript Viewer -->
            <main class="transcript-viewer">
                <div v-if="selectedConv" class="viewer-content">
                    <div class="viewer-header">
                        <h2>Discussion : {{ selectedConv.listing_title }}</h2>
                        <div class="viewer-meta">
                            <span>Acheteur : <strong>{{ selectedConv.buyer_name }}</strong></span>
                            <span>Vendeur : <strong>{{ selectedConv.seller_name }}</strong></span>
                        </div>
                    </div>

                    <div v-if="loading" class="viewer-loading">Chargement des messages...</div>
                    <div v-else-if="messages.length === 0" class="viewer-empty">Cette conversation ne contient aucun message.</div>
                    
                    <div v-else class="transcript-list">
                        <div v-for="msg in messages" :key="msg.id" class="transcript-item">
                            <div class="item-meta">
                                <span class="sender-type" :class="msg.sender_id === selectedConv.buyer_id ? 'buyer' : 'seller'">
                                    {{ msg.sender_id === selectedConv.buyer_id ? 'ACHETEUR' : 'VENDEUR' }}
                                </span>
                                <span class="msg-time">{{ new Date(msg.created_at.Time ?? msg.created_at).toLocaleString() }}</span>
                                <span v-if="msg.is_edited" class="edited-badge">Modifié</span>
                            </div>
                            
                            <div class="item-content">
                                <div class="message-text" v-if="msg.content">{{ msg.content }}</div>
                                
                                <div v-if="msg.message_type === 'price_proposal'" class="proposal-info">
                                    <div class="proposal-label">Proposition de prix</div>
                                    <div class="proposal-value">{{ msg.proposed_price }}€</div>
                                    <div class="proposal-status" :class="msg.proposal_status">Statut : {{ msg.proposal_status }}</div>
                                </div>
                            </div>

                            <!-- History Log -->
                            <div v-if="getMessageHistory(msg.id).length > 0" class="message-history">
                                <div class="history-header">Historique des modifications :</div>
                                <div v-for="h in getMessageHistory(msg.id)" :key="h.id" class="history-record">
                                    <div class="record-time">Le {{ new Date(h.edited_at.Time ?? h.edited_at).toLocaleString() }}</div>
                                    <div class="record-content">
                                        <div v-if="h.old_content">Ancien texte : <em>{{ h.old_content }}</em></div>
                                        <div v-if="h.old_proposed_price">Ancien prix : <em>{{ h.old_proposed_price }}€</em></div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div v-else class="viewer-placeholder">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="48" height="48">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                    </svg>
                    <p>Sélectionnez une conversation pour examiner les messages.</p>
                </div>
            </main>
        </div>
    </div>
</template>

<style scoped>
.admin-chat {
    display: flex;
    flex-direction: column;
    height: 100%;
    color: #2d3748;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.btn-refresh {
    background: #4a5568;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.9em;
}

.admin-layout-grid {
    display: grid;
    grid-template-columns: 320px 1fr;
    gap: 24px;
    flex: 1;
    min-height: 0; /* Important for inner scrolling */
}

/* Sidebar */
.conversation-sidebar {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    overflow-y: auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.conv-card {
    padding: 16px;
    border: 1px solid #edf2f7;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s;
}

.conv-card:hover {
    border-color: #cbd5e0;
    background: #f8fafc;
}

.conv-card.active {
    border-color: #3182ce;
    background: #ebf8ff;
    box-shadow: 0 2px 4px rgba(66, 153, 225, 0.1);
}

.conv-card-header {
    display: flex;
    justify-content: space-between;
    font-size: 0.75rem;
    color: #718096;
    margin-bottom: 8px;
}

.conv-listing {
    font-weight: 700;
    font-size: 0.95rem;
    color: #2d3748;
    margin-bottom: 8px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.conv-users {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
    color: #4a5568;
}

.arrow { color: #a0aec0; }

/* Main Viewer */
.transcript-viewer {
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.viewer-header {
    padding: 24px;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
}

.viewer-header h2 {
    margin: 0 0 12px;
    font-size: 1.25rem;
}

.viewer-meta {
    display: flex;
    gap: 24px;
    font-size: 0.9rem;
    color: #718096;
}

.viewer-meta strong { color: #2d3748; }

.transcript-list {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.transcript-item {
    border-left: 4px solid #e2e8f0;
    padding-left: 20px;
    position: relative;
}

.item-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
}

.sender-type {
    font-size: 0.7rem;
    font-weight: 800;
    padding: 2px 8px;
    border-radius: 4px;
    text-transform: uppercase;
}

.sender-type.buyer { background: #bee3f8; color: #2c5282; }
.sender-type.seller { background: #c6f6d5; color: #22543d; }

.msg-time { font-size: 0.75rem; color: #a0aec0; }

.edited-badge {
    font-size: 0.7rem;
    font-style: italic;
    color: #b7791f;
}

.item-content {
    background: #f7fafc;
    padding: 12px 16px;
    border-radius: 8px;
    color: #2d3748;
    line-height: 1.5;
}

.proposal-info {
    margin-top: 12px;
    padding: 12px;
    background: #edf2f7;
    border-radius: 6px;
    border: 1px dashed #cbd5e0;
}

.proposal-label { font-size: 0.8rem; color: #718096; margin-bottom: 4px; }
.proposal-value { font-size: 1.2rem; font-weight: 800; color: #2d3748; }
.proposal-status { font-size: 0.8rem; margin-top: 4px; font-weight: 600; }
.proposal-status.accepted { color: #38a169; }
.proposal-status.declined { color: #e53e3e; }
.proposal-status.pending { color: #d69e2e; }

/* History */
.message-history {
    margin-top: 12px;
    background: #fffaf0;
    border: 1px solid #feebc8;
    border-radius: 6px;
    padding: 12px;
}

.history-header {
    font-size: 0.75rem;
    font-weight: 800;
    color: #9c4221;
    margin-bottom: 8px;
    text-transform: uppercase;
}

.history-record {
    font-size: 0.85rem;
    padding-top: 8px;
    border-top: 1px solid #fbd38d;
    margin-top: 8px;
}

.record-time { font-size: 0.75rem; color: #718096; margin-bottom: 4px; }

/* Placeholders */
.viewer-placeholder {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #a0aec0;
    text-align: center;
    padding: 40px;
}

.viewer-placeholder p { margin-top: 16px; font-size: 1.1rem; }

.status-msg, .viewer-loading, .viewer-empty {
    padding: 40px;
    text-align: center;
    color: #718096;
    font-style: italic;
}
</style>
