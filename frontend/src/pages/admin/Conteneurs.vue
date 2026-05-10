<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Conteneur {
    id: number
    code_barres: string
    localisation: string
    etat: string
    capacite: number
    objets: number
}

interface Locker {
    id: number
    container_id: number
    label: string
    status: string
    size: string
}

const conteneurs = ref<Conteneur[]>([])
const selectedContainer = ref<Conteneur | null>(null)
const lockers = ref<Locker[]>([])
const showModal = ref(false)
const showLockerForm = ref(false)
const editingLocker = ref<Locker | null>(null)

const lockerForm = ref({
    label: '',
    size: 'M',
    status: 'Available'
})

async function fetchContainers() {
    try {
        const res = await fetch('http://localhost:8081/admin/conteneurs', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) conteneurs.value = await res.json()
    } catch {}
}

onMounted(fetchContainers)

async function selectContainer(c: Conteneur) {
    selectedContainer.value = c
    try {
        const res = await fetch(`http://localhost:8081/container/${c.id}/lockers`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            lockers.value = await res.json()
            showModal.value = true
        }
    } catch {}
}

function openLockerForm(l: Locker | null = null) {
    if (l) {
        editingLocker.value = l
        lockerForm.value = { label: l.label, size: l.size, status: l.status }
    } else {
        editingLocker.value = null
        lockerForm.value = { label: '', size: 'M', status: 'Available' }
    }
    showLockerForm.value = true
}

