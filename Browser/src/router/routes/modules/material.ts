import { DEFAULT_LAYOUT } from '../base'
import { AppRouteRecordRaw } from '../types'

const MATERIAL: AppRouteRecordRaw = {
  path: '/material',
  name: 'material',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.material',
    icon: 'icon-file',
    requiresAuth: true,
    order: 8,
  },
  children: [
    {
      path: 'import',
      name: 'MaterialImport',
      component: DEFAULT_LAYOUT,
      meta: {
        locale: 'menu.user.import',
        requiresAuth: true,
        roles: ['*'],
      },
      children: [
        {
          path: 'qicai',
          name: 'MaterialImportQicai',
          component: () => import('@/views/user/import/qicai.vue'),
          meta: {
            locale: 'menu.user.import.qicai',
            requiresAuth: true,
            roles: ['*'],
          },
        },
        {
          path: 'jiaowu',
          name: 'MaterialImportJiaowu',
          component: () => import('@/views/user/import/jiaowu.vue'),
          meta: {
            locale: 'menu.user.import.jiaowu',
            requiresAuth: true,
            roles: ['*'],
          },
        },
      ],
    },
    {
      path: 'upload',
      name: 'MaterialUpload',
      component: () => import('@/views/material/upload/index.vue'),
      meta: {
        locale: 'menu.material.upload',
        requiresAuth: true,
        roles: ['user'],
      },
    },
    {
      path: 'list',
      name: 'MaterialList',
      component: () => import('@/views/material/list/index.vue'),
      meta: {
        locale: 'menu.material.list',
        requiresAuth: true,
        roles: ['admin', 'reviewer', 'user'],
      },
    },
    {
      path: 'review',
      name: 'MaterialReview',
      component: () => import('@/views/material/review/index.vue'),
      meta: {
        locale: 'menu.material.review',
        requiresAuth: true,
        roles: ['reviewer', 'admin'],
      },
    },
    {
      path: 'statistics',
      name: 'MaterialStatistics',
      component: () => import('@/views/material/statistics/index.vue'),
      meta: {
        locale: 'menu.material.statistics',
        requiresAuth: true,
        roles: ['admin', 'user'],
      },
    },
  ],
}

export default MATERIAL
