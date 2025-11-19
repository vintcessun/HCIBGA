<template>
  <div class="material-statistics">
    <a-card :title="$t('menu.material.statistics')" class="statistics-card">
      <!-- 统计概览 -->
      <div class="overview-section">
        <a-row :gutter="16">
          <a-col :xs="12" :sm="6" v-for="stat in overviewStats" :key="stat.label">
            <a-card class="stat-card" :bordered="false">
              <div class="stat-content">
                <div class="stat-value" :style="{ color: stat.color }">
                  {{ stat.value }}
                </div>
                <div class="stat-label">
                  {{ stat.label }}
                </div>
                <div class="stat-trend" v-if="stat.trend">
                  <icon-caret-up v-if="stat.trend === 'up'" style="color: #00b42a" />
                  <icon-caret-down v-else style="color: #f53f3f" />
                  <span :style="{ color: stat.trend === 'up' ? '#00b42a' : '#f53f3f' }">{{ stat.percent }}%</span>
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>
      </div>

      <!-- 根据角色显示不同标题 -->
      <div class="category-section">
        <a-card
          :title="userStore.role === 'user' ? $t('material.statistics.bonus') : $t('material.statistics.category')"
          class="category-card"
        >
          <a-row :gutter="16">
            <a-col :xs="24" :lg="12">
              <div class="chart-container">
                <v-chart :option="categoryChartOption" autoresize />
              </div>
            </a-col>
            <a-col :xs="24" :lg="12">
              <div class="category-list">
                <div v-for="category in categoryStats" :key="category.name" class="category-item">
                  <div class="category-info">
                    <div class="category-name">{{ category.name }}</div>
                    <div class="category-count">{{ category.value }} {{ $t('material.statistics.unit') }}</div>
                  </div>
                  <div class="category-percent">{{ category.percent }}%</div>
                </div>
              </div>
            </a-col>
          </a-row>
        </a-card>
      </div>

      <!-- 时间趋势 -->
      <div class="trend-section">
        <a-card :title="$t('material.statistics.trend')" class="trend-card">
          <v-chart :option="trendChartOption" autoresize style="height: 300px" />
        </a-card>
      </div>

      <!-- 审核效率 -->
      <div class="efficiency-section">
        <a-card :title="$t('material.statistics.efficiency')" class="efficiency-card">
          <a-row :gutter="16">
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.averageReviewTime }} {{ $t('material.statistics.hours') }}</div>
                <div class="efficiency-label">{{ $t('material.statistics.avgReviewTime') }}</div>
              </div>
            </a-col>
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.reviewCount }}</div>
                <div class="efficiency-label">{{ $t('material.statistics.todayReviewCount') }}</div>
              </div>
            </a-col>
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.approvalRate }}%</div>
                <div class="efficiency-label">{{ $t('material.statistics.approvalRate') }}</div>
              </div>
            </a-col>
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.rejectionRate }}%</div>
                <div class="efficiency-label">{{ $t('material.statistics.rejectionRate') }}</div>
              </div>
            </a-col>
          </a-row>
        </a-card>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { getMaterialStatistics } from '@/api/material'
import { useUserStore } from '@/store'
import { LineChart, PieChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TitleComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { onMounted, ref } from 'vue'
import VChart from 'vue-echarts'
import { useI18n } from 'vue-i18n'

use([CanvasRenderer, PieChart, LineChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent])

interface OverviewStat {
  label: string
  value: number
  color: string
  trend?: 'up' | 'down'
  percent?: number
}

interface CategoryStat {
  name: string
  value: number
  percent: number
}

interface EfficiencyStat {
  averageReviewTime: number
  reviewCount: number
  approvalRate: number
  rejectionRate: number
}

const { t } = useI18n()
const userStore = useUserStore()
const overviewStats = ref<OverviewStat[]>([])
const categoryStats = ref<CategoryStat[]>([])
const efficiencyStats = ref<EfficiencyStat>({
  averageReviewTime: 2.5,
  reviewCount: 0,
  approvalRate: 0,
  rejectionRate: 0,
})

const categoryChartOption = ref({})
const trendChartOption = ref({})

const getCategoryName = (category: string): string => {
  const categoryMap: Record<string, string> = {
    '学术专长成绩-科研成果': t('material.category.academic.research'),
    '学术专长成绩-学业竞赛': t('material.category.academic.competition'),
    '学术专长成绩-创新创业训练': t('material.category.academic.innovation'),
    '综合表现加分-国际组织实习': t('material.category.comprehensive.internship'),
    '综合表现加分-参军入伍服兵役': t('material.category.comprehensive.military'),
    '综合表现加分-志愿服务': t('material.category.comprehensive.volunteer'),
    '综合表现加分-荣誉称号': t('material.category.comprehensive.honor'),
    '综合表现加分-社会工作': t('material.category.comprehensive.social'),
    '综合表现加分-体育比赛': t('material.category.comprehensive.sports'),
  }
  return categoryMap[category] || category
}

const updateCategoryChart = () => {
  categoryChartOption.value = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      data: categoryStats.value.map((item) => item.name),
    },
    series: [
      {
        name: t('material.statistics.chart.category'),
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2,
        },
        label: {
          show: false,
          position: 'center',
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold',
          },
        },
        labelLine: {
          show: false,
        },
        data: categoryStats.value.map((item) => ({
          value: item.value,
          name: item.name,
        })),
      },
    ],
  }
}

