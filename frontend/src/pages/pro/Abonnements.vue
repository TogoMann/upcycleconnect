<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { API_BASE } from '@/config'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const route = useRoute()
const checkoutMessage = ref('')

interface Subscription {
    id: { Int64: number; Valid: boolean }
    subscriber_id: { Int64: number; Valid: boolean }
    price: number
    tier: string
    created_at: { Time: string; Valid: boolean }
    until: { Time: string; Valid: boolean }
}

interface PlanOption {
    id: number
    name: string
    description: string
    price: number
    billing_cycle: string
    features: string[]
    is_active: boolean
}

const currentSub = ref<Subscription | null>(null)
const currentTier = ref('Free')
const plans = ref<PlanOption[]>([])
const hasPaymentMethod = ref(true)
const loading = ref(true)
const changingPlan = ref(false)
const showConfirmModal = ref(false)
const selectedPlan = ref<PlanOption | null>(null)
const siretInput = ref('')
const siretError = ref('')
const siretStatus = ref<'idle' | 'checking' | 'valid' | 'invalid'>('idle')

watch(siretInput, async (val) => {
    if (val.length !== 14) {
        siretStatus.value = 'idle'
        return
    }
    siretStatus.value = 'checking'
    try {
        const res = await fetch(`${API_BASE}/siret/verify?siret=${val}`)
        const data = await res.json()
        siretStatus.value = data.valid ? 'valid' : 'invalid'
    } catch {
        siretStatus.value = 'invalid'
    }
})

