<template>
  <div class="bonus-summary-page">
    <a-space direction="vertical" size="large" style="width: 100%">
      <a-card class="general-card" :title="$t('bonus.summary.title')">
        <a-table row-key="category" :data="items" :loading="loading" :pagination="false">
          <template #columns>
            <a-table-column :title="$t('bonus.summary.column.category')" data-index="category">
              <template #cell="{ record }">
                {{ getCategoryLabel(record.category) }}
              </template>
            </a-table-column>
            <a-table-column :title="$t('bonus.summary.column.itemCount')" data-index="itemCount" width="120" />
            <a-table-column :title="$t('bonus.summary.column.totalScore')" data-index="totalScore" width="120" />
          </template>
        </a-table>
      </a-card>
      <a-result v-if="!loading" status="success" :title="$t('bonus.summary.total.title')">
        <template #subtitle>{{ $t('bonus.summary.total.subtitle', { score: summary?.totalScore ?? 0 }) }}</template>
      </a-result>
    </a-space>
  </div>
</template>

<script lang="ts" setup>
import { getBonusSummary, type BonusSummaryItem, type BonusSummaryResponse } from '@/api/bonus'
import { Message } from '@arco-design/web-vue'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const loading = ref(false)
const items = ref<BonusSummaryItem[]>([])
const summary = ref<BonusSummaryResponse | null>(null)

async function fetchSummary() {
  loading.value = true
  try {
    const { data } = await getBonusSummary()
    summary.value = data
    items.value = data.items
  } catch (error) {
    Message.error(t('bonus.summary.error.fetchFailed'))
  } finally {
    loading.value = false
  }
}

function getCategoryLabel(category: string) {
  if (category === '学术专长') {
    return t('bonus.category.academic')
  }
  if (category === '综合素质') {
    return t('bonus.category.comprehensive')
  }
  return category
}

onMounted(fetchSummary)
</script>

<style scoped>
.bonus-summary-page {
  padding: 20px;
}
</style>
