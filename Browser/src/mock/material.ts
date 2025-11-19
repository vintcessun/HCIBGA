import { MockParams } from '@/types/mock'
import setupMock, { successResponseWrap } from '@/utils/setup-mock'
import Mock from 'mockjs'

interface MaterialItem {
  id: string
  title: string
  description: string
  category: string
  tags: string[]
  fileUrl: string
  fileName: string
  fileSize: number
  status: string
  uploader: string
  uploadTime: string
  reviewer?: string
  reviewTime?: string
  reviewComment?: string
  aiReviewResult?: any
}

// 模拟材料数据
const mockMaterials = Mock.mock({
  'list|20-50': [
    {
      id: '@guid',
      title: '@ctitle(5, 20)',
      description: '@cparagraph(1, 3)',
      'category|1': [
        '学术专长成绩-科研成果',
        '学术专长成绩-学业竞赛',
        '学术专长成绩-创新创业训练',
        '综合表现加分-国际组织实习',
        '综合表现加分-参军入伍服兵役',
        '综合表现加分-志愿服务',
        '综合表现加分-荣誉称号',
        '综合表现加分-社会工作',
        '综合表现加分-体育比赛',
      ],
      'tags|1-3': ['important', 'urgent', 'normal'],
      fileUrl: '@url',
      fileName: '@word(3, 10).@extension',
      'fileSize|1024-10485760': 1, // 1KB - 10MB
      'status|1': ['pending', 'approved', 'rejected'],
      uploader: '@cname',
      uploadTime: '@datetime',
      reviewer(this: any) {
        return this['status|1'] !== 'pending' ? '@cname' : undefined
      },
      reviewTime(this: any) {
        return this['status|1'] !== 'pending' ? '@datetime' : undefined
      },
      reviewComment(this: any) {
        return this['status|1'] === 'rejected' ? '@csentence' : undefined
      },
      'aiReviewResult|1': function aiReviewResult() {
        if (Math.random() > 0.7) {
          return {
            'score|60-95': 1,
            'confidence|0.7-0.95': 0.01,
            'suggestions|1-3': ['@csentence'],
            'riskLevel|1': ['low', 'medium', 'high'],
          }
        }
        return undefined
      },
    },
  ],
})

