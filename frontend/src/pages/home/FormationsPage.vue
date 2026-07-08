<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const clientStore = useClientStore()
const authStore = useAuthStore()
const router = useRouter()

const filterType = ref<'' | 'presentiel' | 'en_ligne'>('')
const showToast = ref(false)
const toastMessage = ref('')

function getNumericId(val: any): number | null {
    if (val === null || val === undefined) return null
    if (typeof val === 'object' && 'Int64' in val) return Number(val.Int64)
    if (typeof val === 'object' && 'id' in val) return Number(val.id)
    const n = Number(val)
    return isNaN(n) ? null : n
}

const filteredCourses = computed(() => {
    let items = clientStore.courses || []
    if (filterType.value) items = items.filter((c: any) => c.type === filterType.value)

    const cartItems = clientStore.cart || []
    const userCourseOrders = clientStore.courseOrders || []

    return items.filter((course: any) => {
        const id = getNumericId(course.id)
        if (id === null) return false

        const inCart = cartItems.some((cartItem: any) => getNumericId(cartItem.course_id) === id)
        if (inCart) return false

        return !userCourseOrders.some((co: any) => getNumericId(co.course_id) === id)
    })
})

function rawDateVal(d: any): string | null {
    if (!d) return null
    if (typeof d === 'object') return d.Valid ? d.Time : null
    return d
}

function formatDate(d: any): string {
    if (!d) return '—'
    return new Date(d).toLocaleDateString(undefined, { day: '2-digit', month: 'long', year: 'numeric' })
}

function getDateRange(course: any): string | null {
    const start = rawDateVal(course.date)
    if (!start) return null
    const end = rawDateVal(course.end_date)
    if (!end || end === start) return formatDate(start)
    return t('formationsPage.dateRange', { start: formatDate(start), end: formatDate(end) })
}

function formatPrice(price: any): string {
    const num = Number(price)
    return num > 0 ? `${num.toFixed(2)} €` : t('listingDetail.free')
}

async function handleJoin(course: any) {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
    }
    try {
        await clientStore.addToCart({ courseId: course.id })
        toastMessage.value = t('formationsPage.joined', { name: course.nom })
        showToast.value = true
        setTimeout(() => { showToast.value = false }, 3000)
    } catch {
        toastMessage.value = t('formationsPage.joinError')
        showToast.value = true
        setTimeout(() => { showToast.value = false }, 3000)
    }
}

function handleContact(course: any) {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
    }
    router.push({ path: '/particulier/chat', query: { courseId: String(course.id) } })
}

onMounted(() => {
    clientStore.fetchCatalogue()
    clientStore.fetchCart()
    if (authStore.isAuthenticated) {
        clientStore.fetchCourseOrders()
    }
})
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container hero-inner">
                <h1 class="hero-title">{{ t('formationsPage.pageTitle') }}</h1>
                <p class="hero-subtitle">{{ t('formationsPage.subtitle') }}</p>
            </div>
        </section>

        <section class="formations-section">
            <div class="container">
                <div class="filter-chips">
                    <button class="filter-chip" :class="{ active: filterType === '' }" @click="filterType = ''">{{ t('formationsPage.filterAll') }}</button>
                    <button class="filter-chip" :class="{ active: filterType === 'presentiel' }" @click="filterType = 'presentiel'">{{ t('formationsPage.filterPresentiel') }}</button>
                    <button class="filter-chip" :class="{ active: filterType === 'en_ligne' }" @click="filterType = 'en_ligne'">{{ t('formationsPage.filterEnLigne') }}</button>
                </div>

                <div v-if="clientStore.isLoading" class="state-empty">
                    <p>{{ t('formationsPage.loading') }}</p>
                </div>

                <div v-else-if="filteredCourses.length === 0" class="state-empty">
                    <p>{{ t('formationsPage.empty') }}</p>
                </div>

                <div v-else class="formations-grid">
                    <div v-for="course in filteredCourses" :key="course.id" class="formation-card">
                        <div class="card-header">
                            <span class="type-badge">
                                {{ course.type === 'en_ligne' ? t('formationsPage.enLigne') : t('formationsPage.presentiel') }}
                            </span>
                            <span class="card-price">{{ formatPrice(course.prix) }}</span>
                        </div>
                        <h3 class="card-name">{{ course.nom }}</h3>
                        <p v-if="course.organisateur" class="card-organizer">{{ t('formationsPage.by', { name: course.organisateur }) }}</p>
                        <p class="card-desc">{{ course.description }}</p>
                        <div class="card-footer" v-if="getDateRange(course)">
                            <span class="card-date">{{ getDateRange(course) }}</span>
                        </div>
                        <div class="card-actions">
                            <button class="btn-join btn-split" @click="handleJoin(course)">{{ t('formationsPage.join') }}</button>
                            <button class="btn-contact btn-split" @click="handleContact(course)">{{ t('formationsPage.contactOrganizer') }}</button>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <Transition name="toast">
            <div v-if="showToast" class="toast-card">
                <div class="toast-content">
                    <span class="toast-icon">✅</span>
                    <span class="toast-text">{{ toastMessage }}</span>
                </div>
                <router-link to="/particulier/panier" class="toast-link">{{ t('events.viewCart') }}</router-link>
            </div>
        </Transition>
    </div>
