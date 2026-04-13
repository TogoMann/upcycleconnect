<script setup lang="ts">
import { onMounted } from 'vue'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()

const statusLabels: Record<string, string> = {
    deposited: 'Déposé',
    validated: 'Validé',
    collected: 'Collecté',
}

const statusClass: Record<string, string> = {
    deposited: 'badge--deposited',
    validated: 'badge--validated',
    collected: 'badge--collected',
}

function formatDate(ts: any): string {
    if (!ts) return '—'
    const date = new Date(ts.Time ?? ts)
    return date.toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatWeight(w: any): string {
    if (!w) return '—'
    const val = typeof w === 'object' ? w.Float64 ?? w.Int64 : w
    return `${Number(val).toFixed(2)} kg`
}

onMounted(() => {
    clientStore.fetchDepots()
})
</script>

<template>
    <div class="page">
        <div class="page-header">
            <h1 class="page-title">Mes Dépôts.</h1>
            <router-link to="/particulier/conteneurs/deposer" class="btn-primary">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19" />
                    <line x1="5" y1="12" x2="19" y2="12" />
                </svg>
                Déposer un objet
            </router-link>
        </div>

        <div v-if="clientStore.isLoading" class="state-empty">
            <p>Chargement…</p>
        </div>

        <div v-else-if="clientStore.depots.length === 0" class="state-empty">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                </svg>
            </div>
            <p class="empty-title">Aucun objet déposé</p>
            <p class="empty-sub">Déposez vos premiers objets dans nos conteneurs partenaires pour contribuer à l'upcycling.</p>
            <router-link to="/particulier/conteneurs/deposer" class="btn-primary btn-primary--mt">
                Planifier un dépôt
            </router-link>
        </div>

        <div v-else>
            <div class="stats-row">
                <div class="stat-chip">
                    <span class="stat-value">{{ clientStore.depots.length }}</span>
                    <span class="stat-label">objet(s) total</span>
                </div>
                <div class="stat-chip">
                    <span class="stat-value">{{ clientStore.depots.filter((d: any) => d.status === 'validated').length }}</span>
                    <span class="stat-label">validé(s)</span>
                </div>
                <div class="stat-chip">
                    <span class="stat-value">{{ clientStore.depots.filter((d: any) => d.status === 'collected').length }}</span>
                    <span class="stat-label">collecté(s)</span>
                </div>
            </div>

            <div class="depots-list">
                <div v-for="depot in clientStore.depots" :key="depot.id?.Int64" class="depot-card">
                    <div class="depot-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                        </svg>
                    </div>
                    <div class="depot-info">
                        <div class="depot-top">
                            <span class="depot-material">{{ depot.material_type || 'Matériau non précisé' }}</span>
                            <span class="badge" :class="statusClass[depot.status] ?? 'badge--deposited'">
                                {{ statusLabels[depot.status] ?? depot.status }}
                            </span>
                        </div>
                        <div class="depot-details">
                            <span v-if="depot.physical_state" class="detail-item">État : {{ depot.physical_state }}</span>
                            <span v-if="depot.weight" class="detail-item">{{ formatWeight(depot.weight) }}</span>
                            <span class="detail-item">{{ formatDate(depot.created_at) }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 32px;
    flex-wrap: wrap;
    gap: 16px;
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0;
    line-height: 1.08;
}
.btn-primary {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: var(--green-dark);
    color: var(--white);
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    transition: background 0.2s;
}
.btn-primary svg {
    width: 16px;
    height: 16px;
}
.btn-primary:hover {
    background: var(--green-mid);
}
.btn-primary--mt {
    margin-top: 16px;
    display: inline-flex;
}

.state-empty {
    text-align: center;
    padding: 64px 32px;
}
.empty-icon {
    width: 64px;
    height: 64px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    color: var(--green-mid);
}
.empty-icon svg {
    width: 28px;
    height: 28px;
}
.empty-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 8px;
}
.empty-sub {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.6;
    max-width: 360px;
    margin: 0 auto;
    line-height: 1.6;
}

.stats-row {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
    flex-wrap: wrap;
}
.stat-chip {
    background: var(--green-pale);
    border-radius: 10px;
    padding: 12px 20px;
    display: flex;
    flex-direction: column;
    gap: 2px;
}
.stat-value {
    font-size: 1.5rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.03em;
    line-height: 1;
}
.stat-label {
    font-size: 0.75rem;
    color: var(--green-dark);
    opacity: 0.7;
    font-weight: 500;
}

.depots-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.depot-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 12px;
    padding: 16px 20px;
    display: flex;
    align-items: center;
    gap: 16px;
}
.depot-icon {
    width: 40px;
    height: 40px;
    background: var(--green-pale);
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--green-mid);
    flex-shrink: 0;
}
.depot-icon svg {
    width: 20px;
    height: 20px;
}
.depot-info {
    flex: 1;
    min-width: 0;
}
.depot-top {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 6px;
    flex-wrap: wrap;
}
.depot-material {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--charcoal);
}
.badge {
    display: inline-block;
    padding: 2px 9px;
    border-radius: 20px;
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
}
.badge--deposited {
    background: rgba(139, 189, 148, 0.25);
    color: var(--green-mid);
}
.badge--validated {
    background: var(--green-pale);
    color: var(--green-dark);
}
.badge--collected {
    background: rgba(8, 106, 53, 0.12);
    color: var(--green-dark);
}
.depot-details {
    display: flex;
    gap: 16px;
    flex-wrap: wrap;
}
.detail-item {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.55;
}
</style>
