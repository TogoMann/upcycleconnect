<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Ad {
    id: { Int64: number; Valid: boolean }
    announcer_id: { Int64: number; Valid: boolean }
    target_id: number
    target_type: string
    ad_type: string
    budget: { Int: any; Exp: number; Status: number } | number
    status: string
    start_date: { Time: string; Valid: boolean }
    end_date: { Time: string; Valid: boolean }
    created_at: { Time: string; Valid: boolean }
}

interface Listing {
    id: { Int64: number; Valid: boolean }
    name: string
}

interface Project {
    id: { Int64: number; Valid: boolean }
    title: string
}

interface Brand {
    id: { Int64: number; Valid: boolean }
    name: string
}

const ads = ref<Ad[]>([])
const listings = ref<Listing[]>([])
const projects = ref<Project[]>([])
const brands = ref<Brand[]>([])
const loading = ref(true)
const showCreateModal = ref(false)
const creating = ref(false)
const createError = ref('')
const filter = ref('all')

const form = ref({
    target_type: 'listing' as 'listing' | 'project' | 'brand',
    target_id: 0,
    ad_type: 'partnership' as 'partnership' | 'other',
    budget: '',
    start_date: '',
    end_date: '',
})

const filteredAds = computed(() => {
    if (filter.value === 'all') return ads.value
    return ads.value.filter(a => a.status === filter.value)
})

const stats = computed(() => ({
    total: ads.value.length,
    active: ads.value.filter(a => a.status === 'validated').length,
    pending: ads.value.filter(a => a.status === 'pending').length,
    totalBudget: ads.value.reduce((sum, a) => sum + getBudget(a), 0),
}))

function getBudget(ad: Ad): number {
    if (typeof ad.budget === 'number') return ad.budget
    if (ad.budget && typeof ad.budget === 'object') {
        try {
            const val = Number(ad.budget.Int?.Int64 || ad.budget.Int || 0)
            return val * Math.pow(10, ad.budget.Exp || 0)
        } catch { return 0 }
    }
    return 0
}

function formatDate(d: { Time: string; Valid: boolean } | undefined) {
    if (!d?.Valid) return '—'
    return new Date(d.Time).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function statusConfig(s: string) {
    const map: Record<string, { label: string; class: string }> = {
        pending: { label: 'En attente', class: 'badge--pending' },
        validated: { label: 'Active', class: 'badge--active' },
        rejected: { label: 'Refusée', class: 'badge--rejected' },
        expired: { label: 'Expirée', class: 'badge--expired' },
    }
    return map[s] || { label: s, class: 'badge--default' }
}

function targetTypeLabel(t: string) {
    const map: Record<string, string> = { listing: 'Annonce', project: 'Projet', brand: 'Marque' }
    return map[t] || t
}

function adTypeLabel(t: string) {
    return t === 'partnership' ? 'Partenariat' : 'Autre'
}

const targetOptions = computed(() => {
    if (form.value.target_type === 'listing') return listings.value.map(l => ({ id: l.id.Int64, name: l.name }))
    if (form.value.target_type === 'project') return projects.value.map(p => ({ id: p.id.Int64, name: p.title }))
    if (form.value.target_type === 'brand') return brands.value.map(b => ({ id: b.id.Int64, name: b.name }))
    return []
})

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    const headers = { Authorization: `Bearer ${token}` }

    const [adsRes, listingsRes, projectsRes, brandsRes] = await Promise.all([
        fetch(`${API_BASE}/advertisement/me`, { headers }).catch(() => null),
        fetch(`${API_BASE}/listing`, { headers }).catch(() => null),
        fetch(`${API_BASE}/project/me`, { headers }).catch(() => null),
        fetch(`${API_BASE}/brands/me`, { headers }).catch(() => null),
    ])

    if (adsRes?.ok) {
        const data = await adsRes.json()
        ads.value = Array.isArray(data) ? data : []
    }
    if (listingsRes?.ok) {
        const data = await listingsRes.json()
        listings.value = Array.isArray(data) ? data : []
    }
    if (projectsRes?.ok) {
        const data = await projectsRes.json()
        projects.value = Array.isArray(data) ? data : []
    }
    if (brandsRes?.ok) {
        const data = await brandsRes.json()
        brands.value = Array.isArray(data) ? data : []
    }

    loading.value = false
})

