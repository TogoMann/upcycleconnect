<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { API_BASE } from '@/config'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()
const invoices = ref<any[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)

async function fetchInvoices() {
    isLoading.value = true
    error.value = null
    try {
        const res = await fetch(`${API_BASE}/financial/my-invoices`, {
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })
        if (!res.ok) throw new Error(t('pro.facturation.loadError'))
        invoices.value = await res.json()
    } catch (err: any) {
        error.value = err.message
    } finally {
        isLoading.value = false
    }
}

async function downloadPDF(id: number, number: string) {
    try {
        const res = await fetch(`${API_BASE}/financial/invoices/${id}/pdf`, {
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })
        if (!res.ok) throw new Error(t('pro.facturation.downloadError'))
        
        const blob = await res.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `Facture_${number}.pdf`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
    } catch (err: any) {
        alert(err.message)
    }
}

onMounted(fetchInvoices)
</script>

<template>
    <div class="facturation-page">
        <div class="page-header">
            <h1 class="page-title">{{ t('pro.facturation.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('pro.facturation.subtitle') }}</p>
        </div>

        <div v-if="isLoading" class="state-msg">{{ t('pro.facturation.loading') }}</div>
        <div v-else-if="error" class="state-msg error">{{ error }}</div>

        <div class="table-container" v-else>
            <table class="invoice-table">
                <thead>
                    <tr>
                        <th>{{ t('pro.facturation.invoiceNumber') }}</th>
                        <th>{{ t('pro.facturation.type') }}</th>
                        <th>{{ t('pro.facturation.role') }}</th>
                        <th>{{ t('pro.facturation.date') }}</th>
                        <th>{{ t('pro.facturation.amountExclTax') }}</th>
                        <th>{{ t('pro.facturation.totalInclTax') }}</th>
                        <th>{{ t('pro.facturation.action') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="inv in invoices" :key="inv.id">
                        <td class="font-bold">{{ inv.invoice_number }}</td>
                        <td class="capitalize">{{ inv.order_type }}</td>
                        <td>
                            <span :class="['badge-role', inv.user_id === authStore.user?.id ? 'purchase' : 'sale']">
                                {{ inv.user_id === authStore.user?.id ? t('pro.facturation.purchase') : t('pro.facturation.sale') }}
                            </span>
                        </td>
                        <td>{{ new Date(inv.created_at).toLocaleDateString() }}</td>
                        <td>{{ inv.amount.toFixed(2) }}€</td>
                        <td class="font-bold">{{ inv.total_amount.toFixed(2) }}€</td>
                        <td>
                            <button class="btn-download" @click="downloadPDF(inv.id, inv.invoice_number)">
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                                    <polyline points="7 10 12 15 17 10" />
                                    <line x1="12" y1="15" x2="12" y2="3" />
                                </svg>
                                PDF
                            </button>
                        </td>
                    </tr>
                    <tr v-if="invoices.length === 0">
                        <td colspan="8" class="empty-cell">{{ t('pro.facturation.noInvoices') }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.facturation-page {
    padding: 20px 0;
}

.page-header {
    margin-bottom: 32px;
}

.page-title {
    font-size: 2.2rem;
    font-weight: 800;
    color: var(--charcoal);
    margin: 0 0 8px;
    letter-spacing: -0.02em;
}

.page-subtitle {
    font-size: 1rem;
    color: var(--charcoal);
    opacity: 0.6;
    margin: 0;
}

.table-container {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    overflow: hidden;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.invoice-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
    font-size: 0.9rem;
}

.invoice-table th {
    background: #fafafa;
    padding: 16px 24px;
    font-weight: 700;
    color: var(--green-mid);
    text-transform: uppercase;
    font-size: 0.75rem;
    letter-spacing: 0.05em;
    border-bottom: 1px solid rgba(53, 53, 53, 0.05);
}

.invoice-table td {
    padding: 18px 24px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.03);
    color: var(--charcoal);
}

.font-bold { font-weight: 700; }
.capitalize { text-transform: capitalize; }

.badge {
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 700;
}

.paid {
    background: var(--green-pale);
    color: var(--green-dark);
}

.badge-role {
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
}
.badge-role.purchase {
    background: #ebf8ff;
    color: #2b6cb0;
}
.badge-role.sale {
    background: #f0fff4;
    color: #2f855a;
}

.btn-download {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: var(--green-pale);
    color: var(--green-dark);
    border: none;
    padding: 8px 14px;
    border-radius: 8px;
    font-weight: 700;
    cursor: pointer;
    transition: transform 0.2s;
}

.btn-download:hover {
    transform: translateY(-1px);
}

.state-msg {
    padding: 40px;
    text-align: center;
    color: rgba(53, 53, 53, 0.5);
}

.error { color: #c0392b; }

.empty-cell {
    padding: 60px !important;
    text-align: center;
    color: rgba(53, 53, 53, 0.4);
}
</style>
