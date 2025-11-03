<template>
  <div class="info-export">
    <!-- 信息选择表格 -->
    <a-card class="general-card" :title="'信息选择'">
      <a-table :columns="columns" :data="tableData" row-key="id" :pagination="false" bordered :row-selection="rowSelection">
        <template #name="{ record }">
          <span>{{ record.name }}</span>
        </template>
        <template #studentId="{ record }">
          <span>{{ record.studentId }}</span>
        </template>
        <template #major="{ record }">
          <span>{{ record.major }}</span>
        </template>
        <template #class="{ record }">
          <span>{{ record.class }}</span>
        </template>
        <template #score="{ record }">
          <span>{{ record.score }}</span>
        </template>
      </a-table>
    </a-card>

    <!-- 学生信息导出方框 Excel -->
    <a-card class="general-card" :title="'学生信息导出为 Excel'" style="margin-top: 16px">
      <template #extra>
        <a-link>帮助</a-link>
      </template>
      <div class="upload-section">
        <div class="upload-area" @click="handleDownloadExcel">
          <icon-download />
          <div class="upload-text">点击下载 Excel 文件</div>
          <div class="upload-hint">支持 .xls, .xlsx 格式</div>
        </div>
      </div>
    </a-card>

    <!-- 学生信息导出方框 TXT -->
    <a-card class="general-card" :title="'学生信息导出为 TXT'" style="margin-top: 16px">
      <template #extra>
        <a-link>帮助</a-link>
      </template>
      <div class="upload-section">
        <div class="upload-area" @click="handleDownloadTxt">
          <icon-download />
          <div class="upload-text">点击下载 TXT 文件</div>
          <div class="upload-hint">支持 .txt 格式</div>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { Message } from '@arco-design/web-vue'
import { ref } from 'vue'

const columns = [
  { title: '姓名', dataIndex: 'name', slotName: 'name' },
  { title: '学号', dataIndex: 'studentId', slotName: 'studentId' },
  { title: '专业', dataIndex: 'major', slotName: 'major' },
  { title: '班级', dataIndex: 'class', slotName: 'class' },
  { title: '加分', dataIndex: 'score', slotName: 'score' },
]

const tableData = ref([
  { id: 1, name: '张三', studentId: '20230001', major: '计算机', class: '1班', score: 5 },
  { id: 2, name: '李四', studentId: '20230002', major: '数学', class: '2班', score: 3 },
])

const rowSelection = {
  type: 'checkbox',
  showCheckedAll: true,
}

const handleDownloadExcel = () => {
  // TODO: 实现Excel文件生成与下载
  Message.success('开始下载Excel导出文件...')
}

const handleDownloadTxt = () => {
  // TODO: 实现txt文件生成与下载
  Message.success('开始下载TXT导出文件...')
}
</script>

<style scoped lang="less">
.info-export {
  text-align: center;
  padding: 20px;

  .upload-section {
    text-align: center;
    margin-bottom: 24px;
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
