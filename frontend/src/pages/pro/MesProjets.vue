<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface Project {
    id: { Int64: number; Valid: boolean }
    listing_id: { Int64: number; Valid: boolean }
    creator_id: { Int64: number; Valid: boolean }
    title: string
    description: string
    final_score: { Int32: number; Valid: boolean }
    status: string
    created_at: { Time: string; Valid: boolean }
    completed_at: { Time: string; Valid: boolean }
}

const projects = ref<Project[]>([])
const loading = ref(true)
const filter = ref('all')

const filteredProjects = computed(() => {
    if (filter.value === 'all') return projects.value
    return projects.value.filter(p => p.status === filter.value)
})

const stats = computed(() => ({
    total: projects.value.length,
    inProgress: projects.value.filter(p => p.status === 'in progress').length,
    done: projects.value.filter(p => p.status === 'done').length,
    featured: projects.value.filter(p => p.status === 'featured').length,
}))

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch(`${API_BASE}/project/me`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) {
            const data = await res.json()
            projects.value = Array.isArray(data) ? data : []
        }
    } catch {}
    loading.value = false
})

function formatDate(d: { Time: string; Valid: boolean } | undefined) {
    if (!d?.Valid) return '—'
    return new Date(d.Time).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function statusConfig(s: string) {
    const map: Record<string, { label: string; class: string; icon: string }> = {
        'in progress': { label: t('pro.mesProjets.statusInProgress'), class: 'badge--progress', icon: 'progress' },
        'done': { label: t('pro.mesProjets.statusDone'), class: 'badge--done', icon: 'check' },
        'featured': { label: t('pro.mesProjets.statusFeatured'), class: 'badge--featured', icon: 'star' },
        'cancelled': { label: t('pro.mesProjets.statusCancelled'), class: 'badge--cancelled', icon: 'x' },
    }
    return map[s] || { label: s, class: 'badge--default', icon: 'default' }
}

function truncate(text: string, max: number) {
    if (!text) return t('pro.mesProjets.noDescription')
    return text.length > max ? text.slice(0, max) + '...' : text
}
</script>

<template>
    <div class="projets">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">{{ t('pro.mesProjets.pageTitle') }}</h1>
                    <p class="page-subtitle">{{ t('pro.mesProjets.subtitle') }}</p>
                </div>
                <router-link to="/pro/projets/nouveau" class="btn-primary">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                    {{ t('pro.mesProjets.newProject') }}
                </router-link>
            </div>
        </div>

        <div v-if="loading" class="loading-state">{{ t('pro.mesProjets.loading') }}</div>

        <template v-else>
            
            <div class="kpi-row">
                <div class="kpi-card">
                    <div class="kpi-value">{{ stats.total }}</div>
                    <div class="kpi-label">{{ t('pro.mesProjets.total') }}</div>
                </div>
                <div class="kpi-card kpi-card--blue">
                    <div class="kpi-value">{{ stats.inProgress }}</div>
                    <div class="kpi-label">{{ t('pro.mesProjets.inProgress') }}</div>
                </div>
                <div class="kpi-card kpi-card--green">
                    <div class="kpi-value">{{ stats.done }}</div>
                    <div class="kpi-label">{{ t('pro.mesProjets.done') }}</div>
                </div>
                <div class="kpi-card kpi-card--yellow">
                    <div class="kpi-value">{{ stats.featured }}</div>
                    <div class="kpi-label">{{ t('pro.mesProjets.featured') }}</div>
                </div>
            </div>

            
            <div class="filter-row">
                <button v-for="f in ['all', 'in progress', 'done', 'featured', 'cancelled']" :key="f"
                    class="filter-btn" :class="{ 'filter-btn--active': filter === f }"
                    @click="filter = f">
                    {{ f === 'all' ? t('pro.mesProjets.all') : statusConfig(f).label }}
                </button>
            </div>

            
            <div v-if="filteredProjects.length === 0" class="empty-state">
                <div class="empty-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                    </svg>
                </div>
                <p class="empty-text">{{ filter === 'all' ? t('pro.mesProjets.noProjects') : t('pro.mesProjets.noProjectsWithStatus') }}</p>
                <router-link v-if="filter === 'all'" to="/pro/projets/nouveau" class="btn-primary">{{ t('pro.mesProjets.createProject') }}</router-link>
            </div>

            
            <div v-else class="projets-grid">
                <router-link
                    v-for="p in filteredProjects"
                    :key="p.id.Int64"
                    :to="`/pro/projets/${p.id.Int64}`"
                    class="projet-card"
                >
                    <div class="projet-header">
                        <span class="projet-title">{{ p.title }}</span>
                        <span class="badge" :class="statusConfig(p.status).class">{{ statusConfig(p.status).label }}</span>
                    </div>
                    <p class="projet-desc">{{ truncate(p.description, 120) }}</p>
                    <div class="projet-meta">
                        <div class="meta-item">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                            {{ formatDate(p.created_at) }}
                        </div>
                        <div v-if="p.final_score?.Valid" class="meta-item meta-score">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                            {{ t('pro.mesProjets.score', { score: p.final_score.Int32 }) }}
                        </div>
                    </div>
                    <div v-if="p.status === 'featured'" class="featured-badge-bar">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                        {{ t('pro.mesProjets.featuredBanner') }}
                    </div>
                </router-link>
            </div>
        </template>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.header-row { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { display: inline-flex; align-items: center; gap: 8px; padding: 11px 20px; background: var(--green-dark); color: var(--white); border-radius: 8px; font-size: 0.88rem; font-weight: 600; text-decoration: none; transition: background 0.2s; white-space: nowrap; border: none; cursor: pointer; }
.btn-primary:hover { background: var(--green-mid); }
.btn-primary svg { width: 16px; height: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }


.kpi-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 20px; }
.kpi-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.08); border-radius: 12px; padding: 18px 20px; }
.kpi-card--green { background: var(--green-pale); border-color: transparent; }
.kpi-card--blue { background: #eff6ff; border-color: transparent; }
.kpi-card--yellow { background: #fffbeb; border-color: transparent; }
.kpi-value { font-size: 1.6rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; line-height: 1; }
.kpi-card--green .kpi-value { color: var(--green-dark); }
.kpi-card--blue .kpi-value { color: #1e40af; }
.kpi-card--yellow .kpi-value { color: #92400e; }
.kpi-label { font-size: 0.78rem; font-weight: 500; color: var(--charcoal); opacity: 0.5; margin-top: 4px; }


.filter-row { display: flex; gap: 8px; margin-bottom: 20px; flex-wrap: wrap; }
.filter-btn { padding: 7px 16px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 20px; background: transparent; color: var(--charcoal); font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.filter-btn:hover { border-color: var(--green-mid); color: var(--green-dark); }
.filter-btn--active { background: var(--green-dark); color: var(--white); border-color: var(--green-dark); }


.empty-state { display: flex; flex-direction: column; align-items: center; gap: 16px; padding: 60px 0; }
.empty-icon { width: 64px; height: 64px; background: var(--green-pale); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: var(--green-mid); }
.empty-icon svg { width: 32px; height: 32px; }
.empty-text { font-size: 0.95rem; color: var(--charcoal); opacity: 0.5; margin: 0; }


.projets-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; }
.projet-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 14px; padding: 22px; text-decoration: none; display: flex; flex-direction: column; gap: 12px; transition: border-color 0.2s, transform 0.2s; }
.projet-card:hover { border-color: var(--green-mid); transform: translateY(-2px); }
.projet-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 10px; }
.projet-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); line-height: 1.3; }
.projet-desc { font-size: 0.84rem; color: var(--charcoal); opacity: 0.55; line-height: 1.5; margin: 0; flex: 1; }
.projet-meta { display: flex; gap: 16px; align-items: center; }
.meta-item { display: flex; align-items: center; gap: 6px; font-size: 0.8rem; color: var(--charcoal); opacity: 0.5; }
.meta-item svg { width: 14px; height: 14px; }
.meta-score { color: var(--green-dark); opacity: 1; font-weight: 600; }
.meta-score svg { color: #f59e0b; }
.featured-badge-bar { display: flex; align-items: center; gap: 8px; padding: 8px 12px; background: #fffbeb; border-radius: 8px; font-size: 0.78rem; font-weight: 600; color: #92400e; }
.featured-badge-bar svg { width: 14px; height: 14px; color: #f59e0b; }


.badge { display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 0.73rem; font-weight: 600; white-space: nowrap; }
.badge--progress { background: #eff6ff; color: #1e40af; }
.badge--done { background: var(--green-pale); color: var(--green-dark); }
.badge--featured { background: #fef3c7; color: #92400e; }
.badge--cancelled { background: rgba(53,53,53,0.08); color: var(--charcoal); opacity: 0.6; }
.badge--default { background: rgba(53,53,53,0.08); color: var(--charcoal); }

@media (max-width: 800px) {
    .kpi-row { grid-template-columns: repeat(2, 1fr); }
    .projets-grid { grid-template-columns: 1fr; }
}
</style>
