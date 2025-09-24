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
      'category|1': ['document', 'image', 'video', 'audio', 'other'],
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
        category: data.category || 'document',
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
          document: mockMaterials.list.filter((item: MaterialItem) => item.category === 'document').length,
          image: mockMaterials.list.filter((item: MaterialItem) => item.category === 'image').length,
          video: mockMaterials.list.filter((item: MaterialItem) => item.category === 'video').length,
          audio: mockMaterials.list.filter((item: MaterialItem) => item.category === 'audio').length,
          other: mockMaterials.list.filter((item: MaterialItem) => item.category === 'other').length,
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
