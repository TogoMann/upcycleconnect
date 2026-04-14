<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Formation {
    id: number
    titre: string
    categorie: string
    duree: string
    statut: string
    inscrits: number
}

const formations = ref<Formation[]>([])

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/salarie/formations', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) formations.value = await res.json()
    } catch {}
})

async function supprimer(id: number) {
    if (!confirm('Supprimer cette formation ?')) return
    await fetch(`http://localhost:8081/salarie/formations/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    formations.value = formations.value.filter(f => f.id !== id)
}
</script>

<template>
    <div class="formations">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">Mes Formations.</h1>
                    <p class="page-subtitle">Modules de formation créés et gérés par vous.</p>
                </div>
                <router-link to="/salarie/formations/nouvelle" class="btn-primary">
                    + Nouvelle formation
                </router-link>
            </div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Catégorie</th>
                        <th>Durée</th>
                        <th>Inscrits</th>
                        <th>Statut</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="formations.length === 0">
                        <td colspan="6" class="empty">Aucune formation.</td>
                    </tr>
                    <tr v-for="f in formations" :key="f.id">
                        <td class="td-bold">{{ f.titre }}</td>
                        <td class="td-muted">{{ f.categorie }}</td>
                        <td>{{ f.duree }}</td>
                        <td>{{ f.inscrits }}</td>
                        <td>
                            <span class="badge" :class="f.statut === 'publiee' ? 'badge--active' : 'badge--draft'">
                                {{ f.statut === 'publiee' ? 'Publiée' : 'Brouillon' }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <router-link :to="`/salarie/formations/${f.id}/edit`" class="btn-icon" title="Modifier">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                                </svg>
                            </router-link>
                            <button class="btn-icon btn-icon--danger" title="Supprimer" @click="supprimer(f.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="3 6 5 6 21 6" />
                                    <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                                    <path d="M10 11v6M14 11v6" />
                                    <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" />
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
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border-radius: 8px; font-size: 0.88rem; font-weight: 600; text-decoration: none; transition: background 0.2s; white-space: nowrap; }
.btn-primary:hover { background: var(--green-mid); }
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
.btn-icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); text-decoration: none; transition: border-color 0.2s, color 0.2s; }
.btn-icon svg { width: 14px; height: 14px; }
.btn-icon:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-icon--danger:hover { border-color: #dc2626; color: #dc2626; }
</style>
