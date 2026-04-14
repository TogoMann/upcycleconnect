<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const stats = ref({ formations: 0, creneaux: 0, articles: 0, threads: 0 })

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/salarie', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) stats.value = { ...stats.value, ...await res.json() }
    } catch {}
})
</script>

<template>
    <div class="dashboard">
        <div class="page-header">
            <h1 class="page-title">
                Bonjour, {{ authStore.user?.first_name || authStore.user?.username || 'vous' }}.
            </h1>
            <p class="page-subtitle">Vue d'ensemble de votre activité.</p>
        </div>

        <div class="kpi-grid">
            <router-link to="/salarie/formations" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
                        <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.formations }}</div>
                <div class="kpi-label">Formations créées</div>
            </router-link>

            <router-link to="/salarie/planning" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                        <line x1="16" y1="2" x2="16" y2="6" />
                        <line x1="8" y1="2" x2="8" y2="6" />
                        <line x1="3" y1="10" x2="21" y2="10" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.creneaux }}</div>
                <div class="kpi-label">Créneaux planifiés</div>
            </router-link>

            <router-link to="/salarie/conseils" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M12 20h9" />
                        <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.articles }}</div>
                <div class="kpi-label">Articles publiés</div>
            </router-link>

            <router-link to="/salarie/forum" class="kpi-card kpi-card--link">
                <div class="kpi-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                    </svg>
                </div>
                <div class="kpi-value">{{ stats.threads }}</div>
                <div class="kpi-label">Discussions actives</div>
            </router-link>
        </div>

        <div class="actions-section">
            <h2 class="section-title">Actions rapides</h2>
            <div class="actions-grid">
                <router-link to="/salarie/formations/nouvelle" class="action-card">
                    <div class="action-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <circle cx="12" cy="12" r="10" />
                            <line x1="12" y1="8" x2="12" y2="16" />
                            <line x1="8" y1="12" x2="16" y2="12" />
                        </svg>
                    </div>
                    <span class="action-label">Nouvelle formation</span>
                    <span class="action-desc">Créer un module de formation</span>
                </router-link>

                <router-link to="/salarie/conseils" class="action-card">
                    <div class="action-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M12 20h9" />
                            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
                        </svg>
                    </div>
                    <span class="action-label">Écrire un conseil</span>
                    <span class="action-desc">Publier un article</span>
                </router-link>

                <router-link to="/salarie/forum" class="action-card">
                    <div class="action-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                        </svg>
                    </div>
                    <span class="action-label">Modérer le forum</span>
                    <span class="action-desc">Gérer les discussions</span>
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
.kpi-card--link:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(8,106,53,0.12); }
.kpi-icon { width: 38px; height: 38px; background: var(--green-mid); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: var(--white); margin-bottom: 4px; }
.kpi-icon svg { width: 18px; height: 18px; }
.kpi-value { font-size: 2rem; font-weight: 800; color: var(--green-dark); line-height: 1; letter-spacing: -0.03em; }
.kpi-label { font-size: 0.82rem; color: var(--green-dark); font-weight: 500; opacity: 0.75; }
.section-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 16px; letter-spacing: -0.01em; }
.actions-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.action-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 12px; padding: 20px; text-decoration: none; display: flex; flex-direction: column; gap: 8px; transition: border-color 0.2s, transform 0.2s; }
.action-card:hover { border-color: var(--green-mid); transform: translateY(-2px); }
.action-icon { width: 36px; height: 36px; background: var(--green-pale); border-radius: 8px; display: flex; align-items: center; justify-content: center; color: var(--green-dark); margin-bottom: 4px; }
.action-icon svg { width: 18px; height: 18px; }
.action-label { font-size: 0.9rem; font-weight: 700; color: var(--charcoal); line-height: 1.3; }
.action-desc { font-size: 0.78rem; color: var(--charcoal); opacity: 0.55; line-height: 1.4; }
@media (max-width: 900px) { .kpi-grid { grid-template-columns: repeat(2, 1fr); } .actions-grid { grid-template-columns: 1fr; } }
@media (max-width: 480px) { .kpi-grid { grid-template-columns: 1fr; } }
</style>
