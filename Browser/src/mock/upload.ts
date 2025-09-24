import setupMock, { successResponseWrap } from '@/utils/setup-mock'
import Mock from 'mockjs'

interface UploadedFile {
  id: string
  name: string
  size: number
  url: string
  md5: string
  upload_time: string
  uploader: string
}

// 模拟已上传文件数据
const mockUploadedFiles = Mock.mock({
  'list|10-30': [
    {
      id: '@guid',
      name: '@word(3, 10).@extension',
      'size|1024-10485760': 1, // 1KB - 10MB
      url: '@url',
      md5: '@string("hex", 32)',
      upload_time: '@datetime',
      uploader: '@cname',
    },
  ],
})

setupMock({
  setup() {
    // 检查文件是否存在 - 总是返回不存在
    Mock.mock(new RegExp('/api/upload/check'), () => {
      return successResponseWrap({
        exists: false,
      })
    })

    // 上传单个文件 - 总是返回成功
    Mock.mock(new RegExp('/api/upload/file'), () => {
      const newFile = {
        id: Mock.mock('@guid'),
        name: 'uploaded-file.pdf',
        size: Mock.mock('@integer(1024, 10485760)'),
        url: Mock.mock('@url'),
        md5: Mock.mock('@string("hex", 32)'),
        upload_time: new Date().toISOString(),
        uploader: '当前用户',
      }

      mockUploadedFiles.list.push(newFile)

      return successResponseWrap({
        file_id: newFile.id,
        url: newFile.url,
        md5: newFile.md5,
      })
    })

    // 批量上传文件 - 总是返回成功
    Mock.mock(new RegExp('/api/upload/batch'), () => {
      const results = []

      // 模拟上传2个文件
      for (let index = 0; index < 2; index += 1) {
        const newFile = {
          id: Mock.mock('@guid'),
          name: `uploaded-file-${index + 1}.pdf`,
          size: Mock.mock('@integer(1024, 10485760)'),
          url: Mock.mock('@url'),
          md5: Mock.mock('@string("hex", 32)'),
          upload_time: new Date().toISOString(),
          uploader: '当前用户',
        }

        mockUploadedFiles.list.push(newFile)
        results.push({
          file_id: newFile.id,
          url: newFile.url,
          md5: newFile.md5,
          status: 'uploaded',
        })
      }

      return successResponseWrap(results)
    })

    // 获取文件列表
    Mock.mock(new RegExp('/api/upload/files'), () => {
      return successResponseWrap(mockUploadedFiles.list)
    })
  },
})
