<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { API_BASE } from '@/config'

const route = useRoute()
const authStore = useAuthStore()
const maintenanceMode = ref(false)
const siteName = ref('upcycleconnect')

const showGlobalWarning = ref(false)
const nextPaymentDateStr = ref('')
const nextPaymentDays = ref(0)
const paymentMethodMissing = ref(false)

const pmRedirectPath = computed(() => {
    return authStore.userRole === 'pro' ? '/pro/moyens-de-paiement' : '/particulier/moyens-de-paiement'
})

async function checkSubscriptionWarning() {
    if (!authStore.token || authStore.userRole === 'admin') {
        showGlobalWarning.value = false
        return
    }
    try {
        const token = authStore.token
        const headers = { Authorization: `Bearer ${token}` }
        
        const [subRes, pmRes] = await Promise.all([
            fetch(`${API_BASE}/subscriptions/me`, { headers }).catch(() => null),
            fetch(`${API_BASE}/payment-methods/check`, { headers }).catch(() => null),
        ])
        
        let hasSub = false
        let untilTime = ''
        let tier = 'Free'
        if (subRes && subRes.ok) {
            const data = await subRes.json()
            if (data.until) {
                hasSub = true
                if (typeof data.until === 'string') {
                    untilTime = data.until
                } else if (typeof data.until === 'object' && data.until.Valid) {
                    untilTime = data.until.Time
                } else {
                    hasSub = false
                }
                tier = data.tier || 'Free'
            }
        }
        
        let hasPm = true
        if (pmRes && pmRes.ok) {
            const data = await pmRes.json()
            hasPm = data.has_payment_method !== false
        }
        
        if (hasSub && tier !== 'Free') {
            const until = new Date(untilTime)
            const now = new Date()
            const diff = Math.ceil((until.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
            
            if (diff <= 7 && !hasPm) {
                showGlobalWarning.value = true
                nextPaymentDays.value = Math.max(0, diff)
                paymentMethodMissing.value = true
                nextPaymentDateStr.value = until.toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
            } else {
                showGlobalWarning.value = false
            }
        } else {
            showGlobalWarning.value = false
        }
    } catch {
        showGlobalWarning.value = false
    }
}

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/parametres/public`)
        if (res.ok) {
            const data = await res.json()
            maintenanceMode.value = !!data.maintenance
            siteName.value = data.nom_site || 'upcycleconnect'
        }
    } catch (e) {
        console.error('Failed to load settings', e)
    }
})

watch(() => authStore.user, () => {
    checkSubscriptionWarning()
}, { immediate: true })

watch(() => route.path, () => {
    checkSubscriptionWarning()
})

const isMaintenanceActive = computed(() => {
    if (!maintenanceMode.value) return false
    
    if (authStore.userRole === 'admin') return false
    
    const path = route.path
    if (path.startsWith('/admin') || path === '/auth/login' || path === '/connexion') return false
    return true
})
</script>

<template>
    <div v-if="isMaintenanceActive" class="maintenance-container">
        <div class="maintenance-card">
            <div class="maintenance-logo">{{ siteName }}</div>
            <div class="maintenance-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z" />
                </svg>
            </div>
            <h1>Maintenance en cours</h1>
            <p>Notre plateforme est actuellement en cours de maintenance pour amélioration. Nous serons de retour très bientôt !</p>
        </div>
    </div>
    <div v-else class="app-container">
        <div v-if="showGlobalWarning" class="global-warning-banner">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
            <span>
                Attention : Votre moyen de paiement n'est pas à jour. Votre abonnement expire dans {{ nextPaymentDays }} jour(s) (le {{ nextPaymentDateStr }}). Sans renouvellement de votre carte, votre accès prendra fin le 1er du mois prochain.
            </span>
            <router-link :to="pmRedirectPath" class="warning-banner-btn">Mettre à jour</router-link>
        </div>
        <RouterView />
    </div>
</template>

<style>
.maintenance-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: linear-gradient(135deg, #f7f9f6 0%, #e2ebd9 100%);
    font-family: 'Inter', sans-serif;
    padding: 20px;
}
.maintenance-card {
    background: #ffffff;
    border-radius: 20px;
    box-shadow: 0 10px 30px rgba(0,0,0,0.05);
    padding: 40px;
    text-align: center;
    max-width: 500px;
    border: 1px solid rgba(8, 106, 53, 0.08);
}
.maintenance-logo {
    font-weight: 800;
    font-size: 1.8rem;
    color: #086a35;
    margin-bottom: 24px;
    text-transform: lowercase;
}
.maintenance-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto 24px;
    background: rgba(8, 106, 53, 0.05);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #086a35;
}
.maintenance-icon svg {
    width: 40px;
    height: 40px;
}
.maintenance-card h1 {
    font-size: 1.8rem;
    font-weight: 700;
    color: #353535;
    margin: 0 0 12px;
}
.maintenance-card p {
    font-size: 1rem;
    color: #353535;
    opacity: 0.7;
    line-height: 1.6;
    margin: 0;
}

.global-warning-banner {
    background: #e53e3e;
    color: #ffffff;
    padding: 10px 20px;
    font-size: 0.85rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    flex-wrap: wrap;
    text-align: center;
    box-shadow: 0 4px 12px rgba(229, 62, 62, 0.25);
    position: relative;
    z-index: 1001;
}
.global-warning-banner svg {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
}
.warning-banner-btn {
    background: #ffffff;
    color: #e53e3e;
    padding: 4px 12px;
    border-radius: 4px;
    text-decoration: none;
    font-size: 0.8rem;
    font-weight: 700;
    transition: background 0.2s;
}
.warning-banner-btn:hover {
    background: #f7f9f6;
}
.app-container {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}
</style>
