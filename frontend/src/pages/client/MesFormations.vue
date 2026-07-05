<script setup lang="ts">
import { API_BASE } from '@/config'
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()

interface CourseOrder {
    id: number
    course_id: number
    price: any
    booked_at: any
    course_name: string
}

interface CourseDocument {
    id: number
    filename: string
    original_name: string
}

const orders = ref<CourseOrder[]>([])
const documentsByCourse = ref<Record<number, CourseDocument[]>>({})
const loading = ref(true)
const error = ref('')

function rawVal(v: any): any {
    if (v === null || v === undefined) return null
    if (typeof v === 'object') return v.Valid ? (v.Time ?? v.Float64 ?? v.String) : null
    return v
}

function fmtDate(d: any): string {
    const val = rawVal(d)
    if (!val) return '—'
    return new Date(val).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}

function fmtPrice(p: any): string {
    const val = rawVal(p)
    const num = Number(val)
    return num > 0 ? `${num.toFixed(2)} €` : t('client.mesFormations.free')
}

async function fetchDocuments(courseId: number) {
    try {
        const res = await fetch(`${API_BASE}/course/${courseId}/documents`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (res.ok) documentsByCourse.value[courseId] = await res.json()
    } catch {}
}

onMounted(async () => {
    loading.value = true
    error.value = ''
    try {
        const res = await fetch(`${API_BASE}/course-order/me`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        })
        if (!res.ok) throw new Error(t('client.mesFormations.loadError'))
        orders.value = await res.json()
        await Promise.all(orders.value.map(o => fetchDocuments(o.course_id)))
    } catch (e: any) {
        error.value = e.message || t('client.mesFormations.loadError')
    } finally {
        loading.value = false
    }
})
</script>

<template>
    <div class="mes-formations">
        <div class="page-header">
            <h1 class="page-title">{{ t('client.mesFormations.pageTitle') }}</h1>
            <p class="page-subtitle">{{ t('client.mesFormations.subtitle') }}</p>
        </div>

        <div v-if="error" class="error-banner">{{ error }}</div>
        <div v-if="loading" class="loading-state">{{ t('client.mesFormations.loading') }}</div>

        <div v-else-if="orders.length === 0" class="empty-state">
            <p>{{ t('client.mesFormations.empty') }}</p>
            <router-link to="/particulier/catalogue" class="btn-primary">{{ t('client.mesFormations.viewCatalogue') }}</router-link>
        </div>

        <div v-else class="formations-list">
            <div v-for="o in orders" :key="o.id" class="formation-card">
                <div class="formation-header">
                    <h3 class="formation-name">{{ o.course_name }}</h3>
                    <span class="formation-price">{{ fmtPrice(o.price) }}</span>
                </div>
                <p class="formation-meta">{{ t('client.mesFormations.joinedOn', { date: fmtDate(o.booked_at) }) }}</p>

                <div class="documents-section">
                    <h4 class="documents-title">{{ t('client.mesFormations.documentsTitle') }}</h4>
                    <ul v-if="(documentsByCourse[o.course_id] || []).length > 0" class="documents-list">
                        <li v-for="doc in documentsByCourse[o.course_id]" :key="doc.id">
                            <a :href="`${API_BASE}/uploads/${doc.filename}`" target="_blank" rel="noopener">
                                📄 {{ doc.original_name }}
                            </a>
                        </li>
                    </ul>
                    <p v-else class="no-documents">{{ t('client.mesFormations.noDocuments') }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header { margin-bottom: 28px; }
.page-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); font-weight: 800; color: var(--charcoal); letter-spacing: -0.03em; margin: 0 0 8px; line-height: 1.08; }
.page-subtitle { font-size: 0.9rem; color: var(--charcoal); opacity: 0.6; margin: 0; }
.error-banner { background: rgba(229, 62, 62, 0.08); border: 1px solid rgba(229, 62, 62, 0.25); border-radius: 8px; padding: 12px 16px; font-size: 0.85rem; color: #c53030; margin-bottom: 16px; }
.loading-state { text-align: center; padding: 60px 0; opacity: 0.5; font-size: 0.9rem; }
.empty-state { text-align: center; padding: 60px 32px; display: flex; flex-direction: column; align-items: center; gap: 16px; opacity: 0.8; }
.btn-primary { padding: 11px 22px; background: var(--green-dark); color: var(--white); border-radius: 8px; font-size: 0.88rem; font-weight: 600; text-decoration: none; transition: background 0.2s; }
.btn-primary:hover { background: var(--green-mid); }
.formations-list { display: flex; flex-direction: column; gap: 16px; }
.formation-card { background: var(--white); border: 1.5px solid rgba(53,53,53,0.1); border-radius: 14px; padding: 22px; }
.formation-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.formation-name { font-size: 1.05rem; font-weight: 700; color: var(--charcoal); margin: 0; }
.formation-price { font-size: 0.9rem; font-weight: 800; color: var(--green-dark); flex-shrink: 0; }
.formation-meta { font-size: 0.8rem; color: var(--charcoal); opacity: 0.55; margin: 4px 0 16px; }
.documents-section { border-top: 1px solid rgba(53,53,53,0.08); padding-top: 14px; }
.documents-title { font-size: 0.82rem; font-weight: 700; color: var(--charcoal); opacity: 0.7; text-transform: uppercase; letter-spacing: 0.04em; margin: 0 0 10px; }
.documents-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 8px; }
.documents-list a { color: var(--green-dark); font-weight: 600; font-size: 0.88rem; text-decoration: none; }
.documents-list a:hover { text-decoration: underline; }
.no-documents { font-size: 0.85rem; color: var(--charcoal); opacity: 0.5; margin: 0; }
</style>
