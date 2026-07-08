<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

interface Project {
    id: { Int64: number; Valid: boolean }
    listing_id: { Int64: number; Valid: boolean }
    creator_id: { Int64: number; Valid: boolean }
    title: string
    description: string
    final_score: { Int32: number; Valid: boolean }
    status: string
    created_at: { Time: string; Valid: boolean }
    completed_at: { Time: string; Valid: boolean }
}

interface Step {
    id: { Int64: number; Valid: boolean }
    project_id: { Int64: number; Valid: boolean }
    step_number: number
    description: string
    created_at: { Time: string; Valid: boolean }
}

const project = ref<Project | null>(null)
const steps = ref<Step[]>([])
const loading = ref(true)
const editing = ref(false)
const editForm = ref({ title: '', description: '' })
const saving = ref(false)

const newStepDesc = ref('')
const addingStep = ref(false)
const editingStepId = ref<number | null>(null)
const editStepDesc = ref('')
const savingStep = ref(false)

const showStatusModal = ref(false)
const showDeleteModal = ref(false)

const projectId = computed(() => Number(route.params.id))

onMounted(async () => {
    await loadProject()
    loading.value = false
})

async function loadProject() {
    const token = authStore.token
    if (!token) return
    const headers = { Authorization: `Bearer ${token}` }

    const [projRes, stepsRes] = await Promise.all([
        fetch(`${API_BASE}/project/${projectId.value}`, { headers }).catch(() => null),
        fetch(`${API_BASE}/project/${projectId.value}/steps`, { headers }).catch(() => null),
    ])

    if (projRes?.ok) {
        project.value = await projRes.json()
    } else {
        router.push('/pro/projets')
        return
    }

    if (stepsRes?.ok) {
        const data = await stepsRes.json()
        steps.value = Array.isArray(data) ? data.sort((a: Step, b: Step) => a.step_number - b.step_number) : []
    }
}

function formatDate(d: { Time: string; Valid: boolean } | undefined) {
    if (!d?.Valid) return '—'
    return new Date(d.Time).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}

function statusConfig(s: string) {
    const map: Record<string, { label: string; class: string }> = {
        'in progress': { label: t('pro.projetDetail.statusInProgress'), class: 'badge--progress' },
        'done': { label: t('pro.projetDetail.statusDone'), class: 'badge--done' },
        'featured': { label: t('pro.projetDetail.statusFeatured'), class: 'badge--featured' },
        'cancelled': { label: t('pro.projetDetail.statusCancelled'), class: 'badge--cancelled' },
    }
    return map[s] || { label: s, class: 'badge--default' }
}

function startEdit() {
    if (!project.value) return
    editForm.value = { title: project.value.title, description: project.value.description }
    editing.value = true
}

async function saveEdit() {
    if (!editForm.value.title.trim()) return
    saving.value = true
    try {
        const res = await fetch(`${API_BASE}/project/${projectId.value}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({
                title: editForm.value.title.trim(),
                description: editForm.value.description.trim(),
                status: project.value?.status || 'in progress',
            }),
        })
        if (res.ok) {
            if (project.value) {
                project.value.title = editForm.value.title.trim()
                project.value.description = editForm.value.description.trim()
            }
            editing.value = false
        }
    } catch {}
    saving.value = false
}

async function changeStatus(newStatus: string) {
    try {
        await fetch(`${API_BASE}/project/${projectId.value}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({
                title: project.value?.title,
                description: project.value?.description,
                status: newStatus,
            }),
        })
        if (project.value) project.value.status = newStatus
    } catch {}
    showStatusModal.value = false
}

