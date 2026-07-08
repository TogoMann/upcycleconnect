<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
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
        const res = await fetch(`${API_BASE}/admin/abonnements`, {
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

function planLabel(p: string): string {
    if (p === 'premium') return t('admin.abonnements.planPremium')
    if (p === 'pro') return t('admin.abonnements.planPro')
    return t('admin.abonnements.planFree')
}

function statutLabel(s: string): string {
    return s === 'actif' ? t('admin.abonnements.statusActive') : t('admin.abonnements.statusInactive')
}
</script>

<template>
    <div class="abonnements">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.abonnements.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.abonnements.subtitle') }}</p>
        </div>

        <div class="kpi-row">
            <div class="kpi-sm">
                <div class="kpi-sm-value">{{ abonnements.filter(a => a.plan === 'premium').length }}</div>
                <div class="kpi-sm-label">{{ t('admin.abonnements.kpiPremium') }}</div>
            </div>
            <div class="kpi-sm">
                <div class="kpi-sm-value">{{ abonnements.filter(a => a.plan === 'pro').length }}</div>
                <div class="kpi-sm-label">{{ t('admin.abonnements.kpiPro') }}</div>
            </div>
            <div class="kpi-sm">
                <div class="kpi-sm-value">{{ abonnements.filter(a => a.statut === 'actif').length }}</div>
                <div class="kpi-sm-label">{{ t('admin.abonnements.kpiActive') }}</div>
            </div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.abonnements.colCompany') }}</th>
                        <th>{{ t('admin.abonnements.colUser') }}</th>
                        <th>{{ t('admin.abonnements.colPlan') }}</th>
                        <th>{{ t('admin.abonnements.colPeriod') }}</th>
                        <th>{{ t('admin.abonnements.colAmount') }}</th>
                        <th>{{ t('admin.abonnements.colStatus') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="abonnements.length === 0">
                        <td colspan="6" class="empty">{{ t('admin.abonnements.empty') }}</td>
                    </tr>
                    <tr v-for="a in abonnements" :key="a.id">
                        <td class="td-bold">{{ a.entreprise }}</td>
                        <td class="td-muted">{{ a.utilisateur }}</td>
                        <td><span :class="planClass(a.plan)">{{ planLabel(a.plan) }}</span></td>
                        <td class="td-muted">{{ a.debut }} → {{ a.fin }}</td>
                        <td>{{ t('admin.abonnements.perMonth', { amount: a.montant.toFixed(2) }) }}</td>
                        <td>
                            <span class="badge" :class="a.statut === 'actif' ? 'badge--active' : 'badge--inactive'">
                                {{ statutLabel(a.statut) }}
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
