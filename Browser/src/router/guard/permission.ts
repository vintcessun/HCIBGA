import NProgress from 'nprogress' // progress bar
import type { Router } from 'vue-router'

import usePermission from '@/hooks/permission'
import { useAppStore, useUserStore } from '@/store'
import { NOT_FOUND } from '../constants'
import { appRoutes } from '../routes'

export default function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    console.log('[PermissionGuard] 进入路由守卫 from:', from.fullPath, 'to:', to.fullPath)
    const appStore = useAppStore()
    const userStore = useUserStore()
    const Permission = usePermission()
    // eslint-disable-next-line no-console
    console.info('[permission-guard] 当前路由:', to.name, '角色:', userStore.role, 'meta.roles:', to.meta?.roles)
    console.log('[PermissionGuard] 调用 Permission.accessRouter 检查权限')
    const permissionsAllow = Permission.accessRouter(to)
    console.log('[PermissionGuard] accessRouter 返回:', permissionsAllow)
    // eslint-disable-next-line no-console
    console.info('[permission-guard] accessRouter 返回:', permissionsAllow)
    // eslint-disable-next-line no-lonely-if
    if (permissionsAllow) {
      console.log('[PermissionGuard] 权限允许，调用 next() 进入路由')
      next()
    } else {
      console.warn('[PermissionGuard] 权限不允许，查找第一个可访问路由')
      const destination = Permission.findFirstPermissionRoute(appRoutes, userStore.role) || NOT_FOUND
      console.log('[PermissionGuard] 跳转到:', destination)
      next(destination)
    }
    NProgress.done()
  })
}
