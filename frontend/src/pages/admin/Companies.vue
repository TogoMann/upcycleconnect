<template>
    <div class="admin-page">
        <div class="section-header">
            <h2 class="section-title">{{ t('admin.companies.pageTitle') }}</h2>
            <button class="btn-refresh" @click="adminStore.fetchCompanies()">
                <svg
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <polyline points="23 4 23 10 17 10" />
                    <polyline points="1 20 1 14 7 14" />
                    <path
                        d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"
                    />
                </svg>
                {{ t('admin.companies.refresh') }}
            </button>
        </div>

        <div v-if="adminStore.isLoading" class="state-loading">
            {{ t('admin.companies.loading') }}
        </div>

        <div v-if="adminStore.error" class="state-error">
            {{ adminStore.error }}
        </div>

        <div class="table-wrap" v-if="!adminStore.isLoading">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.companies.colId') }}</th>
                        <th>{{ t('admin.companies.colName') }}</th>
                        <th>{{ t('admin.companies.colSiret') }}</th>
                        <th>{{ t('admin.companies.colAddress') }}</th>
                        <th>{{ t('admin.companies.colDateAdded') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="company in adminStore.companies" :key="company.id" class="table-row">
                        <td class="cell-id">{{ company.id }}</td>
                        <td class="cell-name">{{ company.name }}</td>
                        <td class="cell-siret"><code>{{ company.siret }}</code></td>
                        <td class="cell-address">{{ company.address }}</td>
                        <td class="cell-date">{{ formatDate(company.created_at) }}</td>
                    </tr>
                </tbody>
            </table>

            <div v-if="adminStore.companies.length === 0" class="empty-state">
                {{ t('admin.companies.empty') }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const adminStore = useAdminStore()

const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    })
}

onMounted(() => {
    adminStore.fetchCompanies()
})
</script>

<style scoped>
.admin-page {
    --green-dark: #086a35;
    --green-mid: #34895b;
    --green-pale: #d7ece1;
    --charcoal: #353535;
    --white: #ffffff;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

.section-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0;
}

.btn-refresh {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    background: var(--white);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
    cursor: pointer;
    transition: all 0.2s;
}

.btn-refresh:hover {
    background: var(--green-pale);
    border-color: var(--green-mid);
}

.table-wrap {
    background: var(--white);
    border-radius: 12px;
    border: 1px solid rgba(0, 0, 0, 0.08);
    overflow: hidden;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
}

.data-table th {
    background: #f8f9fa;
    padding: 14px 20px;
    font-size: 0.85rem;
    font-weight: 700;
    color: #666;
    text-transform: uppercase;
    letter-spacing: 0.03em;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.table-row {
    border-bottom: 1px solid rgba(0, 0, 0, 0.03);
    transition: background 0.15s;
}

.table-row:hover {
    background: #fcfdfc;
}

.data-table td {
    padding: 16px 20px;
    font-size: 0.95rem;
}

.cell-id {
    color: #999;
    font-weight: 600;
    width: 60px;
}

.cell-name {
    font-weight: 700;
    color: var(--green-dark);
}

.cell-siret code {
    background: #f1f1f1;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: monospace;
}

.cell-address {
    color: #666;
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.cell-date {
    color: #888;
}

.state-loading, .state-error, .empty-state {
    padding: 40px;
    text-align: center;
    color: #666;
}

.state-error {
    color: #e53e3e;
}
</style>
