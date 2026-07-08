<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement,
  PointElement,
  LineElement
} from 'chart.js'
import { Bar, Pie, Doughnut } from 'vue-chartjs'
import { useAuthStore } from '@/stores/auth'
import { API_BASE } from '@/config'
import { useI18n } from 'vue-i18n'

ChartJS.register(
  Title, 
  Tooltip, 
  Legend, 
  BarElement, 
  CategoryScale, 
  LinearScale, 
  ArcElement, 
  PointElement, 
  LineElement
)

const { t } = useI18n()
const authStore = useAuthStore()

const predictions = ref<any[]>([])
const loading = ref(true)
const currentPage = ref(1)
const totalPages = ref(1)
const mlStatus = ref<any>(null)

const actorStats = ref<any[]>([])
const prestationStats = ref<any[]>([])
const mlDistStats = ref<Record<string, number>>({})

const actorChartData = computed(() => ({
  labels: actorStats.value.map((s: any) => s.role),
  datasets: [{
    backgroundColor: ['#41B883', '#E46651', '#00D8FF', '#DD1B16'],
    data: actorStats.value.map((s: any) => s.count)
  }]
}))

const prestationChartData = computed(() => ({
  labels: prestationStats.value.map((s: any) => s.type),
  datasets: [{
    label: t('admin.reporting.registrationsCount'),
    backgroundColor: '#3B82F6',
    data: prestationStats.value.map((s: any) => s.count)
  }]
}))

const mlChartData = computed(() => ({
  labels: Object.keys(mlDistStats.value),
  datasets: [{
    backgroundColor: ['#6366F1', '#F59E0B', '#10B981'],
    data: Object.values(mlDistStats.value)
  }]
}))

const fetchReporting = async () => {
  try {
    const token = localStorage.getItem('auth_token')
    const headers = { 'Authorization': `Bearer ${token}` }

    const [actorsRes, prestationsRes, predictionsRes, mlStatusRes, mlDistRes] = await Promise.all([
      fetch(`${API_BASE}/reporting/actors`, { headers }),
      fetch(`${API_BASE}/reporting/prestations`, { headers }),
      fetch(`${API_BASE}/reporting/predictions?page=${currentPage.value}&limit=10`, { headers }),
      fetch(`${API_BASE}/reporting/ml-status`, { headers }),
      fetch(`${API_BASE}/reporting/predictions/distribution`, { headers })
    ])

    if (!actorsRes.ok || !prestationsRes.ok || !predictionsRes.ok || !mlStatusRes.ok || !mlDistRes.ok) {
        throw new Error("One or more requests failed")
    }

    actorStats.value = await actorsRes.json()
    prestationStats.value = await prestationsRes.json()
    const predictionsData = await predictionsRes.json()
    mlStatus.value = await mlStatusRes.json()
    mlDistStats.value = await mlDistRes.json()

    predictions.value = predictionsData.data
    totalPages.value = predictionsData.total_pages
    
  } catch (error) {
    console.error('Failed to fetch reporting data', error)
  } finally {
    loading.value = false
  }
}

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchReporting()
  }
}

onMounted(fetchReporting)
</script>

<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">{{ t('admin.reporting.pageTitle') }}</h1>

    <div v-if="loading" class="text-center py-10">
      {{ t('admin.reporting.loading') }}
    </div>

    <div v-else>
      <div class="bg-white p-4 rounded-lg shadow mb-8 border-l-4 border-indigo-500 flex justify-between items-center">
        <div>
          <h2 class="text-sm font-bold text-indigo-600 uppercase tracking-wider">{{ t('admin.reporting.aiStatus') }}</h2>
          <p class="text-gray-600 text-sm" v-if="mlStatus">
            {{ t('admin.reporting.lastTraining', { date: new Date(mlStatus.last_run).toLocaleString() }) }}
          </p>
        </div>
        <div class="text-right">
          <p class="text-2xl font-bold text-indigo-600">{{ mlStatus?.total_predictions || 0 }}</p>
          <p class="text-xs text-gray-500">{{ t('admin.reporting.profilesAnalyzed') }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
        <div class="bg-white p-6 rounded-lg shadow">
          <h2 class="text-xl font-semibold mb-4">{{ t('admin.reporting.actorsBreakdown') }}</h2>
          <div class="h-64 flex items-center justify-center">
            <Pie v-if="actorStats.length > 0" :data="actorChartData" :options="{ maintainAspectRatio: false }" />
            <p v-else class="text-gray-400 text-sm italic">{{ t('admin.reporting.noDataAvailable') }}</p>
          </div>
        </div>

        <div class="bg-white p-6 rounded-lg shadow">
          <h2 class="text-xl font-semibold mb-4">{{ t('admin.reporting.prestationsSuccess') }}</h2>
          <div class="h-64 flex items-center justify-center">
            <Bar v-if="prestationStats.length > 0" :data="prestationChartData" :options="{ maintainAspectRatio: false }" />
            <p v-else class="text-gray-400 text-sm italic">{{ t('admin.reporting.noDataAvailable') }}</p>
          </div>
        </div>

        <div class="bg-white p-6 rounded-lg shadow">
          <h2 class="text-xl font-semibold mb-4">{{ t('admin.reporting.predictiveTrends') }}</h2>
          <div class="h-64 flex items-center justify-center">
            <Doughnut v-if="Object.keys(mlDistStats).length > 0" :data="mlChartData" :options="{ maintainAspectRatio: false }" />
            <p v-else class="text-gray-400 text-sm italic">{{ t('admin.reporting.computingPredictions') }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h2 class="text-xl font-semibold mb-4">{{ t('admin.reporting.predictionsTitle') }}</h2>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ t('admin.reporting.colUser') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ t('admin.reporting.colEmail') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ t('admin.reporting.colPrediction') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ t('admin.reporting.colProbability') }}</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="pred in predictions" :key="pred.user_id">
                <td class="px-6 py-4 whitespace-nowrap">{{ pred.username }}</td>
                <td class="px-6 py-4 whitespace-nowrap">{{ pred.email }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                    {{ pred.predicted_service_type }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">{{ (pred.probability * 100).toFixed(1) }}%</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="mt-4 flex items-center justify-between border-t border-gray-200 px-4 py-3 sm:px-6">
          <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
            <div class="flex items-center gap-3">
              <p class="text-sm text-gray-700 flex items-center gap-2">
                {{ t('admin.reporting.page') }}
                <input
                  type="number"
                  min="1"
                  :max="totalPages"
                  v-model.lazy="currentPage"
                  @change="changePage(currentPage)"
                  class="w-16 px-2 py-1 text-sm border border-gray-300 rounded-md text-center focus:outline-none focus:border-green-600 focus:ring-1 focus:ring-green-600"
                />
                {{ t('admin.reporting.of') }} <span class="font-medium">{{ totalPages }}</span>
              </p>
            </div>
            <div>
              <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                <button
                  @click="changePage(currentPage - 1)"
                  :disabled="currentPage === 1"
                  class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                >
                  <span class="sr-only">{{ t('admin.reporting.previous') }}</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
                  </svg>
                </button>
                <button
                  @click="changePage(currentPage + 1)"
                  :disabled="currentPage === totalPages"
                  class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                >
                  <span class="sr-only">{{ t('admin.reporting.next') }}</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
                  </svg>
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
