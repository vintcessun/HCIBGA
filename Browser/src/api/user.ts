import { API_CONFIG, getApiUrl } from '@/config/api'
import { UserState } from '@/store/modules/user/types'
import axios from 'axios'
import type { RouteRecordNormalized } from 'vue-router'

export interface LoginData {
  username: string
  password: string
}

export interface LoginRes {
  token: string
  role: string
  username: string
  message: string
}
export function login(data: LoginData) {
  return axios.post<LoginRes>(getApiUrl(API_CONFIG.PATHS.USER.LOGIN), data)
}

export function logout() {
  return axios.post<LoginRes>('/api/user/logout')
}

export function getUserInfo() {
  return axios.post<UserState>('/api/user/info')
}

export function getMenuList() {
  return axios.post<RouteRecordNormalized[]>('/api/user/menu')
}

export interface SaveInfoData {
  name?: string
  email?: string
  nickname?: string
  countryRegion?: string
  area?: string
  address?: string
  profile?: string
}

export function saveUserInfo(data: SaveInfoData) {
  return axios.post('/api/user/save-info', data)
}
