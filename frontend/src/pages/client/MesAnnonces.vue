<script setup lang="ts">
import { onMounted } from 'vue'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()

const statusLabels: Record<string, string> = {
    active: 'Active',
    sold: 'Vendue',
    cancelled: 'Annulée',
}

const statusClass: Record<string, string> = {
    active: 'badge--active',
    sold: 'badge--sold',
    cancelled: 'badge--cancelled',
}

function formatPrice(price: any): string {
    if (!price) return '—'
    const val = typeof price === 'object' ? (price.Float64 ?? price.Int64) : price
    return `${Number(val).toFixed(2)} €`
}

function formatDate(ts: any): string {
    if (!ts) return '—'
    const date = new Date(ts.Time ?? ts)
    return date.toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

async function handleDelete(id: number) {
    await clientStore.deleteAnnonce(id)
}

onMounted(() => {
    clientStore.fetchAnnonces()
})
</script>

<template>
    <div class="page">
        <div class="page-header">
            <h1 class="page-title">Mes Annonces.</h1>
            <router-link to="/particulier/annonces/nouvelle" class="btn-primary">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19" />
                    <line x1="5" y1="12" x2="19" y2="12" />
                </svg>
                Nouvelle annonce
            </router-link>
        </div>

        <div v-if="clientStore.isLoading" class="state-empty">
            <p>Chargement…</p>
        </div>

        <div v-else-if="clientStore.annonces.length === 0" class="state-empty">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                    <polyline points="14 2 14 8 20 8" />
                </svg>
            </div>
            <p class="empty-title">Aucune annonce pour l'instant</p>
            <p class="empty-sub">
                Publiez votre première annonce pour proposer vos objets à la communauté.
            </p>
            <router-link
                to="/particulier/annonces/nouvelle"
                class="btn-primary btn-primary--centered"
            >
                Créer une annonce
            </router-link>
        </div>

        <div v-else class="annonces-list">
            <div
                v-for="annonce in clientStore.annonces"
                :key="annonce.id?.Int64"
                class="annonce-card"
            >
                <div class="annonce-main">
                    <div v-if="annonce.image_url" class="annonce-thumb">
                        <img :src="'http://localhost:8081' + annonce.image_url" alt="" />
                    </div>
                    <div class="annonce-info">
                        <div class="badge-row">
                            <span class="badge" :class="statusClass[annonce.status] ?? 'badge--active'">
                                {{ statusLabels[annonce.status] ?? annonce.status }}
                            </span>
                            <span v-if="!annonce.approved" class="badge badge--pending">
                                En attente
                            </span>
                        </div>
                        <h3 class="annonce-name">{{ annonce.name }}</h3>
                        <p class="annonce-desc">{{ annonce.description }}</p>
                    </div>
                    <div class="annonce-meta">
                        <span class="annonce-price">{{ formatPrice(annonce.price) }}</span>
                        <span class="annonce-date">{{ formatDate(annonce.created_at) }}</span>
                    </div>
                </div>
                <div class="annonce-actions">
                    <button class="btn-danger" @click="handleDelete(annonce.id?.Int64)">
                        Supprimer
                    </button>
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
    border: none;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
}
.btn-primary svg {
    width: 16px;
    height: 16px;
}
.btn-primary:hover {
    background: var(--green-mid);
}
.btn-primary--centered {
    margin-top: 16px;
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
    margin: 0;
    max-width: 360px;
    margin: 0 auto;
    line-height: 1.6;
}

.annonces-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.annonce-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 12px;
    padding: 20px 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
}
.annonce-main {
    display: flex;
    align-items: center;
    gap: 20px;
    flex: 1;
    min-width: 0;
}
.annonce-thumb {
    width: 64px;
    height: 64px;
    border-radius: 8px;
    background: var(--cream);
    overflow: hidden;
    flex-shrink: 0;
}
.annonce-thumb img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.annonce-info {
    flex: 1;
    min-width: 0;
}
.badge-row {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
}
.badge {
    display: inline-block;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
}
.badge--active {
    background: var(--green-pale);
    color: var(--green-dark);
}
.badge--sold {
    background: rgba(52, 137, 91, 0.12);
    color: var(--green-mid);
}
.badge--cancelled {
    background: rgba(53, 53, 53, 0.08);
    color: rgba(53, 53, 53, 0.55);
}
.badge--pending {
    background: rgba(246, 173, 85, 0.15);
    color: #c05621;
}
.annonce-name {
    font-size: 0.95rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.annonce-desc {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.6;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.annonce-meta {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
    flex-shrink: 0;
}
.annonce-price {
    font-size: 1rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
}
.annonce-date {
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.45;
}
.annonce-actions {
    flex-shrink: 0;
}
.btn-danger {
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.18);
    color: rgba(53, 53, 53, 0.55);
    padding: 7px 14px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        border-color 0.2s,
        color 0.2s,
        background 0.2s;
}
.btn-danger:hover {
    border-color: #e53e3e;
    color: #e53e3e;
    background: rgba(229, 62, 62, 0.06);
}

@media (max-width: 640px) {
    .annonce-card {
        flex-direction: column;
        align-items: flex-start;
    }
    .annonce-main {
        flex-direction: column;
        gap: 12px;
    }
    .annonce-meta {
        align-items: flex-start;
    }
}
</style>
