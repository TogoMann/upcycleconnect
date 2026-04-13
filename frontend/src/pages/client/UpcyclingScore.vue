<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()

const score = computed(() => clientStore.score)

const level = computed(() => {
    if (score.value >= 500) return { label: 'Expert', color: '#086a35', next: null, threshold: 500 }
    if (score.value >= 200) return { label: 'Avancé', color: '#34895b', next: 500, threshold: 200 }
    if (score.value >= 50) return { label: 'Intermédiaire', color: '#8bbd94', next: 200, threshold: 50 }
    return { label: 'Débutant', color: '#d7ece1', next: 50, threshold: 0 }
})

const progressToNextLevel = computed(() => {
    if (!level.value.next) return 100
    const range = level.value.next - level.value.threshold
    const current = score.value - level.value.threshold
    return Math.min(100, Math.round((current / range) * 100))
})

const milestones = [
    { label: 'Premier dépôt', points: 10, icon: 'deposit' },
    { label: 'Première annonce', points: 15, icon: 'listing' },
    { label: 'Score 50 pts', points: 50, icon: 'star' },
    { label: 'Score 200 pts', points: 200, icon: 'trophy' },
    { label: 'Score 500 pts', points: 500, icon: 'expert' },
]

onMounted(() => {
    clientStore.fetchScore()
})
</script>

<template>
    <div class="page">
        <h1 class="page-title">Mon Score.</h1>

        <div class="score-hero">
            <div class="score-circle">
                <span class="score-number">{{ score }}</span>
                <span class="score-unit">pts</span>
            </div>
            <div class="score-level-info">
                <span class="level-badge" :style="{ background: level.color + '22', color: level.color }">
                    Niveau : {{ level.label }}
                </span>
                <p class="score-desc">
                    Chaque action éco-responsable vous rapporte des points.
                    Déposez des objets, créez des annonces, participez à des événements pour progresser.
                </p>
                <div v-if="level.next" class="progress-section">
                    <div class="progress-labels">
                        <span>{{ score }} pts</span>
                        <span>{{ level.next }} pts pour niveau suivant</span>
                    </div>
                    <div class="progress-bar">
                        <div class="progress-fill" :style="{ width: progressToNextLevel + '%' }"></div>
                    </div>
                    <p class="progress-caption">
                        {{ progressToNextLevel }}% vers le niveau <strong>suivant</strong>
                    </p>
                </div>
                <div v-else class="progress-section">
                    <div class="progress-bar">
                        <div class="progress-fill progress-fill--max"></div>
                    </div>
                    <p class="progress-caption">Niveau maximum atteint — vous êtes un expert de l'upcycling !</p>
                </div>
            </div>
        </div>

        <div class="how-section">
            <h2 class="section-title">Comment gagner des points ?</h2>
            <div class="how-grid">
                <div class="how-card">
                    <div class="how-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <polyline points="16 16 12 12 8 16" />
                            <line x1="12" y1="12" x2="12" y2="21" />
                            <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3" />
                        </svg>
                    </div>
                    <div>
                        <strong class="how-title">Déposer un objet</strong>
                        <p class="how-pts">+ 10 pts par dépôt validé</p>
                    </div>
                </div>

                <div class="how-card">
                    <div class="how-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                            <polyline points="14 2 14 8 20 8" />
                        </svg>
                    </div>
                    <div>
                        <strong class="how-title">Créer une annonce</strong>
                        <p class="how-pts">+ 15 pts par annonce publiée</p>
                    </div>
                </div>

                <div class="how-card">
                    <div class="how-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                            <line x1="16" y1="2" x2="16" y2="6" />
                            <line x1="8" y1="2" x2="8" y2="6" />
                            <line x1="3" y1="10" x2="21" y2="10" />
                        </svg>
                    </div>
                    <div>
                        <strong class="how-title">Participer à un atelier</strong>
                        <p class="how-pts">+ 25 pts par atelier suivi</p>
                    </div>
                </div>

                <div class="how-card">
                    <div class="how-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                            <circle cx="9" cy="7" r="4" />
                            <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                        </svg>
                    </div>
                    <div>
                        <strong class="how-title">Parrainer un ami</strong>
                        <p class="how-pts">+ 50 pts par parrainage</p>
                    </div>
                </div>
            </div>
        </div>

        <div class="milestones-section">
            <h2 class="section-title">Jalons</h2>
            <div class="milestones-list">
                <div
                    v-for="m in milestones"
                    :key="m.label"
                    class="milestone-item"
                    :class="{ 'milestone-item--done': score >= m.points }"
                >
                    <div class="milestone-check">
                        <svg v-if="score >= m.points" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                            <polyline points="20 6 9 17 4 12" />
                        </svg>
                        <span v-else class="milestone-pts-sm">{{ m.points }}</span>
                    </div>
                    <span class="milestone-label">{{ m.label }}</span>
                    <span class="milestone-pts">{{ m.points }} pts</span>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 32px;
    line-height: 1.08;
}

