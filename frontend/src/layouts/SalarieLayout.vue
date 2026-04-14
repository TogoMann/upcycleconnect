<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

function logout() {
    authStore.logout()
    router.push('/auth/login')
}
</script>

<template>
    <div class="salarie-layout">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <div class="nav-profile">
                    <div class="profile-info">
                        <span class="profile-name">
                            {{ authStore.user?.first_name || authStore.user?.username || 'Salarié' }}
                        </span>
                        <span class="profile-role">Espace salarié</span>
                    </div>
                    <button class="btn-logout" @click="logout">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                            <polyline points="16 17 21 12 16 7" />
                            <line x1="21" y1="12" x2="9" y2="12" />
                        </svg>
                    </button>
                </div>
            </div>
        </header>

        <div class="layout-body">
            <aside class="sidebar">
                <nav class="sidebar-nav">
                    <router-link to="/salarie" class="sidebar-item" exact-active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="12" width="4" height="9" />
                            <rect x="10" y="7" width="4" height="14" />
                            <rect x="17" y="3" width="4" height="18" />
                        </svg>
                        <span>Tableau de bord</span>
                    </router-link>

                    <router-link to="/salarie/formations" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
                            <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
                        </svg>
                        <span>Mes Formations</span>
                    </router-link>

                    <router-link to="/salarie/planning" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                            <line x1="16" y1="2" x2="16" y2="6" />
                            <line x1="8" y1="2" x2="8" y2="6" />
                            <line x1="3" y1="10" x2="21" y2="10" />
                        </svg>
                        <span>Planning</span>
                    </router-link>

                    <router-link to="/salarie/conseils" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M12 20h9" />
                            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
                        </svg>
                        <span>Conseils</span>
                    </router-link>

                    <router-link to="/salarie/forum" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                        </svg>
                        <span>Forum</span>
                    </router-link>

                    <router-link to="/salarie/profil" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                            <circle cx="12" cy="7" r="4" />
                        </svg>
                        <span>Mon Profil</span>
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
.salarie-layout {
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
.nav-profile {
    display: flex;
    align-items: center;
    gap: 12px;
}
.profile-info {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 2px;
}
.profile-name {
    font-size: 0.95rem;
    font-weight: 700;
    color: var(--green-dark);
    line-height: 1.2;
}
.profile-role {
    font-size: 0.75rem;
    color: var(--green-light);
    line-height: 1;
}
.btn-logout {
    background: none;
    border: none;
    cursor: pointer;
    padding: 6px;
    color: var(--green-mid);
    display: flex;
    align-items: center;
    border-radius: 6px;
    transition: background 0.2s, color 0.2s;
}
.btn-logout:hover {
    background: var(--green-pale);
    color: var(--green-dark);
}
.btn-logout svg {
    width: 20px;
    height: 20px;
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
