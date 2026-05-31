import axios from 'axios'
import type {
  AuthResponse,
  ScheduleEntry,
  Pair,
  Teacher,
  Subject,
  Classroom,
  Group,
  PairFormData,
  TeacherFormData,
  SubjectFormData,
  ClassroomFormData,
  GroupFormData,
} from '@/types'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

export const authApi = {
  login: (login: string, password: string) =>
    api.post<AuthResponse>('/auth/login', { login, password }),
}

export const scheduleApi = {
  getSchedule: (params?: {
    group_id?: number
    teacher_id?: number
    date?: string
    start_date?: string
    end_date?: string
  }) => api.get<ScheduleEntry[]>('/schedule', { params }),

  createEntry: (data: PairFormData) => api.post<{ entry: Pair }>('/schedule/entries', data),

  updateEntry: (id: number, data: PairFormData) =>
    api.put<{ entry: Pair }>(`/schedule/entries/${id}`, data),

  deleteEntry: (id: number) => api.delete(`/schedule/entries/${id}`),
}

export const teachersApi = {
  getAll: () => api.get<Teacher[]>('/teachers'),

  create: (data: TeacherFormData) => api.post<Teacher>('/teachers', data),

  update: (id: number, data: Partial<TeacherFormData>) => api.put<Teacher>(`/teachers/${id}`, data),

  delete: (id: number) => api.delete(`/teachers/${id}`),
}

export const subjectsApi = {
  getAll: () => api.get<Subject[]>('/subjects'),

  create: (data: SubjectFormData) => api.post<Subject>('/subjects', data),

  update: (id: number, data: SubjectFormData) => api.put<Subject>(`/subjects/${id}`, data),

  delete: (id: number) => api.delete(`/subjects/${id}`),
}

export const classroomsApi = {
  getAll: () => api.get<Classroom[]>('/classrooms'),

  create: (data: ClassroomFormData) => api.post<Classroom>('/classrooms', data),

  update: (id: number, data: ClassroomFormData) => api.put<Classroom>(`/classrooms/${id}`, data),

  delete: (id: number) => api.delete(`/classrooms/${id}`),
}

export const groupsApi = {
  getAll: () => api.get<Group[]>('/groups'),

  create: (data: GroupFormData) => api.post<Group>('/groups', data),

  update: (id: number, data: GroupFormData) => api.put<Group>(`/groups/${id}`, data),

  delete: (id: number) => api.delete(`/groups/${id}`),
}

export default api
