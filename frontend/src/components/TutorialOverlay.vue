<script setup lang="ts">
import { ref, computed } from 'vue'

const emit = defineEmits<{ done: [] }>()

const currentStep = ref(0)

const steps = [
    {
        title: 'Bienvenue sur UpcycleConnect',
        description:
            'Votre espace particulier vous permet de gérer toutes vos activités de recyclage et d\'upcycling en un seul endroit.',
        icon: 'home',
    },
    {
        title: 'Votre tableau de bord',
        description:
            'Le tableau de bord vous donne une vue d\'ensemble de votre activité : score, annonces actives et dépôts en cours.',
        icon: 'dashboard',
    },
    {
        title: 'Déposer un objet',
        description:
            'Planifiez un créneau pour déposer vos objets dans nos conteneurs partenaires. Chaque dépôt contribue à votre score.',
        icon: 'deposit',
    },
    {
        title: 'Créer une annonce',
        description:
            'Mettez en vente ou proposez gratuitement vos objets upcyclés à la communauté via l\'espace Annonces.',
        icon: 'listing',
    },
    {
        title: 'Votre score upcycling',
        description:
            'Chaque action éco-responsable vous rapporte des points. Suivez votre progression et débloquez de nouveaux avantages.',
        icon: 'score',
    },
    {
        title: 'Le catalogue',
        description:
            'Inscrivez-vous aux ateliers et événements organisés par notre réseau d\'artisans et d\'experts du réemploi.',
        icon: 'catalogue',
    },
]

const isLast = computed(() => currentStep.value === steps.length - 1)
const progress = computed(() => ((currentStep.value + 1) / steps.length) * 100)

function next() {
    if (!isLast.value) {
        currentStep.value++
    } else {
        emit('done')
    }
}

function skip() {
    emit('done')
}
</script>

<template>
    <div class="overlay-backdrop">
        <div class="overlay-modal">
            <div class="overlay-progress-bar">
                <div class="overlay-progress-fill" :style="{ width: progress + '%' }"></div>
            </div>

            <div class="overlay-step-count">
                {{ currentStep + 1 }} / {{ steps.length }}
            </div>

            <div class="overlay-icon">
                <svg v-if="steps[currentStep].icon === 'home'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                    <polyline points="9 22 9 12 15 12 15 22" />
                </svg>
                <svg v-else-if="steps[currentStep].icon === 'dashboard'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <rect x="3" y="12" width="4" height="9" />
                    <rect x="10" y="7" width="4" height="14" />
                    <rect x="17" y="3" width="4" height="18" />
                </svg>
                <svg v-else-if="steps[currentStep].icon === 'deposit'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <polyline points="16 16 12 12 8 16" />
                    <line x1="12" y1="12" x2="12" y2="21" />
                    <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3" />
                </svg>
                <svg v-else-if="steps[currentStep].icon === 'listing'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                    <polyline points="14 2 14 8 20 8" />
                    <line x1="16" y1="13" x2="8" y2="13" />
                    <line x1="16" y1="17" x2="8" y2="17" />
                </svg>
                <svg v-else-if="steps[currentStep].icon === 'score'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <rect x="3" y="3" width="7" height="7" />
                    <rect x="14" y="3" width="7" height="7" />
                    <rect x="14" y="14" width="7" height="7" />
                    <rect x="3" y="14" width="7" height="7" />
                </svg>
            </div>

            <h2 class="overlay-title">{{ steps[currentStep].title }}</h2>
            <p class="overlay-description">{{ steps[currentStep].description }}</p>

            <div class="overlay-dots">
                <button
                    v-for="(_, i) in steps"
                    :key="i"
                    class="overlay-dot"
                    :class="{ 'overlay-dot--active': i === currentStep }"
                    @click="currentStep = i"
                ></button>
            </div>

            <div class="overlay-actions">
                <button class="btn-skip" @click="skip">Passer</button>
                <button class="btn-next" @click="next">
                    {{ isLast ? 'Commencer' : 'Suivant' }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.overlay-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(8, 106, 53, 0.55);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
}

.overlay-modal {
    background: var(--cream);
    border-radius: 18px;
    width: 100%;
    max-width: 480px;
    margin: 24px;
    padding: 0 0 32px;
    box-shadow: 0 24px 60px rgba(8, 106, 53, 0.22);
    overflow: hidden;
}

.overlay-progress-bar {
    height: 4px;
    background: var(--green-pale);
    border-radius: 0;
}
.overlay-progress-fill {
    height: 100%;
    background: var(--green-mid);
    transition: width 0.35s ease;
}

.overlay-step-count {
    text-align: right;
    font-size: 0.78rem;
    color: var(--green-light);
    font-weight: 600;
    padding: 16px 28px 0;
}

.overlay-icon {
    width: 64px;
    height: 64px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 12px auto 20px;
    color: var(--green-dark);
}
.overlay-icon svg {
    width: 28px;
    height: 28px;
}

.overlay-title {
    font-size: 1.4rem;
    font-weight: 800;
    color: var(--charcoal);
    text-align: center;
    margin: 0 28px 12px;
    line-height: 1.2;
    letter-spacing: -0.02em;
}

.overlay-description {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.72;
    text-align: center;
    line-height: 1.65;
    margin: 0 32px 24px;
}

.overlay-dots {
    display: flex;
    justify-content: center;
    gap: 8px;
    margin-bottom: 28px;
}
.overlay-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--green-pale);
    border: none;
    cursor: pointer;
    padding: 0;
    transition: background 0.2s, transform 0.2s;
}
.overlay-dot--active {
    background: var(--green-mid);
    transform: scale(1.3);
}

.overlay-actions {
    display: flex;
    gap: 12px;
    padding: 0 28px;
}
.btn-skip {
    flex: 1;
    padding: 12px;
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    color: rgba(53, 53, 53, 0.55);
    cursor: pointer;
    font-family: inherit;
    transition: border-color 0.2s, color 0.2s;
}
.btn-skip:hover {
    border-color: rgba(53, 53, 53, 0.4);
    color: var(--charcoal);
}
.btn-next {
    flex: 2;
    padding: 12px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
}
.btn-next:hover {
    background: var(--green-mid);
}
</style>
