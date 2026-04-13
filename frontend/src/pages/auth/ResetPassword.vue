<script setup lang="ts">
import { ref, reactive } from 'vue'

const done = ref(false)
const form = reactive({
    password: '',
    confirmPassword: '',
})

function handleSubmit() {
    if (form.password !== form.confirmPassword) return
    console.log('reset password')
    done.value = true
}
</script>

<template>
    <div class="page-content">
        <div class="container">
            <h1 class="page-title">Nouveau mot de passe.</h1>

            <div v-if="!done" class="form-wrap">
                <p class="form-desc">Choisissez un nouveau mot de passe pour votre compte.</p>

                <form class="reset-form" @submit.prevent="handleSubmit">
                    <input
                        v-model="form.password"
                        type="password"
                        placeholder="Nouveau mot de passe"
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
                    <p
                        v-if="form.confirmPassword && form.password !== form.confirmPassword"
                        class="error-msg"
                    >
                        Les mots de passe ne correspondent pas.
                    </p>
                    <button
                        type="submit"
                        class="btn-submit"
                        :disabled="form.password !== form.confirmPassword && !!form.confirmPassword"
                    >
                        Réinitialiser le mot de passe
                    </button>
                </form>
            </div>

            <div v-else class="success-wrap">
                <div class="success-icon">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12" />
                    </svg>
                </div>
                <h2 class="success-title">Mot de passe mis à jour !</h2>
                <p class="success-desc">
                    Votre mot de passe a été réinitialisé avec succès. Vous pouvez maintenant vous
                    connecter.
                </p>
                <router-link to="/auth/login" class="btn-login">Se connecter</router-link>
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
    font-size: clamp(2.4rem, 6vw, 5rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.05;
    letter-spacing: -0.035em;
    margin: 0 0 48px;
    text-align: center;
}

.form-wrap {
    max-width: 440px;
    margin: 0 auto;
}

.form-desc {
    font-size: 0.95rem;
    color: var(--charcoal);
    opacity: 0.75;
    line-height: 1.6;
    margin: 0 0 28px;
    text-align: center;
}

.reset-form {
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

.error-msg {
    font-size: 0.82rem;
    color: #c0392b;
    margin: -4px 0 0;
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
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-submit:hover:not(:disabled) {
    background: var(--green-mid);
    transform: translateY(-1px);
}
.btn-submit:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.success-wrap {
    max-width: 440px;
    margin: 0 auto;
    text-align: center;
}
.success-icon {
    color: var(--green-mid);
    margin-bottom: 24px;
}
.success-title {
    font-size: 1.8rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 16px;
}
.success-desc {
    font-size: 0.95rem;
    color: var(--charcoal);
    opacity: 0.75;
    line-height: 1.6;
    margin: 0 0 32px;
}
.btn-login {
    display: inline-block;
    padding: 14px 32px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.95rem;
    font-weight: 600;
    text-decoration: none;
    transition: background 0.2s;
}
.btn-login:hover {
    background: var(--green-mid);
}

@media (max-width: 700px) {
    .page-title {
        font-size: 2.4rem;
    }
    .form-wrap,
    .success-wrap {
        max-width: 100%;
    }
}
</style>
