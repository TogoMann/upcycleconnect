<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

interface Thread {
    id: number
    created_by: number | null
    title: string
    content: string
    upvotes: number
    downvotes: number
    views: number
    created_at: string | null
    last_post_at: string | null
}

const threads = ref<Thread[]>([])
const searchQuery = ref('')
const loading = ref(true)

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/thread`)
        if (res.ok) threads.value = await res.json()
    } catch {}
    loading.value = false
})

const threadsFiltres = computed(() =>
    threads.value.filter(
        (th) =>
            !searchQuery.value || th.title.toLowerCase().includes(searchQuery.value.toLowerCase()),
    ),
)

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    })
}

function timeAgo(iso: string | null): string {
    if (!iso) return '—'
    const now = new Date()
    const date = new Date(iso)

    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const target = new Date(date.getFullYear(), date.getMonth(), date.getDate())

    const diffTime = today.getTime() - target.getTime()
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

    if (diffDays === 0) return t('forum.today')
    if (diffDays === 1) return t('forum.yesterday')
    return t('forum.daysAgo', { days: diffDays })
}

function initials(title: string): string {
    return title
        .split(' ')
        .slice(0, 2)
        .map((w) => w[0]?.toUpperCase() ?? '')
        .join('')
}
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container hero-inner">
                <div class="hero-text">
                    <h1 class="hero-title">{{ t('forum.pageTitle') }}</h1>
                    <p class="hero-subtitle">
                        {{ t('forum.subtitle') }}
                    </p>
                </div>
                <router-link to="/forum/nouveau" class="btn-nouveau">{{ t('forum.newThread') }}</router-link>
            </div>
        </section>

        <section class="forum-section">
            <div class="container">
                <div class="forum-toolbar">
                    <div class="search-bar">
                        <svg
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <circle cx="11" cy="11" r="8" />
                            <line x1="21" y1="21" x2="16.65" y2="16.65" />
                        </svg>
                        <input
                            v-model="searchQuery"
                            type="text"
                            :placeholder="t('forum.searchPlaceholder')"
                            class="search-input"
                        />
                    </div>
                </div>

                <div v-if="loading" class="loading">{{ t('forum.loading') }}</div>

                <div v-else class="forum-table">
                    <div class="forum-header">
                        <span class="col-sujet">{{ t('forum.colSubject') }}</span>
                        <span class="col-stats">{{ t('forum.colVotes') }}</span>
                        <span class="col-stats">{{ t('forum.colViews') }}</span>
                        <span class="col-activite">{{ t('forum.colActivity') }}</span>
                    </div>

                    <router-link
                        v-for="th in threadsFiltres"
                        :key="th.id"
                        :to="`/forum/${th.id}`"
                        class="forum-row"
                    >
                        <div class="col-sujet">
                            <div class="sujet-avatar">{{ initials(th.title) }}</div>
                            <div class="sujet-info">
                                <h3 class="sujet-titre">{{ th.title }}</h3>
                                <div class="sujet-meta">
                                    <span class="sujet-date">{{ fmtDate(th.created_at) }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="col-stats">
                            <span class="stat-value">{{ th.upvotes }}</span>
                        </div>
                        <div class="col-stats">
                            <span class="stat-value">{{ th.views }}</span>
                        </div>
                        <div class="col-activite">
                            <span class="activite-text">{{ timeAgo(th.last_post_at) }}</span>
                        </div>
                    </router-link>

                    <div v-if="threadsFiltres.length === 0" class="empty-state">
                        <p>{{ t('forum.noResults') }}</p>
                    </div>
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
    padding: 56px 0 40px;
}
.hero-inner {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    flex-wrap: wrap;
}
.hero-title {
    font-size: clamp(2.4rem, 5vw, 3.8rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0 0 10px;
}
.hero-subtitle {
    font-size: 0.95rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
}
.btn-nouveau {
    background: var(--green-dark);
    color: var(--white);
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    white-space: nowrap;
    transition: background 0.2s;
    flex-shrink: 0;
}
.btn-nouveau:hover {
    background: var(--green-mid);
}
.forum-section {
    flex: 1;
    padding: 0 0 80px;
}
.forum-toolbar {
    margin-bottom: 24px;
}
.search-bar {
    display: flex;
    align-items: center;
    gap: 10px;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    padding: 0 16px;
    transition: border-color 0.2s;
    color: var(--charcoal);
    opacity: 0.6;
}
.search-bar:focus-within {
    border-color: var(--green-mid);
    opacity: 1;
}
.search-input {
    flex: 1;
    padding: 12px 0;
    font-size: 0.9rem;
    color: var(--charcoal);
    background: transparent;
    border: none;
    outline: none;
    font-family: inherit;
}
.search-input::placeholder {
    color: rgba(53, 53, 53, 0.45);
}
.loading {
    opacity: 0.5;
    font-size: 0.9rem;
    padding: 40px 0;
}
.forum-table {
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 10px;
    overflow: hidden;
}
.forum-header {
    display: grid;
    grid-template-columns: 1fr 80px 80px 120px;
    gap: 16px;
    padding: 12px 20px;
    background: var(--green-pale);
    font-size: 0.78rem;
    font-weight: 700;
    color: var(--charcoal);
    letter-spacing: 0.04em;
    text-transform: uppercase;
}
.forum-row {
    display: grid;
    grid-template-columns: 1fr 80px 80px 120px;
    gap: 16px;
    padding: 18px 20px;
    text-decoration: none;
    color: inherit;
    border-top: 1px solid rgba(53, 53, 53, 0.08);
    transition: background 0.15s;
    align-items: center;
}
.forum-row:hover {
    background: rgba(215, 236, 225, 0.4);
}
.col-sujet {
    display: flex;
    align-items: flex-start;
    gap: 14px;
    min-width: 0;
}
.col-stats {
    text-align: center;
}
.col-activite {
    text-align: right;
}
.sujet-avatar {
    flex-shrink: 0;
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: var(--green-pale);
    color: var(--green-dark);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.02em;
}
.sujet-info {
    min-width: 0;
}
.sujet-titre {
    font-size: 0.92rem;
    font-weight: 600;
    color: var(--charcoal);
    margin: 0 0 6px;
    line-height: 1.4;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.forum-row:hover .sujet-titre {
    color: var(--green-dark);
}
.sujet-meta {
    display: flex;
    align-items: center;
    gap: 8px;
}
.sujet-date {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.5;
}
.stat-value {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
}
.activite-text {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.6;
}
.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
}
@media (max-width: 700px) {
    .forum-header {
        grid-template-columns: 1fr;
    }
    .forum-header .col-stats,
    .forum-header .col-activite {
        display: none;
    }
    .forum-row {
        grid-template-columns: 1fr;
    }
    .forum-row .col-stats,
    .forum-row .col-activite {
        display: none;
    }
    .sujet-titre {
        white-space: normal;
    }
}
</style>
