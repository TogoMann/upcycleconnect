<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()

const searchQuery = ref('')
const filterType = ref('')
const filterCategorie = ref('')
const filterLocalisation = ref('')

onMounted(async () => {
    await clientStore.fetchAllAnnonces()
})

const annonces = computed(() => {
    return clientStore.allAnnonces.map((a: any) => ({
        id: a.id.Int64,
        titre: a.name,
        description: a.description,
        type: a.price?.Exp === 0 && a.price?.Int === 0 ? 'Don' : 'Vente',
        prix: a.price?.Exp === 0 && a.price?.Int === 0 ? null : parseFloat(a.price?.Int) / 100, // Very simple parsing for Numeric
        localisation: a.city_name || 'Inconnue',
        date: new Date(a.created_at).toLocaleDateString('fr-FR', {
            day: 'numeric',
            month: 'short',
        }),
        vendeur: 'Utilisateur', // We don't have the seller name in the listing model yet
        categorie: a.category || 'Non classé',
        img: 'https://images.unsplash.com/photo-1592078615290-033ee584e267?w=400&q=80', // Placeholder
    }))
})

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
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">A vendre ou à donner.</h1>
            </div>
        </section>

        <section class="search-section">
            <div class="container">
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

        <section class="annonces-section">
            <div class="container">
                <div class="annonces-grid">
                    <router-link
                        v-for="annonce in annoncesFiltrees"
                        :key="annonce.id"
                        :to="`/annonces/${annonce.id}`"
                        class="annonce-card"
                    >
                        <div class="annonce-img-wrap">
                            <img :src="annonce.img" :alt="annonce.titre" class="annonce-img" />
                        </div>

                        <div class="annonce-body">
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

                            <span class="btn-annonce">Voir l'annonce</span>
                        </div>
                    </router-link>
                </div>

                <div v-if="annoncesFiltrees.length === 0" class="empty-state">
                    <p>Aucune annonce ne correspond à votre recherche.</p>
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

.search-section {
    padding: 0 0 32px;
}

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

.annonces-section {
    flex: 1;
    padding: 0 0 80px;
}
.annonces-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
}

.annonce-card {
    background: var(--green-pale);
    border-radius: 10px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    text-decoration: none;
    color: inherit;
    transition:
        transform 0.2s,
        box-shadow 0.2s;
}
.annonce-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(8, 106, 53, 0.12);
}

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

.annonce-body {
    padding: 12px 14px 16px;
    display: flex;
    flex-direction: column;
    flex: 1;
}

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

.annonce-titre {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 4px;
    line-height: 1.3;
}

.annonce-desc {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0 0 8px;
    line-height: 1.4;
}

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

.annonce-vendeur {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.7;
    margin: 0 0 12px;
}

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
    text-align: center;
    display: block;
}
.annonce-card:hover .btn-annonce {
    background: var(--green-mid);
}

.empty-state {
    text-align: center;
    padding: 60px 0;
    color: var(--charcoal);
    opacity: 0.6;
    font-size: 1rem;
}

@media (max-width: 900px) {
    .annonces-grid {
        grid-template-columns: repeat(2, 1fr);
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
}
</style>
