<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'
import { useAuthStore } from '@/stores/auth'

const clientStore = useClientStore()
const authStore = useAuthStore()
const router = useRouter()
type Tab = 'all' | 'events' | 'courses' | 'annonces'
const activeTab = ref<Tab>('all')

const showToast = ref(false)
const toastMessage = ref('')

const allItems = computed(() => {
    const evts = clientStore.events.map((e: any) => ({ ...e, _type: 'event' }))
    const crs = clientStore.courses.map((c: any) => ({ ...c, _type: 'course' }))
    const ans = clientStore.allAnnonces.map((a: any) => ({ ...a, _type: 'annonce' }))
    return [...evts, ...crs, ...ans]
})

function getNumericId(val: any): number | null {
    if (val === null || val === undefined) return null
    if (typeof val === 'object' && 'Int64' in val) return Number(val.Int64)
    if (typeof val === 'object' && 'id' in val) return Number(val.id)
    const n = Number(val)
    return isNaN(n) ? null : n
}

const filtered = computed(() => {
    try {
        let items = allItems.value || []
        if (activeTab.value === 'events') items = items.filter(i => i._type === 'event')
        else if (activeTab.value === 'courses') items = items.filter(i => i._type === 'course')
        else if (activeTab.value === 'annonces') items = items.filter(i => i._type === 'annonce')

        const cartItems = clientStore.cart || []
        const userParticipations = clientStore.participations || []
        const userCourseOrders = clientStore.courseOrders || []

        return items.filter(item => {
            const id = getNumericId(item.id)
            if (id === null) return false
            
            // Check if in cart
            const inCart = cartItems.some(cartItem => {
                if (item._type === 'annonce') return getNumericId(cartItem.listing_id) === id
                if (item._type === 'event') return getNumericId(cartItem.event_id) === id
                if (item._type === 'course') return getNumericId(cartItem.course_id) === id
                return false
            })
            if (inCart) return false

            // Check if already registered (for events and courses)
            if (item._type === 'event') {
                return !userParticipations.some(p => getNumericId(p.event_id) === id)
            }
            if (item._type === 'course') {
                return !userCourseOrders.some(co => getNumericId(co.course_id) === id)
            }

            return true
        })
    } catch (e) {
        console.error('Filtering error:', e)
        return []
    }
})

function getItemId(item: any): number {
    return getNumericId(item.id) || 0
}

function formatDate(ts: any): string {
    if (!ts) return '—'
    const date = new Date(ts.Time ?? ts)
    return date.toLocaleDateString('fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}

function formatPrice(price: any): string {
    if (!price) return 'Gratuit'
    const val = typeof price === 'object' ? price.Float64 ?? price.Int64 : price
    const num = Number(val)
    return num > 0 ? `${num.toFixed(2)} €` : 'Gratuit'
}

function getItemName(item: any): string {
    const name = item.name ?? item.nom
    if (name) return name
    if (item._type === 'event') {
        const d = item.date ?? item.created_at
        return d ? `Événement du ${formatDate(d)}` : 'Événement'
    }
    return 'Atelier sans titre'
}

function getItemDesc(item: any): string {
    const desc = item.description ?? item.desc
    if (desc) return desc
    if (item._type === 'event') return 'Aucune description disponible for cet événement.'
    return 'Aucune description disponible for cet atelier.'
}

function getItemPrice(item: any): any {
    return item.price ?? item.prix
}

function handleContactSeller(item: any) {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
    }
    
    let listingId = ''
    const id = getItemId(item)
    if (typeof id === 'object' && id !== null) {
        listingId = id.Int64?.toString() || id.id?.toString() || ''
    } else {
        listingId = id?.toString() || ''
    }

    if (!listingId) {
        console.error('No listing ID found')
        return
    }

    router.push({
        path: '/particulier/chat',
        query: { listingId }
    })
}

async function handleAddToCart(item: any) {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
    }
    
    try {
        const id = getItemId(item)
        const payload: any = {}
        if (item._type === 'annonce') payload.listingId = id
        else if (item._type === 'event') payload.eventId = id
        else if (item._type === 'course') payload.courseId = id
        
        await clientStore.addToCart(payload)
        toastMessage.value = `"${getItemName(item)}" ajouté au panier !`
        showToast.value = true
        setTimeout(() => { showToast.value = false }, 3000)
    } catch (e: any) {
        toastMessage.value = "Erreur lors de l'ajout"
        showToast.value = true
        setTimeout(() => { showToast.value = false }, 3000)
    }
}

