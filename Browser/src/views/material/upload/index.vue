<template>
  <div class="material-upload">
    <a-card :title="$t('menu.material.upload')" class="upload-card">
      <!-- 文件上传区域 -->
      <div class="upload-section">
        <h3 class="section-title">{{ $t('material.upload.files') }}</h3>
        <a-upload
          :multiple="true"
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
              <div class="upload-text">
                {{ $t('material.upload.clickToUpload') }}
              </div>
              <div class="upload-hint">
                {{ $t('material.upload.supportFormat') }}
              </div>
            </div>
          </template>
        </a-upload>

        <!-- 已上传文件列表 -->
        <div v-if="uploadedFiles.length > 0" class="uploaded-files">
          <h4 class="files-title">{{ $t('material.upload.uploadedFiles') }}</h4>
          <div v-for="file in uploadedFiles" :key="file.file_id" class="file-item">
            <icon-file />
            <span class="file-name">{{ file.name }}</span>
            <span class="file-status" :class="file.status">
              {{ file.status === 'uploaded' ? '已上传' : '已存在' }}
            </span>
          </div>
        </div>
      </div>

      <!-- 材料信息表单 -->
      <div class="form-section">
        <h3 class="section-title">{{ $t('material.upload.materialInfo') }}</h3>
        <a-form
          :model="form"
          :label-col-props="{ span: 6 }"
          :wrapper-col-props="{ span: 18 }"
          class="upload-form"
          :rules="formRules"
          ref="formRef"
        >
          <a-form-item :label="$t('material.upload.title')" field="title">
            <a-input v-model="form.title" :placeholder="$t('material.upload.titlePlaceholder')" :max-length="100" />
          </a-form-item>

          <a-form-item :label="$t('material.upload.description')" field="description">
            <a-textarea
              v-model="form.description"
              :placeholder="$t('material.upload.descriptionPlaceholder')"
              :auto-size="{
                minRows: 3,
                maxRows: 5,
              }"
              :max-length="500"
            />
          </a-form-item>

          <a-form-item :label="$t('material.upload.category')" field="category">
            <a-select v-model="form.category" :placeholder="$t('material.upload.categoryPlaceholder')" allow-clear>
              <a-option value="document">文档</a-option>
              <a-option value="image">图片</a-option>
              <a-option value="video">视频</a-option>
              <a-option value="audio">音频</a-option>
              <a-option value="other">其他</a-option>
            </a-select>
          </a-form-item>

          <a-form-item :label="$t('material.upload.tags')" field="tags">
            <a-select v-model="form.tags" multiple :placeholder="$t('material.upload.tagsPlaceholder')" allow-clear>
              <a-option value="important">重要</a-option>
              <a-option value="urgent">紧急</a-option>
              <a-option value="normal">普通</a-option>
            </a-select>
          </a-form-item>
        </a-form>
      </div>

      <!-- 操作按钮 -->
      <div class="actions-section">
        <a-form-item :wrapper-col-props="{ span: 24 }">
          <div class="upload-actions">
            <a-button type="primary" :loading="uploading" :disabled="uploadedFiles.length === 0 || !formValid" @click="handleSubmit">
              {{ $t('material.upload.submit') }}
            </a-button>
            <a-button @click="handleReset">
              {{ $t('material.upload.reset') }}
            </a-button>
          </div>
        </a-form-item>
      </div>
    </a-card>

    <!-- 成功模态框 -->
    <a-modal v-model:visible="showSuccessModal" :title="$t('material.upload.successTitle')" :footer="false">
      <div class="success-content">
        <icon-check-circle style="color: #00b42a; font-size: 48px; margin-bottom: 16px" />
        <p>{{ $t('material.upload.successMessage') }}</p>
        <a-button type="primary" @click="showSuccessModal = false">
          {{ $t('material.upload.confirm') }}
        </a-button>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { uploadMaterial } from '@/api/material'
import { type UploadFileResponse, checkFileExists, uploadFile } from '@/api/upload'
import { Message } from '@arco-design/web-vue'
import type { UploadRequestOption } from '@arco-design/web-vue/es/upload'
import SparkMD5 from 'spark-md5'
import { computed, reactive, ref } from 'vue'

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
const showSuccessModal = ref(false)
const formRef = ref()

const form = reactive({
  title: '',
  description: '',
  category: '',
  tags: [] as string[],
})

const formRules = {
  title: [{ required: true, message: '请输入材料标题' }],
  category: [{ required: true, message: '请选择材料类别' }],
}

const formValid = computed(() => {
  return form.title.trim() && form.category
})

