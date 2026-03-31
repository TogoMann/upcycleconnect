<script setup lang="ts">
const categories = [
    {
        id: 'recycler',
        label: 'RECYCLER',
        description:
            'Donnez une nouvelle utilité à vos matériaux en fin de vie grâce à notre réseau de collecte et de tri optimisé.',
        img: 'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=400&q=80',
        alt: 'Recycler',
    },
    {
        id: 'reparer',
        label: 'RÉPARER',
        description:
            "Prolongez la durée de vie de vos objets en accédant à l'expertise de nos artisans spécialisés dans la restauration.",
        img: 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=400&q=80',
        alt: 'Réparer',
    },
    {
        id: 'transformer',
        label: 'TRANSFORMER',
        description: 'Métamorphosez vos anciens produits en pièces uniques et design.',
        img: 'https://images.unsplash.com/photo-1565193566173-7a0ee3dbe261?w=400&q=80',
        alt: 'Transformer',
    },
    {
        id: 'seformer',
        label: 'SE FORMER',
        description:
            'Apprenez les gestes de demain avec nos ateliers et développez vos compétences pour devenir un acteur du changement.',
        img: 'https://images.unsplash.com/photo-1523301343968-6a6ebf63c672?w=400&q=80',
        alt: 'Se former',
    },
]

