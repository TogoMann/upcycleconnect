<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Annonce {
    id: number
    titre: string
    categorie: string
    prix: number
    statut: string
    date: string
}

const annonces = ref<Annonce[]>([])
const search = ref('')
const filterStatut = ref('')

const filtered = computed(() =>
    annonces.value.filter(a => {
        const matchSearch = a.titre.toLowerCase().includes(search.value.toLowerCase())
        const matchStatut = !filterStatut.value || a.statut === filterStatut.value
        return matchSearch && matchStatut
    })
)

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro/annonces', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) annonces.value = await res.json()
    } catch {}
})

function badgeClass(s: string) {
    if (s === 'active') return 'badge badge--active'
    if (s === 'sold') return 'badge badge--sold'
    return 'badge badge--draft'
}
function badgeLabel(s: string) {
    if (s === 'active') return 'Active'
    if (s === 'sold') return 'Vendu'
    return 'Brouillon'
}
</script>

<template>
    <div class="annonces">
        <div class="page-header">
            <h1 class="page-title">Annonces.</h1>
            <p class="page-subtitle">Catalogue de vos annonces avec filtres avancés.</p>
        </div>

        <div class="filters-row">
            <input
                v-model="search"
                type="text"
                class="filter-input"
                placeholder="Rechercher une annonce…"
            />
            <select v-model="filterStatut" class="filter-select">
                <option value="">Tous les statuts</option>
                <option value="active">Active</option>
                <option value="sold">Vendu</option>
                <option value="draft">Brouillon</option>
            </select>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Catégorie</th>
                        <th>Prix</th>
                        <th>Date</th>
                        <th>Statut</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="filtered.length === 0">
                        <td colspan="5" class="empty">Aucune annonce trouvée.</td>
                    </tr>
                    <tr v-for="a in filtered" :key="a.id">
                        <td class="td-bold">{{ a.titre }}</td>
                        <td class="td-muted">{{ a.categorie }}</td>
                        <td>{{ a.prix.toFixed(2) }} €</td>
                        <td class="td-muted">{{ a.date }}</td>
                        <td><span :class="badgeClass(a.statut)">{{ badgeLabel(a.statut) }}</span></td>
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
</style>
