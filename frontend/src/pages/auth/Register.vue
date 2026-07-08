<script setup lang="ts">
import { reactive, ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { usePlansStore, type Plan } from '@/stores/plans'
import { API_BASE } from '@/config'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const plansStore = usePlansStore()

const form = reactive({
    username: '',
    prenom: '',
    nom: '',
    email: '',
    password: '',
    confirmPassword: '',
    planId: null as number | null,
    siret: '',
    acceptTerms: false,
})

const error = ref('')
const loading = ref(false)
const plans = ref<Plan[]>([])
const siretStatus = ref<'idle' | 'checking' | 'valid' | 'invalid'>('idle')

watch(() => form.siret, async (val) => {
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

const selectedPlanName = computed(() => {
    const plan = plans.value.find((p) => p.id === form.planId)
    return plan ? plan.name : ''
})

onMounted(async () => {
    try {
        const res = await plansStore.getPlans()
        plans.value = res.filter((p) => p.is_active)
        if (plans.value.length > 0) {
            form.planId = plans.value[0].id
        }
    } catch (e) {
        console.error('Failed to load plans', e)
    }
})

function isPasswordStrong(pw: string): boolean {
    if (pw.length < 10) return false
    const hasLetter = /[a-zA-Z]/.test(pw)
    const hasNumber = /[0-9]/.test(pw)
    return hasLetter && hasNumber
}

async function handleRegister() {
    error.value = ''

    if (!isPasswordStrong(form.password)) {
        error.value = t('auth.register.weakPassword')
        return
    }

    if (form.password !== form.confirmPassword) {
        error.value = t('auth.register.passwordMismatch')
        return
    }

    if (selectedPlanName.value === 'Pro') {
        if (!form.siret) {
            error.value = t('auth.register.siretRequired')
            return
        }
        if (!/^\d{14}$/.test(form.siret)) {
            error.value = t('auth.register.siretInvalidFormat')
            return
        }
        if (siretStatus.value !== 'valid') {
            error.value = t('auth.register.siretNotFound')
            return
        }
    }

    if (!form.acceptTerms) {
        error.value = t('auth.register.termsRequired')
        return
    }

    loading.value = true
    try {
        await authStore.register(
            form.username,
            form.prenom,
            form.nom,
            form.email,
            form.password,
            form.planId || undefined,
            form.siret || undefined,
        )
        const role = authStore.userRole
        if (role === 'admin') router.push('/admin')
        else if (role === 'pro') router.push('/pro')
        else if (role === 'interne') router.push('/salarie')
        else router.push('/particulier')
    } catch (e: any) {
        error.value = e.message || t('auth.register.genericError')
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="page-content">
        <div class="container">
            <h1 class="page-title">{{ t('auth.register.title') }}</h1>

            <form class="register-form" @submit.prevent="handleRegister">
                <input
                    v-model="form.username"
                    type="text"
                    :placeholder="t('auth.register.username')"
                    class="form-input"
                    autocomplete="username"
                    required
                />

                <div class="form-row">
                    <input
                        v-model="form.prenom"
                        type="text"
                        :placeholder="t('auth.register.firstName')"
                        class="form-input"
                        autocomplete="given-name"
                        required
                    />
                    <input
                        v-model="form.nom"
                        type="text"
                        :placeholder="t('auth.register.lastName')"
                        class="form-input"
                        autocomplete="family-name"
                        required
                    />
                </div>

                <input
                    v-model="form.email"
                    type="email"
                    :placeholder="t('auth.register.email')"
                    class="form-input"
                    autocomplete="email"
                    required
                />

                <input
                    v-model="form.password"
                    type="password"
                    :placeholder="t('auth.register.password')"
                    class="form-input"
                    autocomplete="new-password"
                    required
                />
                <p class="form-hint">{{ t('auth.register.passwordHint') }}</p>

                <input
                    v-model="form.confirmPassword"
                    type="password"
                    :placeholder="t('auth.register.confirmPassword')"
                    class="form-input"
                    autocomplete="new-password"
                    required
                />

                <div class="plan-select-wrap">
                    <label class="plan-label">{{ t('auth.register.choosePlan') }}</label>
                    <select v-model="form.planId" class="form-input plan-select">
                        <option v-for="plan in plans" :key="plan.id" :value="plan.id">
                            {{ plan.name }} - {{ plan.price }}€ / {{ plan.billing_cycle }}
                        </option>
                    </select>
                </div>

                <div v-if="selectedPlanName === 'Pro'" class="siret-wrap">
                    <input
                        v-model="form.siret"
                        type="text"
                        :placeholder="t('auth.register.siret')"
                        class="form-input"
                        :class="{ 'form-input--valid': siretStatus === 'valid', 'form-input--invalid': siretStatus === 'invalid' }"
                        maxlength="14"
                        required
                    />
                    <span v-if="siretStatus === 'checking'" class="siret-status siret-status--checking">{{ t('auth.register.siretChecking') }}</span>
                    <span v-else-if="siretStatus === 'valid'" class="siret-status siret-status--valid">{{ t('auth.register.siretValid') }}</span>
                    <span v-else-if="siretStatus === 'invalid'" class="siret-status siret-status--invalid">{{ t('auth.register.siretNotFoundShort') }}</span>
                </div>

                <p v-if="error" class="error-msg">{{ error }}</p>

                <label class="terms-label">
                    <input type="checkbox" v-model="form.acceptTerms" class="terms-checkbox" />
                    <span>
                        {{ t('auth.register.acceptTerms1') }}
                        <router-link to="/mentions-legales#cgu" class="terms-link">{{ t('auth.register.termsLink') }}</router-link>
                        {{ t('auth.register.acceptTerms2') }}
                        <router-link to="/mentions-legales#confidentialite" class="terms-link">{{ t('auth.register.privacyLink') }}</router-link>
                    </span>
                </label>

                <button type="submit" class="btn-submit" :disabled="loading">
                    {{ loading ? t('auth.register.creating') : t('auth.register.submit') }}
                </button>
            </form>

            <div class="login-link-wrap">
                <router-link to="/auth/login" class="login-link">
                    {{ t('auth.register.alreadyHaveAccount') }}
                </router-link>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-content {
    flex: 1;
    padding: 72px 0 80px;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}

.page-title {
    font-size: clamp(3rem, 7vw, 5.5rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.05;
    letter-spacing: -0.035em;
    margin: 0 0 60px;
    text-align: center;
}

.register-form {
    max-width: 440px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.form-row {
    display: flex;
    gap: 12px;
}

.form-input {
    width: 100%;
    padding: 16px 18px;
    font-size: 0.95rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.35);
    border-radius: 8px;
    outline: none;
    transition:
        border-color 0.2s,
        box-shadow 0.2s;
    box-sizing: border-box;
}
.form-input::placeholder {
    color: rgba(53, 53, 53, 0.45);
}
.form-input:focus {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.12);
}

.plan-select-wrap {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.plan-label {
    font-size: 0.88rem;
    font-weight: 600;
    color: var(--green-dark);
}
.plan-select {
    cursor: pointer;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%23353535'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 14px center;
    background-size: 16px;
    appearance: none;
    padding-right: 40px;
}

.error-msg {
    font-size: 0.88rem;
    color: #c0392b;
    margin: 0;
}

.form-hint {
    font-size: 0.78rem;
    color: rgba(53, 53, 53, 0.5);
    margin: -8px 0 0;
}

.siret-wrap {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.form-input--valid {
    border-color: var(--green-mid);
}
.form-input--invalid {
    border-color: #c0392b;
}
.siret-status {
    font-size: 0.8rem;
    font-weight: 600;
}
.siret-status--checking {
    color: rgba(53, 53, 53, 0.5);
}
.siret-status--valid {
    color: var(--green-mid);
}
.siret-status--invalid {
    color: #c0392b;
}

.terms-label {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    font-size: 0.85rem;
    color: var(--charcoal);
    cursor: pointer;
    user-select: none;
    line-height: 1.5;
    margin-top: 2px;
}
.terms-checkbox {
    width: 15px;
    height: 15px;
    margin-top: 2px;
    border: 1.5px solid rgba(53, 53, 53, 0.4);
    border-radius: 3px;
    cursor: pointer;
    accent-color: var(--green-dark);
    flex-shrink: 0;
}
.terms-link {
    color: var(--green-mid);
    text-decoration: none;
    transition: color 0.2s;
}
.terms-link:hover {
    color: var(--green-dark);
}

.btn-submit {
    width: 100%;
    padding: 16px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    letter-spacing: 0.01em;
    margin-top: 6px;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-submit:disabled {
    opacity: 0.65;
    cursor: not-allowed;
    transform: none;
}
.btn-submit:not(:disabled):hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}
.btn-submit:not(:disabled):active {
    transform: translateY(0);
}

.login-link-wrap {
    max-width: 440px;
    margin: 48px auto 0;
    text-align: center;
}
.login-link {
    font-size: 1rem;
    color: var(--green-light);
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s;
}
.login-link:hover {
    color: var(--green-dark);
}

@media (max-width: 700px) {
    .page-title {
        font-size: 2.8rem;
    }
    .register-form,
    .login-link-wrap {
        max-width: 100%;
    }
    .form-row {
        flex-direction: column;
    }
}
</style>