async function deleteProject() {
    try {
        const res = await fetch(`${API_BASE}/project/${projectId.value}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) router.push('/pro/projets')
    } catch {}
}

async function addStep() {
    if (!newStepDesc.value.trim()) return
    addingStep.value = true
    try {
        const nextNumber = steps.value.length > 0 ? Math.max(...steps.value.map(s => s.step_number)) + 1 : 1
        const res = await fetch(`${API_BASE}/project/steps`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({
                project_id: { Int64: projectId.value, Valid: true },
                step_number: nextNumber,
                description: newStepDesc.value.trim(),
            }),
        })
        if (res.ok) {
            newStepDesc.value = ''
            const stepsRes = await fetch(`${API_BASE}/project/${projectId.value}/steps`, {
                headers: { Authorization: `Bearer ${authStore.token}` },
            })
            if (stepsRes.ok) {
                const data = await stepsRes.json()
                steps.value = Array.isArray(data) ? data.sort((a: Step, b: Step) => a.step_number - b.step_number) : []
            }
        }
    } catch {}
    addingStep.value = false
}

function startEditStep(step: Step) {
    editingStepId.value = step.id.Int64
    editStepDesc.value = step.description
}

async function saveStep(step: Step) {
    if (!editStepDesc.value.trim()) return
    savingStep.value = true
    try {
        await fetch(`${API_BASE}/project/steps/${step.id.Int64}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${authStore.token}` },
            body: JSON.stringify({
                step_number: step.step_number,
                description: editStepDesc.value.trim(),
            }),
        })
        step.description = editStepDesc.value.trim()
        editingStepId.value = null
    } catch {}
    savingStep.value = false
}

async function deleteStep(stepId: number) {
    try {
        const res = await fetch(`${API_BASE}/project/steps/${stepId}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) {
            steps.value = steps.value.filter(s => s.id.Int64 !== stepId)
        }
    } catch {}
}
</script>

