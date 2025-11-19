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
          <a-tag color="orange">{{ $t('material.review.status.pending') }}</a-tag>
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
              {{ $t('material.review.action.view') }}
            </a-button>

            <a-button type="text" size="small" status="success" @click="handleQuickApprove(record)">
              <template #icon>
                <icon-check />
              </template>
              {{ $t('material.review.action.quickApprove') }}
            </a-button>

            <a-button type="text" size="small" status="warning" @click="handleReview(record, 'rejected')">
              <template #icon>
                <icon-close />
              </template>
              {{ $t('material.review.action.reject') }}
            </a-button>

            <a-button type="text" size="small" @click="handleReview(record, 'approved')">
              <template #icon>
                <icon-edit />
              </template>
              {{ $t('material.review.action.detailReview') }}
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
            <a-button type="secondary" @click="viewPrevious">{{ $t('material.review.detail.prev') }}</a-button>
            <a-button type="primary" @click="approveCurrent">{{ $t('material.review.detail.approve') }}</a-button>
            <a-button status="danger" @click="rejectCurrent">{{ $t('material.review.detail.reject') }}</a-button>
            <a-button type="secondary" @click="viewNext">{{ $t('material.review.detail.next') }}</a-button>
          </a-space>
        </div>
      </template>
      <template v-else>
        <a-empty :description="$t('material.review.empty')" />
      </template>
    </a-modal>

    <!-- 审核模态框 -->
    <a-modal
      v-model:visible="showReviewModal"
      :title="$t('material.review.modal.title', { title: currentMaterial?.title })"
      @ok="handleConfirmReview"
      @cancel="showReviewModal = false"
      width="500px"
    >
      <a-form :model="reviewForm" layout="vertical">
        <a-form-item :label="$t('material.review.form.result')" field="status" required>
          <a-radio-group v-model="reviewForm.status">
            <a-radio value="approved">{{ $t('material.review.form.approved') }}</a-radio>
            <a-radio value="rejected">{{ $t('material.review.form.rejected') }}</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item :label="$t('material.review.form.comment')" field="comment">
          <a-textarea
            v-model="reviewForm.comment"
            :placeholder="
              reviewForm.status === 'approved' ? $t('material.review.placeholder.approved') : $t('material.review.placeholder.rejected')
            "
            :auto-size="{ minRows: 3, maxRows: 6 }"
          />
        </a-form-item>

        <a-form-item v-if="reviewForm.status === 'rejected'" :label="$t('material.review.form.template')">
          <a-select
            v-model="selectedTemplate"
            :placeholder="$t('material.review.placeholder.template')"
            @change="applyTemplate"
            allow-clear
          >
            <a-option value="content_issue">{{ $t('material.review.template.content') }}</a-option>
            <a-option value="quality_issue">{{ $t('material.review.template.quality') }}</a-option>
            <a-option value="format_issue">{{ $t('material.review.template.format') }}</a-option>
            <a-option value="copyright_issue">{{ $t('material.review.template.copyright') }}</a-option>
            <a-option value="other">{{ $t('material.review.template.other') }}</a-option>
          </a-select>
        </a-form-item>

        <!-- 审核信息 -->
        <a-form-item :label="$t('material.review.form.reviewer')" v-if="userStore.userInfo">
          <a-input :value="userStore.userInfo.name" disabled />
        </a-form-item>

        <a-form-item :label="$t('material.review.form.time')">
          <a-input :value="new Date().toLocaleString('zh-CN')" disabled />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 预览抽屉 -->
    <a-drawer v-model:visible="showPreviewDrawer" :title="currentMaterial?.title" :width="drawerWidth" placement="right">
      <iframe v-if="currentMaterial && previewFileUrl" :src="previewFileUrl" style="width: 100%; height: 100%; border: none"></iframe>
      <a-empty v-else :description="$t('material.review.preview.empty')" />
    </a-drawer>

    <!-- 批量操作区域 -->
    <div class="batch-actions" v-if="selectedMaterials.length > 0">
      <a-space>
        <span>{{ $t('material.review.batch.selected', { count: selectedMaterials.length }) }}</span>
        <a-button type="primary" @click="handleBatchApprove">
          <template #icon>
            <icon-check />
          </template>
          {{ $t('material.review.batch.approve') }}
        </a-button>
        <a-button status="warning" @click="handleBatchReject">
          <template #icon>
            <icon-close />
          </template>
          {{ $t('material.review.batch.reject') }}
        </a-button>
        <a-button @click="clearSelection">
          <template #icon>
            <icon-close-circle />
          </template>
          {{ $t('material.review.batch.cancel') }}
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
import { useI18n } from 'vue-i18n'
import MaterialDetail from '../list/components/MaterialDetail.vue'

