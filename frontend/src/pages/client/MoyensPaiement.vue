<script setup lang="ts">
import { reactive, ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { API_BASE } from '@/config'

const { t } = useI18n()
const authStore = useAuthStore()

interface PaymentMethod {
    id: { Int64: number; Valid: boolean } | number
    card_brand: string
    card_last4: string
    card_exp_month: number
    card_exp_year: number
    is_default: boolean
}

const methods = ref<PaymentMethod[]>([])
const loading = ref(true)
const error = ref('')
const success = ref('')

const form = reactive({
    cardNumber: '',
    cardName: '',
    expiry: '',
    cvc: '',
    brand: 'Visa',
})

const submitting = ref(false)

async function fetchMethods() {
    loading.value = true
    try {
        const res = await fetch(`${API_BASE}/payment-methods/me`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const data = await res.json()
            methods.value = Array.isArray(data) ? data : []
        }
    } catch (e: any) {
        error.value = e.message
    } finally {
        loading.value = false
    }
}

function getMethodId(pm: PaymentMethod): number {
    if (typeof pm.id === 'object' && pm.id !== null) {
        return pm.id.Int64
    }
    return Number(pm.id)
}

async function handleSetDefault(pm: PaymentMethod) {
    const id = getMethodId(pm)
    error.value = ''
    success.value = ''
    try {
        const res = await fetch(`${API_BASE}/payment-methods/${id}/default`, {
            method: 'PATCH',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            success.value = 'Moyen de paiement par défaut mis à jour.'
            await fetchMethods()
        } else {
            error.value = 'Impossible de définir cette carte par défaut.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
}

async function handleDelete(pm: PaymentMethod) {
    const id = getMethodId(pm)
    error.value = ''
    success.value = ''
    try {
        const res = await fetch(`${API_BASE}/payment-methods/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            success.value = 'Moyen de paiement supprimé.'
            await fetchMethods()
        } else {
            error.value = 'Impossible de supprimer cette carte.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
}

function formatCardNumber(val: string): string {
    return val
        .replace(/\D/g, '')
        .slice(0, 16)
        .replace(/(.{4})/g, '$1 ')
        .trim()
}

function formatExpiry(val: string): string {
    const clean = val.replace(/\D/g, '').slice(0, 4)
    if (clean.length > 2) return clean.slice(0, 2) + '/' + clean.slice(2)
    return clean
}

function onCardNumberInput(e: Event) {
    const input = e.target as HTMLInputElement
    form.cardNumber = formatCardNumber(input.value)
    input.value = form.cardNumber
}

function onExpiryInput(e: Event) {
    const input = e.target as HTMLInputElement
    form.expiry = formatExpiry(input.value)
    input.value = form.expiry
}

async function handleAddCard() {
    const rawCard = form.cardNumber.replace(/\s/g, '')
    if (rawCard.length !== 16) {
        error.value = 'Numéro de carte invalide.'
        return
    }
    if (!/^\d{2}\/\d{2}$/.test(form.expiry)) {
        error.value = 'Date d\'expiration invalide (MM/AA).'
        return
    }
    if (!/^\d{3,4}$/.test(form.cvc)) {
        error.value = 'Code CVC invalide.'
        return
    }

    submitting.value = true
    error.value = ''
    success.value = ''

    const parts = form.expiry.split('/')
    const month = parseInt(parts[0], 10)
    const year = 2000 + parseInt(parts[1], 10)

    try {
        const body = {
            stripe_payment_method_id: 'pm_mock_' + Math.random().toString(36).substring(2, 12),
            card_last4: rawCard.slice(-4),
            card_brand: form.brand,
            card_exp_month: month,
            card_exp_year: year,
            is_default: methods.value.length === 0,
        }

        const res = await fetch(`${API_BASE}/payment-methods`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify(body),
        })

        if (res.ok) {
            success.value = 'Moyen de paiement ajouté avec succès.'
            form.cardNumber = ''
            form.cardName = ''
            form.expiry = ''
            form.cvc = ''
            await fetchMethods()
        } else {
            error.value = 'Impossible d\'ajouter la carte.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    fetchMethods()
})
</script>

<template>
    <div class="payment-methods-page">
        <div class="page-header">
            <h1 class="page-title">Moyens de paiement</h1>
            <p class="page-subtitle">Gérez vos cartes enregistrées pour vos règlements et abonnements.</p>
        </div>

        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="methods-layout">
            <div class="methods-list-section">
                <h2>Cartes enregistrées</h2>

                <div v-if="loading" class="loading-state">
                    Chargement des moyens de paiement...
                </div>

                <div v-else-if="methods.length === 0" class="empty-state">
                    <div class="empty-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <rect x="1" y="4" width="22" height="16" rx="2" />
                            <line x1="1" y1="10" x2="23" y2="10" />
                        </svg>
                    </div>
                    <p class="empty-text">Aucun moyen de paiement enregistré.</p>
                </div>

                <div v-else class="cards-list">
                    <div v-for="pm in methods" :key="getMethodId(pm)" class="card-item" :class="{ 'card-item--default': pm.is_default }">
                        <div class="card-logo-wrap">
                            <span class="card-brand">{{ pm.card_brand }}</span>
                        </div>
                        <div class="card-details">
                            <div class="card-number">•••• •••• •••• {{ pm.card_last4 }}</div>
                            <div class="card-expiry">Expire le {{ String(pm.card_exp_month).padStart(2, '0') }}/{{ String(pm.card_exp_year).slice(-2) }}</div>
                        </div>
                        <div class="card-actions">
                            <span v-if="pm.is_default" class="badge-default">Par défaut</span>
                            <button v-else class="btn-action-link" @click="handleSetDefault(pm)">Définir par défaut</button>
                            <button class="btn-delete-card" @click="handleDelete(pm)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="3 6 5 6 21 6"></polyline>
                                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div class="add-card-section">
                <h2>Ajouter une carte</h2>
                <form class="add-card-form" @submit.prevent="handleAddCard">
                    <div class="form-group">
                        <label class="form-label">Type de carte</label>
                        <select v-model="form.brand" class="form-input">
                            <option value="Visa">Visa</option>
                            <option value="Mastercard">Mastercard</option>
                            <option value="American Express">American Express</option>
                        </select>
                    </div>

                    <div class="form-group">
                        <label class="form-label">Numéro de carte</label>
                        <input
                            :value="form.cardNumber"
                            type="text"
                            inputmode="numeric"
                            class="form-input"
                            placeholder="1234 5678 9012 3456"
                            maxlength="19"
                            @input="onCardNumberInput"
                            required
                        />
                    </div>

                    <div class="form-group">
                        <label class="form-label">Nom du titulaire</label>
                        <input
                            v-model="form.cardName"
                            type="text"
                            class="form-input"
                            placeholder="JEAN DUPONT"
                            required
                        />
                    </div>

                    <div class="form-row">
                        <div class="form-group">
                            <label class="form-label">Date d'expiration</label>
                            <input
                                :value="form.expiry"
                                type="text"
                                inputmode="numeric"
                                class="form-input"
                                placeholder="MM/AA"
                                maxlength="5"
                                @input="onExpiryInput"
                                required
                            />
                        </div>

                        <div class="form-group">
                            <label class="form-label">Code CVC</label>
                            <input
                                v-model="form.cvc"
                                type="text"
                                inputmode="numeric"
                                class="form-input"
                                placeholder="123"
                                maxlength="4"
                                required
                            />
                        </div>
                    </div>

                    <button type="submit" class="btn-submit" :disabled="submitting">
                        {{ submitting ? 'Enregistrement...' : 'Enregistrer la carte' }}
                    </button>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
.payment-methods-page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.page-header {
    margin-bottom: 32px;
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 8px;
    line-height: 1.08;
}
.page-subtitle {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.6;
    margin: 0;
}

.alert {
    padding: 14px 20px;
    border-radius: 10px;
    font-size: 0.9rem;
    font-weight: 600;
    margin-bottom: 24px;
}
.alert-error {
    background: rgba(229, 62, 62, 0.08);
    border: 1.5px solid rgba(229, 62, 62, 0.25);
    color: #e53e3e;
}
.alert-success {
    background: var(--green-pale);
    border: 1.5px solid rgba(8, 106, 53, 0.25);
    color: var(--green-dark);
}

.methods-layout {
    display: grid;
    grid-template-columns: 1.2fr 0.8fr;
    gap: 40px;
    align-items: flex-start;
}

.methods-list-section h2,
.add-card-section h2 {
    font-size: 1.15rem;
    font-weight: 700;
    margin: 0 0 20px;
    letter-spacing: -0.01em;
}

.loading-state {
    text-align: center;
    padding: 40px 0;
    color: var(--charcoal);
    opacity: 0.5;
}

.empty-state {
    text-align: center;
    padding: 48px 32px;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.08);
    border-radius: 14px;
}
.empty-icon {
    width: 48px;
    height: 48px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    color: var(--green-mid);
}
.empty-icon svg {
    width: 22px;
    height: 22px;
}
.empty-text {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.5;
    margin: 0;
}

.cards-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.card-item {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 12px;
    padding: 20px 24px;
    display: flex;
    align-items: center;
    gap: 20px;
    transition: border-color 0.2s;
}
.card-item--default {
    border-color: var(--green-mid);
    background: rgba(215, 236, 225, 0.15);
}

.card-logo-wrap {
    width: 60px;
    height: 40px;
    background: var(--cream);
    border: 1px solid rgba(53, 53, 53, 0.15);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    font-size: 0.72rem;
    letter-spacing: 0.04em;
    color: var(--charcoal);
}

.card-details {
    flex: 1;
}
.card-number {
    font-size: 0.95rem;
    font-weight: 700;
    letter-spacing: 0.03em;
    margin-bottom: 2px;
}
.card-expiry {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.55;
}

.card-actions {
    display: flex;
    align-items: center;
    gap: 16px;
}

.badge-default {
    background: var(--green-mid);
    color: var(--white);
    font-size: 0.75rem;
    font-weight: 700;
    padding: 4px 10px;
    border-radius: 20px;
}

.btn-action-link {
    background: none;
    border: none;
    color: var(--green-mid);
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: color 0.2s;
}
.btn-action-link:hover {
    color: var(--green-dark);
    text-decoration: underline;
}

.btn-delete-card {
    background: none;
    border: none;
    color: rgba(53, 53, 53, 0.4);
    cursor: pointer;
    padding: 6px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    transition: background 0.2s, color 0.2s;
}
.btn-delete-card:hover {
    background: rgba(229, 62, 62, 0.08);
    color: #e53e3e;
}
.btn-delete-card svg {
    width: 16px;
    height: 16px;
}

.add-card-section {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    padding: 24px;
}

.add-card-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
}

.form-label {
    font-size: 0.82rem;
    font-weight: 600;
    color: var(--charcoal);
}

.form-input {
    width: 100%;
    padding: 10px 14px;
    font-size: 0.88rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
}
.form-input:focus {
    border-color: var(--green-mid);
    background: var(--white);
}

.btn-submit {
    width: 100%;
    padding: 12px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
    margin-top: 8px;
}
.btn-submit:hover:not(:disabled) {
    background: var(--green-mid);
}
.btn-submit:disabled {
    opacity: 0.65;
    cursor: not-allowed;
}

@media (max-width: 860px) {
    .methods-layout {
        grid-template-columns: 1fr;
    }
}
</style>
