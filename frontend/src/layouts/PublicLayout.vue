<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const footerLinks = ['À propos', 'Mentions légales', 'Politique de confidentialité']

function espaceUrl(): string {
    const role = authStore.userRole
    if (role === 'pro') return '/pro'
    if (role === 'interne') return '/salarie'
    if (role === 'admin') return '/admin'
    return '/particulier'
}

function profilUrl(): string {
    const role = authStore.userRole
    if (role === 'pro') return '/pro/profil'
    if (role === 'interne') return '/salarie/profil'
    if (role === 'admin') return '/admin'
    return '/particulier/profil'
}

function userInitials(): string {
    const u = authStore.user
    if (!u) return '?'
    return ((u.first_name?.[0] ?? '') + (u.last_name?.[0] ?? u.username?.[0] ?? '')).toUpperCase() || '?'
}

function logout() {
    authStore.logout()
    dropdownOpen.value = false
    router.push('/auth/login')
}

function handleClickOutside(e: MouseEvent) {
    if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
        dropdownOpen.value = false
    }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
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

                <router-link v-if="!authStore.isAuthenticated" to="/auth/login" class="btn-nav">
                    S'inscrire / Se connecter
                </router-link>

                <div v-else class="user-dropdown" ref="dropdownRef">
                    <button class="user-trigger" @click.stop="dropdownOpen = !dropdownOpen">
                        <span class="user-avatar">{{ userInitials() }}</span>
                        <span class="user-name">{{ authStore.user?.first_name || authStore.user?.username }}</span>
                        <svg class="chevron" :class="{ 'chevron--open': dropdownOpen }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <polyline points="6 9 12 15 18 9" />
                        </svg>
                    </button>

                    <div v-if="dropdownOpen" class="dropdown-menu">
                        <router-link :to="espaceUrl()" class="dropdown-item" @click="dropdownOpen = false">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                                <rect x="3" y="12" width="4" height="9" />
                                <rect x="10" y="7" width="4" height="14" />
                                <rect x="17" y="3" width="4" height="18" />
                            </svg>
                            Mon espace
                        </router-link>
                        <router-link :to="profilUrl()" class="dropdown-item" @click="dropdownOpen = false">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                                <circle cx="12" cy="7" r="4" />
                            </svg>
                            Mon profil
                        </router-link>
                        <div class="dropdown-divider"></div>
                        <button class="dropdown-item dropdown-item--danger" @click="logout">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                                <polyline points="16 17 21 12 16 7" />
                                <line x1="21" y1="12" x2="9" y2="12" />
                            </svg>
                            Se déconnecter
                        </button>
                    </div>
                </div>
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

.user-dropdown {
    position: relative;
    flex-shrink: 0;
}
.user-trigger {
    display: flex;
    align-items: center;
    gap: 9px;
    background: var(--green-pale);
    border: 1.5px solid rgba(8, 106, 53, 0.15);
    border-radius: 40px;
    padding: 7px 14px 7px 8px;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s, border-color 0.2s;
}
.user-trigger:hover {
    background: rgba(139, 189, 148, 0.35);
    border-color: var(--green-light);
}
.user-avatar {
    width: 28px;
    height: 28px;
    background: var(--green-mid);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.7rem;
    font-weight: 800;
    color: var(--white);
    letter-spacing: -0.02em;
    flex-shrink: 0;
}
.user-name {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--green-dark);
    max-width: 120px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.chevron {
    width: 14px;
    height: 14px;
    color: var(--green-mid);
    transition: transform 0.2s;
    flex-shrink: 0;
}
.chevron--open {
    transform: rotate(180deg);
}

.dropdown-menu {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(8, 106, 53, 0.12), 0 2px 8px rgba(53, 53, 53, 0.08);
    min-width: 200px;
    padding: 6px;
    z-index: 200;
}
.dropdown-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--charcoal);
    text-decoration: none;
    cursor: pointer;
    background: none;
    border: none;
    width: 100%;
    font-family: inherit;
    transition: background 0.15s, color 0.15s;
    text-align: left;
}
.dropdown-item svg {
    width: 16px;
    height: 16px;
    color: var(--green-mid);
    flex-shrink: 0;
}
.dropdown-item:hover {
    background: var(--green-pale);
    color: var(--green-dark);
}
.dropdown-item:hover svg {
    color: var(--green-dark);
}
.dropdown-item--danger {
    color: rgba(53, 53, 53, 0.7);
}
.dropdown-item--danger svg {
    color: rgba(53, 53, 53, 0.4);
}
.dropdown-item--danger:hover {
    background: rgba(229, 62, 62, 0.07);
    color: #c53030;
}
.dropdown-item--danger:hover svg {
    color: #c53030;
}
.dropdown-divider {
    height: 1px;
    background: rgba(53, 53, 53, 0.08);
    margin: 4px 6px;
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
