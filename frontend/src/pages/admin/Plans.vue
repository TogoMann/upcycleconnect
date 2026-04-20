<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Plan {
    id: number
    name: string
    description: string
    price: number
    billing_cycle: string
    features: string[]
    is_active: boolean
}

const plans = ref<Plan[]>([])
const showForm = ref(false)
const editId = ref<number | null>(null)
const loading = ref(false)

const form = ref({
    name: '',
    description: '',
    price: 0,
    billing_cycle: 'monthly',
    features: '',
    is_active: true
})

onMounted(fetchPlans)

async function fetchPlans() {
    try {
        const res = await fetch('http://localhost:8081/plans')
        if (res.ok) plans.value = await res.json()
    } catch {}
}

function openCreate() {
    editId.value = null
    form.value = { name: '', description: '', price: 0, billing_cycle: 'monthly', features: '', is_active: true }
    showForm.value = true
}

function openEdit(p: Plan) {
    editId.value = p.id
    form.value = { 
        name: p.name, 
        description: p.description, 
        price: p.price, 
        billing_cycle: p.billing_cycle, 
        features: p.features.join('\n'), 
        is_active: p.is_active 
    }
    showForm.value = true
}

async function save() {
    loading.value = true
    const payload = { 
        ...form.value, 
        features: form.value.features.split('\n').filter(f => f.trim() !== '') 
    }
    
    try {
        const url = editId.value 
            ? `http://localhost:8081/admin/plans/${editId.value}` 
            : 'http://localhost:8081/admin/plans'
        const method = editId.value ? 'PUT' : 'POST'
        
        const res = await fetch(url, {
            method,
            headers: { 
                'Content-Type': 'application/json', 
                Authorization: `Bearer ${authStore.token}` 
            },
            body: JSON.stringify(payload),
        })
        
        if (res.ok) {
            await fetchPlans()
            showForm.value = false
        }
    } catch {}
    loading.value = false
}

async function supprimer(id: number) {
    if (!confirm('Supprimer ce plan ?')) return
    try {
        const res = await fetch(`http://localhost:8081/admin/plans/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) fetchPlans()
    } catch {}
}
</script>

<template>
    <div class="plans-admin">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">Plans & Tiers.</h1>
                    <p class="page-subtitle">Gérez les offres d'abonnement et bonus.</p>
                </div>
                <button class="btn-primary" @click="openCreate">+ Nouveau Plan</button>
            </div>
        </div>

        <div v-if="showForm" class="form-overlay" @click.self="showForm = false">
            <div class="form-modal">
                <h3 class="modal-title">{{ editId ? 'Modifier le plan' : 'Nouveau plan' }}</h3>
                
                <div class="form-group">
                    <label class="form-label">Nom</label>
                    <input v-model="form.name" type="text" class="form-input" placeholder="ex: Premium" />
                </div>

                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">Prix (€)</label>
                        <input v-model="form.price" type="number" step="0.01" class="form-input" />
                    </div>
                    <div class="form-group">
                        <label class="form-label">Cycle</label>
                        <select v-model="form.billing_cycle" class="form-input">
                            <option value="monthly">Mensuel</option>
                            <option value="yearly">Annuel</option>
                            <option value="once">Une fois</option>
                        </select>
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">Description</label>
                    <textarea v-model="form.description" class="form-input form-textarea" rows="2"></textarea>
                </div>

                <div class="form-group">
                    <label class="form-label">Avantages (un par ligne)</label>
                    <textarea v-model="form.features" class="form-input form-textarea" rows="4" placeholder="Bonus 1&#10;Bonus 2"></textarea>
                </div>

                <div class="form-group form-check">
                    <input id="active" v-model="form.is_active" type="checkbox" />
                    <label for="active" class="form-label">Plan actif</label>
                </div>

                <div class="modal-actions">
                    <button class="btn-secondary" @click="showForm = false">Annuler</button>
                    <button class="btn-primary" :disabled="loading" @click="save">
                        {{ loading ? '…' : 'Enregistrer' }}
                    </button>
                </div>
            </div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Plan</th>
                        <th>Prix</th>
                        <th>Cycle</th>
                        <th>Statut</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="plans.length === 0">
                        <td colspan="5" class="empty">Aucun plan configuré.</td>
                    </tr>
                    <tr v-for="p in plans" :key="p.id">
                        <td class="td-bold">{{ p.name }}</td>
                        <td>{{ p.price.toFixed(2) }} €</td>
                        <td class="td-muted">{{ p.billing_cycle }}</td>
                        <td>
                            <span class="badge" :class="p.is_active ? 'badge--active' : 'badge--inactive'">
                                {{ p.is_active ? 'Actif' : 'Inactif' }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button class="btn-icon" @click="openEdit(p)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                </svg>
                            </button>
                            <button class="btn-icon btn-icon--danger" @click="supprimer(p.id)">
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
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 28px; }
.header-row { display: flex; justify-content: space-between; align-items: flex-start; }
.page-title { font-size: 2.6rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; }
.form-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.35); z-index: 200; display: flex; align-items: center; justify-content: center; padding: 20px; }
.form-modal { background: var(--white); border-radius: 16px; padding: 32px; width: 100%; max-width: 500px; display: flex; flex-direction: column; gap: 18px; }
.modal-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-check { flex-direction: row; align-items: center; gap: 10px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; }
.modal-actions { display: flex; gap: 12px; justify-content: flex-end; }
.btn-secondary { padding: 11px 20px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
</style>
