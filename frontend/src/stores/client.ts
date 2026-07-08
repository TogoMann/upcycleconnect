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
    const { score, scoreHistory, quests } = storeToRefs(profileStore)
    const { events, courses, participations, courseOrders } = storeToRefs(catalogueStore)
    const { cities } = storeToRefs(locationStore)
    const { cart } = storeToRefs(cartStore)

    
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
        quests,
        events,
        courses,
        participations,
        courseOrders,
        cities,
        cart,
        
        
        
        isLoading: computed(() => 
            cartStore.isLoading || 
            catalogueStore.isLoading || 
            depositStore.isLoading || 
            listingStore.isLoading || 
            locationStore.isLoading || 
            planningStore.isLoading || 
            profileStore.isLoading
        ),
        
        
        error: computed(() => 
            cartStore.error || 
            catalogueStore.error || 
            depositStore.error || 
            listingStore.error || 
            locationStore.error || 
            planningStore.error || 
            profileStore.error
        ),

        

        
        fetchAnnonces: listingStore.fetchAnnonces,
        fetchAllAnnonces: listingStore.fetchAllAnnonces,
        uploadImage: listingStore.uploadImage,
        fetchSitesWithLockers: listingStore.fetchSitesWithLockers,
        createAnnonce: listingStore.createAnnonce,
        createOrderCheckout: listingStore.createOrderCheckout,
        deleteAnnonce: listingStore.deleteAnnonce,
        fetchConversations,
        isChattingWith,

        
        fetchDepots: depositStore.fetchDepots,
        fetchSites: depositStore.fetchSites,
        fetchLockerAccesses: depositStore.fetchLockerAccesses,
        createItem: depositStore.createItem,

        
        fetchScore: profileStore.fetchScore,
        fetchScoreHistory: profileStore.fetchScoreHistory,
        fetchQuests: profileStore.fetchQuests,
        markTutorialSeen: profileStore.markTutorialSeen,
        updateProfile: profileStore.updateProfile,

        
        fetchEntries: planningStore.fetchEntries,
        fetchPlanning: planningStore.fetchPlanning,
        createPersonalEvent: planningStore.createPersonalEvent,
        deletePersonalEvent: planningStore.deletePersonalEvent,
        createEntry: planningStore.createEntry,
        deleteEntry: planningStore.deleteEntry,

        
        fetchCatalogue: catalogueStore.fetchCatalogue,
        fetchParticipations: catalogueStore.fetchParticipations,
        fetchCourseOrders: catalogueStore.fetchCourseOrders,
        createCourseOrder: catalogueStore.createCourseOrder,
        createEventParticipation: catalogueStore.createEventParticipation,

        
        fetchCities: locationStore.fetchCities,

        
        fetchCart: cartStore.fetchCart,
        addToCart: cartStore.addToCart,
        removeFromCart: cartStore.removeFromCart,
        clearCart: cartStore.clearCart,
        checkoutCart: cartStore.checkoutCart,
    }
})