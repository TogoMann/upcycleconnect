<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'
import { useCartStore } from '@/stores/cart'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const clientStore = useClientStore()
const cartStore = useCartStore()

const itemName = computed(() => (route.query.name as string) || t('client.paiement.defaultItem'))
const itemPrice = computed(() => Number(route.query.price) || 0)
const itemId = computed(() => Number(route.query.id) || 0)
const paymentType = computed(() => (route.query.type as string) || '')
const isStripeCheckout = computed(() => paymentType.value === 'listing')
const redirecting = ref(false)
const stripeError = ref('')

async function handleStripeCheckout() {
    redirecting.value = true
    stripeError.value = ''
    try {
        const data = await clientStore.createOrderCheckout(itemId.value)
        if (data.free) {
            router.push({
                path: '/particulier/paiement/confirmation',
                query: { name: itemName.value, price: '0', type: 'listing', order_id: data.order_id },
            })
            return
        }
        if (data.url) window.location.href = data.url
    } catch (e: any) {
        stripeError.value = e.message
        redirecting.value = false
    }
}

const form = reactive({
    cardNumber: '',
    cardName: '',
    expiry: '',
    cvc: '',
})

const errors = reactive({
    cardNumber: '',
    cardName: '',
    expiry: '',
    cvc: '',
    global: '',
})

const submitting = ref(false)

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

function validate(): boolean {
    const raw = form.cardNumber.replace(/\s/g, '')
    errors.cardNumber = raw.length === 16 ? '' : t('client.paiement.errorCardNumber')
    errors.cardName = form.cardName.trim().length >= 2 ? '' : t('client.paiement.errorCardName')
    errors.expiry = /^\d{2}\/\d{2}$/.test(form.expiry) ? '' : t('client.paiement.errorExpiry')
    errors.cvc = /^\d{3,4}$/.test(form.cvc) ? '' : t('client.paiement.errorCvc')
    return !errors.cardNumber && !errors.cardName && !errors.expiry && !errors.cvc
}

async function handleSubmit() {
    if (!validate()) return
    submitting.value = true
    errors.global = ''
    try {
        if (paymentType.value === 'cart') {
            await cartStore.checkoutCart()
        } else if (paymentType.value === 'course') {
            await clientStore.createCourseOrder(itemId.value, itemPrice.value)
        } else if (paymentType.value === 'event') {
            await clientStore.createEventParticipation(itemId.value)
        }

        router.push({
            path: '/particulier/paiement/confirmation',
            query: { name: itemName.value, price: itemPrice.value },
        })
    } catch (e: any) {
        errors.global = e.message
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    const type = route.query.type as string
    if (type !== 'cart' && !itemId.value) {
        router.push('/particulier/catalogue')
    }
})
</script>

<template>
    <div class="page">
        <div class="page-header">
            <router-link to="/particulier/catalogue" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="19" y1="12" x2="5" y2="12" />
                    <polyline points="12 19 5 12 12 5" />
                </svg>
                {{ t('client.paiement.catalogue') }}
            </router-link>
            <h1 class="page-title">{{ t('client.paiement.pageTitle') }}</h1>
        </div>

        <div class="checkout-layout">
            <div v-if="isStripeCheckout" class="checkout-form checkout-form--stripe">
                <div class="form-section-title">{{ t('client.paiement.stripeSectionTitle') }}</div>

                <div class="stripe-badge">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                        <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
                        <line x1="1" y1="10" x2="23" y2="10" />
                    </svg>
                    {{ t('client.paiement.stripeRedirect') }}
                </div>

                <p class="stripe-explain">
                    {{ t('client.paiement.stripeExplain') }}
                </p>

                <div v-if="stripeError" class="error-banner">{{ stripeError }}</div>

                <button type="button" class="btn-pay" :disabled="redirecting" @click="handleStripeCheckout">
                    <svg v-if="!redirecting" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
                        <line x1="1" y1="10" x2="23" y2="10" />
                    </svg>
                    {{ redirecting ? t('client.paiement.redirecting') : t('client.paiement.pay', { price: itemPrice.toFixed(2) }) }}
                </button>
            </div>

            <form v-else class="checkout-form" @submit.prevent="handleSubmit">
                <div class="form-section-title">{{ t('client.paiement.paymentInfoTitle') }}</div>

                <div class="stripe-badge">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                        <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
                        <line x1="1" y1="10" x2="23" y2="10" />
                    </svg>
                    {{ t('client.paiement.securePayment') }}
                </div>

                <div class="form-group">
                    <label class="form-label">{{ t('client.paiement.cardNumber') }}</label>
                    <input
                        :value="form.cardNumber"
                        type="text"
                        inputmode="numeric"
                        class="form-input form-input--card"
                        :class="{ 'form-input--error': errors.cardNumber }"
                        placeholder="1234 5678 9012 3456"
                        maxlength="19"
                        @input="onCardNumberInput"
                    />
                    <span v-if="errors.cardNumber" class="form-error">{{ errors.cardNumber }}</span>
                </div>

                <div class="form-group">
                    <label class="form-label">{{ t('client.paiement.cardName') }}</label>
                    <input
                        v-model="form.cardName"
                        type="text"
                        class="form-input"
                        :class="{ 'form-input--error': errors.cardName }"
                        placeholder="JEAN DUPONT"
                        autocomplete="cc-name"
                    />
                    <span v-if="errors.cardName" class="form-error">{{ errors.cardName }}</span>
                </div>

                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">{{ t('client.paiement.expiryDate') }}</label>
                        <input
                            :value="form.expiry"
                            type="text"
                            inputmode="numeric"
                            class="form-input"
                            :class="{ 'form-input--error': errors.expiry }"
                            placeholder="MM/AA"
                            maxlength="5"
                            @input="onExpiryInput"
                        />
                        <span v-if="errors.expiry" class="form-error">{{ errors.expiry }}</span>
                    </div>

                    <div class="form-group">
                        <label class="form-label">{{ t('client.paiement.cvc') }}</label>
                        <input
                            v-model="form.cvc"
                            type="text"
                            inputmode="numeric"
                            class="form-input"
                            :class="{ 'form-input--error': errors.cvc }"
                            placeholder="123"
                            maxlength="4"
                        />
                        <span v-if="errors.cvc" class="form-error">{{ errors.cvc }}</span>
                    </div>
                </div>

                <div v-if="errors.global" class="error-banner">{{ errors.global }}</div>

                <button type="submit" class="btn-pay" :disabled="submitting">
                    <svg v-if="!submitting" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
                        <line x1="1" y1="10" x2="23" y2="10" />
                    </svg>
                    {{ submitting ? t('client.paiement.processing') : t('client.paiement.pay', { price: itemPrice.toFixed(2) }) }}
                </button>

                <p class="secure-note">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                        <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                    </svg>
                    {{ t('client.paiement.secureNote') }}
                </p>
            </form>

            <aside class="order-summary">
                <div class="summary-title">{{ t('client.paiement.summary') }}</div>
                <div class="summary-item">
                    <span class="summary-item-name">{{ itemName }}</span>
                    <span class="summary-item-price">{{ itemPrice.toFixed(2) }} €</span>
                </div>
                <div class="summary-divider"></div>
                <div class="summary-total">
                    <span>{{ t('client.paiement.total') }}</span>
                    <span class="summary-total-amount">{{ itemPrice.toFixed(2) }} €</span>
                </div>
            </aside>
        </div>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}
