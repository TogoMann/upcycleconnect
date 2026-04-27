<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()

const showModal = ref(false)
const form = reactive({
    title: '',
    description: '',
    date: '',
    start_time: '',
    end_time: '',
    all_day: false
})

const sortedItems = computed(() => {
    return [...clientStore.planning].sort((a, b) => {
        const da = new Date(a.date)
        const db = new Date(b.date)
        if (da.getTime() !== db.getTime()) return da.getTime() - db.getTime()
        return a.start_time.localeCompare(b.start_time)
    })
})

const upcomingItems = computed(() =>
    sortedItems.value.filter(e => {
        const d = new Date(e.date)
        const today = new Date()
        today.setHours(0, 0, 0, 0)
        return d >= today
    })
)

const pastItems = computed(() =>
    sortedItems.value.filter(e => {
        const d = new Date(e.date)
        const today = new Date()
        today.setHours(0, 0, 0, 0)
        return d < today
    })
)

function formatDate(dateStr: string): string {
    if (!dateStr) return '—'
    const date = new Date(dateStr)
    return date.toLocaleDateString('fr-FR', {
        weekday: 'long',
        day: '2-digit',
        month: 'long',
        year: 'numeric',
    })
}

function getTypeLabel(type: string): string {
    switch (type) {
        case 'depot': return 'Dépôt'
        case 'workshop': return 'Atelier'
        case 'event': return 'Événement'
        case 'personal': return 'Personnel'
        default: return 'Autre'
    }
}

async function handleDelete(item: any) {
    if (item.type === 'personal') {
        if (confirm('Supprimer cet événement personnel ?')) {
            await clientStore.deletePersonalEvent(item.id)
        }
    } else if (item.type === 'depot') {
        if (confirm('Annuler ce créneau de dépôt ?')) {
            await clientStore.deleteEntry(item.id)
            await clientStore.fetchPlanning()
        }
    } else {
        alert('Cette action ne peut être annulée depuis le planning. Veuillez contacter le support.')
    }
}

async function submitPersonalEvent() {
    try {
        const payload = { ...form }
        if (payload.all_day) {
            payload.start_time = '00:00'
            payload.end_time = '23:59'
        }
        
        if (!payload.start_time || !payload.end_time) {
            throw new Error('Veuillez renseigner les horaires ou cocher "Toute la journée"')
        }

        await clientStore.createPersonalEvent(payload)
        showModal.value = false
        // Reset form
        form.title = ''
        form.description = ''
        form.date = ''
        form.start_time = ''
        form.end_time = ''
        form.all_day = false
    } catch (e: any) {
        alert(e.message)
    }
}

onMounted(() => {
    clientStore.fetchPlanning()
})
</script>

