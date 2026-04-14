<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Conteneur {
    id: number
    code_barres: string
    localisation: string
    etat: string
    capacite: number
    objets: number
}

const conteneurs = ref<Conteneur[]>([])

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/conteneurs', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) conteneurs.value = await res.json()
    } catch {}
})

function etatClass(e: string) {
    if (e === 'actif') return 'badge badge--active'
    if (e === 'plein') return 'badge badge--warn'
    return 'badge badge--inactive'
}
</script>

<template>
    <div class="conteneurs">
        <div class="page-header">
            <h1 class="page-title">Conteneurs.</h1>
            <p class="page-subtitle">Liste, état et codes-barres de tous les conteneurs.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Code-barres</th>
                        <th>Localisation</th>
                        <th>Objets / Capacité</th>
                        <th>État</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="conteneurs.length === 0">
                        <td colspan="4" class="empty">Aucun conteneur.</td>
                    </tr>
                    <tr v-for="c in conteneurs" :key="c.id">
                        <td class="td-mono">{{ c.code_barres }}</td>
                        <td>{{ c.localisation }}</td>
                        <td>
                            <div class="progress-wrap">
                                <div class="progress-bar">
                                    <div class="progress-fill" :style="{ width: Math.min((c.objets / c.capacite) * 100, 100) + '%' }"></div>
                                </div>
                                <span class="progress-label">{{ c.objets }} / {{ c.capacite }}</span>
                            </div>
                        </td>
                        <td><span :class="etatClass(c.etat)">{{ c.etat }}</span></td>
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
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-mono { font-family: 'Courier New', monospace; font-size: 0.85rem; font-weight: 600; color: var(--green-dark); }
.progress-wrap { display: flex; align-items: center; gap: 10px; }
.progress-bar { flex: 1; max-width: 100px; height: 6px; background: rgba(53,53,53,0.1); border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; background: var(--green-mid); border-radius: 3px; transition: width 0.3s; }
.progress-label { font-size: 0.8rem; color: var(--charcoal); opacity: 0.6; white-space: nowrap; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--warn { background: #fef3c7; color: #92400e; }
.badge--inactive { background: rgba(53,53,53,0.08); color: var(--charcoal); }
</style>