// 计算文件MD5
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
    Message.error('文件大小不能超过10MB')
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

  // 同时从已上传文件列表中移除
  const uploadedIndex = uploadedFiles.value.findIndex((item) => item.name === file.name)
  if (uploadedIndex !== -1) {
    uploadedFiles.value.splice(uploadedIndex, 1)
  }
}

const handleFileUpload = async (options: UploadRequestOption) => {
  const { file, onProgress, onSuccess, onError } = options
  const fileObj = file as File

  try {
    // 计算文件MD5
    const md5 = await calculateFileMD5(fileObj)

    // 先检查文件是否存在
    const checkResponse = await checkFileExists({ md5, filename: fileObj.name })
    const fileExists = checkResponse.data.exists

    if (fileExists) {
      // 文件已存在，直接返回成功
      const response: UploadFileResponse = {
        file_id: checkResponse.data.file_id!,
        url: checkResponse.data.url!,
        md5,
      }

      // 添加到已上传文件列表
      uploadedFiles.value.push({
        file_id: response.file_id,
        name: fileObj.name,
        size: fileObj.size,
        status: 'exists',
      })

      onSuccess?.(response)
      return
    }

    // 文件不存在，调用真正的上传API
    const uploadResponse = await uploadFile(fileObj, (progress: number) => {
      onProgress?.(progress)
    })

    // 添加到已上传文件列表
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

const handleReset = () => {
  fileList.value = []
  uploadedFiles.value = []
  form.title = ''
  form.description = ''
  form.category = ''
  form.tags = []
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  if (uploadedFiles.value.length === 0) {
    Message.error('请先上传文件')
    return
  }

  if (!formValid.value) {
    Message.error('请填写完整的材料信息')
    return
  }

  uploading.value = true
  try {
    // 调用材料上传API
    await uploadMaterial({
      title: form.title,
      description: form.description,
      category: form.category,
      tags: form.tags,
      files: uploadedFiles.value.map((file) => file.file_id),
    })

    Message.success('材料上传成功！')
    showSuccessModal.value = true
    handleReset()
  } catch (error) {
    Message.error('上传失败，请重试')
  } finally {
    uploading.value = false
  }
}
</script>

<style lang="less" scoped>
.material-upload {
  padding: 20px;

  .upload-card {
    max-width: 800px;
    margin: 0 auto;
  }

  .section-title {
    margin-bottom: 16px;
    color: #1d2129;
    font-size: 16px;
    font-weight: 600;
  }

  .upload-section {
    margin-bottom: 24px;
  }

  .upload-area {
    padding: 40px 0;
    text-align: center;
    border: 2px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    transition: border-color 0.3s;
    margin: 0 auto;
    max-width: 400px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    &:hover {
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
      margin-bottom: 8px;
    }

    .upload-hint {
      font-size: 14px;
      color: #86909c;
    }
  }

  .uploaded-files {
    margin-top: 16px;
    padding: 16px;
    border: 1px solid #e5e6eb;
    border-radius: 6px;
    background: #f7f8fa;

    .files-title {
      margin-bottom: 12px;
      color: #1d2129;
      font-size: 14px;
      font-weight: 600;
    }

    .file-item {
      display: flex;
      align-items: center;
      padding: 8px;
      margin-bottom: 8px;
      background: white;
      border-radius: 4px;
      border: 1px solid #e5e6eb;

      &:last-child {
        margin-bottom: 0;
      }

      .file-name {
        flex: 1;
        margin-left: 8px;
        color: #1d2129;
      }

      .file-status {
        padding: 2px 8px;
        border-radius: 12px;
        font-size: 12px;

        &.uploaded {
          background: #e8ffea;
          color: #00b42a;
        }

        &.exists {
          background: #e8f4ff;
          color: #165dff;
        }
      }
    }
  }

  .form-section {
    margin-bottom: 24px;
  }

  .upload-form {
    margin-top: 16px;
  }

  .actions-section {
    margin-top: 24px;
    display: flex;
    justify-content: center;
  }

  .upload-actions {
    display: flex;
    justify-content: center;
    gap: 12px;

    .arco-btn {
      min-width: 100px;
    }
  }

  .success-content {
    text-align: center;
    padding: 20px;

    p {
      margin: 16px 0 24px;
      color: #1d2129;
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .material-upload {
    padding: 12px;

    .upload-card {
      margin: 0;
    }

    .upload-area {
      padding: 24px 0;
    }

    .upload-form {
      :deep(.arco-form-item) {
        margin-bottom: 16px;
      }
    }

    .upload-actions {
      .arco-btn {
        width: 100%;
        margin: 8px 0;
      }
    }
  }
}
</style>
