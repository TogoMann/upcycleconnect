<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useClientStore } from '@/stores/client'

const clientStore = useClientStore()

const sortedEntries = computed(() => {
    return [...clientStore.entries].sort((a, b) => {
        const da = new Date(a.schedule?.Time ?? a.schedule)
        const db = new Date(b.schedule?.Time ?? b.schedule)
        return da.getTime() - db.getTime()
    })
})

const upcomingEntries = computed(() =>
    sortedEntries.value.filter(e => {
        const d = new Date(e.schedule?.Time ?? e.schedule)
        return d >= new Date(new Date().toDateString())
    })
)

const pastEntries = computed(() =>
    sortedEntries.value.filter(e => {
        const d = new Date(e.schedule?.Time ?? e.schedule)
        return d < new Date(new Date().toDateString())
    })
)

function formatDate(schedule: any): string {
    if (!schedule) return '—'
    const date = new Date(schedule?.Time ?? schedule)
    return date.toLocaleDateString('fr-FR', {
        weekday: 'long',
        day: '2-digit',
        month: 'long',
        year: 'numeric',
    })
}

function formatTime(t: any): string {
    if (!t) return '—'
    const raw = t?.Microseconds ?? t?.Time ?? t
    if (typeof raw === 'number') {
        const totalSeconds = Math.floor(raw / 1_000_000)
        const h = Math.floor(totalSeconds / 3600)
        const m = Math.floor((totalSeconds % 3600) / 60)
        return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
    }
    return String(raw).slice(0, 5)
}

async function handleDelete(id: number) {
    await clientStore.deleteEntry(id)
}

onMounted(() => {
    clientStore.fetchEntries()
})
</script>

<template>
    <div class="page">
        <div class="page-header">
            <h1 class="page-title">Planning.</h1>
            <router-link to="/particulier/conteneurs/deposer" class="btn-primary">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19" />
                    <line x1="5" y1="12" x2="19" y2="12" />
                </svg>
                Nouveau créneau
            </router-link>
        </div>

        <div v-if="clientStore.isLoading" class="state-empty">
            <p>Chargement…</p>
        </div>

        <div v-else-if="clientStore.entries.length === 0" class="state-empty">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                    <line x1="16" y1="2" x2="16" y2="6" />
                    <line x1="8" y1="2" x2="8" y2="6" />
                    <line x1="3" y1="10" x2="21" y2="10" />
                </svg>
            </div>
            <p class="empty-title">Aucun créneau planifié</p>
            <p class="empty-sub">Réservez un créneau pour déposer vos objets dans nos conteneurs partenaires.</p>
            <router-link to="/particulier/conteneurs/deposer" class="btn-primary btn-primary--mt">
                Planifier un créneau
            </router-link>
        </div>

        <template v-else>
            <div v-if="upcomingEntries.length > 0" class="entries-section">
                <h2 class="section-title">À venir ({{ upcomingEntries.length }})</h2>
                <div class="entries-list">
                    <div v-for="entry in upcomingEntries" :key="entry.id?.Int64" class="entry-card entry-card--upcoming">
                        <div class="entry-date-block">
                            <span class="entry-day">{{ new Date(entry.schedule?.Time ?? entry.schedule).toLocaleDateString('fr-FR', { day: '2-digit' }) }}</span>
                            <span class="entry-month">{{ new Date(entry.schedule?.Time ?? entry.schedule).toLocaleDateString('fr-FR', { month: 'short' }) }}</span>
                        </div>
                        <div class="entry-info">
                            <p class="entry-full-date">{{ formatDate(entry.schedule) }}</p>
                            <p class="entry-time">
                                {{ formatTime(entry.start) }} — {{ formatTime(entry.ending) }}
                            </p>
                        </div>
                        <div class="entry-tag">À venir</div>
                        <button class="btn-delete" @click="handleDelete(entry.id?.Int64)" title="Annuler">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="3 6 5 6 21 6" />
                                <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
                                <path d="M10 11v6" />
                                <path d="M14 11v6" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="pastEntries.length > 0" class="entries-section">
                <h2 class="section-title">Passés ({{ pastEntries.length }})</h2>
                <div class="entries-list">
                    <div v-for="entry in pastEntries" :key="entry.id?.Int64" class="entry-card entry-card--past">
                        <div class="entry-date-block entry-date-block--past">
                            <span class="entry-day">{{ new Date(entry.schedule?.Time ?? entry.schedule).toLocaleDateString('fr-FR', { day: '2-digit' }) }}</span>
                            <span class="entry-month">{{ new Date(entry.schedule?.Time ?? entry.schedule).toLocaleDateString('fr-FR', { month: 'short' }) }}</span>
                        </div>
                        <div class="entry-info">
                            <p class="entry-full-date">{{ formatDate(entry.schedule) }}</p>
                            <p class="entry-time">
                                {{ formatTime(entry.start) }} — {{ formatTime(entry.ending) }}
                            </p>
                        </div>
                        <div class="entry-tag entry-tag--past">Passé</div>
                    </div>
                </div>
            </div>
        </template>
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
.btn-primary {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: var(--green-dark);
    color: var(--white);
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    transition: background 0.2s;
}
.btn-primary svg {
    width: 16px;
    height: 16px;
}
.btn-primary:hover {
    background: var(--green-mid);
}
.btn-primary--mt {
    margin-top: 16px;
    display: inline-flex;
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
.entry-full-date {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--charcoal);
    margin: 0 0 3px;
    text-transform: capitalize;
}
.entry-time {
    font-size: 0.8rem;
    color: var(--charcoal);
    opacity: 0.55;
    margin: 0;
}
.entry-tag {
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    background: var(--green-pale);
    color: var(--green-dark);
    padding: 4px 10px;
    border-radius: 20px;
    flex-shrink: 0;
}
.entry-tag--past {
    background: rgba(53, 53, 53, 0.08);
    color: rgba(53, 53, 53, 0.45);
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
</style>
