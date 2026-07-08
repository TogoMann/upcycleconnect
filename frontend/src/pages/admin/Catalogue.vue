<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, reactive, watch, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface Offre {
    id: number
    nom: string
    categorie: string
    prix: number
    description: string
    actif: boolean
    type?: string
}

interface SessionRow {
    date: string
    start_time: string
    end_time: string
}

const offres = ref<Offre[]>([])
const showForm = ref(false)
const editId = ref<number | null>(null)
const loading = ref(false)
const formError = ref('')

const form = reactive({
    nom: '',
    categorie: '',
    prix: '',
    description: '',
    actif: true,
    type: 'presentiel' as 'presentiel' | 'en_ligne',
    session_link: '',
    max_capacity: '',
    multiDay: false,
    startDate: '',
    endDate: '',
})

const sessions = ref<SessionRow[]>([])
const sessionsReady = ref(false)

const showProposeModal = ref(false)
const proposeComment = ref('')
const proposeId = ref<number | null>(null)

function getTodayDate(): string {
    const d = new Date()
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
}

const todayDate = getTodayDate()

function dateRange(start: string, end: string): string[] {
    if (!start) return []
    const endDate = end && end >= start ? end : start
    const dates: string[] = []
    let cursor = new Date(start + 'T00:00:00')
    const last = new Date(endDate + 'T00:00:00')
    while (cursor <= last) {
        const y = cursor.getFullYear()
        const m = String(cursor.getMonth() + 1).padStart(2, '0')
        const d = String(cursor.getDate()).padStart(2, '0')
        dates.push(`${y}-${m}-${d}`)
        cursor.setDate(cursor.getDate() + 1)
    }
    return dates
}

watch(() => [form.startDate, form.endDate, form.multiDay], () => {
    if (!sessionsReady.value) return
    const effectiveEnd = form.multiDay ? form.endDate : form.startDate
    const dates = dateRange(form.startDate, effectiveEnd)
    const existingByDate = new Map(sessions.value.map(s => [s.date, s]))
    sessions.value = dates.map(date => existingByDate.get(date) || { date, start_time: '', end_time: '' })
})

