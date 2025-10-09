<template>
  <div class="submit-create-page">
    <a-card class="general-card" title="发起提交">
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题">
          <a-input v-model="form.title" placeholder="请输入标题" />
        </a-form-item>
        <a-form-item label="简介">
          <a-textarea v-model="form.description" placeholder="请输入简介" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSubmit">提交</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'

const form = reactive({
  title: '',
  description: '',
})

function handleSubmit() {
  console.log('提交数据:', form)
  // TODO: 发送提交请求
}
</script>

<style scoped>
.submit-create-page {
  padding: 20px;
}
</style>
