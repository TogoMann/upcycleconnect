<script setup lang="ts">
import { RouterView } from 'vue-router'

const footerLinks = ['À propos', 'Mentions légales', 'Politique de confidentialité']
</script>

<template>
    <div class="public-layout">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <nav class="nav-links">
                    <router-link to="/" class="nav-link" exact-active-class="nav-link--active">Accueil</router-link>
                    <router-link to="/prestations" class="nav-link" active-class="nav-link--active">Prestations</router-link>
                    <router-link to="/evenements" class="nav-link" active-class="nav-link--active">Évènements</router-link>
                    <router-link to="/annonces" class="nav-link" active-class="nav-link--active">Annonces</router-link>
                    <router-link to="/forum" class="nav-link" active-class="nav-link--active">Forum</router-link>
                    <router-link to="/a-propos" class="nav-link" active-class="nav-link--active">À propos</router-link>
                </nav>
                <router-link to="/auth/login" class="btn-nav">
                    S'inscrire / Se connecter
                </router-link>
            </div>
        </header>

        <main class="layout-main">
            <RouterView />
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
                        <span class="lang-sep"> &nbsp;·&nbsp; </span>
                        <span>Français</span>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<style scoped>
.public-layout {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
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
    flex-shrink: 0;
    text-decoration: none;
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
    color: var(--green-mid);
    text-decoration: none;
    font-weight: 400;
    transition: color 0.2s;
}
.nav-link:hover {
    color: var(--green-dark);
}
.nav-link--active {
    color: var(--green-dark);
    font-weight: 600;
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

.layout-main {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.footer {
    background: var(--green-dark);
    color: var(--white);
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

@media (max-width: 900px) {
    .nav-links {
        display: none;
    }
}
@media (max-width: 560px) {
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
