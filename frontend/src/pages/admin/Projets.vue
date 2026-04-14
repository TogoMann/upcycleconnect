<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Projet {
    id: number
    nom: string
    type: string
    auteur: string
    statut: string
    date: string
    mis_en_avant: boolean
}

const projets = ref<Projet[]>([])

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/projets', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) projets.value = await res.json()
    } catch {}
})

async function toggleMisEnAvant(p: Projet) {
    await fetch(`http://localhost:8081/admin/projets/${p.id}/mise-en-avant`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
        body: JSON.stringify({ mis_en_avant: !p.mis_en_avant }),
    })
    p.mis_en_avant = !p.mis_en_avant
}
</script>

<template>
    <div class="projets">
        <div class="page-header">
            <h1 class="page-title">Projets.</h1>
            <p class="page-subtitle">Supervision des projets et mise en avant.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Nom</th>
                        <th>Type</th>
                        <th>Auteur</th>
                        <th>Date</th>
                        <th>Statut</th>
                        <th>Mise en avant</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="projets.length === 0">
                        <td colspan="6" class="empty">Aucun projet.</td>
                    </tr>
                    <tr v-for="p in projets" :key="p.id">
                        <td class="td-bold">{{ p.nom }}</td>
                        <td class="td-muted">{{ p.type }}</td>
                        <td class="td-muted">{{ p.auteur }}</td>
                        <td class="td-muted">{{ p.date }}</td>
                        <td>
                            <span class="badge" :class="p.statut === 'en_cours' ? 'badge--active' : 'badge--done'">
                                {{ p.statut === 'en_cours' ? 'En cours' : 'Terminé' }}
                            </span>
                        </td>
                        <td>
                            <button
                                class="toggle-btn"
                                :class="{ 'toggle-btn--on': p.mis_en_avant }"
                                @click="toggleMisEnAvant(p)"
                                :title="p.mis_en_avant ? 'Retirer la mise en avant' : 'Mettre en avant'"
                            >
                                <span class="toggle-track">
                                    <span class="toggle-thumb"></span>
                                </span>
                            </button>
                        </td>
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
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--done { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.toggle-btn { background: none; border: none; cursor: pointer; padding: 4px; }
.toggle-track { display: flex; width: 40px; height: 22px; background: rgba(53,53,53,0.15); border-radius: 11px; padding: 2px; transition: background 0.2s; }
.toggle-thumb { width: 18px; height: 18px; background: var(--white); border-radius: 50%; box-shadow: 0 1px 3px rgba(0,0,0,0.2); transition: transform 0.2s; }
.toggle-btn--on .toggle-track { background: var(--green-dark); }
.toggle-btn--on .toggle-thumb { transform: translateX(18px); }
</style>
