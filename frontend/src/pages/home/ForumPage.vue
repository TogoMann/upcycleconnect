<script setup lang="ts">
import { ref, computed } from 'vue'

const searchQuery = ref('')
const filterCategorie = ref('')

const sujets = ref([
    {
        id: 1,
        titre: 'Comment upcycler des palettes de bois en mobilier de jardin ?',
        categorie: 'Bricolage',
        auteur: 'Sophie M.',
        avatar: 'SM',
        reponses: 14,
        vues: 312,
        date: '12 avr. 2026',
        derniereActivite: 'il y a 2h',
    },
    {
        id: 2,
        titre: 'Partage de bons plans pour trouver des matériaux gratuits à Paris',
        categorie: 'Ressources',
        auteur: 'Thomas L.',
        avatar: 'TL',
        reponses: 8,
        vues: 187,
        date: '10 avr. 2026',
        derniereActivite: 'il y a 5h',
    },
    {
        id: 3,
        titre: 'Quels outils indispensables pour débuter dans le upcycling ?',
        categorie: 'Débutants',
        auteur: 'Camille R.',
        avatar: 'CR',
        reponses: 23,
        vues: 541,
        date: '8 avr. 2026',
        derniereActivite: 'il y a 1j',
    },
    {
        id: 4,
        titre: 'Retour d\'expérience : transformer des vêtements en accessoires de mode',
        categorie: 'Textile',
        auteur: 'Lucie D.',
        avatar: 'LD',
        reponses: 6,
        vues: 98,
        date: '7 avr. 2026',
        derniereActivite: 'il y a 2j',
    },
    {
        id: 5,
        titre: 'Trouver des artisans spécialisés dans ma région — comment faire ?',
        categorie: 'Communauté',
        auteur: 'Marc B.',
        avatar: 'MB',
        reponses: 11,
        vues: 224,
        date: '5 avr. 2026',
        derniereActivite: 'il y a 3j',
    },
    {
        id: 6,
        titre: 'Projet : lampe industrielle à partir de tuyaux de plomberie',
        categorie: 'Bricolage',
        auteur: 'Antoine V.',
        avatar: 'AV',
        reponses: 19,
        vues: 403,
        date: '3 avr. 2026',
        derniereActivite: 'il y a 4j',
    },
])

const categories = ['Bricolage', 'Textile', 'Ressources', 'Débutants', 'Communauté']

