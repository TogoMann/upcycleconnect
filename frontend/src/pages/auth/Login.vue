<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const error = ref('')
const loading = ref(false)
const showPassword = ref(false)
const form = reactive({
    username: '',
    password: '',
})

async function handleLogin() {
    if (!form.username || !form.password) return
    loading.value = true
    error.value = ''
    try {
        await authStore.login(form.username, form.password)
        const role = authStore.userRole
        if (role === 'admin') router.push('/admin')
        else if (role === 'pro') router.push('/pro')
        else if (role === 'salarie') router.push('/salarie')
        else router.push('/particulier')
    } catch (e: any) {
        error.value = 'Identifiants invalides.'
    }
    loading.value = false
}
</script>

<template>
    <div class="page-content">
        <div class="container">
            <h1 class="page-title">Se connecter.</h1>

            <form class="login-form" @submit.prevent="handleLogin">
                <div v-if="error" class="form-error">{{ error }}</div>
                <input
                    id="username"
                    v-model="form.username"
                    type="text"
                    placeholder="Nom d'utilisateur"
                    class="form-input"
                    autocomplete="username"
                />

                <div class="password-field">
                    <input
                        id="password"
                        v-model="form.password"
                        :type="showPassword ? 'text' : 'password'"
                        placeholder="Mot de passe"
                        class="form-input"
                        autocomplete="current-password"
                    />
                    <router-link to="/auth/forgot-password" class="forgot-link">
                        Mot de passe oublié
                    </router-link>
                </div>

                <button type="submit" class="btn-submit" :disabled="loading">
                    {{ loading ? 'Connexion…' : 'Se connecter' }}
                </button>
            </form>

            <div class="create-account">
                <router-link to="/auth/register" class="create-link">
                    Créer un compte
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

.login-form {
    max-width: 440px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 16px;
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

.password-field {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 0;
}
.password-field .form-input {
    padding-right: 18px;
}
.forgot-link {
    align-self: flex-end;
    margin-top: 6px;
    font-size: 0.8rem;
    color: var(--green-light);
    text-decoration: none;
    transition: color 0.2s;
}
.forgot-link:hover {
    color: var(--green-dark);
}

.remember-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.875rem;
    color: var(--charcoal);
    cursor: pointer;
    user-select: none;
    margin-top: 2px;
}
.remember-checkbox {
    width: 15px;
    height: 15px;
    border: 1.5px solid rgba(53, 53, 53, 0.4);
    border-radius: 3px;
    cursor: pointer;
    accent-color: var(--green-dark);
    flex-shrink: 0;
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
.btn-submit:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}
.btn-submit:active {
    transform: translateY(0);
}

.create-account {
    max-width: 440px;
    margin: 48px auto 0;
    text-align: center;
}
.create-link {
    font-size: 1rem;
    color: var(--green-light);
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s;
}
.create-link:hover {
    color: var(--green-dark);
}

.form-error {
    background: #fee2e2;
    color: #991b1b;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 0.88rem;
    font-weight: 500;
}

@media (max-width: 700px) {
    .page-title {
        font-size: 2.8rem;
    }
    .login-form,
    .create-account {
        max-width: 100%;
    }
}
</style>
