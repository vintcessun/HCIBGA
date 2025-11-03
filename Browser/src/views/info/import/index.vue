<template>
  <div class="info-import">
    <a-card class="general-card" :title="'从 Excel 导入'">
      <template #extra>
        <a-link>帮助</a-link>
      </template>
      <div class="upload-section">
        <a-upload :multiple="false" :before-upload="beforeExcelUpload" :custom-request="handleExcelUpload">
          <template #upload-button>
            <div class="upload-area">
              <icon-upload />
              <div class="upload-text">点击选择 Excel 文件</div>
              <div class="upload-hint">支持 .xls, .xlsx 格式</div>
            </div>
          </template>
        </a-upload>
      </div>
    </a-card>

    <a-card class="general-card" :title="'从文本导入'" style="margin-top: 16px">
      <template #extra>
        <a-link>帮助</a-link>
      </template>
      <div class="upload-section">
        <a-textarea v-model="textContent" placeholder="请输入文本内容" :auto-size="{ minRows: 3, maxRows: 5 }" />
        <a-upload :multiple="false" :before-upload="beforeTxtUpload" :custom-request="handleTxtUpload">
          <template #upload-button>
            <div class="upload-area">
              <icon-upload />
              <div class="upload-text">点击上传 TXT 文件</div>
              <div class="upload-hint">支持 .txt 格式</div>
            </div>
          </template>
        </a-upload>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { Message } from '@arco-design/web-vue'
import type { UploadRequestOption } from '@arco-design/web-vue/es/upload'
import { ref } from 'vue'

const textContent = ref('')

const beforeExcelUpload = (file: File) => {
  const isExcel = /\.(xls|xlsx)$/i.test(file.name)
  if (!isExcel) {
    Message.error('只能上传 Excel 文件')
    return false
  }
  return true
}

const handleExcelUpload = async (options: UploadRequestOption) => {
  const { file, onSuccess } = options
  // TODO: 解析 Excel
  Message.success(`Excel 文件 ${file.name} 已选择`)
  onSuccess?.({})
}

const beforeTxtUpload = (file: File) => {
  const isTxt = /\.txt$/i.test(file.name)
  if (!isTxt) {
    Message.error('只能上传 TXT 文件')
    return false
  }
  return true
}

const handleTxtUpload = async (options: UploadRequestOption) => {
  const { file, onSuccess } = options
  // TODO: 解析 TXT
  Message.success(`TXT 文件 ${file.name} 已选择`)
  onSuccess?.({})
}
</script>

<style lang="less" scoped>
.info-import {
  text-align: center;
  padding: 20px;

  .upload-section {
    text-align: center;
    margin-bottom: 24px;
  }

  .section-title {
    margin-bottom: 12px;
    font-size: 16px;
    font-weight: 600;
  }

  .upload-area {
    padding: 12px 0;
    text-align: center;
    border: 2px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    transition: border-color 0.3s;
    width: 100%;
    margin: 0;
    min-width: 600px;
    height: 50px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    &:hover {
      border-color: #165dff;
    }

    .upload-text {
      font-size: 14px;
      color: #1d2129;
      margin-top: 8px;
      text-align: center;
    }

    .upload-hint {
      font-size: 12px;
      color: #86909c;
      margin-top: 4px;
      text-align: center;
    }
  }
}
</style>