function minTimeForDate(date: string): string {
    if (date !== todayDate) return ''
    const now = new Date()
    return `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
}

function sessionDuration(row: SessionRow): string | null {
    if (!row.start_time || !row.end_time) return null
    const [sh, sm] = row.start_time.split(':').map(Number)
    const [eh, em] = row.end_time.split(':').map(Number)
    const totalMinutes = (eh * 60 + em) - (sh * 60 + sm)
    if (totalMinutes <= 0) return null
    const h = Math.floor(totalMinutes / 60)
    const m = totalMinutes % 60
    if (h === 0) return `${m} min`
    if (m === 0) return `${h} h`
    return `${h} h ${m} min`
}

function applyFirstRowToAll() {
    if (sessions.value.length === 0) return
    const { start_time, end_time } = sessions.value[0]
    sessions.value = sessions.value.map(row => ({ ...row, start_time, end_time }))
}

function rawTimeStr(time: any): string {
    if (!time) return ''
    if (typeof time === 'string') return time.substring(0, 5)
    if (typeof time === 'object' && time.Valid && typeof time.Microseconds === 'number') {
        const totalSeconds = Math.floor(time.Microseconds / 1_000_000)
        const h = Math.floor(totalSeconds / 3600)
        const m = Math.floor((totalSeconds % 3600) / 60)
        return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
    }
    return ''
}

function rawDateStr(d: any): string {
    if (!d) return ''
    const val = typeof d === 'object' ? (d.Valid ? d.Time : null) : d
    return val ? val.substring(0, 10) : ''
}

function openPropose(id: number) {
    proposeId.value = id
    proposeComment.value = ''
    showProposeModal.value = true
}

async function submitPropose() {
    if (!proposeComment.value) return
    loading.value = true
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
            alert(t('admin.catalogue.proposalSent'))
            showProposeModal.value = false
            await fetchOffres()
        }
    } catch {}
    loading.value = false
}

async function fetchOffres() {
    try {
        const res = await fetch(`${API_BASE}/admin/catalogue`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) offres.value = await res.json()
    } catch {}
}

onMounted(fetchOffres)

function resetForm() {
    form.nom = ''
    form.categorie = ''
    form.prix = ''
    form.description = ''
    form.actif = true
    form.type = 'presentiel'
    form.session_link = ''
    form.max_capacity = ''
    form.multiDay = false
    form.startDate = ''
    form.endDate = ''
    sessions.value = []
    sessionsReady.value = false
    formError.value = ''
}

function openCreate() {
    editId.value = null
    resetForm()
    sessionsReady.value = true
    showForm.value = true
}

async function openEdit(o: Offre) {
    editId.value = o.id
    resetForm()
    showForm.value = true
    loading.value = true
    try {
        const res = await fetch(`${API_BASE}/course/${o.id}`)
        const data = res.ok ? await res.json() : {}

        form.nom = o.nom
        form.categorie = o.categorie
        form.prix = String(o.prix)
        form.description = o.description
        form.actif = o.actif
        form.type = data.type === 'en_ligne' ? 'en_ligne' : 'presentiel'
        form.session_link = data.session_link || ''
        form.max_capacity = data.max_capacity?.Int32 ?? data.max_capacity ?? ''

        let loadedSessions: SessionRow[] = []
        try {
            const sessRes = await fetch(`${API_BASE}/course/${o.id}/sessions`)
            if (sessRes.ok) {
                const raw = await sessRes.json()
                loadedSessions = (raw || []).map((s: any) => ({
                    date: rawDateStr(s.session_date),
                    start_time: rawTimeStr(s.start_time),
                    end_time: rawTimeStr(s.end_time),
                }))
            }
        } catch {}

        if (loadedSessions.length === 0 && data.date) {
            loadedSessions = [{
                date: rawDateStr(data.date),
                start_time: rawTimeStr(data.start_time),
                end_time: rawTimeStr(data.end_time),
            }]
        }

        loadedSessions.sort((a, b) => a.date.localeCompare(b.date))
        sessions.value = loadedSessions
        form.startDate = loadedSessions[0]?.date || ''
        form.multiDay = loadedSessions.length > 1
        form.endDate = form.multiDay ? loadedSessions[loadedSessions.length - 1].date : ''
    } catch {}
    sessionsReady.value = true
    loading.value = false
}

async function save() {
    formError.value = ''

    if (sessions.value.length === 0) {
        formError.value = t('admin.catalogue.errorSelectDate')
        return
    }
    for (const row of sessions.value) {
        if (!row.start_time || !row.end_time) {
            formError.value = t('admin.catalogue.errorScheduleRequired', { date: row.date })
            return
        }
        if (row.start_time >= row.end_time) {
            formError.value = t('admin.catalogue.errorEndAfterStart', { date: row.date })
            return
        }
        if (row.date === todayDate && row.start_time < minTimeForDate(row.date)) {
            formError.value = t('admin.catalogue.errorPastTime', { date: row.date })
            return
        }
    }

    loading.value = true
    const body = {
        nom: form.nom,
        categorie: form.categorie,
        prix: parseFloat(form.prix) || 0,
        description: form.description,
        actif: form.actif,
        type: form.type,
        session_link: form.session_link,
        max_capacity: parseInt(form.max_capacity) || null,
        sessions: sessions.value,
    }
    try {
        if (editId.value) {
            const res = await fetch(`${API_BASE}/admin/catalogue/${editId.value}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
                body: JSON.stringify(body),
            })
            if (res.ok) {
                await fetchOffres()
            } else {
                const d = await res.json().catch(() => null)
                formError.value = d?.message || t('admin.catalogue.errorUpdate')
                loading.value = false
                return
            }
        } else {
            const res = await fetch(`${API_BASE}/admin/catalogue`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
                body: JSON.stringify(body),
            })
            if (res.ok) {
                await fetchOffres()
            } else {
                const d = await res.json().catch(() => null)
                formError.value = d?.message || t('admin.catalogue.errorCreate')
                loading.value = false
                return
            }
        }
    } catch {
        formError.value = t('admin.catalogue.errorNetwork')
        loading.value = false
        return
    }
    showForm.value = false
    loading.value = false
}

async function approuver(id: number) {
    try {
        const res = await fetch(`${API_BASE}/admin/catalogue/${id}/approve`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const idx = offres.value.findIndex(o => o.id === id)
            if (idx !== -1) offres.value[idx].actif = true
        }
    } catch {}
}

async function desapprouver(id: number) {
    try {
        const res = await fetch(`${API_BASE}/admin/catalogue/${id}/disapprove`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const idx = offres.value.findIndex(o => o.id === id)
            if (idx !== -1) offres.value[idx].actif = false
        }
    } catch {}
}

