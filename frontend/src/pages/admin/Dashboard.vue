<template>
    <div class="page">
        <!-- ════════ NAVBAR ADMIN ════════ -->
        <header class="navbar">
            <div class="nav-container">
                <!-- Logo -->
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>

                <!-- Barre de recherche centrale -->
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

                <!-- Profil admin -->
                <div class="nav-profile">
                    <span class="profile-name">Marie</span>
                    <span class="profile-role">Admin</span>
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

        <!-- ════════ LAYOUT : SIDEBAR + MAIN ════════ -->
        <div class="layout">
            <!-- ── SIDEBAR ── -->
            <aside class="sidebar">
                <nav class="sidebar-nav">
                    <router-link to="/backoffice/revenus" class="sidebar-item">
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
                        <span>Revenus</span>
                    </router-link>
                    <router-link to="/backoffice/commissions" class="sidebar-item">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <circle cx="9" cy="21" r="1" />
                            <circle cx="20" cy="21" r="1" />
                            <path
                                d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"
                            />
                        </svg>
                        <span>Commissions</span>
                    </router-link>
                    <router-link to="/backoffice/inventaire" class="sidebar-item">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
                            <line x1="3" y1="6" x2="21" y2="6" />
                            <path d="M16 10a4 4 0 0 1-8 0" />
                        </svg>
                        <span>Inventaire</span>
                    </router-link>
                    <router-link to="/backoffice/utilisateurs" class="sidebar-item active">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
                        </svg>
                        <span>Utilisateurs</span>
                    </router-link>
                    <router-link to="/backoffice/publicites" class="sidebar-item">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <path
                                d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                            />
                        </svg>
                        <span>Publicités</span>
                    </router-link>
                    <router-link to="/backoffice/notifications" class="sidebar-item">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <circle cx="12" cy="12" r="3" />
                            <path
                                d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"
                            />
                        </svg>
                        <span>Notifications</span>
                    </router-link>
                    <router-link to="/backoffice/logistique" class="sidebar-item">
                        <svg
                            class="sidebar-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="1.8"
                        >
                            <path
                                d="M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2V9M9 21H5a2 2 0 0 1-2-2V9m0 0h18"
                            />
                        </svg>
                        <span>Logistique</span>
                    </router-link>
                </nav>
            </aside>

            <!-- ── CONTENU PRINCIPAL ── -->
            <main class="main-content">
                <!-- Titre -->
                <h1 class="page-title">Tableau de bord.</h1>

                <!-- ── CARTE DONNÉES DU JOUR ── -->
                <div class="data-card">
                    <div class="data-card-header">
                        <div>
                            <h2 class="data-card-title">Données du jour</h2>
                            <p class="data-card-sub">Résumé des ventes</p>
                        </div>
                        <button class="btn-export">
                            <svg
                                width="14"
                                height="14"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                                <polyline points="7 10 12 15 17 10" />
                                <line x1="12" y1="15" x2="12" y2="3" />
                            </svg>
                            Export
                        </button>
                    </div>

                    <!-- KPI Cards — données du store + données simulées -->
                    <div class="kpi-grid">
                        <!-- Revenus du mois (simulé, à brancher sur adminStore) -->
                        <div class="kpi-card">
                            <div class="kpi-icon">
                                <svg
                                    width="20"
                                    height="20"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <rect x="3" y="12" width="4" height="9" />
                                    <rect x="10" y="7" width="4" height="14" />
                                    <rect x="17" y="3" width="4" height="18" />
                                </svg>
                            </div>
                            <div class="kpi-value">11 450 €</div>
                            <div class="kpi-label">Revenus du mois</div>
                        </div>

                        <!-- Annonces à valider = courses (depuis adminStore) -->
                        <div class="kpi-card">
                            <div class="kpi-icon">
                                <svg
                                    width="20"
                                    height="20"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path
                                        d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
                                    />
                                    <polyline points="14 2 14 8 20 8" />
                                </svg>
                            </div>
                            <div class="kpi-value">{{ adminStore.courses.length }}</div>
                            <div class="kpi-label">Annonces à valider</div>
                        </div>

                        <!-- Abonnements actifs = users (depuis adminStore) -->
                        <div class="kpi-card">
                            <div class="kpi-icon">
                                <svg
                                    width="20"
                                    height="20"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path d="M12 20h9" />
                                    <path
                                        d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"
                                    />
                                </svg>
                            </div>
                            <div class="kpi-value">{{ adminStore.users.length }}</div>
                            <div class="kpi-label">Abonnements actifs</div>
                        </div>

                        <!-- Box disponibles = events (depuis adminStore) -->
                        <div class="kpi-card">
                            <div class="kpi-icon">
                                <svg
                                    width="20"
                                    height="20"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                >
                                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                                    <circle cx="9" cy="7" r="4" />
                                    <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                                    <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                                </svg>
                            </div>
                            <div class="kpi-value">{{ adminStore.events.length }}</div>
                            <div class="kpi-label">Box disponibles</div>
                        </div>
                    </div>
                </div>

                <!-- ── TABLEAU ACTIVITÉS ── -->
                <div class="activity-card">
                    <h2 class="activity-title">Dernières activités de validation</h2>
                    <table class="activity-table">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>Type</th>
                                <th>Date</th>
                                <th>Statut</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td class="num">01</td>
                                <td>Don</td>
                                <td>16/06/26</td>
                                <td><span class="statut statut--attente">En attente</span></td>
                            </tr>
                            <tr>
                                <td class="num">02</td>
                                <td>Vente</td>
                                <td>13/06/26</td>
                                <td><span class="statut statut--rejete">Rejeté</span></td>
                            </tr>
                            <tr>
                                <td class="num">03</td>
                                <td>Atelier</td>
                                <td>13/06/26</td>
                                <td><span class="statut statut--valide">Validé</span></td>
                            </tr>
                            <tr class="last-row">
                                <td class="num">04</td>
                                <td>Don</td>
                                <td>11/06/26</td>
                                <td><span class="statut statut--valide">Validé</span></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </main>
        </div>

        <!-- ════════ FOOTER ════════ -->
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

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

