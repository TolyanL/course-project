<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { classroomsApi } from '@/api/client'
import EntityTable from '@/components/EntityTable.vue'
import EntityForm from '@/components/EntityForm.vue'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { Classroom, ClassroomFormData } from '@/types'

const confirm = useConfirm()

const classrooms = ref<Classroom[]>([])
const showForm = ref(false)
const editingItem = ref<Classroom | undefined>(undefined)

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'number', label: 'Номер' },
]

const fields = [{ key: 'number', label: 'Номер аудитории', type: 'text' as const, required: true }]

onMounted(async () => {
  await loadClassrooms()
})

async function loadClassrooms() {
  const { data } = await classroomsApi.getAll()
  classrooms.value = data
}

function handleAdd() {
  editingItem.value = undefined
  showForm.value = true
}

function handleEdit(item: Record<string, unknown>) {
  editingItem.value = item as unknown as Classroom
  showForm.value = true
}

async function handleDelete(id: number) {
  confirm.require({
    message: 'Вы уверены, что хотите удалить эту аудиторию?',
    header: 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await classroomsApi.delete(id)
        await loadClassrooms()
      } catch (e: unknown) {
        const error = e as { response?: { data?: { error?: string } } }
        if (error.response?.data?.error) {
          alert(error.response.data.error)
        }
      }
    },
  })
}

async function handleSubmit(data: Record<string, unknown>) {
  const formData = data as unknown as ClassroomFormData
  if (editingItem.value) {
    await classroomsApi.update(editingItem.value.id, formData)
  } else {
    await classroomsApi.create(formData)
  }
  showForm.value = false
  await loadClassrooms()
}

function handleCancel() {
  showForm.value = false
  editingItem.value = undefined
}
</script>

<template>
  <div class="admin-classrooms-view">
    <h1 class="h4 mb-3">Управление аудиториями</h1>
    <EntityTable
      :columns="columns"
      :data="classrooms"
      @add="handleAdd"
      @edit="handleEdit"
      @delete="handleDelete"
    />
    <EntityForm
      :fields="fields"
      :initial-data="editingItem"
      :is-edit="!!editingItem"
      v-model:visible="showForm"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
    <ConfirmDialog />
  </div>
</template>