function openCreate() {
    form.value = { target_type: 'listing', target_id: 0, ad_type: 'partnership', budget: '', start_date: '', end_date: '' }
    createError.value = ''
    showCreateModal.value = true
}

async function submitCreate() {
    if (!form.value.target_id) { createError.value = 'Veuillez sélectionner une cible.'; return }
    if (!form.value.budget || parseFloat(form.value.budget) <= 0) { createError.value = 'Le budget est requis.'; return }
    if (!form.value.start_date || !form.value.end_date) { createError.value = 'Les dates sont requises.'; return }
    if (new Date(form.value.start_date) >= new Date(form.value.end_date)) { createError.value = 'La date de fin doit être après la date de début.'; return }

    creating.value = true
    createError.value = ''
    try {
        const body = {
            target_id: form.value.target_id,
            target_type: form.value.target_type,
            ad_type: form.value.ad_type,
            budget: parseFloat(form.value.budget),
            status: 'pending',
            start_date: { Time: new Date(form.value.start_date).toISOString(), Valid: true },
            end_date: { Time: new Date(form.value.end_date).toISOString(), Valid: true },
        }
        const res = await fetch(`${API_BASE}/advertisement`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify(body),
        })
        if (res.ok) {
            showCreateModal.value = false
            const adsRes = await fetch(`${API_BASE}/advertisement/me`, {
                headers: { Authorization: `Bearer ${authStore.token}` },
            })
            if (adsRes.ok) {
                const data = await adsRes.json()
                ads.value = Array.isArray(data) ? data : []
            }
        } else {
            const d = await res.text()
            createError.value = d || 'Erreur lors de la création.'
        }
    } catch {
        createError.value = 'Erreur réseau.'
    }
    creating.value = false
}
</script>

