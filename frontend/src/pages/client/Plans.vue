<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { usePlansStore, type Plan } from '@/stores/plans';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const plansStore = usePlansStore();
const authStore = useAuthStore();
const router = useRouter();

const plans = ref<Plan[]>([]);
const selectedPlan = ref<Plan | null>(null);
const siret = ref('');
const loading = ref(false);
const error = ref('');
const success = ref('');

async function loadPlans() {
    try {
        const res = await plansStore.getPlans();
        plans.value = res.filter((p: Plan) => p.is_active);
    } catch (e: any) {
        error.value = e.message;
    }
}

async function handleChoosePlan(plan: Plan) {
    if (plan.name === 'Pro' && !siret.value) {
        selectedPlan.value = plan;
        return;
    }

    loading.value = true;
    error.value = '';
    success.value = '';

    try {
        await plansStore.choosePlan(plan.id, plan.name === 'Pro' ? siret.value : undefined);
        success.value = 'Votre plan a été mis à jour avec succès !';
        setTimeout(() => {
            router.push('/particulier');
        }, 2000);
    } catch (e: any) {
        error.value = e.message;
    } finally {
        loading.value = false;
    }
}

onMounted(loadPlans);
</script>

<template>
    <div class="plans-page">
        <header class="page-header">
            <h1>Choisissez votre forfait</h1>
            <p>Passez à la vitesse supérieure avec nos offres adaptées à vos besoins.</p>
        </header>

        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="plans-grid">
            <div v-for="plan in plans" :key="plan.id" class="plan-card" :class="{ 'plan-pro': plan.name === 'Pro' }">
                <div class="plan-header">
                    <h2>{{ plan.name }}</h2>
                    <div class="plan-price">
                        <span class="amount">{{ plan.price }}€</span>
                        <span class="cycle">/{{ plan.billing_cycle === 'monthly' ? 'mois' : 'an' }}</span>
                    </div>
                </div>
                
                <p class="plan-desc">{{ plan.description }}</p>
                
                <ul class="plan-features">
                    <li v-for="feature in plan.features" :key="feature">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" width="16" height="16">
                            <polyline points="20 6 9 17 4 12" />
                        </svg>
                        {{ feature }}
                    </li>
                </ul>

                <div v-if="selectedPlan?.id === plan.id && plan.name === 'Pro'" class="siret-field">
                    <label for="siret">Numéro SIRET (14 chiffres)</label>
                    <input type="text" id="siret" v-model="siret" placeholder="12345678901234" maxlength="14">
                    <p class="siret-hint">Requis pour les professionnels</p>
                </div>

                <button 
                    @click="handleChoosePlan(plan)" 
                    class="btn-choose"
                    :disabled="loading || (authStore.user?.role === 'pro' && plan.name !== 'Pro')"
                >
                    {{ loading ? 'Traitement...' : 'Sélectionner ce plan' }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.plans-page {
    padding: 40px 0;
}

.page-header {
    text-align: center;
    margin-bottom: 50px;
}

.page-header h1 {
    font-size: 2.5rem;
    color: var(--green-dark);
    margin-bottom: 10px;
}

.page-header p {
    color: #718096;
    font-size: 1.1rem;
}

.alert {
    max-width: 600px;
    margin: 0 auto 30px;
    padding: 15px;
    border-radius: 8px;
    text-align: center;
}

.alert-error { background: #fed7d7; color: #c53030; }
.alert-success { background: #c6f6d5; color: #22543d; }

.plans-grid {
    display: flex;
    justify-content: center;
    gap: 30px;
    flex-wrap: wrap;
}

.plan-card {
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 20px;
    padding: 30px;
    width: 320px;
    display: flex;
    flex-direction: column;
    transition: transform 0.3s, border-color 0.3s;
}

.plan-card:hover {
    transform: translateY(-5px);
    border-color: var(--green-mid);
}

.plan-pro {
    border-color: var(--green-mid);
    box-shadow: 0 10px 25px rgba(8, 106, 53, 0.1);
    position: relative;
}

.plan-pro::before {
    content: "RECOMMANDÉ";
    position: absolute;
    top: -12px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--green-mid);
    color: white;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 800;
}

.plan-header h2 {
    margin: 0 0 15px;
    font-size: 1.5rem;
    color: var(--charcoal);
}

.plan-price {
    margin-bottom: 20px;
}

.amount {
    font-size: 2.5rem;
    font-weight: 800;
    color: var(--green-dark);
}

.cycle {
    color: #a0aec0;
}

.plan-desc {
    color: #4a5568;
    margin-bottom: 25px;
    min-height: 50px;
}

.plan-features {
    list-style: none;
    padding: 0;
    margin: 0 0 30px;
    flex: 1;
}

.plan-features li {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 12px;
    font-size: 0.95rem;
    color: #2d3748;
}

.plan-features svg {
    color: #38a169;
}

.siret-field {
    margin-bottom: 20px;
}

.siret-field label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    margin-bottom: 5px;
}

.siret-field input {
    width: 100%;
    padding: 10px;
    border: 1px solid #cbd5e0;
    border-radius: 6px;
    font-family: inherit;
}

.siret-hint {
    font-size: 0.75rem;
    color: #718096;
    margin-top: 5px;
}

.btn-choose {
    background: var(--green-pale);
    color: var(--green-dark);
    border: none;
    padding: 14px;
    border-radius: 10px;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-choose:hover:not(:disabled) {
    background: var(--green-mid);
    color: white;
}

.btn-choose:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.plan-pro .btn-choose {
    background: var(--green-dark);
    color: white;
}

.plan-pro .btn-choose:hover:not(:disabled) {
    background: #064d35;
}
</style>
