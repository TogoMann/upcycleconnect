import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

const API_BASE = 'http://localhost:8081'

export const useAdminStore = defineStore('admin', () => {
  const authStore = useAuthStore()
  const users = ref<any[]>([])
  const courses = ref<any[]>([])
  const events = ref<any[]>([])

  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const fetchUsers = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/users`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch users')
      users.value = await res.json()
    } catch (err: any) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  const deleteUser = async (id: number) => {
    try {
      const res = await fetch(`${API_BASE}/users/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (res.ok) {
        users.value = users.value.filter(u => u.id !== id)
      }
    } catch (err: any) {
      error.value = err.message
    }
  }

  const updateUser = async (id: number, data: any) => {
    try {
      const res = await fetch(`${API_BASE}/users/${id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${authStore.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      })
      if (!res.ok) throw new Error('Failed to update user')
      
      const updatedUser = await res.json()
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        // If the backend returns a message instead of the full user object, 
        // we might need to fetch users again or update the local object manually.
        // The current backend handler returns `{"message": "user updated successfully"}`
        await fetchUsers()
      }
    } catch (err: any) {
      error.value = err.message
      throw err
    }
  }

  const getScoreHistory = async (id: number) => {
    try {
      const res = await fetch(`${API_BASE}/users/${id}/score/history`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch score history')
      return await res.json()
    } catch (err: any) {
      error.value = err.message
      throw err
    }
  }

  const fetchCourses = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/course`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch courses')
      courses.value = await res.json()
    } catch (err: any) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  const deleteCourse = async (id: number) => {
    try {
      const res = await fetch(`${API_BASE}/course/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (res.ok) {
        courses.value = courses.value.filter(c => c.id !== id)
      }
    } catch (err: any) {
      error.value = err.message
    }
  }

  const fetchEvents = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/event`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch events')
      events.value = await res.json()
    } catch (err: any) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  const deleteEvent = async (id: number) => {
    try {
      const res = await fetch(`${API_BASE}/event/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (res.ok) {
        events.value = events.value.filter(e => e.id !== id)
      }
    } catch (err: any) {
      error.value = err.message
    }
  }

  return {
    users,
    courses,
    events,
    isLoading,
    error,
    fetchUsers,
    deleteUser,
    updateUser,
    getScoreHistory,
    fetchCourses,
    deleteCourse,
    fetchEvents,
    deleteEvent
  }
})