<template>
    <div class="page">
        <div class="page-header">
            <h1 class="page-title">Planning.</h1>
            <div class="header-actions">
                <router-link to="/particulier/conteneurs/deposer" class="btn-secondary">
                    Nouveau dépôt
                </router-link>
                <button class="btn-primary" @click="showModal = true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="12" y1="5" x2="12" y2="19" />
                        <line x1="5" y1="12" x2="19" y2="12" />
                    </svg>
                    Nouvel événement
                </button>
            </div>
        </div>

        <div v-if="clientStore.isLoading" class="state-empty">
            <p>Chargement…</p>
        </div>

        <div v-else-if="clientStore.planning.length === 0" class="state-empty">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                    <line x1="16" y1="2" x2="16" y2="6" />
                    <line x1="8" y1="2" x2="8" y2="6" />
                    <line x1="3" y1="10" x2="21" y2="10" />
                </svg>
            </div>
            <p class="empty-title">Votre planning est vide</p>
            <p class="empty-sub">Planifiez vos activités d'upcycling ou ajoutez vos propres événements.</p>
        </div>

        <template v-else>
            <div v-if="upcomingItems.length > 0" class="entries-section">
                <h2 class="section-title">À venir ({{ upcomingItems.length }})</h2>
                <div class="entries-list">
                    <div v-for="item in upcomingItems" :key="item.type + '-' + item.id" class="entry-card" :class="'entry-card--' + item.type">
                        <div class="entry-date-block">
                            <span class="entry-day">{{ new Date(item.date).toLocaleDateString('fr-FR', { day: '2-digit' }) }}</span>
                            <span class="entry-month">{{ new Date(item.date).toLocaleDateString('fr-FR', { month: 'short' }) }}</span>
                        </div>
                        <div class="entry-info">
                            <div class="entry-header-row">
                                <span class="entry-type-tag">{{ getTypeLabel(item.type) }}</span>
                                <p class="entry-full-date">{{ formatDate(item.date) }}</p>
                            </div>
                            <h3 class="entry-title">{{ item.title }}</h3>
                            <p class="entry-time">
                                {{ item.start_time }} — {{ item.end_time }}
                                <span v-if="item.location" class="entry-location"> • {{ item.location }}</span>
                            </p>
                            <p v-if="item.description" class="entry-desc">{{ item.description }}</p>
                        </div>
                        <button v-if="item.type === 'personal' || item.type === 'depot'" class="btn-delete" @click="handleDelete(item)" title="Supprimer">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="3 6 5 6 21 6" />
                                <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="pastItems.length > 0" class="entries-section">
                <h2 class="section-title">Passés ({{ pastItems.length }})</h2>
                <div class="entries-list">
                    <div v-for="item in pastItems" :key="item.type + '-' + item.id" class="entry-card entry-card--past">
                        <div class="entry-date-block entry-date-block--past">
                            <span class="entry-day">{{ new Date(item.date).toLocaleDateString('fr-FR', { day: '2-digit' }) }}</span>
                            <span class="entry-month">{{ new Date(item.date).toLocaleDateString('fr-FR', { month: 'short' }) }}</span>
                        </div>
                        <div class="entry-info">
                            <div class="entry-header-row">
                                <span class="entry-type-tag">{{ getTypeLabel(item.type) }}</span>
                                <p class="entry-full-date">{{ formatDate(item.date) }}</p>
                            </div>
                            <h3 class="entry-title">{{ item.title }}</h3>
                        </div>
                    </div>
                </div>
            </div>
        </template>

        <!-- Modal Nouvel Evénement -->
        <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
            <div class="modal-content">
                <div class="modal-header">
                    <h2 class="modal-title">Nouvel événement</h2>
                    <button class="btn-close" @click="showModal = false">×</button>
                </div>
                <form @submit.prevent="submitPersonalEvent" class="modal-body">
                    <div class="form-group">
                        <label>Titre</label>
                        <input v-model="form.title" type="text" required placeholder="Ex: Déplacement Lyon" />
                    </div>
                    <div class="form-group">
                        <label>Description (optionnel)</label>
                        <textarea v-model="form.description" rows="3"></textarea>
                    </div>
                    <div class="form-group">
                        <label class="checkbox-label">
                            <input v-model="form.all_day" type="checkbox" />
                            Toute la journée
                        </label>
                    </div>
                    <div class="form-row">
                        <div class="form-group">
                            <label>Date</label>
                            <input v-model="form.date" type="date" required :min="new Date().toISOString().split('T')[0]" />
                        </div>
                    </div>
                    <div v-if="!form.all_day" class="form-row">
                        <div class="form-group">
                            <label>Début</label>
                            <input v-model="form.start_time" type="time" :required="!form.all_day" />
                        </div>
                        <div class="form-group">
                            <label>Fin</label>
                            <input v-model="form.end_time" type="time" :required="!form.all_day" />
                        </div>
                    </div>
                    <div class="modal-actions">
                        <button type="button" class="btn-cancel" @click="showModal = false">Annuler</button>
                        <button type="submit" class="btn-submit">Enregistrer</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}
.page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 32px;
    flex-wrap: wrap;
    gap: 16px;
}
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0;
    line-height: 1.08;
}
.header-actions {
    display: flex;
    gap: 12px;
}
.btn-primary, .btn-secondary {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    transition: all 0.2s;
    cursor: pointer;
    border: none;
    font-family: inherit;
}
.btn-primary {
    background: var(--green-dark);
    color: var(--white);
}
.btn-primary:hover {
    background: var(--green-mid);
}
.btn-secondary {
    background: var(--green-pale);
    color: var(--green-dark);
}
.btn-secondary:hover {
    background: var(--green-light);
    color: var(--white);
}

.state-empty {
    text-align: center;
    padding: 64px 32px;
}
.empty-icon {
    width: 64px;
    height: 64px;
    background: var(--green-pale);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    color: var(--green-mid);
}
.empty-icon svg {
    width: 28px;
    height: 28px;
}
.empty-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 8px;
}
.empty-sub {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.6;
    max-width: 360px;
    margin: 0 auto;
    line-height: 1.6;
}

