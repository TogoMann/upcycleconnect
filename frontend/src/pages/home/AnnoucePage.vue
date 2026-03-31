<script setup lang="ts">
import { ref, computed } from 'vue'

// ── Filtres
const searchQuery = ref('')
const filterType = ref('')
const filterCategorie = ref('')
const filterLocalisation = ref('')

// ── Données des annonces
const annonces = ref([
    {
        id: 1,
        titre: 'Chaise en chêne',
        description: 'Ponçage et vernis mat ...',
        type: 'Don',
        prix: null,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Mobilier',
        img: 'https://images.unsplash.com/photo-1592078615290-033ee584e267?w=400&q=80',
    },
    {
        id: 2,
        titre: 'Vase',
        description: 'Ponçage et vernis mat ...',
        type: 'Vente',
        prix: 25,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Décoration',
        img: 'https://images.unsplash.com/photo-1612196808214-b8e1d6145a8c?w=400&q=80',
    },
    {
        id: 3,
        titre: 'Panier à fruits',
        description: 'Ponçage et vernis mat ...',
        type: 'Don',
        prix: null,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Cuisine',
        img: 'https://images.unsplash.com/photo-1611735341450-74d61e660ad2?w=400&q=80',
    },
    {
        id: 4,
        titre: 'Pull gris',
        description: 'Ponçage et vernis mat ...',
        type: 'Don',
        prix: null,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Vêtement',
        img: 'https://images.unsplash.com/photo-1620799140408-edc6dcb6d633?w=400&q=80',
    },
    {
        id: 5,
        titre: 'Lampe dorée',
        description: 'Ponçage et vernis mat ...',
        type: 'Vente',
        prix: 4,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Décoration',
        img: 'https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=400&q=80',
    },
    {
        id: 6,
        titre: 'Tableau enfant',
        description: 'Ponçage et vernis mat ...',
        type: 'Don',
        prix: null,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Jouet',
        img: 'https://images.unsplash.com/photo-1596464716127-f2a82984de30?w=400&q=80',
    },
    {
        id: 7,
        titre: 'Chaise en chêne',
        description: 'Ponçage et vernis mat ...',
        type: 'Vente',
        prix: 7,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Mobilier',
        img: 'https://images.unsplash.com/photo-1598300042247-d088f8ab3a91?w=400&q=80',
    },
    {
        id: 8,
        titre: 'Table basse',
        description: 'Ponçage et vernis mat ...',
        type: 'Vente',
        prix: 40,
        localisation: 'Nantes',
        date: '14 Mai',
        vendeur: 'Mégane K.',
        categorie: 'Mobilier',
        img: 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400&q=80',
    },
])

// ── Annonces filtrées
const annoncesFiltrees = computed(() => {
    return annonces.value.filter((a) => {
        const matchSearch =
            !searchQuery.value || a.titre.toLowerCase().includes(searchQuery.value.toLowerCase())
        const matchType =
            !filterType.value ||
            (filterType.value === 'Don' && a.type === 'Don') ||
            (filterType.value === 'Vente' && a.type === 'Vente')
        const matchCat = !filterCategorie.value || a.categorie === filterCategorie.value
        const matchLoc = !filterLocalisation.value || a.localisation === filterLocalisation.value
        return matchSearch && matchType && matchCat && matchLoc
    })
})

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
                <h1 class="hero-title">A vendre ou à donner.</h1>
            </div>
        </section>

        <!-- ════════ BARRE DE RECHERCHE ════════ -->
        <section class="search-section">
            <div class="container">
                <!-- Search input -->
                <div class="search-bar">
                    <input
                        v-model="searchQuery"
                        type="text"
                        placeholder="Rechercher un objet..."
                        class="search-input"
                    />
                    <button class="search-btn" aria-label="Rechercher">
                        <svg
                            width="18"
                            height="18"
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

                <!-- Filtres -->
                <div class="filters-row">
                    <span class="filters-label">Type d'offre</span>
                    <div class="select-wrap">
                        <select v-model="filterType" class="filter-select">
                            <option value="">Don/Vente</option>
                            <option value="Don">Don</option>
                            <option value="Vente">Vente</option>
                        </select>
                        <span class="chevron">&#8964;</span>
                    </div>
                    <div class="select-wrap">
                        <select v-model="filterCategorie" class="filter-select">
                            <option value="">Catégorie</option>
                            <option>Mobilier</option>
                            <option>Décoration</option>
                            <option>Vêtement</option>
                            <option>Cuisine</option>
                            <option>Jouet</option>
                            <option>Électronique</option>
                            <option>Outil</option>
                        </select>
                        <span class="chevron">&#8964;</span>
                    </div>
                    <div class="select-wrap">
                        <select v-model="filterLocalisation" class="filter-select">
                            <option value="">Localisation</option>
                            <option>Nantes</option>
                            <option>Paris</option>
                            <option>Lyon</option>
                            <option>Bordeaux</option>
                            <option>Marseille</option>
                        </select>
                        <span class="chevron">&#8964;</span>
                    </div>
                </div>
            </div>
        </section>

        <!-- ════════ GRILLE D'ANNONCES ════════ -->
        <section class="annonces-section">
            <div class="container">
                <div class="annonces-grid">
                    <div v-for="annonce in annoncesFiltrees" :key="annonce.id" class="annonce-card">
                        <!-- Image -->
                        <div class="annonce-img-wrap">
                            <img :src="annonce.img" :alt="annonce.titre" class="annonce-img" />
                        </div>

                        <!-- Infos -->
                        <div class="annonce-body">
                            <!-- Badge type -->
                            <div class="annonce-badges">
                                <span
                                    class="badge"
                                    :class="annonce.type === 'Don' ? 'badge--don' : 'badge--vente'"
                                >
                                    {{ annonce.type }}
                                    <template v-if="annonce.type === 'Vente' && annonce.prix">
                                        &nbsp;- {{ annonce.prix }}€
                                    </template>
                                </span>
                            </div>

                            <h3 class="annonce-titre">{{ annonce.titre }}</h3>
                            <p class="annonce-desc">{{ annonce.description }}</p>

                            <!-- Localisation + date -->
                            <div class="annonce-meta">
                                <svg
                                    class="pin-icon"
                                    viewBox="0 0 24 24"
                                    fill="currentColor"
                                    width="13"
                                    height="13"
                                >
                                    <path
                                        d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
                                    />
                                </svg>
                                <span>{{ annonce.localisation }} - {{ annonce.date }}</span>
                            </div>

                            <p class="annonce-vendeur">{{ annonce.vendeur }}</p>

                            <!-- CTA -->
                            <button class="btn-annonce">Voir l'annonce</button>
                        </div>
                    </div>
                </div>

                <!-- État vide -->
                <div v-if="annoncesFiltrees.length === 0" class="empty-state">
                    <p>Aucune annonce ne correspond à votre recherche.</p>
                </div>
            </div>
        </section>

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
.nav-link:hover,
.nav-link.active {
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

/* ══ HERO ══ */
.hero {
    padding: 64px 0 36px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.6rem, 6vw, 4.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0;
}

/* ══ RECHERCHE ══ */
.search-section {
    padding: 0 0 32px;
}

/* Barre de recherche */
.search-bar {
    display: flex;
    align-items: center;
    max-width: 520px;
    margin: 0 auto 28px;
    background: var(--white);
    border: 2px solid var(--charcoal);
    border-radius: 8px;
    overflow: hidden;
    transition:
        border-color 0.2s,
        box-shadow 0.2s;
}
.search-bar:focus-within {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.15);
}
.search-input {
    flex: 1;
    padding: 14px 18px;
    font-size: 0.95rem;
    color: var(--charcoal);
    background: transparent;
    border: none;
    outline: none;
    font-family: inherit;
}
.search-input::placeholder {
    color: rgba(53, 53, 53, 0.45);
}
.search-btn {
    background: var(--green-dark);
    color: var(--white);
    border: none;
    padding: 0 18px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background 0.2s;
    min-height: 50px;
}
.search-btn:hover {
    background: var(--green-mid);
}

