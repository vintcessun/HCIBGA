<template>
  <div class="bonus-list-page">
    <a-card class="general-card" :title="t('bonus.academic.title')">
      <a-table row-key="id" :loading="loading" :data="records" :pagination="false">
        <template #columns>
          <a-table-column :title="t('bonus.academic.column.project')" data-index="project" />
          <a-table-column :title="t('bonus.academic.column.awardDate')" data-index="awardDate" width="140" />
          <a-table-column :title="t('bonus.academic.column.awardLevel')" data-index="awardLevel" width="120" />
          <a-table-column :title="t('bonus.academic.column.awardType')" data-index="awardType" width="130" />
          <a-table-column :title="t('bonus.academic.column.teamRank')" data-index="teamRank" width="130" />
          <a-table-column :title="t('bonus.academic.column.selfScore')" data-index="selfScore" width="110" />
          <a-table-column :title="t('bonus.academic.column.scoreBasis')" data-index="scoreBasis" />
          <a-table-column :title="t('bonus.academic.column.collegeScore')" data-index="collegeScore" width="140" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { getAcademicBonusList, type BonusRecord } from '@/api/bonus'
import { Message } from '@arco-design/web-vue'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const loading = ref(false)
const records = ref<BonusRecord[]>([])

async function fetchRecords() {
  loading.value = true
  try {
    const { data } = await getAcademicBonusList()
    records.value = data
  } catch (error) {
    Message.error(t('bonus.academic.error.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(fetchRecords)
</script>

<style scoped>
.bonus-list-page {
  padding: 20px;
}
</style>
