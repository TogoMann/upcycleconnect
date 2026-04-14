<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Facture {
    id: number
    numero: string
    date: string
    montant: number
    statut: string
    pdf_url: string
}

const factures = ref<Facture[]>([])

onMounted(async () => {
    const token = authStore.token
    if (!token) return
    try {
        const res = await fetch('http://localhost:8081/pro/facturation', {
            headers: { Authorization: `Bearer ${token}` },
        })
        if (res.ok) factures.value = await res.json()
    } catch {}
})

function downloadPdf(url: string, numero: string) {
    const a = document.createElement('a')
    a.href = url
    a.download = `facture-${numero}.pdf`
    a.click()
}
</script>

<template>
    <div class="facturation">
        <div class="page-header">
            <h1 class="page-title">Facturation.</h1>
            <p class="page-subtitle">Historique de vos factures et téléchargement PDF.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>N° Facture</th>
                        <th>Date</th>
                        <th>Montant</th>
                        <th>Statut</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="factures.length === 0">
                        <td colspan="5" class="empty">Aucune facture.</td>
                    </tr>
                    <tr v-for="f in factures" :key="f.id">
                        <td class="td-bold">{{ f.numero }}</td>
                        <td class="td-muted">{{ f.date }}</td>
                        <td>{{ f.montant.toFixed(2) }} €</td>
                        <td>
                            <span class="badge" :class="f.statut === 'payée' ? 'badge--paid' : 'badge--pending'">
                                {{ f.statut }}
                            </span>
                        </td>
                        <td>
                            <button class="btn-dl" @click="downloadPdf(f.pdf_url, f.numero)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                                    <polyline points="7 10 12 15 17 10" />
                                    <line x1="12" y1="15" x2="12" y2="3" />
                                </svg>
                                PDF
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
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53, 53, 53, 0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53, 53, 53, 0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53, 53, 53, 0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215, 236, 225, 0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--paid { background: var(--green-pale); color: var(--green-dark); }
.badge--pending { background: #fef3c7; color: #92400e; }
.btn-dl { display: flex; align-items: center; gap: 6px; padding: 7px 12px; border-radius: 7px; border: 1.5px solid rgba(53, 53, 53, 0.15); background: transparent; font-size: 0.82rem; font-weight: 600; color: var(--charcoal); cursor: pointer; transition: border-color 0.2s, color 0.2s; }
.btn-dl:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-dl svg { width: 14px; height: 14px; }
</style>