.entries-section {
    margin-bottom: 36px;
}
.section-title {
    font-size: 1rem;
    font-weight: 700;
    color: var(--charcoal);
    margin: 0 0 14px;
    opacity: 0.7;
}
.entries-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.entry-card {
    display: flex;
    align-items: center;
    gap: 16px;
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 12px;
    padding: 16px 20px;
}
.entry-card--upcoming {
    border-color: var(--green-light);
}
.entry-card--depot { border-left: 4px solid #3498db; }
.entry-card--workshop { border-left: 4px solid #9b59b6; }
.entry-card--event { border-left: 4px solid #f1c40f; }
.entry-card--personal { border-left: 4px solid #1abc9c; }

.entry-card--past {
    opacity: 0.6;
}
.entry-date-block {
    width: 48px;
    text-align: center;
    flex-shrink: 0;
    background: var(--green-pale);
    border-radius: 10px;
    padding: 8px 4px;
}
.entry-date-block--past {
    background: rgba(53, 53, 53, 0.06);
}
.entry-day {
    display: block;
    font-size: 1.3rem;
    font-weight: 800;
    color: var(--green-dark);
    letter-spacing: -0.02em;
    line-height: 1;
}
.entry-date-block--past .entry-day {
    color: var(--charcoal);
    opacity: 0.5;
}
.entry-month {
    display: block;
    font-size: 0.7rem;
    font-weight: 600;
    color: var(--green-mid);
    text-transform: uppercase;
    letter-spacing: 0.04em;
}
.entry-date-block--past .entry-month {
    color: var(--charcoal);
    opacity: 0.45;
}
.entry-info {
    flex: 1;
    min-width: 0;
}
.entry-header-row {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 4px;
}
.entry-type-tag {
    font-size: 0.65rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding: 2px 8px;
    border-radius: 4px;
    background: rgba(53, 53, 53, 0.05);
    color: rgba(53, 53, 53, 0.6);
}
.entry-card--depot .entry-type-tag { background: rgba(52, 152, 219, 0.1); color: #2980b9; }
.entry-card--workshop .entry-type-tag { background: rgba(155, 89, 182, 0.1); color: #8e44ad; }
.entry-card--event .entry-type-tag { background: rgba(241, 196, 15, 0.1); color: #f39c12; }
.entry-card--personal .entry-type-tag { background: rgba(26, 188, 156, 0.1); color: #16a085; }

.entry-full-date {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--charcoal);
    opacity: 0.5;
    margin: 0;
    text-transform: capitalize;
}
.entry-title {
    font-size: 1rem;
    font-weight: 700;
    margin: 0 0 4px;
    color: var(--charcoal);
}
.entry-time {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.55;
    margin: 0;
}
.entry-location {
    font-weight: 600;
    color: var(--green-dark);
}
.entry-desc {
    font-size: 0.85rem;
    margin: 8px 0 0;
    line-height: 1.5;
    opacity: 0.7;
}

.btn-delete {
    background: none;
    border: none;
    cursor: pointer;
    padding: 6px;
    color: rgba(53, 53, 53, 0.3);
    border-radius: 6px;
    transition: color 0.2s, background 0.2s;
    flex-shrink: 0;
}
.btn-delete:hover {
    color: #e53e3e;
    background: rgba(229, 62, 62, 0.08);
}
.btn-delete svg {
    width: 16px;
    height: 16px;
}

/* Modal Styles */
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
}
.modal-content {
    background: var(--white);
    width: 100%;
    max-width: 500px;
    border-radius: 16px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
    overflow: hidden;
}
.modal-header {
    padding: 20px 24px;
    border-bottom: 1px solid rgba(53, 53, 53, 0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.modal-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 800;
}
.btn-close {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: rgba(53, 53, 53, 0.4);
}
.modal-body {
    padding: 24px;
}
.form-group {
    margin-bottom: 16px;
}
.form-group label {
    display: block;
    font-size: 0.85rem;
    font-weight: 700;
    margin-bottom: 6px;
}
.checkbox-label {
    display: flex !important;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    user-select: none;
}
.checkbox-label input {
    width: auto !important;
}
.form-group input, .form-group textarea {
    width: 100%;
    padding: 10px 12px;
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 8px;
    font-family: inherit;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
}
.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
}
.btn-cancel {
    background: none;
    border: none;
    font-weight: 600;
    cursor: pointer;
}
.btn-submit {
    background: var(--green-dark);
    color: white;
    border: none;
    padding: 10px 24px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
}
</style>