.score-hero {
    display: flex;
    align-items: flex-start;
    gap: 40px;
    margin-bottom: 48px;
    flex-wrap: wrap;
}
.score-circle {
    flex-shrink: 0;
    width: 140px;
    height: 140px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 4px solid var(--green-mid);
}
.score-number {
    font-size: 2.4rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.04em;
    line-height: 1;
}
.score-unit {
    font-size: 0.85rem;
    color: var(--green-mid);
    font-weight: 600;
}
.score-level-info {
    flex: 1;
    min-width: 240px;
}
.level-badge {
    display: inline-block;
    padding: 5px 14px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 700;
    margin-bottom: 12px;
}
.score-desc {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.7;
    line-height: 1.65;
    margin: 0 0 20px;
}
.progress-section {}
.progress-labels {
    display: flex;
    justify-content: space-between;
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.55;
    margin-bottom: 6px;
}
.progress-bar {
    height: 8px;
    background: rgba(53, 53, 53, 0.12);
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 6px;
}
.progress-fill {
    height: 100%;
    background: var(--green-mid);
    border-radius: 4px;
    transition: width 0.5s ease;
}
.progress-fill--max {
    width: 100%;
    background: var(--green-dark);
}
.progress-caption {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.5;
    margin: 0;
}

.section-title {
    font-size: 1.05rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 16px;
    letter-spacing: -0.01em;
}

.how-section {
    margin-bottom: 40px;
}
.how-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
}
.how-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 12px;
    padding: 16px 18px;
    display: flex;
    align-items: flex-start;
    gap: 14px;
}
.how-icon {
    width: 36px;
    height: 36px;
    background: var(--green-pale);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--green-mid);
    flex-shrink: 0;
}
.how-icon svg {
    width: 18px;
    height: 18px;
}
.how-title {
    font-size: 0.875rem;
    font-weight: 700;
    color: var(--charcoal);
    display: block;
    margin-bottom: 3px;
}
.how-pts {
    font-size: 0.78rem;
    color: var(--green-mid);
    margin: 0;
    font-weight: 600;
}

.milestones-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.milestone-item {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 18px;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 10px;
    opacity: 0.6;
}
.milestone-item--done {
    opacity: 1;
    border-color: var(--green-light);
    background: var(--green-pale);
}
.milestone-check {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: rgba(53, 53, 53, 0.08);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: var(--green-dark);
    font-size: 0.7rem;
    font-weight: 700;
}
.milestone-item--done .milestone-check {
    background: var(--green-mid);
    color: var(--white);
}
.milestone-check svg {
    width: 14px;
    height: 14px;
}
.milestone-pts-sm {
    font-size: 0.65rem;
    font-weight: 700;
    color: var(--charcoal);
    opacity: 0.55;
}
.milestone-label {
    flex: 1;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--charcoal);
}
.milestone-pts {
    font-size: 0.78rem;
    color: var(--green-mid);
    font-weight: 700;
}

@media (max-width: 640px) {
    .score-hero {
        flex-direction: column;
        align-items: center;
        text-align: center;
    }
    .progress-labels {
        flex-direction: column;
        gap: 2px;
    }
    .how-grid {
        grid-template-columns: 1fr;
    }
}
</style>
