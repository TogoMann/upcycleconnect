<script setup lang="ts">
import { API_BASE } from '@/config'
import { onMounted, computed } from 'vue'
import { useClientStore } from '@/stores/client'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const clientStore = useClientStore()
const router = useRouter()

onMounted(() => {
    clientStore.fetchCart()
})

function fmtTime(time: any): string {
    if (!time) return '—'
    if (typeof time === 'string') {
        return time.slice(0, 5)
    }
    if (typeof time === 'object') {
        if (time.Valid === false || time.valid === false) return '—'
        const micros = time.Microseconds ?? time.microseconds ?? 0
        const totalSeconds = Math.floor(micros / 1_000_000)
        const h = Math.floor(totalSeconds / 3600).toString().padStart(2, '0')
        const m = Math.floor((totalSeconds % 3600) / 60).toString().padStart(2, '0')
        return `${h}:${m}`
    }
    return '—'
}

function fmtDate(date: any): string {
    if (!date) return '—'
    if (typeof date === 'string') {
        return new Date(date).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
    }
    if (typeof date === 'object') {
        if (date.Valid === false || date.valid === false) return '—'
        const t = date.Time ?? date.time
        if (t) {
            return new Date(t).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
        }
    }
    return '—'
}

const total = computed(() => {
    return clientStore.cart.reduce((acc, item) => {
        const price = item.listing?.price || item.event?.price || item.course?.price || 0
        return acc + (parseFloat(price) || 0)
    }, 0)
})

function goToChat(id: any) {
    let listingId = ''
    if (typeof id === 'object' && id !== null) {
        listingId = id.Int64?.toString() || id.id?.toString() || ''
    } else {
        listingId = id?.toString() || ''
    }

    if (!listingId) return

    router.push({
        path: '/particulier/chat',
        query: { listingId }
    })
}

async function handleCheckout() {
    if (clientStore.cart.length === 0) return

    const hasListing = clientStore.cart.some((item: any) => item.listing_id != null)
    const hasDirectPay = clientStore.cart.some((item: any) => item.event_id != null || item.course_id != null)

    if (hasListing && !hasDirectPay) {
        alert(t('client.panier.listingsNeedAgreement'))
        return
    }

    if (hasListing && hasDirectPay) {
        if (!confirm(t('client.panier.mixedCartConfirm'))) {
            return
        }
    }

    const directPayTotal = clientStore.cart.reduce((acc: number, item: any) => {
        if (item.listing_id != null) return acc
        const price = item.event?.price || item.course?.price || 0
        return acc + (parseFloat(price) || 0)
    }, 0)

    router.push({
        path: '/particulier/paiement',
        query: {
            type: 'cart',
            price: directPayTotal.toFixed(2),
            name: hasListing ? t('client.panier.workshopsAndEvents') : t('client.panier.cart'),
        },
    })
}

async function handleRemove(item: any) {
    try {
        if (item.listing_id != null) {
            await clientStore.removeFromCart('listing', item.listing_id)
        } else if (item.event_id != null) {
            await clientStore.removeFromCart('event', item.event_id)
        } else if (item.course_id != null) {
            await clientStore.removeFromCart('course', item.course_id)
        }
    } catch (e: any) {
        alert(e.message)
    }
}

async function handleClearCart() {
    if (!confirm(t('client.panier.confirmClearCart'))) return
    try {
        await clientStore.clearCart()
    } catch (e: any) {
        alert(e.message)
    }
}
</script>

