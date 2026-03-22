<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">Manage Events</h2>
      <button class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition" @click="adminStore.fetchEvents()">Refresh</button>
    </div>

    <div v-if="adminStore.isLoading" class="text-gray-500 py-4">Loading events...</div>
    <div v-if="adminStore.error" class="bg-red-100 text-red-700 p-4 rounded mb-4">{{ adminStore.error }}</div>

    <div class="bg-white rounded shadow overflow-hidden" v-if="!adminStore.isLoading">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-gray-50 text-gray-700 border-b">
            <th class="p-4 font-semibold">ID</th>
            <th class="p-4 font-semibold">Created By</th>
            <th class="p-4 font-semibold">Date</th>
            <th class="p-4 font-semibold">Price</th>
            <th class="p-4 font-semibold">Status</th>
            <th class="p-4 font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in adminStore.events" :key="event.id" class="border-b hover:bg-gray-50">
            <td class="p-4">{{ event.id }}</td>
            <td class="p-4">User {{ event.created_by }}</td>
            <td class="p-4">{{ new Date(event.date).toLocaleDateString() }}</td>
            <td class="p-4">{{ event.price }} €</td>
            <td class="p-4">
              <span v-if="event.approved" class="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full font-medium">Approved</span>
              <span v-else class="px-2 py-1 bg-yellow-100 text-yellow-800 text-xs rounded-full font-medium">Pending</span>
            </td>
            <td class="p-4 space-x-4">
              <button class="text-red-600 hover:text-red-800 font-medium" @click="deleteEvent(event.id)">Delete</button>
            </td>
          </tr>
          <tr v-if="adminStore.events.length === 0">
            <td colspan="6" class="p-8 text-center text-gray-500">No events found.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

const deleteEvent = (id: number) => {
  if (confirm('Are you sure you want to delete this event?')) {
    adminStore.deleteEvent(id)
  }
}

onMounted(() => {
  adminStore.fetchEvents()
})
</script>