</template>

<style scoped>
.page-content { flex: 1; display: flex; flex-direction: column; }
.container { max-width: 1060px; margin: 0 auto; padding: 0 32px; }
.hero { background: var(--cream); padding: 72px 0 52px; text-align: center; }
.hero-inner { display: flex; flex-direction: column; align-items: center; }
.hero-title { font-size: clamp(2.4rem, 5vw, 3.8rem); font-weight: 800; color: var(--charcoal); line-height: 1.1; letter-spacing: -0.025em; margin: 0 0 18px; }
.hero-subtitle { font-size: clamp(1.1rem, 2.2vw, 1.5rem); font-weight: 400; color: var(--charcoal); margin: 0; line-height: 1.4; max-width: 640px; }

.formations-section { padding: 40px 0 80px; flex: 1; }

.filter-chips { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 28px; }
.filter-chip { padding: 8px 18px; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 20px; background: none; font-size: 0.85rem; font-weight: 600; cursor: pointer; transition: all 0.15s; font-family: inherit; color: var(--charcoal); }
.filter-chip:hover { border-color: var(--green-mid); color: var(--green-dark); }
.filter-chip.active { background: var(--green-dark); color: white; border-color: var(--green-dark); }

.state-empty { text-align: center; padding: 64px 32px; opacity: 0.6; }

.formations-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.formation-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 14px; padding: 20px; display: flex; flex-direction: column; gap: 12px; transition: border-color 0.2s, transform 0.2s; }
.formation-card:hover { border-color: var(--green-light); transform: translateY(-2px); }
.card-header { display: flex; align-items: center; justify-content: space-between; }
.type-badge { display: inline-block; padding: 3px 10px; border-radius: 20px; font-size: 0.7rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase; background: var(--green-pale); color: var(--green-mid); }
.card-price { font-size: 0.9rem; font-weight: 800; color: var(--green-dark); letter-spacing: -0.02em; }
.card-name { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0; line-height: 1.3; }
.card-organizer { font-size: 0.78rem; color: var(--green-dark); font-weight: 600; margin: -8px 0 0; }
.card-desc { font-size: 0.82rem; color: var(--charcoal); opacity: 0.65; line-height: 1.6; margin: 0; flex: 1; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
.card-footer { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; }
.card-date { font-size: 0.75rem; color: var(--charcoal); opacity: 0.45; }
.card-actions { display: flex; gap: 8px; margin-top: auto; }
.btn-split { width: auto; flex: 1; }
.btn-join { padding: 11px; background: var(--green-pale); color: var(--green-dark); border: 1.5px solid rgba(8,106,53,0.1); border-radius: 8px; font-size: 0.85rem; font-weight: 700; cursor: pointer; font-family: inherit; transition: all 0.2s; }
.btn-join:hover { background: var(--green-mid); color: var(--white); }
.btn-contact { padding: 11px; background: #eef4fb; color: #4183d7; border: 1.5px solid rgba(65,131,215,0.25); border-radius: 8px; font-size: 0.85rem; font-weight: 700; cursor: pointer; font-family: inherit; transition: all 0.2s; }
.btn-contact:hover { background: #4183d7; color: white; }

.toast-card { position: fixed; bottom: 30px; left: 50%; transform: translateX(-50%); background: var(--white); border: 1.5px solid var(--green-mid); border-radius: 12px; padding: 12px 20px; display: flex; align-items: center; gap: 20px; box-shadow: 0 10px 25px rgba(0,0,0,0.1); z-index: 2000; }
.toast-content { display: flex; align-items: center; gap: 10px; }
.toast-text { font-size: 0.9rem; font-weight: 600; }
.toast-link { color: var(--green-dark); font-weight: 700; font-size: 0.85rem; text-decoration: underline; }
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, 20px); }

@media (max-width: 560px) {
    .formations-grid { grid-template-columns: 1fr; }
    .hero-title { font-size: 2rem; }
}
</style>
