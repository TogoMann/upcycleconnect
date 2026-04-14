<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Notification {
    id: number
    titre: string
    message: string
    cible: string
    date: string
    envoyes: number
}

const historique = ref<Notification[]>([])
const form = ref({ titre: '', message: '', cible: 'tous' })
const loading = ref(false)
const success = ref(false)

onMounted(async () => {
    try {
        const res = await fetch('http://localhost:8081/admin/notifications', {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) historique.value = await res.json()
    } catch {}
})

async function envoyer() {
    if (!form.value.titre || !form.value.message) return
    loading.value = true
    success.value = false
    try {
        const res = await fetch('http://localhost:8081/admin/notifications', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify(form.value),
        })
        if (res.ok) {
            const created = await res.json()
            historique.value.unshift(created)
            form.value = { titre: '', message: '', cible: 'tous' }
            success.value = true
        }
    } catch {}
    loading.value = false
}
</script>

<template>
    <div class="notifications">
        <div class="page-header">
            <h1 class="page-title">Notifications.</h1>
            <p class="page-subtitle">Envoi push OneSignal et historique.</p>
        </div>

        <div class="send-card">
            <h3 class="card-title">Envoyer une notification</h3>
            <div v-if="success" class="alert alert--success">Notification envoyée.</div>

            <div class="form-group">
                <label class="form-label">Titre</label>
                <input v-model="form.titre" type="text" class="form-input" placeholder="Titre de la notification" />
            </div>
            <div class="form-group">
                <label class="form-label">Message</label>
                <textarea v-model="form.message" class="form-input form-textarea" rows="3" placeholder="Corps du message…"></textarea>
            </div>
            <div class="form-group">
                <label class="form-label">Cible</label>
                <select v-model="form.cible" class="form-input">
                    <option value="tous">Tous les utilisateurs</option>
                    <option value="client">Particuliers</option>
                    <option value="pro">Professionnels</option>
                    <option value="salarie">Salariés</option>
                </select>
            </div>
            <button class="btn-primary" :disabled="loading || !form.titre || !form.message" @click="envoyer">
                {{ loading ? 'Envoi…' : 'Envoyer la notification' }}
            </button>
        </div>

        <h3 class="section-title">Historique</h3>
        <div class="table-wrap">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Titre</th>
                        <th>Cible</th>
                        <th>Date</th>
                        <th>Envoyés</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="historique.length === 0">
                        <td colspan="4" class="empty">Aucune notification envoyée.</td>
                    </tr>
                    <tr v-for="n in historique" :key="n.id">
                        <td class="td-bold">{{ n.titre }}</td>
                        <td class="td-muted">{{ n.cible }}</td>
                        <td class="td-muted">{{ n.date }}</td>
                        <td>{{ n.envoyes }}</td>
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
.send-card { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); padding: 28px; margin-bottom: 32px; max-width: 560px; display: flex; flex-direction: column; gap: 16px; }
.card-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-label { font-size: 0.85rem; font-weight: 600; color: var(--charcoal); opacity: 0.75; }
.form-input { padding: 11px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea { resize: vertical; min-height: 80px; }
.btn-primary { padding: 12px 24px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.9rem; font-weight: 600; cursor: pointer; transition: background 0.2s; align-self: flex-start; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.5; cursor: default; }
.alert { padding: 12px 16px; border-radius: 8px; font-size: 0.88rem; font-weight: 500; background: var(--green-pale); color: var(--green-dark); }
.section-title { font-size: 1rem; font-weight: 700; color: var(--charcoal); margin: 0 0 16px; }
.table-wrap { background: var(--white); border-radius: 14px; border: 1.5px solid rgba(53,53,53,0.08); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { text-align: left; padding: 14px 20px; font-size: 0.8rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.06em; border-bottom: 1px solid rgba(53,53,53,0.08); }
.data-table td { padding: 14px 20px; font-size: 0.9rem; color: var(--charcoal); border-bottom: 1px solid rgba(53,53,53,0.05); }
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover { background: rgba(215,236,225,0.3); }
.td-bold { font-weight: 600; }
.td-muted { opacity: 0.55; font-size: 0.85rem; }
.empty { text-align: center; opacity: 0.4; padding: 40px !important; }
</style>
