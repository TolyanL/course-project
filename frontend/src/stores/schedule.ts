import { defineStore } from 'pinia'
import { ref } from 'vue'
import { scheduleApi } from '@/api/client'
import type { ScheduleEntry, PairFormData } from '@/types'

export const useScheduleStore = defineStore('schedule', () => {
  const entries = ref<ScheduleEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSchedule(params?: {
    group_id?: number
    teacher_id?: number
    date?: string
    start_date?: string
    end_date?: string
  }) {
    loading.value = true
    error.value = null
    entries.value = []
    try {
      const { data } = await scheduleApi.getSchedule(params)
      entries.value = data || []
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  async function addEntry(data: PairFormData) {
    const { data: response } = await scheduleApi.createEntry(data)
    await fetchSchedule()
    return response.entry
  }

  async function updateEntry(id: number, data: PairFormData) {
    const { data: response } = await scheduleApi.updateEntry(id, data)
    await fetchSchedule()
    return response.entry
  }

  async function deleteEntry(id: number) {
    await scheduleApi.deleteEntry(id)
    await fetchSchedule()
  }

  return {
    entries,
    loading,
    error,
    fetchSchedule,
    addEntry,
    updateEntry,
    deleteEntry,
  }
})