.page-header {
    margin-bottom: 32px;
}
.back-link {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.83rem;
    color: var(--green-mid);
    text-decoration: none;
    font-weight: 500;
    margin-bottom: 16px;
    transition: color 0.2s;
}
.back-link:hover {
    color: var(--green-dark);
}
.back-link svg {
    width: 16px;
    height: 16px;
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0;
    line-height: 1.08;
}

.checkout-layout {
    display: flex;
    gap: 32px;
    align-items: flex-start;
}
.checkout-form {
    flex: 1;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    padding: 28px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.form-section-title {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin-bottom: 4px;
}

.stripe-badge {
    display: inline-flex;
    align-items: center;
    gap: 7px;
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--green-mid);
    background: var(--green-pale);
    padding: 6px 12px;
    border-radius: 20px;
    width: fit-content;
}
.stripe-badge svg {
    width: 14px;
    height: 14px;
}
.checkout-form--stripe {
    align-items: center;
    text-align: center;
}
.stripe-explain {
    font-size: 0.88rem;
    color: var(--charcoal);
    opacity: 0.7;
    line-height: 1.6;
    margin: 0;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 7px;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
}
.form-label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--charcoal);
}
.form-input {
    width: 100%;
    padding: 13px 16px;
    font-size: 0.9rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
}
.form-input::placeholder {
    color: rgba(53, 53, 53, 0.35);
}
.form-input:focus {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.1);
}
.form-input--error {
    border-color: #e53e3e;
}
.form-input--card {
    letter-spacing: 0.08em;
    font-size: 1rem;
}
.form-error {
    font-size: 0.78rem;
    color: #e53e3e;
    font-weight: 500;
}
.error-banner {
    background: rgba(229, 62, 62, 0.08);
    border: 1px solid rgba(229, 62, 62, 0.25);
    border-radius: 8px;
    padding: 12px 16px;
    font-size: 0.83rem;
    color: #e53e3e;
}

.btn-pay {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    width: 100%;
    padding: 15px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 10px;
    font-size: 1rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
    margin-top: 4px;
}
.btn-pay svg {
    width: 18px;
    height: 18px;
}
.btn-pay:hover:not(:disabled) {
    background: var(--green-mid);
}
.btn-pay:disabled {
    opacity: 0.65;
    cursor: not-allowed;
}
.secure-note {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.4;
    margin: 0;
    justify-content: center;
}
.secure-note svg {
    width: 13px;
    height: 13px;
}

.order-summary {
    width: 260px;
    flex-shrink: 0;
    background: var(--green-pale);
    border-radius: 14px;
    padding: 24px;
}
.summary-title {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--green-dark);
    margin-bottom: 16px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    font-size: 0.75rem;
}
.summary-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 16px;
}
.summary-item-name {
    font-size: 0.875rem;
    color: var(--charcoal);
    font-weight: 500;
    line-height: 1.4;
    flex: 1;
}
.summary-item-price {
    font-size: 0.875rem;
    font-weight: 700;
    color: var(--green-dark);
    flex-shrink: 0;
}
.summary-divider {
    height: 1px;
    background: rgba(8, 106, 53, 0.15);
    margin-bottom: 16px;
}
.summary-total {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.summary-total span {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--green-dark);
}
.summary-total-amount {
    font-size: 1.2rem;
    font-weight: 800;
    letter-spacing: -0.02em;
}

@media (max-width: 700px) {
    .checkout-layout {
        flex-direction: column;
    }
    .order-summary {
        width: 100%;
    }
    .form-row {
        grid-template-columns: 1fr;
    }
}
</style>
