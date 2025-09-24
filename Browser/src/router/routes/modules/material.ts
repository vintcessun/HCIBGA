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
      path: 'upload',
      name: 'MaterialUpload',
      component: () => import('@/views/material/upload/index.vue'),
      meta: {
        locale: 'menu.material.upload',
        requiresAuth: true,
        roles: ['user'], // 只有普通用户可以上传
      },
    },
    {
      path: 'list',
      name: 'MaterialList',
      component: () => import('@/views/material/list/index.vue'),
      meta: {
        locale: 'menu.material.list',
        requiresAuth: true,
        roles: ['admin', 'reviewer'],
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
        roles: ['admin'],
      },
    },
  ],
}

export default MATERIAL
