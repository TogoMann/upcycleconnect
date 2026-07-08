<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePlansStore, type Plan } from '@/stores/plans'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const plansStore = usePlansStore()
const authStore = useAuthStore()

const plans = ref<Plan[]>([])
const loading = ref(true)

onMounted(async () => {
    try {
        const res = await plansStore.getPlans()
        plans.value = res.filter((p: Plan) => p.is_active)
    } catch {
        plans.value = []
    } finally {
        loading.value = false
    }
})

function priceLabel(plan: Plan): string {
    if (!plan.price) return t('tarifs.free')
    return plan.billing_cycle === 'yearly'
        ? t('tarifs.perYear', { price: plan.price })
        : t('tarifs.perMonth', { price: plan.price })
}

function prerequisite(plan: Plan): string {
    return plan.name === 'Pro' ? t('tarifs.siretRequired') : t('tarifs.noPrerequisite')
}

const ctaTarget = computed(() => {
    if (!authStore.isAuthenticated) return '/auth/register'
    const role = authStore.userRole
    if (role === 'pro') return '/pro/abonnements'
    if (role === 'client') return '/particulier/plans'
    return null
})

const ctaLabel = computed(() => (authStore.isAuthenticated ? t('tarifs.manage') : t('tarifs.cta')))
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">{{ t('tarifs.pageTitle') }}</h1>
                <p class="hero-subtitle">{{ t('tarifs.subtitle') }}</p>
            </div>
        </section>

        <section class="plans-section">
            <div class="container">
                <div v-if="loading" class="empty-state">
                    <p>{{ t('tarifs.loading') }}</p>
                </div>

                <div v-else class="plans-grid">
                    <div v-for="plan in plans" :key="plan.id" class="plan-card" :class="{ 'plan-card--pro': plan.name === 'Pro' }">
                        <span v-if="plan.name === 'Pro'" class="badge-recommended">{{ t('tarifs.recommended') }}</span>

                        <h2 class="plan-name">{{ plan.name }}</h2>
                        <div class="plan-price">{{ priceLabel(plan) }}</div>
                        <p class="plan-desc">{{ plan.description }}</p>

                        <div class="plan-prereq">
                            <span class="plan-prereq-label">{{ t('tarifs.prerequisiteLabel') }}</span>
                            <span>{{ prerequisite(plan) }}</span>
                        </div>

                        <ul class="plan-features">
                            <li v-for="feature in plan.features" :key="feature">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" width="16" height="16">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                                {{ feature }}
                            </li>
                        </ul>

                        <router-link v-if="ctaTarget" :to="ctaTarget" class="btn-cta">{{ ctaLabel }}</router-link>
                    </div>
                </div>

                <div v-if="!loading && plans.length === 0" class="empty-state">
                    <p>{{ t('tarifs.noResults') }}</p>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page-content {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}

.hero {
    padding: 64px 0 48px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.6rem, 5.5vw, 4.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0 0 16px;
}
.hero-subtitle {
    font-size: 1rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
}

.plans-section {
    flex: 1;
    padding: 0 0 80px;
}

.plans-grid {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 28px;
}

.plan-card {
    position: relative;
    background: var(--white);
    border: 2px solid rgba(53, 53, 53, 0.1);
    border-radius: 16px;
    padding: 32px;
    width: 320px;
    display: flex;
    flex-direction: column;
    transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
}
.plan-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(8, 106, 53, 0.1);
}
.plan-card--pro {
    border-color: var(--green-mid);
    box-shadow: 0 10px 30px rgba(8, 106, 53, 0.12);
}

.badge-recommended {
    position: absolute;
    top: -13px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--green-mid);
    color: var(--white);
    padding: 4px 14px;
    border-radius: 20px;
    font-size: 0.72rem;
    font-weight: 800;
    letter-spacing: 0.02em;
}

.plan-name {
    font-size: 1.4rem;
    font-weight: 800;
    color: var(--charcoal);
    margin: 8px 0 12px;
    letter-spacing: -0.01em;
}
.plan-price {
    font-size: 1.6rem;
    font-weight: 800;
    color: var(--green-dark);
    margin-bottom: 16px;
    letter-spacing: -0.02em;
}
.plan-desc {
    font-size: 0.88rem;
    color: var(--charcoal);
    opacity: 0.75;
    line-height: 1.6;
    margin: 0 0 20px;
    min-height: 44px;
}

.plan-prereq {
    display: flex;
    flex-direction: column;
    gap: 2px;
    background: var(--green-pale);
    border-radius: 8px;
    padding: 10px 14px;
    margin-bottom: 20px;
    font-size: 0.82rem;
    color: var(--green-dark);
}
.plan-prereq-label {
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    opacity: 0.75;
}

.plan-features {
    list-style: none;
    padding: 0;
    margin: 0 0 28px;
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.plan-features li {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 0.9rem;
    color: var(--charcoal);
}
.plan-features svg {
    color: var(--green-mid);
    flex-shrink: 0;
}

.btn-cta {
    display: block;
    text-align: center;
    padding: 13px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 700;
    text-decoration: none;
    transition: background 0.2s, transform 0.15s;
}
.btn-cta:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}

.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
}

@media (max-width: 560px) {
    .hero-title {
        font-size: 2.2rem;
    }
    .plan-card {
        width: 100%;
    }
}
</style>