const arguments_ = [
    {
        num: '01',
        title: 'Matériaux sourcés localement',
        text: 'Nous récupérons des objets destinés au rebut pour leur offrir une seconde vie.',
    },
    {
        num: '02',
        title: 'Design unique et artisanal',
        text: 'Chaque pièce est fabriquée à la main, garantissant un modèle que vous ne trouverez nulle part ailleurs.',
    },
    {
        num: '03',
        title: 'Impact carbone réduit',
        text: "En réutilisant l'existant, nous économisons les ressources naturelles et limitons les déchets.",
    },
    {
        num: '04',
        title: 'Communauté de créateurs',
        text: "Accédez à un réseau d'artisans passionnés prêts à réaliser vos projets sur mesure.",
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
                    <router-link to="/" class="nav-link active">Accueil</router-link>
                    <router-link to="/prestations" class="nav-link">Prestations</router-link>
                    <router-link to="/evenements" class="nav-link">Évènements</router-link>
                    <router-link to="/annonces" class="nav-link">Annonces</router-link>
                    <router-link to="/reparer" class="nav-link">Réparer</router-link>
                </nav>
                <router-link to="/auth/login" class="btn-nav">
                    S'inscrire / Se connecter
                </router-link>
            </div>
        </header>

        <section class="hero">
            <div class="hero-text-wrap">
                <h1 class="hero-title">L'art de transformer vos objets.</h1>
            </div>
            <div class="hero-image-wrap">
                <img
                    src="https://images.unsplash.com/photo-1504701954957-2010ec3bcec1?w=1400&q=85"
                    alt="Vue aérienne nature verte"
                    class="hero-img"
                />
            </div>
        </section>

        <section class="categories-section">
            <div class="container">
                <p class="tagline">
                    Valorisons nos ressources pour un avenir plus
                    <span class="tagline-accent">durable.</span>
                </p>
                <div class="categories-grid">
                    <div v-for="cat in categories" :key="cat.id" class="category-card">
                        <div class="category-img-wrap">
                            <img :src="cat.img" :alt="cat.alt" class="category-img" />
                        </div>
                        <h3 class="category-label">{{ cat.label }}</h3>
                        <p class="category-desc">{{ cat.description }}</p>
                    </div>
                </div>
            </div>
        </section>

        <section class="arguments-section">
            <div class="container arguments-inner">
                <div class="arguments-content">
                    <h2 class="arguments-title">
                        Transformez l'existant, préservez l'<span class="green">avenir.</span>
                    </h2>
                    <p class="arguments-subtitle">
                        UpCycleConnect accompagne chaque étape de votre démarche éco-responsable, de
                        l'idée à la réalisation, dans chaque région.
                    </p>
                    <ul class="arguments-list">
                        <li v-for="arg in arguments_" :key="arg.num" class="argument-item">
                            <span class="arg-num">{{ arg.num }}</span>
                            <div class="arg-body">
                                <strong class="arg-title">{{ arg.title }} :</strong>
                                <span class="arg-text"> {{ arg.text }}</span>
                            </div>
                        </li>
                    </ul>
                    <div class="arguments-ctas">
                        <router-link to="/prestations" class="btn btn-primary-dark">
                            Nos prestations
                        </router-link>
                        <router-link to="/evenements" class="btn btn-outline-dark">
                            Nos événements
                        </router-link>
                    </div>
                </div>
                <div class="arguments-image-wrap">
                    <img
                        src="https://images.unsplash.com/photo-1582738412745-8660e269a986?w=700&q=80"
                        alt="Étagère skateboards upcyclés"
                        class="arguments-img"
                    />
                </div>
            </div>
        </section>

        <div class="divider-wrap">
            <div class="container">
                <hr class="divider" />
            </div>
        </div>

        <section class="why-section">
            <div class="container why-inner">
                <h2 class="why-title">
                    Pourquoi rejoindre
                    <span class="green">UpcycleConnect</span> ?
                </h2>
                <p class="why-subtitle">
                    Parce que chaque objet mérite une seconde chance. Nous avons créé un écosystème
                    où la créativité rencontre la durabilité pour transformer notre façon de
                    consommer et d'agir ensemble.
                </p>
                <router-link to="/auth/login" class="btn btn-green-filled">
                    Découvrir la communauté
                </router-link>
            </div>
        </section>

        <section class="testimonial-section">
            <div class="container testimonial-inner">
                <div class="testimonial-img-wrap">
                    <img
                        src="https://images.unsplash.com/photo-1559827291-72ee739d0d9a?w=700&q=80"
                        alt="Installation lumineuse upcyclée"
                        class="testimonial-img"
                    />
                </div>
                <blockquote class="testimonial-quote">
                    <p>
                        "J'étais à la recherche d'une manière unique de décorer mon café sans
                        acheter du neuf. Grâce à UpcycleConnect, j'ai trouvé l'inspiration et les
                        artisans pour transformer de simples tasses en une installation lumineuse
                        qui fascine mes clients tous les jours !!"
                    </p>
                    <footer class="testimonial-author">Mario L.</footer>
                </blockquote>
            </div>
        </section>

        <footer class="footer">
            <div class="footer-container">
                <nav class="footer-links">
                    <a v-for="link in footerLinks" :key="link" href="#" class="footer-link">{{
                        link
                    }}</a>
                </nav>
                <div class="footer-bottom">
                    <span class="footer-logo">UpCycleConnect</span>
                    <div class="footer-socials">
                        <a href="#" class="social-icon" aria-label="Instagram">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                                <path
                                    d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"
                                />
                            </svg>
                        </a>
                        <a href="#" class="social-icon" aria-label="Facebook">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                                <path
                                    d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"
                                />
                            </svg>
                        </a>
                        <a href="#" class="social-icon" aria-label="Twitter">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                                <path
                                    d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"
                                />
                            </svg>
                        </a>
                    </div>
                    <div class="footer-lang">
                        <span>Français (FR)</span>
                        <svg
                            width="12"
                            height="12"
                            viewBox="0 0 12 12"
                            fill="currentColor"
                            style="margin-left: 4px"
                        >
                            <path d="M6 8L1 3h10z" />
                        </svg>
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
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}
.green {
    color: var(--green-dark);
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
    flex-shrink: 0;
    text-decoration: none;
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
    background: var(--cream);
    padding: 56px 32px 0;
}
.hero-text-wrap {
    max-width: 1060px;
    margin: 0 auto;
    text-align: center;
    padding-bottom: 32px;
}
.hero-title {
    font-size: clamp(2.2rem, 5vw, 3.5rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.1;
    letter-spacing: -0.02em;
    margin: 0;
}
.hero-image-wrap {
    max-width: 1060px;
    margin: 0 auto;
    border-radius: 12px;
    overflow: hidden;
    height: 340px;
}
.hero-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.categories-section {
    padding: 64px 0 80px;
}
.tagline {
    text-align: center;
    font-size: 1.2rem;
    font-weight: 400;
    color: var(--charcoal);
    margin: 0 0 52px;
    line-height: 1.5;
}
.tagline-accent {
    color: var(--green-dark);
    font-weight: 600;
}
.categories-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
}
.category-card {
    display: flex;
    flex-direction: column;
}
.category-img-wrap {
    width: 100%;
    aspect-ratio: 4/3;
    border-radius: 8px;
    overflow: hidden;
    margin-bottom: 12px;
}
.category-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    transition: transform 0.4s ease;
}
.category-img-wrap:hover .category-img {
    transform: scale(1.06);
}
.category-label {
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    color: var(--charcoal);
    margin: 0 0 7px;
    text-transform: uppercase;
}
.category-desc {
    font-size: 0.8rem;
    color: var(--charcoal);
    line-height: 1.6;
    margin: 0;
    opacity: 0.8;
}

