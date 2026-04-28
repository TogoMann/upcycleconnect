<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()
const router = useRouter()
type Tab = 'all' | 'events' | 'courses' | 'annonces'
const activeTab = ref<Tab>('all')

const allItems = computed(() => {
    const evts = clientStore.events.map((e: any) => ({ ...e, _type: 'event' }))
    const crs = clientStore.courses.map((c: any) => ({ ...c, _type: 'course' }))
    const ans = clientStore.allAnnonces.map((a: any) => ({ ...a, _type: 'annonce' }))
    return [...evts, ...crs, ...ans]
})

const filtered = computed(() => {
    if (activeTab.value === 'events') return allItems.value.filter(i => i._type === 'event')
    if (activeTab.value === 'courses') return allItems.value.filter(i => i._type === 'course')
    if (activeTab.value === 'annonces') return allItems.value.filter(i => i._type === 'annonce')
    return allItems.value
})

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

function getItemId(item: any): number {
    if (item.id && typeof item.id === 'object' && 'Int64' in item.id) {
        return item.id.Int64
    }
    return Number(item.id)
}

function getItemPrice(item: any): any {
    return item.price ?? item.prix
}

function handleBuy(item: any) {
    const price = getItemPrice(item)
    router.push({
        path: '/particulier/paiement',
        query: {
            id: getItemId(item),
            name: getItemName(item),
            price: typeof price === 'object' ? price?.Float64 ?? price?.Int64 : price,
            type: item._type,
        },
    })
}

onMounted(() => {
    clientStore.fetchCatalogue()
    clientStore.fetchAllAnnonces()
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
                <div v-if="item.image_url" class="card-img-wrap">
                    <img :src="'http://localhost:8081' + item.image_url" alt="" class="card-img" />
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

                <button class="btn-book" @click="handleBuy(item)">
                    {{ formatPrice(getItemPrice(item)) === 'Gratuit' ? "S'inscrire" : "Acheter" }}
                </button>
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
.btn-book {
    width: 100%;
    padding: 11px;
    background: var(--green-pale);
    color: var(--green-dark);
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s, color 0.2s;
    margin-top: auto;
}
.btn-book:hover {
    background: var(--green-dark);
    color: var(--white);
}
</style>