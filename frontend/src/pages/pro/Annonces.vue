<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClientStore } from '@/stores/client'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const clientStore = useClientStore()

const CATEGORIES = ['Mobilier', 'Décoration', 'Vêtements', 'Jouet', 'Electronique', 'Outils']

const search = ref('')
const filterCategorie = ref('')
const buyingId = ref<number | null>(null)
const errorMsg = ref('')

onMounted(async () => {
    await clientStore.fetchAllAnnonces()
})

const annonces = computed(() => {
    return clientStore.allAnnonces.map((a: any) => {
        const p = a.price
        const priceVal = p ? (typeof p === 'object' ? (p.Float64 ?? p.Int64) : Number(p)) : 0
        return {
            id: a.id?.Int64 ?? a.id,
            titre: a.name,
            categorie: a.category || t('pro.annonces.uncategorized'),
            prix: Number(priceVal) || 0,
            statut: a.status,
            date: new Date(a.created_at?.Time ?? a.created_at).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'short', year: 'numeric' }),
            handoffMode: a.handoff_mode === 'casier' ? t('pro.annonces.handoffLocker') : t('pro.annonces.handoffInPerson'),
        }
    })
})

const filtered = computed(() =>
    annonces.value.filter(a => {
        const matchSearch = a.titre.toLowerCase().includes(search.value.toLowerCase())
        const matchCategorie = !filterCategorie.value || a.categorie === filterCategorie.value
        return matchSearch && matchCategorie
    })
)

async function handleBuy(a: { id: number; prix: number; statut: string }) {
    if (a.statut !== 'active') return
    buyingId.value = a.id
    errorMsg.value = ''
    try {
        const data = await clientStore.createOrderCheckout(a.id)
        if (data.free) {
            await clientStore.fetchAllAnnonces()
        } else if (data.url) {
            window.location.href = data.url
        }
    } catch (e: any) {
        errorMsg.value = e.message || t('pro.annonces.buyError')
    } finally {
        buyingId.value = null
    }
}

function badgeClass(s: string) {
    if (s === 'active') return 'badge badge--active'
    if (s === 'sold') return 'badge badge--sold'
    return 'badge badge--draft'
}
function badgeLabel(s: string) {
    if (s === 'active') return t('pro.annonces.statusActive')
    if (s === 'sold') return t('pro.annonces.statusSold')
    return t('pro.annonces.statusCancelled')
}
</script>

<template>
    <div class="annonces">
        <div class="page-header">
            <h1 class="page-title">{{ t('pro.annonces.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('pro.annonces.subtitle') }}</p>
        </div>

        <div v-if="errorMsg" class="error-banner">{{ errorMsg }}</div>

        <div class="filters-row">
            <input
                v-model="search"
                type="text"
                class="filter-input"
                :placeholder="t('pro.annonces.searchPlaceholder')"
            />
            <select v-model="filterCategorie" class="filter-select">
                <option value="">{{ t('pro.annonces.allCategories') }}</option>
                <option v-for="cat in CATEGORIES" :key="cat" :value="cat">{{ cat }}</option>
            </select>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('pro.annonces.title') }}</th>
                        <th>{{ t('pro.annonces.category') }}</th>
                        <th>{{ t('pro.annonces.transaction') }}</th>
                        <th>{{ t('pro.annonces.price') }}</th>
                        <th>{{ t('pro.annonces.date') }}</th>
                        <th>{{ t('pro.annonces.status') }}</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="filtered.length === 0">
                        <td colspan="7" class="empty">{{ t('pro.annonces.noResults') }}</td>
                    </tr>
                    <tr v-for="a in filtered" :key="a.id">
                        <td class="td-bold">{{ a.titre }}</td>
                        <td class="td-muted">{{ a.categorie }}</td>
                        <td class="td-muted">{{ a.handoffMode }}</td>
                        <td>{{ a.prix === 0 ? t('pro.annonces.don') : a.prix.toFixed(2) + ' €' }}</td>
                        <td class="td-muted">{{ a.date }}</td>
                        <td><span :class="badgeClass(a.statut)">{{ badgeLabel(a.statut) }}</span></td>
                        <td>
                            <button
                                v-if="a.statut === 'active'"
                                class="btn-buy"
                                :disabled="buyingId === a.id"
                                @click="handleBuy(a)"
                            >
                                {{ buyingId === a.id ? t('pro.annonces.buying') : (a.prix === 0 ? t('pro.annonces.retrieve') : t('pro.annonces.buy')) }}
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.error-banner { background: rgba(229, 62, 62, 0.08); border: 1px solid rgba(229, 62, 62, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: #e53e3e; margin-bottom: 16px; }
.filters-row { display: flex; gap: 12px; margin-bottom: 20px; }
.filter-input { flex: 1; padding: 10px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.filter-input:focus { border-color: var(--green-mid); }
.filter-select { padding: 10px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; cursor: pointer; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--sold { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.badge--draft { background: #fef3c7; color: #92400e; }
.btn-buy { padding: 7px 16px; background: var(--green-dark); color: var(--white); border: none; border-radius: 6px; font-size: 0.82rem; font-weight: 600; cursor: pointer; font-family: inherit; transition: background 0.2s; }
.btn-buy:hover:not(:disabled) { background: var(--green-mid); }
.btn-buy:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
