<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useScheduleStore } from '@/stores/schedule'
import { useAuthStore } from '@/stores/auth'
import { groupsApi, subjectsApi, classroomsApi, teachersApi } from '@/api/client'
import ScheduleTable from '@/components/ScheduleTable.vue'
import PairEditor from '@/components/PairEditor.vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { Group, Subject, Classroom, Teacher, Pair, PairFormData, ConflictResponse } from '@/types'
import axios from 'axios'

const confirm = useConfirm()

const scheduleStore = useScheduleStore()
const authStore = useAuthStore()

const groups = ref<Group[]>([])
const subjects = ref<Subject[]>([])
const classrooms = ref<Classroom[]>([])
const teachers = ref<Teacher[]>([])

const showEditor = ref(false)
const editingPair = ref<Pair | undefined>(undefined)
const editingGroupId = ref<number | undefined>(undefined)
const editingDate = ref<string | undefined>(undefined)
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

  const currentTeacher = teachers.value.find((t) => t.login === localStorage.getItem('token'))
  if (currentTeacher) {
    authStore.setUserData({ role: authStore.role!, teacherId: currentTeacher.id, name: currentTeacher.name })
  }

  await scheduleStore.fetchSchedule({ teacher_id: authStore.teacherId ?? undefined })
})

async function handleSave(data: PairFormData) {
  try {
    if (editingPair.value) {
      await scheduleStore.updateEntry(editingPair.value.id, data)
    } else {
      await scheduleStore.addEntry(data)
    }
    showEditor.value = false
    editingPair.value = undefined
  } catch (e) {
    if (axios.isAxiosError(e) && e.response?.status === 409) {
      const response = e.response.data as ConflictResponse
      pairEditorRef.value?.showConflict(response.details)
    }
  }
}

function handleEdit(pair: Pair) {
  const entry = scheduleStore.entries.find((e) => e.pairs.some((p) => p.id === pair.id))
  editingPair.value = pair
  editingGroupId.value = entry?.group_id
  editingDate.value = entry?.date
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
    },
  })
}

function openAddModal() {
  editingPair.value = undefined
  showEditor.value = true
}

function closeModal() {
  showEditor.value = false
  editingPair.value = undefined
  editingGroupId.value = undefined
  editingDate.value = undefined
}

const myPairs = computed(() => {
  if (!scheduleStore.entries) return []
  return scheduleStore.entries.filter(
    (entry) => entry.pairs.some((p) => p.teacher_id === authStore.teacherId)
  )
})
</script>

<template>
  <div class="teacher-view container py-3" style="max-width: 1200px">
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h1 class="h4 mb-0">Панель преподавателя</h1>
      <Button label="Добавить пару" icon="pi pi-plus" @click="openAddModal" />
    </div>

    <div v-if="scheduleStore.loading" class="text-center py-4">Загрузка...</div>
    <div v-else-if="scheduleStore.error" class="text-center py-4 text-danger">{{ scheduleStore.error }}</div>
    <ScheduleTable
      v-else
      :entries="myPairs"
      :show-actions="true"
      @edit="handleEdit"
      @delete="handleDelete"
    />

    <Dialog v-model:visible="showEditor" :header="editingPair ? 'Редактирование пары' : 'Добавление пары'" :modal="true" :style="{ width: '500px' }">
      <PairEditor
        ref="pairEditorRef"
        :is-edit="!!editingPair"
        :entry="editingPair"
        :groups="groups"
        :subjects="subjects"
        :classrooms="classrooms"
        :teachers="teachers"
        :initial-teacher-id="authStore.teacherId ?? undefined"
        :initial-group-id="editingGroupId"
        :initial-date="editingDate"
        @save="handleSave"
        @cancel="closeModal"
      />
    </Dialog>
    <ConfirmDialog />
  </div>
</template>