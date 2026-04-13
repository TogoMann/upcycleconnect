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
</script>

<template>
    <div class="home-page">
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
    </div>
</template>

<style scoped>
.home-page {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.container {
    max-width: 1060px;
    margin: 0 auto;
    padding: 0 32px;
}
.green {
    color: var(--green-dark);
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
}
</style>
