<template>
  <div class="material-review">
    <a-card :title="$t('menu.material.review')" class="review-card">
      <!-- 待审核材料列表 -->
      <a-table
        :columns="columns"
        :data="pendingMaterials"
        :loading="loading"
        :pagination="pagination"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #status>
          <a-tag color="orange">待审核</a-tag>
        </template>

        <template #uploadTime="{ record }">
          {{ formatDate(record.uploadTime) }}
        </template>

        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="handleView(record)">
              <template #icon>
                <icon-eye />
              </template>
              查看详情
            </a-button>

            <a-button type="text" size="small" status="success" @click="handleQuickApprove(record)">
              <template #icon>
                <icon-check />
              </template>
              快速通过
            </a-button>

            <a-button type="text" size="small" status="warning" @click="handleReview(record, 'rejected')">
              <template #icon>
                <icon-close />
              </template>
              拒绝
            </a-button>

            <a-button type="text" size="small" @click="handleReview(record, 'approved')">
              <template #icon>
                <icon-edit />
              </template>
              详细审核
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 查看详情模态框 -->
    <a-modal v-model:visible="showDetailModal" :title="currentMaterial?.title" :footer="false" :width="detailModalWidth">
      <template v-if="currentMaterial">
        <div class="detail-modal-content">
          <MaterialDetail :material="currentMaterial" />
        </div>
        <div class="detail-footer">
          <a-space size="large">
            <a-button type="secondary" @click="viewPrevious">上一条</a-button>
            <a-button type="primary" @click="approveCurrent">同意</a-button>
            <a-button status="danger" @click="rejectCurrent">拒绝</a-button>
            <a-button type="secondary" @click="viewNext">下一条</a-button>
          </a-space>
        </div>
      </template>
      <template v-else>
        <a-empty description="暂无材料" />
      </template>
    </a-modal>

    <!-- 审核模态框 -->
    <a-modal
      v-model:visible="showReviewModal"
      :title="`审核材料 - ${currentMaterial?.title}`"
      @ok="handleConfirmReview"
      @cancel="showReviewModal = false"
      width="500px"
    >
      <a-form :model="reviewForm" layout="vertical">
        <a-form-item label="审核结果" field="status" required>
          <a-radio-group v-model="reviewForm.status">
            <a-radio value="approved">通过</a-radio>
            <a-radio value="rejected">拒绝</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="审核意见" field="comment">
          <a-textarea
            v-model="reviewForm.comment"
            :placeholder="reviewForm.status === 'approved' ? '请输入通过意见（可选）' : '请输入拒绝原因'"
            :auto-size="{ minRows: 3, maxRows: 6 }"
          />
        </a-form-item>

        <a-form-item v-if="reviewForm.status === 'rejected'" label="拒绝原因模板">
          <a-select v-model="selectedTemplate" placeholder="选择拒绝原因模板" @change="applyTemplate" allow-clear>
            <a-option value="content_issue">内容不符合要求</a-option>
            <a-option value="quality_issue">材料质量不达标</a-option>
            <a-option value="format_issue">格式不正确</a-option>
            <a-option value="copyright_issue">版权问题</a-option>
            <a-option value="other">其他原因</a-option>
          </a-select>
        </a-form-item>

        <!-- 审核信息 -->
        <a-form-item label="审核人" v-if="userStore.userInfo">
          <a-input :value="userStore.userInfo.name" disabled />
        </a-form-item>

        <a-form-item label="审核时间">
          <a-input :value="new Date().toLocaleString('zh-CN')" disabled />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 预览抽屉 -->
    <a-drawer v-model:visible="showPreviewDrawer" :title="currentMaterial?.title" :width="drawerWidth" placement="right">
      <iframe v-if="currentMaterial && previewFileUrl" :src="previewFileUrl" style="width: 100%; height: 100%; border: none"></iframe>
      <a-empty v-else description="暂无可预览文件" />
    </a-drawer>

    <!-- 批量操作区域 -->
    <div class="batch-actions" v-if="selectedMaterials.length > 0">
      <a-space>
        <span>已选择 {{ selectedMaterials.length }} 个材料</span>
        <a-button type="primary" @click="handleBatchApprove">
          <template #icon>
            <icon-check />
          </template>
          批量通过
        </a-button>
        <a-button status="warning" @click="handleBatchReject">
          <template #icon>
            <icon-close />
          </template>
          批量拒绝
        </a-button>
        <a-button @click="clearSelection">
          <template #icon>
            <icon-close-circle />
          </template>
          取消选择
        </a-button>
      </a-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { batchReviewMaterials, getPendingMaterials, reviewMaterial, type Material } from '@/api/material'
