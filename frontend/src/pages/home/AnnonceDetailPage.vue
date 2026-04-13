<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const annonces = [
    {
        id: 1,
        titre: 'Chaise en chêne',
        description: 'Belle chaise ancienne en chêne massif. Légèrement usée mais très solide, idéale pour un projet de ponçage et de vernissage. Pieds stables, assise en bon état. À récupérer sur place.',
        type: 'Don',
        prix: null,
        localisation: 'Nantes',
        quartier: 'Centre-ville',
        date: '14 Mai 2026',
        vendeur: 'Mégane K.',
        vendeurdepuis: 'Membre depuis 2024',
        categorie: 'Mobilier',
        etat: 'Bon état',
        img: 'https://images.unsplash.com/photo-1592078615290-033ee584e267?w=800&q=80',
        imgs: [
            'https://images.unsplash.com/photo-1592078615290-033ee584e267?w=800&q=80',
            'https://images.unsplash.com/photo-1598300042247-d088f8ab3a91?w=800&q=80',
        ],
    },
    {
        id: 2,
        titre: 'Vase',
        description: 'Vase en céramique artisanale, couleur terracotta. Quelques micro-éclats à la base, pas visibles une fois posé. Hauteur 28 cm. Idéal pour une décoration vintage ou bohème.',
        type: 'Vente',
        prix: 25,
        localisation: 'Nantes',
        quartier: 'Île de Nantes',
        date: '14 Mai 2026',
        vendeur: 'Mégane K.',
        vendeurdepuis: 'Membre depuis 2024',
        categorie: 'Décoration',
        etat: 'Très bon état',
        img: 'https://images.unsplash.com/photo-1612196808214-b8e1d6145a8c?w=800&q=80',
        imgs: [
            'https://images.unsplash.com/photo-1612196808214-b8e1d6145a8c?w=800&q=80',
        ],
    },
]

const annonce = computed(() => {
    const id = Number(route.params.id)
    return annonces.find((a) => a.id === id) ?? annonces[0]
})
</script>

<template>
    <div class="page-content">
        <section class="breadcrumb-bar">
            <div class="container">
                <router-link to="/annonces" class="breadcrumb-link">Annonces</router-link>
                <span class="breadcrumb-sep">›</span>
                <span class="breadcrumb-current">{{ annonce.titre }}</span>
            </div>
        </section>

        <section class="detail-section">
            <div class="container">
                <div class="detail-layout">
                    <div class="detail-left">
                        <div class="main-img-wrap">
                            <img :src="annonce.img" :alt="annonce.titre" class="main-img" />
                        </div>
                        <div v-if="annonce.imgs.length > 1" class="thumbs-row">
                            <div
                                v-for="(img, i) in annonce.imgs"
                                :key="i"
                                class="thumb-wrap"
                            >
                                <img :src="img" :alt="`${annonce.titre} ${i + 1}`" class="thumb-img" />
                            </div>
                        </div>
                    </div>

                    <div class="detail-right">
                        <div class="detail-badges">
                            <span class="badge" :class="annonce.type === 'Don' ? 'badge--don' : 'badge--vente'">
                                {{ annonce.type }}
                            </span>
                            <span class="badge badge--cat">{{ annonce.categorie }}</span>
                        </div>

                        <h1 class="detail-titre">{{ annonce.titre }}</h1>

                        <div v-if="annonce.type === 'Vente' && annonce.prix" class="detail-prix">
                            {{ annonce.prix }} €
                        </div>
                        <div v-else class="detail-gratuit">Gratuit</div>

                        <div class="detail-meta">
                            <div class="meta-item">
                                <svg viewBox="0 0 24 24" fill="currentColor" width="14" height="14">
                                    <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
                                </svg>
                                <span>{{ annonce.localisation }}, {{ annonce.quartier }}</span>
                            </div>
                            <div class="meta-item">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                                    <line x1="16" y1="2" x2="16" y2="6"/>
                                    <line x1="8" y1="2" x2="8" y2="6"/>
                                    <line x1="3" y1="10" x2="21" y2="10"/>
                                </svg>
                                <span>Publié le {{ annonce.date }}</span>
                            </div>
                            <div class="meta-item">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                                </svg>
                                <span>État : {{ annonce.etat }}</span>
                            </div>
                        </div>

                        <div class="detail-desc">
                            <h2 class="desc-title">Description</h2>
                            <p class="desc-text">{{ annonce.description }}</p>
                        </div>

                        <div class="vendeur-card">
                            <div class="vendeur-avatar">{{ annonce.vendeur.slice(0, 2).toUpperCase() }}</div>
                            <div class="vendeur-info">
                                <span class="vendeur-nom">{{ annonce.vendeur }}</span>
                                <span class="vendeur-since">{{ annonce.vendeurdepuis }}</span>
                            </div>
                        </div>

                        <div class="detail-actions">
                            <router-link to="/auth/login" class="btn-contact">
                                Contacter le vendeur
                            </router-link>
                            <button class="btn-save">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18">
                                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                                </svg>
                                Sauvegarder
                            </button>
                        </div>
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

