<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useClientStore } from '@/stores/client'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const clientStore = useClientStore()

const listing = ref<any>(null)
const loading = ref(true)

onMounted(async () => {
    try {
        const id = route.params.id
        const res = await fetch(`${API_BASE}/listing/${id}`)
        if (res.ok) {
            listing.value = await res.json()
        }
        if (authStore.isAuthenticated) {
            await clientStore.fetchConversations()
        }
    } catch (e) {
        console.error('Failed to fetch listing:', e)
    } finally {
        loading.value = false
    }
})

const isFree = computed(() => {
    if (!listing.value) return true
    const p = listing.value.price
    const val = typeof p === 'object' ? p.Float64 ?? p.Int64 : p
    return Number(val) === 0
})

const formattedPrice = computed(() => {
    if (!listing.value) return '0.00'
    const p = listing.value.price
    const val = typeof p === 'object' ? p.Float64 ?? p.Int64 : p
    return Number(val).toFixed(2)
})

const formattedDate = computed(() => {
    if (!listing.value?.created_at) return '—'
    return new Date(listing.value.created_at).toLocaleDateString('fr-FR', {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
    })
})

const mainImage = computed(() => {
    if (!listing.value?.image_url) return 'https://images.unsplash.com/photo-1592078615290-033ee584e267?w=800&q=80'
    return `${API_BASE}${listing.value.image_url}`
})

const isLockerHandoff = computed(() => listing.value?.handoff_mode === 'casier')

const sellerLevel = computed(() => {
    if (!listing.value) return 1
    return listing.value.seller_level || Math.floor(Math.max(listing.value.seller_score || 0, 0) / 100) + 1
})

const sellerName = computed(() => listing.value?.created_by_name || t('listings.unknownUser'))
const sellerInitial = computed(() => sellerName.value.charAt(0).toUpperCase())

function getItemId(item: any): number {
    if (!item) return 0
    if (item.id && typeof item.id === 'object' && 'Int64' in item.id) {
        return item.id.Int64
    }
    return Number(item.id)
}

function handleAction() {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
    }

    if (!listing.value) return

    let listingId = ''
    if (typeof listing.value.id === 'object' && listing.value.id !== null) {
        listingId = listing.value.id.Int64?.toString() || listing.value.id.id?.toString() || ''
    } else {
        listingId = listing.value.id?.toString() || ''
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
</script>

<template>
    <div class="page-content">
        <div v-if="loading" class="container loading-state">
            {{ t('listingDetail.loading') }}
        </div>

        <div v-else-if="!listing" class="container error-state">
            {{ t('listingDetail.notFound') }}
        </div>

        <template v-else>
            <section class="breadcrumb-bar">
                <div class="container">
                    <router-link to="/annonces" class="breadcrumb-link">{{ t('layout.nav.listings') }}</router-link>
                    <span class="breadcrumb-sep">›</span>
                    <span class="breadcrumb-current">{{ listing.name }}</span>
                </div>
            </section>

            <section class="detail-section">
                <div class="container">
                    <div class="detail-layout">
                        <div class="detail-left">
                            <div class="main-img-wrap">
                                <img :src="mainImage" :alt="listing.name" class="main-img" />
                            </div>
                        </div>

                        <div class="detail-right">
                            <div class="detail-badges">
                                <span class="badge" :class="isFree ? 'badge--don' : 'badge--vente'">
                                    {{ isFree ? t('listings.filterDon') : t('listings.filterVente') }}
                                </span>
                                <span class="badge badge--cat">{{ listing.category || t('listingDetail.uncategorized') }}</span>
                            </div>

                            <h1 class="detail-titre">{{ listing.name }}</h1>

                            <div v-if="!isFree" class="detail-prix">
                                {{ formattedPrice }} €
                            </div>
                            <div v-else class="detail-gratuit">{{ t('listingDetail.free') }}</div>

                            <div class="detail-meta">
                                <div class="meta-item">
                                    <svg viewBox="0 0 24 24" fill="currentColor" width="14" height="14">
                                        <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
                                    </svg>
                                    <span>{{ listing.city_name || t('listingDetail.unknownLocation') }}</span>
                                </div>
                                <div class="meta-item">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                                        <line x1="16" y1="2" x2="16" y2="6"/>
                                        <line x1="8" y1="2" x2="8" y2="6"/>
                                        <line x1="3" y1="10" x2="21" y2="10"/>
                                    </svg>
                                    <span>{{ t('listingDetail.publishedOn', { date: formattedDate }) }}</span>
                                </div>
                                <div class="meta-item">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                                        <rect x="1" y="7" width="22" height="13" rx="2" ry="2"/>
                                        <path d="M16 3H8a2 2 0 0 0-2 2v2h12V5a2 2 0 0 0-2-2z"/>
                                    </svg>
                                    <span v-if="isLockerHandoff">{{ t('listings.handoffLocker') }}</span>
                                    <span v-else>{{ t('listings.handoffInPerson') }}{{ listing.address ? ' — ' + listing.address : '' }}</span>
                                </div>
                            </div>

                            <div class="vendeur-card">
                                <div class="vendeur-avatar">{{ sellerInitial }}</div>
                                <div class="vendeur-info">
                                    <span class="vendeur-nom">{{ sellerName }}</span>
                                    <span class="vendeur-since">Niveau {{ sellerLevel }}</span>
                                </div>
                            </div>

                            <div class="detail-desc">
                                <h2 class="desc-title">{{ t('listingDetail.description') }}</h2>
                                <p class="desc-text">{{ listing.description || t('listingDetail.noDescription') }}</p>
                            </div>

                            <div class="detail-actions">
                                <button
                                    class="btn-contact"
                                    :class="{ 'btn-contact--active': listing && clientStore.isChattingWith(getItemId(listing)) }"
                                    @click="handleAction"
                                >
                                    {{ listing && clientStore.isChattingWith(getItemId(listing)) ? t('listingDetail.continueChat') : t('listingDetail.contactSeller') }}
                                </button>
                                <button v-if="!listing || !clientStore.isChattingWith(getItemId(listing))" class="btn-save">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18">
                                        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                                    </svg>
                                    {{ t('listingDetail.save') }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </template>
    </div>
</template>

<style scoped>
.loading-state, .error-state {
    padding: 80px 0;
    text-align: center;
    opacity: 0.6;
}
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

.breadcrumb-bar {
    padding: 20px 0;
    border-bottom: 1px solid rgba(53, 53, 53, 0.08);
}
.breadcrumb-bar .container {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
}
.breadcrumb-link {
    color: var(--green-mid);
    text-decoration: none;
    transition: color 0.2s;
}
.breadcrumb-link:hover {
    color: var(--green-dark);
}
.breadcrumb-sep {
    color: var(--charcoal);
    opacity: 0.4;
}
.breadcrumb-current {
    color: var(--charcoal);
    opacity: 0.7;
}

.detail-section {
    flex: 1;
    padding: 40px 0 80px;
}

.detail-layout {
    display: grid;
    grid-template-columns: 1fr 400px;
    gap: 48px;
    align-items: flex-start;
}

.main-img-wrap {
    border-radius: 12px;
    overflow: hidden;
    aspect-ratio: 4/3;
    margin-bottom: 12px;
}
.main-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.thumbs-row {
    display: flex;
    gap: 10px;
}
.thumb-wrap {
    width: 80px;
    height: 60px;
    border-radius: 6px;
    overflow: hidden;
    border: 2px solid transparent;
    cursor: pointer;
    transition: border-color 0.2s;
}
.thumb-wrap:hover {
    border-color: var(--green-mid);
}
.thumb-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}

.detail-badges {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
}
.badge {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 5px;
    font-size: 0.75rem;
    font-weight: 700;
}
.badge--don {
    background: var(--green-mid);
    color: var(--white);
}
.badge--vente {
    background: var(--green-mid);
    color: var(--white);
}
.badge--cat {
    background: var(--green-pale);
    color: var(--green-dark);
}

.detail-titre {
    font-size: clamp(1.6rem, 3vw, 2.2rem);
    font-weight: 800;
    color: var(--charcoal);
    line-height: 1.2;
    letter-spacing: -0.02em;
    margin: 0 0 16px;
}

.detail-prix {
    font-size: 2rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
    margin-bottom: 20px;
}
.detail-gratuit {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--green-mid);
    margin-bottom: 20px;
}

