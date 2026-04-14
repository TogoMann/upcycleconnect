<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const stats = ref({ annonces: 0, projets: 0, vues: 0, score: 0 })

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) {
            const data = await res.json()
            stats.value = { ...stats.value, ...data }
        }
    } catch {}
})
</script>

<template>
    <div class="dashboard">
        <div class="page-header">
            <h1 class="page-title">
                Bonjour, {{ authStore.user?.first_name || authStore.user?.username || 'vous' }}.
            </h1>
            <p class="page-subtitle">Vue d'ensemble de votre espace professionnel.</p>
        </div>

        <div class="kpi-grid">
            <router-link to="/pro/annonces" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="3" width="7" height="7" />
                        <rect x="14" y="3" width="7" height="7" />
                        <rect x="14" y="14" width="7" height="7" />
                        <rect x="3" y="14" width="7" height="7" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.annonces }}</div>
                <div class="kpi-label">Annonces actives</div>
            </router-link>

            <router-link to="/pro/projets" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.projets }}</div>
                <div class="kpi-label">Projets en cours</div>
            </router-link>

            <router-link to="/pro/publicites" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.vues }}</div>
                <div class="kpi-label">Vues ce mois</div>
            </router-link>

            <router-link to="/pro/abonnements" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.score }}</div>
                <div class="kpi-label">Score upcycling</div>
            </router-link>
        </div>

        <div class="actions-section">
            <h2 class="section-title">Actions rapides</h2>
            <div class="actions-grid">
                <router-link to="/pro/projets/nouveau" class="action-card">
                    <div class="action-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <circle cx="12" cy="12" r="10" />
                            <line x1="12" y1="8" x2="12" y2="16" />
                            <line x1="8" y1="12" x2="16" y2="12" />
                        </svg>
                    </div>
                    <span class="action-label">Nouveau projet</span>
                    <span class="action-desc">Créer un projet upcycling</span>
                </router-link>

                <router-link to="/pro/conteneurs/recuperer" class="action-card">
                    <div class="action-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                        </svg>
                    </div>
                    <span class="action-label">Récupérer un objet</span>
                    <span class="action-desc">Scanner un code-barres</span>
                </router-link>

                <router-link to="/pro/facturation" class="action-card">
                    <div class="action-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                            <polyline points="14 2 14 8 20 8" />
                        </svg>
                    </div>
                    <span class="action-label">Mes factures</span>
                    <span class="action-desc">Télécharger vos PDF</span>
                </router-link>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.kpi-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 48px; }
.kpi-card { background: var(--green-pale); border-radius: 14px; padding: 22px 20px; display: flex; flex-direction: column; gap: 8px; }
.kpi-card--link { text-decoration: none; transition: transform 0.2s, box-shadow 0.2s; }
.kpi-card--link:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(8, 106, 53, 0.12); }
.kpi-icon { width: 38px; height: 38px; background: var(--green-mid); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: var(--white); margin-bottom: 4px; }
.kpi-icon svg { width: 18px; height: 18px; }
.kpi-value { font-size: 2rem; font-weight: 800; color: var(--green-dark); line-height: 1; letter-spacing: -0.03em; }
.kpi-label { font-size: 0.82rem; color: var(--green-dark); font-weight: 500; opacity: 0.75; }
.section-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 16px; letter-spacing: -0.01em; }
.actions-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.action-card { background: var(--white); border: 1.5px solid rgba(53, 53, 53, 0.1); border-radius: 12px; padding: 20px; text-decoration: none; display: flex; flex-direction: column; gap: 8px; transition: border-color 0.2s, transform 0.2s; }
.action-card:hover { border-color: var(--green-mid); transform: translateY(-2px); }
.action-icon { width: 36px; height: 36px; background: var(--green-pale); border-radius: 8px; display: flex; align-items: center; justify-content: center; color: var(--green-dark); margin-bottom: 4px; }
.action-icon svg { width: 18px; height: 18px; }
.action-label { font-size: 0.9rem; font-weight: 700; color: var(--charcoal); line-height: 1.3; }
.action-desc { font-size: 0.78rem; color: var(--charcoal); opacity: 0.55; line-height: 1.4; }
@media (max-width: 900px) { .kpi-grid { grid-template-columns: repeat(2, 1fr); } .actions-grid { grid-template-columns: 1fr; } }
@media (max-width: 480px) { .kpi-grid { grid-template-columns: 1fr; } }
</style>
