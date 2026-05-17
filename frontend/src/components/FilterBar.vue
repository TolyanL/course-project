<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import Select from 'primevue/select'
import type { Group } from '@/types'

defineProps<{
  groups: Group[]
}>()

const emit = defineEmits<{
  filterChange: [{ start_date: string; end_date: string; group_id: number | null }]
}>()

interface WeekOption {
  label: string
  start_date: string
  end_date: string
  isCurrent: boolean
}

const selectedGroupId = ref<number | null>(null)
const selectedWeek = ref<WeekOption | null>(null)

const weeks = computed<WeekOption[]>(() => {
  const result: WeekOption[] = []
  const today = new Date()
  const todayStr = today.toISOString().split('T')[0] ?? ''
  
  for (let i = -5; i <= 5; i++) {
    const date = new Date(today)
    date.setDate(today.getDate() + i * 7)
    
    const day = date.getDay()
    const mondayDiff = day === 0 ? -6 : 1 - day
    
    const monday = new Date(date)
    monday.setDate(date.getDate() + mondayDiff)
    
    const sunday = new Date(monday)
    sunday.setDate(monday.getDate() + 6)
    
    const formatDate = (d: Date): string => {
      return d.toISOString().split('T')[0] ?? ''
    }
    
    const startStr = formatDate(monday)
    const endStr = formatDate(sunday)
    
    const dayOfMonth = monday.getDate()
    const monthNames = ['янв', 'фев', 'мар', 'апр', 'май', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек']
    const month = monthNames[monday.getMonth()]
    
    const weekNum = getWeekNumber(monday)
    const year = monday.getFullYear()
    
    result.push({
      label: `${dayOfMonth} ${month} (${weekNum} нед., ${year})`,
      start_date: startStr,
      end_date: endStr,
      isCurrent: startStr <= todayStr && todayStr <= endStr
    })
  }
  
  return result
})

function getWeekNumber(date: Date): number {
  const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
  const dayNum = d.getUTCDay() || 7
  d.setUTCDate(d.getUTCDate() + 4 - dayNum)
  const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1))
  return Math.ceil((((d.getTime() - yearStart.getTime()) / 86400000) + 1) / 7)
}

if (weeks.value.length > 5) {
  selectedWeek.value = weeks.value[5] ?? null
}

watch([selectedWeek, selectedGroupId], () => {
  if (!selectedWeek.value) return
  
  emit('filterChange', {
    start_date: selectedWeek.value.start_date,
    end_date: selectedWeek.value.end_date,
    group_id: selectedGroupId.value,
  })
}, { immediate: true })
</script>

<template>
  <div class="filter-bar p-3 mb-3 rounded">
    <div class="row g-3 align-items-end">
      <div class="col-auto">
        <label class="form-label small text-muted">Неделя</label>
        <Select
          v-model="selectedWeek"
          :options="weeks"
          optionLabel="label"
          placeholder="Выберите неделю"
          class="w-100"
        >
          <template #option="{ option }">
            <span :class="{ 'current-week': option.isCurrent }">
              {{ option.label }}{{ option.isCurrent ? ' (текущая)' : '' }}
            </span>
          </template>
        </Select>
      </div>
      <div class="col-auto">
        <label class="form-label small text-muted">Группа</label>
        <Select
          v-model="selectedGroupId"
          :options="groups"
          optionLabel="name"
          optionValue="id"
          placeholder="Все группы"
          class="w-100"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.filter-bar {
  background: #f8f9fa;
}

.current-week {
  font-weight: 600;
  color: #0d6efd;
}
</style>