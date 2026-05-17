<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { teachersApi } from '@/api/client'
import EntityTable from '@/components/EntityTable.vue'
import EntityForm from '@/components/EntityForm.vue'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import type { Teacher, TeacherFormData } from '@/types'

const confirm = useConfirm()

const teachers = ref<Teacher[]>([])
const showForm = ref(false)
const editingItem = ref<Teacher | undefined>(undefined)

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'name', label: 'Имя' },
  { key: 'login', label: 'Логин' },
  { key: 'role', label: 'Роль' },
]

const createFields = [
  { key: 'name', label: 'Имя', type: 'text' as const, required: true },
  { key: 'login', label: 'Логин', type: 'text' as const, required: true },
  { key: 'password', label: 'Пароль', type: 'password' as const, required: true },
  {
    key: 'role',
    label: 'Роль',
    type: 'select' as const,
    required: true,
    options: [
      { value: 'teacher', label: 'Преподаватель' },
      { value: 'admin', label: 'Админ' },
    ],
  },
]

const editFields = [
  { key: 'name', label: 'Имя', type: 'text' as const, required: true },
  {
    key: 'role',
    label: 'Роль',
    type: 'select' as const,
    required: true,
    options: [
      { value: 'teacher', label: 'Преподаватель' },
      { value: 'admin', label: 'Админ' },
    ],
  },
]

onMounted(async () => {
  await loadTeachers()
})

async function loadTeachers() {
  const { data } = await teachersApi.getAll()
  teachers.value = data
}

function handleAdd() {
  editingItem.value = undefined
  showForm.value = true
}

function handleEdit(item: Record<string, unknown>) {
  editingItem.value = item as unknown as Teacher
  showForm.value = true
}

async function handleDelete(id: number) {
  confirm.require({
    message: 'Вы уверены, что хотите удалить этого преподавателя?',
    header: 'Подтверждение',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await teachersApi.delete(id)
        await loadTeachers()
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
  const formData = data as unknown as TeacherFormData
  console.log('Creating teacher - URL: http://localhost:8080/api/teachers')
  console.log('Creating teacher with data:', JSON.stringify(formData))
  if (editingItem.value) {
    await teachersApi.update(editingItem.value.id, formData)
  } else {
    await teachersApi.create(formData)
  }
  showForm.value = false
  await loadTeachers()
}

function handleCancel() {
  showForm.value = false
  editingItem.value = undefined
}
</script>

<template>
  <div class="admin-teachers-view">
    <h1 class="h4 mb-3">Управление преподавателями</h1>
    <EntityTable
      :columns="columns"
      :data="teachers"
      @add="handleAdd"
      @edit="handleEdit"
      @delete="handleDelete"
    />
    <EntityForm
      :fields="editingItem ? editFields : createFields"
      :initial-data="editingItem"
      :is-edit="!!editingItem"
      v-model:visible="showForm"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
    <ConfirmDialog />
  </div>
</template>