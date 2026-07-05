<script setup lang="ts">
import { reactive, ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useClientStore } from '@/stores/client'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const clientStore = useClientStore()

const ETATS = ['Neuf', 'Bon état', 'Abimé', 'Cassé'] as const
const MATERIAUX = ['Bois', 'Métal', 'Textile', 'Plastique', 'Verre', 'Céramique', 'Cuir', 'Autre']

const etatLabels = computed<Record<string, string>>(() => ({
    'Neuf': t('client.deposerObjet.stateNew'),
    'Bon état': t('client.deposerObjet.stateGood'),
    'Abimé': t('client.deposerObjet.stateDamaged'),
    'Cassé': t('client.deposerObjet.stateBroken'),
}))
const materiauLabels = computed<Record<string, string>>(() => ({
    'Bois': t('client.deposerObjet.materialWood'),
    'Métal': t('client.deposerObjet.materialMetal'),
    'Textile': t('client.deposerObjet.materialTextile'),
    'Plastique': t('client.deposerObjet.materialPlastic'),
    'Verre': t('client.deposerObjet.materialGlass'),
    'Céramique': t('client.deposerObjet.materialCeramic'),
    'Cuir': t('client.deposerObjet.materialLeather'),
    'Autre': t('client.deposerObjet.materialOther'),
}))

const form = reactive({
    material_type: '',
    physical_state: '',
    site_id: '',
    schedule: '',
    start: '',
    ending: '',
})

const errors = reactive({
    material_type: '',
    physical_state: '',
    schedule: '',
    start: '',
    ending: '',
    global: '',
})

const submitting = ref(false)
const success = ref(false)

function validate(): boolean {
    errors.material_type = form.material_type ? '' : t('client.deposerObjet.errorMaterialRequired')
    errors.physical_state = form.physical_state ? '' : t('client.deposerObjet.errorStateRequired')
    errors.schedule = form.schedule ? '' : t('client.deposerObjet.errorDateRequired')
    errors.start = form.start ? '' : t('client.deposerObjet.errorStartRequired')
    errors.ending = form.ending ? '' : t('client.deposerObjet.errorEndRequired')
    if (form.start && form.ending && form.start >= form.ending) {
        errors.ending = t('client.deposerObjet.errorEndAfterStart')
    }
    return !errors.material_type && !errors.physical_state && !errors.schedule && !errors.start && !errors.ending
}

function getTodayDate(): string {
    const d = new Date()
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
}

async function handleSubmit() {
    if (!validate()) return
    submitting.value = true
    errors.global = ''
    try {
        try {
            await clientStore.createItem({
                material_type: form.material_type,
                physical_state: form.physical_state,
                site_id: form.site_id ? Number(form.site_id) : undefined,
            })
        } catch {
            // Endpoint POST /items/ non encore implémenté côté backend
        }
        await clientStore.createEntry({
            schedule: form.schedule,
            start: form.start + ':00',
            ending: form.ending + ':00',
        })
        success.value = true
    } catch (e: any) {
        errors.global = e.message
    } finally {
        submitting.value = false
    }
}

function resetForm() {
    form.material_type = ''
    form.physical_state = ''
    form.site_id = ''
    form.schedule = ''
    form.start = ''
    form.ending = ''
    success.value = false
}

onMounted(() => {
    clientStore.fetchSites()
})
</script>

