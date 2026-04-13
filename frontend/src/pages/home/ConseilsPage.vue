<script setup lang="ts">
import { ref, computed } from 'vue'

const filterCategorie = ref('')

const articles = ref([
    {
        id: 1,
        titre: '5 techniques pour poncer et traiter le bois récupéré',
        categorie: 'Bois & Mobilier',
        auteur: 'Sophie M.',
        date: '10 avr. 2026',
        duree: '5 min',
        img: 'https://images.unsplash.com/photo-1504148455328-c376907d081c?w=600&q=80',
        extrait: "Le bois de récupération est une mine d'or pour les projets d'upcycling. Voici comment le préparer correctement avant de le transformer.",
    },
    {
        id: 2,
        titre: 'Comment identifier les palettes réutilisables en toute sécurité',
        categorie: 'Bois & Mobilier',
        auteur: 'Thomas L.',
        date: '7 avr. 2026',
        duree: '3 min',
        img: 'https://images.unsplash.com/photo-1565516965946-4f9046d42d07?w=600&q=80',
        extrait: "Toutes les palettes ne se valent pas. Apprenez à lire les marquages pour choisir celles qui sont sans danger pour votre intérieur.",
    },
    {
        id: 3,
        titre: 'Transformer un jean usé en sac à dos : guide étape par étape',
        categorie: 'Textile',
        auteur: 'Lucie D.',
        date: '3 avr. 2026',
        duree: '7 min',
        img: 'https://images.unsplash.com/photo-1602143407151-7111542de6e8?w=600&q=80',
        extrait: "Votre vieux jean a vécu ? Il peut devenir un sac solide et unique. Tout ce qu'il vous faut : une machine à coudre, du fil et un peu de patience.",
    },
    {
        id: 4,
        titre: 'Les 10 meilleurs sites pour trouver des matériaux gratuits en France',
        categorie: 'Ressources',
        auteur: 'Marc B.',
        date: '1 avr. 2026',
        duree: '4 min',
        img: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600&q=80',
        extrait: "Annonces, groupes Facebook, applications mobiles... La liste complète des endroits où trouver des matériaux de récupération sans débourser un euro.",
    },
    {
        id: 5,
        titre: 'Customiser ses meubles IKEA avec des matériaux récupérés',
        categorie: 'Bois & Mobilier',
        auteur: 'Camille R.',
        date: '28 mars 2026',
        duree: '6 min',
        img: 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=600&q=80',
        extrait: "Le phénomène IKEAhack passe au vert. Comment customiser vos meubles standards avec des chutes de bois, du cuir et du tissu de récupération.",
    },
    {
        id: 6,
        titre: 'Réparer plutôt que jeter : comment trouver les bonnes pièces',
        categorie: 'Réparation',
        auteur: 'Antoine V.',
        date: '25 mars 2026',
        duree: '5 min',
        img: 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=600&q=80',
        extrait: "Obsolescence programmée, pièces introuvables... Voici les astuces pour dénicher les pièces détachées dont vous avez besoin, neuves ou d'occasion.",
    },
])

const categories = ['Bois & Mobilier', 'Textile', 'Ressources', 'Réparation']

