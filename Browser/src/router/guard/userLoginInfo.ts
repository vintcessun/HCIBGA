import NProgress from 'nprogress' // progress bar
import type { LocationQueryRaw, Router } from 'vue-router'

import { useUserStore } from '@/store'
import { isLogin } from '@/utils/auth'

export default function setupUserLoginInfoGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    console.log('[UserLoginInfoGuard] 进入路由守卫 from:', from.fullPath, 'to:', to.fullPath)
    NProgress.start()
    const userStore = useUserStore()
    console.log('[UserLoginInfoGuard] 当前用户状态:', userStore.$state)
    if (isLogin()) {
      console.log('[UserLoginInfoGuard] 检测到已登录')
      if (userStore.role) {
        console.log('[UserLoginInfoGuard] 已有角色:', userStore.role)
        next()
      } else {
        try {
          console.log('[UserLoginInfoGuard] 用户未设置角色，调用 userStore.info() 获取用户信息')
          // eslint-disable-next-line no-console
          console.info('[route-guard] 用户未设置role，调用userStore.info() 获取用户信息...')
          await userStore.info()
          // eslint-disable-next-line no-console
          console.info('[route-guard] userStore.info() 返回用户信息:', JSON.stringify(userStore.userInfo, null, 2))
          console.log('[UserLoginInfoGuard] 获取用户信息成功，调用 next()')
          next()
        } catch (error) {
          console.error('[UserLoginInfoGuard] 获取用户信息失败:', error)
          // eslint-disable-next-line no-console
          console.error('[route-guard] 获取用户信息失败:', error)
          await userStore.logout()
          next({
            name: 'login',
            query: {
              redirect: to.name,
              ...to.query,
            } as LocationQueryRaw,
          })
        }
      }
    } else {
      console.log('[UserLoginInfoGuard] 未登录')
      if (to.name === 'login') {
        console.log('[UserLoginInfoGuard] 目标是登录页，直接进入')
        next()
        return
      }
      console.log('[UserLoginInfoGuard] 跳转到登录页，附带 redirect 参数:', to.name)
      next({
        name: 'login',
        query: {
          redirect: to.name,
          ...to.query,
        } as LocationQueryRaw,
      })
    }
  })
}
