<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useChatStore, type Conversation, type Message, type EditHistory } from '@/stores/chat';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n();
const chatStore = useChatStore();

const CENSORED_MARKER = '[Message censuré par la modération]';

const conversations = ref<Conversation[]>([]);
const selectedListingId = ref<string | null>(null);
const selectedConv = ref<Conversation | null>(null);
const messages = ref<Message[]>([]);
const history = ref<EditHistory[]>([]);
const loading = ref(false);
const showAudit = ref(false);

async function loadConversations() {
    loading.value = true;
    try {
        conversations.value = await chatStore.adminGetConversations();
        // If a listing was selected, keep it, otherwise clear selections
        if (selectedListingId.value && !listingsWithConversations.value.some(l => l.group_key === selectedListingId.value)) {
            selectedListingId.value = null;
            selectedConv.value = null;
        }
    } catch (e) {
        console.error('Failed to load conversations:', e);
    } finally {
        loading.value = false;
    }
}

function conversationGroupKey(conv: Conversation): string {
    return conv.course_id ? `course-${conv.course_id}` : `listing-${conv.listing_id}`;
}

// Group conversations by listing or course
const listingsWithConversations = computed(() => {
    const groups: { [key: string]: { group_key: string; listing_title: string; conversations: Conversation[] } } = {};
    for (const conv of conversations.value) {
        const key = conversationGroupKey(conv);
        const title = conv.course_id
            ? (conv.course_title || t('admin.chatReview.trainingFallback', { id: conv.course_id }))
            : (conv.listing_title || t('admin.chatReview.listingFallback', { id: conv.listing_id }));
        if (!groups[key]) {
            groups[key] = {
                group_key: key,
                listing_title: title,
                conversations: []
            };
        }
        groups[key].conversations.push(conv);
    }
    return Object.values(groups);
});

// Selected listing object
const selectedListing = computed(() => {
    return listingsWithConversations.value.find(l => l.group_key === selectedListingId.value) || null;
});

function selectListing(groupKey: string) {
    selectedListingId.value = groupKey;
    const group = listingsWithConversations.value.find(l => l.group_key === groupKey);
    if (group && group.conversations.length > 0) {
        selectConversation(group.conversations[0]);
    } else {
        selectedConv.value = null;
        messages.value = [];
        history.value = [];
    }
}

