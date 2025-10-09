<template>
  <div class="submit-list-page">
    <a-card class="general-card" title="提交列表">
      <a-table :columns="columns" :data="data" row-selection="multiple">
        <template #columns>
          <a-table-column title="我的材料" data-index="myMaterial" />
          <a-table-column v-for="col in otherColumns" :key="col.dataIndex" v-bind="col" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

// 表格数据
const data = ref([
  { myMaterial: '材料1', category: '类别A', status: '已提交' },
  { myMaterial: '材料2', category: '类别B', status: '已提交' },
])

// 不含“操作”列
const otherColumns = [
  { title: '类别', dataIndex: 'category' },
  { title: '状态', dataIndex: 'status' },
]

const columns = [{ title: '我的材料', dataIndex: 'myMaterial' }, ...otherColumns]
</script>

<style scoped>
.submit-list-page {
  padding: 20px;
}
</style>