onMounted(() => {
    clientStore.fetchCart()
    clientStore.fetchCatalogue()
    clientStore.fetchAllAnnonces()
    if (authStore.isAuthenticated) {
        clientStore.fetchConversations()
        clientStore.fetchParticipations()
        clientStore.fetchCourseOrders()
    }
})
</script>

<template>
    <div class="page">
        <h1 class="page-title">Catalogue.</h1>

        <div class="tabs">
            <button
                class="tab-btn"
                :class="{ 'tab-btn--active': activeTab === 'all' }"
                @click="activeTab = 'all'"
            >
                Tout ({{ allItems.length }})
            </button>
            <button
                class="tab-btn"
                :class="{ 'tab-btn--active': activeTab === 'events' }"
                @click="activeTab = 'events'"
            >
                Événements ({{ clientStore.events.length }})
            </button>
            <button
                class="tab-btn"
                :class="{ 'tab-btn--active': activeTab === 'courses' }"
                @click="activeTab = 'courses'"
            >
                Ateliers ({{ clientStore.courses.length }})
            </button>
            <button
                class="tab-btn"
                :class="{ 'tab-btn--active': activeTab === 'annonces' }"
                @click="activeTab = 'annonces'"
            >
                Annonces ({{ clientStore.allAnnonces.length }})
            </button>
        </div>

        <div v-if="clientStore.isLoading" class="state-empty">
            <p>Chargement…</p>
        </div>

        <div v-else-if="filtered.length === 0" class="state-empty">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="7" height="7" />
                    <rect x="14" y="3" width="7" height="7" />
                    <rect x="14" y="14" width="7" height="7" />
                    <rect x="3" y="14" width="7" height="7" />
                </svg>
            </div>
            <p class="empty-title">Aucun élément disponible</p>
            <p class="empty-sub">Revenez bientôt for découvrir nos prochains événements et ateliers.</p>
        </div>

        <div v-else class="catalogue-grid">
            <div
                v-for="item in filtered"
                :key="`${item._type}-${getItemId(item)}`"
                class="catalogue-card"
            >
                <div v-if="item.image_url?.String" class="card-img-wrap">
                    <img :src="`${API_BASE}` + item.image_url.String" alt="" class="card-img" />
                </div>
                <div v-else-if="item.image_url && typeof item.image_url === 'string'" class="card-img-wrap">
                    <img :src="`${API_BASE}` + item.image_url" alt="" class="card-img" />
                </div>
                <div class="card-header">
                    <span class="type-badge" :class="{
                        'type-badge--event': item._type === 'event',
                        'type-badge--course': item._type === 'course',
                        'type-badge--annonce': item._type === 'annonce'
                    }">
                        {{ item._type === 'event' ? 'Événement' : item._type === 'course' ? 'Atelier' : 'Annonce' }}
                    </span>
                    <span class="card-price">{{ formatPrice(getItemPrice(item)) }}</span>
                </div>

                <h3 class="card-name">{{ getItemName(item) }}</h3>
                <p class="card-desc">{{ getItemDesc(item) }}</p>

                <div class="card-footer">
                    <span v-if="item.created_at" class="card-date">
                        {{ formatDate(item.created_at) }}
                    </span>
                    <span v-if="item.max_capacity" class="card-capacity">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                            <circle cx="9" cy="7" r="4" />
                        </svg>
                        {{ item.max_capacity?.Int32 ?? item.max_capacity }} places
                    </span>
                </div>

                <div class="card-actions">
                    <template v-if="item._type === 'annonce'">
                        <button v-if="clientStore.isChattingWith(getItemId(item))" class="btn-chat btn-chat--active" @click="handleContactSeller(item)">
                            💬 Continuer la discussion
                        </button>
                        <button v-else class="btn-cart btn-cart--full" @click="handleAddToCart(item)">
                            🛒 Ajouter au panier
                        </button>
                    </template>
                    <button v-else class="btn-cart btn-cart--full" @click="handleAddToCart(item)">
                        🛒 Ajouter au panier
                    </button>
                </div>
            </div>
        </div>

        <!-- Toast Notification -->
        <Transition name="toast">
            <div v-if="showToast" class="toast-card">
                <div class="toast-content">
                    <span class="toast-icon">✅</span>
                    <span class="toast-text">{{ toastMessage }}</span>
                </div>
                <router-link to="/particulier/panier" class="toast-link">Voir panier</router-link>
            </div>
        </Transition>
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
    margin: 0 0 28px;
    line-height: 1.08;
}

