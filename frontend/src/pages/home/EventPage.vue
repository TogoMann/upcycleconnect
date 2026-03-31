<script setup lang="ts">
const events = [
    {
        id: 'conference',
        title: "Conférence : L'avenir circulaire",
        description:
            "Explorez les enjeux et les opportunités de l'économie circulaire avec des experts passionnés. Une immersion au cœur des solutions durables pour transformer nos modes de consommation.",
        date: '13 septembre 2026 – Paris',
        img: 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=1200&q=85',
        alt: 'Conférence économie circulaire',
    },
    {
        id: 'popup',
        title: 'Pop up Store : nos créateurs',
        description:
            "Découvrez une sélection exclusive d'objets uniques et recyclés. Rencontrez les artisans locaux qui redonnent vie aux matériaux et repartez avec une pièce chargée d'histoire.",
        date: '30 juillet 2026 – Paris',
        img: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=1200&q=85',
        alt: 'Pop up store créateurs',
    },
    {
        id: 'atelier',
        title: "Atelier d'upcycling créatif",
        description:
            "Passez de la théorie à la pratique ! Apprenez à transformer vos propres objets lors d'une session conviviale encadrée par un artisan designer. Repartez avec votre création.",
        date: '12 aout 2026 – Paris',
        img: 'https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=1200&q=85',
        alt: "Atelier d'upcycling créatif",
    },
]

const footerLinks = ['À propos', 'Mentions légales', 'Politique de confidentialité']
</script>

<template>
    <div class="page">
        <header class="navbar">
            <div class="nav-container">
                <router-link to="/" class="nav-logo">UpCycleConnect</router-link>
                <nav class="nav-links">
                    <router-link to="/" class="nav-link">Accueil</router-link>
                    <router-link to="/prestations" class="nav-link">Prestations</router-link>
                    <router-link to="/evenements" class="nav-link active">Évènements</router-link>
                    <router-link to="/forum" class="nav-link">Forum</router-link>
                    <router-link to="/a-propos" class="nav-link">À propos</router-link>
                </nav>
                <router-link to="/auth/register" class="btn-nav">
                    S'inscrire / Se connecter
                </router-link>
            </div>
        </header>

        <section class="hero">
            <div class="container">
                <h1 class="hero-title">Découvrez nos évènements.</h1>
            </div>
        </section>

        <section class="events-section">
            <div class="container">
                <div
                    v-for="(event, index) in events"
                    :key="event.id"
                    class="event-block"
                    :class="{ 'event-block--last': index === events.length - 1 }"
                >
                    <div class="event-img-wrap">
                        <img :src="event.img" :alt="event.alt" class="event-img" />
                    </div>

                    <div class="event-info">
                        <h2 class="event-title">{{ event.title }}</h2>
                        <p class="event-desc">{{ event.description }}</p>
                    </div>

                    <div class="event-footer">
                        <span class="event-date">{{ event.date }}</span>
                        <button class="btn-reserver">Réserver</button>
                    </div>
                </div>
            </div>
        </section>

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
                        <span class="lang-sep"> &nbsp;·&nbsp; </span>
                        <span>Français</span>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<style scoped>
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

.hero {
    padding: 56px 0 36px;
    background: var(--cream);
}
.hero-title {
    font-size: clamp(2.4rem, 5.5vw, 4rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.08;
    letter-spacing: -0.03em;
    margin: 0;
}

.events-section {
    flex: 1;
    padding: 0 0 80px;
}

.event-block {
    padding-bottom: 56px;
    margin-bottom: 56px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.12);
}
.event-block--last {
    border-bottom: none;
    margin-bottom: 0;
}

.event-img-wrap {
    width: 100%;
    height: 380px;
    border-radius: 10px;
    overflow: hidden;
    margin-bottom: 24px;
}
.event-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.5s ease;
}
.event-img-wrap:hover .event-img {
    transform: scale(1.03);
}

.event-title {
    font-size: clamp(1.4rem, 2.8vw, 1.9rem);
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.2;
    margin: 0 0 12px;
    letter-spacing: -0.01em;
}

.event-info {
    max-width: 680px;
}
.event-desc {
    font-size: 0.9rem;
    color: var(--charcoal);
    line-height: 1.65;
    margin: 0;
    opacity: 0.8;
}

.event-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 16px;
    margin-top: 24px;
}

.event-date {
    background: var(--green-pale);
    color: var(--green-dark);
    font-size: 0.82rem;
    font-weight: 600;
    padding: 8px 16px;
    border-radius: 6px;
    white-space: nowrap;
    letter-spacing: 0.01em;
}

.btn-reserver {
    background: var(--green-mid);
    color: var(--white);
    border: none;
    padding: 10px 24px;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition:
        background 0.2s,
        transform 0.15s;
    font-family: inherit;
    white-space: nowrap;
}
.btn-reserver:hover {
    background: var(--green-dark);
    transform: translateY(-1px);
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

@media (max-width: 860px) {
    .event-img-wrap {
        height: 260px;
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
    .hero-title {
        font-size: 2.2rem;
    }
    .event-img-wrap {
        height: 200px;
    }
    .event-footer {
        flex-direction: column;
        align-items: flex-start;
        gap: 12px;
    }
    .footer-container {
        flex-direction: column;
        gap: 12px;
        text-align: center;
    }
}
</style>
