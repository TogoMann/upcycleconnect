<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { API_BASE } from '@/config'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const token = ref('')
const password = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)
const success = ref(false)

onMounted(() => {
    token.value = route.query.token as string || ''
    if (!token.value) {
        error.value = t('auth.resetPassword.missingToken')
    }
})

function isPasswordStrong(pw: string): boolean {
    if (pw.length < 10) return false
    return /[a-zA-Z]/.test(pw) && /[0-9]/.test(pw)
}

const handleSubmit = async () => {
    if (!isPasswordStrong(password.value)) {
        error.value = t('auth.register.weakPassword')
        return
    }

    if (password.value !== confirmPassword.value) {
        error.value = t('auth.register.passwordMismatch')
        return
    }

    isLoading.value = true
    error.value = null

    try {
        const res = await fetch(`${API_BASE}/auth/reset-password`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                token: token.value,
                password: password.value
            })
        })

        if (!res.ok) {
            const data = await res.text()
            throw new Error(data || t('auth.resetPassword.error'))
        }

        success.value = true
        setTimeout(() => {
            router.push('/auth/login')
        }, 3000)
    } catch (err: any) {
        error.value = err.message
    } finally {
        isLoading.value = false
    }
}
</script>

<template>
    <div class="reset-page">
        <div class="card">
            <h2 class="title">{{ t('auth.resetPassword.title') }}</h2>
            <p class="subtitle">{{ t('auth.resetPassword.subtitle') }}</p>

            <div v-if="success" class="state-success">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="24" height="24">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                    <polyline points="22 4 12 14.01 9 11.01" />
                </svg>
                {{ t('auth.resetPassword.success') }}
            </div>

            <form v-else-if="token" @submit.prevent="handleSubmit" class="form">
                <div v-if="error" class="state-error">{{ error }}</div>

                <div class="form-group">
                    <label>{{ t('auth.resetPassword.newPassword') }}</label>
                    <input
                        type="password"
                        v-model="password"
                        required
                        class="form-input"
                        placeholder="••••••••"
                    />
                    <p class="form-hint">{{ t('auth.register.passwordHint') }}</p>
                </div>

                <div class="form-group">
                    <label>{{ t('auth.register.confirmPassword') }}</label>
                    <input
                        type="password"
                        v-model="confirmPassword"
                        required
                        class="form-input"
                        placeholder="••••••••"
                    />
                </div>

                <button type="submit" class="btn-submit" :disabled="isLoading">
                    {{ isLoading ? t('auth.resetPassword.submitting') : t('auth.resetPassword.submit') }}
                </button>
            </form>

            <div v-else class="state-error">
                {{ error || t('auth.resetPassword.missingToken') }}
            </div>

            <router-link to="/auth/login" class="back-link">{{ t('auth.forgotPassword.backToLogin') }}</router-link>
        </div>
    </div>
</template>

<style scoped>
.reset-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f8f5ee;
    padding: 20px;
    font-family: 'Inter', sans-serif;
}

.card {
    background: white;
    width: 100%;
    max-width: 400px;
    padding: 40px;
    border-radius: 16px;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.05);
}

.title {
    margin: 0 0 8px;
    font-size: 1.5rem;
    font-weight: 800;
    color: #353535;
}

.subtitle {
    margin: 0 0 32px;
    font-size: 0.9rem;
    color: rgba(53, 53, 53, 0.6);
}

.form {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.form-group label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #353535;
}

.form-hint {
    font-size: 0.78rem;
    color: rgba(53, 53, 53, 0.5);
    margin: 0;
}

.form-input {
    width: 100%;
    padding: 12px 16px;
    border: 1.5px solid rgba(53, 53, 53, 0.15);
    border-radius: 8px;
    font-family: inherit;
    box-sizing: border-box;
}

.form-input:focus {
    outline: none;
    border-color: #34895b;
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.1);
}

.btn-submit {
    background: #086a35;
    color: white;
    border: none;
    padding: 14px;
    border-radius: 8px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
    margin-top: 10px;
}

.btn-submit:hover {
    background: #34895b;
}

.btn-submit:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.state-error {
    background: #fdecea;
    color: #c0392b;
    padding: 12px;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 500;
}

.state-success {
    background: #d7ece1;
    color: #086a35;
    padding: 20px;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    text-align: center;
    line-height: 1.5;
}

.back-link {
    display: block;
    text-align: center;
    margin-top: 24px;
    font-size: 0.85rem;
    color: #34895b;
    text-decoration: none;
    font-weight: 600;
}

.back-link:hover {
    text-decoration: underline;
}
</style>
