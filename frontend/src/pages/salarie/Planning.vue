<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Creneau {
    id: number
    titre: string
    date: string
    heure_debut: string
    heure_fin: string
    lieu: string
    type: string
}

const creneaux = ref<Creneau[]>([])

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/salarie/planning', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) creneaux.value = await res.json()
    } catch {}
})

function typeClass(t: string) {
    if (t === 'formation') return 'badge badge--formation'
    if (t === 'atelier') return 'badge badge--atelier'
    return 'badge badge--autre'
}
</script>

<template>
    <div class="planning">
        <div class="page-header">
            <h1 class="page-title">Planning.</h1>
            <p class="page-subtitle">Vos créneaux et interventions planifiés.</p>
        </div>

        <div v-if="creneaux.length === 0" class="empty-state">
            <div class="empty-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                    <line x1="16" y1="2" x2="16" y2="6" />
                    <line x1="8" y1="2" x2="8" y2="6" />
                    <line x1="3" y1="10" x2="21" y2="10" />
                </svg>
            </div>
            <p class="empty-text">Aucun créneau planifié.</p>
        </div>

        <div v-else class="creneaux-list">
            <div v-for="c in creneaux" :key="c.id" class="creneau-card">
                <div class="creneau-date">
                    <div class="date-day">{{ new Date(c.date).getDate() }}</div>
                    <div class="date-month">{{ new Date(c.date).toLocaleDateString('fr-FR', { month: 'short' }) }}</div>
                </div>
                <div class="creneau-body">
                    <div class="creneau-titre">{{ c.titre }}</div>
                    <div class="creneau-meta">
                        <span>{{ c.heure_debut }} – {{ c.heure_fin }}</span>
                        <span class="meta-sep">·</span>
                        <span>{{ c.lieu }}</span>
                    </div>
                </div>
                <span :class="typeClass(c.type)">{{ c.type }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.empty-state { display: flex; flex-direction: column; align-items: center; gap: 16px; padding: 80px 0; }
.empty-icon { width: 64px; height: 64px; background: var(--green-pale); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: var(--green-mid); }
.empty-icon svg { width: 32px; height: 32px; }
.empty-text { font-size: 0.95rem; color: var(--charcoal); opacity: 0.5; margin: 0; }
.creneaux-list { display: flex; flex-direction: column; gap: 12px; }
.creneau-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.08); border-radius: 12px; padding: 18px 20px; display: flex; align-items: center; gap: 20px; }
.creneau-card:hover { border-color: var(--green-light); }
.creneau-date { width: 48px; flex-shrink: 0; text-align: center; background: var(--green-pale); border-radius: 10px; padding: 8px 0; }
.date-day { font-size: 1.4rem; font-weight: 800; color: var(--green-dark); line-height: 1; }
.date-month { font-size: 0.72rem; font-weight: 600; color: var(--green-mid); text-transform: uppercase; }
.creneau-body { flex: 1; min-width: 0; }
.creneau-titre { font-size: 0.95rem; font-weight: 700; color: var(--charcoal); margin-bottom: 4px; }
.creneau-meta { font-size: 0.82rem; color: var(--charcoal); opacity: 0.5; display: flex; gap: 8px; }
.meta-sep { opacity: 0.3; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; white-space: nowrap; }
.badge--formation { background: var(--green-pale); color: var(--green-dark); }
.badge--atelier { background: #dbeafe; color: #1e40af; }
.badge--autre { background: rgba(53,53,53,0.08); color: var(--charcoal); }
</style>
