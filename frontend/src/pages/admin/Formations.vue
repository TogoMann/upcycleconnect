<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface AdminCourse {
    id: number
    name: string
    description: string
    creator_name: string
    max_capacity: { Int32: number; Valid: boolean } | number | null
    approved: boolean
    status: string
    correction_comment: { String: string; Valid: boolean } | null
    price: { Float64: number; Valid: boolean } | number | null
    date: { Time: string; Valid: boolean } | string | null
    end_date: { Time: string; Valid: boolean } | string | null
    type: string
    session_link: string
    session_count: number
}

const courses = ref<AdminCourse[]>([])
const loading = ref(true)
const error = ref('')
const filterStatus = ref('')

const showProposeModal = ref(false)
const proposeComment = ref('')
const proposeId = ref<number | null>(null)
const actionLoading = ref(false)

async function fetchCourses() {
    loading.value = true
    error.value = ''
    try {
        const res = await fetch(`${API_BASE}/admin/formations`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (!res.ok) throw new Error(t('admin.formations.errorLoad'))
        courses.value = await res.json()
    } catch (e: any) {
        error.value = e.message || t('admin.formations.errorLoad')
    } finally {
        loading.value = false
    }
}

onMounted(fetchCourses)

function rawVal(v: any): any {
    if (v === null || v === undefined) return null
    if (typeof v === 'object') return v.Valid ? (v.Time ?? v.Float64 ?? v.Int32 ?? v.String) : null
    return v
}

function fmtDate(d: AdminCourse['date']): string {
    const val = rawVal(d)
    if (!val) return '—'
    return new Date(val).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function fmtDateRange(c: AdminCourse): string {
    const start = rawVal(c.date)
    const end = rawVal(c.end_date)
    if (!start) return '—'
    if (!end || end === start) return fmtDate(c.date)
    return `${fmtDate(c.date)} → ${fmtDate(c.end_date)}`
}

function fmtPrice(p: AdminCourse['price']): string {
    const val = rawVal(p)
    return val ? `${Number(val).toFixed(2)} €` : t('admin.formations.free')
}

function fmtCapacity(c: AdminCourse['max_capacity']): string {
    const val = rawVal(c)
    return val ? t('admin.formations.seats', { count: val }) : t('admin.formations.unlimited')
}

function statutLabel(c: AdminCourse): string {
    switch (c.status) {
        case 'brouillon': return t('admin.formations.statusDraft')
        case 'pending': return t('admin.formations.statusPending')
        case 'approved': return t('admin.formations.statusApproved')
        case 'needs_modification': return t('admin.formations.statusNeedsModification')
        case 'rejected': return t('admin.formations.statusRejected')
        default: return t('admin.formations.statusDraft')
    }
}

function statutBadgeClass(c: AdminCourse): string {
    switch (c.status) {
        case 'approved': return 'badge--active'
        case 'pending': return 'badge--pending'
        case 'needs_modification': return 'badge--warn'
        case 'rejected': return 'badge--danger'
        default: return 'badge--draft'
    }
}

const displayedCourses = computed(() =>
    filterStatus.value ? courses.value.filter(c => c.status === filterStatus.value) : courses.value
)

function setFilter(status: string) {
    filterStatus.value = filterStatus.value === status ? '' : status
}

async function approuver(id: number) {
    actionLoading.value = true
    try {
        const res = await fetch(`${API_BASE}/admin/catalogue/${id}/approve`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) await fetchCourses()
    } catch {}
    actionLoading.value = false
}

async function refuser(id: number) {
    if (!confirm(t('admin.formations.confirmReject'))) return
    actionLoading.value = true
    try {
        const res = await fetch(`${API_BASE}/admin/catalogue/${id}/disapprove`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) await fetchCourses()
    } catch {}
    actionLoading.value = false
}

function openPropose(id: number) {
    proposeId.value = id
    proposeComment.value = ''
    showProposeModal.value = true
}

async function submitPropose() {
    if (!proposeComment.value) return
    actionLoading.value = true
    try {
        const res = await fetch(`${API_BASE}/admin/catalogue/${proposeId.value}/propose`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({ comment: proposeComment.value }),
        })
        if (res.ok) {
            showProposeModal.value = false
            await fetchCourses()
        }
    } catch {}
    actionLoading.value = false
}
</script>

<template>
    <div class="formations-admin">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.formations.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.formations.subtitle') }}</p>
        </div>

        <div class="filter-chips">
            <button class="filter-chip" :class="{ active: filterStatus === '' }" @click="setFilter('')">{{ t('admin.formations.filterAll') }}</button>
            <button class="filter-chip" :class="{ active: filterStatus === 'pending' }" @click="setFilter('pending')">{{ t('admin.formations.filterPending') }}</button>
            <button class="filter-chip" :class="{ active: filterStatus === 'needs_modification' }" @click="setFilter('needs_modification')">{{ t('admin.formations.filterNeedsModification') }}</button>
            <button class="filter-chip" :class="{ active: filterStatus === 'approved' }" @click="setFilter('approved')">{{ t('admin.formations.filterApproved') }}</button>
            <button class="filter-chip" :class="{ active: filterStatus === 'brouillon' }" @click="setFilter('brouillon')">{{ t('admin.formations.filterDraft') }}</button>
            <button class="filter-chip" :class="{ active: filterStatus === 'rejected' }" @click="setFilter('rejected')">{{ t('admin.formations.filterRejected') }}</button>
        </div>

        <div v-if="error" class="error-banner">{{ error }}</div>
        <div v-if="loading" class="loading-state">{{ t('admin.formations.loading') }}</div>

        <div v-else class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.formations.colTitle') }}</th>
                        <th>{{ t('admin.formations.colTrainer') }}</th>
                        <th>{{ t('admin.formations.colType') }}</th>
                        <th>{{ t('admin.formations.colDates') }}</th>
                        <th>{{ t('admin.formations.colPrice') }}</th>
                        <th>{{ t('admin.formations.colCapacity') }}</th>
                        <th>{{ t('admin.formations.colStatus') }}</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="displayedCourses.length === 0">
                        <td colspan="8" class="empty">{{ t('admin.formations.empty') }}</td>
                    </tr>
                    <tr v-for="c in displayedCourses" :key="c.id">
                        <td class="td-bold">
                            <div>{{ c.name }}</div>
                            <div v-if="c.status === 'needs_modification' && c.correction_comment?.Valid" class="correction-note">
                                {{ t('admin.formations.correctionRequested', { comment: c.correction_comment.String }) }}
                            </div>
                            <div v-if="c.session_count > 1" class="session-note">{{ t('admin.formations.sessionDays', { count: c.session_count }) }}</div>
                        </td>
                        <td class="td-muted">{{ c.creator_name }}</td>
                        <td>
                            <span class="badge" :class="c.type === 'en_ligne' ? 'badge--online' : 'badge--draft'">
                                {{ c.type === 'en_ligne' ? t('admin.formations.online') : t('admin.formations.inPerson') }}
                            </span>
                        </td>
                        <td class="td-muted">{{ fmtDateRange(c) }}</td>
                        <td class="td-muted">{{ fmtPrice(c.price) }}</td>
                        <td class="td-muted">{{ fmtCapacity(c.max_capacity) }}</td>
                        <td>
                            <span class="badge" :class="statutBadgeClass(c)">{{ statutLabel(c) }}</span>
                        </td>
                        <td class="td-actions">
                            <button v-if="c.status !== 'approved'" class="btn-icon" :title="t('admin.formations.approve')" :disabled="actionLoading" @click="approuver(c.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                            </button>
                            <button v-if="c.status !== 'rejected' && c.status !== 'brouillon'" class="btn-icon" :title="t('admin.formations.proposeModification')" :disabled="actionLoading" @click="openPropose(c.id)">
                                <span style="font-size: 1.1rem; line-height: 1;">📝</span>
                            </button>
                            <button v-if="c.status !== 'rejected'" class="btn-icon btn-icon--danger" :title="t('admin.formations.reject')" :disabled="actionLoading" @click="refuser(c.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <line x1="18" y1="6" x2="6" y2="18"></line>
                                    <line x1="6" y1="6" x2="18" y2="18"></line>
                                </svg>
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div v-if="showProposeModal" class="form-overlay" @click.self="showProposeModal = false">
            <div class="form-modal">
                <h3 class="modal-title">{{ t('admin.formations.proposeModalTitle') }}</h3>
                <div class="form-group">
                    <label class="form-label">{{ t('admin.formations.commentLabel') }}</label>
                    <textarea v-model="proposeComment" class="form-input form-textarea" rows="4" :placeholder="t('admin.formations.commentPlaceholder')"></textarea>
                </div>
                <div class="modal-actions">
                    <button class="btn-secondary" @click="showProposeModal = false">{{ t('admin.formations.cancel') }}</button>
                    <button class="btn-primary" :disabled="actionLoading || !proposeComment" @click="submitPropose">
                        {{ actionLoading ? t('admin.formations.sending') : t('admin.formations.sendProposal') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 20px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.filter-chips { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 20px; }
.filter-chip { padding: 6px 14px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 20px; background: none; font-size: 0.8rem; font-weight: 600; cursor: pointer; transition: all 0.15s; font-family: inherit; color: var(--charcoal); }
.filter-chip:hover { border-color: var(--green-mid); color: var(--green-dark); }
.filter-chip.active { background: var(--green-dark); color: white; border-color: var(--green-dark); }
.error-banner { background: rgba(229, 62, 62, 0.08); border: 1px solid rgba(229, 62, 62, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: #c53030; margin-bottom: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); white-space: nowrap; }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; white-space: nowrap; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.correction-note { font-size: 0.8rem; color: #b91c1c; font-weight: 500; margin-top: 4px; }
.session-note { font-size: 0.78rem; color: var(--green-dark); font-weight: 600; margin-top: 4px; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; white-space: nowrap; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--draft { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.badge--pending { background: #fef9c3; color: #a16207; }
.badge--warn { background: #fef3c7; color: #d97706; }
.badge--danger { background: #fee2e2; color: #b91c1c; }
.badge--online { background: #e0f2fe; color: #0369a1; }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover:not(:disabled) { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon:disabled { opacity: 0.4; cursor: default; }
.btn-icon--danger:hover:not(:disabled) { border-color: #dc2626; color: #dc2626; }
.form-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.35); z-index: 200; display: flex; align-items: center; justify-content: center; padding: 20px; }
.form-modal { background: var(--white); border-radius: 16px; padding: 32px; width: 100%; max-width: 480px; display: flex; flex-direction: column; gap: 18px; }
.modal-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 80px; }
.modal-actions { display: flex; gap: 12px; justify-content: flex-end; }
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; }
.btn-secondary { padding: 11px 20px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; }
.btn-secondary:hover { border-color: var(--charcoal); }
</style>
