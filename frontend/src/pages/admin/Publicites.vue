<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()

interface Pub {
    id: number
    titre: string
    annonceur: string
    type: string
    debut: string
    fin: string
    statut: string
    budget: number
}

const pubs = ref<Pub[]>([])
const filter = ref('all')
const loading = ref(true)

const filteredPubs = computed(() => {
    if (filter.value === 'all') return pubs.value
    return pubs.value.filter(p => p.statut === filter.value)
})

const stats = computed(() => ({
    total: pubs.value.length,
    active: pubs.value.filter(p => p.statut === 'active').length,
    pending: pubs.value.filter(p => p.statut === 'inactive').length,
    revenue: pubs.value.filter(p => p.statut === 'active').reduce((s, p) => s + p.budget, 0),
}))

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/admin/publicites`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) pubs.value = await res.json()
    } catch {}
    loading.value = false
})

async function approuver(p: Pub) {
    try {
        const res = await fetch(`${API_BASE}/advertisement/${p.id}/approve`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) p.statut = 'active'
    } catch {}
}

async function rejeter(p: Pub) {
    try {
        const res = await fetch(`${API_BASE}/advertisement/${p.id}/reject`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) p.statut = 'inactive'
    } catch {}
}

async function supprimer(id: number) {
    if (!confirm(t('admin.publicites.confirmDelete'))) return
    try {
        const res = await fetch(`${API_BASE}/admin/publicites/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) pubs.value = pubs.value.filter(p => p.id !== id)
    } catch {}
}

function filterLabel(f: string): string {
    if (f === 'all') return t('admin.publicites.filterAll')
    if (f === 'active') return t('admin.publicites.filterActive')
    return t('admin.publicites.filterPendingRejected')
}
</script>

<template>
    <div class="publicites">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.publicites.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.publicites.subtitle') }}</p>
        </div>

        <div v-if="loading" class="loading-state">{{ t('admin.publicites.loading') }}</div>

        <template v-else>
            <div class="kpi-row">
                <div class="kpi-sm">
                    <div class="kpi-sm-value">{{ stats.total }}</div>
                    <div class="kpi-sm-label">{{ t('admin.publicites.kpiTotal') }}</div>
                </div>
                <div class="kpi-sm kpi-sm--green">
                    <div class="kpi-sm-value">{{ stats.active }}</div>
                    <div class="kpi-sm-label">{{ t('admin.publicites.kpiActive') }}</div>
                </div>
                <div class="kpi-sm kpi-sm--yellow">
                    <div class="kpi-sm-value">{{ stats.pending }}</div>
                    <div class="kpi-sm-label">{{ t('admin.publicites.kpiPending') }}</div>
                </div>
                <div class="kpi-sm">
                    <div class="kpi-sm-value">{{ stats.revenue.toFixed(0) }}€</div>
                    <div class="kpi-sm-label">{{ t('admin.publicites.kpiRevenue') }}</div>
                </div>
            </div>

            <div class="filter-row">
                <button v-for="f in ['all', 'active', 'inactive']" :key="f"
                    class="filter-btn" :class="{ 'filter-btn--active': filter === f }"
                    @click="filter = f">
                    {{ filterLabel(f) }}
                </button>
            </div>

            <div class="table-wrap">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>{{ t('admin.publicites.colTitle') }}</th>
                            <th>{{ t('admin.publicites.colAdvertiser') }}</th>
                            <th>{{ t('admin.publicites.colType') }}</th>
                            <th>{{ t('admin.publicites.colPeriod') }}</th>
                            <th>{{ t('admin.publicites.colBudget') }}</th>
                            <th>{{ t('admin.publicites.colStatus') }}</th>
                            <th>{{ t('admin.publicites.colActions') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="filteredPubs.length === 0">
                            <td colspan="7" class="empty">{{ t('admin.publicites.empty') }}</td>
                        </tr>
                        <tr v-for="p in filteredPubs" :key="p.id">
                            <td class="td-bold">{{ p.titre }}</td>
                            <td class="td-muted">{{ p.annonceur }}</td>
                            <td>{{ p.type }}</td>
                            <td class="td-muted">{{ p.debut }} → {{ p.fin }}</td>
                            <td class="td-budget">{{ p.budget.toFixed(0) }} €</td>
                            <td>
                                <span class="badge" :class="p.statut === 'active' ? 'badge--active' : 'badge--inactive'">
                                    {{ p.statut === 'active' ? t('admin.publicites.active') : t('admin.publicites.inactive') }}
                                </span>
                            </td>
                            <td>
                                <div class="action-group">
                                    <button v-if="p.statut !== 'active'" class="btn-icon btn-icon--approve" @click="approuver(p)" :title="t('admin.publicites.approve')">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                                    </button>
                                    <button v-if="p.statut === 'active'" class="btn-icon btn-icon--reject" @click="rejeter(p)" :title="t('admin.publicites.deactivate')">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
                                    </button>
                                    <button class="btn-icon btn-icon--danger" @click="supprimer(p.id)" :title="t('admin.publicites.delete')">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/></svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </template>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; }

.kpi-row { display: flex; gap: 14px; margin-bottom: 20px; }
.kpi-sm { background: var(--white); border: 1.5px solid rgba(53,53,53,0.08); border-radius: 12px; padding: 16px 24px; }
.kpi-sm--green { background: var(--green-pale); border-color: transparent; }
.kpi-sm--yellow { background: #fffbeb; border-color: transparent; }
.kpi-sm-value { font-size: 1.8rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; line-height: 1; }
.kpi-sm--green .kpi-sm-value { color: var(--green-dark); }
.kpi-sm--yellow .kpi-sm-value { color: #92400e; }
.kpi-sm-label { font-size: 0.78rem; font-weight: 500; color: var(--charcoal); opacity: 0.5; margin-top: 4px; }

.filter-row { display: flex; gap: 8px; margin-bottom: 20px; }
.filter-btn { padding: 7px 16px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 20px; background: transparent; color: var(--charcoal); font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.filter-btn:hover { border-color: var(--green-mid); }
.filter-btn--active { background: var(--green-dark); color: var(--white); border-color: var(--green-dark); }

.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-budget { font-weight: 700; color: var(--green-dark); }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--inactive { background: #fef3c7; color: #92400e; }

.action-group { display: flex; gap: 6px; }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; transition: all 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon--approve { color: var(--green-dark); border-color: var(--green-dark); }
.btn-icon--approve:hover { background: var(--green-pale); }
.btn-icon--reject { color: #d97706; border-color: #fde68a; }
.btn-icon--reject:hover { background: #fffbeb; }
.btn-icon--danger { color: var(--charcoal); opacity: 0.5; }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; opacity: 1; }
</style>
