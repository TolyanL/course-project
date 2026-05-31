<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useScheduleStore } from '@/stores/schedule'
import { groupsApi, subjectsApi, classroomsApi, teachersApi } from '@/api/client'
import ScheduleTable from '@/components/ScheduleTable.vue'
import PairEditor from '@/components/PairEditor.vue'
import Select from 'primevue/select'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { Group, Subject, Classroom, Teacher, Pair, PairFormData, ConflictResponse } from '@/types'
import axios from 'axios'

const confirm = useConfirm()

const scheduleStore = useScheduleStore()

const groups = ref<Group[]>([])
const subjects = ref<Subject[]>([])
const classrooms = ref<Classroom[]>([])
const teachers = ref<Teacher[]>([])

const selectedGroupId = ref<number | null>(null)

interface WeekOption {
  label: string
  start_date: string
  end_date: string
  isCurrent: boolean
}

const selectedWeek = ref<WeekOption | null>(null)

const weeks = computed<WeekOption[]>(() => {
  const result: WeekOption[] = []
  const today = new Date()
  const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`

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
      const y = d.getFullYear()
      const m = String(d.getMonth() + 1).padStart(2, '0')
      const day = String(d.getDate()).padStart(2, '0')
      return `${y}-${m}-${day}`
    }
    const startStr = formatDate(monday)
    const endStr = formatDate(sunday)
    const monthNames = ['янв', 'фев', 'мар', 'апр', 'май', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек']
    const dayOfMonth = monday.getDate()
    const month = monthNames[monday.getMonth()]
    const weekNum = getWeekNumber(monday)
    const year = monday.getFullYear()
    result.push({
      label: `${dayOfMonth} ${month} (${weekNum} нед., ${year})`,
      start_date: startStr,
      end_date: endStr,
      isCurrent: startStr <= todayStr && todayStr <= endStr,
    })
  }
  return result
})

function getWeekNumber(date: Date): number {
  const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
  const dayNum = d.getUTCDay() || 7
  d.setUTCDate(d.getUTCDate() + 4 - dayNum)
  const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1))
  return Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7)
}

const showEditor = ref(false)
const editingPair = ref<Pair | undefined>(undefined)
const editingPairDate = ref<string | undefined>(undefined)
const pairEditorRef = ref<InstanceType<typeof PairEditor> | null>(null)

onMounted(async () => {
  const [groupsRes, subjectsRes, classroomsRes, teachersRes] = await Promise.all([
    groupsApi.getAll(),
    subjectsApi.getAll(),
    classroomsApi.getAll(),
    teachersApi.getAll(),
  ])
  groups.value = groupsRes.data
  subjects.value = subjectsRes.data
  classrooms.value = classroomsRes.data
  teachers.value = teachersRes.data

  if (weeks.value.length > 5) {
    selectedWeek.value = weeks.value[5] ?? null
  }
})

watch(
  [selectedGroupId, selectedWeek],
  async () => {
    if (!selectedGroupId.value || !selectedWeek.value) {
      scheduleStore.entries = []
      return
    }
    await scheduleStore.fetchSchedule({
      group_id: selectedGroupId.value,
      start_date: selectedWeek.value.start_date,
      end_date: selectedWeek.value.end_date,
    })
  },
)

async function handleSave(data: PairFormData) {
  try {
    if (editingPair.value) {
      await scheduleStore.updateEntry(editingPair.value.id, data)
    } else {
      await scheduleStore.addEntry(data)
    }
    showEditor.value = false
    editingPair.value = undefined
    editingPairDate.value = undefined
    if (selectedGroupId.value && selectedWeek.value) {
      await scheduleStore.fetchSchedule({
        group_id: selectedGroupId.value,
        start_date: selectedWeek.value.start_date,
        end_date: selectedWeek.value.end_date,
      })
    }
  } catch (e) {
    if (axios.isAxiosError(e) && e.response?.status === 409) {
      const response = e.response.data as ConflictResponse
      pairEditorRef.value?.showConflict(response.details)
    }
  }
}

function handleEdit(pair: Pair) {
  editingPair.value = pair
  const entry = scheduleStore.entries.find((e) => e.pairs.some((p) => p.id === pair.id))
  editingPairDate.value = entry?.date ?? selectedWeek.value?.start_date
  showEditor.value = true
}

async function handleDelete(id: number) {
  confirm.require({
    message: 'Вы уверены, что хотите удалить эту пару?',
    header: 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      await scheduleStore.deleteEntry(id)
      if (selectedGroupId.value && selectedWeek.value) {
        await scheduleStore.fetchSchedule({
          group_id: selectedGroupId.value,
          start_date: selectedWeek.value.start_date,
          end_date: selectedWeek.value.end_date,
        })
      }
    },
  })
}

function openAddModal() {
  editingPair.value = undefined
  editingPairDate.value = undefined
  showEditor.value = true
}

function closeModal() {
  showEditor.value = false
  editingPair.value = undefined
  editingPairDate.value = undefined
}

const scheduleEntries = computed(() => {
  if (!scheduleStore.entries) return []
  return scheduleStore.entries
})
</script>

<template>
  <div class="admin-schedule-view py-3">
    <h1 class="h4 mb-3">Управление расписанием</h1>

    <div class="schedule-filters mb-3">
      <div class="filter-item">
        <label class="form-label small text-muted">Неделя</label>
        <Select
          v-model="selectedWeek"
          :options="weeks"
          optionLabel="label"
          placeholder="Выберите неделю"
          class="w-100"
        >
          <template #option="{ option }">
            <span :class="{ 'fw-semibold': option.isCurrent, 'text-primary': option.isCurrent }">
              {{ option.label }}{{ option.isCurrent ? ' (текущая)' : '' }}
            </span>
          </template>
        </Select>
      </div>
      <div class="filter-item">
        <label class="form-label small text-muted">Группа</label>
        <Select
          v-model="selectedGroupId"
          :options="groups"
          optionLabel="name"
          optionValue="id"
          placeholder="Выберите группу"
          class="w-100"
        />
      </div>
      <div class="filter-item filter-item--button" v-if="selectedGroupId">
        <Button label="Добавить пару" icon="pi pi-plus" @click="openAddModal" class="w-100" />
      </div>
    </div>

    <template v-if="selectedGroupId">
      <div v-if="scheduleStore.loading" class="text-center py-4">Загрузка...</div>
      <div v-else-if="scheduleStore.error" class="text-center py-4 text-danger">{{ scheduleStore.error }}</div>
      <ScheduleTable
        v-else
        :entries="scheduleEntries"
        :show-actions="true"
        :admin-mode="true"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </template>
    <div v-else class="text-center text-muted py-5">
      Выберите группу для просмотра расписания
    </div>

    <Dialog v-model:visible="showEditor" :header="editingPair ? 'Редактирование пары' : 'Добавление пары'" :modal="true" class="schedule-dialog">
      <PairEditor
        ref="pairEditorRef"
        :is-edit="!!editingPair"
        :entry="editingPair"
        :groups="groups"
        :subjects="subjects"
        :classrooms="classrooms"
        :teachers="teachers"
        :initial-group-id="selectedGroupId ?? undefined"
        :initial-date="editingPair ? editingPairDate : selectedWeek?.start_date"
        @save="handleSave"
        @cancel="closeModal"
      />
    </Dialog>
    <ConfirmDialog />
  </div>
</template>

<style scoped>
.schedule-filters {
  display: flex;
  gap: 1rem;
  align-items: flex-end;
  flex-wrap: wrap;
}

.filter-item {
  flex: 1;
  min-width: 200px;
}

.filter-item--button {
  flex: 0 0 auto;
  min-width: auto;
}

@media (max-width: 768px) {
  .schedule-filters {
    flex-direction: column;
    gap: 0.75rem;
  }

  .filter-item {
    min-width: 100%;
  }

  .filter-item--button {
    width: 100%;
  }
}
</style>

<style>
.schedule-dialog .p-dialog {
  width: 90vw;
  max-width: 500px;
  margin: 0.5rem;
}

@media (max-width: 768px) {
  .schedule-dialog .p-dialog {
    width: 95vw;
    max-width: none;
    margin: 0.25rem;
  }

  .schedule-dialog .p-dialog-content {
    padding: 0.75rem;
    overflow-y: auto;
    max-height: 80vh;
  }
}
</style>
