<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'

const router = useRouter()
const authStore = useAuthStore()

const login = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true

  const result = await authStore.login(login.value, password.value)

  if (result.success) {
    if (authStore.isAdmin) {
      router.push('/admin')
    } else {
      router.push('/teacher')
    }
  } else {
    error.value = result.error || 'Неверные учетные данные'
    console.error('Login error:', authStore.lastError)
    console.log('Login:', login.value, password.value)
  }

  loading.value = false
}
</script>

<template>
  <div class="login-view d-flex justify-content-center align-items-center min-vh-80">
    <Card class="w-100" style="max-width: 400px">
      <template #title>
        <div class="text-center mb-3">Вход в систему</div>
      </template>
      <template #content>
        <form @submit.prevent="handleLogin">
          <div class="mb-3">
            <label for="login" class="form-label">Логин</label>
            <InputText id="login" v-model="login" class="w-100" required />
          </div>
          <div class="mb-3">
            <label for="password" class="form-label">Пароль</label>
            <Password id="password" v-model="password" class="w-100" :feedback="false" toggleMask required inputClass="w-100" />
          </div>
          <Message v-if="error" severity="error" :closable="false" class="mb-3">{{ error }}</Message>
          <Button type="submit" :loading="loading" :label="loading ? 'Вход...' : 'Войти'" class="w-100" />
        </form>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.min-vh-80 {
  min-height: 80vh;
}
:deep(.p-password) {
  width: 100%;
}
:deep(.p-password-input) {
  width: 100%;
}
</style>