// ── Logique existante conservée intacte ──
const adminStore = useAdminStore()
onMounted(() => {
    adminStore.fetchUsers()
    adminStore.fetchCourses()
    adminStore.fetchEvents()
})
</script>

<style scoped>
/* ══ Charte graphique UCC ══ */
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

/* ══ NAVBAR ══ */
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

/* Barre de recherche */
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

/* Profil admin */
.nav-profile {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-left: auto;
    cursor: pointer;
}
.profile-name {
    font-size: 0.95rem;
    font-weight: 700;
    color: var(--green-dark);
}
.profile-role {
    font-size: 0.78rem;
    color: var(--green-light);
    display: block;
    line-height: 1;
    margin-top: -2px;
}
.nav-profile {
    flex-direction: column;
    align-items: flex-end;
    gap: 0;
}
.profile-chevron {
    color: var(--green-dark);
    margin-top: -14px;
    margin-left: 4px;
}

/* ══ LAYOUT ══ */
.layout {
    display: flex;
    flex: 1;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
    padding: 0 32px;
    gap: 0;
}

/* ── SIDEBAR ── */
.sidebar {
    flex: 0 0 220px;
    padding: 40px 0 40px;
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
.sidebar-item.active {
    color: var(--green-dark);
    font-weight: 600;
}
.sidebar-icon {
    width: 20px;
    height: 20px;
    flex-shrink: 0;
    color: var(--charcoal);
}
.sidebar-item:hover .sidebar-icon,
.sidebar-item.active .sidebar-icon {
    color: var(--green-dark);
}

/* ── CONTENU PRINCIPAL ── */
.main-content {
    flex: 1;
    padding: 40px 0 60px 32px;
    min-width: 0;
}

.page-title {
    font-size: clamp(2.4rem, 4vw, 3.4rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 32px;
    line-height: 1.08;
}

/* ── CARTE DONNÉES DU JOUR ── */
.data-card {
    background: var(--green-dark);
    border-radius: 14px;
    padding: 28px 28px 24px;
    margin-bottom: 24px;
}
.data-card-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 24px;
}
.data-card-title {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--white);
    margin: 0 0 4px;
}
.data-card-sub {
    font-size: 0.85rem;
    color: rgba(255, 255, 255, 0.65);
    margin: 0;
}
.btn-export {
    display: flex;
    align-items: center;
    gap: 7px;
    background: transparent;
    color: var(--white);
    border: 1.5px solid rgba(255, 255, 255, 0.5);
    padding: 9px 18px;
    border-radius: 7px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        background 0.2s,
        border-color 0.2s;
    flex-shrink: 0;
}
.btn-export:hover {
    background: rgba(255, 255, 255, 0.1);
    border-color: var(--white);
}

/* KPI Grid */
.kpi-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 14px;
}
.kpi-card {
    background: var(--green-pale);
    border-radius: 10px;
    padding: 20px 18px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.kpi-icon {
    width: 38px;
    height: 38px;
    background: var(--green-mid);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--white);
    margin-bottom: 4px;
}
.kpi-value {
    font-size: 1.6rem;
    font-weight: 800;
    color: var(--green-dark);
    line-height: 1;
    letter-spacing: -0.02em;
}
.kpi-label {
    font-size: 0.8rem;
    color: var(--green-dark);
    font-weight: 500;
    opacity: 0.8;
}

/* ── TABLEAU ACTIVITÉS ── */
.activity-card {
    background: var(--green-pale);
    border-radius: 14px;
    padding: 24px 28px;
}
.activity-title {
    font-size: 1.05rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 20px;
}
.activity-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.875rem;
}
.activity-table thead tr {
    border-bottom: 1px solid rgba(53, 53, 53, 0.12);
}
.activity-table th {
    text-align: left;
    padding: 0 16px 12px 0;
    font-size: 0.78rem;
    font-weight: 700;
    color: var(--green-mid);
    letter-spacing: 0.02em;
    text-transform: capitalize;
}
.activity-table th:first-child {
    width: 60px;
}
.activity-table td {
    padding: 16px 16px 16px 0;
    color: var(--charcoal);
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
}
.activity-table tr.last-row td {
    border-bottom: none;
}
.activity-table td.num {
    font-weight: 600;
    color: rgba(53, 53, 53, 0.55);
    font-size: 0.82rem;
}

/* Statuts */
.statut {
    font-size: 0.82rem;
    font-weight: 500;
}
.statut--attente {
    color: var(--charcoal);
}
.statut--rejete {
    color: #c0392b;
}
.statut--valide {
    color: var(--green-mid);
}

/* ══ FOOTER ══ */
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

/* ══ RESPONSIVE ══ */
@media (max-width: 960px) {
    .kpi-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .sidebar {
        flex: 0 0 180px;
    }
}
@media (max-width: 700px) {
    .layout {
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
        gap: 8px;
    }
    .sidebar-item {
        padding: 8px 12px;
        font-size: 0.8rem;
    }
    .main-content {
        padding: 24px 0 40px;
    }
    .kpi-grid {
        grid-template-columns: repeat(2, 1fr);
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