const updateTrendChart = () => {
  // 模拟最近7天的上传数据
  const dates = Array.from({ length: 7 }, (_, i) => {
    const date = new Date()
    date.setDate(date.getDate() - i)
    return date.toLocaleDateString('zh-CN')
  }).reverse()

  const uploadData = dates.map(() => Math.floor(Math.random() * 20) + 5)

  trendChartOption.value = {
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      data: dates,
    },
    yAxis: {
      type: 'value',
      name: t('material.statistics.chart.uploadCount'),
    },
    series: [
      {
        name: t('material.statistics.chart.dailyUpload'),
        type: 'line',
        data: uploadData,
        smooth: true,
        lineStyle: {
          width: 3,
          color: '#165dff',
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(22, 93, 255, 0.6)' },
              { offset: 1, color: 'rgba(22, 93, 255, 0.1)' },
            ],
          },
        },
      },
    ],
  }
}

const fetchStatistics = async () => {
  try {
    const response = await getMaterialStatistics()
    const stats = response.data

    // 更新概览统计
    overviewStats.value = [
      {
        label: t('material.statistics.total'),
        value: stats.total,
        color: '#165dff',
      },
      {
        label: t('material.statistics.pending'),
        value: stats.pending,
        color: '#ff7d00',
        trend: stats.pending > 10 ? 'up' : 'down',
        percent: Math.round((stats.pending / stats.total) * 100),
      },
      {
        label: t('material.statistics.approved'),
        value: stats.approved,
        color: '#00b42a',
        trend: 'up',
        percent: Math.round((stats.approved / stats.total) * 100),
      },
      {
        label: t('material.statistics.rejected'),
        value: stats.rejected,
        color: '#f53f3f',
        trend: 'down',
        percent: Math.round((stats.rejected / stats.total) * 100),
      },
    ]

    // 更新分类统计
    const categoryValues = Object.values(stats.byCategory) as number[]
    const totalByCategory = categoryValues.reduce((sum, count) => sum + count, 0)
    categoryStats.value = Object.entries(stats.byCategory).map(([name, value]) => ({
      name: getCategoryName(name),
      value: value as number,
      percent: Math.round(((value as number) / totalByCategory) * 100),
    }))

    // 更新分类图表
    updateCategoryChart()

    // 更新趋势图表
    updateTrendChart()

    // 更新审核效率
    efficiencyStats.value = {
      ...efficiencyStats.value,
      reviewCount: stats.pending + stats.approved + stats.rejected,
      approvalRate: Math.round((stats.approved / (stats.approved + stats.rejected || 1)) * 100),
      rejectionRate: Math.round((stats.rejected / (stats.approved + stats.rejected || 1)) * 100),
    }
  } catch (error) {
    // 获取统计信息失败
  }
}

onMounted(() => {
  fetchStatistics()
})
</script>

<style lang="less" scoped>
.material-statistics {
  padding: 20px;

  .statistics-card {
    min-height: 500px;
  }

  .overview-section {
    margin-bottom: 24px;

    .stat-card {
      text-align: center;
      background: linear-gradient(135deg, #f7f8fa 0%, #e5e6eb 100%);
      border-radius: 8px;

      .stat-content {
        padding: 16px;

        .stat-value {
          font-size: 24px;
          font-weight: bold;
          margin-bottom: 8px;
        }

        .stat-label {
          color: #86909c;
          font-size: 14px;
          margin-bottom: 4px;
        }

        .stat-trend {
          font-size: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 4px;
        }
      }
    }
  }

  .category-section {
    margin-bottom: 24px;

    .chart-container {
      height: 300px;
    }

    .category-list {
      .category-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 0;
        border-bottom: 1px solid #e5e6eb;

        &:last-child {
          border-bottom: none;
        }

        .category-info {
          display: flex;
          align-items: center;
          gap: 8px;

          .category-name {
            font-weight: 500;
          }

          .category-count {
            color: #86909c;
            font-size: 14px;
          }
        }

        .category-percent {
          color: #165dff;
          font-weight: bold;
        }
      }
    }
  }

  .efficiency-section {
    .efficiency-item {
      text-align: center;
      padding: 16px;

      .efficiency-value {
        font-size: 24px;
        font-weight: bold;
        color: #165dff;
        margin-bottom: 8px;
      }

      .efficiency-label {
        color: #86909c;
        font-size: 14px;
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .material-statistics {
    padding: 12px;

    .overview-section {
      .stat-card {
        margin-bottom: 12px;
      }
    }

    .category-section {
      .category-list {
        margin-top: 16px;
      }
    }

    .efficiency-section {
      .efficiency-item {
        margin-bottom: 12px;
      }
    }
  }
}
</style>
