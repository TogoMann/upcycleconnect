<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Event {
    id: number
    approved: boolean
    price: number | null
    date: string | null
    created_by: number | null
    created_at: string | null
}

const events = ref<Event[]>([])
const loading = ref(true)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/event')
        if (res.ok) events.value = await res.json()
    } catch {}
    loading.value = false
})

function fmtDate(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleDateString('fr-FR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

function fmtPrice(price: number | null): string {
    if (price == null) return 'Gratuit'
    if (price === 0) return 'Gratuit'
    return `${price} €`
}
</script>

<template>
    <div class="page-content">
        <section class="hero">
            <div class="container">
                <h1 class="hero-title">Découvrez nos évènements.</h1>
            </div>
        </section>

        <section class="events-section">
            <div class="container">
                <div v-if="loading" class="loading">Chargement des évènements…</div>

                <div v-else-if="events.length === 0" class="empty">Aucun évènement à venir.</div>

                <div
                    v-else
                    v-for="(event, index) in events"
                    :key="event.id"
                    class="event-block"
                    :class="{ 'event-block--last': index === events.length - 1 }"
                >
                    <div class="event-info">
                        <div class="event-badges">
                            <span class="badge-date">{{ fmtDate(event.date) }}</span>
                            <span v-if="!event.approved" class="badge-pending">En attente de validation</span>
                        </div>
                        <h2 class="event-title">Évènement #{{ event.id }}</h2>
                    </div>

                    <div class="event-footer">
                        <span class="event-price">{{ fmtPrice(event.price) }}</span>
                        <button class="btn-reserver">Réserver</button>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page-content { flex: 1; display: flex; flex-direction: column; }
.container { max-width: 1060px; margin: 0 auto; padding: 0 32px; }
.hero { padding: 56px 0 36px; background: var(--cream); }
.hero-title { font-size: clamp(2.4rem, 5.5vw, 4rem); font-weight: 800; color: var(--charcoal); line-height: 1.08; letter-spacing: -0.03em; margin: 0; }
.events-section { flex: 1; padding: 48px 0 80px; }
.loading, .empty { opacity: 0.5; font-size: 0.9rem; padding: 40px 0; }
.event-block { padding-bottom: 48px; margin-bottom: 48px; border-bottom: 1px solid rgba(53,53,53,0.12); }
.event-block--last { border-bottom: none; margin-bottom: 0; }
.event-info { margin-bottom: 24px; }
.event-badges { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; margin-bottom: 12px; }
.badge-date { background: var(--green-pale); color: var(--green-dark); font-size: 0.82rem; font-weight: 600; padding: 6px 14px; border-radius: 6px; text-transform: capitalize; }
.badge-pending { background: #fef3c7; color: #92400e; font-size: 0.78rem; font-weight: 600; padding: 4px 10px; border-radius: 6px; }
.event-title { font-size: clamp(1.4rem, 2.8vw, 1.9rem); font-weight: 700; color: var(--charcoal); line-height: 1.2; margin: 0; letter-spacing: -0.01em; }
.event-footer { display: flex; align-items: center; justify-content: flex-end; gap: 16px; }
.event-price { font-size: 1.1rem; font-weight: 700; color: var(--green-dark); }
.btn-reserver { background: var(--green-mid); color: var(--white); border: none; padding: 10px 24px; border-radius: 6px; font-size: 0.875rem; font-weight: 600; cursor: pointer; transition: background 0.2s, transform 0.15s; font-family: inherit; }
.btn-reserver:hover { background: var(--green-dark); transform: translateY(-1px); }
</style>
