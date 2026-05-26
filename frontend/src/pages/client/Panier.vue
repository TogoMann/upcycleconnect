<script setup lang="ts">
import { API_BASE } from '@/config'
import { onMounted, computed } from 'vue'
import { useClientStore } from '@/stores/client'
import { useRouter } from 'vue-router'

const clientStore = useClientStore()
const router = useRouter()

onMounted(() => {
    clientStore.fetchCart()
})

const total = computed(() => {
    return clientStore.cart.reduce((acc, item) => acc + (parseFloat(item.listing.price) || 0), 0)
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
    alert("Les annonces nécessitent un accord préalable avec le vendeur via la messagerie. Veuillez finaliser le paiement depuis vos conversations.")
}

async function handleRemove(listingId: number) {
    try {
        await clientStore.removeFromCart(listingId)
    } catch (e: any) {
        alert(e.message)
    }
}
</script>

<template>
    <div class="panier">
        <div class="page-header">
            <h1 class="page-title">Mon Panier.</h1>
            <p class="page-subtitle">Centralisez vos trouvailles avant de valider.</p>
        </div>

        <div v-if="clientStore.isLoading && clientStore.cart.length === 0" class="loading">
            Chargement du panier...
        </div>

        <div v-else-if="clientStore.cart.length === 0" class="empty-cart">
            <div class="empty-icon">🛒</div>
            <h2>Votre panier est vide</h2>
            <p>Parcourez le catalogue pour trouver des objets uniques.</p>
            <router-link to="/client/catalogue" class="btn-primary">Voir le catalogue</router-link>
        </div>

        <div v-else class="cart-content">
            <div class="items-list">
                <div v-for="item in clientStore.cart" :key="item.id" class="cart-item">
                    <img v-if="item.listing.image_url?.String" :src="`${API_BASE}` + item.listing.image_url.String" :alt="item.listing.name" class="item-img">
                    <img v-else-if="item.listing.image_url && typeof item.listing.image_url === 'string'" :src="`${API_BASE}` + item.listing.image_url" :alt="item.listing.name" class="item-img">
                    <img v-else src="https://via.placeholder.com/100" :alt="item.listing.name" class="item-img">
                    <div class="item-details">
                        <h3 class="item-name">{{ item.listing.name }}</h3>
                        <p class="item-category">{{ item.listing.category }}</p>
                        <p class="item-warning">⚠️ Accord avec le vendeur requis</p>
                    </div>
                    <div class="item-price">{{ item.listing.price }}€</div>
                    <div class="item-actions">
                        <button class="btn-chat" @click="goToChat(item.listing.id)">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                                <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
                            </svg>
                            Envoyer un message
                        </button>
                        <button class="btn-remove" @click="handleRemove(item.listing.id)">Supprimer</button>
                    </div>
                </div>
            </div>

            <div class="cart-summary">
                <div class="summary-card">
                    <h3>Récapitulatif</h3>
                    <div class="summary-line">
                        <span>Articles ({{ clientStore.cart.length }})</span>
                        <span>{{ total.toFixed(2) }}€</span>
                    </div>
                    <div class="summary-line total">
                        <span>Total estimé</span>
                        <span>{{ total.toFixed(2) }}€</span>
                    </div>
                    <p class="summary-note">Les prix finaux sont fixés après dialogue avec les vendeurs.</p>
                    <button class="btn-checkout" @click="handleCheckout">Comment payer ?</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: 2.6rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; }
.page-subtitle { font-size: 1rem; color: var(--charcoal); opacity: 0.6; }

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
