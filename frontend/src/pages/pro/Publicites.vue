<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Pub {
    id: number
    titre: string
    type: string
    debut: string
    fin: string
    statut: string
    impressions: number
}

const pubs = ref<Pub[]>([])

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro/publicites', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) pubs.value = await res.json()
    } catch {}
})

function statusClass(s: string) {
    return s === 'active' ? 'badge badge--active' : s === 'expired' ? 'badge badge--expired' : 'badge badge--pending'
}

function statusLabel(s: string) {
    return s === 'active' ? 'Active' : s === 'expired' ? 'Expirée' : 'En attente'
}
</script>

<template>
    <div class="publicites">
        <div class="page-header">
            <h1 class="page-title">Publicités.</h1>
            <p class="page-subtitle">Suivez la performance de vos annonces publicitaires.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Type</th>
                        <th>Période</th>
                        <th>Impressions</th>
                        <th>Statut</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="pubs.length === 0">
                        <td colspan="5" class="empty">Aucune publicité.</td>
                    </tr>
                    <tr v-for="pub in pubs" :key="pub.id">
                        <td class="td-bold">{{ pub.titre }}</td>
                        <td>{{ pub.type }}</td>
                        <td class="td-muted">{{ pub.debut }} → {{ pub.fin }}</td>
                        <td>{{ pub.impressions.toLocaleString('fr-FR') }}</td>
                        <td><span :class="statusClass(pub.statut)">{{ statusLabel(pub.statut) }}</span></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53, 53, 53, 0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53, 53, 53, 0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53, 53, 53, 0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215, 236, 225, 0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--expired { background: rgba(53, 53, 53, 0.08); color: var(--charcoal); opacity: 0.6; }
.badge--pending { background: #fef3c7; color: #92400e; }
</style>
