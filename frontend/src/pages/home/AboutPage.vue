<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { API_BASE } from '@/config'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const valeurs = computed(() => [
    { num: '01', titre: t('about.values.item1.title'), texte: t('about.values.item1.text') },
    { num: '02', titre: t('about.values.item2.title'), texte: t('about.values.item2.text') },
    { num: '03', titre: t('about.values.item3.title'), texte: t('about.values.item3.text') },
    { num: '04', titre: t('about.values.item4.title'), texte: t('about.values.item4.text') },
])

const stats = ref({
    active_members: 0,
    items_renewed: 0,
    partner_artisans: 0,
    regions_covered: 0,
})

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/stats/public`)
        if (res.ok) stats.value = await res.json()
    } catch {}
})

const chiffres = computed(() => [
    { valeur: `${stats.value.active_members}`, label: t('about.stats.activeMembers') },
    { valeur: `${stats.value.items_renewed}`, label: t('about.stats.itemsRenewed') },
    { valeur: `${stats.value.partner_artisans}`, label: t('about.stats.partnerArtisans') },
    { valeur: `${stats.value.regions_covered}`, label: t('about.stats.regionsCovered') },
])
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">{{ t('about.hero.title') }}</h1>
                <p class="hero-subtitle">
                    {{ t('about.hero.subtitle') }}
                </p>
            </div>
        </section>

        <section class="image-section">
            <div class="container">
                <div class="hero-img-wrap">
                    <img
                        src="https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=1400&q=85"
                        :alt="t('about.hero.imageAlt')"
                        class="hero-img"
                    />
                </div>
            </div>
        </section>

        <section class="chiffres-section">
            <div class="container">
                <div class="chiffres-grid">
                    <div v-for="c in chiffres" :key="c.label" class="chiffre-card">
                        <span class="chiffre-valeur">{{ c.valeur }}</span>
                        <span class="chiffre-label">{{ c.label }}</span>
                    </div>
                </div>
            </div>
        </section>

        <section class="valeurs-section">
            <div class="container valeurs-inner">
                <div class="valeurs-header">
                    <h2 class="valeurs-title">
                        {{ t('about.values.title') }} <span class="green">{{ t('about.values.titleAccent') }}</span>
                    </h2>
                    <p class="valeurs-subtitle">
                        {{ t('about.values.subtitle') }}
                    </p>
                </div>

                <ul class="valeurs-list">
                    <li v-for="v in valeurs" :key="v.num" class="valeur-item">
                        <span class="valeur-num">{{ v.num }}</span>
                        <div class="valeur-body">
                            <strong class="valeur-titre">{{ v.titre }} :</strong>
                            <span class="valeur-texte"> {{ v.texte }}</span>
                        </div>
                    </li>
                </ul>
            </div>
        </section>

        <section class="story-section">
            <div class="container story-inner">
                <div class="story-img-wrap">
                    <img
                        src="https://images.unsplash.com/photo-1552664730-d307ca884978?w=700&q=80"
                        :alt="t('about.story.imageAlt')"
                        class="story-img"
                    />
                </div>
                <div class="story-content">
                    <h2 class="story-title">{{ t('about.story.title') }}</h2>
                    <p class="story-text">
                        {{ t('about.story.text1') }}
                    </p>
                    <p class="story-text">
                        {{ t('about.story.text2') }}
                    </p>
                    <router-link v-if="!authStore.isAuthenticated" to="/auth/register" class="btn-join">
                        {{ t('about.story.cta') }}
                    </router-link>
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

.green {
    color: var(--green-dark);
}

.hero {
    padding: 72px 0 48px;
    text-align: center;
}
.hero-title {
    font-size: clamp(2.8rem, 6vw, 4.8rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.05;
    letter-spacing: -0.035em;
    margin: 0 0 24px;
}
.hero-subtitle {
    font-size: clamp(1rem, 2vw, 1.25rem);
    color: var(--charcoal);
    opacity: 0.75;
    line-height: 1.7;
    max-width: 640px;
    margin: 0 auto;
}

.image-section {
    padding: 0 0 72px;
}
.hero-img-wrap {
    border-radius: 12px;
    overflow: hidden;
    height: 420px;
}
.hero-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.chiffres-section {
    padding: 56px 0;
    background: var(--green-pale);
}
.chiffres-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 32px;
    text-align: center;
}
.chiffre-card {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.chiffre-valeur {
    font-size: clamp(2rem, 4vw, 3rem);
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
}
.chiffre-label {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.75;
    font-weight: 500;
}

.valeurs-section {
    padding: 80px 0;
}
.valeurs-inner {
    display: flex;
    gap: 80px;
    align-items: flex-start;
}
.valeurs-header {
    flex: 0 0 320px;
}
.valeurs-title {
    font-size: clamp(1.6rem, 3vw, 2.2rem);
    font-weight: 700;
    color: var(--charcoal);
    line-height: 1.2;
    margin: 0 0 16px;
}
.valeurs-subtitle {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.7;
    line-height: 1.65;
    margin: 0;
}
.valeurs-list {
    flex: 1;
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 20px;
}
.valeur-item {
    display: flex;
    align-items: flex-start;
    gap: 16px;
    font-size: 0.88rem;
    line-height: 1.65;
}
.valeur-num {
    font-weight: 700;
    color: var(--green-mid);
    font-size: 0.82rem;
    flex-shrink: 0;
    min-width: 24px;
    margin-top: 1px;
}
.valeur-titre {
    font-weight: 700;
    color: var(--charcoal);
}
.valeur-texte {
    color: var(--charcoal);
    opacity: 0.8;
}

.story-section {
    padding: 0 0 80px;
}
.story-inner {
    display: flex;
    gap: 60px;
    align-items: center;
}
.story-img-wrap {
    flex: 0 0 400px;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
}
.story-img {
    width: 100%;
    height: 480px;
    object-fit: cover;
    display: block;
}
.story-content {
    flex: 1;
}
.story-title {
    font-size: clamp(1.6rem, 3vw, 2.2rem);
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 20px;
    letter-spacing: -0.02em;
}
.story-text {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.8;
    line-height: 1.7;
    margin: 0 0 16px;
}
.btn-join {
    display: inline-block;
    margin-top: 8px;
    padding: 13px 32px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    text-decoration: none;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-join:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}

@media (max-width: 900px) {
    .chiffres-grid {
        grid-template-columns: repeat(2, 1fr);
    }
    .valeurs-inner {
        flex-direction: column;
        gap: 40px;
    }
    .valeurs-header {
        flex: none;
    }
    .story-inner {
        flex-direction: column;
    }
    .story-img-wrap {
        flex: none;
        width: 100%;
    }
    .story-img {
        height: 280px;
    }
    .hero-img-wrap {
        height: 280px;
    }
}
@media (max-width: 560px) {
    .chiffres-grid {
        grid-template-columns: repeat(2, 1fr);
        gap: 20px;
    }
}
</style>
