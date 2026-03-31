<script setup lang="ts">
import { ref, reactive } from 'vue'

const showPassword = ref(false)
const form = reactive({
    email: '',
    password: '',
    rememberMe: false,
})

function handleLogin() {
    console.log('form:', form)
}

const footerLinks = ['À propos', 'Mentions légales', 'Politique de confidentialité']
</script>

<template>
    <div class="page">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <nav class="nav-links">
                    <router-link to="/" class="nav-link">Accueil</router-link>
                    <router-link to="/prestations" class="nav-link">Prestations</router-link>
                    <router-link to="/evenements" class="nav-link">Évènements</router-link>
                    <router-link to="/forum" class="nav-link">Forum</router-link>
                    <router-link to="/a-propos" class="nav-link">À propos</router-link>
                </nav>
                <router-link to="/auth/register" class="btn-nav">
                    S'inscrire / Se connecter
                </router-link>
            </div>
        </header>

        <main class="main">
            <div class="container">
                <h1 class="page-title">Se connecter.</h1>

                <form class="login-form" @submit.prevent="handleLogin">
                    <input
                        id="email"
                        v-model="form.email"
                        type="email"
                        placeholder="Email"
                        class="form-input"
                        autocomplete="email"
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

                    <label class="remember-label">
                        <input
                            type="checkbox"
                            v-model="form.rememberMe"
                            class="remember-checkbox"
                        />
                        <span>Se rappeler de moi</span>
                    </label>

                    <button type="submit" class="btn-submit">Se connecter</button>
                </form>

                <div class="create-account">
                    <router-link to="/auth/register" class="create-link">
                        Créer un compte
                    </router-link>
                </div>
            </div>
        </main>

        <footer class="footer">
            <div class="footer-top">
                <div class="footer-links-wrap">
                    <a v-for="link in footerLinks" :key="link" href="#" class="footer-link">
                        {{ link }}
                    </a>
                </div>
            </div>
            <div class="footer-bottom">
                <div class="footer-container">
                    <span class="footer-logo">UpCycleConnect</span>
                    <div class="footer-lang">
                        <span>Choisir la langue</span>
                        <span class="lang-sep"> - </span>
                        <span>Français</span>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<style scoped>
.page {
    --cream: #f8f5ee;
    --green-dark: #086a35;
    --green-mid: #34895b;
    --green-light: #8bbd94;
    --green-pale: #d7ece1;
    --charcoal: #353535;
    --white: #ffffff;

    background-color: var(--cream);
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
    overflow-x: hidden;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}

.navbar {
    background: var(--cream);
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
    position: sticky;
    top: 0;
    z-index: 100;
}
.nav-container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
    height: 68px;
    display: flex;
    align-items: center;
    gap: 40px;
}
.nav-logo {
    font-weight: 800;
    font-size: 1.1rem;
    color: var(--green-dark);
    text-decoration: none;
    flex-shrink: 0;
    letter-spacing: -0.01em;
}
.nav-links {
    display: flex;
    gap: 32px;
    flex: 1;
    justify-content: center;
}
.nav-link {
    font-size: 0.875rem;
    color: var(--green-light);
    text-decoration: none;
    font-weight: 400;
    transition: color 0.2s;
}
.nav-link:hover {
    color: var(--green-dark);
}
.btn-nav {
    background: var(--green-dark);
    color: var(--white);
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    white-space: nowrap;
    transition: background 0.2s;
    flex-shrink: 0;
}
.btn-nav:hover {
    background: var(--green-mid);
}

.main {
    flex: 1;
    padding: 72px 0 80px;
    display: flex;
    flex-direction: column;
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

.footer {
    background: var(--green-dark);
    color: var(--white);
    margin-top: auto;
}
.footer-top {
    display: flex;
    justify-content: center;
    padding: 32px 32px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.12);
}
.footer-links-wrap {
    display: flex;
    gap: 40px;
}
.footer-link {
    color: rgba(255, 255, 255, 0.75);
    text-decoration: none;
    font-size: 0.85rem;
    transition: color 0.2s;
}
.footer-link:hover {
    color: var(--white);
}
.footer-bottom {
    padding: 20px 32px 28px;
}
.footer-container {
    max-width: 1060px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.footer-logo {
    font-weight: 800;
    font-size: 1.2rem;
    color: var(--white);
    letter-spacing: -0.01em;
}
.footer-lang {
    font-size: 0.85rem;
    color: rgba(255, 255, 255, 0.75);
}
.lang-sep {
    opacity: 0.5;
}

@media (max-width: 700px) {
    .nav-links {
        display: none;
    }
    .page-title {
        font-size: 2.8rem;
    }
    .login-form,
    .create-account {
        max-width: 100%;
    }
    .footer-links-wrap {
        flex-direction: column;
        align-items: center;
        gap: 12px;
    }
    .footer-container {
        flex-direction: column;
        gap: 12px;
        text-align: center;
    }
}
</style>
