<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Pub {
    id: number
    titre: string
    annonceur: string
    type: string
    debut: string
    fin: string
    statut: string
    budget: number
}

const pubs = ref<Pub[]>([])

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/publicites', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) pubs.value = await res.json()
    } catch {}
})

async function toggleStatut(p: Pub) {
    const newStatut = p.statut === 'active' ? 'inactive' : 'active'
    await fetch(`http://localhost:8081/admin/publicites/${p.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
        body: JSON.stringify({ statut: newStatut }),
    })
    p.statut = newStatut
}

async function supprimer(id: number) {
    if (!confirm('Supprimer cette publicité ?')) return
    await fetch(`http://localhost:8081/admin/publicites/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    pubs.value = pubs.value.filter(p => p.id !== id)
}
</script>

<template>
    <div class="publicites">
        <div class="page-header">
            <h1 class="page-title">Publicités.</h1>
            <p class="page-subtitle">Gestion des publicités sur la plateforme.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Annonceur</th>
                        <th>Type</th>
                        <th>Période</th>
                        <th>Budget</th>
                        <th>Statut</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="pubs.length === 0">
                        <td colspan="7" class="empty">Aucune publicité.</td>
                    </tr>
                    <tr v-for="p in pubs" :key="p.id">
                        <td class="td-bold">{{ p.titre }}</td>
                        <td class="td-muted">{{ p.annonceur }}</td>
                        <td>{{ p.type }}</td>
                        <td class="td-muted">{{ p.debut }} → {{ p.fin }}</td>
                        <td>{{ p.budget.toFixed(0) }} €</td>
                        <td>
                            <button
                                class="badge badge-btn"
                                :class="p.statut === 'active' ? 'badge--active' : 'badge--inactive'"
                                @click="toggleStatut(p)"
                            >
                                {{ p.statut === 'active' ? 'Active' : 'Inactive' }}
                            </button>
                        </td>
                        <td>
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
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge-btn { border: none; cursor: pointer; transition: opacity 0.2s; }
.badge-btn:hover { opacity: 0.7; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
</style>
