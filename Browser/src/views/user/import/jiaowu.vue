<template>
  <div class="import-page">
    <a-card class="general-card" :title="t('user.import.jiaowu.title')">
      <a-upload
        :file-list="fileList"
        :before-upload="beforeUpload"
        @change="handleFileChange"
        @remove="handleFileRemove"
        :show-upload-list="true"
        :disabled="uploading"
        :custom-request="handleFileUpload"
      >
        <template #upload-button>
          <div class="upload-area">
            <div class="upload-icon">
              <icon-upload />
            </div>
            <div class="upload-text">{{ t('user.import.jiaowu.upload.text') }}</div>
          </div>
        </template>
      </a-upload>

      <div class="actions-section">
        <a-button type="primary" :loading="uploading" :disabled="uploadedFiles.length === 0" @click="handleSubmit">
          {{ t('user.import.jiaowu.submit') }}
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { checkFileExists, uploadFile, type UploadFileResponse } from '@/api/upload'
import { Message } from '@arco-design/web-vue'
import type { UploadRequestOption } from '@arco-design/web-vue/es/upload'
import SparkMD5 from 'spark-md5'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface FileItem {
  uid: string
  name: string
  url?: string
  status?: 'uploading' | 'done' | 'error'
  percent?: number
}

interface UploadedFile {
  file_id: string
  name: string
  size: number
  status: 'uploaded' | 'exists'
}

const fileList = ref<FileItem[]>([])
const uploadedFiles = ref<UploadedFile[]>([])
const uploading = ref(false)

const calculateFileMD5 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const spark = new SparkMD5.ArrayBuffer()
    const reader = new FileReader()
    reader.onload = (e) => {
      spark.append(e.target?.result as ArrayBuffer)
      resolve(spark.end())
    }
    reader.onerror = reject
    reader.readAsArrayBuffer(file)
  })
}

const beforeUpload = (file: File) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    Message.error(t('user.import.jiaowu.error.fileSize'))
    return false
  }
  return true
}

const handleFileChange = (files: FileItem[]) => {
  fileList.value = files
}

const handleFileRemove = (file: FileItem) => {
  const index = fileList.value.findIndex((item) => item.uid === file.uid)
  if (index !== -1) {
    fileList.value.splice(index, 1)
  }
  const uploadedIndex = uploadedFiles.value.findIndex((item) => item.name === file.name)
  if (uploadedIndex !== -1) {
    uploadedFiles.value.splice(uploadedIndex, 1)
  }
}

const handleFileUpload = async (options: UploadRequestOption) => {
  const { file, onProgress, onSuccess, onError } = options
  const fileObj = file as File
  try {
    const md5 = await calculateFileMD5(fileObj)
    const checkResponse = await checkFileExists({ md5, filename: fileObj.name })
    const fileExists = checkResponse.data.exists
    if (fileExists) {
      const response: UploadFileResponse = {
        file_id: checkResponse.data.file_id!,
        url: checkResponse.data.url!,
        md5,
      }
      uploadedFiles.value.push({
        file_id: response.file_id,
        name: fileObj.name,
        size: fileObj.size,
        status: 'exists',
      })
      onSuccess?.(response)
      return
    }
    const uploadResponse = await uploadFile(fileObj, (progress: number) => {
      onProgress?.(progress)
    })
    uploadedFiles.value.push({
      file_id: uploadResponse.data.file_id,
      name: fileObj.name,
      size: fileObj.size,
      status: 'uploaded',
    })
    onSuccess?.(uploadResponse.data)
  } catch (error) {
    onError?.(error as Error)
  }
}

const handleSubmit = async () => {
  if (uploadedFiles.value.length === 0) {
    Message.error(t('user.import.jiaowu.error.noFile'))
    return
  }
  uploading.value = true
  try {
    Message.success(t('user.import.jiaowu.success.submit'))
    fileList.value = []
    uploadedFiles.value = []
  } catch (error) {
    Message.error(t('user.import.jiaowu.error.submit'))
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.import-page {
  text-align: center;
  padding: 20px;
}
.upload-area {
  padding: 40px 0;
  text-align: center;
  border: 2px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  transition: border-color 0.3s;
  margin: 0 auto;
  min-width: 500px;
  max-width: 800px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}
.upload-area .upload-text {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 40px;
}
.upload-area:hover {
  border-color: #165dff;
}
.upload-icon {
  font-size: 48px;
  color: #165dff;
  margin-bottom: 16px;
}
.upload-text {
  font-size: 16px;
  color: #1d2129;
}
.actions-section {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}
.actions-section .arco-btn {
  min-width: 280px;
}
</style>
