<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, onUnmounted } from 'vue'
import Button from 'primevue/button'

const router = useRouter()
const route = useRoute()

const isMobile = ref(false)

function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

const menuItems = [
  {
    label: 'Расписание',
    icon: 'pi pi-calendar',
    route: '/admin/schedule',
    command: () => router.push('/admin/schedule'),
  },
  {
    label: 'Преподаватели',
    icon: 'pi pi-users',
    route: '/admin/teachers',
    command: () => router.push('/admin/teachers'),
  },
  {
    label: 'Предметы',
    icon: 'pi pi-book',
    route: '/admin/subjects',
    command: () => router.push('/admin/subjects'),
  },
  {
    label: 'Аудитории',
    icon: 'pi pi-home',
    route: '/admin/classrooms',
    command: () => router.push('/admin/classrooms'),
  },
  {
    label: 'Группы',
    icon: 'pi pi-users',
    route: '/admin/groups',
    command: () => router.push('/admin/groups'),
  },
]

function isActive(itemRoute: string) {
  return route.path === itemRoute
}
</script>

<template>
  <div class="admin-layout">
    <aside v-if="!isMobile" class="sidebar d-flex flex-column">
      <h2 class="h5 mb-3 text-white sidebar-title">Админ-панель</h2>
      <nav class="d-flex flex-column gap-2 flex-grow-1">
        <Button
          v-for="item in menuItems"
          :key="item.label"
          :label="item.label"
          :icon="item.icon"
          severity="secondary"
          text
          class="text-white text-start sidebar-link"
          :class="{ 'sidebar-link--active': isActive(item.route) }"
          @click="item.command"
        />
      </nav>
    </aside>
    <main class="content flex-grow-1 p-3 overflow-auto" :class="{ 'content--mobile': isMobile }">
      <router-view />
    </main>
    <nav v-if="isMobile" class="bottom-nav">
      <button
        v-for="item in menuItems"
        :key="item.label"
        class="bottom-nav-item"
        :class="{ 'bottom-nav-item--active': isActive(item.route) }"
        @click="item.command"
      >
        <i :class="item.icon"></i>
        <span>{{ item.label }}</span>
      </button>
    </nav>
  </div>
</template>

<style scoped>
.admin-layout {
  height: 100%;
  display: flex;
  flex-direction: row;
}

.sidebar {
  width: 250px;
  background: #212529;
  padding: 0.75rem;
  flex-shrink: 0;
}

.sidebar-title {
  padding: 0.5rem;
}

.sidebar-link {
  border-radius: 8px;
}

.sidebar-link--active {
  background: rgba(255, 255, 255, 0.15) !important;
}

.content {
  background: #f8f9fa;
}

.content--mobile {
  padding: 0.75rem;
  padding-bottom: 5.5rem;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #212529;
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 0.5rem 0;
  z-index: 100;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.15);
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.625rem;
  padding: 0.375rem 0.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 0;
}

.bottom-nav-item i {
  font-size: 1.125rem;
}

.bottom-nav-item--active {
  color: white;
  background: rgba(255, 255, 255, 0.15);
}
</style>