/* Ligne filtres */
.filters-row {
    display: flex;
    align-items: center;
    gap: 12px;
    justify-content: center;
    flex-wrap: wrap;
}
.filters-label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
    margin-right: 4px;
}
.select-wrap {
    position: relative;
    display: inline-flex;
    align-items: center;
}
.filter-select {
    appearance: none;
    -webkit-appearance: none;
    padding: 10px 36px 10px 16px;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--charcoal);
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    cursor: pointer;
    font-family: inherit;
    outline: none;
    transition: border-color 0.2s;
    min-width: 130px;
}
.filter-select:focus {
    border-color: var(--green-mid);
}
.chevron {
    position: absolute;
    right: 12px;
    font-size: 1.1rem;
    color: var(--charcoal);
    pointer-events: none;
    line-height: 1;
}

/* ══ GRILLE ANNONCES ══ */
.annonces-section {
    flex: 1;
    padding: 0 0 80px;
}
.annonces-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
}

/* Carte annonce */
.annonce-card {
    background: var(--green-pale);
    border-radius: 10px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition:
        transform 0.2s,
        box-shadow 0.2s;
}
.annonce-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(8, 106, 53, 0.12);
}

/* Image */
.annonce-img-wrap {
    width: 100%;
    aspect-ratio: 4/3;
    overflow: hidden;
}
.annonce-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.4s ease;
}
.annonce-card:hover .annonce-img {
    transform: scale(1.05);
}

/* Corps */
.annonce-body {
    padding: 12px 14px 16px;
    display: flex;
    flex-direction: column;
    flex: 1;
}

/* Badge Don / Vente */
.annonce-badges {
    margin-bottom: 8px;
}
.badge {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 5px;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.02em;
}
.badge--don {
    background: var(--green-mid);
    color: var(--white);
}
.badge--vente {
    background: var(--green-mid);
    color: var(--white);
}

/* Titre */
.annonce-titre {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 4px;
    line-height: 1.3;
}

/* Description */
.annonce-desc {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0 0 8px;
    line-height: 1.4;
}

/* Localisation */
.annonce-meta {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.75;
    margin-bottom: 3px;
}
.pin-icon {
    color: var(--green-mid);
    flex-shrink: 0;
}

/* Vendeur */
.annonce-vendeur {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0 0 12px;
}

/* Bouton */
.btn-annonce {
    width: 100%;
    padding: 10px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    margin-top: auto;
    transition: background 0.2s;
}
.btn-annonce:hover {
    background: var(--green-mid);
}

/* État vide */
.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
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
}
.footer-lang {
    font-size: 0.85rem;
    color: rgba(255, 255, 255, 0.75);
}
.lang-sep {
    opacity: 0.5;
}

/* ══ RESPONSIVE ══ */
@media (max-width: 900px) {
    .annonces-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .nav-links {
        display: none;
    }
}
@media (max-width: 560px) {
    .annonces-grid {
        grid-template-columns: 1fr;
    }
    .hero-title {
        font-size: 2.2rem;
    }
    .filters-row {
        flex-direction: column;
        align-items: flex-start;
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
