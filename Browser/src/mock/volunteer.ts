import setupMock, { successResponseWrap } from '@/utils/setup-mock'
import Mock from 'mockjs'

setupMock({
  setup() {
    const mockResponse = {
      credit_hours: 135.29,
      honor_hours: 80.4,
      total_hours: 215.69,
    }

    Mock.mock(new RegExp('/api/volunteer/credit'), () => {
      return successResponseWrap(mockResponse)
    })
  },
})
