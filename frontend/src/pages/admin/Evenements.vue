<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface PgtypeTime {
    Microseconds: number
    Valid: boolean
}

interface Event {
    id: number
    title: string
    description: string | null
    approved: boolean
    approved_by: number | null
    approved_at: string | null
    price: number | null
    date: string | null
    start_time: PgtypeTime | null
    end_time: PgtypeTime | null
    location: string | null
    max_capacity: any | null
    created_by: number | null
    created_at: string | null
    creator_username: string | null
    creator_email: string | null
    approver_username: string | null
    approver_email: string | null
    premium?: boolean
}

const events = ref<Event[]>([])
const showForm = ref(false)
const form = ref({ title: '', description: '', date: '', start_time: '', end_time: '', location: '', price: '', max_capacity: '', premium: false })
const saving = ref(false)
const selectedEvent = ref<Event | null>(null)

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/admin/events`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) events.value = await res.json()
    } catch {}
})

async function creer() {
    saving.value = true
    try {
        const body: Record<string, unknown> = {
            title: form.value.title,
            description: form.value.description,
            date: form.value.date,
            start_time: form.value.start_time + ':00',
            end_time: form.value.end_time + ':00',
            location: form.value.location,
            premium: form.value.premium
        }
        if (form.value.price) body.price = parseFloat(form.value.price)
        if (form.value.max_capacity) body.max_capacity = parseInt(form.value.max_capacity, 10)
        const res = await fetch(`${API_BASE}/event/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(body),
        })
        if (res.ok) {
            const created = await res.json()
            events.value.unshift(created)
            form.value = { title: '', description: '', date: '', start_time: '', end_time: '', location: '', price: '', max_capacity: '', premium: false }
            showForm.value = false
        }
    } catch {}
    saving.value = false
}

async function approuver(event: Event) {
    try {
        const res = await fetch(`${API_BASE}/event/${event.id}/approve`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({ approved: true }),
        })
        if (res.ok) event.approved = true
    } catch {}
}

async function desapprouver(event: Event) {
    try {
        const res = await fetch(`${API_BASE}/event/${event.id}/disapprove`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
        })
        if (res.ok) event.approved = false
    } catch {}
}

