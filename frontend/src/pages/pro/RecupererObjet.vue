<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const code = ref('')
const result = ref<{ titre: string; categorie: string; poids: string } | null>(null)
const error = ref('')
const loading = ref(false)

async function recuperer() {
    if (!code.value.trim()) return
    loading.value = true
    error.value = ''
    result.value = null
    try {
        const res = await fetch('http://localhost:8081/pro/conteneurs/recuperer', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${authStore.token}`,
            },
            body: JSON.stringify({ code: code.value.trim() }),
        })
        if (res.ok) {
            result.value = await res.json()
            code.value = ''
        } else {
            const d = await res.json()
            error.value = d.message ?? 'Code introuvable.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
    loading.value = false
}
</script>

<template>
    <div class="recuperer">
        <div class="page-header">
            <h1 class="page-title">Récupérer un objet.</h1>
            <p class="page-subtitle">Saisissez ou scannez le code-barres du conteneur.</p>
        </div>

        <div class="scan-card">
            <div class="scan-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                </svg>
            </div>
            <form class="scan-form" @submit.prevent="recuperer">
                <input
                    v-model="code"
                    type="text"
                    class="scan-input"
                    placeholder="Code-barres (ex: UC-2024-00123)"
                    autofocus
                />
                <button type="submit" class="scan-btn" :disabled="loading || !code.trim()">
                    {{ loading ? 'Recherche…' : 'Récupérer' }}
                </button>
            </form>

            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <div v-if="result" class="result-card">
                <div class="result-title">Objet récupéré avec succès</div>
                <div class="result-row">
                    <span class="result-label">Titre</span>
                    <span class="result-value">{{ result.titre }}</span>
                </div>
                <div class="result-row">
                    <span class="result-label">Catégorie</span>
                    <span class="result-value">{{ result.categorie }}</span>
                </div>
                <div class="result-row">
                    <span class="result-label">Poids</span>
                    <span class="result-value">{{ result.poids }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 40px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.scan-card { max-width: 520px; background: var(--white); border-radius: 16px; border: 1.5px solid rgba(53,53,53,0.1); padding: 40px; display: flex; flex-direction: column; gap: 24px; }
.scan-icon { width: 64px; height: 64px; background: var(--green-pale); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: var(--green-dark); }
.scan-icon svg { width: 32px; height: 32px; }
.scan-form { display: flex; gap: 10px; }
.scan-input { flex: 1; padding: 12px 16px; font-size: 0.95rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.scan-input:focus { border-color: var(--green-mid); background: var(--white); }
.scan-btn { padding: 12px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; white-space: nowrap; }
.scan-btn:hover:not(:disabled) { background: var(--green-mid); }
.scan-btn:disabled { opacity: 0.5; cursor: default; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; }
.alert--error { background: #fee2e2; color: #991b1b; }
.result-card { background: var(--green-pale); border-radius: 12px; padding: 20px; display: flex; flex-direction: column; gap: 12px; }
.result-title { font-size: 0.88rem; font-weight: 700; color: var(--green-dark); text-transform: uppercase; letter-spacing: 0.06em; }
.result-row { display: flex; justify-content: space-between; align-items: center; }
.result-label { font-size: 0.85rem; color: var(--charcoal); opacity: 0.55; }
.result-value { font-size: 0.9rem; font-weight: 600; color: var(--charcoal); }
</style>
