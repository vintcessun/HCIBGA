import { getToken } from '@/utils/auth'
import axios from 'axios'

export interface AIReviewResult {
  score: number
  confidence: number
  suggestions: string[]
  riskLevel: 'low' | 'medium' | 'high'
}

export interface Material {
  id: string
  title: string
  description: string
  category: string
  tags: string[]
  files: {
    fileUrl: string
    fileName: string
    fileSize: number
  }[]
  status: 'pending' | 'approved' | 'rejected'
  uploader: string
  uploadTime: string
  reviewer?: string
  reviewTime?: string
  reviewComment?: string
  aiReviewResult?: AIReviewResult
}

export interface UploadMaterialRequest {
  title: string
  description: string
  category: string
  tags: string[]
  files: string[] // 文件UUID列表
}

export interface ReviewMaterialRequest {
  materialId: string
  status: 'approved' | 'rejected'
  comment?: string
}

export interface MaterialFilter {
  status?: string
  category?: string
  uploader?: string
  startDate?: string
  endDate?: string
}

export interface CheckFileExistsRequest {
  md5: string
  filename: string
}

export interface CheckFileExistsResponse {
  exists: boolean
  fileId: string
  url: string
}

export interface UploadFileResponse {
  fileId: string
  url: string
  md5: string
}

export interface LlmFillRequest {
  file_id: string
  metadata: {
    title?: string
    description?: string
    category?: string
    tags?: string[]
  }
}

export interface LlmFillResponse {
  title: string
  description: string
  category: string
  tags: string[]
}

// 检查文件是否存在
export function checkFileExists(data: CheckFileExistsRequest) {
  const accountId = getToken() || ''
  const finalData = { ...data, accountId }
  return axios.post<CheckFileExistsResponse>('/api/upload/check', finalData)
}

// 上传文件
export function uploadFile(formData: FormData) {
  const accountId = getToken() || ''
  formData.append('accountId', accountId)
  return axios.post<UploadFileResponse>('/api/upload/file', formData)
}

// 调用LLM自动填写材料信息
export function llmFillMaterialInfo(data: LlmFillRequest) {
  const accountId = getToken() || ''
  const finalData = { ...data, accountId }
  return axios.post<LlmFillResponse>('/api/material/llm-fill', finalData)
}

// 上传材料
export function uploadMaterial(data: UploadMaterialRequest) {
  const accountId = getToken() || ''
  const finalData = { ...data, accountId }
  return axios.post<Material>('/api/material/upload', finalData)
}

// 获取材料列表
export function getMaterialList(params?: MaterialFilter) {
  const accountId = getToken() || ''
  const finalParams = { ...params, accountId }
  return axios.get<Material[]>('/api/material/list', { params: finalParams })
}

// 获取待审核材料
export function getPendingMaterials() {
  const accountId = getToken() || ''
  return axios.post<Material[]>('/api/material/pending', { accountId })
}

// 审核材料
export function reviewMaterial(data: ReviewMaterialRequest) {
  const accountId = getToken() || ''
  const finalData = { ...data, accountId }
  return axios.post<Material>('/api/material/review', finalData)
}

// 获取材料统计
export function getMaterialStatistics() {
  const accountId = getToken() || ''

  return axios.post<any>('/api/material/statistics', { accountId })
}

// 删除材料
export function deleteMaterial(materialId: string) {
  const accountId = getToken() || ''
  return axios.delete(`/api/material/${materialId}`, { data: { accountId } })
}

export function batchReviewMaterials(materialIds: string[], status: 'approved' | 'rejected', comment?: string) {
  const accountId = getToken() || ''
  const finalData = { materialIds, status, comment, accountId }
  return axios.post('/api/material/batch-review', finalData)
}

// 获取材料文件下载/预览地址（返回完整链接）
export function getMaterialFileUrl(fileId: string) {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || ''
  return `${baseUrl}/upload/${fileId}`
}
