import { API_CONFIG, getApiUrl } from '@/config/api'
import { RoleType, UserState } from '@/store/modules/user/types'
import axios from 'axios'
import type { RouteRecordNormalized } from 'vue-router'

import { getToken } from '@/utils/auth'

export interface LoginData {
  username: string
  password: string
}

export interface LoginRes {
  code: number
  data: {
    id: string
    username: string
    ip: string
    login_time: string
    role: RoleType
  }
  msg: string
}

export function login(data: LoginData): Promise<LoginRes> {
  const url = getApiUrl(API_CONFIG.PATHS.USER.LOGIN)
  // 在 Vite + TS 环境下，使用 console.info 保证在浏览器和 Node.js 控制台都能输出
  // eslint-disable-next-line no-console
  console.info('[login] 请求地址:', url)
  // eslint-disable-next-line no-console
  console.info('[login] 请求参数:', data)

  return axios
    .post<LoginRes>(url, data)
    .then((res) => {
      // eslint-disable-next-line no-console
      console.info('[login] 完整响应对象:', JSON.stringify(res, null, 2))
      // eslint-disable-next-line no-console
      console.info('[login] 响应数据对象:', JSON.stringify(res.data, null, 2))
      return res as unknown as LoginRes
    })
    .catch((err) => {
      // eslint-disable-next-line no-console
      console.error('[login] 请求出错:', err)
      throw err
    })
}

export function logout() {
  const accountId = getToken() || ''
  return axios.post<LoginRes>('/api/user/logout', { accountId })
}

export function getUserInfo() {
  const accountId = getToken() || ''
  return axios.post<UserState>('/api/user/info', { accountId })
}

export function getMenuList() {
  const accountId = getToken() || ''
  return axios.post<RouteRecordNormalized[]>('/api/user/menu', { accountId })
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
  const accountId = getToken() || ''
  return axios.post('/api/user/save-info', {
    accountId,
    ...data,
  })
}

export interface QrCodeRes {
  qrUrl: string
  qrId: string
}

// 获取二维码链接和ID
export function getQrCode() {
  return axios.get<QrCodeRes>('/api/user/qr-code')
}

export interface QrStatusRes {
  status: 'pending' | 'done' | 'expired'
}

// 轮询二维码状态
export function pollQrStatus(qrId: string) {
  return axios.get<QrStatusRes>('/api/user/qr-status', { params: { qrId } })
}

export interface QrAuthRes {
  token: string
  role: RoleType
}

// 根据二维码ID获取认证结果
export function getQrAuthResult(qrId: string) {
  return axios.post<QrAuthRes>('/api/user/qr-auth-result', { qrId })
}
