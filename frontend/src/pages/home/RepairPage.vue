<script setup lang="ts">
import { ref } from 'vue'

const typeObjet = ref('')
const codePostal = ref('')

const typesObjets = [
    'Électronique',
    'Mobilier',
    'Vêtement',
    'Électroménager',
    'Jouet',
    'Outil',
    'Autre',
]

const codesPostaux = [
    '75001',
    '75002',
    '75003',
    '75004',
    '75005',
    '75006',
    '75007',
    '75008',
    '75009',
    '75010',
    '75011',
    '75012',
    '75013',
    '75014',
    '75015',
    '75016',
    '75017',
    '75018',
    '75019',
    '75020',
    '69001',
    '13001',
    '31000',
    '33000',
    '59000',
    '67000',
    '06000',
    '44000',
    '34000',
    '76000',
]

const handleRecherche = () => {
    console.log('Recherche :', { typeObjet: typeObjet.value, codePostal: codePostal.value })
}

const footerLinks = ['À propos', 'Mentions légales', 'Politique de confidentialité']
</script>

<template>
    <div class="page">
        <!-- ════════ NAVBAR ════════ -->
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

        <!-- ════════ HERO ════════ -->
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">Réparer vos objets.</h1>
            </div>
        </section>

        <!-- ════════ FORMULAIRE DE RECHERCHE ════════ -->
        <section class="search-section">
            <div class="container">
                <div class="search-card">
                    <!-- En-tête de la carte -->
                    <div class="search-card-header">
                        <h2 class="search-card-title">Trouver une ressource de réparation</h2>
                    </div>

                    <!-- Corps de la carte -->
                    <div class="search-card-body">
                        <!-- Ligne avec les 2 selects -->
                        <div class="search-fields">
                            <div class="field-group">
                                <label class="field-label">Type d'objet</label>
                                <div class="select-wrap">
                                    <select v-model="typeObjet" class="field-select">
                                        <option value="" disabled>
                                            Électronique, Mobilier, Vêtement...
                                        </option>
                                        <option v-for="t in typesObjets" :key="t" :value="t">
                                            {{ t }}
                                        </option>
                                    </select>
                                    <span class="select-arrow">&#8964;</span>
                                </div>
                            </div>

                            <div class="field-group field-group--sm">
                                <label class="field-label">Code Postal</label>
                                <div class="select-wrap">
                                    <select v-model="codePostal" class="field-select">
                                        <option value="" disabled>75001</option>
                                        <option v-for="cp in codesPostaux" :key="cp" :value="cp">
                                            {{ cp }}
                                        </option>
                                    </select>
                                    <span class="select-arrow">&#8964;</span>
                                </div>
                            </div>
                        </div>

                        <!-- Bouton Rechercher -->
                        <button class="btn-search" @click="handleRecherche">Rechercher</button>
                    </div>
                </div>
            </div>
        </section>

        <!-- Espace pour les résultats futurs -->
        <div class="results-placeholder" />

        <!-- ════════ FOOTER ════════ -->
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

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
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
.nav-link.active {
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

/* ══ HERO ══ */
.hero {
    padding: 72px 0 48px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.8rem, 6vw, 4.4rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0;
}

/* ══ FORMULAIRE ══ */
.search-section {
    padding: 0 0 40px;
}

.search-card {
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 24px rgba(8, 106, 53, 0.1);
}

/* En-tête verte foncée */
.search-card-header {
    background: var(--green-mid);
    padding: 24px 28px;
}
.search-card-title {
    font-size: 1.35rem;
    font-weight: 700;
    color: var(--white);
    margin: 0;
    letter-spacing: -0.01em;
}

/* Corps vert clair */
.search-card-body {
    background: var(--green-light);
    padding: 28px 28px 28px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

/* Ligne des 2 selects */
.search-fields {
    display: flex;
    gap: 14px;
    align-items: flex-end;
}

.field-group {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.field-group--sm {
    flex: 0 0 220px;
}

.field-label {
    font-size: 1.15rem;
    font-weight: 700;
    color: rgba(255, 255, 255, 0.85);
    letter-spacing: -0.01em;
}

/* Select custom */
.select-wrap {
    position: relative;
    display: flex;
    align-items: center;
}
.field-select {
    width: 100%;
    padding: 16px 48px 16px 18px;
    font-size: 1.05rem;
    font-weight: 400;
    color: var(--charcoal);
    background: var(--cream);
    border: 2px solid rgba(255, 255, 255, 0.4);
    border-radius: 8px;
    appearance: none;
    -webkit-appearance: none;
    cursor: pointer;
    font-family: inherit;
    outline: none;
    transition:
        border-color 0.2s,
        box-shadow 0.2s;
}
.field-select:focus {
    border-color: var(--green-dark);
    box-shadow: 0 0 0 3px rgba(8, 106, 53, 0.15);
}
.field-select option[value=''][disabled] {
    color: rgba(53, 53, 53, 0.45);
}
.select-arrow {
    position: absolute;
    right: 16px;
    font-size: 1.3rem;
    color: var(--charcoal);
    pointer-events: none;
    line-height: 1;
}

/* Bouton Rechercher */
.btn-search {
    width: 100%;
    padding: 22px;
    background: var(--green-mid);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 1.35rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    letter-spacing: -0.01em;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-search:hover {
    background: var(--green-dark);
    transform: translateY(-1px);
}
.btn-search:active {
    transform: translateY(0);
}

/* Espace résultats */
.results-placeholder {
    flex: 1;
    min-height: 120px;
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

/* ══ RESPONSIVE ══ */
@media (max-width: 700px) {
    .search-fields {
        flex-direction: column;
    }
    .field-group--sm {
        flex: none;
        width: 100%;
    }
    .nav-links {
        display: none;
    }
    .hero-title {
        font-size: 2.4rem;
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
