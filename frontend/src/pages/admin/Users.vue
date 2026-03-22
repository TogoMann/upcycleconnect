<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold">Manage Users</h2>
      <button class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition" @click="adminStore.fetchUsers()">Refresh</button>
    </div>

    <div v-if="adminStore.isLoading" class="text-gray-500 py-4">Loading users...</div>
    <div v-if="adminStore.error" class="bg-red-100 text-red-700 p-4 rounded mb-4">{{ adminStore.error }}</div>

    <div class="bg-white rounded shadow overflow-hidden" v-if="!adminStore.isLoading">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-gray-50 text-gray-700 border-b">
            <th class="p-4 font-semibold">ID</th>
            <th class="p-4 font-semibold">Name</th>
            <th class="p-4 font-semibold">Email</th>
            <th class="p-4 font-semibold">Role</th>
            <th class="p-4 font-semibold">Score</th>
            <th class="p-4 font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in adminStore.users" :key="user.id" class="border-b hover:bg-gray-50">
            <td class="p-4">{{ user.id }}</td>
            <td class="p-4">{{ user.first_name }} {{ user.last_name }}</td>
            <td class="p-4">{{ user.email }}</td>
            <td class="p-4">
              <span class="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full font-medium">{{ user.role }}</span>
            </td>
            <td class="p-4">{{ user.score }}</td>
            <td class="p-4">
              <button class="text-red-600 hover:text-red-800 font-medium" @click="deleteUser(user.id)">Delete</button>
            </td>
          </tr>
          <tr v-if="adminStore.users.length === 0">
            <td colspan="6" class="p-8 text-center text-gray-500">No users found.</td>
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

const deleteUser = (id: number) => {
  if (confirm('Are you sure you want to delete this user?')) {
    adminStore.deleteUser(id)
  }
}

onMounted(() => {
  adminStore.fetchUsers()
})
</script>
