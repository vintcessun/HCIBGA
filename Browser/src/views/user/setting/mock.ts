import setupMock, { successResponseWrap } from '@/utils/setup-mock'
import Mock from 'mockjs'

setupMock({
  setup() {
    Mock.mock(new RegExp('/api/user/save-info'), () => {
      return successResponseWrap('ok')
    })
    Mock.mock(new RegExp('/api/user/certification'), () => {
      return successResponseWrap({
        enterpriseInfo: {
          accountType: '学校账号',
          status: 0,
          time: '2018-10-22 14:53:12',
          legalPerson: '李**',
          certificateType: '身份证',
          authenticationNumber: '130************123',
          enterpriseName: '某某学校',
        },
        record: [
          {
            certificationType: 1,
            certificationContent: '学校认证，姓名：李**',
            status: 0,
            time: '2021-02-28 10:30:50',
          },
          {
            certificationType: 1,
            certificationContent: '学校认证，姓名：李**',
            status: 1,
            time: '2020-05-13 08:00:00',
          },
        ],
      })
    })
    Mock.mock(new RegExp('/api/user/upload'), () => {
      return successResponseWrap('ok')
    })
  },
})