async function saveLocker() {
    const url = editingLocker.value 
        ? `http://localhost:8081/admin/lockers/${editingLocker.value.id}`
        : 'http://localhost:8081/admin/lockers'
    
    const method = editingLocker.value ? 'PUT' : 'POST'
    const body = editingLocker.value 
        ? lockerForm.value 
        : { ...lockerForm.value, container_id: selectedContainer.value?.id }

    try {
        const res = await fetch(url, {
            method,
            headers: { 
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}` 
            },
            body: JSON.stringify(body)
        })
        if (res.ok) {
            showLockerForm.value = false
            if (selectedContainer.value) selectContainer(selectedContainer.value)
            fetchContainers()
        } else {
            console.error('Failed to save locker:', await res.text())
        }
    } catch (err) {
        console.error('Error saving locker:', err)
    }
}

async function deleteLocker(id: number) {
    if (!confirm('Supprimer ce casier ?')) return
    try {
        const res = await fetch(`http://localhost:8081/admin/lockers/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            if (selectedContainer.value) selectContainer(selectedContainer.value)
            fetchContainers()
        }
    } catch {}
}

function etatClass(e: string) {
    if (e === 'actif' || e === 'Available') return 'badge badge--active'
    if (e === 'plein' || e === 'Occupied') return 'badge badge--warn'
    return 'badge badge--inactive'
}
</script>

<template>
    <div class="conteneurs">
        <div class="page-header">
            <div class="header-main">
                <h1 class="page-title">Conteneurs.</h1>
                <p class="page-subtitle">Gestion des box et des casiers individuels.</p>
            </div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Code-barres</th>
                        <th>Localisation</th>
                        <th>Objets / Capacité</th>
                        <th>État</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="conteneurs.length === 0">
                        <td colspan="5" class="empty">Aucun conteneur.</td>
                    </tr>
                    <tr v-for="c in conteneurs" :key="c.id">
                        <td class="td-mono">{{ c.code_barres }}</td>
                        <td>{{ c.localisation }}</td>
                        <td>
                            <div class="progress-wrap">
                                <div class="progress-bar">
                                    <div class="progress-fill" :style="{ width: Math.min((c.objets / c.capacite) * 100, 100) + '%' }"></div>
                                </div>
                                <span class="progress-label">{{ c.objets }} / {{ c.capacite }}</span>
                            </div>
                        </td>
                        <td><span :class="etatClass(c.etat)">{{ c.etat }}</span></td>
                        <td>
                            <button class="btn-text" @click="selectContainer(c)">Gérer les casiers</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Locker Management Modal -->
        <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>Casiers : {{ selectedContainer?.code_barres }}</h2>
                    <button class="btn-close" @click="showModal = false">&times;</button>
                </div>

                <div class="modal-body">
                    <div class="locker-grid">
                        <div v-for="l in lockers" :key="l.id" class="locker-card">
                            <div class="locker-info">
                                <span class="locker-label">{{ l.label }}</span>
                                <span class="locker-size">Taille: {{ l.size }}</span>
                                <span :class="etatClass(l.status)">{{ l.status }}</span>
                            </div>
                            <div class="locker-actions">
                                <button class="btn-icon" @click="openLockerForm(l)">✎</button>
                                <button class="btn-icon btn-icon--danger" @click="deleteLocker(l.id)">🗑</button>
                            </div>
                        </div>
                        <button class="locker-add" @click="openLockerForm()">
                            <span>+ Ajouter un casier</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Locker Form Modal -->
        <div v-if="showLockerForm" class="modal-overlay modal-overlay--sub" @click.self="showLockerForm = false">
            <div class="modal-content modal-content--small">
                <h3>{{ editingLocker ? 'Modifier' : 'Nouveau' }} casier</h3>
                <div class="form-group">
                    <label>Label</label>
                    <input v-model="lockerForm.label" type="text" placeholder="Ex: A1, B2...">
                </div>
                <div class="form-group">
                    <label>Taille</label>
                    <select v-model="lockerForm.size">
                        <option value="S">Small (10 pts)</option>
                        <option value="M">Medium (20 pts)</option>
                        <option value="L">Large (50 pts)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Statut</label>
                    <select v-model="lockerForm.status">
                        <option value="Available">Disponible</option>
                        <option value="Occupied">Occupé</option>
                        <option value="HS">Hors-service</option>
                    </select>
                </div>
                <div class="modal-footer">
                    <button class="btn-secondary" @click="showLockerForm = false">Annuler</button>
                    <button class="btn-primary" @click="saveLocker">Enregistrer</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-mono { font-family: 'Courier New', monospace; font-size: 0.85rem; font-weight: 600; color: var(--green-dark); }
.progress-wrap { display: flex; align-items: center; gap: 10px; }
.progress-bar { flex: 1; max-width: 100px; height: 6px; background: rgba(53,53,53,0.1); border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; background: var(--green-mid); border-radius: 3px; transition: width 0.3s; }
.progress-label { font-size: 0.8rem; color: var(--charcoal); opacity: 0.6; white-space: nowrap; }
.btn-text { background: none; border: none; color: var(--green-dark); font-weight: 600; font-size: 0.85rem; cursor: pointer; padding: 0; text-decoration: underline; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 20px; }
.modal-overlay--sub { z-index: 1001; background: rgba(0,0,0,0.2); }
.modal-content { background: var(--white); border-radius: 20px; width: 100%; max-width: 800px; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column; }
.modal-content--small { max-width: 400px; padding: 24px; }
.modal-header { padding: 20px 24px; border-bottom: 1px solid rgba(53,53,53,0.1); display: flex; justify-content: space-between; align-items: center; }
.modal-header h2 { margin: 0; font-size: 1.4rem; color: var(--charcoal); }
.btn-close { background: none; border: none; font-size: 2rem; cursor: pointer; color: var(--charcoal); opacity: 0.5; }
.modal-body { padding: 24px; overflow-y: auto; }

.locker-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 16px; }
.locker-card { background: #f9f9f9; border: 1.5px solid rgba(53,53,53,0.05); border-radius: 12px; padding: 16px; display: flex; justify-content: space-between; align-items: flex-start; }
.locker-info { display: flex; flex-direction: column; gap: 4px; }
.locker-label { font-weight: 700; color: var(--charcoal); font-size: 1.1rem; }
.locker-size { font-size: 0.8rem; opacity: 0.6; }
.locker-actions { display: flex; gap: 4px; }
.btn-icon { background: none; border: 1.5px solid rgba(53,53,53,0.1); border-radius: 6px; padding: 4px 8px; cursor: pointer; font-size: 0.9rem; }
.btn-icon--danger { color: #dc2626; border-color: rgba(220,38,38,0.1); }
.locker-add { border: 2px dashed rgba(53,53,53,0.1); border-radius: 12px; background: none; padding: 20px; cursor: pointer; display: flex; align-items: center; justify-content: center; font-weight: 600; color: var(--green-dark); transition: all 0.2s; }
.locker-add:hover { background: rgba(215,236,225,0.3); border-color: var(--green-mid); }

.form-group { margin-bottom: 16px; display: flex; flex-direction: column; gap: 6px; }
.form-group label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.7; }
.form-group input, .form-group select { padding: 10px; border-radius: 8px; border: 1.5px solid rgba(53,53,53,0.1); font-size: 0.9rem; }
.modal-footer { margin-top: 24px; display: flex; justify-content: flex-end; gap: 12px; }
.btn-primary { background: var(--green-dark); color: var(--white); border: none; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
.btn-secondary { background: none; border: 1.5px solid rgba(53,53,53,0.1); color: var(--charcoal); padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }

.badge { display: inline-block; padding: 2px 8px; border-radius: 20px; font-size: 0.7rem; font-weight: 600; text-align: center; width: fit-content; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--warn { background: #fef3c7; color: #92400e; }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
</style>