const sujetsFiltres = computed(() => {
    return sujets.value.filter((s) => {
        const matchSearch =
            !searchQuery.value ||
            s.titre.toLowerCase().includes(searchQuery.value.toLowerCase())
        const matchCat = !filterCategorie.value || s.categorie === filterCategorie.value
        return matchSearch && matchCat
    })
})
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container hero-inner">
                <div class="hero-text">
                    <h1 class="hero-title">Forum communautaire.</h1>
                    <p class="hero-subtitle">
                        Échangez, partagez et apprenez avec la communauté UpCycleConnect.
                    </p>
                </div>
                <router-link to="/auth/login" class="btn-nouveau">+ Nouveau sujet</router-link>
            </div>
        </section>

        <section class="forum-section">
            <div class="container">
                <div class="forum-toolbar">
                    <div class="search-bar">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="11" cy="11" r="8" />
                            <line x1="21" y1="21" x2="16.65" y2="16.65" />
                        </svg>
                        <input
                            v-model="searchQuery"
                            type="text"
                            placeholder="Rechercher un sujet..."
                            class="search-input"
                        />
                    </div>

                    <div class="select-wrap">
                        <select v-model="filterCategorie" class="filter-select">
                            <option value="">Toutes les catégories</option>
                            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                        </select>
                        <span class="chevron">&#8964;</span>
                    </div>
                </div>

                <div class="forum-table">
                    <div class="forum-header">
                        <span class="col-sujet">Sujet</span>
                        <span class="col-stats">Réponses</span>
                        <span class="col-stats">Vues</span>
                        <span class="col-activite">Activité</span>
                    </div>

                    <router-link
                        v-for="sujet in sujetsFiltres"
                        :key="sujet.id"
                        :to="`/forum/${sujet.id}`"
                        class="forum-row"
                    >
                        <div class="col-sujet">
                            <div class="sujet-avatar">{{ sujet.avatar }}</div>
                            <div class="sujet-info">
                                <h3 class="sujet-titre">{{ sujet.titre }}</h3>
                                <div class="sujet-meta">
                                    <span class="badge-cat">{{ sujet.categorie }}</span>
                                    <span class="sujet-auteur">par {{ sujet.auteur }}</span>
                                    <span class="sujet-date">· {{ sujet.date }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="col-stats">
                            <span class="stat-value">{{ sujet.reponses }}</span>
                        </div>
                        <div class="col-stats">
                            <span class="stat-value">{{ sujet.vues }}</span>
                        </div>
                        <div class="col-activite">
                            <span class="activite-text">{{ sujet.derniereActivite }}</span>
                        </div>
                    </router-link>

                    <div v-if="sujetsFiltres.length === 0" class="empty-state">
                        <p>Aucun sujet ne correspond à votre recherche.</p>
                    </div>
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
    padding: 56px 0 40px;
}
.hero-inner {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 24px;
    flex-wrap: wrap;
}
.hero-title {
    font-size: clamp(2.4rem, 5vw, 3.8rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0 0 10px;
}
.hero-subtitle {
    font-size: 0.95rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0;
    line-height: 1.5;
}
.btn-nouveau {
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
.btn-nouveau:hover {
    background: var(--green-mid);
}

.forum-section {
    flex: 1;
    padding: 0 0 80px;
}

.forum-toolbar {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
    flex-wrap: wrap;
}

.search-bar {
    flex: 1;
    min-width: 200px;
    display: flex;
    align-items: center;
    gap: 10px;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    padding: 0 16px;
    transition: border-color 0.2s;
    color: var(--charcoal);
    opacity: 0.6;
}
.search-bar:focus-within {
    border-color: var(--green-mid);
    opacity: 1;
}
.search-input {
    flex: 1;
    padding: 12px 0;
    font-size: 0.9rem;
    color: var(--charcoal);
    background: transparent;
    border: none;
    outline: none;
    font-family: inherit;
}
.search-input::placeholder {
    color: rgba(53, 53, 53, 0.45);
}

.select-wrap {
    position: relative;
    display: flex;
    align-items: center;
}
.filter-select {
    appearance: none;
    -webkit-appearance: none;
    padding: 12px 40px 12px 16px;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--charcoal);
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    cursor: pointer;
    font-family: inherit;
    outline: none;
    transition: border-color 0.2s;
    min-width: 200px;
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

.forum-table {
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 10px;
    overflow: hidden;
}

.forum-header {
    display: grid;
    grid-template-columns: 1fr 80px 80px 120px;
    gap: 16px;
    padding: 12px 20px;
    background: var(--green-pale);
    font-size: 0.78rem;
    font-weight: 700;
    color: var(--charcoal);
    letter-spacing: 0.04em;
    text-transform: uppercase;
}

.forum-row {
    display: grid;
    grid-template-columns: 1fr 80px 80px 120px;
    gap: 16px;
    padding: 18px 20px;
    text-decoration: none;
    color: inherit;
    border-top: 1px solid rgba(53, 53, 53, 0.08);
    transition: background 0.15s;
    align-items: center;
}
.forum-row:hover {
    background: rgba(215, 236, 225, 0.4);
}

.col-sujet {
    display: flex;
    align-items: flex-start;
    gap: 14px;
    min-width: 0;
}
.col-stats {
    text-align: center;
}
.col-activite {
    text-align: right;
}

.sujet-avatar {
    flex-shrink: 0;
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: var(--green-pale);
    color: var(--green-dark);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.02em;
}

.sujet-info {
    min-width: 0;
}
.sujet-titre {
    font-size: 0.92rem;
    font-weight: 600;
    color: var(--charcoal);
    margin: 0 0 6px;
    line-height: 1.4;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.forum-row:hover .sujet-titre {
    color: var(--green-dark);
}
.sujet-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}
.badge-cat {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 4px;
    background: var(--green-pale);
    color: var(--green-dark);
    font-size: 0.72rem;
    font-weight: 600;
}
.sujet-auteur {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.7;
}
.sujet-date {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.5;
}

.stat-value {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--charcoal);
}

.activite-text {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.6;
}

.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
}

@media (max-width: 700px) {
    .forum-header {
        grid-template-columns: 1fr;
    }
    .forum-header .col-stats,
    .forum-header .col-activite {
        display: none;
    }
    .forum-row {
        grid-template-columns: 1fr;
    }
    .forum-row .col-stats,
    .forum-row .col-activite {
        display: none;
    }
    .sujet-titre {
        white-space: normal;
    }
}
</style>