<template>
    <div class="publicites">
        <div class="page-header">
            <div class="header-row">
                <div>
                    <h1 class="page-title">Publicités.</h1>
                    <p class="page-subtitle">Créez et suivez vos campagnes publicitaires sur la plateforme.</p>
                </div>
                <button class="btn-primary" @click="openCreate">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                    Nouvelle publicité
                </button>
            </div>
        </div>

        <div v-if="loading" class="loading-state">Chargement...</div>

        <template v-else>
            <!-- KPIs -->
            <div class="kpi-row">
                <div class="kpi-card">
                    <div class="kpi-value">{{ stats.total }}</div>
                    <div class="kpi-label">Total</div>
                </div>
                <div class="kpi-card kpi-card--green">
                    <div class="kpi-value">{{ stats.active }}</div>
                    <div class="kpi-label">Actives</div>
                </div>
                <div class="kpi-card kpi-card--yellow">
                    <div class="kpi-value">{{ stats.pending }}</div>
                    <div class="kpi-label">En attente</div>
                </div>
                <div class="kpi-card">
                    <div class="kpi-value">{{ stats.totalBudget.toFixed(0) }}€</div>
                    <div class="kpi-label">Budget total</div>
                </div>
            </div>

            <!-- Filtres -->
            <div class="filter-row">
                <button v-for="f in ['all', 'pending', 'validated', 'rejected', 'expired']" :key="f"
                    class="filter-btn" :class="{ 'filter-btn--active': filter === f }"
                    @click="filter = f">
                    {{ f === 'all' ? 'Toutes' : statusConfig(f).label }}
                </button>
            </div>

            <!-- Empty state -->
            <div v-if="filteredAds.length === 0" class="empty-state">
                <div class="empty-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
                </div>
                <p class="empty-text">{{ filter === 'all' ? 'Aucune publicité pour le moment.' : 'Aucune publicité avec ce statut.' }}</p>
                <button v-if="filter === 'all'" class="btn-primary" @click="openCreate">Créer une publicité</button>
            </div>

            <!-- Liste des pubs -->
            <div v-else class="ads-list">
                <div v-for="ad in filteredAds" :key="ad.id.Int64" class="ad-card">
                    <div class="ad-header">
                        <div class="ad-types">
                            <span class="ad-target-type">{{ targetTypeLabel(ad.target_type) }}</span>
                            <span class="ad-separator">-</span>
                            <span class="ad-ad-type">{{ adTypeLabel(ad.ad_type) }}</span>
                        </div>
                        <span class="badge" :class="statusConfig(ad.status).class">{{ statusConfig(ad.status).label }}</span>
                    </div>
                    <div class="ad-body">
                        <div class="ad-info">
                            <div class="ad-info-item">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>
                                <span>{{ getBudget(ad).toFixed(2) }}€</span>
                            </div>
                            <div class="ad-info-item">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                                <span>{{ formatDate(ad.start_date) }} → {{ formatDate(ad.end_date) }}</span>
                            </div>
                        </div>
                    </div>
                    <div class="ad-footer">
                        <span class="ad-created">Créée le {{ formatDate(ad.created_at) }}</span>
                    </div>
                </div>
            </div>
        </template>

        <!-- Modal création -->
        <Teleport to="body">
            <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
                <div class="modal-card">
                    <h3 class="modal-title">Nouvelle publicité</h3>
                    <p class="modal-subtitle">Votre publicité sera soumise à validation par l'administration.</p>

                    <form class="modal-form" @submit.prevent="submitCreate">
                        <div class="form-row">
                            <div class="form-group">
                                <label class="form-label">Type de cible *</label>
                                <select v-model="form.target_type" class="form-input" @change="form.target_id = 0">
                                    <option value="listing">Annonce</option>
                                    <option value="project">Projet</option>
                                    <option value="brand">Marque</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Cible *</label>
                                <select v-model.number="form.target_id" class="form-input">
                                    <option :value="0" disabled>Sélectionner...</option>
                                    <option v-for="opt in targetOptions" :key="opt.id" :value="opt.id">{{ opt.name }}</option>
                                </select>
                                <p v-if="targetOptions.length === 0" class="form-hint">Aucun(e) {{ targetTypeLabel(form.target_type).toLowerCase() }} disponible. Créez-en un(e) d'abord.</p>
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label class="form-label">Type de publicité *</label>
                                <select v-model="form.ad_type" class="form-input">
                                    <option value="partnership">Partenariat</option>
                                    <option value="other">Autre</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Budget (€) *</label>
                                <input v-model="form.budget" type="number" min="100" max="500" step="10" class="form-input" placeholder="100 - 500" />
                                <p class="form-hint">Entre 100€ et 500€/mois</p>
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label class="form-label">Date de début *</label>
                                <input v-model="form.start_date" type="date" class="form-input" />
                            </div>
                            <div class="form-group">
                                <label class="form-label">Date de fin *</label>
                                <input v-model="form.end_date" type="date" class="form-input" />
                            </div>
                        </div>

                        <div v-if="createError" class="alert alert--error">{{ createError }}</div>

                        <div class="modal-actions">
                            <button type="button" class="btn-secondary" @click="showCreateModal = false" :disabled="creating">Annuler</button>
                            <button type="submit" class="btn-primary" :disabled="creating">
                                {{ creating ? 'Envoi...' : 'Soumettre la publicité' }}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </Teleport>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 24px; }