import { useUserStore } from '@/store'
import { Message } from '@arco-design/web-vue'
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import MaterialDetail from '../list/components/MaterialDetail.vue'

const userStore = useUserStore()
const loading = ref(false)
const pendingMaterials = ref<Material[]>([])
const selectedMaterials = ref<string[]>([])
const showDetailModal = ref(false)
const showReviewModal = ref(false)
const showPreviewDrawer = ref(false)
const currentMaterial = ref<Material | null>(null)
const previewFileUrl = ref('')
const selectedTemplate = ref('')

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const reviewForm = reactive({
  status: 'approved' as 'approved' | 'rejected',
  comment: '',
})

const detailModalWidth = computed(() => {
  return window.innerWidth < 768 ? '95%' : '600px'
})

const drawerWidth = computed(() => {
  return window.innerWidth < 768 ? '95%' : '800px'
})

const columns = computed(() => [
  {
    title: '选择',
    type: 'checkbox',
    width: 60,
  },
  {
    title: '标题',
    dataIndex: 'title',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '分类',
    dataIndex: 'category',
    width: 100,
    render: ({ record }: { record: Material }) => {
      const categoryMap: Record<string, string> = {
        document: '文档',
        image: '图片',
        video: '视频',
        audio: '音频',
        other: '其他',
      }
      return categoryMap[record.category] || record.category
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 100,
    slotName: 'status',
  },
  {
    title: '上传者',
    dataIndex: 'uploader',
    width: 120,
  },
  {
    title: '上传时间',
    dataIndex: 'uploadTime',
    width: 180,
    slotName: 'uploadTime',
  },
  {
    title: '操作',
    slotName: 'actions',
    width: 300,
    fixed: 'right',
  },
])

const rowSelection = computed(() => ({
  type: 'checkbox',
  selectedRowKeys: selectedMaterials.value,
  onChange: (keys: string[]) => {
    selectedMaterials.value = keys
  },
}))

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

const fetchPendingMaterials = async () => {
  loading.value = true
  try {
    const response = await getPendingMaterials()
    pendingMaterials.value = response.data
    pagination.total = pendingMaterials.value.length
  } catch (error) {
    Message.error('获取待审核材料失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pagination.current = page
}

const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
}

const handleView = (material: Material) => {
  currentMaterial.value = material
  showDetailModal.value = true
}

const handleQuickApprove = async (material: Material) => {
  try {
    await reviewMaterial({
      materialId: material.id,
      status: 'approved',
      comment: '快速审核通过',
    })
    Message.success('审核通过')
    fetchPendingMaterials()
  } catch (error) {
    Message.error('审核失败')
  }
}

const handleReview = (material: Material, status: 'approved' | 'rejected') => {
  currentMaterial.value = material
  reviewForm.status = status
  reviewForm.comment = ''
  selectedTemplate.value = ''
  showReviewModal.value = true
}

const applyTemplate = (template: string) => {
  const templates: Record<string, string> = {
    content_issue: '内容不符合平台要求，请修改后重新提交',
    quality_issue: '材料质量不达标，无法通过审核',
    format_issue: '文件格式不正确，请使用支持的格式',
    copyright_issue: '存在版权问题，请确保拥有合法授权',
    other: '其他原因，请详细说明',
  }
  reviewForm.comment = templates[template] || ''
}

const handleConfirmReview = async () => {
  if (!currentMaterial.value) return

  try {
    await reviewMaterial({
      materialId: currentMaterial.value.id,
      status: reviewForm.status,
      comment: reviewForm.comment,
    })
    Message.success('审核完成')
    showReviewModal.value = false
    fetchPendingMaterials()
  } catch (error) {
    Message.error('审核失败')
  }
}

const clearSelection = () => {
  selectedMaterials.value = []
}

const handleBatchApprove = async () => {
  if (selectedMaterials.value.length === 0) return

  try {
    await batchReviewMaterials(selectedMaterials.value, 'approved', '批量审核通过')
    Message.success(`已通过 ${selectedMaterials.value.length} 个材料`)
    clearSelection()
    fetchPendingMaterials()
  } catch (error) {
    Message.error('批量审核失败')
  }
}

const handleBatchReject = async () => {
  if (selectedMaterials.value.length === 0) return

  try {
    await batchReviewMaterials(selectedMaterials.value, 'rejected', '批量审核拒绝')
    Message.success(`已拒绝 ${selectedMaterials.value.length} 个材料`)
    clearSelection()
    fetchPendingMaterials()
  } catch (error) {
    Message.error('批量审核失败')
  }
}

const handleOpenPreviewDrawer = (e: Event) => {
  const url = (e as CustomEvent<string>).detail
  previewFileUrl.value = url
  showPreviewDrawer.value = true
}

onMounted(() => {
  fetchPendingMaterials()
  window.addEventListener('openPreviewDrawer', handleOpenPreviewDrawer)
})

onUnmounted(() => {
  window.removeEventListener('openPreviewDrawer', handleOpenPreviewDrawer)
})

const viewPrevious = () => {
  const index = pendingMaterials.value.findIndex((m) => m.id === currentMaterial.value?.id)
  if (index > 0) {
    currentMaterial.value = pendingMaterials.value[index - 1]
  }
}

const viewNext = () => {
  const index = pendingMaterials.value.findIndex((m) => m.id === currentMaterial.value?.id)
  if (index >= 0 && index < pendingMaterials.value.length - 1) {
    currentMaterial.value = pendingMaterials.value[index + 1]
  }
}

const approveCurrent = async () => {
  if (!currentMaterial.value) return
  try {
    await reviewMaterial({
      materialId: currentMaterial.value.id,
      status: 'approved',
      comment: '同意审核',
    })
    Message.success('已同意该材料')
    fetchPendingMaterials()
  } catch (error) {
    Message.error('同意操作失败')
  }
}

const rejectCurrent = () => {
  if (!currentMaterial.value) return
  reviewForm.status = 'rejected'
  reviewForm.comment = ''
  selectedTemplate.value = ''
  showReviewModal.value = true
}
</script>

<style lang="less" scoped>
.detail-modal-content {
  max-height: 65vh;
  overflow-y: auto;
  margin-bottom: 16px;
  padding-right: 4px;
}

.detail-footer {
  text-align: center;
  padding: 8px 0 4px;
  border-top: 1px solid #f2f3f5;
}

.material-review {
  padding: 20px;

  .review-card {
    min-height: 500px;
  }

  .batch-actions {
    position: fixed;
    bottom: 20px;
    right: 20px;
    padding: 16px;
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 1000;
  }
}

// 移动端适配
@media (max-width: 768px) {
  .material-review {
    padding: 12px;

    .batch-actions {
      position: static;
      margin-top: 16px;
      box-shadow: none;
      border: 1px solid #e5e6eb;

      .arco-space {
        width: 100%;
        justify-content: center;
      }

      .arco-btn {
        width: 100%;
        margin-bottom: 8px;
      }
    }
  }

  .detail-modal-content {
    max-height: 55vh;
  }

  .detail-footer {
    .arco-btn {
      width: 100%;
      margin-bottom: 8px;
    }

    :deep(.arco-space) {
      width: 100%;
      flex-direction: column;
    }
  }
}
</style>