.detail-meta {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 28px;
}
.meta-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
    color: var(--charcoal);
    opacity: 0.75;
}
.meta-item svg {
    color: var(--green-mid);
    flex-shrink: 0;
}

.detail-desc {
    border-top: 1px solid rgba(53, 53, 53, 0.1);
    padding-top: 20px;
    margin-bottom: 24px;
}
.desc-title {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--charcoal);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin: 0 0 12px;
}
.desc-text {
    font-size: 0.9rem;
    color: var(--charcoal);
    opacity: 0.82;
    line-height: 1.7;
    margin: 0;
}

.vendeur-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: var(--green-pale);
    border-radius: 10px;
    margin-bottom: 20px;
}
.vendeur-avatar {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: var(--green-dark);
    color: var(--white);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 700;
    flex-shrink: 0;
}
.vendeur-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}
.vendeur-nom {
    font-size: 0.92rem;
    font-weight: 700;
    color: var(--charcoal);
}
.vendeur-since {
    font-size: 0.78rem;
    color: var(--charcoal);
    opacity: 0.6;
}

.detail-actions {
    display: flex;
    gap: 12px;
}
.btn-contact {
    flex: 1;
    display: block;
    padding: 14px 24px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 700;
    text-decoration: none;
    text-align: center;
    transition:
        background 0.2s,
        transform 0.15s;
}
.btn-contact:hover {
    background: var(--green-mid);
    transform: translateY(-1px);
}
.btn-contact--active {
    background: #4183d7;
}
.btn-contact--active:hover {
    background: #3569ad;
}
.btn-save {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 18px;
    background: transparent;
    color: var(--charcoal);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition:
        border-color 0.2s,
        color 0.2s;
    white-space: nowrap;
}
.btn-save:hover {
    border-color: var(--green-mid);
    color: var(--green-mid);
}

@media (max-width: 860px) {
    .detail-layout {
        grid-template-columns: 1fr;
    }
}
@media (max-width: 560px) {
    .detail-actions {
        flex-direction: column;
    }
    .btn-save {
        justify-content: center;
    }
}
</style>
