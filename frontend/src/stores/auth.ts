import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/client'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const role = ref<'admin' | 'teacher' | null>(null)
  const teacherId = ref<number | null>(null)
  const name = ref<string | null>(null)
  const lastError = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => role.value === 'admin')
  const isTeacher = computed(() => role.value === 'teacher')

  async function login(
    login: string,
    password: string,
  ): Promise<{ success: boolean; error?: string }> {
    lastError.value = null

    try {
      const { data } = await authApi.login(login, password)
      token.value = data.token
      role.value = data.role
      localStorage.setItem('token', data.token)
      localStorage.setItem('role', data.role)
      return { success: true }
    } catch (e) {
      let errorMessage = 'Ошибка входа'
      if (axios.isAxiosError(e)) {
        if (e.response?.data?.error) {
          errorMessage = e.response.data.error
        } else if (e.response?.status === 401) {
          errorMessage = 'Неверные учетные данные'
        }
      }
      lastError.value = errorMessage
      return { success: false, error: errorMessage }
    }
  }

  function setUserData(data: { role: 'admin' | 'teacher'; teacherId?: number; name?: string }) {
    role.value = data.role
    if (data.teacherId) teacherId.value = data.teacherId
    if (data.name) name.value = data.name
  }

  function logout() {
    token.value = null
    role.value = null
    teacherId.value = null
    name.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('role')
  }

  function initFromStorage() {
    const storedRole = localStorage.getItem('role') as 'admin' | 'teacher' | null
    if (storedRole) {
      role.value = storedRole
    }
  }

  return {
    token,
    role,
    teacherId,
    name,
    lastError,
    isAuthenticated,
    isAdmin,
    isTeacher,
    login,
    logout,
    setUserData,
    initFromStorage,
  }
})
