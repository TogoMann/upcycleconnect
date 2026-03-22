<template>
  <div>
    <h2 class="text-2xl font-bold mb-4">Dashboard</h2>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white p-6 rounded shadow border-t-4 border-blue-500">
        <div class="text-gray-500 text-sm font-semibold uppercase tracking-wider">Total Users</div>
        <div class="text-3xl font-bold text-gray-800 mt-2">{{ adminStore.users.length }}</div>
      </div>
      <div class="bg-white p-6 rounded shadow border-t-4 border-green-500">
        <div class="text-gray-500 text-sm font-semibold uppercase tracking-wider">Total Courses</div>
        <div class="text-3xl font-bold text-gray-800 mt-2">{{ adminStore.courses.length }}</div>
      </div>
      <div class="bg-white p-6 rounded shadow border-t-4 border-purple-500">
        <div class="text-gray-500 text-sm font-semibold uppercase tracking-wider">Total Events</div>
        <div class="text-3xl font-bold text-gray-800 mt-2">{{ adminStore.events.length }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

onMounted(() => {
  adminStore.fetchUsers()
  adminStore.fetchCourses()
  adminStore.fetchEvents()
})
</script>
