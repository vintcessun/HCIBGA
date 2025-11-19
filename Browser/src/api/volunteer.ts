import { getToken } from '@/utils/auth'
import axios from 'axios'

export interface VolunteerHoursRequest {
  username: string
  password: string
}

export interface VolunteerHoursResponse {
  credit_hours: number
  honor_hours: number
  total_hours: number
}

export function fetchVolunteerHours(data: VolunteerHoursRequest) {
  const accountId = getToken() || ''
  return axios.post<VolunteerHoursResponse>('/api/volunteer/credit', { ...data, accountId })
}
