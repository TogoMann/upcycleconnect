<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Abonnement {
    id: number
    entreprise: string
    utilisateur: string
    plan: string
    debut: string
    fin: string
    statut: string
    montant: number
}

const abonnements = ref<Abonnement[]>([])

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/abonnements', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) abonnements.value = await res.json()
    } catch {}
})

function planClass(p: string) {
    if (p === 'premium') return 'badge badge--premium'
    if (p === 'pro') return 'badge badge--pro'
    return 'badge badge--free'
}
</script>

<template>
    <div class="abonnements">
        <div class="page-header">
            <h1 class="page-title">Abonnements.</h1>
            <p class="page-subtitle">Tous les abonnements professionnels actifs.</p>
        </div>

        <div class="kpi-row">
            <div class="kpi-sm">
                <div class="kpi-sm-value">{{ abonnements.filter(a => a.plan === 'premium').length }}</div>
                <div class="kpi-sm-label">Premium</div>
            </div>
            <div class="kpi-sm">
                <div class="kpi-sm-value">{{ abonnements.filter(a => a.plan === 'pro').length }}</div>
                <div class="kpi-sm-label">Pro</div>
            </div>
            <div class="kpi-sm">
                <div class="kpi-sm-value">{{ abonnements.filter(a => a.statut === 'actif').length }}</div>
                <div class="kpi-sm-label">Actifs</div>
            </div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Entreprise</th>
                        <th>Utilisateur</th>
                        <th>Plan</th>
                        <th>Période</th>
                        <th>Montant</th>
                        <th>Statut</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="abonnements.length === 0">
                        <td colspan="6" class="empty">Aucun abonnement.</td>
                    </tr>
                    <tr v-for="a in abonnements" :key="a.id">
                        <td class="td-bold">{{ a.entreprise }}</td>
                        <td class="td-muted">{{ a.utilisateur }}</td>
                        <td><span :class="planClass(a.plan)">{{ a.plan }}</span></td>
                        <td class="td-muted">{{ a.debut }} → {{ a.fin }}</td>
                        <td>{{ a.montant.toFixed(2) }} €/mois</td>
                        <td>
                            <span class="badge" :class="a.statut === 'actif' ? 'badge--active' : 'badge--inactive'">
                                {{ a.statut }}
                            </span>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.kpi-row { display: flex; gap: 14px; margin-bottom: 24px; }
.kpi-sm { background: var(--green-pale); border-radius: 12px; padding: 16px 24px; }
.kpi-sm-value { font-size: 1.8rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.03em; line-height: 1; }
.kpi-sm-label { font-size: 0.78rem; font-weight: 500; color: var(--green-dark); opacity: 0.65; margin-top: 4px; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.badge--premium { background: #fef3c7; color: #92400e; }
.badge--pro { background: #dbeafe; color: #1e40af; }
.badge--free { background: rgba(53,53,53,0.08); color: var(--charcoal); }
</style>
