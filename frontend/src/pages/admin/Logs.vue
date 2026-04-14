<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Log {
    id: number
    utilisateur: string
    action: string
    ressource: string
    ip: string
    date: string
    niveau: string
}

const logs = ref<Log[]>([])
const filterNiveau = ref('')

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/logs', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) logs.value = await res.json()
    } catch {}
})

const filtered = computed(() => logs.value.filter(l => !filterNiveau.value || l.niveau === filterNiveau.value))

function niveauClass(n: string) {
    if (n === 'error') return 'badge badge--error'
    if (n === 'warning') return 'badge badge--warning'
    return 'badge badge--info'
}
</script>

<template>
    <div class="logs">
        <div class="page-header">
            <h1 class="page-title">Logs.</h1>
            <p class="page-subtitle">Audit trail : qui a fait quoi et quand.</p>
        </div>

        <div class="filter-row">
            <select v-model="filterNiveau" class="filter-select">
                <option value="">Tous les niveaux</option>
                <option value="info">Info</option>
                <option value="warning">Warning</option>
                <option value="error">Error</option>
            </select>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Utilisateur</th>
                        <th>Action</th>
                        <th>Ressource</th>
                        <th>IP</th>
                        <th>Niveau</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="filtered.length === 0">
                        <td colspan="6" class="empty">Aucun log.</td>
                    </tr>
                    <tr v-for="l in filtered" :key="l.id">
                        <td class="td-mono">{{ l.date }}</td>
                        <td class="td-bold">{{ l.utilisateur }}</td>
                        <td>{{ l.action }}</td>
                        <td class="td-muted">{{ l.ressource }}</td>
                        <td class="td-mono">{{ l.ip }}</td>
                        <td><span :class="niveauClass(l.niveau)">{{ l.niveau }}</span></td>
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
.filter-row { margin-bottom: 16px; }
.filter-select { padding: 9px 14px; font-size: 0.88rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; cursor: pointer; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; min-width: 700px; }
.data-table th { text-align: left; padding: 14px 16px; font-size: 0.78rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); white-space: nowrap; }
.data-table td { padding: 12px 16px; font-size: 0.85rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); white-space: nowrap; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; }
.td-mono { font-family: 'Courier New', monospace; font-size: 0.78rem; opacity: 0.65; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; white-space: normal; }
.badge { display: inline-block; padding: 3px 8px; border-radius: 20px; font-size: 0.72rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; }
.badge--info { background: var(--green-pale); color: var(--green-dark); }
.badge--warning { background: #fef3c7; color: #92400e; }
.badge--error { background: #fee2e2; color: #991b1b; }
</style>
