import { getToken } from '@/utils/auth'
import axios from 'axios'

export interface FileCheckRequest {
  md5: string
  filename: string
}

export interface FileCheckResponse {
  exists: boolean
  file_id?: string
  url?: string
}

export interface UploadFileResponse {
  file_id: string
  url: string
  md5: string
}

export interface BatchUploadResponseItem {
  file_id: string
  url: string
  md5: string
  status: 'uploaded' | 'exists'
}

export interface FileInfo {
  id: string
  name: string
  size: number
  url: string
  md5: string
  upload_time: string
  uploader: string
}

// 检查文件是否存在
export function checkFileExists(data: FileCheckRequest) {
  const accountId = getToken() || ''
  const finalData = { ...data, accountId }
  return axios.post<FileCheckResponse>('/api/upload/check', finalData)
}

// 上传单个文件
export function uploadFile(file: File, onProgress?: (progress: number) => void) {
  const accountId = getToken() || ''
  const formData = new FormData()
  formData.append('file', file)
  formData.append('accountId', accountId)

  return axios.post<UploadFileResponse>('/api/upload/file', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress(progress)
      }
    },
  })
}

// 批量上传文件
export function batchUploadFiles(files: File[], onProgress?: (progress: number) => void) {
  const formData = new FormData()
  files.forEach((file, index) => {
    formData.append(`files`, file)
  })

  return axios.post<BatchUploadResponseItem[]>('/api/upload/batch', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress(progress)
      }
    },
  })
}

// 获取文件列表
export function getFileList() {
  return axios.get<FileInfo[]>('/api/upload/files')
}