.arguments-section {
    padding: 72px 0;
}
.arguments-inner {
    display: flex;
    align-items: flex-start;
    gap: 60px;
}
.arguments-content {
    flex: 1;
    min-width: 0;
}
.arguments-title {
    font-size: clamp(1.5rem, 2.8vw, 2rem);
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.2;
    margin: 0 0 14px;
}
.arguments-subtitle {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.72;
    line-height: 1.65;
    margin: 0 0 28px;
}
.arguments-list {
    list-style: none;
    padding: 0;
    margin: 0 0 32px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.argument-item {
    display: flex;
    align-items: flex-start;
    gap: 14px;
    font-size: 0.86rem;
    line-height: 1.6;
}
.arg-num {
    font-weight: 700;
    color: var(--green-mid);
    font-size: 0.82rem;
    flex-shrink: 0;
    min-width: 20px;
    margin-top: 1px;
}
.arg-title {
    font-weight: 700;
    color: var(--charcoal);
}
.arg-text {
    color: var(--charcoal);
    opacity: 0.8;
}
.arguments-ctas {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}
.btn {
    display: inline-block;
    padding: 10px 20px;
    border-radius: 6px;
    font-size: 0.84rem;
    font-weight: 600;
    text-decoration: none;
    cursor: pointer;
    transition: all 0.2s ease;
    border: 2px solid transparent;
}
.btn-primary-dark {
    background: var(--green-mid);
    color: var(--white);
    border-color: var(--green-mid);
}
.btn-primary-dark:hover {
    background: var(--green-dark);
    border-color: var(--green-dark);
}
.btn-outline-dark {
    background: transparent;
    color: var(--green-mid);
    border-color: var(--green-mid);
}
.btn-outline-dark:hover {
    background: var(--green-pale);
}
.btn-green-filled {
    background: var(--green-mid);
    color: var(--white);
    border-color: var(--green-mid);
    padding: 13px 36px;
    font-size: 0.9rem;
}
.btn-green-filled:hover {
    background: var(--green-dark);
    border-color: var(--green-dark);
    transform: translateY(-1px);
    box-shadow: 0 6px 18px rgba(52, 137, 91, 0.3);
}
.arguments-image-wrap {
    flex: 0 0 380px;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
}
.arguments-img {
    width: 100%;
    height: 460px;
    object-fit: cover;
    display: block;
}

.divider-wrap {
    padding: 0 0 12px;
}
.divider {
    border: none;
    border-top: 1px solid rgba(53, 53, 53, 0.15);
    margin: 0;
}

.why-section {
    padding: 80px 0 72px;
}
.why-inner {
    text-align: center;
}
.why-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.2;
    margin: 0 0 18px;
    letter-spacing: -0.02em;
}
.why-subtitle {
    font-size: 0.93rem;
    color: var(--charcoal);
    opacity: 0.72;
    line-height: 1.7;
    max-width: 540px;
    margin: 0 auto 32px;
}

.testimonial-section {
    padding: 0 0 80px;
}
.testimonial-inner {
    display: flex;
    align-items: center;
    gap: 60px;
}
.testimonial-img-wrap {
    flex: 0 0 295px;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
}
.testimonial-img {
    width: 100%;
    height: 350px;
    object-fit: cover;
    display: block;
}
.testimonial-quote {
    flex: 1;
    margin: 0;
    padding: 0;
    border: none;
}
.testimonial-quote p {
    font-size: clamp(1.1rem, 2.2vw, 1.5rem);
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.45;
    margin: 0 0 18px;
}
.testimonial-author {
    font-size: 0.88rem;
    color: var(--charcoal);
    opacity: 0.6;
    font-style: normal;
}

.footer {
    background: var(--green-dark);
    color: var(--white);
    padding: 40px 0;
}
.footer-container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 24px;
}
.footer-links {
    display: flex;
    gap: 32px;
}
.footer-link {
    color: rgba(255, 255, 255, 0.72);
    text-decoration: none;
    font-size: 0.84rem;
    transition: color 0.2s;
}
.footer-link:hover {
    color: var(--white);
}
.footer-bottom {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.footer-logo {
    font-weight: 700;
    font-size: 1rem;
    color: var(--white);
}
.footer-socials {
    display: flex;
    gap: 16px;
}
.social-icon {
    color: rgba(255, 255, 255, 0.72);
    transition: color 0.2s;
    display: flex;
    align-items: center;
}
.social-icon:hover {
    color: var(--white);
}
.footer-lang {
    display: flex;
    align-items: center;
    font-size: 0.84rem;
    color: rgba(255, 255, 255, 0.85);
    cursor: pointer;
}

@media (max-width: 900px) {
    .categories-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .arguments-inner,
    .testimonial-inner {
        flex-direction: column;
    }
    .arguments-image-wrap {
        flex: none;
        width: 100%;
    }
    .arguments-img {
        height: 260px;
    }
    .testimonial-img-wrap {
        flex: none;
        width: 100%;
    }
    .testimonial-img {
        height: 230px;
    }
    .nav-links {
        display: none;
    }
    .hero-image-wrap {
        height: 200px;
    }
}
@media (max-width: 560px) {
    .categories-grid {
        grid-template-columns: 1fr;
    }
    .hero-title {
        font-size: 1.9rem;
    }
    .footer-bottom {
        flex-direction: column;
        gap: 16px;
        text-align: center;
    }
    .footer-links {
        flex-direction: column;
        align-items: center;
        gap: 12px;
    }
}
</style>
