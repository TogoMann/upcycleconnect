<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Event {
    id: number
    approved: boolean
    price: number | null
    date: string | null
    created_by: number | null
    created_at: string | null
}

const events = ref<Event[]>([])
const showForm = ref(false)
const form = ref({ date: '', price: '' })
const saving = ref(false)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/event', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) events.value = await res.json()
    } catch {}
})

async function creer() {
    saving.value = true
    try {
        const body: Record<string, unknown> = { date: form.value.date }
        if (form.value.price) body.price = parseFloat(form.value.price)
        const res = await fetch('http://localhost:8081/event/', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify(body),
        })
        if (res.ok) {
            const created = await res.json()
            events.value.unshift(created)
            form.value = { date: '', price: '' }
            showForm.value = false
        }
    } catch {}
    saving.value = false
}

async function approuver(event: Event) {
    try {
        const res = await fetch(`http://localhost:8081/event/${event.id}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ approved: true }),
        })
        if (res.ok) event.approved = true
    } catch {}
}

async function supprimer(id: number) {
    if (!confirm('Supprimer cet évènement ?')) return
    try {
        const res = await fetch(`http://localhost:8081/event/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) events.value = events.value.filter(e => e.id !== id)
    } catch {}
}

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
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
                    <input v-model="form.date" type="datetime-local" class="field-input" required />
                </div>
                <div class="field">
                    <label class="field-label">Prix (€)</label>
                    <input v-model="form.price" type="number" step="0.01" min="0" class="field-input" placeholder="Gratuit si vide" />
                </div>
                <button type="submit" class="btn-save" :disabled="saving">
                    {{ saving ? 'Création…' : 'Créer' }}
                </button>
            </form>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Date</th>
                        <th>Prix</th>
                        <th>Statut</th>
                        <th>Créé par</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="events.length === 0">
                        <td colspan="6" class="empty">Aucun évènement.</td>
                    </tr>
                    <tr v-for="e in events" :key="e.id">
                        <td class="td-id">#{{ e.id }}</td>
                        <td>{{ fmtDate(e.date) }}</td>
                        <td>{{ e.price != null ? e.price + ' €' : 'Gratuit' }}</td>
                        <td>
                            <span class="badge" :class="e.approved ? 'badge--ok' : 'badge--pending'">
                                {{ e.approved ? 'Approuvé' : 'En attente' }}
                            </span>
                        </td>
                        <td class="td-muted">{{ e.created_by != null ? `#${e.created_by}` : '—' }}</td>
                        <td class="td-actions">
                            <button v-if="!e.approved" class="btn-sm btn-sm--approve" @click="approuver(e)">Approuver</button>
                            <button class="btn-sm btn-sm--danger" @click="supprimer(e.id)">Supprimer</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.toolbar { margin-bottom: 16px; }
.btn-create { padding: 10px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; font-family: inherit; transition: background 0.2s; }
.btn-create:hover { background: var(--green-mid); }
.form-card { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); padding: 24px; margin-bottom: 20px; }
.form-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 20px; }
.form-grid { display: flex; gap: 16px; align-items: flex-end; flex-wrap: wrap; }
.field { display: flex; flex-direction: column; gap: 6px; }
.field-label { font-size: 0.78rem; font-weight: 600; color: var(--charcoal); opacity: 0.6; text-transform: uppercase; letter-spacing: 0.04em; }
.field-input { padding: 9px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.field-input:focus { border-color: var(--green-mid); }
.btn-save { padding: 9px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; font-family: inherit; }
.btn-save:disabled { opacity: 0.5; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-id { font-weight: 600; color: rgba(53,53,53,0.45); font-size: 0.82rem; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 3px 8px; border-radius: 20px; font-size: 0.72rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; }
.badge--ok { background: var(--green-pale); color: var(--green-dark); }
.badge--pending { background: #fef3c7; color: #92400e; }
.btn-sm { padding: 5px 12px; border-radius: 6px; font-size: 0.78rem; font-weight: 600; cursor: pointer; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; white-space: nowrap; }
.btn-sm--approve { border-color: rgba(8,106,53,0.3); color: var(--green-dark); }
.btn-sm--approve:hover { border-color: var(--green-dark); background: var(--green-pale); }
.btn-sm--danger { border-color: rgba(220,38,38,0.3); color: #dc2626; }
.btn-sm--danger:hover { border-color: #dc2626; background: #fee2e2; }
</style>
