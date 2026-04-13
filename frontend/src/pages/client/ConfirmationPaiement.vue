<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const itemName = computed(() => (route.query.name as string) || 'votre article')
const itemPrice = computed(() => Number(route.query.price) || 0)
</script>

<template>
    <div class="page">
        <div class="confirmation-card">
            <div class="success-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="20 6 9 17 4 12" />
                </svg>
            </div>

            <h1 class="confirm-title">Paiement confirmé !</h1>
            <p class="confirm-subtitle">
                Votre commande pour <strong>{{ itemName }}</strong> a bien été enregistrée.
            </p>

            <div class="receipt">
                <div class="receipt-row">
                    <span class="receipt-label">Article</span>
                    <span class="receipt-value">{{ itemName }}</span>
                </div>
                <div class="receipt-row">
                    <span class="receipt-label">Montant</span>
                    <span class="receipt-value receipt-value--price">{{ itemPrice.toFixed(2) }} €</span>
                </div>
                <div class="receipt-row">
                    <span class="receipt-label">Statut</span>
                    <span class="receipt-status">Payé</span>
                </div>
            </div>

            <div class="confirm-actions">
                <router-link to="/particulier/catalogue" class="btn-secondary">
                    Retour au catalogue
                </router-link>
                <router-link to="/particulier" class="btn-primary">
                    Tableau de bord
                </router-link>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding: 40px 0;
}
.confirmation-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 18px;
    padding: 48px 40px;
    text-align: center;
    max-width: 480px;
    width: 100%;
}
.success-icon {
    width: 64px;
    height: 64px;
    background: var(--green-dark);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
    color: var(--white);
}
.success-icon svg {
    width: 28px;
    height: 28px;
}
.confirm-title {
    font-size: 1.8rem;
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 12px;
    line-height: 1.1;
}
.confirm-subtitle {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.65;
    line-height: 1.6;
    margin: 0 0 32px;
}

.receipt {
    background: var(--cream);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 28px;
    text-align: left;
}
.receipt-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid rgba(53, 53, 53, 0.07);
}
.receipt-row:last-child {
    border-bottom: none;
}
.receipt-label {
    font-size: 0.82rem;
    color: var(--charcoal);
    opacity: 0.55;
    font-weight: 500;
}
.receipt-value {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--charcoal);
}
.receipt-value--price {
    font-size: 1rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
}
.receipt-status {
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    background: var(--green-pale);
    color: var(--green-dark);
    padding: 3px 10px;
    border-radius: 20px;
}

.confirm-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    flex-wrap: wrap;
}
.btn-primary {
    padding: 11px 22px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    transition: background 0.2s;
}
.btn-primary:hover {
    background: var(--green-mid);
}
.btn-secondary {
    padding: 11px 22px;
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    color: rgba(53, 53, 53, 0.65);
    text-decoration: none;
    transition: border-color 0.2s, color 0.2s;
}
.btn-secondary:hover {
    border-color: rgba(53, 53, 53, 0.4);
    color: var(--charcoal);
}
</style>
