import axios from 'axios'

import { HttpResponse } from '@/api/interceptor'
import { getToken } from '@/utils/auth'

export interface MyProjectRecord {
  id: number
  name: string
  description: string
  peopleNumber: number
  contributors: {
    name: string
    email: string
    avatar: string
  }[]
}
export function queryMyProjectList() {
  return axios.post<HttpResponse<MyProjectRecord[]>>('/api/user/my-project/list', {
    accountId: getToken(),
  })
}

export interface MyTeamRecord {
  id: number
  avatar: string
  name: string
  peopleNumber: number
}
export function queryMyTeamList() {
  return axios.post<HttpResponse<MyTeamRecord[]>>('/api/user/my-team/list', {
    accountId: getToken(),
  })
}

export interface LatestActivity {
  id: number
  title: string
  description: string
  avatar: string
}
export function queryLatestActivity() {
  return axios.post<HttpResponse<LatestActivity[]>>('/api/user/latest-activity', {
    accountId: getToken(),
  })
}

export interface BasicInfoModel {
  email: string
  nickname: string
  countryRegion: string
  area: string
  address: string
  profile: string
}

export interface EnterpriseCertificationModel {
  accountType: number
  status: number
  time: string
  legalPerson: string
  certificateType: string
  authenticationNumber: string
  enterpriseName: string
}

export type CertificationRecord = Array<{
  certificationType: number
  certificationContent: string
  status: number
  time: string
}>

export interface UnitCertification {
  enterpriseInfo: EnterpriseCertificationModel
  record: CertificationRecord
}

export function queryCertification() {
  return axios.post<UnitCertification>('/api/user/certification', {
    accountId: getToken(),
  })
}

export function userUploadApi(
  data: FormData,
  config: {
    controller: AbortController
    onUploadProgress?: (progressEvent: any) => void
  }
) {
  return axios.post('/api/user/upload', data, config)
}