<template>
    <div class="page">
        <div class="page-header">
            <router-link to="/particulier/conteneurs" class="back-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="19" y1="12" x2="5" y2="12" />
                    <polyline points="12 19 5 12 12 5" />
                </svg>
                {{ t('client.deposerObjet.backToDeposits') }}
            </router-link>
            <h1 class="page-title">{{ t('client.deposerObjet.pageTitle') }}</h1>
        </div>

        <div v-if="success" class="success-card">
            <div class="success-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="20 6 9 17 4 12" />
                </svg>
            </div>
            <h2 class="success-title">{{ t('client.deposerObjet.successTitle') }}</h2>
            <p class="success-desc">
                {{ t('client.deposerObjet.successDesc') }}
            </p>
            <div class="success-actions">
                <router-link to="/particulier/planning" class="btn-primary">{{ t('client.deposerObjet.viewPlanning') }}</router-link>
                <button class="btn-secondary" @click="resetForm">{{ t('client.deposerObjet.depositAnother') }}</button>
            </div>
        </div>

        <template v-else>
            <div class="info-banner">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <circle cx="12" cy="12" r="10" />
                    <line x1="12" y1="8" x2="12" y2="12" />
                    <line x1="12" y1="16" x2="12.01" y2="16" />
                </svg>
                <p>
                    {{ t('client.deposerObjet.infoBanner') }}
                </p>
            </div>

            <form class="form-card" @submit.prevent="handleSubmit">
                <div class="form-section-title">{{ t('client.deposerObjet.itemSection') }}</div>

                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">{{ t('client.deposerObjet.materialType') }}</label>
                        <select
                            v-model="form.material_type"
                            class="form-input"
                            :class="{ 'form-input--error': errors.material_type }"
                        >
                            <option value="" disabled>{{ t('client.deposerObjet.select') }}</option>
                            <option v-for="m in MATERIAUX" :key="m" :value="m">{{ materiauLabels[m] }}</option>
                        </select>
                        <span v-if="errors.material_type" class="form-error">{{ errors.material_type }}</span>
                    </div>

                    <div class="form-group">
                        <label class="form-label">{{ t('client.deposerObjet.physicalState') }}</label>
                        <select
                            v-model="form.physical_state"
                            class="form-input"
                            :class="{ 'form-input--error': errors.physical_state }"
                        >
                            <option value="" disabled>{{ t('client.deposerObjet.select') }}</option>
                            <option v-for="e in ETATS" :key="e" :value="e">{{ etatLabels[e] }}</option>
                        </select>
                        <span v-if="errors.physical_state" class="form-error">{{ errors.physical_state }}</span>
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">{{ t('client.deposerObjet.depositLocation') }}</label>
                    <select v-model="form.site_id" class="form-input">
                        <option value="">{{ t('client.deposerObjet.selectSiteOptional') }}</option>
                        <option v-for="site in clientStore.sites" :key="site.id?.Int64" :value="site.id?.Int64">
                            {{ t('client.deposerObjet.site', { id: site.id?.Int64, type: site.type_site || t('client.deposerObjet.collectionPoint') }) }}
                        </option>
                    </select>
                </div>

                <div class="form-section-title form-section-title--mt">{{ t('client.deposerObjet.slotSection') }}</div>

                <div class="form-group">
                    <label class="form-label">{{ t('client.deposerObjet.depositDate') }}</label>
                    <input
                        v-model="form.schedule"
                        type="date"
                        class="form-input"
                        :class="{ 'form-input--error': errors.schedule }"
                        :min="getTodayDate()"
                    />
                    <span v-if="errors.schedule" class="form-error">{{ errors.schedule }}</span>
                </div>

                <div class="form-row">
                    <div class="form-group">
                        <label class="form-label">{{ t('client.deposerObjet.startTime') }}</label>
                        <input
                            v-model="form.start"
                            type="time"
                            class="form-input"
                            :class="{ 'form-input--error': errors.start }"
                        />
                        <span v-if="errors.start" class="form-error">{{ errors.start }}</span>
                    </div>

                    <div class="form-group">
                        <label class="form-label">{{ t('client.deposerObjet.endTime') }}</label>
                        <input
                            v-model="form.ending"
                            type="time"
                            class="form-input"
                            :class="{ 'form-input--error': errors.ending }"
                        />
                        <span v-if="errors.ending" class="form-error">{{ errors.ending }}</span>
                    </div>
                </div>

                <div v-if="errors.global" class="error-banner">{{ errors.global }}</div>

                <div class="form-actions">
                    <router-link to="/particulier/conteneurs" class="btn-cancel">{{ t('client.deposerObjet.cancel') }}</router-link>
                    <button type="submit" class="btn-submit" :disabled="submitting">
                        {{ submitting ? t('client.deposerObjet.saving') : t('client.deposerObjet.validateDeposit') }}
                    </button>
                </div>
            </form>
        </template>
    </div>
</template>

<style scoped>
.page {
    font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
    color: var(--charcoal);
    max-width: 600px;
}

.page-header {
    margin-bottom: 28px;
}
.back-link {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.83rem;
    color: var(--green-mid);
    text-decoration: none;
    font-weight: 500;
    margin-bottom: 16px;
    transition: color 0.2s;
}
.back-link:hover { color: var(--green-dark); }
.back-link svg { width: 16px; height: 16px; }
.page-title {
    font-size: clamp(1.8rem, 3.5vw, 2.6rem);
    font-weight: 800;
    color: var(--charcoal);
    letter-spacing: -0.03em;
    margin: 0;
    line-height: 1.08;
}

