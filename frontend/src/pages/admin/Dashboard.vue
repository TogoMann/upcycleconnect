<template>
    <div class="dashboard">
        <h1 class="page-title">Tableau de bord.</h1>

        <div class="kpi-grid">
            <div class="kpi-card">
                <div class="kpi-icon">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                        <circle cx="9" cy="7" r="4" />
                        <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                        <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                    </svg>
                </div>
                <div class="kpi-value">{{ adminStore.users.length }}</div>
                <div class="kpi-label">Utilisateurs</div>
            </div>

            <div class="kpi-card">
                <div class="kpi-icon">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                        <polyline points="14 2 14 8 20 8" />
                    </svg>
                </div>
                <div class="kpi-value">{{ adminStore.courses.length }}</div>
                <div class="kpi-label">Annonces</div>
            </div>

            <div class="kpi-card">
                <div class="kpi-icon">
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                        <line x1="16" y1="2" x2="16" y2="6" />
                        <line x1="8" y1="2" x2="8" y2="6" />
                        <line x1="3" y1="10" x2="21" y2="10" />
                    </svg>
                </div>
                <div class="kpi-value">{{ adminStore.events.length }}</div>
                <div class="kpi-label">Événements</div>
            </div>
        </div>

        <div class="audit-section">
            <h2 class="section-title">Audits & Rapports</h2>
            <div class="audit-card">
                <div class="audit-info">
                    <h3>Extraction des dépôts (PDF)</h3>
                    <p>Téléchargez le rapport du nombre total d'objets déposés sur une période donnée.</p>
                </div>
                <div class="audit-actions">
                    <div class="date-inputs">
                        <div class="input-group">
                            <label for="startDate">Date de début</label>
                            <input type="date" id="startDate" v-model="auditStartDate" />
                        </div>
                        <div class="input-group">
                            <label for="endDate">Date de fin</label>
                            <input type="date" id="endDate" v-model="auditEndDate" />
                        </div>
                    </div>
                    <button class="btn-download" @click="downloadAudit" :disabled="!auditStartDate || !auditEndDate || isDownloading">
                        <svg v-if="isDownloading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
                        </svg>
                        <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                            <polyline points="7 10 12 15 17 10" />
                            <line x1="12" y1="15" x2="12" y2="3" />
                        </svg>
                        {{ isDownloading ? 'Génération...' : 'Télécharger l\'audit' }}
                    </button>
                    <p v-if="auditError" class="error-msg">{{ auditError }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useAuthStore } from '@/stores/auth'
import { API_BASE } from '@/config'

const adminStore = useAdminStore()
const authStore = useAuthStore()

const auditStartDate = ref('')
const auditEndDate = ref('')
const isDownloading = ref(false)
const auditError = ref('')

const downloadAudit = async () => {
    if (!auditStartDate.value || !auditEndDate.value) return;

    isDownloading.value = true
    auditError.value = ''

    try {
        const response = await fetch(`${API_BASE}/reporting/audit/items/pdf?start=${auditStartDate.value}&end=${auditEndDate.value}`, {
            headers: {
                'Authorization': `Bearer ${authStore.token}`
            }
        })

        if (!response.ok) {
            const errText = await response.text()
            throw new Error(errText || 'Erreur lors de la génération du PDF')
        }

        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `audit_depots_${auditStartDate.value}_au_${auditEndDate.value}.pdf`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
    } catch (error: any) {
        console.error("Audit download error:", error)
        auditError.value = error.message
    } finally {
        isDownloading.value = false
    }
}

onMounted(() => {
    adminStore.fetchUsers()
    adminStore.fetchCourses()
    adminStore.fetchEvents()
})
</script>

<style scoped>
.dashboard {
    --green-dark: #086a35;
    --green-mid: #34895b;
    --green-light: #8bbd94;
    --green-pale: #d7ece1;
    --charcoal: #353535;
    --white: #ffffff;
    --red-error: #e53e3e;

    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
}

.page-title {
    font-size: clamp(2.2rem, 4vw, 3.2rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0 0 32px;
    line-height: 1.08;
}

.kpi-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 40px;
}

.kpi-card {
    background: var(--green-pale);
    border-radius: 14px;
    padding: 24px 22px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.kpi-icon {
    width: 40px;
    height: 40px;
    background: var(--green-mid);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--white);
    margin-bottom: 4px;
}

.kpi-value {
    font-size: 2.2rem;
    font-weight: 800;
    color: var(--green-dark);
    line-height: 1;
    letter-spacing: -0.03em;
}

.kpi-label {
    font-size: 0.85rem;
    color: var(--green-dark);
    font-weight: 500;
    opacity: 0.75;
}

/* Audit Section Styles */
.audit-section {
    margin-top: 40px;
}

.section-title {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 20px;
}

.audit-card {
    background: var(--white);
    border: 1px solid rgba(0, 0, 0, 0.08);
    border-radius: 14px;
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.audit-info h3 {
    margin: 0 0 8px 0;
    font-size: 1.1rem;
    color: var(--green-dark);
}

.audit-info p {
    margin: 0;
    font-size: 0.9rem;
    color: #666;
}

.audit-actions {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.date-inputs {
    display: flex;
    gap: 16px;
}

.input-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.input-group label {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--charcoal);
}

.input-group input {
    padding: 10px 12px;
    border: 1px solid #ccc;
    border-radius: 8px;
    font-family: inherit;
    font-size: 0.95rem;
}

.btn-download {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    padding: 12px 20px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s ease;
    align-self: flex-start;
}

.btn-download:hover:not(:disabled) {
    background: var(--green-mid);
}

.btn-download:disabled {
    background: #ccc;
    cursor: not-allowed;
}

.spinner {
    animation: spin 1s linear infinite;
    width: 18px;
    height: 18px;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

.error-msg {
    color: var(--red-error);
    font-size: 0.85rem;
    margin: 0;
}

@media (max-width: 600px) {
    .kpi-grid {
        grid-template-columns: 1fr;
    }
    
    .date-inputs {
        flex-direction: column;
    }
    
    .btn-download {
        align-self: stretch;
    }
}
</style>
