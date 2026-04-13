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
</script>

<template>
    <div class="page-content">
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

@media (max-width: 860px) {
    .event-img-wrap {
        height: 260px;
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
}
</style>
