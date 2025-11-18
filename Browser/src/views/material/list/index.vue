<template>
  <div class="material-list">
    <a-card :title="$t('menu.material.list')" class="list-card">
      <!-- 搜索和筛选区域 -->
      <div class="filter-section">
        <a-form :model="filterForm" layout="inline" class="filter-form">
          <a-form-item :label="$t('material.list.status')">
            <a-select v-model="filterForm.status" :placeholder="$t('material.list.statusPlaceholder')" allow-clear style="width: 120px">
              <a-option value="pending">待审核</a-option>
              <a-option value="approved">已通过</a-option>
              <a-option value="rejected">已拒绝</a-option>
            </a-select>
          </a-form-item>

          <a-form-item :label="$t('material.list.category')">
            <a-select v-model="filterForm.category" :placeholder="$t('material.list.categoryPlaceholder')" allow-clear style="width: 120px">
              <a-option value="document">文档</a-option>
              <a-option value="image">图片</a-option>
              <a-option value="video">视频</a-option>
              <a-option value="audio">音频</a-option>
              <a-option value="other">其他</a-option>
            </a-select>
          </a-form-item>

          <a-form-item :label="$t('material.list.uploader')">
            <a-input v-model="filterForm.uploader" :placeholder="$t('material.list.uploaderPlaceholder')" style="width: 150px" />
          </a-form-item>

          <a-form-item :label="$t('material.list.dateRange')">
            <a-range-picker v-model="filterForm.dateRange" style="width: 240px" />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" @click="handleSearch">
              <template #icon>
                <icon-search />
              </template>
              {{ $t('material.list.search') }}
            </a-button>
            <a-button @click="handleReset" style="margin-left: 8px">
              <template #icon>
                <icon-refresh />
              </template>
              {{ $t('material.list.reset') }}
            </a-button>
          </a-form-item>
        </a-form>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-section">
        <a-space>
          <a-button
            v-if="userStore.role === 'admin' || userStore.role === 'user'"
            type="primary"
            @click="handleBatchDelete"
            :disabled="selectedMaterials.length === 0"
          >
            <template #icon>
              <icon-delete />
            </template>
            {{ $t('material.list.batchDelete') }}
          </a-button>
          <a-button v-if="userStore.role === 'admin' || userStore.role === 'user'" @click="handleExport">
            <template #icon>
              <icon-download />
            </template>
            {{ $t('material.list.export') }}
          </a-button>
        </a-space>
      </div>

      <!-- 材料列表 -->
      <a-table
        :columns="columns"
        :data="materials"
        :pagination="pagination"
        :loading="loading"
        :row-selection="rowSelection"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      >
        <template #status="{ record }">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
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
              查看
            </a-button>

            <a-button
              v-if="userStore.role === 'admin' || userStore.role === 'reviewer'"
              type="text"
              size="small"
              status="warning"
              @click="handleEdit(record)"
            >
              <template #icon>
                <icon-edit />
              </template>
              修改
            </a-button>

            <a-button
              v-if="userStore.role === 'admin' || userStore.role === 'user'"
              type="text"
              size="small"
              status="danger"
              @click="handleDelete(record)"
            >
              <template #icon>
                <icon-delete />
              </template>
              删除
            </a-button>

            <a-button
              v-if="userStore.role === 'reviewer' && record.status === 'pending'"
              type="text"
              size="small"
              status="success"
              @click="handleReview(record, 'approved')"
            >
              <template #icon>
                <icon-check />
              </template>
              通过
            </a-button>

            <a-button
              v-if="userStore.role === 'reviewer' && record.status === 'pending'"
              type="text"
              size="small"
              status="warning"
              @click="handleReview(record, 'rejected')"
            >
              <template #icon>
                <icon-close />
              </template>
              拒绝
            </a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 查看详情模态框 -->
    <a-modal v-model:visible="showDetailModal" :title="currentMaterial?.title" :footer="false" :width="modalWidth">
      <MaterialDetail :material="currentMaterial" v-if="currentMaterial" />
    </a-modal>

    <!-- 审核模态框 -->
    <a-modal
      v-model:visible="showReviewModal"
      :title="`审核材料 - ${currentMaterial?.title}`"
      @ok="handleConfirmReview"
      @cancel="showReviewModal = false"
    >
      <a-form :model="reviewForm" layout="vertical">
        <a-form-item label="审核结果" field="status">
          <a-radio-group v-model="reviewForm.status">
            <a-radio value="approved">通过</a-radio>
            <a-radio value="rejected">拒绝</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="审核意见" field="comment">
          <a-textarea v-model="reviewForm.comment" placeholder="请输入审核意见" :auto-size="{ minRows: 3 }" />
        </a-form-item>
      </a-form>
    </a-modal>
    <!-- 新增右侧预览抽屉 -->
    <a-drawer v-model:visible="showPreviewDrawer" :title="currentMaterial?.title" :width="drawerWidth" placement="right">
      <iframe
        v-if="currentMaterial && currentMaterial.files && currentMaterial.files.length > 0"
        :src="currentMaterial.files[0].fileUrl"
        style="width: 100%; height: 100%; border: none"
      ></iframe>
    </a-drawer>
  </div>
</template>

