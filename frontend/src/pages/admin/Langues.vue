<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

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
        const res = await fetch('http://localhost:8081/admin/langues', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) traductions.value = await res.json()
    } catch {}
})

const filtered = computed(() => traductions.value.filter(t =>
    !search.value || t.cle.toLowerCase().includes(search.value.toLowerCase()) || t.fr.toLowerCase().includes(search.value.toLowerCase())
))

async function save(t: Traduction) {
    saving.value = t.cle
    try {
        await fetch(`http://localhost:8081/admin/langues/${t.cle}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ fr: t.fr, en: t.en }),
        })
        t.modifie = false
    } catch {}
    saving.value = null
}
</script>

<template>
    <div class="langues">
        <div class="page-header">
            <h1 class="page-title">Langues.</h1>
            <p class="page-subtitle">Gérez les traductions de l'interface sans toucher au code.</p>
        </div>

        <div class="search-row">
            <input v-model="search" type="text" class="search-input" placeholder="Rechercher une clé ou un texte…" />
        </div>

        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Clé</th>
                        <th>Français</th>
                        <th>English</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="filtered.length === 0">
                        <td colspan="4" class="empty">Aucune traduction trouvée.</td>
                    </tr>
                    <tr v-for="t in filtered" :key="t.cle" :class="{ 'row-modified': t.modifie }">
                        <td class="td-key">{{ t.cle }}</td>
                        <td>
                            <input
                                v-model="t.fr"
                                type="text"
                                class="inline-input"
                                @input="t.modifie = true"
                            />
                        </td>
                        <td>
                            <input
                                v-model="t.en"
                                type="text"
                                class="inline-input"
                                @input="t.modifie = true"
                            />
                        </td>
                        <td>
                            <button
                                v-if="t.modifie"
                                class="btn-save"
                                :disabled="saving === t.cle"
                                @click="save(t)"
                            >
                                {{ saving === t.cle ? '…' : 'Sauv.' }}
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
