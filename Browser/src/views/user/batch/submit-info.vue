<template>
  <div class="submit-info-page">
    <a-card class="general-card" title="提交批次材料">
      <a-table :columns="columns" :data="data" :row-selection="{ type: 'checkbox' }" />
    </a-card>
    <a-card class="general-card" title="提交批次">
      <div style="display: flex; align-items: center">
        <a-select v-model="selectedBatch" placeholder="选择批次" style="width: 200px; margin-right: 10px">
          <a-option v-for="item in data" :key="item.title" :value="item.title">{{ item.title }}</a-option>
        </a-select>
        <a-button type="primary" @click="handleSubmit">提交</a-button>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const data = ref([
  { title: '批次1', description: '简介1', status: '待提交' },
  { title: '批次2', description: '简介2', status: '已提交' },
])

const selectedBatch = ref('')

const columns = [
  { title: '标题', dataIndex: 'title' },
  { title: '简介', dataIndex: 'description' },
  { title: '状态', dataIndex: 'status' },
]

function handleSubmit() {
  if (!selectedBatch.value) {
    alert('请选择批次')
    return false
  }
  alert(`已提交批次: ${selectedBatch.value}`)
  return true
}
</script>

<style scoped>
.submit-info-page {
  padding: 20px;
}
</style>
