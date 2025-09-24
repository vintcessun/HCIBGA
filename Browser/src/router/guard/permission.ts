import NProgress from 'nprogress' // progress bar
import type { Router } from 'vue-router'

import usePermission from '@/hooks/permission'
import { useAppStore, useUserStore } from '@/store'
import { NOT_FOUND } from '../constants'
import { appRoutes } from '../routes'

export default function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const appStore = useAppStore()
    const userStore = useUserStore()
    const Permission = usePermission()
    const permissionsAllow = Permission.accessRouter(to)
    // eslint-disable-next-line no-lonely-if
    if (permissionsAllow) next()
    else {
      const destination = Permission.findFirstPermissionRoute(appRoutes, userStore.role) || NOT_FOUND
      next(destination)
    }
    NProgress.done()
  })
}
