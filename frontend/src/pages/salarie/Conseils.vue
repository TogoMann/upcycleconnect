<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Article {
    id: number
    titre: string
    categorie: string
    date: string
    statut: string
}

const articles = ref<Article[]>([])
const showForm = ref(false)
const editId = ref<number | null>(null)
const form = ref({ titre: '', categorie: '', statut: 'publie' })
const loading = ref(false)

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/salarie/conseils', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) articles.value = await res.json()
    } catch {}
})

function openCreate() {
    editId.value = null
    form.value = { titre: '', categorie: '', statut: 'publie' }
    showForm.value = true
}

function openEdit(a: Article) {
    editId.value = a.id
    form.value = { titre: a.titre, categorie: a.categorie, statut: a.statut }
    showForm.value = true
}

async function save() {
    loading.value = true
    try {
        if (editId.value) {
            const res = await fetch(`http://localhost:8081/salarie/conseils/${editId.value}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
                body: JSON.stringify(form.value),
            })
            if (res.ok) {
                const updated = await res.json()
                const idx = articles.value.findIndex(a => a.id === editId.value)
                if (idx !== -1) articles.value[idx] = updated
            }
        } else {
            const res = await fetch('http://localhost:8081/salarie/conseils', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
                body: JSON.stringify(form.value),
            })
            if (res.ok) articles.value.unshift(await res.json())
        }
    } catch {}
    showForm.value = false
    loading.value = false
}

async function supprimer(id: number) {
    if (!confirm('Supprimer cet article ?')) return
    await fetch(`http://localhost:8081/salarie/conseils/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    articles.value = articles.value.filter(a => a.id !== id)
}
</script>

<template>
    <div class="conseils">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">Conseils.</h1>
                    <p class="page-subtitle">Créez, modifiez et supprimez vos articles de conseils.</p>
                </div>
                <button class="btn-primary" @click="openCreate">+ Nouvel article</button>
            </div>
        </div>

        <div v-if="showForm" class="form-overlay" @click.self="showForm = false">
            <div class="form-modal">
                <h3 class="modal-title">{{ editId ? 'Modifier l\'article' : 'Nouvel article' }}</h3>
                <div class="form-group">
                    <label class="form-label">Titre</label>
                    <input v-model="form.titre" type="text" class="form-input" />
                </div>
                <div class="form-group">
                    <label class="form-label">Catégorie</label>
                    <select v-model="form.categorie" class="form-input">
                        <option value="">Choisir</option>
                        <option value="upcycling">Upcycling</option>
                        <option value="ecologie">Écologie</option>
                        <option value="diy">DIY</option>
                    </select>
                </div>
                <div class="form-group">
                    <label class="form-label">Statut</label>
                    <select v-model="form.statut" class="form-input">
                        <option value="publie">Publié</option>
                        <option value="brouillon">Brouillon</option>
                    </select>
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
                        <th>Titre</th>
                        <th>Catégorie</th>
                        <th>Date</th>
                        <th>Statut</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="articles.length === 0">
                        <td colspan="5" class="empty">Aucun article.</td>
                    </tr>
                    <tr v-for="a in articles" :key="a.id">
                        <td class="td-bold">{{ a.titre }}</td>
                        <td class="td-muted">{{ a.categorie }}</td>
                        <td class="td-muted">{{ a.date }}</td>
                        <td>
                            <span class="badge" :class="a.statut === 'publie' ? 'badge--active' : 'badge--draft'">
                                {{ a.statut === 'publie' ? 'Publié' : 'Brouillon' }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button class="btn-icon" @click="openEdit(a)" title="Modifier">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                </svg>
                            </button>
                            <button class="btn-icon btn-icon--danger" @click="supprimer(a.id)" title="Supprimer">
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
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; white-space: nowrap; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; }
.form-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.35); z-index: 200; display: flex; align-items: center; justify-content: center; padding: 20px; }
.form-modal { background: var(--white); border-radius: 16px; padding: 32px; width: 100%; max-width: 480px; display: flex; flex-direction: column; gap: 18px; }
.modal-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
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
.td-actions { display: flex; gap: 8px; align-items: center; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--draft { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
</style>
