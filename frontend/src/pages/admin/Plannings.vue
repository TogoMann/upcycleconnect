<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Creneau {
    id: number
    titre: string
    type: string
    responsable: string
    date: string
    heure_debut: string
    heure_fin: string
    participants: number
}

const creneaux = ref<Creneau[]>([])
const filterType = ref('')

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/plannings', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) creneaux.value = await res.json()
    } catch {}
})

const types = ['formation', 'atelier', 'depot', 'collecte']

function typeClass(t: string) {
    const map: Record<string, string> = {
        formation: 'badge--formation',
        atelier: 'badge--atelier',
        depot: 'badge--depot',
        collecte: 'badge--collecte',
    }
    return 'badge ' + (map[t] ?? 'badge--autre')
}

const filtered = computed(() => creneaux.value.filter(c => !filterType.value || c.type === filterType.value))
</script>

<template>
    <div class="plannings">
        <div class="page-header">
            <h1 class="page-title">Plannings.</h1>
            <p class="page-subtitle">Vue globale de tous les créneaux planifiés.</p>
        </div>

        <div class="filter-row">
            <button
                v-for="t in ['', ...types]"
                :key="t"
                class="filter-btn"
                :class="{ 'filter-btn--active': filterType === t }"
                @click="filterType = t"
            >{{ t || 'Tous' }}</button>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Type</th>
                        <th>Responsable</th>
                        <th>Date</th>
                        <th>Horaire</th>
                        <th>Participants</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="filtered.length === 0">
                        <td colspan="6" class="empty">Aucun créneau.</td>
                    </tr>
                    <tr v-for="c in filtered" :key="c.id">
                        <td class="td-bold">{{ c.titre }}</td>
                        <td><span :class="typeClass(c.type)">{{ c.type }}</span></td>
                        <td class="td-muted">{{ c.responsable }}</td>
                        <td class="td-muted">{{ c.date }}</td>
                        <td>{{ c.heure_debut }} – {{ c.heure_fin }}</td>
                        <td>{{ c.participants }}</td>
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
.filter-row { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 20px; }
.filter-btn { padding: 7px 14px; border-radius: 8px; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; font-size: 0.82rem; font-weight: 600; color: var(--charcoal); cursor: pointer; transition: all 0.15s; text-transform: capitalize; }
.filter-btn:hover { border-color: var(--green-mid); color: var(--green-dark); }
.filter-btn--active { background: var(--green-dark); color: var(--white); border-color: var(--green-dark); }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; text-transform: capitalize; }
.badge--formation { background: var(--green-pale); color: var(--green-dark); }
.badge--atelier { background: #dbeafe; color: #1e40af; }
.badge--depot { background: #fef3c7; color: #92400e; }
.badge--collecte { background: #ede9fe; color: #5b21b6; }
.badge--autre { background: rgba(53,53,53,0.08); color: var(--charcoal); }
</style>