<template>
    <div class="projet-detail">
        <router-link to="/pro/projets" class="back-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
            {{ t('pro.projetDetail.backToProjects') }}
        </router-link>

        <div v-if="loading" class="loading-state">{{ t('pro.projetDetail.loading') }}</div>

        <template v-else-if="project">
            <!-- En-tête projet -->
            <div class="project-header-card">
                <div class="project-top">
                    <div class="project-title-area">
                        <template v-if="!editing">
                            <h1 class="page-title">{{ project.title }}</h1>
                            <p class="project-description">{{ project.description || t('pro.projetDetail.noDescription') }}</p>
                        </template>
                        <template v-else>
                            <input v-model="editForm.title" class="edit-title-input" :placeholder="t('pro.projetDetail.titlePlaceholder')" />
                            <textarea v-model="editForm.description" class="edit-desc-input" rows="3" :placeholder="t('pro.projetDetail.descriptionPlaceholder')"></textarea>
                            <div class="edit-actions">
                                <button class="btn-sm btn-sm--secondary" @click="editing = false">{{ t('pro.projetDetail.cancel') }}</button>
                                <button class="btn-sm btn-sm--primary" @click="saveEdit" :disabled="saving">{{ t('pro.projetDetail.save') }}</button>
                            </div>
                        </template>
                    </div>
                    <div class="project-actions">
                        <span class="badge badge-lg" :class="statusConfig(project.status).class">{{ statusConfig(project.status).label }}</span>
                    </div>
                </div>

                <!-- Info cards -->
                <div class="info-grid">
                    <div class="info-card">
                        <div class="info-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                        </div>
                        <div>
                            <div class="info-label">{{ t('pro.projetDetail.createdOn') }}</div>
                            <div class="info-value">{{ formatDate(project.created_at) }}</div>
                        </div>
                    </div>
                    <div class="info-card" v-if="project.completed_at?.Valid">
                        <div class="info-icon info-icon--green">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
                        </div>
                        <div>
                            <div class="info-label">{{ t('pro.projetDetail.completedOn') }}</div>
                            <div class="info-value">{{ formatDate(project.completed_at) }}</div>
                        </div>
                    </div>
                    <div class="info-card" v-if="project.final_score?.Valid">
                        <div class="info-icon info-icon--yellow">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
                        </div>
                        <div>
                            <div class="info-label">{{ t('pro.projetDetail.upcyclingScore') }}</div>
                            <div class="info-value">{{ project.final_score.Int32 }} pts</div>
                        </div>
                    </div>
                    <div class="info-card">
                        <div class="info-icon info-icon--blue">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                        </div>
                        <div>
                            <div class="info-label">{{ t('pro.projetDetail.steps') }}</div>
                            <div class="info-value">{{ steps.length }}</div>
                        </div>
                    </div>
                </div>

                <!-- Action buttons -->
                <div class="action-bar">
                    <button class="btn-action" @click="startEdit" v-if="!editing">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                        {{ t('pro.projetDetail.edit') }}
                    </button>
                    <button class="btn-action" @click="showStatusModal = true">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                        {{ t('pro.projetDetail.changeStatus') }}
                    </button>
                    <button class="btn-action btn-action--danger" @click="showDeleteModal = true">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/></svg>
                        {{ t('pro.projetDetail.delete') }}
                    </button>
                </div>
            </div>

            <!-- Section étapes -->
            <div class="steps-section">
                <div class="steps-header">
                    <h2 class="section-title">{{ t('pro.projetDetail.stepsTitle') }}</h2>
                    <span class="steps-count">{{ t('pro.projetDetail.stepsCount', { count: steps.length, plural: steps.length !== 1 ? 's' : '' }) }}</span>
                </div>

                <!-- Timeline des étapes -->
                <div class="steps-timeline" v-if="steps.length > 0">
                    <div v-for="step in steps" :key="step.id.Int64" class="step-item">
                        <div class="step-marker">
                            <div class="step-number">{{ step.step_number }}</div>
                            <div class="step-line"></div>
                        </div>
                        <div class="step-content">
                            <template v-if="editingStepId !== step.id.Int64">
                                <p class="step-desc">{{ step.description }}</p>
                                <div class="step-meta">
                                    <span class="step-date">{{ formatDate(step.created_at) }}</span>
                                    <div class="step-actions">
                                        <button class="step-btn" @click="startEditStep(step)" :title="t('pro.projetDetail.edit')">
                                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                                        </button>
                                        <button class="step-btn step-btn--danger" @click="deleteStep(step.id.Int64)" :title="t('pro.projetDetail.delete')">
                                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                                        </button>
                                    </div>
                                </div>
                            </template>
                            <template v-else>
                                <textarea v-model="editStepDesc" class="step-edit-input" rows="2"></textarea>
                                <div class="step-edit-actions">
                                    <button class="btn-sm btn-sm--secondary" @click="editingStepId = null">{{ t('pro.projetDetail.cancel') }}</button>
                                    <button class="btn-sm btn-sm--primary" @click="saveStep(step)" :disabled="savingStep">{{ t('pro.projetDetail.save') }}</button>
                                </div>
                            </template>
                        </div>
                    </div>
                </div>

                <div v-else class="steps-empty">
                    <p>{{ t('pro.projetDetail.noSteps') }}</p>
                </div>

                <!-- Ajouter une étape -->
                <div class="add-step-form">
                    <div class="add-step-marker">
                        <span>+</span>
                    </div>
                    <div class="add-step-content">
                        <textarea v-model="newStepDesc" class="form-input form-textarea-sm" :placeholder="t('pro.projetDetail.addStepPlaceholder')" rows="2"></textarea>
                        <button class="btn-primary btn-sm" @click="addStep" :disabled="addingStep || !newStepDesc.trim()">
                            {{ addingStep ? t('pro.projetDetail.adding') : t('pro.projetDetail.addStep') }}
                        </button>
                    </div>
                </div>
            </div>
        </template>

        <!-- Modal statut -->
        <Teleport to="body">
            <div v-if="showStatusModal" class="modal-overlay" @click.self="showStatusModal = false">
                <div class="modal-card">
                    <h3 class="modal-title">{{ t('pro.projetDetail.changeStatusTitle') }}</h3>
                    <div class="status-options">
                        <button v-for="s in ['in progress', 'done', 'cancelled']" :key="s"
                            class="status-option" :class="{ 'status-option--active': project?.status === s }"
                            @click="changeStatus(s)" :disabled="project?.status === s">
                            <span class="badge" :class="statusConfig(s).class">{{ statusConfig(s).label }}</span>
                            <span v-if="project?.status === s" class="status-current">{{ t('pro.projetDetail.currentStatus') }}</span>
                        </button>
                    </div>
                    <div class="modal-actions">
                        <button class="btn-secondary" @click="showStatusModal = false">{{ t('pro.projetDetail.close') }}</button>
                    </div>
                </div>
            </div>
        </Teleport>

        <!-- Modal suppression -->
        <Teleport to="body">
            <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
                <div class="modal-card">
                    <h3 class="modal-title modal-title--danger">{{ t('pro.projetDetail.deleteProjectTitle') }}</h3>
                    <p class="modal-text">{{ t('pro.projetDetail.deleteProjectText') }}</p>
                    <div class="modal-actions">
                        <button class="btn-secondary" @click="showDeleteModal = false">{{ t('pro.projetDetail.cancel') }}</button>
                        <button class="btn-danger" @click="deleteProject">{{ t('pro.projetDetail.deletePermanently') }}</button>
                    </div>
                </div>
            </div>
        </Teleport>
    </div>
</template>

<style scoped>
.back-link { display: inline-flex; align-items: center; gap: 6px; font-size: 0.85rem; color: var(--green-mid); text-decoration: none; margin-bottom: 24px; transition: color 0.2s; }
.back-link:hover { color: var(--green-dark); }
.back-link svg { width: 16px; height: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }

