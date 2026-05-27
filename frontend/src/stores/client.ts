import { defineStore, storeToRefs } from 'pinia'
import { computed } from 'vue'

import { useAuthStore } from './auth'
import { useChatStore } from './chat'
import { useCartStore } from './cart'
import { useCatalogueStore } from './catalogue'
import { useDepositStore } from './deposit'
import { useListingStore } from './listing'
import { useLocationStore } from './location'
import { usePlanningStore } from './planning'
import { useProfileStore } from './profile'

export const useClientStore = defineStore('client', () => {
    // Initialize all domain stores
    const authStore = useAuthStore()
    const chatStore = useChatStore()
    const cartStore = useCartStore()
    const catalogueStore = useCatalogueStore()
    const depositStore = useDepositStore()
    const listingStore = useListingStore()
    const locationStore = useLocationStore()
    const planningStore = usePlanningStore()
    const profileStore = useProfileStore()

    const { annonces, allAnnonces, conversations } = storeToRefs(listingStore)
    const { depots, sites, lockerAccesses } = storeToRefs(depositStore)
    const { entries, planning } = storeToRefs(planningStore)
    const { score, scoreHistory } = storeToRefs(profileStore)
    const { events, courses, participations, courseOrders } = storeToRefs(catalogueStore)
    const { cities } = storeToRefs(locationStore)
    const { cart } = storeToRefs(cartStore)

    // Aggregate Chat (Conversations used directly in client store previously)
    async function fetchConversations() {
        try {
            const data = await chatStore.getConversations()
            conversations.value = Array.isArray(data) ? data : []
        } catch (e) {
            console.error('Fetch Conversations Error:', e)
            conversations.value = []
        }
    }

    function isChattingWith(listingId: number): boolean {
        return conversations.value.some(c => {
            const cLid = c.listing_id && typeof c.listing_id === 'object' ? (c.listing_id as any).Int64 : c.listing_id
            return Number(cLid) === listingId
        })
    }

    return {
        // --- State ---
        annonces,
        allAnnonces,
        conversations,
        depots,
        sites,
        lockerAccesses,
        entries,
        planning,
        score,
        scoreHistory,
        events,
        courses,
        participations,
        courseOrders,
        cities,
        cart,
        
        // Global UI state
        // For simplicity, we bind isLoading to a computed that checks if ANY store is loading.
        isLoading: computed(() => 
            cartStore.isLoading || 
            catalogueStore.isLoading || 
            depositStore.isLoading || 
            listingStore.isLoading || 
            locationStore.isLoading || 
            planningStore.isLoading || 
            profileStore.isLoading
        ),
        
        // Same for error (returns the first non-null error)
        error: computed(() => 
            cartStore.error || 
            catalogueStore.error || 
            depositStore.error || 
            listingStore.error || 
            locationStore.error || 
            planningStore.error || 
            profileStore.error
        ),

        // --- Methods (Re-exported directly) ---

        // Listing
        fetchAnnonces: listingStore.fetchAnnonces,
        fetchAllAnnonces: listingStore.fetchAllAnnonces,
        uploadImage: listingStore.uploadImage,
        createAnnonce: listingStore.createAnnonce,
        createOrder: listingStore.createOrder,
        fetchConversations,
        isChattingWith,

        // Deposit
        fetchDepots: depositStore.fetchDepots,
        fetchSites: depositStore.fetchSites,
        fetchLockerAccesses: depositStore.fetchLockerAccesses,
        createItem: depositStore.createItem,

        // Profile
        fetchScore: profileStore.fetchScore,
        fetchScoreHistory: profileStore.fetchScoreHistory,
        markTutorialSeen: profileStore.markTutorialSeen,
        updateProfile: profileStore.updateProfile,

        // Planning
        fetchEntries: planningStore.fetchEntries,
        fetchPlanning: planningStore.fetchPlanning,
        createPersonalEvent: planningStore.createPersonalEvent,
        deletePersonalEvent: planningStore.deletePersonalEvent,
        createEntry: planningStore.createEntry,
        deleteEntry: planningStore.deleteEntry,

        // Catalogue
        fetchCatalogue: catalogueStore.fetchCatalogue,
        fetchParticipations: catalogueStore.fetchParticipations,
        fetchCourseOrders: catalogueStore.fetchCourseOrders,
        createCourseOrder: catalogueStore.createCourseOrder,
        createEventParticipation: catalogueStore.createEventParticipation,

        // Location
        fetchCities: locationStore.fetchCities,

        // Cart
        fetchCart: cartStore.fetchCart,
        addToCart: cartStore.addToCart,
        removeFromCart: cartStore.removeFromCart,
        checkoutCart: cartStore.checkoutCart,
    }
})