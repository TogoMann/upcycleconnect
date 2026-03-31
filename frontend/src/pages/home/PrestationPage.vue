<script setup lang="ts">
const prestations = [
    {
        id: 'recycler',
        label: 'RECYCLER',
        description:
            'Cette prestation inclut un audit de vos matériaux, la mise à disposition de bacs de collecte spécifiques et le transport vers nos centres de tri partenaires où chaque matière est traitée pour réintégrer un cycle de production.',
        img: 'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=600&q=80',
        alt: 'Recycler — bacs de tri colorés',
    },
    {
        id: 'reparer',
        label: 'RÉPARER',
        description:
            "Un service de remise en état sur mesure. Nous diagnostiquons l'origine de la panne ou de l'usure de votre objet et nous le confions à un expert qualifié. La prestation comprend la recherche de pièces détachées de seconde main et une garantie sur la réparation effectuée.",
        img: 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=600&q=80',
        alt: 'Réparer — électronique',
    },
    {
        id: 'transformer',
        label: 'TRANSFORMER',
        description:
            'Un accompagnement créatif pour vos projets de design. De la conception du croquis à la fabrication finale, nos artisans utilisent des techniques de surcyclage pour créer du mobilier ou des accessoires uniques à partir de vos anciens objets.',
        img: 'https://images.unsplash.com/photo-1565193566173-7a0ee3dbe261?w=600&q=80',
        alt: 'Transformer — travail du bois',
    },
    {
        id: 'seformer',
        label: 'SE FORMER',
        description:
            "Des modules d'apprentissage concrets de 2 à 4 heures. Chaque session vous permet de manipuler les outils, d'apprendre des techniques spécifiques et de repartir avec votre propre création ou un guide méthodologique complet.",
        img: 'https://images.unsplash.com/photo-1523301343968-6a6ebf63c672?w=600&q=80',
        alt: 'Se former — nature et écologie',
    },
]

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
                    <router-link to="/prestations" class="nav-link active">Prestations</router-link>
                    <router-link to="/evenements" class="nav-link">Évènements</router-link>
                    <router-link to="/annonces" class="nav-link">Annonces</router-link>
                    <router-link to="/reparer" class="nav-link">Réparer</router-link>
                </nav>
                <router-link to="/auth/login" class="btn-nav">
                    S'inscrire / Se connecter
                </router-link>
            </div>
        </header>

        <!-- ════════ HERO ════════ -->
        <section class="hero">
            <div class="container hero-inner">
                <h1 class="hero-title">Découvrez nos prestations.</h1>
                <p class="hero-subtitle">
                    Valorisons nos ressources pour un avenir plus
                    <span class="green-accent">durable</span>.
                </p>
            </div>
        </section>

        <!-- ════════ PRESTATIONS GRID ════════ -->
        <section class="prestations-section">
            <div class="container">
                <div class="prestations-grid">
                    <div
                        v-for="(p, index) in prestations"
                        :key="p.id"
                        class="prestation-row"
                        :class="{ 'row-border': index < prestations.length - 2 }"
                    >
                        <!-- Gauche : image -->
                        <div class="prestation-img-wrap">
                            <img :src="p.img" :alt="p.alt" class="prestation-img" />
                        </div>
                        <!-- Droite : contenu -->
                        <div class="prestation-content">
                            <h2 class="prestation-label">{{ p.label }}</h2>
                            <p class="prestation-desc">{{ p.description }}</p>
                            <button class="btn-reserver">Réserver</button>
                        </div>
                    </div>
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
    color: var(--green-mid);
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
    background: var(--cream);
    padding: 72px 0 52px;
    text-align: center;
}
.hero-inner {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.hero-title {
    font-size: clamp(2.4rem, 5vw, 3.8rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.1;
    letter-spacing: -0.025em;
    margin: 0 0 18px;
}
.hero-subtitle {
    font-size: clamp(1.1rem, 2.2vw, 1.5rem);
    font-weight: 400;
    color: var(--charcoal);
    margin: 0;
    line-height: 1.4;
}
.green-accent {
    color: var(--green-dark);
    font-weight: 500;
}

/* ══ PRESTATIONS ══ */
.prestations-section {
    padding: 8px 0 80px;
    flex: 1;
}

/* Grille 2 colonnes de rangées */
.prestations-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0;
}

/* Chaque rangée = image + texte côte à côte */
.prestation-row {
    display: flex;
    align-items: flex-start;
    gap: 0;
    padding: 36px 32px 36px 0;
    border-bottom: 1px solid rgba(53, 53, 53, 0.1);
}

/* rangées du bas (3ème & 4ème) : pas de bordure en bas */
.prestation-row:nth-child(3),
.prestation-row:nth-child(4) {
    border-bottom: none;
}

/* Séparateur vertical entre les 2 colonnes */
.prestation-row:nth-child(odd) {
    padding-right: 48px;
    border-right: 1px solid rgba(53, 53, 53, 0.1);
}
.prestation-row:nth-child(even) {
    padding-left: 48px;
    padding-right: 0;
}

.prestation-img-wrap {
    flex: 0 0 240px;
    height: 170px;
    border-radius: 8px;
    overflow: hidden;
    margin-right: 24px;
    flex-shrink: 0;
}
.prestation-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.4s ease;
}
.prestation-img-wrap:hover .prestation-img {
    transform: scale(1.06);
}

.prestation-content {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding-top: 4px;
}

.prestation-label {
    font-size: 0.88rem;
    font-weight: 800;
    letter-spacing: 0.06em;
    color: var(--charcoal);
    text-transform: uppercase;
    margin: 0 0 12px;
}

.prestation-desc {
    font-size: 0.83rem;
    color: var(--charcoal);
    line-height: 1.65;
    margin: 0 0 20px;
    opacity: 0.85;
    flex: 1;
}

.btn-reserver {
    background: var(--green-mid);
    color: var(--white);
    border: none;
    padding: 10px 22px;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    transition:
        background 0.2s,
        transform 0.15s;
    font-family: inherit;
}
.btn-reserver:hover {
    background: var(--green-dark);
    transform: translateY(-1px);
}

/* ══ FOOTER ══ */
.footer {
    background: var(--green-dark);
    color: var(--white);
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
@media (max-width: 860px) {
    .prestations-grid {
        grid-template-columns: 1fr;
    }
    .prestation-row {
        border-right: none !important;
        padding-left: 0 !important;
        padding-right: 0 !important;
        border-bottom: 1px solid rgba(53, 53, 53, 0.1) !important;
    }
    .prestation-row:last-child {
        border-bottom: none !important;
    }
    .prestation-img-wrap {
        flex: 0 0 160px;
        height: 130px;
    }
    .nav-links {
        display: none;
    }
    .footer-links-wrap {
        flex-direction: column;
        align-items: center;
        gap: 12px;
    }
}

@media (max-width: 560px) {
    .prestation-row {
        flex-direction: column;
    }
    .prestation-img-wrap {
        flex: none;
        width: 100%;
        height: 200px;
        margin-right: 0;
        margin-bottom: 16px;
    }
    .hero-title {
        font-size: 2rem;
    }
    .footer-container {
        flex-direction: column;
        gap: 12px;
        text-align: center;
    }
}
</style>
