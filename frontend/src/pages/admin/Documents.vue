<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Document {
    id: number
    nom: string
    type: string
    utilisateur: string
    date: string
    taille: string
    url: string
}

const documents = ref<Document[]>([])

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/documents', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) documents.value = await res.json()
    } catch {}
})

function download(doc: Document) {
    const a = document.createElement('a')
    a.href = doc.url
    a.download = doc.nom
    a.click()
}
</script>

<template>
    <div class="documents">
        <div class="page-header">
            <h1 class="page-title">Documents.</h1>
            <p class="page-subtitle">Tous les PDF générés sur la plateforme.</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Nom</th>
                        <th>Type</th>
                        <th>Utilisateur</th>
                        <th>Date</th>
                        <th>Taille</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="documents.length === 0">
                        <td colspan="6" class="empty">Aucun document.</td>
                    </tr>
                    <tr v-for="d in documents" :key="d.id">
                        <td>
                            <div class="doc-name">
                                <div class="doc-icon">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                                        <polyline points="14 2 14 8 20 8" />
                                    </svg>
                                </div>
                                <span class="td-bold">{{ d.nom }}</span>
                            </div>
                        </td>
                        <td>
                            <span class="badge badge--type">{{ d.type }}</span>
                        </td>
                        <td class="td-muted">{{ d.utilisateur }}</td>
                        <td class="td-muted">{{ d.date }}</td>
                        <td class="td-muted">{{ d.taille }}</td>
                        <td>
                            <button class="btn-dl" @click="download(d)">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                                    <polyline points="7 10 12 15 17 10" />
                                    <line x1="12" y1="15" x2="12" y2="3" />
                                </svg>
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
.data-table td { padding: 12px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.doc-name { display: flex; align-items: center; gap: 10px; }
.doc-icon { width: 28px; height: 28px; background: var(--green-pale); border-radius: 6px; display: flex; align-items: center; justify-content: center; color: var(--green-mid); flex-shrink: 0; }
.doc-icon svg { width: 14px; height: 14px; }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--type { background: rgba(53,53,53,0.08); color: var(--charcoal); }
.btn-dl { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 6px; border: 1.5px solid rgba(53,53,53,0.12); background: transparent; cursor: pointer; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; }
.btn-dl svg { width: 14px; height: 14px; }
.btn-dl:hover { border-color: var(--green-dark); color: var(--green-dark); }
</style>
