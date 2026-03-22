<template>
  <div class="min-h-screen bg-gray-100 flex flex-col md:flex-row">
    <!-- Sidebar -->
    <aside class="w-full md:w-64 bg-gray-800 text-white flex flex-col">
      <div class="p-4 text-2xl font-bold border-b border-gray-700">
        Admin Panel
      </div>
      <nav class="flex-1 p-4 space-y-2">
        <router-link to="/admin" class="block p-2 rounded hover:bg-gray-700" exact-active-class="bg-gray-700">Dashboard</router-link>
        <router-link to="/admin/users" class="block p-2 rounded hover:bg-gray-700" active-class="bg-gray-700">Users</router-link>
        <router-link to="/admin/courses" class="block p-2 rounded hover:bg-gray-700" active-class="bg-gray-700">Courses</router-link>
        <router-link to="/admin/events" class="block p-2 rounded hover:bg-gray-700" active-class="bg-gray-700">Events</router-link>
        <router-link to="/" class="block p-2 rounded hover:bg-gray-700 mt-4 text-gray-400">Back to Site</router-link>
      </nav>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col min-w-0">
      <header class="bg-white shadow p-4 flex justify-between items-center">
        <h1 class="text-xl font-semibold text-gray-800">Back Office</h1>
        <div class="text-sm text-gray-600">Administrator</div>
      </header>
      <div class="p-6 overflow-auto flex-1">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
</script>
