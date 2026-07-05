<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()

interface Depot {
    id: number
    utilisateur: string
    objet: string
    date: string
    statut: string
    code_envoye: boolean
}

const depots = ref<Depot[]>([])

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/admin/depots`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) depots.value = await res.json()
    } catch {}
})

async function valider(d: Depot) {
    await fetch(`${API_BASE}/admin/depots/${d.id}/valider`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    d.statut = 'valide'
}

async function envoyerCode(d: Depot) {
    await fetch(`${API_BASE}/admin/depots/${d.id}/code`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${authStore.token}` },
    })
    d.code_envoye = true
}
</script>

<template>
    <div class="depots">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.depots.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.depots.subtitle') }}</p>
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.depots.colUser') }}</th>
                        <th>{{ t('admin.depots.colItem') }}</th>
                        <th>{{ t('admin.depots.colDate') }}</th>
                        <th>{{ t('admin.depots.colStatus') }}</th>
                        <th>{{ t('admin.depots.colCode') }}</th>
                        <th>{{ t('admin.depots.colActions') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="depots.length === 0">
                        <td colspan="6" class="empty">{{ t('admin.depots.empty') }}</td>
                    </tr>
                    <tr v-for="d in depots" :key="d.id">
                        <td class="td-bold">{{ d.utilisateur }}</td>
                        <td>{{ d.objet }}</td>
                        <td class="td-muted">{{ d.date }}</td>
                        <td>
                            <span class="badge" :class="d.statut === 'valide' ? 'badge--active' : 'badge--pending'">
                                {{ d.statut === 'valide' ? t('admin.depots.validated') : t('admin.depots.pending') }}
                            </span>
                        </td>
                        <td>
                            <span class="badge" :class="d.code_envoye ? 'badge--active' : 'badge--none'">
                                {{ d.code_envoye ? t('admin.depots.sent') : t('admin.depots.notSent') }}
                            </span>
                        </td>
                        <td class="td-actions">
                            <button
                                v-if="d.statut !== 'valide'"
                                class="btn-action"
                                @click="valider(d)"
                            >{{ t('admin.depots.validate') }}</button>
                            <button
                                v-if="d.statut === 'valide' && !d.code_envoye"
                                class="btn-action btn-action--send"
                                @click="envoyerCode(d)"
                            >{{ t('admin.depots.sendCode') }}</button>
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
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.td-actions { display: flex; gap: 8px; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge--active { background: var(--green-pale); color: var(--green-dark); }
.badge--pending { background: #fef3c7; color: #92400e; }
.badge--none { background: rgba(53,53,53,0.08); color: var(--charcoal); opacity: 0.6; }
.btn-action { padding: 6px 14px; border-radius: 6px; font-size: 0.8rem; font-weight: 600; cursor: pointer; border: 1.5px solid rgba(53,53,53,0.15); background: transparent; color: var(--charcoal); transition: border-color 0.2s, color 0.2s; white-space: nowrap; }
.btn-action:hover { border-color: var(--green-dark); color: var(--green-dark); }
.btn-action--send { border-color: var(--green-mid); color: var(--green-dark); }
.btn-action--send:hover { background: var(--green-pale); }
</style>