.header-row { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.btn-primary { display: inline-flex; align-items: center; gap: 8px; padding: 11px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; white-space: nowrap; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.btn-primary svg { width: 16px; height: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }

/* KPIs */
.kpi-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 20px; }
.kpi-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.08); border-radius: 12px; padding: 18px 20px; }
.kpi-card--green { background: var(--green-pale); border-color: transparent; }
.kpi-card--yellow { background: #fffbeb; border-color: transparent; }
.kpi-value { font-size: 1.6rem; font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; line-height: 1; }
.kpi-card--green .kpi-value { color: var(--green-dark); }
.kpi-card--yellow .kpi-value { color: #92400e; }
.kpi-label { font-size: 0.78rem; font-weight: 500; color: var(--charcoal); opacity: 0.5; margin-top: 4px; }

/* Filters */
.filter-row { display: flex; gap: 8px; margin-bottom: 20px; flex-wrap: wrap; }
.filter-btn { padding: 7px 16px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 20px; background: transparent; color: var(--charcoal); font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.filter-btn:hover { border-color: var(--green-mid); color: var(--green-dark); }
.filter-btn--active { background: var(--green-dark); color: var(--white); border-color: var(--green-dark); }

/* Empty state */
.empty-state { display: flex; flex-direction: column; align-items: center; gap: 16px; padding: 60px 0; }
.empty-icon { width: 64px; height: 64px; background: var(--green-pale); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: var(--green-mid); }
.empty-icon svg { width: 32px; height: 32px; }
.empty-text { font-size: 0.95rem; color: var(--charcoal); opacity: 0.5; margin: 0; }

/* Ads list */
.ads-list { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; }
.ad-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 14px; padding: 20px; display: flex; flex-direction: column; gap: 14px; transition: border-color 0.2s; }
.ad-card:hover { border-color: var(--green-mid); }
.ad-header { display: flex; justify-content: space-between; align-items: flex-start; }
.ad-types { display: flex; align-items: center; gap: 6px; }
.ad-target-type { font-size: 0.9rem; font-weight: 700; color: var(--charcoal); }
.ad-separator { opacity: 0.3; }
.ad-ad-type { font-size: 0.84rem; color: var(--charcoal); opacity: 0.6; }
.ad-body { flex: 1; }
.ad-info { display: flex; flex-direction: column; gap: 8px; }
.ad-info-item { display: flex; align-items: center; gap: 8px; font-size: 0.86rem; color: var(--charcoal); opacity: 0.7; }
.ad-info-item svg { width: 16px; height: 16px; flex-shrink: 0; opacity: 0.5; }
.ad-footer { border-top: 1px solid rgba(53,53,53,0.06); padding-top: 10px; }
.ad-created { font-size: 0.78rem; color: var(--charcoal); opacity: 0.4; }

/* Badges */
.badge { display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; white-space: nowrap; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--pending { background: #fef3c7; color: #92400e; }
.badge--rejected { background: #fee2e2; color: #991b1b; }
.badge--expired { background: rgba(53,53,53,0.08); color: var(--charcoal); opacity: 0.6; }
.badge--default { background: rgba(53,53,53,0.08); color: var(--charcoal); }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.modal-card { background: var(--white); border-radius: 20px; padding: 32px; max-width: 580px; width: 90%; box-shadow: 0 20px 60px rgba(0,0,0,0.15); max-height: 90vh; overflow-y: auto; }
.modal-title { font-size: 1.2rem; font-weight: 800; color: var(--charcoal); margin: 0 0 4px; }
.modal-subtitle { font-size: 0.85rem; color: var(--charcoal); opacity: 0.5; margin: 0 0 24px; }
.modal-form { display: flex; flex-direction: column; gap: 18px; }
.modal-actions { display: flex; gap: 12px; justify-content: flex-end; padding-top: 8px; }

/* Form */
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-label { font-size: 0.82rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 10px 14px; font-size: 0.88rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-hint { font-size: 0.76rem; color: var(--charcoal); opacity: 0.4; margin: 0; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.86rem; font-weight: 500; }
.alert--error { background: #fee2e2; color: #991b1b; }
.btn-secondary { padding: 12px 24px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: border-color 0.2s; }
.btn-secondary:hover { border-color: var(--charcoal); }

@media (max-width: 800px) {
    .kpi-row { grid-template-columns: repeat(2, 1fr); }
    .ads-list { grid-template-columns: 1fr; }
    .form-row { grid-template-columns: 1fr; }
}
</style>
