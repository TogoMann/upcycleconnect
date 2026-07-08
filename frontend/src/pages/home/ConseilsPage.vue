<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { API_BASE } from '@/config'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface Conseil {
    id: number
    title: string
    content: string
    created_at: string
    upvotes: number
    downvotes: number
    my_vote: 'up' | 'down' | null
}

const conseils = ref<Conseil[]>([])
const loading = ref(true)
const sortBy = ref<'recent' | 'top' | 'bottom'>('recent')

async function fetchConseils() {
    loading.value = true
    try {
        const headers: Record<string, string> = {}
        if (authStore.token) headers.Authorization = `Bearer ${authStore.token}`
        const res = await fetch(`${API_BASE}/conseils?sort=${sortBy.value}`, { headers })
        if (res.ok) conseils.value = await res.json()
    } catch {
        conseils.value = []
    } finally {
        loading.value = false
    }
}

onMounted(fetchConseils)
watch(sortBy, fetchConseils)

async function vote(article: Conseil, type: 'up' | 'down') {
    if (!authStore.isAuthenticated) return
    try {
        const res = await fetch(`${API_BASE}/conseils/${article.id}/vote`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({ type }),
        })
        if (res.ok) {
            const updated: Conseil = await res.json()
            article.upvotes = updated.upvotes
            article.downvotes = updated.downvotes
            article.my_vote = updated.my_vote
        }
    } catch {}
}

function fmtDate(iso: string): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

function excerpt(content: string): string {
    return content.length > 220 ? content.slice(0, 220) + '…' : content
}
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">{{ t('conseils.pageTitle') }}</h1>
                <p class="hero-subtitle">
                    {{ t('conseils.subtitle') }}
                </p>
            </div>
        </section>

        <section class="articles-section">
            <div class="container">
                <div class="sort-row">
                    <label class="sort-label" for="conseils-sort">{{ t('conseils.sortLabel') }}</label>
                    <select id="conseils-sort" v-model="sortBy" class="sort-select">
                        <option value="recent">{{ t('conseils.sortRecent') }}</option>
                        <option value="top">{{ t('conseils.sortMostLiked') }}</option>
                        <option value="bottom">{{ t('conseils.sortLeastLiked') }}</option>
                    </select>
                </div>

                <div v-if="loading" class="empty-state">
                    <p>{{ t('conseils.loading') }}</p>
                </div>

                <div v-else class="articles-grid">
                    <article v-for="article in conseils" :key="article.id" class="article-card">
                        <div class="article-body">
                            <h2 class="article-titre">{{ article.title }}</h2>
                            <p class="article-extrait">{{ excerpt(article.content) }}</p>
                            <div class="article-meta">
                                <span class="article-date">{{ fmtDate(article.created_at) }}</span>
                                <div class="vote-row">
                                    <button
                                        class="vote-btn vote-btn--up"
                                        :class="{ 'vote-btn--active': article.my_vote === 'up', 'vote-btn--disabled': !authStore.isAuthenticated }"
                                        :title="authStore.isAuthenticated ? t('conseils.likeAction') : t('conseils.loginToVote')"
                                        @click="vote(article, 'up')"
                                    >
                                        ▲ {{ article.upvotes }}
                                    </button>
                                    <button
                                        class="vote-btn vote-btn--down"
                                        :class="{ 'vote-btn--active': article.my_vote === 'down', 'vote-btn--disabled': !authStore.isAuthenticated }"
                                        :title="authStore.isAuthenticated ? t('conseils.dislikeAction') : t('conseils.loginToVote')"
                                        @click="vote(article, 'down')"
                                    >
                                        ▼ {{ article.downvotes }}
                                    </button>
                                </div>
                            </div>
                        </div>
                    </article>
                </div>

                <div v-if="!loading && conseils.length === 0" class="empty-state">
                    <p>{{ t('conseils.noResults') }}</p>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page-content {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}

.hero {
    padding: 64px 0 48px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.6rem, 5.5vw, 4.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0 0 16px;
}
.hero-subtitle {
    font-size: 1rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
}

.articles-section {
    flex: 1;
    padding: 0 0 80px;
}

.articles-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
}

.article-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 10px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition:
        transform 0.2s,
        box-shadow 0.2s;
}
.article-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(8, 106, 53, 0.1);
}

.article-body {
    padding: 22px 22px 24px;
    display: flex;
    flex-direction: column;
    flex: 1;
}

.article-titre {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.4;
    margin: 0 0 10px;
}

.article-extrait {
    font-size: 0.85rem;
    color: var(--charcoal);
    opacity: 0.72;
    line-height: 1.6;
    margin: 0 0 16px;
    flex: 1;
}

.article-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.6;
}

.sort-row {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 10px;
    margin-bottom: 24px;
}
.sort-label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--charcoal);
    opacity: 0.75;
}
.sort-select {
    padding: 8px 14px;
    font-size: 0.85rem;
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    border-radius: 8px;
    background: var(--white);
    color: var(--charcoal);
    font-family: inherit;
    outline: none;
    cursor: pointer;
}

.vote-row { display: flex; gap: 6px; }
.vote-btn { background: none; border: 1px solid rgba(53,53,53,0.15); border-radius: 6px; padding: 3px 8px; font-size: 0.75rem; font-weight: 600; cursor: pointer; transition: all 0.2s; display: flex; align-items: center; gap: 4px; opacity: 1; }
.vote-btn--up { color: var(--green-dark); }
.vote-btn--up:hover:not(.vote-btn--disabled) { background: rgba(8,106,53,0.08); border-color: var(--green-mid); }
.vote-btn--up.vote-btn--active { background: var(--green-pale); border-color: var(--green-dark); }
.vote-btn--down { color: #dc2626; }
.vote-btn--down:hover:not(.vote-btn--disabled) { background: rgba(220,38,38,0.08); border-color: #fca5a5; }
.vote-btn--down.vote-btn--active { background: #fee2e2; border-color: #dc2626; }
.vote-btn--disabled { opacity: 0.45; cursor: default; }

.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
}

@media (max-width: 860px) {
    .articles-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
@media (max-width: 560px) {
    .articles-grid {
        grid-template-columns: 1fr;
    }
    .hero-title {
        font-size: 2.2rem;
    }
}
</style>
