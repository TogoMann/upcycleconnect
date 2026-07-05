<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface Formation {
    id: number
    name: string
    description: string
    price: { Float64: number; Valid: boolean } | number | null
    date: { Time: string; Valid: boolean } | string | null
    end_date: { Time: string; Valid: boolean } | string | null
    max_capacity: { Int32: number; Valid: boolean } | number | null
    type: string
    approved: boolean
    status: string
    correction_comment: { String: string; Valid: boolean } | null
}

const formations = ref<Formation[]>([])
const loading = ref(true)
const error = ref('')

async function fetchFormations() {
    loading.value = true
    error.value = ''
    const token = authStore.token
    if (!token) {
        loading.value = false
        return
    }
    try {
        const res = await fetch(`${API_BASE}/salarie/formations`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (!res.ok) throw new Error(t('salarie.mesFormations.errorLoad'))
        formations.value = await res.json()
    } catch (e: any) {
        error.value = e.message || t('salarie.mesFormations.errorLoad')
    } finally {
        loading.value = false
    }
}

onMounted(fetchFormations)

function fmtPrice(p: Formation['price']): string {
    if (p === null) return t('salarie.mesFormations.free')
    const val = typeof p === 'object' ? (p.Valid ? p.Float64 : null) : p
    return val ? `${Number(val).toFixed(2)} €` : t('salarie.mesFormations.free')
}

function rawDateVal(d: Formation['date']): string | null {
    if (!d) return null
    const val = typeof d === 'object' ? (d.Valid ? d.Time : null) : d
    return val || null
}

function fmtDate(d: Formation['date']): string {
    const val = rawDateVal(d)
    if (!val) return '—'
    return new Date(val).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function fmtDateRange(f: Formation): string {
    const start = rawDateVal(f.date)
    const end = rawDateVal(f.end_date)
    if (!start) return '—'
    if (!end || end === start) return fmtDate(f.date)
    return `${fmtDate(f.date)} → ${fmtDate(f.end_date)}`
}

function fmtCapacity(c: Formation['max_capacity']): string {
    if (c === null) return t('salarie.mesFormations.unlimited')
    const val = typeof c === 'object' ? (c.Valid ? c.Int32 : null) : c
    return val ? t('salarie.mesFormations.seats', { count: val }) : t('salarie.mesFormations.unlimited')
}

function statutLabel(f: Formation): string {
    switch (f.status) {
        case 'brouillon': return t('salarie.mesFormations.statusDraft')
        case 'pending': return t('salarie.mesFormations.statusPending')
        case 'approved': return t('salarie.mesFormations.statusApproved')
        case 'needs_modification': return t('salarie.mesFormations.statusNeedsModification')
        case 'rejected': return t('salarie.mesFormations.statusRejected')
        default: return t('salarie.mesFormations.statusDraft')
    }
}

function statutBadgeClass(f: Formation): string {
    switch (f.status) {
        case 'approved': return 'badge--active'
        case 'pending': return 'badge--pending'
        case 'needs_modification': return 'badge--warn'
        case 'rejected': return 'badge--danger'
        default: return 'badge--draft'
    }
}

async function supprimer(id: number) {
    if (!confirm(t('salarie.mesFormations.confirmDelete'))) return
    try {
        await fetch(`${API_BASE}/salarie/formations/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        formations.value = formations.value.filter(f => f.id !== id)
    } catch {
        error.value = t('salarie.mesFormations.errorDelete')
    }
}
</script>

<template>
    <div class="formations">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">{{ t('salarie.mesFormations.pageTitle') }}</h1>
                    <p class="page-subtitle">{{ t('salarie.mesFormations.subtitle') }}</p>
                </div>
                <router-link to="/salarie/formations/nouvelle" class="btn-primary">
                    {{ t('salarie.mesFormations.newTraining') }}
                </router-link>
            </div>
        </div>

        <div v-if="error" class="error-banner">{{ error }}</div>

        <div v-if="loading" class="loading-state">{{ t('salarie.mesFormations.loading') }}</div>

        <div v-else class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('salarie.mesFormations.colTitle') }}</th>
                        <th>{{ t('salarie.mesFormations.colType') }}</th>
                        <th>{{ t('salarie.mesFormations.colDate') }}</th>
                        <th>{{ t('salarie.mesFormations.colPrice') }}</th>
                        <th>{{ t('salarie.mesFormations.colCapacity') }}</th>
                        <th>{{ t('salarie.mesFormations.colStatus') }}</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="formations.length === 0">
                        <td colspan="7" class="empty">{{ t('salarie.mesFormations.empty') }}</td>
                    </tr>
                    <tr v-for="f in formations" :key="f.id">
                        <td class="td-bold">
                            <div>{{ f.name }}</div>
                            <div v-if="f.status === 'needs_modification' && f.correction_comment?.Valid" class="correction-note">
                                {{ t('salarie.mesFormations.correctionRequired', { comment: f.correction_comment.String }) }}
                            </div>
                        </td>
                        <td>
                            <span class="badge" :class="f.type === 'en_ligne' ? 'badge--online' : 'badge--draft'">
                                {{ f.type === 'en_ligne' ? t('salarie.mesFormations.online') : t('salarie.mesFormations.inPerson') }}
                            </span>
                        </td>
                        <td class="td-muted">{{ fmtDateRange(f) }}</td>
                        <td class="td-muted">{{ fmtPrice(f.price) }}</td>
                        <td class="td-muted">{{ fmtCapacity(f.max_capacity) }}</td>
                        <td>
                            <span class="badge" :class="statutBadgeClass(f)">
                                {{ statutLabel(f) }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <router-link :to="`/salarie/formations/${f.id}/edit`" class="btn-icon" :title="t('salarie.mesFormations.edit')">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                </svg>
                            </router-link>
                            <button class="btn-icon btn-icon--danger" :title="t('salarie.mesFormations.delete')" @click="supprimer(f.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="3 6 5 6 21 6" />
                                    <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                                    <path d="M10 11v6M14 11v6" />
                                    <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
                                </svg>
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 28px; }
.header-row { display: flex; justify-content: space-between; align-items: flex-start; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border-radius: 8px; font-size: 0.88rem; font-weight: 600; text-decoration: none; transition: background 0.2s; white-space: nowrap; }
.btn-primary:hover { background: var(--green-mid); }
.error-banner { background: rgba(229, 62, 62, 0.08); border: 1px solid rgba(229, 62, 62, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: #c53030; margin-bottom: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; align-items: center; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.correction-note { font-size: 0.82rem; color: #b91c1c; font-weight: 500; margin-top: 4px; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--draft { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.badge--pending { background: #fef9c3; color: #a16207; }
.badge--warn { background: #fef3c7; color: #d97706; }
.badge--danger { background: #fee2e2; color: #b91c1c; }
.badge--online { background: #e0f2fe; color: #0369a1; }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); text-decoration: none; transition: border-color 0.2s, color 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
</style>