async function supprimer(id: number) {
    if (!confirm(t('admin.catalogue.confirmDelete'))) return
    await fetch(`${API_BASE}/admin/catalogue/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    offres.value = offres.value.filter(o => o.id !== id)
}
</script>

<template>
    <div class="catalogue">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">{{ t('admin.catalogue.pageTitle') }}</h1>
                    <p class="page-subtitle">{{ t('admin.catalogue.subtitle') }}</p>
                </div>
                <button class="btn-primary" @click="openCreate">{{ t('admin.catalogue.add') }}</button>
            </div>
        </div>

        <div v-if="showForm" class="form-overlay" @click.self="showForm = false">
            <div class="form-modal">
                <h3 class="modal-title">{{ editId ? t('admin.catalogue.editOffer') : t('admin.catalogue.newOffer') }}</h3>

                <div v-if="formError" class="alert alert--error">{{ formError }}</div>

                <div class="form-group">
                    <label class="form-label">{{ t('admin.catalogue.name') }}</label>
                    <input v-model="form.nom" type="text" class="form-input" />
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">{{ t('admin.catalogue.category') }}</label>
                        <select v-model="form.categorie" class="form-input">
                            <option value="">{{ t('admin.catalogue.categoryChoose') }}</option>
                            <option value="atelier">{{ t('admin.catalogue.categoryAtelier') }}</option>
                            <option value="formation">{{ t('admin.catalogue.categoryFormation') }}</option>
                            <option value="service">{{ t('admin.catalogue.categoryService') }}</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label class="form-label">{{ t('admin.catalogue.price') }}</label>
                        <input v-model="form.prix" type="number" step="0.01" class="form-input" />
                    </div>
                </div>

                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">{{ t('admin.catalogue.typeLabel') }}</label>
                        <div class="radio-row">
                            <label class="form-radio-label">
                                <input type="radio" value="presentiel" v-model="form.type" />
                                {{ t('admin.catalogue.inPerson') }}
                            </label>
                            <label class="form-radio-label">
                                <input type="radio" value="en_ligne" v-model="form.type" />
                                {{ t('admin.catalogue.online') }}
                            </label>
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">{{ t('admin.catalogue.maxCapacity') }}</label>
                        <input v-model="form.max_capacity" type="number" class="form-input" :placeholder="t('admin.catalogue.maxCapacityPlaceholder')" />
                    </div>
                </div>

                <div class="form-group" v-if="form.type === 'en_ligne'">
                    <label class="form-label">{{ t('admin.catalogue.sessionLink') }}</label>
                    <input v-model="form.session_link" type="text" class="form-input" :placeholder="t('admin.catalogue.sessionLinkPlaceholder')" />
                </div>

                <div class="form-group form-group--checkbox">
                    <label class="form-checkbox-label">
                        <input type="checkbox" v-model="form.multiDay" />
                        {{ t('admin.catalogue.multiDayLabel') }}
                    </label>
                </div>

                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">{{ form.multiDay ? t('admin.catalogue.startDateLabel') : t('admin.catalogue.dateLabel') }}</label>
                        <input v-model="form.startDate" type="date" class="form-input" :min="todayDate" />
                    </div>
                    <div class="form-group" v-if="form.multiDay">
                        <label class="form-label">{{ t('admin.catalogue.endDateLabel') }}</label>
                        <input v-model="form.endDate" type="date" class="form-input" :min="form.startDate || todayDate" />
                    </div>
                </div>

                <div class="form-group" v-if="sessions.length > 0">
                    <div class="sessions-header">
                        <label class="form-label">
                            {{ form.multiDay ? t('admin.catalogue.scheduleLabelMulti', { count: sessions.length, plural: sessions.length > 1 ? 's' : '' }) : t('admin.catalogue.scheduleLabel') }}
                        </label>
                        <button v-if="form.multiDay && sessions.length > 1" type="button" class="btn-link" @click="applyFirstRowToAll">
                            {{ t('admin.catalogue.applyFirstToAll') }}
                        </button>
                    </div>
                    <div class="sessions-list">
                        <div v-for="row in sessions" :key="row.date" class="session-row">
                            <span v-if="form.multiDay" class="session-date">
                                {{ new Date(row.date + 'T00:00:00').toLocaleDateString(locale === 'en' ? 'en-US' : 'fr-FR', { weekday: 'short', day: '2-digit', month: 'short' }) }}
                            </span>
                            <input v-model="row.start_time" type="time" class="form-input session-time" :min="minTimeForDate(row.date)" />
                            <span class="session-sep">{{ t('admin.catalogue.sessionSep') }}</span>
                            <input v-model="row.end_time" type="time" class="form-input session-time" :min="row.start_time" />
                            <span class="session-duration">{{ sessionDuration(row) || '—' }}</span>
                        </div>
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">{{ t('admin.catalogue.description') }}</label>
                    <textarea v-model="form.description" class="form-input form-textarea" rows="3"></textarea>
                </div>
                <div class="form-group form-check">
                    <input id="actif" v-model="form.actif" type="checkbox" />
                    <label for="actif" class="form-label">{{ t('admin.catalogue.active') }}</label>
                </div>
                <div class="modal-actions">
                    <button class="btn-secondary" @click="showForm = false">{{ t('admin.catalogue.cancel') }}</button>
                    <button class="btn-primary" :disabled="loading" @click="save">
                        {{ loading ? t('admin.catalogue.saving') : t('admin.catalogue.save') }}
                    </button>
                </div>
            </div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.catalogue.colName') }}</th>
                        <th>{{ t('admin.catalogue.colCategory') }}</th>
                        <th>{{ t('admin.catalogue.colPrice') }}</th>
                        <th>{{ t('admin.catalogue.colStatus') }}</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="offres.length === 0">
                        <td colspan="5" class="empty">{{ t('admin.catalogue.empty') }}</td>
                    </tr>
                    <tr v-for="o in offres" :key="o.id">
                        <td class="td-bold">{{ o.nom }}</td>
                        <td class="td-muted">{{ o.categorie }}</td>
                        <td>{{ o.prix.toFixed(2) }} €</td>
                        <td>
                            <span class="badge" :class="o.actif ? 'badge--active' : 'badge--inactive'">
                                {{ o.actif ? t('admin.catalogue.activeStatus') : t('admin.catalogue.inactiveStatus') }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button v-if="!o.actif" class="btn-icon" :title="t('admin.catalogue.approve')" @click="approuver(o.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                            </button>
                            <button v-if="!o.actif" class="btn-icon" :title="t('admin.catalogue.proposeModification')" @click="openPropose(o.id)">
                                <span style="font-size: 1.1rem; line-height: 1;">📝</span>
                            </button>
                            <button v-else class="btn-icon" :title="t('admin.catalogue.disapprove')" @click="desapprouver(o.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <line x1="18" y1="6" x2="6" y2="18"></line>
                                    <line x1="6" y1="6" x2="18" y2="18"></line>
                                </svg>
                            </button>
                            <button class="btn-icon" :title="t('admin.catalogue.edit')" @click="openEdit(o)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                </svg>
                            </button>
                            <button class="btn-icon btn-icon--danger" :title="t('admin.catalogue.delete')" @click="supprimer(o.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="3 6 5 6 21 6" />
                                    <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                                </svg>
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div v-if="showProposeModal" class="form-overlay" @click.self="showProposeModal = false">
            <div class="form-modal">
                <h3 class="modal-title">{{ t('admin.catalogue.proposeModalTitle') }}</h3>
                <div class="form-group">
                    <label class="form-label">{{ t('admin.catalogue.commentLabel') }}</label>
                    <textarea v-model="proposeComment" class="form-input form-textarea" rows="4" :placeholder="t('admin.catalogue.commentPlaceholder')"></textarea>
                </div>
                <div class="modal-actions">
                    <button class="btn-secondary" @click="showProposeModal = false">{{ t('admin.catalogue.cancel') }}</button>
                    <button class="btn-primary" :disabled="loading || !proposeComment" @click="submitPropose">
                        {{ loading ? t('admin.catalogue.sending') : t('admin.catalogue.sendProposal') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 28px; }
.header-row { display: flex; justify-content: space-between; align-items: flex-start; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; white-space: nowrap; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; }
.form-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.35); z-index: 200; display: flex; align-items: center; justify-content: center; padding: 20px; }
.form-modal { background: var(--white); border-radius: 16px; padding: 32px; width: 100%; max-width: 560px; max-height: 90vh; overflow-y: auto; display: flex; flex-direction: column; gap: 18px; }
.modal-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-check { flex-direction: row; align-items: center; gap: 10px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 80px; }
.radio-row { display: flex; gap: 16px; height: 100%; align-items: center; }
.form-radio-label { display: flex; align-items: center; gap: 6px; font-size: 0.88rem; font-weight: 500; color: var(--charcoal); cursor: pointer; }
.form-checkbox-label { display: flex; align-items: center; gap: 10px; font-size: 0.9rem; font-weight: 600; color: var(--charcoal); cursor: pointer; }
.sessions-header { display: flex; justify-content: space-between; align-items: center; }
.btn-link { background: none; border: none; color: var(--green-dark); font-size: 0.78rem; font-weight: 600; cursor: pointer; text-decoration: underline; padding: 0; }
.sessions-list { display: flex; flex-direction: column; gap: 8px; max-height: 260px; overflow-y: auto; padding-right: 4px; }
.session-row { display: flex; align-items: center; gap: 8px; }
.session-date { flex: 0 0 100px; font-size: 0.82rem; font-weight: 600; color: var(--charcoal); text-transform: capitalize; }
.session-time { flex: 1; min-width: 0; }
.session-sep { font-size: 0.8rem; color: var(--charcoal); opacity: 0.5; }
.session-duration { flex: 0 0 70px; font-size: 0.78rem; font-weight: 600; color: var(--green-dark); text-align: right; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; background: #fee2e2; color: #991b1b; }
.modal-actions { display: flex; gap: 12px; justify-content: flex-end; }
.btn-secondary { padding: 11px 20px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; }
.btn-secondary:hover { border-color: var(--charcoal); }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
@media (max-width: 560px) { .form-row { grid-template-columns: 1fr; } }
</style>
