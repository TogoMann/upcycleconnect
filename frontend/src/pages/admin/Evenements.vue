<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface PgtypeTime {
    Microseconds: number
    Valid: boolean
}

interface Event {
    id: number
    approved: boolean
    approved_by: number | null
    approved_at: string | null
    price: number | null
    date: string | null
    start_time: PgtypeTime | null
    end_time: PgtypeTime | null
    location: string | null
    created_by: number | null
    created_at: string | null
    creator_username: string | null
    creator_email: string | null
    approver_username: string | null
    approver_email: string | null
}

const events = ref<Event[]>([])
const showForm = ref(false)
const form = ref({ date: '', start_time: '', end_time: '', location: '', price: '' })
const saving = ref(false)
const selectedEvent = ref<Event | null>(null)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/events', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) events.value = await res.json()
    } catch {}
})

async function creer() {
    saving.value = true
    try {
        const body: Record<string, unknown> = {
            date: form.value.date,
            start_time: form.value.start_time + ':00',
            end_time: form.value.end_time + ':00',
            location: form.value.location
        }
        if (form.value.price) body.price = parseFloat(form.value.price)
        const res = await fetch('http://localhost:8081/event/', {
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
            form.value = { date: '', start_time: '', end_time: '', location: '', price: '' }
            showForm.value = false
        }
    } catch {}
    saving.value = false
}

async function approuver(event: Event) {
    try {
        const res = await fetch(`http://localhost:8081/event/${event.id}/approve`, {
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
        const res = await fetch(`http://localhost:8081/event/${event.id}/disapprove`, {
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
    if (!confirm('Supprimer cet évènement ?')) return
    try {
        const res = await fetch(`http://localhost:8081/event/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) events.value = events.value.filter((e) => e.id !== id)
    } catch {}
}

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    })
}

function fmtTime(t: PgtypeTime | null): string {
    if (!t || !t.Valid) return '—'
    const totalSeconds = Math.floor(t.Microseconds / 1_000_000)
    const h = Math.floor(totalSeconds / 3600).toString().padStart(2, '0')
    const m = Math.floor((totalSeconds % 3600) / 60).toString().padStart(2, '0')
    return `${h}:${m}`
}

function fmtTimestamp(ts: string | null): string {
    if (!ts) return '—'
    return new Date(ts).toLocaleString('fr-FR', {
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
            <h1 class="page-title">Évènements.</h1>
            <p class="page-subtitle">Gérez les évènements de la plateforme.</p>
        </div>

        <div class="toolbar">
            <button class="btn-create" @click="showForm = !showForm">
                {{ showForm ? 'Annuler' : '+ Nouvel évènement' }}
            </button>
        </div>

        <div v-if="showForm" class="form-card">
            <h3 class="form-title">Nouvel évènement</h3>
            <form @submit.prevent="creer" class="form-grid">
                <div class="field">
                    <label class="field-label">Date</label>
                    <input v-model="form.date" type="date" class="field-input" required :min="new Date().toISOString().split('T')[0]" />
                </div>
                <div class="field">
                    <label class="field-label">Début</label>
                    <input v-model="form.start_time" type="time" class="field-input" required />
                </div>
                <div class="field">
                    <label class="field-label">Fin</label>
                    <input v-model="form.end_time" type="time" class="field-input" required />
                </div>
                <div class="field">
                    <label class="field-label">Lieu</label>
                    <input v-model="form.location" type="text" class="field-input" placeholder="ex: Salle Polyvalente" required />
                </div>
                <div class="field">
                    <label class="field-label">Prix (€)</label>
                    <input
                        v-model="form.price"
                        type="number"
                        step="0.01"
                        min="0"
                        class="field-input"
                        placeholder="Gratuit si vide"
                    />
                </div>
                <button type="submit" class="btn-save" :disabled="saving">
                    {{ saving ? 'Création…' : 'Créer' }}
                </button>
            </form>
        </div>

        <Teleport to="body">
            <div v-if="selectedEvent" class="modal-overlay" @click.self="selectedEvent = null">
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title">Détails — Évènement #{{ selectedEvent.id }}</h3>
                        <button class="modal-close" @click="selectedEvent = null">✕</button>
                    </div>
                    <dl class="modal-dl">
                        <dt>Créé par</dt>
                        <dd>
                            <template v-if="selectedEvent.creator_username">
                                {{ selectedEvent.creator_username }}
                                <span class="modal-email">{{ selectedEvent.creator_email }}</span>
                            </template>
                            <template v-else>—</template>
                        </dd>
                        <dt>Créé le</dt>
                        <dd>{{ fmtTimestamp(selectedEvent.created_at) }}</dd>
                        <dt>Approuvé par</dt>
                        <dd>
                            <template v-if="selectedEvent.approver_username">
                                {{ selectedEvent.approver_username }}
                                <span class="modal-email">{{ selectedEvent.approver_email }}</span>
                            </template>
                            <template v-else>—</template>
                        </dd>
                        <dt>Approuvé le</dt>
                        <dd>{{ fmtTimestamp(selectedEvent.approved_at) }}</dd>
                    </dl>
                </div>
            </div>
        </Teleport>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Date</th>
                        <th>Début</th>
                        <th>Fin</th>
                        <th>Lieu</th>
                        <th>Prix</th>
                        <th>Statut</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="events.length === 0">
                        <td colspan="8" class="empty">Aucun évènement.</td>
                    </tr>
                    <tr v-for="e in events" :key="e.id">
                        <td class="td-id">#{{ e.id }}</td>
                        <td>{{ fmtDate(e.date) }}</td>
                        <td>{{ fmtTime(e.start_time) }}</td>
                        <td>{{ fmtTime(e.end_time) }}</td>
                        <td>{{ e.location || '—' }}</td>
                        <td>{{ e.price != null ? e.price + ' €' : 'Gratuit' }}</td>
                        <td>
                            <span
                                class="badge"
                                :class="e.approved ? 'badge--ok' : 'badge--pending'"
                            >
                                {{ e.approved ? 'Approuvé' : 'En attente' }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button class="btn-sm btn-sm--info" @click="selectedEvent = e">Détails</button>
                            <button v-if="!e.approved" class="btn-sm btn-sm--approve" @click="approuver(e)">Approuver</button>
                            <button v-else class="btn-sm" @click="desapprouver(e)">Désapprouver</button>
                            <button class="btn-sm btn-sm--danger" @click="supprimer(e.id)">Supprimer</button>
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
</style>
