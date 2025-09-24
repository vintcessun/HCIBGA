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
  files: string[]  // 文件UUID列表
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

// 上传材料
export function uploadMaterial(data: UploadMaterialRequest) {
  return axios.post<Material>('/api/material/upload', data)
}

// 获取材料列表
export function getMaterialList(params?: MaterialFilter) {
  return axios.post<Material[]>('/api/material/list', params)
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
  return axios.post('/api/material/delete', { materialId })
}

// 批量审核
export function batchReviewMaterials(materialIds: string[], status: 'approved' | 'rejected', comment?: string) {
  return axios.post('/api/material/batch-review', { materialIds, status, comment })
}