<template>
    <div class="panier">
        <div class="page-header">
            <div class="page-header-text">
                <h1 class="page-title">{{ t('client.panier.pageTitle') }}</h1>
                <p class="page-subtitle">{{ t('client.panier.subtitle') }}</p>
            </div>
            <button v-if="clientStore.cart.length > 0" class="btn-clear-cart" @click="handleClearCart">{{ t('client.panier.clearCart') }}</button>
        </div>

        <div v-if="clientStore.isLoading && clientStore.cart.length === 0" class="loading">
            {{ t('client.panier.loading') }}
        </div>

        <div v-else-if="clientStore.cart.length === 0" class="empty-cart">
            <div class="empty-icon">🛒</div>
            <h2>{{ t('client.panier.emptyTitle') }}</h2>
            <p>{{ t('client.panier.emptySubtitle') }}</p>
            <router-link to="/particulier/catalogue" class="btn-primary">{{ t('client.panier.viewCatalogue') }}</router-link>
        </div>

        <div v-else class="cart-content">
            <div class="items-list">
                <div v-for="item in clientStore.cart" :key="item.id" class="cart-item">
                    
                    <template v-if="item.listing">
                        <img v-if="item.listing.image_url?.String" :src="`${API_BASE}` + item.listing.image_url.String" :alt="item.listing.name" class="item-img">
                        <img v-else-if="item.listing.image_url && typeof item.listing.image_url === 'string'" :src="`${API_BASE}` + item.listing.image_url" :alt="item.listing.name" class="item-img">
                        <img v-else src="https://via.placeholder.com/100" :alt="item.listing.name" class="item-img">
                        <div class="item-details">
                            <h3 class="item-name">{{ item.listing.name }}</h3>
                            <p class="item-category">{{ t('client.panier.item') }} • {{ item.listing.category }}</p>
                            <p class="item-warning">{{ t('client.panier.sellerAgreementRequired') }}</p>
                        </div>
                        <div class="item-price">{{ item.listing.price }}€</div>
                        <div class="item-actions">
                            <button class="btn-chat" @click="goToChat(item.listing.id)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                                    <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
                                </svg>
                                {{ t('client.panier.message') }}
                            </button>
                            <button class="btn-remove" @click="handleRemove(item)">{{ t('client.panier.remove') }}</button>
                        </div>
                    </template>

                    
                    <template v-else-if="item.event">
                        <div class="item-icon">📅</div>
                        <div class="item-details">
                            <h3 class="item-name">{{ t('client.panier.event') }}</h3>
                            <p class="item-category">{{ t('client.panier.participation') }} • {{ item.event.location }}</p>
                            <p class="item-info">📅 {{ fmtDate(item.event.date) }} • {{ fmtTime(item.event.start_time) }} - {{ fmtTime(item.event.end_time) }}</p>
                        </div>
                        <div class="item-price">{{ item.event.price }}€</div>
                        <div class="item-actions">
                            <button class="btn-remove" @click="handleRemove(item)">{{ t('client.panier.remove') }}</button>
                        </div>
                    </template>

                    
                    <template v-else-if="item.course">
                        <div class="item-icon">🛠️</div>
                        <div class="item-details">
                            <h3 class="item-name">{{ item.course.name }}</h3>
                            <p class="item-category">{{ t('client.panier.training') }} • {{ t('client.panier.capacity') }}: {{ item.course.max_capacity?.Int32 }}</p>
                            <p class="item-info">📅 {{ fmtDate(item.course.date) }} • {{ fmtTime(item.course.start_time) }} - {{ fmtTime(item.course.end_time) }}</p>
                        </div>
                        <div class="item-price">{{ item.course.price }}€</div>
                        <div class="item-actions">
                            <button class="btn-remove" @click="handleRemove(item)">{{ t('client.panier.remove') }}</button>
                        </div>
                    </template>
                </div>
            </div>

            <div class="cart-summary">
                <div class="summary-card">
                    <h3>{{ t('client.panier.summary') }}</h3>
                    <div class="summary-line">
                        <span>{{ t('client.panier.items', { count: clientStore.cart.length }) }}</span>
                        <span>{{ total.toFixed(2) }}€</span>
                    </div>
                    <div class="summary-line total">
                        <span>{{ t('client.panier.estimatedTotal') }}</span>
                        <span>{{ total.toFixed(2) }}€</span>
                    </div>
                    <p class="summary-note">{{ t('client.panier.summaryNote') }}</p>
                    <button class="btn-checkout" @click="handleCheckout">{{ t('client.panier.payNow') }}</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.item-icon { font-size: 2.5rem; width: 80px; height: 80px; display: flex; align-items: center; justify-content: center; background: #f0f0f0; border-radius: 12px; }
.item-info { font-size: 0.85rem; opacity: 0.7; margin: 4px 0 0; }

.page-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; margin-bottom: 32px; }
.page-title { font-size: 2.6rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; }
.page-subtitle { font-size: 1rem; color: var(--charcoal); opacity: 0.6; }
.btn-clear-cart { flex-shrink: 0; padding: 10px 18px; border: 1px solid #d9534f; border-radius: 8px; background: transparent; color: #d9534f; font-weight: 600; cursor: pointer; transition: background 0.15s, color 0.15s; }
.btn-clear-cart:hover { background: #d9534f; color: #fff; }

.empty-cart { text-align: center; padding: 60px 20px; background: var(--white); border-radius: 20px; border: 1.5px solid rgba(53,53,53,0.08); }
.empty-icon { font-size: 4rem; margin-bottom: 20px; }
.empty-cart h2 { margin-bottom: 10px; color: var(--charcoal); }
.empty-cart p { margin-bottom: 30px; opacity: 0.6; }

.cart-content { display: grid; grid-template-columns: 1fr 350px; gap: 30px; align-items: start; }

.items-list { display: flex; flex-direction: column; gap: 16px; }
.cart-item { background: var(--white); border-radius: 16px; padding: 16px; display: flex; align-items: center; gap: 20px; border: 1.5px solid rgba(53,53,53,0.05); }
.item-img { width: 80px; height: 80px; border-radius: 12px; object-fit: cover; background: #f0f0f0; }
.item-details { flex: 1; }
.item-name { font-size: 1.1rem; font-weight: 700; margin: 0 0 4px; }
.item-category { font-size: 0.85rem; opacity: 0.5; margin: 0; }
.item-warning { font-size: 0.8rem; color: #d97706; font-weight: 600; margin: 8px 0 0; }
.item-price { font-weight: 700; font-size: 1.1rem; color: var(--green-dark); }
.item-actions { display: flex; flex-direction: column; gap: 8px; align-items: flex-end; }
.btn-chat { background: #4183d7; border: none; color: white; font-size: 0.85rem; font-weight: 700; cursor: pointer; padding: 10px 16px; border-radius: 12px; transition: background 0.2s, transform 0.1s; display: flex; align-items: center; gap: 8px; white-space: nowrap; }
.btn-chat:hover { background: #3569ad; transform: translateY(-1px); }
.btn-chat:active { transform: translateY(0); }
.btn-chat svg { flex-shrink: 0; }
.btn-remove { background: none; border: none; color: #dc2626; font-size: 0.85rem; font-weight: 600; cursor: pointer; padding: 8px; }

.summary-card { background: var(--white); border-radius: 20px; padding: 24px; border: 1.5px solid rgba(53,53,53,0.08); position: sticky; top: 20px; }
.summary-card h3 { margin: 0 0 20px; font-size: 1.3rem; }
.summary-line { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 0.95rem; }
.summary-line.total { border-top: 1px solid rgba(53,53,53,0.1); padding-top: 12px; margin-top: 12px; font-weight: 800; font-size: 1.2rem; }
.summary-note { font-size: 0.8rem; color: #718096; margin-top: 12px; line-height: 1.4; }
.btn-checkout { width: 100%; margin-top: 16px; padding: 16px; border-radius: 12px; border: none; background: var(--green-dark); color: var(--white); font-weight: 700; font-size: 1rem; cursor: pointer; transition: transform 0.2s; }
.btn-checkout:hover { transform: translateY(-2px); }

@media (max-width: 900px) {
    .cart-content { grid-template-columns: 1fr; }
}

.btn-primary { display: inline-block; background: var(--green-dark); color: var(--white); padding: 12px 24px; border-radius: 10px; text-decoration: none; font-weight: 600; }
</style>
