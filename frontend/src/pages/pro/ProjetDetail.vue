<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

interface Projet {
    id: number
    nom: string
    description: string
    type: string
    statut: string
    date_debut: string
    budget: number
    materiaux: { nom: string; quantite: string; cout: number }[]
}

const projet = ref<Projet | null>(null)
const loading = ref(true)

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch(`http://localhost:8081/pro/projets/${route.params.id}`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) projet.value = await res.json()
        else router.push('/pro/projets')
    } catch {}
    loading.value = false
})

function badgeClass(s: string) {
    if (s === 'en_cours') return 'badge badge--active'
    if (s === 'termine') return 'badge badge--done'
    return 'badge badge--draft'
}
function badgeLabel(s: string) {
    if (s === 'en_cours') return 'En cours'
    if (s === 'termine') return 'Terminé'
    return 'Brouillon'
}
</script>

<template>
    <div class="projet-detail">
        <router-link to="/pro/projets" class="back-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15 18 9 12 15 6" />
            </svg>
            Retour aux projets
        </router-link>

        <div v-if="loading" class="loading">Chargement…</div>

        <template v-else-if="projet">
            <div class="page-header">
                <div class="header-row">
                    <h1 class="page-title">{{ projet.nom }}.</h1>
                    <span :class="badgeClass(projet.statut)">{{ badgeLabel(projet.statut) }}</span>
                </div>
                <p class="page-subtitle">{{ projet.description || 'Aucune description.' }}</p>
            </div>

            <div class="info-grid">
                <div class="info-card">
                    <div class="info-label">Type</div>
                    <div class="info-value">{{ projet.type }}</div>
                </div>
                <div class="info-card">
                    <div class="info-label">Début</div>
                    <div class="info-value">{{ projet.date_debut }}</div>
                </div>
                <div class="info-card">
                    <div class="info-label">Budget</div>
                    <div class="info-value">{{ projet.budget.toFixed(2) }} €</div>
                </div>
            </div>

            <div class="section" v-if="projet.materiaux?.length">
                <h2 class="section-title">Matériaux</h2>
                <div class="table-wrap">
                    <table class="data-table">
                        <thead>
                            <tr>
                                <th>Nom</th>
                                <th>Quantité</th>
                                <th>Coût</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="m in projet.materiaux" :key="m.nom">
                                <td class="td-bold">{{ m.nom }}</td>
                                <td>{{ m.quantite }}</td>
                                <td>{{ m.cout.toFixed(2) }} €</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </template>
    </div>
</template>

<style scoped>
.back-link { display: inline-flex; align-items: center; gap: 6px; font-size: 0.85rem; color: var(--green-mid); text-decoration: none; margin-bottom: 24px; transition: color 0.2s; }
.back-link:hover { color: var(--green-dark); }
.back-link svg { width: 16px; height: 16px; }
.page-header { margin-bottom: 28px; }
.header-row { display: flex; align-items: flex-start; gap: 16px; margin-bottom: 8px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.badge { display: inline-block; padding: 5px 12px; border-radius: 20px; font-size: 0.78rem; font-weight: 600; white-space: nowrap; align-self: center; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--done { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.badge--draft { background: #fef3c7; color: #92400e; }
.info-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 40px; }
.info-card { background: var(--green-pale); border-radius: 12px; padding: 20px; }
.info-label { font-size: 0.78rem; font-weight: 600; color: var(--green-dark); opacity: 0.7; text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 8px; }
.info-value { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); }
.section { margin-bottom: 32px; }
.section-title { font-size: 1.05rem; font-weight: 700; color: var(--charcoal); margin: 0 0 16px; }
.table-wrap { background: var(--white); border-radius: 12px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 12px 18px; font-size: 0.78rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 12px 18px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.td-bold { font-weight: 600; }
.loading { opacity: 0.5; font-size: 0.9rem; padding: 40px 0; }
@media (max-width: 700px) { .info-grid { grid-template-columns: 1fr; } }
</style>
