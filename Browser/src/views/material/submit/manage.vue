<template>
  <div class="manage-submit-page">
    <a-card class="general-card" title="管理提交">
      <a-button type="primary" style="margin-bottom: 16px" @click="createBatch">新建批次</a-button>
      <a-table :columns="columns" :data="data" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const data = ref([
  { title: '批次1', description: '简介1' },
  { title: '批次2', description: '简介2' },
])

const columns = [
  { title: '标题', dataIndex: 'title' },
  { title: '简介', dataIndex: 'description' },
]

function createBatch() {
  // TODO: 打开新建批次表单
  console.log('新建批次')
}
</script>

<style scoped>
.manage-submit-page {
  padding: 20px;
}
</style>