async function supprimer(id: number) {
    if (!confirm(t('admin.evenements.confirmDelete'))) return
    try {
        const res = await fetch(`${API_BASE}/event/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) events.value = events.value.filter((e) => e.id !== id)
    } catch {}
}

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    })
}

function getLocalDateString(): string {
    const d = new Date()
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
}

function fmtTime(time: PgtypeTime | null): string {
    if (!time || !time.Valid) return '—'
    const totalSeconds = Math.floor(time.Microseconds / 1_000_000)
    const h = Math.floor(totalSeconds / 3600).toString().padStart(2, '0')
    const m = Math.floor((totalSeconds % 3600) / 60).toString().padStart(2, '0')
    return `${h}:${m}`
}

function fmtTimestamp(ts: string | null): string {
    if (!ts) return '—'
    return new Date(ts).toLocaleString(locale.value === 'en' ? 'en-US' : 'fr-FR', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}
</script>

<template>
    <div class="evenements">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.evenements.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.evenements.subtitle') }}</p>
        </div>

        <div class="toolbar">
            <button class="btn-create" @click="showForm = !showForm">
                {{ showForm ? t('admin.evenements.cancel') : t('admin.evenements.newEvent') }}
            </button>
        </div>

        <div v-if="showForm" class="form-card">
            <h3 class="form-title">{{ t('admin.evenements.newEventTitle') }}</h3>
            <form @submit.prevent="creer" class="form-grid">
                <div class="field field--wide">
                    <label class="field-label">{{ t('admin.evenements.titleLabel') }}</label>
                    <input v-model="form.title" type="text" class="field-input" :placeholder="t('admin.evenements.titlePlaceholder')" required />
                </div>
                <div class="field field--wide">
                    <label class="field-label">{{ t('admin.evenements.descriptionLabel') }}</label>
                    <input v-model="form.description" type="text" class="field-input" :placeholder="t('admin.evenements.descriptionPlaceholder')" />
                </div>
                <div class="field">
                    <label class="field-label">{{ t('admin.evenements.dateLabel') }}</label>
                    <input v-model="form.date" type="date" class="field-input" required :min="getLocalDateString()" />
                </div>
                <div class="field">
                    <label class="field-label">{{ t('admin.evenements.startLabel') }}</label>
                    <input v-model="form.start_time" type="time" class="field-input" required />
                </div>
                <div class="field">
                    <label class="field-label">{{ t('admin.evenements.endLabel') }}</label>
                    <input v-model="form.end_time" type="time" class="field-input" required />
                </div>
                <div class="field">
                    <label class="field-label">{{ t('admin.evenements.locationLabel') }}</label>
                    <input v-model="form.location" type="text" class="field-input" :placeholder="t('admin.evenements.locationPlaceholder')" required />
                </div>
                <div class="field">
                    <label class="field-label">{{ t('admin.evenements.priceLabel') }}</label>
                    <input
                        v-model="form.price"
                        type="number"
                        step="0.01"
                        min="0"
                        class="field-input"
                        :placeholder="t('admin.evenements.pricePlaceholder')"
                    />
                </div>
                <div class="field">
                    <label class="field-label">{{ t('admin.evenements.maxCapacityLabel') }}</label>
                    <input
                        v-model="form.max_capacity"
                        type="number"
                        min="1"
                        class="field-input"
                        :placeholder="t('admin.evenements.maxCapacityPlaceholder')"
                    />
                </div>
                <div class="field field--checkbox">
                    <label class="checkbox-label">
                        <input type="checkbox" v-model="form.premium" class="checkbox-input" />
                        <span>{{ t('admin.evenements.premiumLabel') }}</span>
                    </label>
                </div>
                <button type="submit" class="btn-save" :disabled="saving">
                    {{ saving ? t('admin.evenements.creating') : t('admin.evenements.create') }}
                </button>
            </form>
        </div>

        <Teleport to="body">
            <div v-if="selectedEvent" class="modal-overlay" @click.self="selectedEvent = null">
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title">
                            {{ selectedEvent.title }}
                            <span v-if="selectedEvent.premium" class="badge-premium-inline">Premium</span>
                        </h3>
                        <button class="modal-close" @click="selectedEvent = null">✕</button>
                    </div>
                    <dl class="modal-dl">
                        <dt>{{ t('admin.evenements.description') }}</dt>
                        <dd>{{ selectedEvent.description || '—' }}</dd>
                        <dt>{{ t('admin.evenements.createdBy') }}</dt>
                        <dd>
                            <template v-if="selectedEvent.creator_username">
                                {{ selectedEvent.creator_username }}
                                <span class="modal-email">{{ selectedEvent.creator_email }}</span>
                            </template>
                            <template v-else>—</template>
                        </dd>
                        <dt>{{ t('admin.evenements.createdOn') }}</dt>
                        <dd>{{ fmtTimestamp(selectedEvent.created_at) }}</dd>
                        <dt>{{ t('admin.evenements.approvedBy') }}</dt>
                        <dd>
                            <template v-if="selectedEvent.approver_username">
                                {{ selectedEvent.approver_username }}
                                <span class="modal-email">{{ selectedEvent.approver_email }}</span>
                            </template>
                            <template v-else>—</template>
                        </dd>
                        <dt>{{ t('admin.evenements.approvedOn') }}</dt>
                        <dd>{{ fmtTimestamp(selectedEvent.approved_at) }}</dd>
                        <dt>{{ t('admin.evenements.capacity') }}</dt>
                        <dd>{{ selectedEvent.max_capacity != null && selectedEvent.max_capacity.Valid !== false ? t('admin.evenements.maxParticipants', { count: selectedEvent.max_capacity.Int32 ?? selectedEvent.max_capacity }) : t('admin.evenements.unlimited') }}</dd>
                    </dl>
                </div>
            </div>
        </Teleport>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.evenements.colId') }}</th>
                        <th>{{ t('admin.evenements.colTitle') }}</th>
                        <th>{{ t('admin.evenements.colDate') }}</th>
                        <th>{{ t('admin.evenements.colStart') }}</th>
                        <th>{{ t('admin.evenements.colEnd') }}</th>
                        <th>{{ t('admin.evenements.colLocation') }}</th>
                        <th>{{ t('admin.evenements.colCapacity') }}</th>
                        <th>{{ t('admin.evenements.colPrice') }}</th>
                        <th>{{ t('admin.evenements.colStatus') }}</th>
                        <th>{{ t('admin.evenements.colActions') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="events.length === 0">
                        <td colspan="10" class="empty">{{ t('admin.evenements.empty') }}</td>
                    </tr>
                    <tr v-for="e in events" :key="e.id">
                        <td class="td-id">#{{ e.id }}</td>
                        <td class="td-bold">
                            {{ e.title }}
                            <span v-if="e.premium" class="badge-premium-inline">Premium</span>
                        </td>
                        <td>{{ fmtDate(e.date) }}</td>
                        <td>{{ fmtTime(e.start_time) }}</td>
                        <td>{{ fmtTime(e.end_time) }}</td>
                        <td>{{ e.location || '—' }}</td>
                        <td>{{ e.max_capacity != null && e.max_capacity.Valid !== false ? t('admin.evenements.seats', { count: e.max_capacity.Int32 ?? e.max_capacity }) : t('admin.evenements.unlimited') }}</td>
                        <td>{{ e.price != null ? e.price + ' €' : t('admin.evenements.free') }}</td>
                        <td>
                            <span
                                class="badge"
                                :class="e.approved ? 'badge--ok' : 'badge--pending'"
                            >
                                {{ e.approved ? t('admin.evenements.approved') : t('admin.evenements.pending') }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button class="btn-sm btn-sm--info" @click="selectedEvent = e">{{ t('admin.evenements.details') }}</button>
                            <button v-if="!e.approved" class="btn-sm btn-sm--approve" @click="approuver(e)">{{ t('admin.evenements.approve') }}</button>
                            <button v-else class="btn-sm" @click="desapprouver(e)">{{ t('admin.evenements.disapprove') }}</button>
                            <button class="btn-sm btn-sm--danger" @click="supprimer(e.id)">{{ t('admin.evenements.delete') }}</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header {
    margin-bottom: 24px;
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 8px;
    line-height: 1.08;
}
.page-subtitle {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.6;
    margin: 0;
}
.toolbar {
    margin-bottom: 16px;
}
.btn-create {
    padding: 10px 20px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.88rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
}
.btn-create:hover {
    background: var(--green-mid);
}
.form-card {
    background: var(--white);
    border-radius: 14px;
    border: 1.5px solid rgba(53, 53, 53, 0.08);
    padding: 24px;
    margin-bottom: 20px;
}
.form-title {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 20px;
}
.form-grid {
    display: flex;
    gap: 16px;
    align-items: flex-end;
    flex-wrap: wrap;
}
.field {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.field--wide {
    flex-basis: 100%;
}
.field-label {
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--charcoal);
    opacity: 0.6;
    text-transform: uppercase;
    letter-spacing: 0.04em;
}
.field-input {
    padding: 9px 14px;
    font-size: 0.9rem;
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    border-radius: 8px;
    background: var(--white);
    color: var(--charcoal);
    font-family: inherit;
    outline: none;
    transition: border-color 0.2s;
}
.field-input:focus {
    border-color: var(--green-mid);
}
.btn-save {
    padding: 9px 20px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.88rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
}
.btn-save:disabled {
    opacity: 0.5;
}
.table-wrap {
    background: var(--white);
    border-radius: 14px;
    border: 1.5px solid rgba(53, 53, 53, 0.08);
    overflow: hidden;
}
.data-table {
    width: 100%;
    border-collapse: collapse;
}
.data-table th {
    text-align: left;
    padding: 14px 20px;
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--charcoal);
    opacity: 0.5;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
}
.data-table td {
    padding: 14px 20px;
    font-size: 0.9rem;
    color: var(--charcoal);
    border-bottom: 1px solid rgba(53, 53, 53, 0.05);
    vertical-align: middle;
}
.data-table tr:last-child td {
    border-bottom: none;
}
.data-table tbody tr:hover {
    background: rgba(215, 236, 225, 0.3);
}
.td-id {
    font-weight: 600;
    color: rgba(53, 53, 53, 0.45);
    font-size: 0.82rem;
}
.td-muted {
    opacity: 0.55;
    font-size: 0.85rem;
}
.td-bold {
    font-weight: 600;
}
.td-actions {
    display: flex;
    gap: 8px;
}
.empty {
    text-align: center;
    opacity: 0.4;
    padding: 40px !important;
}
.badge {
    display: inline-block;
    padding: 3px 8px;
    border-radius: 20px;
    font-size: 0.72rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
}
.badge--ok {
    background: var(--green-pale);
    color: var(--green-dark);
}
.badge--pending {
    background: #fef3c7;
    color: #92400e;
}
.btn-sm {
    padding: 5px 12px;
    border-radius: 6px;
    font-size: 0.78rem;
    font-weight: 600;
    cursor: pointer;
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    background: transparent;
    color: var(--charcoal);
    transition:
        border-color 0.2s,
        color 0.2s;
    white-space: nowrap;
}
.btn-sm--approve {
    border-color: rgba(8, 106, 53, 0.3);
    color: var(--green-dark);
}
.btn-sm--approve:hover {
    border-color: var(--green-dark);
    background: var(--green-pale);
}
.btn-sm--danger {
    border-color: rgba(220, 38, 38, 0.3);
    color: #dc2626;
}
.btn-sm--danger:hover {
    border-color: #dc2626;
    background: #fee2e2;
}
.btn-sm--info {
    border-color: rgba(37, 99, 235, 0.25);
    color: #2563eb;
}
.btn-sm--info:hover {
    border-color: #2563eb;
    background: #eff6ff;
}

.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}
.modal {
    background: #fff;
    border-radius: 14px;
    padding: 28px 32px;
    min-width: 320px;
    max-width: 480px;
    width: 100%;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}
.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
}
.modal-title {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0;
}
.modal-close {
    background: none;
    border: none;
    font-size: 1rem;
    cursor: pointer;
    color: rgba(53, 53, 53, 0.5);
    padding: 4px;
    line-height: 1;
}
.modal-close:hover {
    color: var(--charcoal);
}
.modal-dl {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 10px 16px;
    margin: 0;
}
.modal-dl dt {
    font-size: 0.78rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: rgba(53, 53, 53, 0.5);
    align-self: center;
}
.modal-dl dd {
    font-size: 0.9rem;
    color: var(--charcoal);
    margin: 0;
}
.modal-email {
    display: block;
    font-size: 0.78rem;
    color: rgba(53, 53, 53, 0.5);
}
.badge-premium-inline {
    display: inline-block;
    background: #ebd3ff;
    color: #6b21a8;
    font-size: 0.72rem;
    font-weight: 700;
    padding: 2px 6px;
    border-radius: 4px;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-left: 6px;
    vertical-align: middle;
}
.field--checkbox {
    grid-column: 1 / -1;
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 4px;
}
.checkbox-label {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--charcoal);
    cursor: pointer;
}
.checkbox-input {
    width: 16px;
    height: 16px;
    cursor: pointer;
}
</style>
