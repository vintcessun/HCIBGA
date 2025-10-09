<template>
  <div class="manage-batch-page">
    <a-card class="general-card" title="管理批次">
      <a-table :columns="columns" :data="data" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const data = ref([
  { batch: '批次1', total: 10, status: '进行中' },
  { batch: '批次2', total: 15, status: '已完成' },
])

const columns = [
  { title: '批次名称', dataIndex: 'batch' },
  { title: '总数', dataIndex: 'total' },
  { title: '状态', dataIndex: 'status' },
]
</script>

<style scoped>
.manage-batch-page {
  padding: 20px;
}
</style>
