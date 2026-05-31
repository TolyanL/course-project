export interface AuthResponse {
  token: string
  role: 'admin' | 'teacher'
}

export interface Teacher {
  id: number
  name: string
  login: string
  role: 'admin' | 'teacher'
  created_at: string
}

export interface Subject {
  id: number
  name: string
}

export interface Classroom {
  id: number
  number: string
}

export interface Group {
  id: number
  name: string
}

export interface Pair {
  id: number
  schedule_id: number
  subject_id: number
  teacher_id: number
  classroom_id: number
  pair_number: number
  subject: Subject
  teacher: Teacher
  classroom: Classroom
}

export interface ScheduleEntry {
  id: number
  group_id: number
  date: string
  pairs: Pair[]
  group: Group
}

export interface ConflictDetail {
  type: 'classroom' | 'teacher'
  message: string
  field: string
  entry_id: number
}

export interface ConflictResponse {
  error: 'conflict detected'
  details: ConflictDetail[]
}

export interface ApiError {
  error: string
}

export interface PairFormData {
  group_id: number
  date: string
  subject_id: number
  teacher_id: number
  classroom_id: number
  pair_number: number
  force_save?: boolean
}

export interface TeacherFormData {
  name: string
  login: string
  password?: string
  role: 'teacher' | 'admin'
}

export interface SubjectFormData {
  name: string
}

export interface ClassroomFormData {
  number: string
}

export interface GroupFormData {
  name: string
}

export const PAIR_TIMES: Record<number, string> = {
  1: '08:00 - 09:35',
  2: '09:45 - 11:20',
  3: '12:20 - 13:55',
  4: '14:05 - 15:40',
  5: '16:00 - 17:35',
  6: '17:45 - 19:20',
  7: '19:30 - 21:05',
}
