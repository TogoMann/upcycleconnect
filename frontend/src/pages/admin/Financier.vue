<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface FinancierData {
    ca_total: number
    ca_mois: number
    charges: number
    marge: number
    evolution: { mois: string; ca: number; charges: number }[]
}

const data = ref<FinancierData | null>(null)
const maxCa = computed(() => Math.max(...(data.value?.evolution.map(x => x.ca) ?? [1])))

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/financier', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) data.value = await res.json()
    } catch {}
})
</script>

<template>
    <div class="financier">
        <div class="page-header">
            <h1 class="page-title">Financier.</h1>
            <p class="page-subtitle">Revenus, chiffre d'affaires et charges.</p>
        </div>

        <div v-if="data">
            <div class="kpi-grid">
                <div class="kpi-card">
                    <div class="kpi-label">CA Total</div>
                    <div class="kpi-value">{{ data.ca_total.toLocaleString('fr-FR') }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">CA ce mois</div>
                    <div class="kpi-value">{{ data.ca_mois.toLocaleString('fr-FR') }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">Charges</div>
                    <div class="kpi-value kpi-value--warn">{{ data.charges.toLocaleString('fr-FR') }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">Marge nette</div>
                    <div class="kpi-value">{{ data.marge.toLocaleString('fr-FR') }} €</div>
                </div>
            </div>

            <div class="chart-card">
                <h3 class="section-title">Évolution CA vs Charges</h3>
                <div class="chart-legend">
                    <span class="legend-item legend-ca">CA</span>
                    <span class="legend-item legend-charges">Charges</span>
                </div>
                <div class="chart-bars">
                    <div v-for="e in data.evolution" :key="e.mois" class="bar-group">
                        <div class="bars-pair">
                            <div
                                class="bar bar--ca"
                                :style="{ height: Math.min((e.ca / maxCa) * 120, 120) + 'px' }"
                            ></div>
                            <div
                                class="bar bar--charges"
                                :style="{ height: Math.min((e.charges / maxCa) * 120, 120) + 'px' }"
                            ></div>
                        </div>
                        <span class="chart-label">{{ e.mois }}</span>
                    </div>
                </div>
            </div>
        </div>
        <div v-else class="loading">Chargement…</div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.kpi-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 28px; }
.kpi-card { background: var(--green-pale); border-radius: 14px; padding: 22px 20px; }
.kpi-label { font-size: 0.8rem; font-weight: 600; color: var(--green-dark); opacity: 0.7; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 8px; }
.kpi-value { font-size: 1.8rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.03em; line-height: 1; }
.kpi-value--warn { color: #dc2626; }
.chart-card { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); padding: 28px; }
.section-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 16px; }
.chart-legend { display: flex; gap: 16px; margin-bottom: 20px; }
.legend-item { font-size: 0.8rem; font-weight: 600; display: flex; align-items: center; gap: 6px; }
.legend-item::before { content: ''; width: 12px; height: 12px; border-radius: 3px; display: inline-block; }
.legend-ca::before { background: var(--green-mid); }
.legend-charges::before { background: #fca5a5; }
.chart-bars { display: flex; align-items: flex-end; gap: 16px; height: 140px; }
.bar-group { display: flex; flex-direction: column; align-items: center; gap: 8px; flex: 1; }
.bars-pair { display: flex; align-items: flex-end; gap: 3px; width: 100%; justify-content: center; }
.bar { width: 14px; border-radius: 4px 4px 0 0; min-height: 2px; transition: height 0.3s; }
.bar--ca { background: var(--green-mid); }
.bar--charges { background: #fca5a5; }
.chart-label { font-size: 0.7rem; color: var(--charcoal); opacity: 0.5; }
.loading { opacity: 0.5; font-size: 0.9rem; padding: 40px 0; }
@media (max-width: 700px) { .kpi-grid { grid-template-columns: repeat(2, 1fr); } }
</style>
