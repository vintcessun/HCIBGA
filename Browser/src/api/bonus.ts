import axios from 'axios'
import { getToken } from '@/utils/auth'

export interface BonusRecord {
  id: string
  project: string
  awardDate: string
  awardLevel: string
  awardType: 'individual' | 'collective'
  teamRank: string
  selfScore: number
  scoreBasis: string
  collegeScore: number
}

export interface BonusSummaryItem {
  category: string
  totalScore: number
  itemCount: number
}

export interface BonusSummaryResponse {
  totalScore: number
  items: BonusSummaryItem[]
}

function withAccountId<T extends Record<string, unknown>>(payload?: T) {
  const accountId = getToken() || ''
  return { ...(payload || {}), accountId }
}

export function getAcademicBonusList() {
  return axios.get<BonusRecord[]>('/api/bonus/academic/list', { params: withAccountId() })
}

export function getComprehensiveBonusList() {
  return axios.get<BonusRecord[]>('/api/bonus/comprehensive/list', { params: withAccountId() })
}

export function getBonusSummary() {
  return axios.get<BonusSummaryResponse>('/api/bonus/summary', { params: withAccountId() })
}