const { t } = useI18n()
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

const getCategoryText = (category: string) => {
  const categoryMap: Record<string, string> = {
    '学术专长成绩-科研成果': t('material.category.academic.research'),
    '学术专长成绩-学业竞赛': t('material.category.academic.competition'),
    '学术专长成绩-创新创业训练': t('material.category.academic.innovation'),
    '综合表现加分-国际组织实习': t('material.category.comprehensive.internship'),
    '综合表现加分-参军入伍服兵役': t('material.category.comprehensive.military'),
    '综合表现加分-志愿服务': t('material.category.comprehensive.volunteer'),
    '综合表现加分-荣誉称号': t('material.category.comprehensive.honor'),
    '综合表现加分-社会工作': t('material.category.comprehensive.social'),
    '综合表现加分-体育比赛': t('material.category.comprehensive.sports'),
  }
  return categoryMap[category] || category
}

const detailModalWidth = computed(() => {
  return window.innerWidth < 768 ? '95%' : '600px'
})

const drawerWidth = computed(() => {
  return window.innerWidth < 768 ? '95%' : '800px'
})

const columns = computed(() => [
  {
    title: t('material.review.column.selection'),
    type: 'checkbox',
    width: 30,
  },
  {
    title: t('material.review.column.title'),
    dataIndex: 'title',
    ellipsis: true,
    tooltip: true,
    width: 80,
  },
  {
    title: t('material.review.column.category'),
    dataIndex: 'category',
    width: 100,
    render: ({ record }: { record: Material }) => {
      return getCategoryText(record.category)
    },
  },
  {
    title: t('material.review.column.status'),
    dataIndex: 'status',
    width: 60,
    slotName: 'status',
  },
  {
    title: t('material.review.column.uploader'),
    dataIndex: 'uploader',
    width: 90,
  },
  {
    title: t('material.review.column.uploadTime'),
    dataIndex: 'uploadTime',
    width: 110,
    slotName: 'uploadTime',
  },
  {
    title: t('material.review.column.actions'),
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
    Message.error(t('material.review.message.fetchFailed'))
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
      comment: t('material.review.comment.quickApprove'),
    })
    Message.success(t('material.review.message.approveSuccess'))
    fetchPendingMaterials()
  } catch (error) {
    Message.error(t('material.review.message.failed'))
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
    content_issue: t('material.review.template.content.detail'),
    quality_issue: t('material.review.template.quality.detail'),
    format_issue: t('material.review.template.format.detail'),
    copyright_issue: t('material.review.template.copyright.detail'),
    other: t('material.review.template.other.detail'),
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
    Message.success(t('material.review.message.complete'))
    showReviewModal.value = false
    fetchPendingMaterials()
  } catch (error) {
    Message.error(t('material.review.message.failed'))
  }
}

const clearSelection = () => {
  selectedMaterials.value = []
}

const handleBatchApprove = async () => {
  if (selectedMaterials.value.length === 0) return

  try {
    await batchReviewMaterials(selectedMaterials.value, 'approved', t('material.review.comment.batchApprove'))
    Message.success(t('material.review.message.batchApproveSuccess', { count: selectedMaterials.value.length }))
    clearSelection()
    fetchPendingMaterials()
  } catch (error) {
    Message.error(t('material.review.message.batchFailed'))
  }
}

const handleBatchReject = async () => {
  if (selectedMaterials.value.length === 0) return

  try {
    await batchReviewMaterials(selectedMaterials.value, 'rejected', t('material.review.comment.batchReject'))
    Message.success(t('material.review.message.batchRejectSuccess', { count: selectedMaterials.value.length }))
    clearSelection()
    fetchPendingMaterials()
  } catch (error) {
    Message.error(t('material.review.message.batchFailed'))
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
      comment: t('material.review.comment.agree'),
    })
    Message.success(t('material.review.message.agreeSuccess'))
    fetchPendingMaterials()
  } catch (error) {
    Message.error(t('material.review.message.agreeFailed'))
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