/* Project header card */
.project-header-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 16px; padding: 28px; margin-bottom: 28px; }
.project-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 20px; margin-bottom: 24px; }
.project-title-area { flex: 1; }
.page-title { font-size: clamp(1.6rem, 3vw, 2.2rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 10px; line-height: 1.15; }
.project-description { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; line-height: 1.6; margin: 0; }

.edit-title-input { width: 100%; padding: 10px 14px; font-size: 1.3rem; font-weight: 700; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; margin-bottom: 10px; }
.edit-title-input:focus { border-color: var(--green-mid); }
.edit-desc-input { width: 100%; padding: 10px 14px; font-size: 0.9rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; resize: vertical; line-height: 1.5; }
.edit-desc-input:focus { border-color: var(--green-mid); }
.edit-actions { display: flex; gap: 8px; margin-top: 10px; }

.badge-lg { font-size: 0.82rem; padding: 6px 16px; }
.badge { display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; white-space: nowrap; }
.badge--progress { background: #eff6ff; color: #1e40af; }
.badge--done { background: var(--green-pale); color: var(--green-dark); }
.badge--featured { background: #fef3c7; color: #92400e; }
.badge--cancelled { background: rgba(53,53,53,0.08); color: var(--charcoal); opacity: 0.6; }
.badge--default { background: rgba(53,53,53,0.08); color: var(--charcoal); }

/* Info grid */
.info-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 20px; }
.info-card { display: flex; align-items: center; gap: 12px; background: var(--cream); border-radius: 10px; padding: 14px 16px; }
.info-icon { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; background: rgba(53,53,53,0.06); flex-shrink: 0; }
.info-icon svg { width: 18px; height: 18px; color: var(--charcoal); opacity: 0.6; }
.info-icon--green { background: var(--green-pale); }
.info-icon--green svg { color: var(--green-dark); opacity: 1; }
.info-icon--yellow { background: #fef3c7; }
.info-icon--yellow svg { color: #f59e0b; opacity: 1; }
.info-icon--blue { background: #eff6ff; }
.info-icon--blue svg { color: #3b82f6; opacity: 1; }
.info-label { font-size: 0.72rem; font-weight: 600; color: var(--charcoal); opacity: 0.5; text-transform: uppercase; letter-spacing: 0.05em; }
.info-value { font-size: 0.92rem; font-weight: 700; color: var(--charcoal); margin-top: 2px; }

/* Action bar */
.action-bar { display: flex; gap: 10px; border-top: 1px solid rgba(53,53,53,0.06); padding-top: 18px; flex-wrap: wrap; }
.btn-action { display: inline-flex; align-items: center; gap: 8px; padding: 9px 16px; border: 1.5px solid rgba(53,53,53,0.12); border-radius: 8px; background: transparent; color: var(--charcoal); font-size: 0.84rem; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.btn-action:hover { border-color: var(--green-mid); color: var(--green-dark); }
.btn-action svg { width: 15px; height: 15px; }
.btn-action--danger { border-color: #fecaca; color: #991b1b; }
.btn-action--danger:hover { border-color: #dc2626; background: #fef2f2; }

/* Steps section */
.steps-section { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 16px; padding: 28px; }
.steps-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.section-title { font-size: 1.1rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.steps-count { font-size: 0.82rem; color: var(--charcoal); opacity: 0.45; font-weight: 500; }

/* Timeline */
.steps-timeline { display: flex; flex-direction: column; gap: 0; }
.step-item { display: flex; gap: 16px; }
.step-marker { display: flex; flex-direction: column; align-items: center; }
.step-number { width: 32px; height: 32px; border-radius: 50%; background: var(--green-dark); color: var(--white); display: flex; align-items: center; justify-content: center; font-size: 0.8rem; font-weight: 700; flex-shrink: 0; }
.step-line { width: 2px; flex: 1; background: rgba(53,53,53,0.1); min-height: 20px; }
.step-item:last-child .step-line { display: none; }
.step-content { flex: 1; padding-bottom: 24px; }
.step-desc { font-size: 0.9rem; color: var(--charcoal); line-height: 1.6; margin: 4px 0 8px; }
.step-meta { display: flex; justify-content: space-between; align-items: center; }
.step-date { font-size: 0.78rem; color: var(--charcoal); opacity: 0.4; }
.step-actions { display: flex; gap: 6px; }
.step-btn { width: 28px; height: 28px; border-radius: 6px; border: 1px solid rgba(53,53,53,0.1); background: transparent; cursor: pointer; display: flex; align-items: center; justify-content: center; color: var(--charcoal); opacity: 0.4; transition: all 0.2s; }
.step-btn:hover { opacity: 1; border-color: var(--green-mid); }
.step-btn svg { width: 14px; height: 14px; }
.step-btn--danger:hover { border-color: #dc2626; color: #dc2626; }
.step-edit-input { width: 100%; padding: 10px 12px; font-size: 0.88rem; border: 1.5px solid var(--green-mid); border-radius: 8px; background: var(--white); color: var(--charcoal); font-family: inherit; outline: none; resize: vertical; line-height: 1.5; margin: 4px 0 8px; }
.step-edit-actions { display: flex; gap: 8px; }

.steps-empty { padding: 32px; text-align: center; }
.steps-empty p { font-size: 0.88rem; color: var(--charcoal); opacity: 0.45; margin: 0; }

/* Add step form */
.add-step-form { display: flex; gap: 16px; margin-top: 20px; padding-top: 20px; border-top: 1px solid rgba(53,53,53,0.06); }
.add-step-marker { width: 32px; height: 32px; border-radius: 50%; border: 2px dashed rgba(53,53,53,0.2); display: flex; align-items: center; justify-content: center; font-size: 1rem; color: var(--charcoal); opacity: 0.3; flex-shrink: 0; }
.add-step-content { flex: 1; display: flex; flex-direction: column; gap: 10px; }
.form-input { padding: 10px 14px; font-size: 0.88rem; border: 1.5px solid rgba(53,53,53,0.15); border-radius: 8px; background: var(--cream); color: var(--charcoal); font-family: inherit; outline: none; transition: border-color 0.2s; }
.form-input:focus { border-color: var(--green-mid); background: var(--white); }
.form-textarea-sm { resize: vertical; min-height: 60px; line-height: 1.5; }

/* Buttons */
.btn-primary { padding: 10px 20px; background: var(--green-dark); color: var(--white); border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; align-self: flex-start; }
.btn-primary:hover:not(:disabled) { background: var(--green-mid); }
.btn-primary:disabled { opacity: 0.4; cursor: default; }
.btn-sm { padding: 7px 14px; font-size: 0.82rem; }
.btn-sm--primary { background: var(--green-dark); color: var(--white); border: none; border-radius: 6px; font-weight: 600; cursor: pointer; padding: 7px 14px; font-size: 0.82rem; }
.btn-sm--primary:hover { background: var(--green-mid); }
.btn-sm--secondary { background: transparent; border: 1.5px solid rgba(53,53,53,0.15); color: var(--charcoal); border-radius: 6px; font-weight: 600; cursor: pointer; padding: 7px 14px; font-size: 0.82rem; }
.btn-secondary { padding: 10px 20px; background: transparent; color: var(--charcoal); border: 1.5px solid rgba(53,53,53,0.2); border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; }
.btn-secondary:hover { border-color: var(--charcoal); }
.btn-danger { padding: 10px 20px; background: #dc2626; color: white; border: none; border-radius: 8px; font-size: 0.88rem; font-weight: 600; cursor: pointer; transition: background 0.2s; }
.btn-danger:hover { background: #b91c1c; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.modal-card { background: var(--white); border-radius: 20px; padding: 32px; max-width: 440px; width: 90%; box-shadow: 0 20px 60px rgba(0,0,0,0.15); }
.modal-title { font-size: 1.15rem; font-weight: 800; color: var(--charcoal); margin: 0 0 16px; }
.modal-title--danger { color: #dc2626; }
.modal-text { font-size: 0.88rem; color: var(--charcoal); opacity: 0.7; line-height: 1.5; margin: 0 0 20px; }
.modal-actions { display: flex; gap: 12px; justify-content: flex-end; margin-top: 20px; }

.status-options { display: flex; flex-direction: column; gap: 8px; }
.status-option { display: flex; align-items: center; justify-content: space-between; padding: 12px 16px; border: 1.5px solid rgba(53,53,53,0.1); border-radius: 10px; background: transparent; cursor: pointer; transition: border-color 0.2s; }
.status-option:hover:not(:disabled) { border-color: var(--green-mid); }
.status-option--active { border-color: var(--green-dark); background: var(--green-pale); }
.status-option:disabled { cursor: default; }
.status-current { font-size: 0.78rem; color: var(--green-dark); font-weight: 600; }

@media (max-width: 800px) {
    .info-grid { grid-template-columns: repeat(2, 1fr); }
    .project-top { flex-direction: column; }
}
@media (max-width: 500px) {
    .info-grid { grid-template-columns: 1fr; }
}
</style>
