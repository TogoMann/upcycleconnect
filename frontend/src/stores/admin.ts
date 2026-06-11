import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import { API_BASE } from '@/config'

export const useAdminStore = defineStore('admin', () => {
  const authStore = useAuthStore()
  const users = ref<any[]>([])
  const usersCount = ref(0)
  const usersTotalPages = ref(1)
  const usersCurrentPage = ref(1)
  const courses = ref<any[]>([])
  const listingsCount = ref(0)
  const events = ref<any[]>([])
  const eventsCount = ref(0)
  const companies = ref<any[]>([])

  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const fetchCompanies = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/companies`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch companies')
      companies.value = await res.json()
    } catch (err: any) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  const fetchUsers = async (page = 1, limit = 10) => {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/users?page=${page}&limit=${limit}`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch users')
      const data = await res.json()
      users.value = data.data || []
      usersCount.value = data.total || 0
      usersTotalPages.value = data.total_pages || 1
      usersCurrentPage.value = data.page || 1
    } catch (err: any) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  const fetchListings = async () => {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE}/listing?page=1&limit=1`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch listings')
      const data = await res.json()
      listingsCount.value = data.total || 0
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

  const requestPasswordReset = async (userId: number) => {
    try {
      const res = await fetch(`${API_BASE}/auth/admin/reset-request`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authStore.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ user_id: userId })
      })
      if (!res.ok) {
        const errData = await res.text()
        throw new Error(errData || 'Failed to request password reset')
      }
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
      const res = await fetch(`${API_BASE}/admin/events`, {
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to fetch events')
      events.value = await res.json()
      eventsCount.value = events.value.length
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

  const approveEvent = async (id: number) => {
    try {
      const res = await fetch(`${API_BASE}/event/${id}/approve`, {
        method: 'PATCH',
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to approve event')
      const ev = events.value.find(e => e.id === id)
      if (ev) ev.approved = true
    } catch (err: any) {
      error.value = err.message
    }
  }

  const disapproveEvent = async (id: number) => {
    try {
      const res = await fetch(`${API_BASE}/event/${id}/disapprove`, {
        method: 'PATCH',
        headers: {
          'Authorization': `Bearer ${authStore.token}`
        }
      })
      if (!res.ok) throw new Error('Failed to disapprove event')
      const ev = events.value.find(e => e.id === id)
      if (ev) ev.approved = false
    } catch (err: any) {
      error.value = err.message
    }
  }

  return {
    users,
    usersCount,
    usersTotalPages,
    usersCurrentPage,
    courses,
    listingsCount,
    events,
    eventsCount,
    companies,
    isLoading,
    error,
    fetchUsers,
    fetchListings,
    fetchCompanies,
    deleteUser,

    updateUser,
    getScoreHistory,
    requestPasswordReset,
    fetchCourses,
    deleteCourse,
    fetchEvents,
    deleteEvent,
    approveEvent,
    disapproveEvent
  }
})
