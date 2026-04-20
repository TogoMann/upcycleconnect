<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Annonce {
    id: number
    name: string
    category: string
    price: number
    description: string
    approved: boolean
}

const annonces = ref<Annonce[]>([])
const loading = ref(false)

onMounted(async () => {
    loading.value = true
    try {
        const res = await fetch('http://localhost:8081/admin/listings', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) annonces.value = await res.json()
    } catch {}
    loading.value = false
})

async function deleteAnnonce(id: number) {
    if (!confirm('Supprimer cette annonce ?')) return
    try {
        const res = await fetch(`http://localhost:8081/listing/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) annonces.value = annonces.value.filter(a => a.id !== id)
    } catch {}
}

async function approveAnnonce(id: number) {
    try {
        const res = await fetch(`http://localhost:8081/listing/${id}/approve`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const idx = annonces.value.findIndex(a => a.id === id)
            if (idx !== -1) annonces.value[idx].approved = true
        }
    } catch {}
}

async function disapproveAnnonce(id: number) {
    try {
        const res = await fetch(`http://localhost:8081/listing/${id}/disapprove`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const idx = annonces.value.findIndex(a => a.id === id)
            if (idx !== -1) annonces.value[idx].approved = false
        }
    } catch {}
}
</script>

<template>
    <div class="admin-annonces">
        <div class="page-header">
            <h1 class="page-title">Annonces.</h1>
            <p class="page-subtitle">Gestion des annonces de la plateforme.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Nom</th>
                        <th>Catégorie</th>
                        <th>Prix</th>
                        <th>Statut</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="loading"><td colspan="5" class="empty">Chargement...</td></tr>
                    <tr v-else-if="annonces.length === 0"><td colspan="5" class="empty">Aucune annonce.</td></tr>
                    <tr v-for="a in annonces" :key="a.id">
                        <td class="td-bold">{{ a.name }}</td>
                        <td class="td-muted">{{ a.category }}</td>
                        <td>{{ a.price.toFixed(2) }} €</td>
                        <td>
                            <span class="badge" :class="a.approved ? 'badge--active' : 'badge--inactive'">
                                {{ a.approved ? 'Approuvée' : 'En attente' }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button v-if="!a.approved" class="btn-icon" title="Approuver" @click="approveAnnonce(a.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                            </button>
                            <button v-else class="btn-icon" title="Désapprouver" @click="disapproveAnnonce(a.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <line x1="18" y1="6" x2="6" y2="18"></line>
                                    <line x1="6" y1="6" x2="18" y2="18"></line>
                                </svg>
                            </button>
                            <router-link :to="`/admin/annonces/${a.id}`" class="btn-icon" title="Voir">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                                    <circle cx="12" cy="12" r="3" />
                                </svg>
                            </router-link>
                            <button class="btn-icon btn-icon--danger" title="Supprimer" @click="deleteAnnonce(a.id)">
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
.page-title { font-size: 2.6rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; border-bottom: 1px solid rgba(53,53,53,0.05); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; padding: 40px !important; opacity: 0.4; }
.badge { padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); transition: all 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
</style>
