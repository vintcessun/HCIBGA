import { getToken } from '@/utils/auth'
import axios from 'axios'

export interface StudentScoreRecord {
  project: string
  awardDate: string
  awardLevel: string
  awardType: 'individual' | 'collective'
  teamRank: string
  selfScore: number
  scoreBasis: string
  collegeScore: number
}

export interface StudentExportRecord {
  name: string
  accountId: string
  scores: StudentScoreRecord[]
}

export interface StudentExportResponse {
  generatedAt: string
  total: number
  list: StudentExportRecord[]
}

export function getStudentExportData(params?: Record<string, unknown>) {
  const accountId = getToken() || ''
  const finalParams = { ...params, accountId }
  return axios.post<StudentExportResponse>('/api/info/export/students', finalParams)
}
