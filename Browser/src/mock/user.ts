import setupMock, { failResponseWrap, successResponseWrap } from '@/utils/setup-mock'
import Mock from 'mockjs'

import { MockParams } from '@/types/mock'
import { isLogin } from '@/utils/auth'

setupMock({
  setup() {
    // Mock.XHR.prototype.withCredentials = true;

    // 用户信息
    Mock.mock(new RegExp('/api/user/info'), () => {
      if (isLogin()) {
        const role = window.localStorage.getItem('userRole') || 'admin'
        return successResponseWrap({
          name: '王立群',
          avatar: 'https://i.gtimg.cn/club/item/face/img/2/15922_100.gif',
          email: 'wangliqun@email.com',
          job: 'frontend',
          jobName: '前端艺术家',
          organization: 'Frontend',
          organizationName: '前端',
          location: 'beijing',
          locationName: '北京',
          introduction: '人潇洒，性温存',
          personalWebsite: 'https://www.arco.design',
          phone: '150****0000',
          registrationDate: '2013-05-10 12:10:00',
          accountId: '15012312300',
          certification: 1,
          role,
        })
      }
      return failResponseWrap(null, '未登录', 50008)
    })

    // 登录 - Mock数据
    Mock.mock(new RegExp('/api/user/login'), (params: MockParams) => {
      const { username, password } = JSON.parse(params.body)

      if (!username) {
        return {
          code: 50000,
          status: 'fail',
          msg: '用户名不能为空',
          data: null,
        }
      }
      if (!password) {
        return {
          code: 50000,
          status: 'fail',
          msg: '密码不能为空',
          data: null,
        }
      }

      // 支持三个用户：user、reviewer、admin，密码与用户名相同
      const validUsers = ['user', 'reviewer', 'admin']

      if (!validUsers.includes(username) || password !== username) {
        return {
          code: 50000,
          status: 'fail',
          msg: '用户名或密码错误',
          data: null,
        }
      }

      // 存储用户角色信息
      window.localStorage.setItem('userRole', username)

      // 返回成功的登录响应
      return {
        code: 20000,
        status: 'ok',
        msg: '登录成功',
        data: {
          token: `fake-jwt-token-for-${username}`,
          role: username,
          username,
          message: '登录成功',
        },
      }
    })

    // 登出
    Mock.mock(new RegExp('/api/user/logout'), () => {
      return successResponseWrap(null)
    })

    // 用户的服务端菜单
    Mock.mock(new RegExp('/api/user/menu'), () => {
      const menuList = [
        {
          path: '/material',
          name: 'material',
          meta: {
            locale: 'menu.material',
            requiresAuth: true,
            icon: 'icon-file',
            order: 1,
          },
          children: [
            {
              path: 'upload',
              name: 'MaterialUpload',
              meta: {
                locale: 'menu.material.upload',
                requiresAuth: true,
              },
            },
            {
              path: 'list',
              name: 'MaterialList',
              meta: {
                locale: 'menu.material.list',
                requiresAuth: true,
              },
            },
            {
              path: 'review',
              name: 'MaterialReview',
              meta: {
                locale: 'menu.material.review',
                requiresAuth: true,
                roles: ['admin', 'reviewer'],
              },
            },
            {
              path: 'statistics',
              name: 'MaterialStatistics',
              meta: {
                locale: 'menu.material.statistics',
                requiresAuth: true,
                roles: ['admin'],
              },
            },
          ],
        },
      ]
      return successResponseWrap(menuList)
    })
  },
})
