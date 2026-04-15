<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
    username: '',
    prenom: '',
    nom: '',
    email: '',
    password: '',
    confirmPassword: '',
    acceptTerms: false,
})

const error = ref('')
const loading = ref(false)

async function handleRegister() {
    error.value = ''

    if (form.password !== form.confirmPassword) {
        error.value = 'Les mots de passe ne correspondent pas.'
        return
    }
    if (!form.acceptTerms) {
        error.value = "Vous devez accepter les conditions d'utilisation."
        return
    }

    loading.value = true
    try {
        await authStore.register(form.username, form.prenom, form.nom, form.email, form.password)
        const role = authStore.userRole
        if (role === 'admin') router.push('/admin')
        else if (role === 'pro') router.push('/pro')
        else if (role === 'interne') router.push('/salarie')
        else router.push('/particulier')
    } catch (e: any) {
        error.value = e.message || 'Erreur lors de la création du compte.'
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="page-content">
        <div class="container">
            <h1 class="page-title">Créer un compte.</h1>

            <form class="register-form" @submit.prevent="handleRegister">
                <input
                    v-model="form.username"
                    type="text"
                    placeholder="Nom d'utilisateur"
                    class="form-input"
                    autocomplete="username"
                    required
                />

                <div class="form-row">
                    <input
                        v-model="form.prenom"
                        type="text"
                        placeholder="Prénom"
                        class="form-input"
                        autocomplete="given-name"
                        required
                    />
                    <input
                        v-model="form.nom"
                        type="text"
                        placeholder="Nom"
                        class="form-input"
                        autocomplete="family-name"
                        required
                    />
                </div>

                <input
                    v-model="form.email"
                    type="email"
                    placeholder="Email"
                    class="form-input"
                    autocomplete="email"
                    required
                />

                <input
                    v-model="form.password"
                    type="password"
                    placeholder="Mot de passe"
                    class="form-input"
                    autocomplete="new-password"
                    required
                />

                <input
                    v-model="form.confirmPassword"
                    type="password"
                    placeholder="Confirmer le mot de passe"
                    class="form-input"
                    autocomplete="new-password"
                    required
                />

                <p v-if="error" class="error-msg">{{ error }}</p>

                <label class="terms-label">
                    <input
                        type="checkbox"
                        v-model="form.acceptTerms"
                        class="terms-checkbox"
                    />
                    <span>
                        J'accepte les
                        <a href="#" class="terms-link">conditions d'utilisation</a>
                        et la
                        <a href="#" class="terms-link">politique de confidentialité</a>
                    </span>
                </label>

                <button type="submit" class="btn-submit" :disabled="loading">
                    {{ loading ? 'Création...' : 'Créer mon compte' }}
                </button>
            </form>

            <div class="login-link-wrap">
                <router-link to="/auth/login" class="login-link">
                    J'ai déjà un compte
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

.error-msg {
    font-size: 0.88rem;
    color: #c0392b;
    margin: 0;
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