.tabs {
    display: flex;
    gap: 8px;
    margin-bottom: 28px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.1);
    padding-bottom: 0;
}
.tab-btn {
    padding: 10px 18px;
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    font-size: 0.875rem;
    font-weight: 600;
    color: rgba(53, 53, 53, 0.5);
    cursor: pointer;
    font-family: inherit;
    transition: color 0.2s, border-color 0.2s;
}
.tab-btn:hover {
    color: var(--charcoal);
}
.tab-btn--active {
    color: var(--green-dark);
    border-bottom-color: var(--green-dark);
}

.state-empty {
    text-align: center;
    padding: 64px 32px;
}
.empty-icon {
    width: 64px;
    height: 64px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    color: var(--green-mid);
}
.empty-icon svg {
    width: 28px;
    height: 28px;
}
.empty-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 8px;
}
.empty-sub {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.6;
    max-width: 320px;
    margin: 0 auto;
    line-height: 1.6;
}

.catalogue-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
}
.catalogue-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    transition: border-color 0.2s, transform 0.2s;
}
.catalogue-card:hover {
    border-color: var(--green-light);
    transform: translateY(-2px);
}
.card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.type-badge {
    display: inline-block;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
}
.type-badge--event {
    background: rgba(8, 106, 53, 0.1);
    color: var(--green-dark);
}
.type-badge--course {
    background: var(--green-pale);
    color: var(--green-mid);
}
.type-badge--annonce {
    background: #fef3c7;
    color: #92400e;
}
.card-img-wrap {
    width: calc(100% + 40px);
    margin: -20px -20px 12px;
    height: 160px;
    overflow: hidden;
}
.card-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.card-price {
    font-size: 0.9rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
}
.card-name {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0;
    line-height: 1.3;
}
.card-desc {
    font-size: 0.82rem;
    color: var(--charcoal);
    opacity: 0.65;
    line-height: 1.6;
    margin: 0;
    flex: 1;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.card-footer {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
}
.card-date {
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.45;
}
.card-capacity {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 0.75rem;
    color: var(--charcoal);
    opacity: 0.45;
}
.card-capacity svg {
    width: 13px;
    height: 13px;
}
.card-actions {
    display: flex;
    gap: 8px;
    margin-top: auto;
}
.btn-cart {
    width: 100%;
    padding: 11px;
    background: #f3f4f6;
    color: var(--charcoal);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: all 0.2s;
}
.btn-cart:hover {
    background: #e5e7eb;
}
.btn-cart--full {
    background: var(--green-pale);
    color: var(--green-dark);
    border-color: rgba(8, 106, 53, 0.1);
}
.btn-cart--full:hover {
    background: var(--green-mid);
    color: var(--white);
}
.btn-chat--active {
    width: 100%;
    padding: 11px;
    background: #4183d7;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: all 0.2s;
}
.btn-chat--active:hover {
    background: #3569ad;
}

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
.toast-content {
    display: flex;
    align-items: center;
    gap: 10px;
}
.toast-text {
    font-size: 0.9rem;
    font-weight: 600;
}
.toast-link {
    color: var(--green-dark);
    font-weight: 700;
    font-size: 0.85rem;
    text-decoration: underline;
}

.toast-enter-active, .toast-leave-active {
    transition: all 0.3s ease;
}
.toast-enter-from, .toast-leave-to {
    opacity: 0;
    transform: translate(-50%, 20px);
}
</style>