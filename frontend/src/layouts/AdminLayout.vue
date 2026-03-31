<template>
    <div class="admin-layout">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <div class="search-bar">
                    <input type="text" placeholder="Rechercher un objet..." class="search-input" />
                    <button class="search-btn" aria-label="Rechercher">
                        <svg
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <circle cx="11" cy="11" r="8" />
                            <line x1="21" y1="21" x2="16.65" y2="16.65" />
                        </svg>
                    </button>
                </div>
                <div class="nav-profile">
                    <span class="profile-name">Admin</span>
                    <span class="profile-role">Administrateur</span>
                    <svg
                        class="profile-chevron"
                        width="14"
                        height="14"
                        viewBox="0 0 12 12"
                        fill="currentColor"
                    >
                        <path d="M6 8L1 3h10z" />
                    </svg>
                </div>
            </div>
        </header>

        <div class="layout-body">
            <aside class="sidebar">
                <nav class="sidebar-nav">
                    <router-link
                        to="/admin"
                        class="sidebar-item"
                        exact-active-class="sidebar-item--active"
                    >
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <rect x="3" y="12" width="4" height="9" />
                            <rect x="10" y="7" width="4" height="14" />
                            <rect x="17" y="3" width="4" height="18" />
                        </svg>
                        <span>Dashboard</span>
                    </router-link>
                    <router-link
                        to="/admin/users"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                            <circle cx="9" cy="7" r="4" />
                            <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                        </svg>
                        <span>Utilisateurs</span>
                    </router-link>
                    <router-link
                        to="/admin/courses"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                            <polyline points="14 2 14 8 20 8" />
                        </svg>
                        <span>Annonces</span>
                    </router-link>
                    <router-link
                        to="/admin/events"
                        class="sidebar-item"
                        active-class="sidebar-item--active"
                    >
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                            <line x1="16" y1="2" x2="16" y2="6" />
                            <line x1="8" y1="2" x2="8" y2="6" />
                            <line x1="3" y1="10" x2="21" y2="10" />
                        </svg>
                        <span>Événements</span>
                    </router-link>

                    <div class="sidebar-divider"></div>

                    <router-link to="/" class="sidebar-item sidebar-item--back">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
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

<script setup lang="ts"></script>

<style scoped>
.admin-layout {
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
    gap: 24px;
}
.nav-logo {
    font-weight: 800;
    font-size: 1.1rem;
    color: var(--green-dark);
    text-decoration: none;
    flex-shrink: 0;
    letter-spacing: -0.01em;
}
.search-bar {
    flex: 1;
    max-width: 440px;
    display: flex;
    align-items: center;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    overflow: hidden;
    transition: border-color 0.2s;
}
.search-bar:focus-within {
    border-color: var(--green-mid);
}
.search-input {
    flex: 1;
    padding: 10px 14px;
    font-size: 0.9rem;
    background: transparent;
    border: none;
    outline: none;
    font-family: inherit;
    color: var(--charcoal);
}
.search-input::placeholder {
    color: rgba(53, 53, 53, 0.4);
}
.search-btn {
    background: var(--green-dark);
    color: var(--white);
    border: none;
    padding: 0 14px;
    min-height: 42px;
    display: flex;
    align-items: center;
    cursor: pointer;
    transition: background 0.2s;
}
.search-btn:hover {
    background: var(--green-mid);
}
.nav-profile {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 0;
    margin-left: auto;
    cursor: pointer;
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
.profile-chevron {
    color: var(--green-dark);
    margin-top: 2px;
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
    padding: 12px 16px;
    border-radius: 8px;
    text-decoration: none;
    font-size: 0.95rem;
    color: var(--charcoal);
    font-weight: 500;
    transition:
        background 0.15s,
        color 0.15s;
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
    width: 20px;
    height: 20px;
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
