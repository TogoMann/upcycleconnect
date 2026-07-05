<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()

interface Traduction {
    cle: string
    fr: string
    en: string
    modifie: boolean
}

const traductions = ref<Traduction[]>([])
const search = ref('')
const saving = ref<string | null>(null)

onMounted(async () => {
    try {
        const res = await fetch(`${API_BASE}/admin/langues`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) traductions.value = await res.json()
    } catch {}
})

const filtered = computed(() => traductions.value.filter(entry =>
    !search.value || entry.cle.toLowerCase().includes(search.value.toLowerCase()) || entry.fr.toLowerCase().includes(search.value.toLowerCase())
))

async function save(entry: Traduction) {
    saving.value = entry.cle
    try {
        await fetch(`${API_BASE}/admin/langues/${entry.cle}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ fr: entry.fr, en: entry.en }),
        })
        entry.modifie = false
    } catch {}
    saving.value = null
}
</script>

<template>
    <div class="langues">
        <div class="page-header">
            <h1 class="page-title">{{ t('admin.langues.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('admin.langues.subtitle') }}</p>
        </div>

        <div class="search-row">
            <input v-model="search" type="text" class="search-input" :placeholder="t('admin.langues.searchPlaceholder')" />
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>{{ t('admin.langues.colKey') }}</th>
                        <th>{{ t('admin.langues.colFrench') }}</th>
                        <th>{{ t('admin.langues.colEnglish') }}</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="filtered.length === 0">
                        <td colspan="4" class="empty">{{ t('admin.langues.empty') }}</td>
                    </tr>
                    <tr v-for="entry in filtered" :key="entry.cle" :class="{ 'row-modified': entry.modifie }">
                        <td class="td-key">{{ entry.cle }}</td>
                        <td>
                            <input
                                v-model="entry.fr"
                                type="text"
                                class="inline-input"
                                @input="entry.modifie = true"
                            />
                        </td>
                        <td>
                            <input
                                v-model="entry.en"
                                type="text"
                                class="inline-input"
                                @input="entry.modifie = true"
                            />
                        </td>
                        <td>
                            <button
                                v-if="entry.modifie"
                                class="btn-save"
                                :disabled="saving === entry.cle"
                                @click="save(entry)"
                            >
                                {{ saving === entry.cle ? '…' : t('admin.langues.save') }}
                            </button>
                        </td>
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
.search-row { margin-bottom: 16px; }
.search-input { width: 100%; max-width: 400px; padding: 10px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.search-input:focus { border-color: var(--green-mid); }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 10px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); vertical-align: middle; }
.data-table tr:last-child td { border-bottom: none; }
.row-modified { background: #fefce8; }
.td-key { font-family: 'Courier New', monospace; font-size: 0.8rem; color: var(--green-dark); opacity: 0.8; font-weight: 600; }
.inline-input { width: 100%; padding: 7px 10px; font-size: 0.88rem; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 6px; background: transparent; color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.inline-input:focus { border-color: var(--green-mid); background: var(--white); }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
.btn-save { padding: 5px 12px; border-radius: 6px; font-size: 0.78rem; font-weight: 600; cursor: pointer; background: var(--green-dark); color: var(--white); border: none; white-space: nowrap; }
.btn-save:disabled { opacity: 0.5; }
</style>