setupMock({
  setup() {
    // 上传材料 - 总是返回成功
    Mock.mock(new RegExp('/api/material/upload'), (params: MockParams) => {
      const data = JSON.parse(params.body)

      const newMaterial = {
        id: Mock.mock('@guid'),
        title: data.title || '示例材料',
        description: data.description || '这是一个示例材料描述',
        category: data.category || '学术专长成绩-科研成果',
        tags: data.tags || ['normal'],
        fileUrl: Mock.mock('@url'),
        fileName: `material-${Date.now()}.pdf`,
        fileSize: Mock.mock('@integer(1024, 10485760)'),
        status: 'pending',
        uploader: '当前用户',
        uploadTime: new Date().toISOString(),
        reviewer: undefined,
        reviewTime: undefined,
        reviewComment: undefined,
        aiReviewResult: undefined,
      }
      mockMaterials.list.push(newMaterial)
      return successResponseWrap(newMaterial)
    })

    // 检查文件是否存在
    Mock.mock(new RegExp('/api/upload/check'), (params: MockParams) => {
      const { md5, filename } = JSON.parse(params.body)
      const existsFile = mockMaterials.list.find((item: MaterialItem) => item.fileName === filename)
      if (existsFile) {
        return successResponseWrap({
          exists: true,
          fileId: existsFile.id,
          url: existsFile.fileUrl,
        })
      }
      return successResponseWrap({
        exists: false,
        fileId: '',
        url: '',
      })
    })

    // 上传文件
    Mock.mock(new RegExp('/api/upload/file'), () => {
      return successResponseWrap({
        fileId: Mock.mock('@guid'),
        url: Mock.mock('@url'),
        md5: Mock.mock('@string("lower", 32)'),
      })
    })

    // 调用LLM自动填写材料信息
    Mock.mock(new RegExp('/api/material/llm-fill'), (params: MockParams) => {
      const { fileId, metadata } = JSON.parse(params.body)
      return successResponseWrap({
        title: metadata?.title || '项目申报材料（自动生成）',
        description: metadata?.description || '包含项目预算与团队信息，由LLM自动生成。',
        category: metadata?.category || '学术专长成绩-科研成果',
        tags: metadata?.tags?.length ? metadata.tags : ['自动生成', '待确认'],
      })
    })

    // 获取材料列表
    Mock.mock(new RegExp('/api/material/list'), (params: MockParams) => {
      const filter = JSON.parse(params.body)
      let filteredList = [...mockMaterials.list]

      if (filter.status) {
        filteredList = filteredList.filter((item: MaterialItem) => item.status === filter.status)
      }
      if (filter.category) {
        filteredList = filteredList.filter((item: MaterialItem) => item.category === filter.category)
      }
      if (filter.uploader) {
        filteredList = filteredList.filter((item: MaterialItem) => item.uploader.includes(filter.uploader))
      }
      if (filter.startDate && filter.endDate) {
        filteredList = filteredList.filter((item: MaterialItem) => {
          const uploadTime = new Date(item.uploadTime)
          return uploadTime >= new Date(filter.startDate) && uploadTime <= new Date(filter.endDate)
        })
      }

      return successResponseWrap(filteredList)
    })

    // 获取待审核材料
    Mock.mock(new RegExp('/api/material/pending'), () => {
      const pendingMaterials = mockMaterials.list.filter((item: MaterialItem) => item.status === 'pending')
      return successResponseWrap(pendingMaterials)
    })

    // 审核材料
    Mock.mock(new RegExp('/api/material/review'), (params: MockParams) => {
      const { materialId, status, comment } = JSON.parse(params.body)
      const material = mockMaterials.list.find((item: MaterialItem) => item.id === materialId)

      if (material) {
        material.status = status
        material.reviewer = '审核员'
        material.reviewTime = new Date().toISOString()
        material.reviewComment = comment
      }

      return successResponseWrap(material)
    })

    // 获取材料统计
    Mock.mock(new RegExp('/api/material/statistics'), () => {
      const stats = {
        total: mockMaterials.list.length,
        pending: mockMaterials.list.filter((item: MaterialItem) => item.status === 'pending').length,
        approved: mockMaterials.list.filter((item: MaterialItem) => item.status === 'approved').length,
        rejected: mockMaterials.list.filter((item: MaterialItem) => item.status === 'rejected').length,
        byCategory: {
          '学术专长成绩-科研成果': mockMaterials.list.filter((item: MaterialItem) => item.category === '学术专长成绩-科研成果').length,
          '学术专长成绩-学业竞赛': mockMaterials.list.filter((item: MaterialItem) => item.category === '学术专长成绩-学业竞赛').length,
          '学术专长成绩-创新创业训练': mockMaterials.list.filter((item: MaterialItem) => item.category === '学术专长成绩-创新创业训练')
            .length,
          '综合表现加分-国际组织实习': mockMaterials.list.filter((item: MaterialItem) => item.category === '综合表现加分-国际组织实习')
            .length,
          '综合表现加分-参军入伍服兵役': mockMaterials.list.filter((item: MaterialItem) => item.category === '综合表现加分-参军入伍服兵役')
            .length,
          '综合表现加分-志愿服务': mockMaterials.list.filter((item: MaterialItem) => item.category === '综合表现加分-志愿服务').length,
          '综合表现加分-荣誉称号': mockMaterials.list.filter((item: MaterialItem) => item.category === '综合表现加分-荣誉称号').length,
          '综合表现加分-社会工作': mockMaterials.list.filter((item: MaterialItem) => item.category === '综合表现加分-社会工作').length,
          '综合表现加分-体育比赛': mockMaterials.list.filter((item: MaterialItem) => item.category === '综合表现加分-体育比赛').length,
        },
      }
      return successResponseWrap(stats)
    })

    // AI自动审核
    Mock.mock(new RegExp('/api/material/ai-review'), (params: MockParams) => {
      const { materialId } = JSON.parse(params.body)
      const material = mockMaterials.list.find((item: MaterialItem) => item.id === materialId)

      if (material) {
        const aiResult = {
          score: Mock.mock('@integer(60, 95)'),
          confidence: Mock.mock('@float(0.7, 0.95, 2, 2)'),
          suggestions: [Mock.mock('@csentence'), Mock.mock('@csentence')],
          riskLevel: Mock.mock('@pick(["low", "medium", "high"])'),
        }

        material.aiReviewResult = aiResult
        return successResponseWrap(aiResult)
      }

      return successResponseWrap(null)
    })

    // 删除材料
    Mock.mock(new RegExp('/api/material/delete'), (params: MockParams) => {
      const { materialId } = JSON.parse(params.body)
      const index = mockMaterials.list.findIndex((item: MaterialItem) => item.id === materialId)

      if (index !== -1) {
        mockMaterials.list.splice(index, 1)
      }

      return successResponseWrap(null)
    })

    // 批量审核
    Mock.mock(new RegExp('/api/material/batch-review'), (params: MockParams) => {
      const { materialIds, status, comment } = JSON.parse(params.body)

      materialIds.forEach((materialId: string) => {
        const material = mockMaterials.list.find((item: MaterialItem) => item.id === materialId)
        if (material) {
          material.status = status
          material.reviewer = '批量审核'
          material.reviewTime = new Date().toISOString()
          material.reviewComment = comment || `批量${status === 'approved' ? '通过' : '拒绝'}`
        }
      })

      return successResponseWrap(null)
    })
  },
})
