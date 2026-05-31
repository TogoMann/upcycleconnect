<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement
} from 'chart.js'
import { Bar, Pie } from 'vue-chartjs'
import { useAuthStore } from '@/stores/auth'
import { API_BASE } from '@/config'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ArcElement)

const authStore = useAuthStore()

const predictions = ref<any[]>([])
const loading = ref(true)
const currentPage = ref(1)
const totalPages = ref(1)

const actorChartData = reactive({
  labels: [] as string[],
  datasets: [{
    backgroundColor: ['#41B883', '#E46651', '#00D8FF', '#DD1B16'],
    data: [] as number[]
  }]
})

const prestationChartData = reactive({
  labels: [] as string[],
  datasets: [{
    label: 'Nombre d\'inscriptions/achats',
    backgroundColor: '#3B82F6',
    data: [] as number[]
  }]
})

const fetchReporting = async () => {
  try {
    const token = localStorage.getItem('auth_token')
    const headers = { 'Authorization': `Bearer ${token}` }

    const [actorsRes, prestationsRes, predictionsRes] = await Promise.all([
      fetch(`${API_BASE}/reporting/actors`, { headers }),
      fetch(`${API_BASE}/reporting/prestations`, { headers }),
      fetch(`${API_BASE}/reporting/predictions?page=${currentPage.value}&limit=10`, { headers })
    ])

    if (!actorsRes.ok || !prestationsRes.ok || !predictionsRes.ok) {
        throw new Error("One or more requests failed")
    }

    const actors = await actorsRes.json()
    const prestations = await prestationsRes.json()
    const predictionsData = await predictionsRes.json()

    predictions.value = predictionsData.data
    totalPages.value = predictionsData.total_pages

    actorChartData.labels = actors.map((s: any) => s.role)
    actorChartData.datasets[0].data = actors.map((s: any) => s.count)

    prestationChartData.labels = prestations.map((s: any) => s.type)
    prestationChartData.datasets[0].data = prestations.map((s: any) => s.count)
    
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
    <h1 class="text-2xl font-bold mb-6">Dashboards & Data Mining</h1>

    <div v-if="loading" class="text-center py-10">
      Chargement des données...
    </div>

    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <div class="bg-white p-6 rounded-lg shadow">
          <h2 class="text-xl font-semibold mb-4">Répartition des Acteurs</h2>
          <div class="h-64">
            <Pie :data="actorChartData" :options="{ maintainAspectRatio: false }" />
          </div>
        </div>

        <div class="bg-white p-6 rounded-lg shadow">
          <h2 class="text-xl font-semibold mb-4">Succès des Prestations</h2>
          <div class="h-64">
            <Bar :data="prestationChartData" :options="{ maintainAspectRatio: false }" />
          </div>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow">
        <h2 class="text-xl font-semibold mb-4">Prédictions de Prochaines Prestations (Modèle de Machine Learning)</h2>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Utilisateur</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Prédiction</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Probabilité</th>
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
                Page 
                <input 
                  type="number" 
                  min="1" 
                  :max="totalPages" 
                  v-model.lazy="currentPage" 
                  @change="changePage(currentPage)" 
                  class="w-16 px-2 py-1 text-sm border border-gray-300 rounded-md text-center focus:outline-none focus:border-green-600 focus:ring-1 focus:ring-green-600" 
                />
                sur <span class="font-medium">{{ totalPages }}</span>
              </p>
            </div>
            <div>
              <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                <button
                  @click="changePage(currentPage - 1)"
                  :disabled="currentPage === 1"
                  class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                >
                  <span class="sr-only">Précédent</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
                  </svg>
                </button>
                <button
                  @click="changePage(currentPage + 1)"
                  :disabled="currentPage === totalPages"
                  class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                >
                  <span class="sr-only">Suivant</span>
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
