<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, reactive, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
    titre: '',
    categorie: '',
    description: '',
    type: 'presentiel' as 'presentiel' | 'en_ligne',
    session_link: '',
    statut: 'brouillon',
    multiDay: false,
    startDate: '',
    endDate: '',
    prix: '',
    max_capacity: '',
})
const error = ref('')
const loading = ref(false)

interface SessionRow {
    date: string
    start_time: string
    end_time: string
}

const sessions = ref<SessionRow[]>([])

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

async function submit() {
    if (!form.titre || !form.categorie) {
        error.value = t('salarie.formationForm.errorTitleCategoryRequired')
        return
    }
    if (sessions.value.length === 0) {
        error.value = t('salarie.formationForm.errorSelectDate')
        return
    }
    for (const row of sessions.value) {
        if (!row.start_time || !row.end_time) {
            error.value = t('salarie.formationForm.errorScheduleRequired', { date: row.date })
            return
        }
        if (row.start_time >= row.end_time) {
            error.value = t('salarie.formationForm.errorEndAfterStart', { date: row.date })
            return
        }
        if (row.date === todayDate && row.start_time < minTimeForDate(row.date)) {
            error.value = t('salarie.formationForm.errorPastTime', { date: row.date })
            return
        }
    }
    loading.value = true
    error.value = ''
    try {
        const payload = {
            titre: form.titre,
            categorie: form.categorie,
            description: form.description,
            type: form.type,
            session_link: form.session_link,
            statut: form.statut,
            prix: parseFloat(form.prix) || 0,
            max_capacity: parseInt(form.max_capacity) || null,
            sessions: sessions.value,
        }
        const res = await fetch(`${API_BASE}/salarie/formations`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(payload),
        })
        if (res.ok) {
            const data = await res.json()
            router.push(`/salarie/formations/${data.id}/edit`)
        } else {
            const d = await res.json()
            error.value = d.message ?? t('salarie.formationForm.errorCreate')
        }
    } catch {
        error.value = t('salarie.formationForm.errorNetwork')
    }
    loading.value = false
}
</script>

