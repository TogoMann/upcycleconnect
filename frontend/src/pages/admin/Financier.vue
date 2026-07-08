<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface FinancierData {
    ca_total: number
    ca_mois: number
    charges: number
    marge: number
    evolution: { mois: string; ca: number; charges: number }[]
}

interface Expense {
    id: number
    label: string
    amount: number
    category: string
    created_at: string
}

const data = ref<FinancierData | null>(null)
const expenses = ref<Expense[]>([])
const maxCa = computed(() => Math.max(...(data.value?.evolution.map(x => x.ca) ?? [1])))
const showExpenseForm = ref(false)
const expenseForm = ref({ label: '', amount: '', category: '' })
const savingExpense = ref(false)

async function fetchFinancier() {
    try {
        const res = await fetch(`${API_BASE}/admin/financier`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) data.value = await res.json()
    } catch {}
}

async function fetchExpenses() {
    try {
        const res = await fetch(`${API_BASE}/admin/expenses`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) expenses.value = await res.json()
    } catch {}
}

async function submitExpense() {
    if (!expenseForm.value.label || !expenseForm.value.amount) return
    savingExpense.value = true
    try {
        const res = await fetch(`${API_BASE}/admin/expenses`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({
                label: expenseForm.value.label,
                amount: parseFloat(expenseForm.value.amount),
                category: expenseForm.value.category,
            }),
        })
        if (res.ok) {
            expenseForm.value = { label: '', amount: '', category: '' }
            showExpenseForm.value = false
            await Promise.all([fetchExpenses(), fetchFinancier()])
        }
    } catch {}
    savingExpense.value = false
}

function fmtDate(iso: string): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

onMounted(() => {
    fetchFinancier()
    fetchExpenses()
})
</script>

<template>
    <div class="financier">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.financier.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.financier.subtitle') }}</p>
        </div>

        <div v-if="data">
            <div class="kpi-grid">
                <div class="kpi-card">
                    <div class="kpi-label">{{ t('admin.financier.caTotal') }}</div>
                    <div class="kpi-value">{{ data.ca_total.toLocaleString(locale === 'en' ? 'en-US' : 'fr-FR') }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">{{ t('admin.financier.caMonth') }}</div>
                    <div class="kpi-value">{{ data.ca_mois.toLocaleString(locale === 'en' ? 'en-US' : 'fr-FR') }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">{{ t('admin.financier.charges') }}</div>
                    <div class="kpi-value kpi-value--warn">{{ data.charges.toLocaleString(locale === 'en' ? 'en-US' : 'fr-FR') }} €</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-label">{{ t('admin.financier.netMargin') }}</div>
                    <div class="kpi-value">{{ data.marge.toLocaleString(locale === 'en' ? 'en-US' : 'fr-FR') }} €</div>
                </div>
            </div>

            <div class="chart-card">
                <h3 class="section-title">{{ t('admin.financier.evolutionTitle') }}</h3>
                <div class="chart-legend">
                    <span class="legend-item legend-ca">{{ t('admin.financier.legendCa') }}</span>
                    <span class="legend-item legend-charges">{{ t('admin.financier.legendCharges') }}</span>
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

            <div class="chart-card expenses-card">
                <div class="expenses-header">
                    <h3 class="section-title">{{ t('admin.financier.expensesTitle') }}</h3>
                    <button class="btn-add-expense" @click="showExpenseForm = !showExpenseForm">
                        {{ showExpenseForm ? t('admin.financier.cancel') : t('admin.financier.addExpense') }}
                    </button>
                </div>

                <form v-if="showExpenseForm" class="expense-form" @submit.prevent="submitExpense">
                    <input v-model="expenseForm.label" type="text" class="field-input" :placeholder="t('admin.financier.labelPlaceholder')" required />
                    <input v-model="expenseForm.amount" type="number" step="0.01" min="0.01" class="field-input" :placeholder="t('admin.financier.amountPlaceholder')" required />
                    <input v-model="expenseForm.category" type="text" class="field-input" :placeholder="t('admin.financier.categoryPlaceholder')" />
                    <button type="submit" class="btn-save-expense" :disabled="savingExpense">
                        {{ savingExpense ? t('admin.financier.saving') : t('admin.financier.save') }}
                    </button>
                </form>

                <table class="expenses-table">
                    <thead>
                        <tr>
                            <th>{{ t('admin.financier.colLabel') }}</th>
                            <th>{{ t('admin.financier.colCategory') }}</th>
                            <th>{{ t('admin.financier.colDate') }}</th>
                            <th>{{ t('admin.financier.colAmount') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="expenses.length === 0">
                            <td colspan="4" class="empty">{{ t('admin.financier.empty') }}</td>
                        </tr>
                        <tr v-for="e in expenses" :key="e.id">
                            <td class="td-bold">{{ e.label }}</td>
                            <td class="td-muted">{{ e.category || '—' }}</td>
                            <td class="td-muted">{{ fmtDate(e.created_at) }}</td>
                            <td>{{ e.amount.toFixed(2) }} €</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div v-else class="loading">{{ t('admin.financier.loading') }}</div>
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
.expenses-card { margin-top: 20px; }
.expenses-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.btn-add-expense { padding: 8px 16px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-add-expense:hover { background: var(--green-mid); }
.expense-form { display: flex; gap: 10px; margin-bottom: 20px; flex-wrap: wrap; }
.field-input { flex: 1; min-width: 140px; padding: 9px 14px; font-size: 0.88rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; }
.field-input:focus { border-color: var(--green-mid); }
.btn-save-expense { padding: 9px 18px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.85rem; font-weight: 600; cursor: pointer; }
.btn-save-expense:disabled { opacity: 0.5; }
.expenses-table { width: 100%; border-collapse: collapse; }
.expenses-table th { text-align: left; padding: 10px 12px; font-size: 0.75rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.05em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.expenses-table td { padding: 10px 12px; font-size: 0.85rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; }
.empty { text-align: center; opacity: 0.4; padding: 24px !important; }
@media (max-width: 700px) { .kpi-grid { grid-template-columns: repeat(2, 1fr); } }
</style>
