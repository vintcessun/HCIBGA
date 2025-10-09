<template>
  <div class="review-submit-page">
    <a-card class="general-card" title="审核提交">
      <a-button type="primary" style="margin-bottom: 16px" @click="publishBatch">发布批次</a-button>
      <a-table :columns="columns" :data="data" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const data = ref([
  { title: '批次1', description: '简介1', status: '待审核' },
  { title: '批次2', description: '简介2', status: '已审核' },
])

const columns = [
  { title: '标题', dataIndex: 'title' },
  { title: '简介', dataIndex: 'description' },
  { title: '状态', dataIndex: 'status' },
]

function publishBatch() {
  // TODO: 发布批次逻辑
  console.log('发布批次')
}
</script>

<style scoped>
.review-submit-page {
  padding: 20px;
}
</style>
