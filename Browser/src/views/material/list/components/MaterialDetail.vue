<template>
  <div class="material-detail">
    <div class="detail-section">
      <h3>基本信息</h3>
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="标题">{{ material.title }}</a-descriptions-item>
        <a-descriptions-item label="分类">{{ getCategoryText(material.category) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(material.status)">
            {{ getStatusText(material.status) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="文件信息">
          <div v-if="material.files && material.files.length > 0">
            <div v-for="(file, idx) in material.files" :key="idx" style="margin-bottom: 8px">
              <div>文件名：{{ file.fileName }}</div>
              <div>大小：{{ formatFileSize(file.fileSize) }}</div>
            </div>
          </div>
          <span v-else>-</span>
        </a-descriptions-item>
        <a-descriptions-item label="上传者">{{ material.uploader }}</a-descriptions-item>
        <a-descriptions-item label="上传时间">{{ formatDate(material.uploadTime) }}</a-descriptions-item>
        <a-descriptions-item label="标签">
          <a-space v-if="material.tags && material.tags.length > 0">
            <a-tag v-for="tag in material.tags" :key="tag">{{ tag }}</a-tag>
          </a-space>
          <span v-else>-</span>
        </a-descriptions-item>
        <a-descriptions-item label="描述">
          {{ material.description || '无描述' }}
        </a-descriptions-item>
      </a-descriptions>
    </div>

    <div class="detail-section" v-if="material.reviewer">
      <h3>审核信息</h3>
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="审核人">{{ material.reviewer }}</a-descriptions-item>
        <a-descriptions-item label="审核时间">{{ material.reviewTime ? formatDate(material.reviewTime) : '-' }}</a-descriptions-item>
        <a-descriptions-item label="审核意见">
          {{ material.reviewComment || '无审核意见' }}
        </a-descriptions-item>
      </a-descriptions>
    </div>

    <div class="detail-section" v-if="material.aiReviewResult">
      <h3>AI审核结果</h3>
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="评分">{{ material.aiReviewResult.score }}/100</a-descriptions-item>
        <a-descriptions-item label="置信度">{{ (material.aiReviewResult.confidence * 100).toFixed(2) }}%</a-descriptions-item>
        <a-descriptions-item label="风险等级">
          <a-tag :color="getRiskLevelColor(material.aiReviewResult.riskLevel)">
            {{ getRiskLevelText(material.aiReviewResult.riskLevel) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="建议">
          <ul v-if="material.aiReviewResult.suggestions.length > 0">
            <li v-for="(suggestion, index) in material.aiReviewResult.suggestions" :key="index">
              {{ suggestion }}
            </li>
          </ul>
          <span v-else>无建议</span>
        </a-descriptions-item>
      </a-descriptions>
    </div>

    <div class="detail-section" v-if="material.files && material.files.length > 0">
      <h3>文件操作</h3>
      <div class="file-actions">
        <div
          v-for="(file, idx) in material.files"
          :key="idx"
          style="margin-bottom: 12px; padding: 12px; background: #fff; border: 1px solid #e5e6eb; border-radius: 6px"
        >
          <div style="font-weight: 600; margin-bottom: 8px">{{ file.fileName }}</div>
          <div style="margin-top: 8px">
            <a-space>
              <a-button type="primary" @click="() => handleDownloadFile(file)">
                <template #icon>
                  <icon-download />
                </template>
                下载
              </a-button>
              <a-button @click="() => handlePreviewFile(file)">
                <template #icon>
                  <icon-eye />
                </template>
                预览
              </a-button>
            </a-space>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import type { Material } from '@/api/material'
import { getMaterialFileUrl } from '@/api/material'
import { Message } from '@arco-design/web-vue'
import { PropType, defineComponent, defineOptions, defineProps } from 'vue'

export default defineComponent({
  name: 'MaterialDetail',
  props: {
    material: {
      type: Object as PropType<Material>,
      required: true,
    },
  },
  setup(props) {
    const getCategoryText = (category: string) => {
      const categoryMap: Record<string, string> = {
        document: '文档',
        image: '图片',
        video: '视频',
        audio: '音频',
        other: '其他',
      }
      return categoryMap[category] || category
    }

    const getStatusColor = (status: string) => {
      const colors: Record<string, string> = {
        pending: 'orange',
        approved: 'green',
        rejected: 'red',
      }
      return colors[status] || 'gray'
    }

    const getStatusText = (status: string) => {
      const texts: Record<string, string> = {
        pending: '待审核',
        approved: '已通过',
        rejected: '已拒绝',
      }
      return texts[status] || status
    }

    const getRiskLevelColor = (riskLevel: string) => {
      const colors: Record<string, string> = {
        low: 'green',
        medium: 'orange',
        high: 'red',
      }
      return colors[riskLevel] || 'gray'
    }

    const getRiskLevelText = (riskLevel: string) => {
      const texts: Record<string, string> = {
        low: '低风险',
        medium: '中风险',
        high: '高风险',
      }
      return texts[riskLevel] || riskLevel
    }

    const formatFileSize = (bytes: number) => {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return `${parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`
    }

    const formatDate = (dateString: string) => {
      return new Date(dateString).toLocaleString('zh-CN')
    }

    const handleDownloadFile = (file: { fileUrl: string; fileName: string }) => {
      if (file.fileName) {
        fetch(getMaterialFileUrl(file.fileName))
          .then((res) => res.blob())
          .then((blob) => {
            const url = window.URL.createObjectURL(blob)
            const link = document.createElement('a')
            link.href = url
            link.download = file.fileName
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
            window.URL.revokeObjectURL(url)
            Message.success('开始下载文件')
          })
          .catch(() => {
            Message.error('下载文件失败')
          })
      } else {
        Message.error('文件链接无效')
      }
    }

    const handlePreviewFile = (file: { fileUrl: string; fileName: string }) => {
      if (file.fileUrl) {
        // 触发父组件的 Drawer 打开逻辑
        const event = new CustomEvent('openPreviewDrawer', { detail: getMaterialFileUrl(file.fileName) })
        window.dispatchEvent(event)
      } else {
        Message.error('文件链接无效')
      }
    }

    const handleDownload = () => {
      if (props.material.files && props.material.files.length > 0) {
        props.material.files.forEach((file) => {
          if (file.fileUrl) {
            const link = document.createElement('a')
            link.href = file.fileUrl
            link.download = file.fileName
            link.click()
          }
        })
        Message.success('开始下载所有文件')
      } else {
        Message.error('没有可下载的文件')
      }
    }

    const handlePreview = () => {
      if (props.material.files && props.material.files.length > 0) {
        props.material.files.forEach((file) => {
          if (file.fileUrl) {
            window.open(file.fileUrl, '_blank')
          }
        })
      } else {
        Message.error('没有可预览的文件')
      }
    }

    return {
      getCategoryText,
      getStatusColor,
      getStatusText,
      getRiskLevelColor,
      getRiskLevelText,
      formatFileSize,
      formatDate,
      handleDownloadFile,
      handlePreviewFile,
      handleDownload,
      handlePreview,
    }
  },
})

defineOptions({
  name: 'MaterialDetail',
})

const props = defineProps({
  material: {
    type: Object as PropType<Material>,
    required: true,
  },
})

const getCategoryText = (category: string) => {
  const categoryMap: Record<string, string> = {
    document: '文档',
    image: '图片',
    video: '视频',
    audio: '音频',
    other: '其他',
  }
  return categoryMap[category] || category
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'orange',
    approved: 'green',
    rejected: 'red',
  }
  return colors[status] || 'gray'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
  }
  return texts[status] || status
}

const getRiskLevelColor = (riskLevel: string) => {
  const colors: Record<string, string> = {
    low: 'green',
    medium: 'orange',
    high: 'red',
  }
  return colors[riskLevel] || 'gray'
}

const getRiskLevelText = (riskLevel: string) => {
  const texts: Record<string, string> = {
    low: '低风险',
    medium: '中风险',
    high: '高风险',
  }
  return texts[riskLevel] || riskLevel
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const handleDownload = () => {
  if (props.material.files && props.material.files.length > 0) {
    props.material.files.forEach((file) => {
      if (file.fileUrl) {
        const link = document.createElement('a')
        link.href = file.fileUrl
        link.download = file.fileName
        link.click()
      }
    })
    Message.success('开始下载所有文件')
  } else {
    Message.error('没有可下载的文件')
  }
}

const handlePreview = () => {
  if (props.material.files && props.material.files.length > 0) {
    props.material.files.forEach((file) => {
      if (file.fileUrl) {
        window.open(file.fileUrl, '_blank')
      }
    })
  } else {
    Message.error('没有可预览的文件')
  }
}
</script>

<style lang="less" scoped>
.material-detail {
  .detail-section {
    margin-bottom: 24px;

    h3 {
      margin-bottom: 12px;
      color: #1d2129;
      font-weight: 600;
    }
  }

  .file-actions {
    padding: 16px;
    background: #f7f8fa;
    border-radius: 6px;
  }
}

// 优化描述列表的文字换行
:deep(.arco-descriptions) {
  .arco-descriptions-item-value {
    word-break: break-word;
    overflow-wrap: break-word;
    hyphens: auto;
  }

  // 长文本内容特殊处理
  .arco-descriptions-item[data-label='描述'],
  .arco-descriptions-item[data-label='审核意见'],
  .arco-descriptions-item[data-label='建议'] {
    .arco-descriptions-item-value {
      white-space: pre-wrap;
      line-height: 1.6;
    }
  }
}

// 移动端按钮适配
@media (max-width: 768px) {
  .material-detail {
    .file-actions {
      .arco-btn {
        width: 100%;
        margin-bottom: 8px;
        margin-left: 0 !important;
      }
    }
  }
}

// 移动端宽度修复 - 确保正确适应屏幕
@media (max-width: 768px) {
  .material-detail {
    width: 100%;
    max-width: 100vw;
    box-sizing: border-box;

    .detail-section {
      margin-bottom: 16px;

      h3 {
        font-size: 16px;
        margin-bottom: 8px;
        padding: 0 12px;
      }
    }

    :deep(.arco-descriptions) {
      width: 100% !important;
      max-width: 100% !important;

      .arco-descriptions-table {
        width: 100% !important;
        table-layout: fixed;
      }

      .arco-descriptions-item {
        width: 100% !important;
        display: flex;
        flex-direction: column;
        margin-bottom: 8px;

        .arco-descriptions-item-label {
          padding: 8px 12px 4px 12px;
          box-sizing: border-box;
          width: 100% !important;
          font-weight: 600;
          background-color: #f7f8fa;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .arco-descriptions-item-value {
          padding: 4px 12px 8px 12px;
          box-sizing: border-box;
          width: 100% !important;
          overflow-wrap: break-word;
          word-break: break-word;
          hyphens: auto;
          min-height: 20px;
        }
      }

      // 确保表格单元格正确显示
      .arco-descriptions-table {
        tbody {
          tr {
            display: flex;
            flex-direction: column;

            td {
              display: block;
              width: 100% !important;
            }
          }
        }
      }
    }
  }
}

// 超小屏幕适配
@media (max-width: 480px) {
  .material-detail {
    .file-actions {
      .arco-btn {
        font-size: 13px;
        padding: 8px 12px;
      }
    }

    :deep(.arco-descriptions) {
      .arco-descriptions-item {
        .arco-descriptions-item-label,
        .arco-descriptions-item-value {
          font-size: 13px;
          padding: 0 2px;
        }
      }
    }
  }
}
</style>