<script lang="ts" setup>
import { deleteMaterial, getMaterialList, reviewMaterial, type Material } from '@/api/material'
import { useUserStore } from '@/store'
import { Message } from '@arco-design/web-vue'
import { computed, onMounted, reactive, ref } from 'vue'
import MaterialDetail from './components/MaterialDetail.vue'

const userStore = useUserStore()
const loading = ref(false)
const materials = ref<Material[]>([])
const selectedMaterials = ref<string[]>([])
const showDetailModal = ref(false)
const showReviewModal = ref(false)
const currentMaterial = ref<Material | null>(null)
const showPreviewDrawer = ref(false)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
})

const filterForm = reactive({
  status: '',
  category: '',
  uploader: '',
  dateRange: [] as string[],
})

const reviewForm = reactive({
  status: 'approved' as 'approved' | 'rejected',
  comment: '',
})

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
    minWidth: 200,
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
  // 已移除文件大小列
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
    title: '审核人',
    dataIndex: 'reviewer',
    width: 120,
    render: ({ record }: { record: Material }) => {
      return record.reviewer || '-'
    },
  },
  {
    title: '审核时间',
    dataIndex: 'reviewTime',
    width: 180,
    render: ({ record }: { record: Material }) => {
      return record.reviewTime ? formatDate(record.reviewTime) : '-'
    },
  },
  {
    title: '操作',
    slotName: 'actions',
    width: userStore.role === 'admin' || userStore.role === 'user' ? 300 : 200,
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

// 响应式模态框宽度
const modalWidth = computed(() => {
  return window.innerWidth < 768 ? '95%' : '600px'
})
const drawerWidth = computed(() => {
  return window.innerWidth < 768 ? '95%' : '800px'
})

const fetchMaterials = async () => {
  loading.value = true
  try {
    const response = await getMaterialList({
      status: filterForm.status || undefined,
      category: filterForm.category || undefined,
      uploader: filterForm.uploader || undefined,
      startDate: filterForm.dateRange[0],
      endDate: filterForm.dateRange[1],
    })
    materials.value = response.data
    pagination.total = materials.value.length
  } catch (error) {
    Message.error('获取材料列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchMaterials()
}

const handleReset = () => {
  filterForm.status = ''
  filterForm.category = ''
  filterForm.uploader = ''
  filterForm.dateRange = []
  pagination.current = 1
  fetchMaterials()
}

const handlePageChange = (page: number) => {
  pagination.current = page
  fetchMaterials()
}

const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  fetchMaterials()
}

const handleView = (material: Material) => {
  currentMaterial.value = material
  showDetailModal.value = true
}

const handlePreviewDrawer = (material: Material) => {
  currentMaterial.value = material
  showPreviewDrawer.value = true
}

const handleEdit = (material: Material) => {
  currentMaterial.value = material
  showReviewModal.value = true
}

const handleDelete = async (material: Material) => {
  try {
    await deleteMaterial(material.id)
    Message.success('删除成功')
    fetchMaterials()
  } catch (error) {
    Message.error('删除失败')
  }
}

const handleBatchDelete = async () => {
  // 实现批量删除逻辑
  Message.info('批量删除功能待实现')
}

const handleExport = () => {
  // 实现导出逻辑
  Message.info('导出功能待实现')
}

const handleReview = (material: Material, status: 'approved' | 'rejected') => {
  currentMaterial.value = material
  reviewForm.status = status
  reviewForm.comment = ''
  showReviewModal.value = true
}

const handleConfirmReview = async () => {
  if (!currentMaterial.value) return

  try {
    await reviewMaterial({
      materialId: currentMaterial.value.id,
      status: reviewForm.status,
      comment: reviewForm.comment,
    })
    Message.success('审核成功')
    showReviewModal.value = false
    fetchMaterials()
  } catch (error) {
    Message.error('审核失败')
  }
}

onMounted(() => {
  fetchMaterials()
  // 监听来自 MaterialDetail 的预览事件
  window.addEventListener('openPreviewDrawer', (e: Event) => {
    const url = (e as CustomEvent).detail
    if (currentMaterial.value) {
      // 替换 iframe src 为点击的文件 URL
      currentMaterial.value = {
        ...currentMaterial.value,
        files: [{ ...currentMaterial.value.files[0], fileUrl: url }],
      }
    }
    showPreviewDrawer.value = true
  })
})
</script>

<style lang="less" scoped>
.material-list {
  padding: 20px;

  .list-card {
    min-height: 500px;
  }

  .filter-section {
    margin-bottom: 16px;

    .filter-form {
      :deep(.arco-form-item) {
        margin-bottom: 0;
        margin-right: 16px;
      }
    }
  }

  .action-section {
    margin-bottom: 16px;
  }
}

// 移动端适配
@media (max-width: 768px) {
  .material-list {
    padding: 12px;

    .filter-section {
      .filter-form {
        :deep(.arco-form-item) {
          margin-right: 0;
          margin-bottom: 12px;
          width: 100%;

          .arco-select,
          .arco-input,
          .arco-range-picker {
            width: 100% !important;
          }
        }
      }
    }

    .action-section {
      .arco-btn {
        width: 100%;
        margin-bottom: 8px;
      }
    }
  }
}
</style>