<template>
    <div class="nouvelle-formation">
        <div class="page-header">
            <router-link to="/salarie/formations" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="15 18 9 12 15 6" />
                </svg>
                {{ t('salarie.formationForm.back') }}
            </router-link>
            <h1 class="page-title">{{ t('salarie.formationForm.newPageTitle') }}</h1>
            <p class="page-subtitle">{{ t('salarie.formationForm.newSubtitle') }}</p>
        </div>

        <form class="form-card" @submit.prevent="submit">
            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <div class="form-group">
                <label class="form-label">{{ t('salarie.formationForm.titleLabelRequired') }}</label>
                <input v-model="form.titre" type="text" class="form-input" :placeholder="t('salarie.formationForm.titlePlaceholder')" />
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">{{ t('salarie.formationForm.categoryLabelRequired') }}</label>
                    <select v-model="form.categorie" class="form-input">
                        <option value="">{{ t('salarie.formationForm.categoryChoose') }}</option>
                        <option value="textile">{{ t('salarie.formationForm.categoryTextile') }}</option>
                        <option value="mobilier">{{ t('salarie.formationForm.categoryMobilier') }}</option>
                        <option value="electronique">{{ t('salarie.formationForm.categoryElectronique') }}</option>
                        <option value="general">{{ t('salarie.formationForm.categoryGeneral') }}</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="form-label">{{ t('salarie.formationForm.typeLabel') }}</label>
                    <div class="radio-row">
                        <label class="form-radio-label">
                            <input type="radio" value="presentiel" v-model="form.type" />
                            {{ t('salarie.formationForm.inPerson') }}
                        </label>
                        <label class="form-radio-label">
                            <input type="radio" value="en_ligne" v-model="form.type" />
                            {{ t('salarie.formationForm.online') }}
                        </label>
                    </div>
                </div>
            </div>

            <div class="form-group" v-if="form.type === 'en_ligne'">
                <label class="form-label">{{ t('salarie.formationForm.sessionLinkLabel') }}</label>
                <input v-model="form.session_link" type="text" class="form-input" :placeholder="t('salarie.formationForm.sessionLinkPlaceholder')" />
                <p class="form-hint">{{ t('salarie.formationForm.sessionLinkHint') }}</p>
            </div>

            <div class="form-group form-group--checkbox">
                <label class="form-checkbox-label">
                    <input type="checkbox" v-model="form.multiDay" />
                    {{ t('salarie.formationForm.multiDayLabel') }}
                </label>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">{{ form.multiDay ? t('salarie.formationForm.startDateLabel') : t('salarie.formationForm.dateLabel') }}</label>
                    <input v-model="form.startDate" type="date" class="form-input" :min="todayDate" />
                </div>
                <div class="form-group" v-if="form.multiDay">
                    <label class="form-label">{{ t('salarie.formationForm.endDateLabel') }}</label>
                    <input v-model="form.endDate" type="date" class="form-input" :min="form.startDate || todayDate" />
                </div>
            </div>

            <div class="form-group" v-if="sessions.length > 0">
                <div class="sessions-header">
                    <label class="form-label">
                        {{ form.multiDay ? t('salarie.formationForm.scheduleLabelMulti', { count: sessions.length, plural: sessions.length > 1 ? 's' : '' }) : t('salarie.formationForm.scheduleLabel') }}
                    </label>
                    <button v-if="form.multiDay && sessions.length > 1" type="button" class="btn-link" @click="applyFirstRowToAll">
                        {{ t('salarie.formationForm.applyFirstToAll') }}
                    </button>
                </div>
                <div class="sessions-list">
                    <div v-for="row in sessions" :key="row.date" class="session-row">
                        <span v-if="form.multiDay" class="session-date">
                            {{ new Date(row.date + 'T00:00:00').toLocaleDateString(locale === 'en' ? 'en-US' : 'fr-FR', { weekday: 'short', day: '2-digit', month: 'short' }) }}
                        </span>
                        <input v-model="row.start_time" type="time" class="form-input session-time" :min="minTimeForDate(row.date)" />
                        <span class="session-sep">{{ t('salarie.formationForm.sessionSep') }}</span>
                        <input v-model="row.end_time" type="time" class="form-input session-time" :min="row.start_time" />
                        <span class="session-duration">{{ sessionDuration(row) || '—' }}</span>
                    </div>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label class="form-label">{{ t('salarie.formationForm.priceLabel') }}</label>
                    <input v-model="form.prix" type="number" step="0.01" class="form-input" placeholder="0.00" />
                </div>
                <div class="form-group">
                    <label class="form-label">{{ t('salarie.formationForm.capacityLabel') }}</label>
                    <input v-model="form.max_capacity" type="number" class="form-input" :placeholder="t('salarie.formationForm.capacityPlaceholder')" />
                </div>
            </div>

            <div class="form-group">
                <label class="form-label">{{ t('salarie.formationForm.descriptionLabel') }}</label>
                <textarea v-model="form.description" class="form-input form-textarea" rows="5" :placeholder="t('salarie.formationForm.descriptionPlaceholder')"></textarea>
            </div>

            <div class="form-group">
                <label class="form-label">{{ t('salarie.formationForm.statusLabel') }}</label>
                <select v-model="form.statut" class="form-input">
                    <option value="brouillon">{{ t('salarie.formationForm.statusDraft') }}</option>
                    <option value="publiee">{{ t('salarie.formationForm.statusSubmit') }}</option>
                </select>
                <p class="form-hint">{{ t('salarie.formationForm.statusHint') }}</p>
            </div>

            <div class="form-actions">
                <router-link to="/salarie/formations" class="btn-secondary">{{ t('salarie.formationForm.cancel') }}</router-link>
                <button type="submit" class="btn-primary" :disabled="loading">
                    {{ loading ? t('salarie.formationForm.creating') : t('salarie.formationForm.create') }}
                </button>
            </div>
        </form>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.back-link { display: inline-flex; align-items: center; gap: 6px; font-size: 0.85rem; color: var(--green-mid); text-decoration: none; margin-bottom: 16px; transition: color 0.2s; }
.back-link:hover { color: var(--green-dark); }
.back-link svg { width: 16px; height: 16px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.form-card { max-width: 640px; background: var(--white); border-radius: 16px; border: 1.5px solid rgba(53,53,53,0.1); padding: 32px; display: flex; flex-direction: column; gap: 20px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-hint { font-size: 0.78rem; color: var(--charcoal); opacity: 0.5; margin: 0; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 120px; }
.radio-row { display: flex; gap: 16px; height: 100%; align-items: center; }
.form-radio-label { display: flex; align-items: center; gap: 6px; font-size: 0.88rem; font-weight: 500; color: var(--charcoal); cursor: pointer; }
.form-checkbox-label { display: flex; align-items: center; gap: 10px; font-size: 0.9rem; font-weight: 600; color: var(--charcoal); cursor: pointer; }
.sessions-header { display: flex; justify-content: space-between; align-items: center; }
.btn-link { background: none; border: none; color: var(--green-dark); font-size: 0.78rem; font-weight: 600; cursor: pointer; text-decoration: underline; padding: 0; }
.sessions-list { display: flex; flex-direction: column; gap: 8px; max-height: 320px; overflow-y: auto; padding-right: 4px; }
.session-row { display: flex; align-items: center; gap: 8px; }
.session-date { flex: 0 0 100px; font-size: 0.82rem; font-weight: 600; color: var(--charcoal); text-transform: capitalize; }
.session-time { flex: 1; min-width: 0; }
.session-sep { font-size: 0.8rem; color: var(--charcoal); opacity: 0.5; }
.session-duration { flex: 0 0 70px; font-size: 0.78rem; font-weight: 600; color: var(--green-dark); text-align: right; }
.form-actions { display: flex; gap: 12px; justify-content: flex-end; padding-top: 8px; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-secondary { padding: 12px 24px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; }
.btn-secondary:hover { border-color: var(--charcoal); }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; background: #fee2e2; color: #991b1b; }
@media (max-width: 560px) { .form-row { grid-template-columns: 1fr; } }
</style>
