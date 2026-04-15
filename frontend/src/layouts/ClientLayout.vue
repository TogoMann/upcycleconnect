<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

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
    <div class="client-layout">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <div class="user-dropdown" ref="dropdownRef">
                    <button class="user-trigger" @click.stop="dropdownOpen = !dropdownOpen">
                        <span class="user-avatar">{{ userInitials() }}</span>
                        <span class="user-name">{{ authStore.user?.first_name || authStore.user?.username }}</span>
                        <svg class="chevron" :class="{ 'chevron--open': dropdownOpen }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <polyline points="6 9 12 15 18 9" />
                        </svg>
                    </button>

                    <div v-if="dropdownOpen" class="dropdown-menu">
                        <router-link to="/particulier" class="dropdown-item" @click="dropdownOpen = false">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                                <rect x="3" y="12" width="4" height="9" />
                                <rect x="10" y="7" width="4" height="14" />
                                <rect x="17" y="3" width="4" height="18" />
                            </svg>
                            Mon espace
                        </router-link>
                        <router-link to="/particulier/profil" class="dropdown-item" @click="dropdownOpen = false">
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

        <div class="layout-body">
            <aside class="sidebar">
                <nav class="sidebar-nav">
                    <router-link
                        to="/particulier"
                        class="sidebar-item"
                        exact-active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="12" width="4" height="9" />
                            <rect x="10" y="7" width="4" height="14" />
                            <rect x="17" y="3" width="4" height="18" />
                        </svg>
                        <span>Tableau de bord</span>
                    </router-link>

                    <router-link
                        to="/particulier/annonces"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                            <polyline points="14 2 14 8 20 8" />
                            <line x1="16" y1="13" x2="8" y2="13" />
                            <line x1="16" y1="17" x2="8" y2="17" />
                            <polyline points="10 9 9 9 8 9" />
                        </svg>
                        <span>Mes Annonces</span>
                    </router-link>

                    <router-link
                        to="/particulier/conteneurs"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                        </svg>
                        <span>Mes Dépôts</span>
                    </router-link>

                    <router-link
                        to="/particulier/conteneurs/deposer"
                        class="sidebar-item sidebar-item--sub"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polyline points="16 16 12 12 8 16" />
                            <line x1="12" y1="12" x2="12" y2="21" />
                            <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3" />
                        </svg>
                        <span>Déposer un objet</span>
                    </router-link>

                    <router-link
                        to="/particulier/score"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
                        </svg>
                        <span>Mon Score</span>
                    </router-link>

                    <router-link
                        to="/particulier/planning"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                            <line x1="16" y1="2" x2="16" y2="6" />
                            <line x1="8" y1="2" x2="8" y2="6" />
                            <line x1="3" y1="10" x2="21" y2="10" />
                        </svg>
                        <span>Planning</span>
                    </router-link>

                    <router-link
                        to="/particulier/catalogue"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="3" width="7" height="7" />
                            <rect x="14" y="3" width="7" height="7" />
                            <rect x="14" y="14" width="7" height="7" />
                            <rect x="3" y="14" width="7" height="7" />
                        </svg>
                        <span>Catalogue</span>
                    </router-link>

                    <router-link
                        to="/particulier/paiement"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
                            <line x1="1" y1="10" x2="23" y2="10" />
                        </svg>
                        <span>Paiement</span>
                    </router-link>

                    <div class="sidebar-divider"></div>

                    <router-link to="/" class="sidebar-item sidebar-item--back">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                            <polyline points="9 22 9 12 15 12 15 22" />
                        </svg>
                        <span>Retour au site</span>
                    </router-link>
                </nav>
            </aside>

            <main class="main-content">
                <RouterView />
            </main>
        </div>

        <footer class="footer">
            <div class="footer-top">
                <div class="footer-links-wrap">
                    <a href="#" class="footer-link">À propos</a>
                    <a href="#" class="footer-link">Mentions légales</a>
                    <a href="#" class="footer-link">Politique de confidentialité</a>
                </div>
            </div>
            <div class="footer-bottom">
                <div class="footer-container">
                    <span class="footer-logo">UpCycleConnect</span>
                    <div class="footer-lang">
                        <span>Choisir la langue</span>
                        <span class="lang-sep">&nbsp;·&nbsp;</span>
                        <span>Français</span>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<style scoped>
.client-layout {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: var(--cream);
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.navbar {
    background: var(--cream);
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
    position: sticky;
    top: 0;
    z-index: 100;
}
.nav-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 32px;
    height: 68px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.nav-logo {
    font-weight: 800;
    font-size: 1.1rem;
    color: var(--green-dark);
    text-decoration: none;
    letter-spacing: -0.01em;
    flex-shrink: 0;
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

.layout-body {
    display: flex;
    flex: 1;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
    padding: 0 32px;
}

.sidebar {
    flex: 0 0 220px;
    padding: 40px 0;
}
.sidebar-nav {
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.sidebar-item {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 11px 16px;
    border-radius: 8px;
    text-decoration: none;
    font-size: 0.92rem;
    color: var(--charcoal);
    font-weight: 500;
    transition: background 0.15s, color 0.15s;
}
.sidebar-item:hover {
    background: var(--green-pale);
    color: var(--green-dark);
}
.sidebar-item--active {
    color: var(--green-dark) !important;
    font-weight: 700;
    background: var(--green-pale);
}
.sidebar-item--sub {
    padding-left: 28px;
    font-size: 0.86rem;
    color: rgba(53, 53, 53, 0.65);
}
.sidebar-item--back {
    color: rgba(53, 53, 53, 0.5);
    font-size: 0.88rem;
}
.sidebar-item--back:hover {
    color: var(--green-dark);
}
.sidebar-icon {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
    color: inherit;
}
.sidebar-divider {
    height: 1px;
    background: rgba(53, 53, 53, 0.1);
    margin: 8px 16px;
}

.main-content {
    flex: 1;
    padding: 40px 0 60px 32px;
    min-width: 0;
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
    max-width: 1200px;
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

@media (max-width: 860px) {
    .layout-body {
        flex-direction: column;
        padding: 0 16px;
    }
    .sidebar {
        flex: none;
        padding: 20px 0 0;
    }
    .sidebar-nav {
        flex-direction: row;
        flex-wrap: wrap;
        gap: 6px;
    }
    .sidebar-item {
        padding: 8px 12px;
        font-size: 0.82rem;
    }
    .sidebar-item--sub {
        padding-left: 12px;
    }
    .main-content {
        padding: 24px 0 40px;
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
