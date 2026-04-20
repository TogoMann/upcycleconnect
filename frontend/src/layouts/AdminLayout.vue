<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

function logout() {
    authStore.logout()
    router.push('/auth/login')
}
</script>

<template>
    <div class="admin-layout">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <div class="nav-profile">
                    <div class="profile-info">
                        <span class="profile-name">
                            {{ authStore.user?.first_name || authStore.user?.username || 'Admin' }}
                        </span>
                        <span class="profile-role">Administrateur</span>
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
                    <router-link to="/admin" class="sidebar-item" exact-active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="12" width="4" height="9" />
                            <rect x="10" y="7" width="4" height="14" />
                            <rect x="17" y="3" width="4" height="18" />
                        </svg>
                        <span>Dashboard</span>
                    </router-link>

                    <router-link to="/admin/utilisateurs" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                            <circle cx="9" cy="7" r="4" />
                            <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                        </svg>
                        <span>Utilisateurs</span>
                    </router-link>

                    <router-link to="/admin/conteneurs" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                        </svg>
                        <span>Conteneurs</span>
                    </router-link>

                    <router-link to="/admin/depots" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polyline points="16 16 12 12 8 16" />
                            <line x1="12" y1="12" x2="12" y2="21" />
                            <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3" />
                        </svg>
                        <span>Dépôts</span>
                    </router-link>

                    <router-link to="/admin/abonnements" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="1" y="4" width="22" height="16" rx="2" ry="2" />
                            <line x1="1" y1="10" x2="23" y2="10" />
                        </svg>
                        <span>Abonnements</span>
                    </router-link>

                    <router-link to="/admin/plans" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polygon points="12 2 2 7 12 12 22 7 12 2" />
                            <polyline points="2 17 12 22 22 17" />
                            <polyline points="2 12 12 17 22 12" />
                        </svg>
                        <span>Plans</span>
                    </router-link>

                    <router-link to="/admin/catalogue" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="3" width="7" height="7" />
                            <rect x="14" y="3" width="7" height="7" />
                            <rect x="14" y="14" width="7" height="7" />
                            <rect x="3" y="14" width="7" height="7" />
                        </svg>
                        <span>Catalogue</span>
                    </router-link>

                    <router-link to="/admin/annonces" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4Z" />
                            <path d="M3 6h18" />
                            <path d="M16 10a4 4 0 0 1-8 0" />
                        </svg>
                        <span>Annonces</span>
                    </router-link>

                    <router-link to="/admin/publicites" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
                        </svg>
                        <span>Publicités</span>
                    </router-link>

                    <router-link to="/admin/commissions" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <line x1="12" y1="1" x2="12" y2="23" />
                            <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
                        </svg>
                        <span>Commissions</span>
                    </router-link>

                    <router-link to="/admin/financier" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polyline points="23 6 13.5 15.5 8.5 10.5 1 18" />
                            <polyline points="17 6 23 6 23 12" />
                        </svg>
                        <span>Financier</span>
                    </router-link>

                    <router-link to="/admin/notifications" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                            <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                        </svg>
                        <span>Notifications</span>
                    </router-link>

                    <router-link to="/admin/documents" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                            <polyline points="14 2 14 8 20 8" />
                        </svg>
                        <span>Documents</span>
                    </router-link>

                    <router-link to="/admin/plannings" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                            <line x1="16" y1="2" x2="16" y2="6" />
                            <line x1="8" y1="2" x2="8" y2="6" />
                            <line x1="3" y1="10" x2="21" y2="10" />
                        </svg>
                        <span>Plannings</span>
                    </router-link>

                    <router-link to="/admin/projets" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
                        </svg>
                        <span>Projets</span>
                    </router-link>

                    <router-link to="/admin/evenements" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                            <line x1="16" y1="2" x2="16" y2="6" />
                            <line x1="8" y1="2" x2="8" y2="6" />
                            <line x1="3" y1="10" x2="21" y2="10" />
                            <line x1="8" y1="14" x2="8" y2="14" />
                        </svg>
                        <span>Évènements</span>
                    </router-link>

                    <router-link to="/admin/actualites" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a4 4 0 0 1-4 4z" />
                            <line x1="10" y1="7" x2="16" y2="7" />
                            <line x1="10" y1="11" x2="16" y2="11" />
                            <line x1="10" y1="15" x2="14" y2="15" />
                        </svg>
                        <span>Actualités</span>
                    </router-link>

                    <router-link to="/admin/forum" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                        </svg>
                        <span>Forum</span>
                    </router-link>

                    <router-link to="/admin/langues" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <circle cx="12" cy="12" r="10" />
                            <line x1="2" y1="12" x2="22" y2="12" />
                            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                        </svg>
                        <span>Langues</span>
                    </router-link>

                    <router-link to="/admin/logs" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
                        </svg>
                        <span>Logs</span>
                    </router-link>

                    <router-link to="/admin/parametres" class="sidebar-item" active-class="sidebar-item--active">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <circle cx="12" cy="12" r="3" />
                            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
                        </svg>
                        <span>Paramètres</span>
                    </router-link>

                    <div class="sidebar-divider"></div>

                    <router-link to="/" class="sidebar-item sidebar-item--back">
                        <svg class="sidebar-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                            <polyline points="16 17 21 12 16 7" />
                            <line x1="21" y1="12" x2="9" y2="12" />
                        </svg>
                        <span>Retour au site</span>
                    </router-link>
                </nav>
            </aside>

            <main class="main-content">
                <router-view />
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
                        <span class="lang-sep"> - </span>
                        <span>Français</span>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<style scoped>
.admin-layout {
    background-color: var(--cream);
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
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
    flex-shrink: 0;
    letter-spacing: -0.01em;
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
    overflow-y: auto;
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
    padding: 12px 16px;
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
