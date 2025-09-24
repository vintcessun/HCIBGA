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

      <!-- 分类统计 -->
      <div class="category-section">
        <a-card title="分类统计" class="category-card">
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
                    <div class="category-count">{{ category.value }} 个</div>
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
        <a-card title="上传趋势" class="trend-card">
          <v-chart :option="trendChartOption" autoresize style="height: 300px" />
        </a-card>
      </div>

      <!-- 审核效率 -->
      <div class="efficiency-section">
        <a-card title="审核效率" class="efficiency-card">
          <a-row :gutter="16">
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.averageReviewTime }}</div>
                <div class="efficiency-label">平均审核时间</div>
              </div>
            </a-col>
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.reviewCount }}</div>
                <div class="efficiency-label">今日审核数量</div>
              </div>
            </a-col>
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.approvalRate }}%</div>
                <div class="efficiency-label">通过率</div>
              </div>
            </a-col>
            <a-col :xs="12" :lg="6">
              <div class="efficiency-item">
                <div class="efficiency-value">{{ efficiencyStats.rejectionRate }}%</div>
                <div class="efficiency-label">拒绝率</div>
              </div>
            </a-col>
          </a-row>
        </a-card>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { getMaterialStatistics, type Material } from '@/api/material'

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
  averageReviewTime: string
  reviewCount: number
  approvalRate: number
  rejectionRate: number
}

const overviewStats = ref<OverviewStat[]>([])
const categoryStats = ref<CategoryStat[]>([])
const efficiencyStats = ref<EfficiencyStat>({
  averageReviewTime: '2.5小时',
  reviewCount: 0,
  approvalRate: 0,
  rejectionRate: 0,
})

const categoryChartOption = ref({})
const trendChartOption = ref({})

const getCategoryName = (key: string): string => {
  const categoryMap: Record<string, string> = {
    document: '文档',
    image: '图片',
    video: '视频',
    audio: '音频',
    other: '其他',
  }
  return categoryMap[key] || key
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
        name: '材料分类',
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
      name: '上传数量',
    },
    series: [
      {
        name: '每日上传',
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
        label: '总材料数',
        value: stats.total,
        color: '#165dff',
      },
      {
        label: '待审核',
        value: stats.pending,
        color: '#ff7d00',
        trend: stats.pending > 10 ? 'up' : 'down',
        percent: Math.round((stats.pending / stats.total) * 100),
      },
      {
        label: '已通过',
        value: stats.approved,
        color: '#00b42a',
        trend: 'up',
        percent: Math.round((stats.approved / stats.total) * 100),
      },
      {
        label: '已拒绝',
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
