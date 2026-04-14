<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Commission {
    id: number
    type: string
    taux: number
    montant_total: number
    nb_transactions: number
    periode: string
}

const commissions = ref<Commission[]>([])
const total = ref(0)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/commissions', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            commissions.value = await res.json()
            total.value = commissions.value.reduce((s, c) => s + c.montant_total, 0)
        }
    } catch {}
})
</script>

<template>
    <div class="commissions">
        <div class="page-header">
            <h1 class="page-title">Commissions.</h1>
            <p class="page-subtitle">Suivi des commissions par type de transaction.</p>
        </div>

        <div class="total-card">
            <div class="total-label">Total des commissions perçues</div>
            <div class="total-value">{{ total.toLocaleString('fr-FR') }} €</div>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Type</th>
                        <th>Taux</th>
                        <th>Transactions</th>
                        <th>Période</th>
                        <th>Montant total</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="commissions.length === 0">
                        <td colspan="5" class="empty">Aucune commission.</td>
                    </tr>
                    <tr v-for="c in commissions" :key="c.id">
                        <td class="td-bold">{{ c.type }}</td>
                        <td>{{ c.taux }} %</td>
                        <td>{{ c.nb_transactions }}</td>
                        <td class="td-muted">{{ c.periode }}</td>
                        <td class="td-amount">{{ c.montant_total.toLocaleString('fr-FR') }} €</td>
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
.total-card { background: var(--green-dark); border-radius: 14px; padding: 24px 28px; margin-bottom: 24px; display: inline-flex; flex-direction: column; gap: 4px; }
.total-label { font-size: 0.8rem; font-weight: 600; color: rgba(255,255,255,0.65); text-transform: uppercase; letter-spacing: 0.06em; }
.total-value { font-size: 2.4rem; font-weight: 800; color: var(--white); letter-spacing: -0.04em; line-height: 1; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-amount { font-weight: 700; color: var(--green-dark); }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
</style>