const articlesFiltres = computed(() => {
    if (!filterCategorie.value) return articles.value
    return articles.value.filter((a) => a.categorie === filterCategorie.value)
})
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">Conseils & inspirations.</h1>
                <p class="hero-subtitle">
                    Guides pratiques, tutoriels et idées pour réussir vos projets d'upcycling.
                </p>
            </div>
        </section>

        <section class="articles-section">
            <div class="container">
                <div class="filters-row">
                    <button
                        class="filter-btn"
                        :class="{ 'filter-btn--active': filterCategorie === '' }"
                        @click="filterCategorie = ''"
                    >
                        Tous
                    </button>
                    <button
                        v-for="cat in categories"
                        :key="cat"
                        class="filter-btn"
                        :class="{ 'filter-btn--active': filterCategorie === cat }"
                        @click="filterCategorie = cat"
                    >
                        {{ cat }}
                    </button>
                </div>

                <div class="articles-grid">
                    <article v-for="article in articlesFiltres" :key="article.id" class="article-card">
                        <div class="article-img-wrap">
                            <img :src="article.img" :alt="article.titre" class="article-img" />
                            <span class="article-cat">{{ article.categorie }}</span>
                        </div>
                        <div class="article-body">
                            <h2 class="article-titre">{{ article.titre }}</h2>
                            <p class="article-extrait">{{ article.extrait }}</p>
                            <div class="article-meta">
                                <span class="article-auteur">{{ article.auteur }}</span>
                                <span class="article-sep">·</span>
                                <span class="article-date">{{ article.date }}</span>
                                <span class="article-sep">·</span>
                                <span class="article-duree">{{ article.duree }} de lecture</span>
                            </div>
                            <button class="btn-lire">Lire l'article</button>
                        </div>
                    </article>
                </div>

                <div v-if="articlesFiltres.length === 0" class="empty-state">
                    <p>Aucun article dans cette catégorie pour le moment.</p>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page-content {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}

.hero {
    padding: 64px 0 48px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.6rem, 5.5vw, 4.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0 0 16px;
}
.hero-subtitle {
    font-size: 1rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
}

.articles-section {
    flex: 1;
    padding: 0 0 80px;
}

.filters-row {
    display: flex;
    gap: 8px;
    margin-bottom: 36px;
    flex-wrap: wrap;
}
.filter-btn {
    padding: 8px 18px;
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 20px;
    background: transparent;
    color: var(--charcoal);
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    font-family: inherit;
    transition:
        background 0.2s,
        border-color 0.2s,
        color 0.2s;
}
.filter-btn:hover {
    border-color: var(--green-mid);
    color: var(--green-mid);
}
.filter-btn--active {
    background: var(--green-dark);
    border-color: var(--green-dark);
    color: var(--white);
}
.filter-btn--active:hover {
    color: var(--white);
}

.articles-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
}

.article-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 10px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    transition:
        transform 0.2s,
        box-shadow 0.2s;
}
.article-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(8, 106, 53, 0.1);
}

.article-img-wrap {
    position: relative;
    width: 100%;
    aspect-ratio: 16/9;
    overflow: hidden;
}
.article-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.4s ease;
}
.article-card:hover .article-img {
    transform: scale(1.04);
}
.article-cat {
    position: absolute;
    top: 12px;
    left: 12px;
    background: var(--green-dark);
    color: var(--white);
    padding: 4px 10px;
    border-radius: 5px;
    font-size: 0.72rem;
    font-weight: 700;
}

.article-body {
    padding: 18px 20px 20px;
    display: flex;
    flex-direction: column;
    flex: 1;
}

.article-titre {
    font-size: 0.95rem;
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.4;
    margin: 0 0 10px;
}

.article-extrait {
    font-size: 0.82rem;
    color: var(--charcoal);
    opacity: 0.72;
    line-height: 1.6;
    margin: 0 0 16px;
    flex: 1;
}

.article-meta {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.6;
    margin-bottom: 14px;
    flex-wrap: wrap;
}
.article-auteur {
    font-weight: 600;
    opacity: 1;
}
.article-sep {
    opacity: 0.4;
}

.btn-lire {
    width: 100%;
    padding: 10px;
    background: var(--green-pale);
    color: var(--green-dark);
    border: 1.5px solid rgba(8, 106, 53, 0.2);
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        background 0.2s,
        color 0.2s;
    margin-top: auto;
}
.btn-lire:hover {
    background: var(--green-dark);
    color: var(--white);
    border-color: var(--green-dark);
}

.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
}

@media (max-width: 860px) {
    .articles-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
@media (max-width: 560px) {
    .articles-grid {
        grid-template-columns: 1fr;
    }
    .hero-title {
        font-size: 2.2rem;
    }
}
</style>
