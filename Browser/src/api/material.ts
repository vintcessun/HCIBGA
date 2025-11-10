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
  fileUrl: string
  fileName: string
  fileSize: number
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
  return axios.post<CheckFileExistsResponse>('/api/upload/check', data)
}

// 上传文件
export function uploadFile(formData: FormData) {
  return axios.post<UploadFileResponse>('/api/upload/file', formData)
}

// 调用LLM自动填写材料信息
export function llmFillMaterialInfo(data: LlmFillRequest) {
  return axios.post<LlmFillResponse>('/api/material/llm-fill', data)
}

// 上传材料
export function uploadMaterial(data: UploadMaterialRequest) {
  return axios.post<Material>('/api/material/upload', data)
}

// 获取材料列表
export function getMaterialList(params?: MaterialFilter) {
  return axios.get<Material[]>('/api/material/list', { params })
}

// 获取待审核材料
export function getPendingMaterials() {
  return axios.post<Material[]>('/api/material/pending')
}

// 审核材料
export function reviewMaterial(data: ReviewMaterialRequest) {
  return axios.post<Material>('/api/material/review', data)
}

// 获取材料统计
export function getMaterialStatistics() {
  return axios.post<any>('/api/material/statistics')
}

// 删除材料
export function deleteMaterial(materialId: string) {
  return axios.delete(`/api/material/${materialId}`)
}

// 批量审核
export function batchReviewMaterials(materialIds: string[], status: 'approved' | 'rejected', comment?: string) {
  return axios.post('/api/material/batch-review', { materialIds, status, comment })
}
