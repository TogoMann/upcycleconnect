<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Stats {
    ca_mois: number
    objets_recycles: number
    taux_conversion: number
    score_eco: number
    evolution: { mois: string; valeur: number }[]
}

const stats = ref<Stats | null>(null)
const locked = ref(false)
const maxValeur = computed(() => Math.max(...(stats.value?.evolution.map(x => x.valeur) ?? [1])))

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro/dashboard-avance', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.status === 403) { locked.value = true; return }
        if (res.ok) stats.value = await res.json()
    } catch {}
})
</script>

<template>
    <div class="tba">
        <div class="page-header">
            <h1 class="page-title">Tableau avancé.</h1>
            <p class="page-subtitle">Statistiques approfondies de votre activité.</p>
        </div>

        <div v-if="locked" class="locked-card">
            <div class="locked-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
            </div>
            <div class="locked-title">Fonctionnalité Premium</div>
            <div class="locked-desc">Le tableau de bord avancé est réservé aux abonnements Premium.</div>
            <router-link to="/pro/abonnements" class="btn-primary">Passer à Premium</router-link>
        </div>

        <template v-else-if="stats">
            <div class="kpi-grid">
                <div class="kpi-card">
                    <div class="kpi-label">CA ce mois</div>
                    <div class="kpi-value">{{ stats.ca_mois.toFixed(0) }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">Objets recyclés</div>
                    <div class="kpi-value">{{ stats.objets_recycles }}</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">Taux conversion</div>
                    <div class="kpi-value">{{ stats.taux_conversion }} %</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">Score éco</div>
                    <div class="kpi-value">{{ stats.score_eco }}</div>
                </div>
            </div>

            <div class="chart-section">
                <h2 class="section-title">Évolution mensuelle</h2>
                <div class="chart-bars">
                    <div
                        v-for="e in stats.evolution"
                        :key="e.mois"
                        class="chart-bar-wrap"
                    >
                        <div
                            class="chart-bar"
                            :style="{ height: Math.min((e.valeur / maxValeur) * 140, 140) + 'px' }"
                        ></div>
                        <span class="chart-label">{{ e.mois }}</span>
                    </div>
                </div>
            </div>
        </template>

        <div v-else class="loading">Chargement des statistiques…</div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.locked-card { max-width: 440px; background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 16px; padding: 48px 40px; display: flex; flex-direction: column; align-items: center; gap: 16px; text-align: center; }
.locked-icon { width: 64px; height: 64px; background: var(--green-pale); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: var(--green-mid); }
.locked-icon svg { width: 32px; height: 32px; }
.locked-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); }
.locked-desc { font-size: 0.9rem; color: var(--charcoal); opacity: 0.55; line-height: 1.5; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border-radius: 8px; font-size: 0.9rem; font-weight: 600; text-decoration: none; transition: background 0.2s; }
.btn-primary:hover { background: var(--green-mid); }
.kpi-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 40px; }
.kpi-card { background: var(--green-pale); border-radius: 14px; padding: 22px 20px; }
.kpi-label { font-size: 0.8rem; font-weight: 600; color: var(--green-dark); opacity: 0.7; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 8px; }
.kpi-value { font-size: 2rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.03em; line-height: 1; }
.section-title { font-size: 1.05rem; font-weight: 700; color: var(--charcoal); margin: 0 0 20px; }
.chart-section { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); padding: 28px; }
.chart-bars { display: flex; align-items: flex-end; gap: 12px; height: 160px; }
.chart-bar-wrap { display: flex; flex-direction: column; align-items: center; gap: 8px; flex: 1; }
.chart-bar { width: 100%; background: var(--green-mid); border-radius: 6px 6px 0 0; transition: height 0.3s; min-height: 4px; }
.chart-label { font-size: 0.72rem; color: var(--charcoal); opacity: 0.5; }
.loading { opacity: 0.5; font-size: 0.9rem; padding: 40px 0; }
@media (max-width: 700px) { .kpi-grid { grid-template-columns: repeat(2, 1fr); } }
</style>
