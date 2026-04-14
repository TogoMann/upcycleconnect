<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const currentPlan = ref('individual')
const loading = ref(false)

const plans = [
    {
        id: 'individual',
        name: 'Individuel',
        price: '0',
        features: ['5 annonces actives', 'Accès catalogue basique', 'Dépôt conteneurs'],
    },
    {
        id: 'pro',
        name: 'Pro',
        price: '29',
        features: ['Annonces illimitées', 'Publicités', 'Facturation PDF', 'Récupération objets'],
    },
    {
        id: 'premium',
        name: 'Premium',
        price: '79',
        features: ['Tout Pro inclus', 'Tableau bord avancé', 'Support prioritaire', 'Statistiques avancées'],
    },
]

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro/abonnements', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) {
            const data = await res.json()
            currentPlan.value = data.plan ?? 'individual'
        }
    } catch {}
})

async function changePlan(planId: string) {
    if (planId === currentPlan.value) return
    loading.value = true
    try {
        const res = await fetch('http://localhost:8081/pro/abonnements', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({ plan: planId }),
        })
        if (res.ok) currentPlan.value = planId
    } catch {}
    loading.value = false
}
</script>

<template>
    <div class="abonnements">
        <div class="page-header">
            <h1 class="page-title">Abonnements.</h1>
            <p class="page-subtitle">Gérez votre formule et passez à un niveau supérieur.</p>
        </div>

        <div class="plans-grid">
            <div
                v-for="plan in plans"
                :key="plan.id"
                class="plan-card"
                :class="{ 'plan-card--active': currentPlan === plan.id }"
            >
                <div class="plan-badge" v-if="currentPlan === plan.id">Actif</div>
                <div class="plan-name">{{ plan.name }}</div>
                <div class="plan-price">
                    <span class="plan-amount">{{ plan.price }}€</span>
                    <span class="plan-period">/mois</span>
                </div>
                <ul class="plan-features">
                    <li v-for="f in plan.features" :key="f" class="plan-feature">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                            <polyline points="20 6 9 17 4 12" />
                        </svg>
                        {{ f }}
                    </li>
                </ul>
                <button
                    class="plan-btn"
                    :class="{ 'plan-btn--active': currentPlan === plan.id }"
                    :disabled="currentPlan === plan.id || loading"
                    @click="changePlan(plan.id)"
                >
                    {{ currentPlan === plan.id ? 'Plan actuel' : currentPlan === 'premium' && plan.id !== 'premium' ? 'Rétrograder' : 'Choisir' }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.plans-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
.plan-card { background: var(--white); border: 1.5px solid rgba(53, 53, 53, 0.1); border-radius: 16px; padding: 28px 24px; display: flex; flex-direction: column; gap: 16px; position: relative; }
.plan-card--active { border-color: var(--green-dark); border-width: 2px; }
.plan-badge { position: absolute; top: -12px; left: 50%; transform: translateX(-50%); background: var(--green-dark); color: var(--white); font-size: 0.72rem; font-weight: 700; padding: 4px 12px; border-radius: 20px; letter-spacing: 0.04em; text-transform: uppercase; white-space: nowrap; }
.plan-name { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); }
.plan-price { display: flex; align-items: baseline; gap: 4px; }
.plan-amount { font-size: 2.4rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.04em; line-height: 1; }
.plan-period { font-size: 0.85rem; color: var(--charcoal); opacity: 0.5; }
.plan-features { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 10px; flex: 1; }
.plan-feature { display: flex; align-items: center; gap: 10px; font-size: 0.88rem; color: var(--charcoal); opacity: 0.8; }
.plan-feature svg { width: 16px; height: 16px; flex-shrink: 0; color: var(--green-mid); }
.plan-btn { padding: 12px; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: 1.5px solid var(--green-dark); background: var(--green-dark); color: var(--white); transition: background 0.2s; }
.plan-btn:hover:not(:disabled) { background: var(--green-mid); border-color: var(--green-mid); }
.plan-btn--active { background: transparent; color: var(--green-dark); cursor: default; }
.plan-btn:disabled { opacity: 0.6; }
@media (max-width: 800px) { .plans-grid { grid-template-columns: 1fr; } }
</style>
