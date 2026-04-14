<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const form = ref({
    nom_site: '',
    logo_url: '',
    email_contact: '',
    telephone: '',
    adresse: '',
    commission_taux: '',
    maintenance: false,
    inscription_ouverte: true,
})
const success = ref(false)
const error = ref('')
const loading = ref(false)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/parametres', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            const data = await res.json()
            form.value = { ...form.value, ...data }
        }
    } catch {}
})

async function save() {
    loading.value = true
    error.value = ''
    success.value = false
    try {
        const res = await fetch('http://localhost:8081/admin/parametres', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({ ...form.value, commission_taux: parseFloat(form.value.commission_taux) || 0 }),
        })
        if (res.ok) success.value = true
        else {
            const d = await res.json()
            error.value = d.message ?? 'Erreur.'
        }
    } catch {
        error.value = 'Erreur réseau.'
    }
    loading.value = false
}
</script>

<template>
    <div class="parametres">
        <div class="page-header">
            <h1 class="page-title">Paramètres.</h1>
            <p class="page-subtitle">Configuration globale de la plateforme.</p>
        </div>

        <form class="form-sections" @submit.prevent="save">
            <div v-if="success" class="alert alert--success">Paramètres enregistrés.</div>
            <div v-if="error" class="alert alert--error">{{ error }}</div>

            <div class="form-card">
                <h3 class="card-title">Identité</h3>
                <div class="form-group">
                    <label class="form-label">Nom du site</label>
                    <input v-model="form.nom_site" type="text" class="form-input" />
                </div>
                <div class="form-group">
                    <label class="form-label">URL du logo</label>
                    <input v-model="form.logo_url" type="text" class="form-input" placeholder="https://…" />
                </div>
            </div>

            <div class="form-card">
                <h3 class="card-title">Contact</h3>
                <div class="form-group">
                    <label class="form-label">Email de contact</label>
                    <input v-model="form.email_contact" type="email" class="form-input" />
                </div>
                <div class="form-group">
                    <label class="form-label">Téléphone</label>
                    <input v-model="form.telephone" type="tel" class="form-input" />
                </div>
                <div class="form-group">
                    <label class="form-label">Adresse</label>
                    <input v-model="form.adresse" type="text" class="form-input" />
                </div>
            </div>

            <div class="form-card">
                <h3 class="card-title">Commercial</h3>
                <div class="form-group">
                    <label class="form-label">Taux de commission (%)</label>
                    <input v-model="form.commission_taux" type="number" step="0.1" min="0" max="100" class="form-input" />
                </div>
            </div>

            <div class="form-card">
                <h3 class="card-title">Système</h3>
                <div class="toggle-row">
                    <div class="toggle-info">
                        <div class="toggle-title">Mode maintenance</div>
                        <div class="toggle-desc">Le site public affiche une page de maintenance.</div>
                    </div>
                    <button
                        type="button"
                        class="toggle-btn"
                        :class="{ 'toggle-btn--on': form.maintenance }"
                        @click="form.maintenance = !form.maintenance"
                    >
                        <span class="toggle-track"><span class="toggle-thumb"></span></span>
                    </button>
                </div>
                <div class="toggle-row">
                    <div class="toggle-info">
                        <div class="toggle-title">Inscriptions ouvertes</div>
                        <div class="toggle-desc">Autoriser les nouveaux comptes.</div>
                    </div>
                    <button
                        type="button"
                        class="toggle-btn"
                        :class="{ 'toggle-btn--on': form.inscription_ouverte }"
                        @click="form.inscription_ouverte = !form.inscription_ouverte"
                    >
                        <span class="toggle-track"><span class="toggle-thumb"></span></span>
                    </button>
                </div>
            </div>

            <button type="submit" class="btn-primary" :disabled="loading">
                {{ loading ? 'Enregistrement…' : 'Enregistrer les paramètres' }}
            </button>
        </form>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.form-sections { display: flex; flex-direction: column; gap: 20px; max-width: 600px; }
.form-card { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); padding: 24px; display: flex; flex-direction: column; gap: 16px; }
.card-title { font-size: 0.88rem; font-weight: 700; color: var(--charcoal); text-transform: uppercase; letter-spacing: 0.06em; opacity: 0.6; margin: 0; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.toggle-row { display: flex; justify-content: space-between; align-items: center; gap: 16px; }
.toggle-info { flex: 1; }
.toggle-title { font-size: 0.9rem; font-weight: 600; color: var(--charcoal); margin-bottom: 2px; }
.toggle-desc { font-size: 0.8rem; color: var(--charcoal); opacity: 0.5; }
.toggle-btn { background: none; border: none; cursor: pointer; padding: 4px; flex-shrink: 0; }
.toggle-track { display: flex; width: 44px; height: 24px; background: rgba(53,53,53,0.15); border-radius: 12px; padding: 2px; transition: background 0.2s; }
.toggle-thumb { width: 20px; height: 20px; background: var(--white); border-radius: 50%; box-shadow: 0 1px 3px rgba(0,0,0,0.2); transition: transform 0.2s; }
.toggle-btn--on .toggle-track { background: var(--green-dark); }
.toggle-btn--on .toggle-thumb { transform: translateX(20px); }
.btn-primary { padding: 13px 32px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; align-self: flex-start; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; }
.alert--success { background: var(--green-pale); color: var(--green-dark); }
.alert--error { background: #fee2e2; color: #991b1b; }
</style>
