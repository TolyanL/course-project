import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'schedule',
      component: () => import('@/views/ScheduleView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/teacher',
      name: 'teacher',
      component: () => import('@/views/TeacherView.vue'),
      meta: { requiresAuth: true, role: 'teacher' },
    },
    {
      path: '/admin',
      component: () => import('@/views/AdminView.vue'),
      meta: { requiresAuth: true, role: 'admin' },
      children: [
        {
          path: '',
          redirect: '/admin/schedule',
        },
        {
          path: 'schedule',
          name: 'admin-schedule',
          component: () => import('@/views/AdminScheduleView.vue'),
        },
        {
          path: 'teachers',
          name: 'admin-teachers',
          component: () => import('@/views/AdminTeachersView.vue'),
        },
        {
          path: 'subjects',
          name: 'admin-subjects',
          component: () => import('@/views/AdminSubjectsView.vue'),
        },
        {
          path: 'classrooms',
          name: 'admin-classrooms',
          component: () => import('@/views/AdminClassroomsView.vue'),
        },
        {
          path: 'groups',
          name: 'admin-groups',
          component: () => import('@/views/AdminGroupsView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, _from) => {
  const authStore = useAuthStore()
  authStore.initFromStorage()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return '/login'
  }

  if (to.meta.guest && authStore.isAuthenticated) {
    if (authStore.isAdmin) {
      return '/admin'
    } else {
      return '/teacher'
    }
  }

  if (to.meta.role) {
    if (authStore.role !== to.meta.role) {
      return '/'
    }
  }

  return true
})

export default router