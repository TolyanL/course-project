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

      const payload = decodeTokenPayload(data.token)
      if (payload?.user_id) {
        teacherId.value = payload.user_id
      }
      if (payload?.login) {
        name.value = payload.login
        localStorage.setItem('name', payload.login)
      }

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

  function decodeTokenPayload(
    token: string,
  ): { user_id?: number; login?: string; role?: string } | null {
    try {
      const parts = token.split('.')
      if (parts.length !== 3 || !parts[1]) return null
      const payload = JSON.parse(atob(parts[1]))
      return payload
    } catch {
      return null
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
    localStorage.removeItem('name')
  }

  function initFromStorage() {
    const storedRole = localStorage.getItem('role') as 'admin' | 'teacher' | null
    if (storedRole) {
      role.value = storedRole
    }

    const storedName = localStorage.getItem('name')
    if (storedName) {
      name.value = storedName
    }

    const storedToken = localStorage.getItem('token')
    if (storedToken && !teacherId.value) {
      const payload = decodeTokenPayload(storedToken)
      if (payload?.user_id) {
        teacherId.value = payload.user_id
      }
      if (payload?.login && !storedName) {
        name.value = payload.login
        localStorage.setItem('name', payload.login)
      }
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
