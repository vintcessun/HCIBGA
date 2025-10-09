import { DEFAULT_LAYOUT } from '../base'
import { AppRouteRecordRaw } from '../types'

const INFO_MANAGE: AppRouteRecordRaw = {
  path: '/info',
  name: 'info',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.info',
    icon: 'icon-folder',
    requiresAuth: true,
    order: 9,
    roles: ['admin'], // 只有管理员可访问
  },
  children: [
    {
      path: 'import',
      name: 'InfoImport',
      component: () => import('@/views/info/import/index.vue'),
      meta: {
        locale: 'menu.info.import',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
    {
      path: 'export',
      name: 'InfoExport',
      component: () => import('@/views/info/export/index.vue'),
      meta: {
        locale: 'menu.info.export',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
  ],
}

export default INFO_MANAGE