const daysRemaining = computed(() => {
    if (!currentSub.value?.until?.Valid) return null
    const until = new Date(currentSub.value.until.Time)
    const now = new Date()
    const diff = Math.ceil((until.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
    return Math.max(0, diff)
})

const progressPercent = computed(() => {
    if (!currentSub.value?.created_at?.Valid || !currentSub.value?.until?.Valid) return 0
    const start = new Date(currentSub.value.created_at.Time).getTime()
    const end = new Date(currentSub.value.until.Time).getTime()
    const now = Date.now()
    if (end <= start) return 100
    return Math.min(100, Math.max(0, Math.round(((now - start) / (end - start)) * 100)))
})

const isExpiringSoon = computed(() => daysRemaining.value !== null && daysRemaining.value <= 7)
const isExpired = computed(() => daysRemaining.value !== null && daysRemaining.value <= 0)

function formatDate(d: { Time: string; Valid: boolean } | undefined) {
    if (!d?.Valid) return '—'
    return new Date(d.Time).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}

function planIcon(name: string) {
    if (name === 'Pro') return 'rocket'
    if (name === 'Premium') return 'star'
    return 'leaf'
}

onMounted(async () => {
    const token = authStore.token
    if (!token) return

    const headers = { Authorization: `Bearer ${token}` }

    const [subRes, plansRes, pmRes] = await Promise.all([
        fetch(`${API_BASE}/subscriptions/me`, { headers }).catch(() => null),
        fetch(`${API_BASE}/plans`).catch(() => null),
        fetch(`${API_BASE}/payment-methods/check`, { headers }).catch(() => null),
    ])

    if (subRes?.ok) {
        const data = await subRes.json()
        if (data.tier && !data.id) {
            currentTier.value = data.tier
        } else if (data.id) {
            currentSub.value = data
            currentTier.value = data.tier || 'Free'
        }
    }

    if (plansRes?.ok) {
        const data = await plansRes.json()
        if (Array.isArray(data)) plans.value = data.filter((p: PlanOption) => p.is_active)
    }

    if (pmRes?.ok) {
        const data = await pmRes.json()
        hasPaymentMethod.value = data.has_payment_method !== false
    } else {
        hasPaymentMethod.value = false
    }

    const sessionId = route.query.session_id as string | undefined
    if (sessionId) {
        try {
            const verifyRes = await fetch(`${API_BASE}/payments/verify?session_id=${sessionId}`, { headers })
            if (verifyRes.ok) {
                const verifyData = await verifyRes.json()
                checkoutMessage.value = verifyData.paid
                    ? t('pro.abonnements.paymentConfirmed')
                    : t('pro.abonnements.paymentNotConfirmed')
            }
        } catch {
            checkoutMessage.value = t('pro.abonnements.paymentCheckError')
        }
    }

    loading.value = false
})

function requestChangePlan(plan: PlanOption) {
    if (plan.name.toLowerCase() === currentTier.value.toLowerCase()) return
    selectedPlan.value = plan
    siretInput.value = ''
    siretError.value = ''
    showConfirmModal.value = true
}

async function confirmChangePlan() {
    if (!selectedPlan.value) return
    if (selectedPlan.value.name === 'Pro' && !siretInput.value.trim()) {
        siretError.value = t('pro.abonnements.errorSiretRequired')
        return
    }
    if (selectedPlan.value.name === 'Pro' && siretStatus.value !== 'valid') {
        siretError.value = t('pro.abonnements.errorSiretInvalid')
        return
    }
    changingPlan.value = true
    siretError.value = ''
    try {
        if (selectedPlan.value.price > 0) {
            const res = await fetch(`${API_BASE}/subscriptions/checkout`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${authStore.token}`,
                },
                body: JSON.stringify({
                    plan_id: selectedPlan.value.id,
                    siret: selectedPlan.value.name === 'Pro' ? siretInput.value.trim() : '',
                    return_path: '/pro/abonnements',
                }),
            })
            if (res.ok) {
                const data = await res.json()
                window.location.href = data.url
                return
            }
            const errText = await res.text()
            siretError.value = errText || t('pro.abonnements.errorPaymentCreate')
            changingPlan.value = false
            return
        }

        const res = await fetch(`${API_BASE}/subscriptions/choose`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({
                plan_id: selectedPlan.value.id,
                siret: '',
            }),
        })
        if (res.ok) {
            window.location.reload()
        } else {
            const d = await res.json().catch(() => ({ message: t('pro.abonnements.errorGeneric') }))
            siretError.value = d.message || d.error || t('pro.abonnements.errorChange')
        }
    } catch {
        siretError.value = t('pro.abonnements.errorNetwork')
    }
    changingPlan.value = false
}

function closeModal() {
    showConfirmModal.value = false
    selectedPlan.value = null
}

function changeLabel(plan: PlanOption) {
    const current = currentTier.value.toLowerCase()
    const target = plan.name.toLowerCase()
    if (current === target) return t('pro.abonnements.currentPlan')
    const order = ['free', 'premium', 'pro']
    return order.indexOf(target) > order.indexOf(current) ? t('pro.abonnements.upgrade') : t('pro.abonnements.downgrade')
}
</script>

<template>
    <div class="abonnements">
        <div class="page-header">
            <h1 class="page-title">{{ t('pro.abonnements.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('pro.abonnements.subtitle') }}</p>
        </div>

        <div v-if="loading" class="loading-state">{{ t('pro.abonnements.loading') }}</div>

        <template v-else>
            <div v-if="checkoutMessage" class="alert-banner alert-banner--warning">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                <div>
                    <strong>{{ t('pro.abonnements.payment') }}</strong>
                    <span>{{ checkoutMessage }}</span>
                </div>
            </div>

            <!-- Bandeau carte bancaire manquante -->
            <div v-if="!hasPaymentMethod && currentTier !== 'Free'" class="alert-banner alert-banner--danger">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>
                <div>
                    <strong>{{ t('pro.abonnements.missingPaymentTitle') }}</strong>
                    <span>{{ t('pro.abonnements.missingPaymentDesc') }}</span>
                </div>
            </div>

            <!-- Bandeau expiration proche -->
            <div v-if="isExpiringSoon && !isExpired" class="alert-banner alert-banner--warning">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                <div>
                    <strong>{{ t('pro.abonnements.expiringSoonTitle') }}</strong>
                    <span>{{ t('pro.abonnements.expiringSoonDesc', { tier: currentTier, days: daysRemaining, plural: daysRemaining !== 1 ? 's' : '' }) }}</span>
                </div>
            </div>

            <!-- Bandeau expiré -->
            <div v-if="isExpired" class="alert-banner alert-banner--danger">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
                <div>
                    <strong>{{ t('pro.abonnements.expiredTitle') }}</strong>
                    <span>{{ t('pro.abonnements.expiredDesc', { tier: currentTier }) }}</span>
                </div>
            </div>

            <!-- Section abonnement actuel -->
            <div v-if="currentSub && currentTier !== 'Free'" class="current-sub-card">
                <div class="sub-header">
                    <div class="sub-tier-badge" :class="`tier--${currentTier.toLowerCase()}`">{{ currentTier }}</div>
                    <div class="sub-price">{{ currentSub.price.toFixed(2) }}€<span>{{ t('pro.abonnements.perMonth') }}</span></div>
                </div>
                <div class="sub-details">
                    <div class="sub-detail">
                        <span class="sub-detail-label">{{ t('pro.abonnements.subscriptionStart') }}</span>
                        <span class="sub-detail-value">{{ formatDate(currentSub.created_at) }}</span>
                    </div>
                    <div class="sub-detail">
                        <span class="sub-detail-label">{{ t('pro.abonnements.nextPayment') }}</span>
                        <span class="sub-detail-value">{{ formatDate(currentSub.until) }}</span>
                    </div>
                    <div class="sub-detail">
                        <span class="sub-detail-label">{{ t('pro.abonnements.timeRemaining') }}</span>
                        <span class="sub-detail-value" :class="{ 'text-warning': isExpiringSoon, 'text-danger': isExpired }">
                            {{ isExpired ? t('pro.abonnements.expired') : t('pro.abonnements.daysRemaining', { days: daysRemaining, plural: daysRemaining !== 1 ? 's' : '' }) }}
                        </span>
                    </div>
                </div>
                <div class="progress-bar-wrap">
                    <div class="progress-bar">
                        <div class="progress-fill" :class="{ 'progress--warning': isExpiringSoon }" :style="{ width: progressPercent + '%' }"></div>
                    </div>
                    <div class="progress-labels">
                        <span>{{ t('pro.abonnements.start') }}</span>
                        <span>{{ t('pro.abonnements.elapsed', { percent: progressPercent }) }}</span>
                        <span>{{ t('pro.abonnements.deadline') }}</span>
                    </div>
                </div>
            </div>

            <div v-else-if="currentTier === 'Free'" class="current-sub-card current-sub-card--free">
                <div class="sub-header">
                    <div class="sub-tier-badge tier--free">{{ t('pro.abonnements.free') }}</div>
                </div>
                <p class="free-message">{{ t('pro.abonnements.freeMessage') }}</p>
            </div>

            <!-- Grille des plans -->
            <h2 class="section-title">{{ t('pro.abonnements.choosePlanTitle') }}</h2>
            <div class="plans-grid">
                <div
                    v-for="plan in plans"
                    :key="plan.id"
                    class="plan-card"
                    :class="{ 'plan-card--active': currentTier.toLowerCase() === plan.name.toLowerCase() }"
                >
                    <div class="plan-badge" v-if="currentTier.toLowerCase() === plan.name.toLowerCase()">{{ t('pro.abonnements.active') }}</div>
                    <div class="plan-icon" :class="`icon--${planIcon(plan.name)}`">
                        <svg v-if="planIcon(plan.name) === 'leaf'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 8C8 10 5.9 16.17 3.82 21.34l1.89.66 .27-.77C8 16 14 14 17 8z"/><path d="M20.49 3.51c-3.11 3.11-6.89 4.19-9.49 4.49 2.5 2.5 5 4.5 9 3 .78-3.5-.53-6.12-1-7.49"/></svg>
                        <svg v-else-if="planIcon(plan.name) === 'star'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 15l6-12 6 12"/><path d="M3 22h18"/><path d="M12 15v7"/></svg>
                    </div>
                    <div class="plan-name">{{ plan.name }}</div>
                    <div class="plan-desc" v-if="plan.description">{{ plan.description }}</div>
                    <div class="plan-price">
                        <span class="plan-amount">{{ plan.price.toFixed(0) }}€</span>
                        <span class="plan-period">{{ plan.billing_cycle === 'yearly' ? t('pro.abonnements.perYear') : t('pro.abonnements.perMonth') }}</span>
                    </div>
                    <ul class="plan-features" v-if="plan.features?.length">
                        <li v-for="f in plan.features" :key="f" class="plan-feature">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                            {{ f }}
                        </li>
                    </ul>
                    <button
                        class="plan-btn"
                        :class="{
                            'plan-btn--active': currentTier.toLowerCase() === plan.name.toLowerCase(),
                            'plan-btn--upgrade': currentTier.toLowerCase() !== plan.name.toLowerCase()
                        }"
                        :disabled="currentTier.toLowerCase() === plan.name.toLowerCase() || changingPlan"
                        @click="requestChangePlan(plan)"
                    >
                        {{ changeLabel(plan) }}
                    </button>
                </div>
            </div>
        </template>

        <!-- Modal de confirmation -->
        <Teleport to="body">
            <div v-if="showConfirmModal" class="modal-overlay" @click.self="closeModal">
                <div class="modal-card">
                    <h3 class="modal-title">{{ t('pro.abonnements.changeSubscriptionTitle') }}</h3>
                    <div class="modal-body">
                        <div class="modal-change-summary">
                            <div class="change-from">
                                <span class="change-label">{{ t('pro.abonnements.currentPlan') }}</span>
                                <span class="change-value">{{ currentTier }}</span>
                            </div>
                            <svg class="change-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>
                            <div class="change-to">
                                <span class="change-label">{{ t('pro.abonnements.newPlan') }}</span>
                                <span class="change-value change-value--new">{{ selectedPlan?.name }}</span>
                            </div>
                        </div>

                        <div class="modal-price-info">
                            <p>{{ t('pro.abonnements.billedMessage', { price: selectedPlan?.price.toFixed(2) }) }}</p>
                            <p v-if="currentTier !== 'Free'" class="modal-note">{{ t('pro.abonnements.replacedNote') }}</p>
                        </div>

                        <div v-if="selectedPlan?.name === 'Pro'" class="form-group">
                            <label class="form-label">{{ t('pro.abonnements.siretLabel') }}</label>
                            <input v-model="siretInput" type="text" class="form-input" placeholder="Ex: 12345678901234" maxlength="14" />
                            <p v-if="siretStatus === 'checking'" class="form-hint">{{ t('pro.abonnements.siretChecking') }}</p>
                            <p v-else-if="siretStatus === 'valid'" class="form-hint form-hint--valid">{{ t('pro.abonnements.siretValid') }}</p>
                            <p v-else-if="siretStatus === 'invalid'" class="form-hint form-hint--invalid">{{ t('pro.abonnements.siretInvalid') }}</p>
                            <p v-else class="form-hint">{{ t('pro.abonnements.siretHint') }}</p>
                        </div>

                        <div v-if="siretError" class="alert alert--error">{{ siretError }}</div>
                    </div>
                    <div class="modal-actions">
                        <button class="btn-secondary" @click="closeModal" :disabled="changingPlan">{{ t('pro.abonnements.cancel') }}</button>
                        <button class="btn-primary" @click="confirmChangePlan" :disabled="changingPlan">
                            {{ changingPlan ? t('pro.abonnements.processing') : t('pro.abonnements.confirmChange') }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }

/* Alert banners */
.alert-banner { display: flex; align-items: flex-start; gap: 14px; padding: 16px 20px; border-radius: 12px; margin-bottom: 20px; }
.alert-banner svg { width: 22px; height: 22px; flex-shrink: 0; margin-top: 1px; }
.alert-banner div { display: flex; flex-direction: column; gap: 2px; }
.alert-banner strong { font-size: 0.9rem; font-weight: 700; }
.alert-banner span { font-size: 0.84rem; opacity: 0.85; }
.alert-banner--danger { background: #fef2f2; color: #991b1b; border: 1px solid #fecaca; }
.alert-banner--danger svg { color: #dc2626; }
.alert-banner--warning { background: #fffbeb; color: #92400e; border: 1px solid #fde68a; }
.alert-banner--warning svg { color: #f59e0b; }

/* Current subscription card */
.current-sub-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 16px; padding: 28px; margin-bottom: 32px; }
.current-sub-card--free { text-align: center; padding: 32px; }
.free-message { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 12px 0 0; }
.sub-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.sub-tier-badge { display: inline-block; padding: 6px 16px; border-radius: 24px; font-size: 0.82rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.06em; }
.tier--free { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.tier--premium { background: #fef3c7; color: #92400e; }
.tier--pro { background: #dbeafe; color: #1e40af; }
.sub-price { font-size: 1.6rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.03em; }
.sub-price span { font-size: 0.85rem; font-weight: 500; opacity: 0.6; }

.sub-details { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 20px; }
.sub-detail { background: var(--cream); border-radius: 10px; padding: 14px 16px; }
.sub-detail-label { display: block; font-size: 0.75rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 6px; }
.sub-detail-value { font-size: 1rem; font-weight: 700; color: var(--charcoal); }
.text-warning { color: #d97706 !important; }
.text-danger { color: #dc2626 !important; }

.progress-bar-wrap { margin-top: 4px; }
.progress-bar { height: 6px; background: rgba(53,53,53,0.08); border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; background: var(--green-dark); border-radius: 3px; transition: width 0.6s ease; }
.progress--warning { background: #f59e0b; }
.progress-labels { display: flex; justify-content: space-between; margin-top: 6px; font-size: 0.72rem; color: var(--charcoal); opacity: 0.4; }

/* Section title */
.section-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 16px; }

/* Plans grid */
.plans-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
.plan-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 16px; padding: 28px 24px; display: flex; flex-direction: column; gap: 14px; position: relative; transition: border-color 0.2s, transform 0.2s; }
.plan-card:hover { transform: translateY(-2px); }
.plan-card--active { border-color: var(--green-dark); border-width: 2px; }
.plan-badge { position: absolute; top: -12px; left: 50%; transform: translateX(-50%); background: var(--green-dark); color: var(--white); font-size: 0.72rem; font-weight: 700; padding: 4px 14px; border-radius: 20px; letter-spacing: 0.04em; text-transform: uppercase; white-space: nowrap; }
.plan-icon { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; }
.plan-icon svg { width: 24px; height: 24px; }
.icon--leaf { background: var(--green-pale); color: var(--green-dark); }
.icon--star { background: #fef3c7; color: #f59e0b; }
.icon--rocket { background: #dbeafe; color: #3b82f6; }
.plan-name { font-size: 1.15rem; font-weight: 700; color: var(--charcoal); }
.plan-desc { font-size: 0.82rem; color: var(--charcoal); opacity: 0.5; line-height: 1.4; }
.plan-price { display: flex; align-items: baseline; gap: 4px; }
.plan-amount { font-size: 2.2rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.04em; line-height: 1; }
.plan-period { font-size: 0.85rem; color: var(--charcoal); opacity: 0.5; }
.plan-features { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 10px; flex: 1; }
.plan-feature { display: flex; align-items: center; gap: 10px; font-size: 0.86rem; color: var(--charcoal); opacity: 0.75; }
.plan-feature svg { width: 16px; height: 16px; flex-shrink: 0; color: var(--green-mid); }
.plan-btn { padding: 12px; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: 1.5px solid var(--green-dark); transition: all 0.2s; width: 100%; }
.plan-btn--upgrade { background: var(--green-dark); color: var(--white); }
.plan-btn--upgrade:hover:not(:disabled) { background: var(--green-mid); border-color: var(--green-mid); }
.plan-btn--active { background: transparent; color: var(--green-dark); cursor: default; }
.plan-btn:disabled { opacity: 0.5; cursor: default; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.modal-card { background: var(--white); border-radius: 20px; padding: 32px; max-width: 480px; width: 90%; box-shadow: 0 20px 60px rgba(0,0,0,0.15); }
.modal-title { font-size: 1.2rem; font-weight: 800; color: var(--charcoal); margin: 0 0 24px; }
.modal-body { display: flex; flex-direction: column; gap: 20px; }
.modal-change-summary { display: flex; align-items: center; justify-content: center; gap: 20px; padding: 20px; background: var(--cream); border-radius: 12px; }
.change-from, .change-to { text-align: center; }
.change-label { display: block; font-size: 0.72rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--charcoal); opacity: 0.5; margin-bottom: 6px; }
.change-value { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); }
.change-value--new { color: var(--green-dark); }
.change-arrow { width: 24px; height: 24px; color: var(--charcoal); opacity: 0.3; }
.modal-price-info p { font-size: 0.88rem; color: var(--charcoal); margin: 0; line-height: 1.5; }
.modal-note { opacity: 0.6; font-size: 0.82rem !important; margin-top: 4px !important; }
.modal-actions { display: flex; gap: 12px; justify-content: flex-end; margin-top: 24px; }

/* Form elements */
.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-hint { font-size: 0.78rem; color: var(--charcoal); opacity: 0.45; margin: 0; }
.form-hint--valid { color: var(--green-mid); opacity: 1; font-weight: 600; }
.form-hint--invalid { color: #dc2626; opacity: 1; font-weight: 600; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.86rem; font-weight: 500; }
.alert--error { background: #fee2e2; color: #991b1b; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-secondary { padding: 12px 24px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: border-color 0.2s; }
.btn-secondary:hover { border-color: var(--charcoal); }

@media (max-width: 800px) {
    .plans-grid { grid-template-columns: 1fr; }
    .sub-details { grid-template-columns: 1fr; }
    .modal-change-summary { flex-direction: column; }
}
</style>