.breadcrumb-bar {
    padding: 20px 0;
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
}
.breadcrumb-bar .container {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
}
.breadcrumb-link {
    color: var(--green-mid);
    text-decoration: none;
    transition: color 0.2s;
}
.breadcrumb-link:hover {
    color: var(--green-dark);
}
.breadcrumb-sep {
    color: var(--charcoal);
    opacity: 0.4;
}
.breadcrumb-current {
    color: var(--charcoal);
    opacity: 0.7;
}

.detail-section {
    flex: 1;
    padding: 40px 0 80px;
}

.detail-layout {
    display: grid;
    grid-template-columns: 1fr 400px;
    gap: 48px;
    align-items: flex-start;
}

.main-img-wrap {
    border-radius: 12px;
    overflow: hidden;
    aspect-ratio: 4/3;
    margin-bottom: 12px;
}
.main-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.thumbs-row {
    display: flex;
    gap: 10px;
}
.thumb-wrap {
    width: 80px;
    height: 60px;
    border-radius: 6px;
    overflow: hidden;
    border: 2px solid transparent;
    cursor: pointer;
    transition: border-color 0.2s;
}
.thumb-wrap:hover {
    border-color: var(--green-mid);
}
.thumb-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.detail-badges {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
}
.badge {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 5px;
    font-size: 0.75rem;
    font-weight: 700;
}
.badge--don {
    background: var(--green-mid);
    color: var(--white);
}
.badge--vente {
    background: var(--green-mid);
    color: var(--white);
}
.badge--cat {
    background: var(--green-pale);
    color: var(--green-dark);
}

.detail-titre {
    font-size: clamp(1.6rem, 3vw, 2.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.2;
    letter-spacing: -0.02em;
    margin: 0 0 16px;
}

.detail-prix {
    font-size: 2rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
    margin-bottom: 20px;
}
.detail-gratuit {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--green-mid);
    margin-bottom: 20px;
}

.detail-meta {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 28px;
}
.meta-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
    color: var(--charcoal);
    opacity: 0.75;
}
.meta-item svg {
    color: var(--green-mid);
    flex-shrink: 0;
}

.detail-desc {
    border-top: 1px solid rgba(53, 53, 53, 0.1);
    padding-top: 20px;
    margin-bottom: 24px;
}
.desc-title {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--charcoal);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin: 0 0 12px;
}
.desc-text {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.82;
    line-height: 1.7;
    margin: 0;
}

.vendeur-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: var(--green-pale);
    border-radius: 10px;
    margin-bottom: 20px;
}
.vendeur-avatar {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: var(--green-dark);
    color: var(--white);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 700;
    flex-shrink: 0;
}
.vendeur-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}
.vendeur-nom {
    font-size: 0.92rem;
    font-weight: 700;
    color: var(--charcoal);
}
.vendeur-since {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.6;
}

.detail-actions {
    display: flex;
    gap: 12px;
}
.btn-contact {
    flex: 1;
    display: block;
    padding: 14px 24px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 700;
    text-decoration: none;
    text-align: center;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-contact:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}
.btn-save {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 18px;
    background: transparent;
    color: var(--charcoal);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        border-color 0.2s,
        color 0.2s;
    white-space: nowrap;
}
.btn-save:hover {
    border-color: var(--green-mid);
    color: var(--green-mid);
}

@media (max-width: 860px) {
    .detail-layout {
        grid-template-columns: 1fr;
    }
}
@media (max-width: 560px) {
    .detail-actions {
        flex-direction: column;
    }
    .btn-save {
        justify-content: center;
    }
}
</style>
