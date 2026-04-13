<script setup lang="ts">
import { ref, reactive } from 'vue'

const sent = ref(false)
const form = reactive({ email: '' })

function handleSubmit() {
    console.log('forgot password:', form.email)
    sent.value = true
}
</script>

<template>
    <div class="page-content">
        <div class="container">
            <h1 class="page-title">Mot de passe oublié.</h1>

            <div v-if="!sent" class="form-wrap">
                <p class="form-desc">
                    Entrez votre adresse email et nous vous enverrons un lien pour réinitialiser
                    votre mot de passe.
                </p>
                <form class="forgot-form" @submit.prevent="handleSubmit">
                    <input
                        v-model="form.email"
                        type="email"
                        placeholder="Email"
                        class="form-input"
                        autocomplete="email"
                        required
                    />
                    <button type="submit" class="btn-submit">Envoyer le lien</button>
                </form>
                <div class="back-wrap">
                    <router-link to="/auth/login" class="back-link">
                        Retour à la connexion
                    </router-link>
                </div>
            </div>

            <div v-else class="success-wrap">
                <div class="success-icon">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.07 13a19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 2.92 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L7.91 9.91a16 16 0 0 0 6.09 6.09l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
                    </svg>
                </div>
                <h2 class="success-title">Email envoyé !</h2>
                <p class="success-desc">
                    Si cette adresse est associée à un compte, vous recevrez un email avec un lien
                    de réinitialisation dans quelques minutes.
                </p>
                <router-link to="/auth/login" class="btn-back">Retour à la connexion</router-link>
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

.forgot-form {
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
.btn-submit:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}
.btn-submit:active {
    transform: translateY(0);
}

.back-wrap {
    margin-top: 28px;
    text-align: center;
}
.back-link {
    font-size: 0.9rem;
    color: var(--green-light);
    text-decoration: none;
    transition: color 0.2s;
}
.back-link:hover {
    color: var(--green-dark);
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
.btn-back {
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
.btn-back:hover {
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
