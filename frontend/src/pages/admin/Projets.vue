<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
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
const filter = ref('all')
const loading = ref(true)
const selectedProjet = ref<Projet | null>(null)

const filteredProjets = computed(() => {
    if (filter.value === 'all') return projets.value
    if (filter.value === 'mis_en_avant') return projets.value.filter(p => p.mis_en_avant)
    return projets.value.filter(p => p.statut === filter.value)
})

const stats = computed(() => ({
    total: projets.value.length,
    enCours: projets.value.filter(p => p.statut === 'en_cours').length,
    termines: projets.value.filter(p => p.statut === 'termine').length,
    featured: projets.value.filter(p => p.mis_en_avant).length,
}))

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/admin/projets`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) projets.value = await res.json()
    } catch {}
    loading.value = false
})

async function toggleMisEnAvant(p: Projet) {
    try {
        const res = await fetch(`${API_BASE}/admin/projets/${p.id}/mise-en-avant`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ mis_en_avant: !p.mis_en_avant }),
        })
        if (res.ok) p.mis_en_avant = !p.mis_en_avant
    } catch {}
}

function statusConfig(s: string) {
    const map: Record<string, { label: string; class: string }> = {
        'en_cours': { label: 'En cours', class: 'badge--progress' },
        'termine': { label: 'Terminé', class: 'badge--done' },
    }
    return map[s] || { label: s, class: 'badge--default' }
}
</script>

<template>
    <div class="projets">
        <div class="page-header">
            <h1 class="page-title">Projets.</h1>
            <p class="page-subtitle">Supervision des projets d'upcycling et mise en avant.</p>
        </div>

        <div v-if="loading" class="loading-state">Chargement...</div>

        <template v-else>
            <!-- KPIs -->
            <div class="kpi-row">
                <div class="kpi-sm">
                    <div class="kpi-sm-value">{{ stats.total }}</div>
                    <div class="kpi-sm-label">Total</div>
                </div>
                <div class="kpi-sm kpi-sm--blue">
                    <div class="kpi-sm-value">{{ stats.enCours }}</div>
                    <div class="kpi-sm-label">En cours</div>
                </div>
                <div class="kpi-sm kpi-sm--green">
                    <div class="kpi-sm-value">{{ stats.termines }}</div>
                    <div class="kpi-sm-label">Terminés</div>
                </div>
                <div class="kpi-sm kpi-sm--yellow">
                    <div class="kpi-sm-value">{{ stats.featured }}</div>
                    <div class="kpi-sm-label">Mis en avant</div>
                </div>
            </div>

            <!-- Filtres -->
            <div class="filter-row">
                <button v-for="f in ['all', 'en_cours', 'termine', 'mis_en_avant']" :key="f"
                    class="filter-btn" :class="{ 'filter-btn--active': filter === f }"
                    @click="filter = f">
                    {{ f === 'all' ? 'Tous' : f === 'mis_en_avant' ? 'Mis en avant' : statusConfig(f).label }}
                </button>
            </div>

            <div class="table-wrap">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>Nom</th>
                            <th>Auteur</th>
                            <th>Date</th>
                            <th>Statut</th>
                            <th>Mise en avant</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="filteredProjets.length === 0">
                            <td colspan="6" class="empty">Aucun projet.</td>
                        </tr>
                        <tr v-for="p in filteredProjets" :key="p.id">
                            <td class="td-bold">{{ p.nom }}</td>
                            <td class="td-muted">{{ p.auteur }}</td>
                            <td class="td-muted">{{ p.date }}</td>
                            <td>
                                <span class="badge" :class="statusConfig(p.statut).class">
                                    {{ statusConfig(p.statut).label }}
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
                            <td>
                                <button class="btn-detail" @click="selectedProjet = p" title="Voir détails">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                                    Détails
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </template>

        <!-- Modal détail -->
        <Teleport to="body">
            <div v-if="selectedProjet" class="modal-overlay" @click.self="selectedProjet = null">
                <div class="modal-card">
                    <div class="modal-header">
                        <h3 class="modal-title">{{ selectedProjet.nom }}</h3>
                        <button class="modal-close" @click="selectedProjet = null">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>
                    <div class="detail-grid">
                        <div class="detail-item">
                            <span class="detail-label">Auteur</span>
                            <span class="detail-value">{{ selectedProjet.auteur }}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">Date de création</span>
                            <span class="detail-value">{{ selectedProjet.date }}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">Statut</span>
                            <span class="badge" :class="statusConfig(selectedProjet.statut).class">{{ statusConfig(selectedProjet.statut).label }}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">Mise en avant</span>
                            <span :class="selectedProjet.mis_en_avant ? 'detail-yes' : 'detail-no'">{{ selectedProjet.mis_en_avant ? 'Oui' : 'Non' }}</span>
                        </div>
                    </div>
                    <div class="modal-actions">
                        <button class="btn-toggle" :class="selectedProjet.mis_en_avant ? 'btn-toggle--off' : 'btn-toggle--on'"
                            @click="toggleMisEnAvant(selectedProjet)">
                            {{ selectedProjet.mis_en_avant ? 'Retirer la mise en avant' : 'Mettre en avant' }}
                        </button>
                        <button class="btn-secondary" @click="selectedProjet = null">Fermer</button>
                    </div>
                </div>
            </div>
        </Teleport>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; }

.kpi-row { display: flex; gap: 14px; margin-bottom: 20px; }
.kpi-sm { background: var(--white); border: 1.5px solid rgba(53,53,53,0.08); border-radius: 12px; padding: 16px 24px; }
.kpi-sm--green { background: var(--green-pale); border-color: transparent; }
.kpi-sm--blue { background: #eff6ff; border-color: transparent; }
.kpi-sm--yellow { background: #fffbeb; border-color: transparent; }
.kpi-sm-value { font-size: 1.8rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; line-height: 1; }
.kpi-sm--green .kpi-sm-value { color: var(--green-dark); }
.kpi-sm--blue .kpi-sm-value { color: #1e40af; }
.kpi-sm--yellow .kpi-sm-value { color: #92400e; }
.kpi-sm-label { font-size: 0.78rem; font-weight: 500; color: var(--charcoal); opacity: 0.5; margin-top: 4px; }

.filter-row { display: flex; gap: 8px; margin-bottom: 20px; }
.filter-btn { padding: 7px 16px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 20px; background: transparent; color: var(--charcoal); font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.filter-btn:hover { border-color: var(--green-mid); }
.filter-btn--active { background: var(--green-dark); color: var(--white); border-color: var(--green-dark); }

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
.badge--progress { background: #eff6ff; color: #1e40af; }
.badge--done { background: var(--green-pale); color: var(--green-dark); }
.badge--default { background: rgba(53,53,53,0.08); color: var(--charcoal); }

.toggle-btn { background: none; border: none; cursor: pointer; padding: 4px; }
.toggle-track { display: flex; width: 40px; height: 22px; background: rgba(53,53,53,0.15); border-radius: 11px; padding: 2px; transition: background 0.2s; }
.toggle-thumb { width: 18px; height: 18px; background: var(--white); border-radius: 50%; box-shadow: 0 1px 3px rgba(0,0,0,0.2); transition: transform 0.2s; }
.toggle-btn--on .toggle-track { background: var(--green-dark); }
.toggle-btn--on .toggle-thumb { transform: translateX(18px); }

.btn-detail { display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 6px; background: transparent; color: var(--charcoal); font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.btn-detail:hover { border-color: var(--green-mid); color: var(--green-dark); }
.btn-detail svg { width: 14px; height: 14px; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.modal-card { background: var(--white); border-radius: 20px; padding: 32px; max-width: 480px; width: 90%; box-shadow: 0 20px 60px rgba(0,0,0,0.15); }
.modal-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px; }
.modal-title { font-size: 1.15rem; font-weight: 800; color: var(--charcoal); margin: 0; }
.modal-close { background: none; border: none; cursor: pointer; color: var(--charcoal); opacity: 0.4; padding: 4px; }
.modal-close:hover { opacity: 1; }
.modal-close svg { width: 20px; height: 20px; }

.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 24px; }
.detail-item { background: var(--cream); border-radius: 10px; padding: 14px 16px; }
.detail-label { display: block; font-size: 0.72rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 6px; }
.detail-value { font-size: 0.95rem; font-weight: 700; color: var(--charcoal); }
.detail-yes { font-size: 0.95rem; font-weight: 700; color: var(--green-dark); }
.detail-no { font-size: 0.95rem; font-weight: 700; color: var(--charcoal); opacity: 0.4; }

.modal-actions { display: flex; gap: 12px; justify-content: flex-end; }
.btn-secondary { padding: 10px 20px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; }
.btn-secondary:hover { border-color: var(--charcoal); }
.btn-toggle { padding: 10px 20px; border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-toggle--on { background: var(--green-dark); color: var(--white); }
.btn-toggle--on:hover { background: var(--green-mid); }
.btn-toggle--off { background: #fef3c7; color: #92400e; }
.btn-toggle--off:hover { background: #fde68a; }
</style>
