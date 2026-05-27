<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useClientStore } from '@/stores/client'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const clientStore = useClientStore()
const router = useRouter()

interface Course {
    id: number
    nom: string
    description: string
    prix: number
    categorie: string
}

const dynamicCourses = ref<Course[]>([])

const filteredCourses = computed(() => {
    return dynamicCourses.value.filter(course => {
        // Check cart
        const inCart = clientStore.cart.some(cartItem => {
            const cartCourseId = cartItem.course_id && typeof cartItem.course_id === 'object' && 'Int64' in cartItem.course_id 
                ? Number(cartItem.course_id.Int64) 
                : Number(cartItem.course_id)
            return cartCourseId === Number(course.id)
        })
        if (inCart) return false

        // Check already registered
        const alreadyRegistered = clientStore.courseOrders.some(co => {
            const cCid = co.course_id && typeof co.course_id === 'object' ? co.course_id.Int64 : co.course_id
            return Number(cCid) === Number(course.id)
        })
        return !alreadyRegistered
    })
})

const showToast = ref(false)

onMounted(async () => {
    clientStore.fetchCart()
    clientStore.fetchCourseOrders()
    try {
        const res = await fetch(`${API_BASE}/course/catalogue`)
        if (res.ok) {
            dynamicCourses.value = await res.json()
        }
    } catch (e) {
        console.error('Failed to fetch courses:', e)
    }
})

async function handleAddToCart(course: Course) {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
    }

    try {
        await clientStore.addToCart({ courseId: course.id })
        showToast.value = true
        setTimeout(() => { showToast.value = false }, 3000)
    } catch (e: any) {
        alert("Erreur lors de l'ajout au panier")
    }
}

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
]
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container hero-inner">
                <h1 class="hero-title">Découvrez nos prestations.</h1>
                <p class="hero-subtitle">
                    Valorisons nos ressources pour un avenir plus
                    <span class="green-accent">durable</span>.
                </p>
            </div>
        </section>

        <section class="prestations-section">
            <div class="container">
                <div class="prestations-grid">
                    <div
                        v-for="(p, index) in prestations"
                        :key="p.id"
                        class="prestation-row"
                        :class="{ 'row-border': index < prestations.length - 2 }"
                    >
                        <div class="prestation-img-wrap">
                            <img :src="p.img" :alt="p.alt" class="prestation-img" />
                        </div>
                        <div class="prestation-content">
                            <h2 class="prestation-label">{{ p.label }}</h2>
                            <p class="prestation-desc">{{ p.description }}</p>
                            <button class="btn-reserver" @click="router.push('/reparer')">En savoir plus</button>
                        </div>
                    </div>

                    <div
                        v-for="course in filteredCourses"
                        :key="course.id"
                        class="prestation-row"
                    >
                        <div class="prestation-img-wrap">
                            <img src="https://images.unsplash.com/photo-1523301343968-6a6ebf63c672?w=600&q=80" alt="Se former" class="prestation-img" />
                        </div>
                        <div class="prestation-content">
                            <h2 class="prestation-label">SE FORMER : {{ course.nom }}</h2>
                            <p class="prestation-desc">{{ course.description }}</p>
                            <div class="course-meta">
                                <span class="course-price">{{ course.prix > 0 ? course.prix + ' €' : 'Gratuit' }}</span>
                                <span class="course-cat">{{ course.categorie }}</span>
                            </div>
                            <button class="btn-reserver" @click="handleAddToCart(course)">Ajouter au panier</button>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Toast Notification -->
        <Transition name="toast">
            <div v-if="showToast" class="toast-card">
                <div class="toast-content">
                    <span class="toast-icon">✅</span>
                    <span class="toast-text">Atelier ajouté au panier !</span>
                </div>
                <router-link to="/particulier/panier" class="toast-link">Voir panier</router-link>
            </div>
        </Transition>
    </div>
</template>

<style scoped>
/* Toast Styles */
.toast-card {
    position: fixed;
    bottom: 30px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--white);
    border: 1.5px solid var(--green-mid);
    border-radius: 12px;
    padding: 12px 20px;
    display: flex;
    align-items: center;
    gap: 20px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.1);
    z-index: 2000;
}
.toast-content { display: flex; align-items: center; gap: 10px; }
.toast-text { font-size: 0.9rem; font-weight: 600; }
.toast-link { color: var(--green-dark); font-weight: 700; font-size: 0.85rem; text-decoration: underline; }
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px); }

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

.prestations-section {
    padding: 8px 0 80px;
    flex: 1;
}

.prestations-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0;
}

.prestation-row {
    display: flex;
    align-items: flex-start;
    gap: 0;
    padding: 36px 32px 36px 0;
    border-bottom: 1px solid rgba(53, 53, 53, 0.1);
}

.prestation-row:nth-child(3),
.prestation-row:nth-child(4) {
    border-bottom: none;
}

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

.course-meta {
    display: flex;
    gap: 12px;
    margin-bottom: 12px;
    font-size: 0.8rem;
    font-weight: 600;
}
.course-price {
    color: var(--green-dark);
}
.course-cat {
    opacity: 0.6;
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
    text-decoration: none;
    display: inline-block;
}
.btn-reserver:hover {
    background: var(--green-dark);
    transform: translateY(-1px);
}

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
}
</style>