async function selectConversation(conv: Conversation) {
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

async function censorMessage(msgId: number) {
    if (!confirm(t('admin.chatReview.confirmCensor'))) return;
    loading.value = true;
    try {
        await chatStore.adminCensorMessage(msgId);
        if (selectedConv.value) {
            await selectConversation(selectedConv.value);
        }
    } catch (e) {
        alert(t('admin.chatReview.errorCensor'));
    } finally {
        loading.value = false;
    }
}

onMounted(loadConversations);
</script>

<template>
    <div class="admin-chat">
        <header class="page-header">
            <div>
                <h1 class="page-title">{{ t('admin.chatReview.pageTitle') }}</h1>
                <p class="page-subtitle">{{ t('admin.chatReview.subtitle') }}</p>
            </div>
            <button @click="loadConversations" class="btn-refresh" :disabled="loading">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
                {{ t('admin.chatReview.refresh') }}
            </button>
        </header>

        <div class="two-column-layout">

            <section class="column column--subjects">
                <h3 class="column-title">{{ t('admin.chatReview.subjects') }}</h3>
                <div class="column-content">
                    <div v-if="loading && listingsWithConversations.length === 0" class="loading-state">{{ t('admin.chatReview.loading') }}</div>
                    <div v-else-if="listingsWithConversations.length === 0" class="empty-state">{{ t('admin.chatReview.noSubjects') }}</div>

                    <div v-for="list in listingsWithConversations" :key="list.group_key"
                         class="list-item list-item--subject"
                         :class="{ active: selectedListingId === list.group_key }"
                         @click="selectListing(list.group_key)">
                        <div class="subject-title">{{ list.listing_title }}</div>
                        <div class="subject-meta">{{ t('admin.chatReview.conversationsCount', { count: list.conversations.length }) }}</div>
                    </div>
                </div>
            </section>

            <section class="column column--transcript">
                <h3 class="column-title">{{ t('admin.chatReview.messages') }}</h3>
                <div class="column-content transcript-container">
                    <div v-if="!selectedListingId" class="select-prompt">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="40" height="40">
                            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                        </svg>
                        <p>{{ t('admin.chatReview.selectSubjectPrompt') }}</p>
                    </div>

                    <div v-else class="transcript-wrapper">
                        <!-- Conversations Tabs Bar -->
                        <div class="conv-tabs-bar">
                            <button v-for="conv in selectedListing?.conversations" :key="conv.id"
                                    class="conv-tab" :class="{ active: selectedConv?.id === conv.id }"
                                    @click="selectConversation(conv)">
                                <span class="role-dot buyer"></span> {{ conv.buyer_name }} <span class="tab-vs">vs</span> {{ conv.seller_name }}
                            </button>
                        </div>

                        <div v-if="selectedConv" class="transcript-grid" :class="{ 'has-audit': showAudit }">
                            <!-- Chat message feed -->
                            <div class="chat-feed">
                                <div class="chat-feed-header">
                                    <div>
                                        <h4 class="chat-feed-title">{{ selectedListing?.listing_title }}</h4>
                                        <p class="chat-feed-subtitle">
                                            {{ selectedConv.course_id ? t('admin.chatReview.trainingRef', { id: selectedConv.course_id }) : t('admin.chatReview.listingRef', { id: selectedConv.listing_id }) }} | Chat #{{ selectedConv.id }}
                                        </p>
                                    </div>
                                    <div style="display: flex; align-items: center; gap: 10px;">
                                        <button class="btn-toggle-audit" @click="showAudit = !showAudit">
                                            {{ showAudit ? t('admin.chatReview.hideHistory') : t('admin.chatReview.showHistory', { count: history.length }) }}
                                        </button>
                                    </div>
                                </div>

                                <div class="messages-scroll">
                                    <div v-if="messages.length === 0" class="no-messages">{{ t('admin.chatReview.noMessages') }}</div>
                                    <div v-for="msg in messages" :key="msg.id"
                                         class="message-row"
                                         :class="msg.sender_id === selectedConv.buyer_id ? 'msg-left' : 'msg-right'">

                                        <div class="msg-meta">
                                            <strong>{{ msg.sender_id === selectedConv.buyer_id ? selectedConv.buyer_name : selectedConv.seller_name }}</strong>
                                            <span class="role-label" :class="msg.sender_id === selectedConv.buyer_id ? 'buyer' : 'seller'">
                                                {{ msg.sender_id === selectedConv.buyer_id ? t('admin.chatReview.buyer') : t('admin.chatReview.seller') }}
                                            </span>
                                            <span class="time-label">{{ new Date(msg.created_at.Time ?? msg.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</span>
                                        </div>

                                        <div class="msg-bubble" :class="{
                                            'bubble-buyer': msg.sender_id === selectedConv.buyer_id,
                                            'bubble-seller': msg.sender_id !== selectedConv.buyer_id,
                                            'bubble-censored': msg.content === CENSORED_MARKER
                                        }">
                                            <div class="msg-text">{{ msg.content === CENSORED_MARKER ? t('admin.chatReview.censoredMessage') : msg.content }}</div>

                                            <div v-if="msg.message_type === 'price_proposal'" class="proposal-box">
                                                <div class="proposal-title font-semibold">{{ t('admin.chatReview.priceProposal') }}</div>
                                                <div class="proposal-price">{{ msg.proposed_price }} €</div>
                                                <span class="proposal-badge" :class="msg.proposal_status">
                                                    {{ t('admin.chatReview.statusLabel', { status: msg.proposal_status }) }}
                                                </span>
                                            </div>

                                            <button v-if="msg.content !== CENSORED_MARKER"
                                                    class="btn-bubble-censor"
                                                    @click.stop="censorMessage(msg.id)">
                                                {{ t('admin.chatReview.censorMessage') }}
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div v-if="showAudit" class="audit-sidebar">
                                <h4 class="audit-title">{{ t('admin.chatReview.editHistory') }}</h4>
                                <div class="audit-scroll">
                                    <div v-if="history.length === 0" class="no-audit">{{ t('admin.chatReview.noEdits') }}</div>
                                    <div v-for="h in history" :key="h.id" class="audit-card">
                                        <div class="audit-card-meta">
                                            <span class="msg-id">{{ t('admin.chatReview.msgNumber', { id: h.message_id }) }}</span>
                                            <span class="audit-date">{{ new Date(h.edited_at.Time ?? h.edited_at).toLocaleString(locale === 'en' ? 'en-US' : 'fr-FR', {month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'}) }}</span>
                                        </div>
                                        <div class="audit-card-body">
                                            <div v-if="h.old_content" class="audit-change-item">
                                                <span class="change-label">{{ t('admin.chatReview.oldText') }}</span>
                                                <span class="change-val">"{{ h.old_content }}"</span>
                                            </div>
                                            <div v-if="h.old_proposed_price" class="audit-change-item">
                                                <span class="change-label">{{ t('admin.chatReview.oldPrice') }}</span>
                                                <span class="change-val">{{ h.old_proposed_price }} €</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div v-else class="select-prompt">
                            <p>{{ t('admin.chatReview.noConversationSelected') }}</p>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<style scoped>
.admin-chat {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 120px);
    color: #1e293b;
    font-family: inherit;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-shrink: 0;
}

.page-title {
    font-size: 1.8rem;
    font-weight: 800;
    color: #0f172a;
    letter-spacing: -0.02em;
    margin: 0 0 4px;
}

.page-subtitle {
    font-size: 0.88rem;
    color: #475569;
    margin: 0;
}

.btn-refresh {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: #f1f5f9;
    color: #0f172a;
    border: 1px solid #cbd5e1;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
    transition: all 0.2s;
}

.btn-refresh:hover:not(:disabled) {
    background: #e2e8f0;
    border-color: #94a3b8;
}

.btn-refresh svg {
    transition: transform 0.3s;
}

.btn-refresh:hover svg {
    transform: rotate(180deg);
}

/* 2-Column Navigation Grid */
.two-column-layout {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 20px;
    flex: 1;
    min-height: 0; /* Ensures columns scroll internally */
}

.transcript-wrapper {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
}

.conv-tabs-bar {
    display: flex;
    gap: 8px;
    padding: 12px 16px;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
    overflow-x: auto;
    flex-shrink: 0;
}

.conv-tab {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 8px 14px;
    background: #ffffff;
    border: 1.5px solid #e2e8f0;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    color: #475569;
    cursor: pointer;
    white-space: nowrap;
    transition: all 0.2s;
}

.conv-tab:hover {
    border-color: #cbd5e1;
    background: #f1f5f9;
}

.conv-tab.active {
    border-color: #10b981;
    background: #f0fdf4;
    color: #047857;
}

.tab-vs {
    font-size: 0.7rem;
    opacity: 0.6;
    font-weight: 500;
}

.column {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    box-shadow: 0 1px 3px rgba(0,0,0,0.02);
}

.column-title {
    font-size: 0.82rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #475569;
    padding: 14px 16px;
    margin: 0;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
}

.column-content {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

/* Column 1: Subjects List (Listings) */
.list-item {
    padding: 12px;
    border: 1px solid #f1f5f9;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    background: #f8fafc;
}

.list-item:hover {
    border-color: #cbd5e1;
    background: #f1f5f9;
}

.list-item.active {
    background: #e0f2fe;
    border-color: #0284c7;
}

.subject-title {
    font-weight: 700;
    font-size: 0.88rem;
    color: #0f172a;
    margin-bottom: 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.subject-meta {
    font-size: 0.75rem;
    color: #64748b;
}

/* Column 2: Conversations List */
.conv-header {
    display: flex;
    justify-content: space-between;
    font-size: 0.72rem;
    color: #64748b;
    margin-bottom: 6px;
}

.conv-tag {
    font-weight: 700;
}

.conv-participants {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 0.8rem;
}

.participant {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #334155;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.role-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
}

.role-dot.buyer { background: #3b82f6; }
.role-dot.seller { background: #10b981; }

.list-item--conv.active {
    background: #f0fdf4;
    border-color: #10b981;
}

/* Column 3: Transcript & Details Split */
.transcript-container {
    padding: 0;
    display: flex;
    flex-direction: column;
}

.transcript-grid {
    display: grid;
    grid-template-columns: 1fr;
    height: 100%;
    min-height: 0;
}

.transcript-grid.has-audit {
    grid-template-columns: 1fr 280px;
}

.btn-toggle-audit {
    background: #f1f5f9;
    color: #475569;
    border: 1px solid #cbd5e1;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-toggle-audit:hover {
    background: #e2e8f0;
    color: #0f172a;
    border-color: #94a3b8;
}

/* Chat Transcript area */
.chat-feed {
    border-right: 1px solid #e2e8f0;
    display: flex;
    flex-direction: column;
    min-height: 0;
    background: #fcfcfd;
}

.chat-feed-header {
    padding: 16px 20px;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.chat-feed-title {
    margin: 0 0 2px;
    font-size: 0.95rem;
    font-weight: 800;
    color: #0f172a;
}

.chat-feed-subtitle {
    margin: 0;
    font-size: 0.75rem;
    color: #64748b;
}

.badge-live {
    font-size: 0.68rem;
    font-weight: 700;
    background: #fee2e2;
    color: #ef4444;
    padding: 2px 8px;
    border-radius: 12px;
    text-transform: uppercase;
}

.messages-scroll {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.message-row {
    display: flex;
    flex-direction: column;
    max-width: 80%;
}

.message-row.msg-left {
    align-self: flex-start;
}

.message-row.msg-right {
    align-self: flex-end;
    align-items: flex-end;
}

.msg-meta {
    font-size: 0.72rem;
    color: #64748b;
    margin-bottom: 4px;
    display: flex;
    align-items: center;
    gap: 6px;
}

.role-label {
    font-size: 0.65rem;
    font-weight: 700;
    padding: 1px 6px;
    border-radius: 10px;
    text-transform: uppercase;
}

.role-label.buyer { background: #dbeafe; color: #1e40af; }
.role-label.seller { background: #d1fae5; color: #065f46; }

.msg-bubble {
    position: relative;
    padding: 12px 14px;
    border-radius: 12px;
    font-size: 0.88rem;
    line-height: 1.45;
    color: #0f172a;
    box-shadow: 0 1px 2px rgba(0,0,0,0.03);
}

.bubble-buyer {
    background: #f1f5f9;
    border: 1px solid #e2e8f0;
    border-top-left-radius: 2px;
}

.bubble-seller {
    background: #e8f8f0;
    border: 1px solid #d1ebd9;
    border-top-right-radius: 2px;
}

.bubble-censored {
    background: #f8fafc !important;
    color: #94a3b8 !important;
    border-color: #cbd5e1 !important;
    font-style: italic;
}

.msg-text {
    word-break: break-word;
}

.proposal-box {
    margin-top: 10px;
    padding: 8px 10px;
    background: rgba(255, 255, 255, 0.5);
    border-radius: 6px;
    border: 1px dashed rgba(0, 0, 0, 0.1);
    font-size: 0.8rem;
}

.proposal-price {
    font-size: 1rem;
    font-weight: 800;
    margin: 2px 0;
}

.proposal-badge {
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
}

.proposal-badge.accepted { color: #15803d; }
.proposal-badge.declined { color: #b91c1c; }
.proposal-badge.pending { color: #b45309; }

.btn-bubble-censor {
    margin-top: 8px;
    background: #fee2e2;
    color: #991b1b;
    border: none;
    padding: 4px 10px;
    border-radius: 5px;
    font-size: 0.72rem;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-bubble-censor:hover {
    background: #fca5a5;
}

/* Audit history logs sidebar */
.audit-sidebar {
    background: #fafaf9;
    display: flex;
    flex-direction: column;
    min-height: 0;
    padding: 16px;
}

.audit-title {
    font-size: 0.78rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #78716c;
    margin: 0 0 12px;
}

.audit-scroll {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.audit-card {
    background: #ffffff;
    border: 1px solid #fed7aa;
    background: #fffbeb;
    border-radius: 8px;
    padding: 10px;
    font-size: 0.78rem;
}

.audit-card-meta {
    display: flex;
    justify-content: space-between;
    font-weight: 700;
    color: #c2410c;
    margin-bottom: 4px;
}

.audit-date {
    font-size: 0.7rem;
    font-weight: 500;
    opacity: 0.75;
}

.audit-change-item {
    margin-top: 4px;
    display: flex;
    flex-direction: column;
    gap: 1px;
}

.change-label {
    font-size: 0.7rem;
    opacity: 0.55;
}

.change-val {
    color: #44403c;
    font-style: italic;
}

.no-audit, .no-messages {
    font-size: 0.78rem;
    color: #64748b;
    font-style: italic;
    text-align: center;
    padding: 20px 0;
}

/* Placeholders & Prompts */
.select-prompt {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #94a3b8;
    text-align: center;
    padding: 30px;
}

.select-prompt p {
    margin-top: 12px;
    font-size: 0.88rem;
}

.loading-state, .empty-state {
    padding: 24px 10px;
    text-align: center;
    color: #64748b;
    font-style: italic;
    font-size: 0.85rem;
}
</style>
