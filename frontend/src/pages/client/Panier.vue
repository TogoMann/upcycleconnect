<script setup lang="ts">
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

async function handleCheckout() {
    if (clientStore.cart.length === 0) return
    try {
        await clientStore.checkoutCart()
        router.push({ name: 'ConfirmationPaiement' })
    } catch (e: any) {
        alert(e.message)
    }
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
                    <img v-if="item.listing.image_url?.String" :src="'http://localhost:8081' + item.listing.image_url.String" :alt="item.listing.name" class="item-img">
                    <img v-else-if="item.listing.image_url && typeof item.listing.image_url === 'string'" :src="'http://localhost:8081' + item.listing.image_url" :alt="item.listing.name" class="item-img">
                    <img v-else src="https://via.placeholder.com/100" :alt="item.listing.name" class="item-img">
                    <div class="item-details">
                        <h3 class="item-name">{{ item.listing.name }}</h3>
                        <p class="item-category">{{ item.listing.category }}</p>
                    </div>
                    <div class="item-price">{{ item.listing.price }}€</div>
                    <button class="btn-remove" @click="handleRemove(item.listing.id)">Supprimer</button>
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
                        <span>Total</span>
                        <span>{{ total.toFixed(2) }}€</span>
                    </div>
                    <button class="btn-checkout" @click="handleCheckout">Passer la commande</button>
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
.item-price { font-weight: 700; font-size: 1.1rem; color: var(--green-dark); }
.btn-remove { background: none; border: none; color: #dc2626; font-size: 0.85rem; font-weight: 600; cursor: pointer; padding: 10px; }

.summary-card { background: var(--white); border-radius: 20px; padding: 24px; border: 1.5px solid rgba(53,53,53,0.08); position: sticky; top: 20px; }
.summary-card h3 { margin: 0 0 20px; font-size: 1.3rem; }
.summary-line { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 0.95rem; }
.summary-line.total { border-top: 1px solid rgba(53,53,53,0.1); padding-top: 12px; margin-top: 12px; font-weight: 800; font-size: 1.2rem; }
.btn-checkout { width: 100%; margin-top: 24px; padding: 16px; border-radius: 12px; border: none; background: var(--green-dark); color: var(--white); font-weight: 700; font-size: 1rem; cursor: pointer; transition: transform 0.2s; }
.btn-checkout:hover { transform: translateY(-2px); }

@media (max-width: 900px) {
    .cart-content { grid-template-columns: 1fr; }
}

.btn-primary { display: inline-block; background: var(--green-dark); color: var(--white); padding: 12px 24px; border-radius: 10px; text-decoration: none; font-weight: 600; }
</style>
