<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Projet {
    id: number
    nom: string
    type: string
    statut: string
    date_debut: string
    budget: number
}

const projets = ref<Projet[]>([])

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro/projets', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) projets.value = await res.json()
    } catch {}
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
    <div class="projets">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">Mes Projets.</h1>
                    <p class="page-subtitle">Gérez vos projets de valorisation.</p>
                </div>
                <router-link to="/pro/projets/nouveau" class="btn-primary">
                    + Nouveau projet
                </router-link>
            </div>
        </div>

        <div v-if="projets.length === 0" class="empty-state">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
                </svg>
            </div>
            <p class="empty-text">Aucun projet pour le moment.</p>
            <router-link to="/pro/projets/nouveau" class="btn-primary">Créer un projet</router-link>
        </div>

        <div v-else class="projets-grid">
            <router-link
                v-for="p in projets"
                :key="p.id"
                :to="`/pro/projets/${p.id}`"
                class="projet-card"
            >
                <div class="projet-header">
                    <span class="projet-nom">{{ p.nom }}</span>
                    <span :class="badgeClass(p.statut)">{{ badgeLabel(p.statut) }}</span>
                </div>
                <div class="projet-type">{{ p.type }}</div>
                <div class="projet-footer">
                    <span class="projet-date">{{ p.date_debut }}</span>
                    <span class="projet-budget">{{ p.budget.toFixed(0) }} €</span>
                </div>
            </router-link>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.header-row { display: flex; justify-content: space-between; align-items: flex-start; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { padding: 11px 20px; background: var(--green-dark); color: var(--white); border-radius: 8px; font-size: 0.88rem; font-weight: 600; text-decoration: none; transition: background 0.2s; white-space: nowrap; }
.btn-primary:hover { background: var(--green-mid); }
.empty-state { display: flex; flex-direction: column; align-items: center; gap: 16px; padding: 80px 0; }
.empty-icon { width: 64px; height: 64px; background: var(--green-pale); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: var(--green-mid); }
.empty-icon svg { width: 32px; height: 32px; }
.empty-text { font-size: 0.95rem; color: var(--charcoal); opacity: 0.5; margin: 0; }
.projets-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.projet-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 12px; padding: 20px; text-decoration: none; display: flex; flex-direction: column; gap: 10px; transition: border-color 0.2s, transform 0.2s; }
.projet-card:hover { border-color: var(--green-mid); transform: translateY(-2px); }
.projet-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 8px; }
.projet-nom { font-size: 0.95rem; font-weight: 700; color: var(--charcoal); line-height: 1.3; }
.projet-type { font-size: 0.82rem; color: var(--charcoal); opacity: 0.5; }
.projet-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 4px; }
.projet-date { font-size: 0.8rem; color: var(--charcoal); opacity: 0.5; }
.projet-budget { font-size: 0.9rem; font-weight: 700; color: var(--green-dark); }
.badge { display: inline-block; padding: 3px 8px; border-radius: 20px; font-size: 0.7rem; font-weight: 600; white-space: nowrap; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--done { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.badge--draft { background: #fef3c7; color: #92400e; }
@media (max-width: 800px) { .projets-grid { grid-template-columns: 1fr; } }
</style>