.info-banner {
    display: flex;
    gap: 12px;
    background: var(--green-pale);
    border-radius: 10px;
    padding: 14px 18px;
    margin-bottom: 24px;
    align-items: flex-start;
}
.info-banner svg {
    width: 18px;
    height: 18px;
    color: var(--green-mid);
    flex-shrink: 0;
    margin-top: 1px;
}
.info-banner p {
    font-size: 0.83rem;
    color: var(--charcoal);
    opacity: 0.8;
    margin: 0;
    line-height: 1.55;
}

.form-card {
    background: var(--white);
    border: 1.5px solid rgba(53, 53, 53, 0.1);
    border-radius: 14px;
    padding: 28px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}
.form-section-title {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--green-dark);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    padding-bottom: 4px;
    border-bottom: 1.5px solid var(--green-pale);
}
.form-section-title--mt {
    margin-top: 4px;
}
.form-group {
    display: flex;
    flex-direction: column;
    gap: 7px;
}
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
}
.form-label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--charcoal);
}
.form-input {
    width: 100%;
    padding: 13px 16px;
    font-size: 0.9rem;
    font-family: inherit;
    color: var(--charcoal);
    background: var(--cream);
    border: 1.5px solid rgba(53, 53, 53, 0.25);
    border-radius: 8px;
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
    appearance: none;
}
.form-input:focus {
    border-color: var(--green-mid);
    box-shadow: 0 0 0 3px rgba(52, 137, 91, 0.1);
}
.form-input--error { border-color: #e53e3e; }
.form-error {
    font-size: 0.78rem;
    color: #e53e3e;
    font-weight: 500;
}
.error-banner {
    background: rgba(229, 62, 62, 0.08);
    border: 1px solid rgba(229, 62, 62, 0.25);
    border-radius: 8px;
    padding: 12px 16px;
    font-size: 0.83rem;
    color: #e53e3e;
}
.form-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    padding-top: 8px;
    border-top: 1px solid rgba(53, 53, 53, 0.08);
}
.btn-cancel {
    padding: 11px 20px;
    background: transparent;
    border: 1.5px solid rgba(53, 53, 53, 0.2);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    color: rgba(53, 53, 53, 0.6);
    text-decoration: none;
    transition: border-color 0.2s, color 0.2s;
}
.btn-cancel:hover {
    border-color: rgba(53, 53, 53, 0.4);
    color: var(--charcoal);
}
.btn-submit {
    padding: 11px 24px;
    background: var(--green-dark);
    color: var(--white);
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s;
}
.btn-submit:hover:not(:disabled) { background: var(--green-mid); }
.btn-submit:disabled {
    opacity: 0.65;
    cursor: not-allowed;
}

.success-card {
    background: var(--green-pale);
    border-radius: 16px;
    padding: 40px 32px;
    text-align: center;
}
.success-icon {
    width: 56px;
    height: 56px;
    background: var(--green-mid);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 20px;
    color: var(--white);
}
.success-icon svg { width: 24px; height: 24px; }
.success-title {
    font-size: 1.4rem;
    font-weight: 800;
    color: var(--green-dark);
    margin: 0 0 12px;
    letter-spacing: -0.02em;
}
.success-desc {
    font-size: 0.875rem;
    color: var(--charcoal);
    opacity: 0.75;
    line-height: 1.65;
    margin: 0 0 28px;
    max-width: 380px;
    margin-left: auto;
    margin-right: auto;
}
.success-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    flex-wrap: wrap;
}
.btn-primary {
    display: inline-block;
    padding: 11px 22px;
    background: var(--green-dark);
    color: var(--white);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    text-decoration: none;
    transition: background 0.2s;
}
.btn-primary:hover { background: var(--green-mid); }
.btn-secondary {
    padding: 11px 22px;
    background: transparent;
    border: 1.5px solid var(--green-mid);
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--green-mid);
    cursor: pointer;
    font-family: inherit;
    transition: background 0.2s, color 0.2s;
}
.btn-secondary:hover {
    background: var(--green-dark);
    color: var(--white);
    border-color: var(--green-dark);
}

@media (max-width: 480px) {
    .form-row {
        grid-template-columns: 1fr;
    }
}
</style>
