<script setup lang="ts">
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Select from 'primevue/select'
import Button from 'primevue/button'

interface Field {
  key: string
  label: string
  type: 'text' | 'password' | 'select'
  options?: { value: string | number; label: string }[]
  required?: boolean
}

const props = defineProps<{
  fields: Field[]
  initialData?: Record<string, unknown>
  isEdit: boolean
  visible: boolean
}>()

const emit = defineEmits<{
  submit: [data: Record<string, unknown>]
  cancel: []
  'update:visible': [value: boolean]
}>()

const form = ref<Record<string, string>>({})

watch(
  () => props.initialData,
  (data) => {
    if (data) {
      const newForm: Record<string, string> = {}
      for (const key in data) {
        newForm[key] = String(data[key] ?? '')
      }
      form.value = newForm
    } else {
      props.fields.forEach((field) => {
        form.value[field.key] = ''
      })
    }
  },
  { immediate: true }
)

function handleSubmit() {
  emit('submit', { ...form.value })
}

function handleCancel() {
  emit('update:visible', false)
  emit('cancel')
}
</script>

<template>
  <Dialog
    :visible="visible"
    :header="isEdit ? 'Редактирование' : 'Добавление'"
    :modal="true"
    :style="{ width: '400px' }"
    @update:visible="emit('update:visible', $event)"
  >
    <form @submit.prevent="handleSubmit">
      <div v-for="field in fields" :key="field.key" class="mb-3">
        <label :for="field.key" class="form-label">{{ field.label }}</label>
        <InputText
          v-if="field.type === 'text'"
          :id="field.key"
          v-model="form[field.key]"
          class="w-100"
          :required="field.required"
        />
        <Password
          v-else-if="field.type === 'password'"
          :id="field.key"
          v-model="form[field.key]"
          class="w-100"
          toggleMask
          :required="field.required"
          inputClass="w-100"
        />
        <Select
          v-else-if="field.type === 'select'"
          :id="field.key"
          v-model="form[field.key]"
          :options="field.options"
          optionLabel="label"
          optionValue="value"
          class="w-100"
          :required="field.required"
        />
      </div>
      <div class="d-flex gap-2 justify-content-end mt-4">
        <Button label="Отмена" severity="secondary" text @click="handleCancel" />
        <Button label="Сохранить" type="submit" />
      </div>
    </form>
  </Dialog>
</template>

<style scoped>
:deep(.p-password) {
  width: 100%;
}
:deep(.p-password-input) {
  width: 100%;
}
</style>