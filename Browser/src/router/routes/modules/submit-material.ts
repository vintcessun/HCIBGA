import { DEFAULT_LAYOUT } from '../base'
import { AppRouteRecordRaw } from '../types'

const SUBMIT_MATERIAL: AppRouteRecordRaw = {
  path: '/submit-material',
  name: 'submit-material',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.material.submit',
    icon: 'icon-upload',
    requiresAuth: true,
    order: 9,
  },
  children: [
    {
      path: 'submit-info',
      name: 'SubmitInfoMaterial',
      component: () => import('@/views/user/batch/submit-info.vue'),
      meta: {
        locale: 'menu.user.submit.list',
        requiresAuth: true,
        roles: ['user'],
      },
    },
    {
      path: 'batch-list',
      name: 'BatchListMaterial',
      component: () => import('@/views/user/batch-list/manage-batch.vue'),
      meta: {
        locale: 'menu.user.batch.list',
        requiresAuth: true,
        roles: ['admin', 'reviewer', 'user'],
      },
    },
    {
      path: 'manage-batch',
      name: 'ManageBatchMaterial',
      component: () => import('@/views/material/submit/review.vue'),
      meta: {
        locale: 'menu.material.submit.manage',
        requiresAuth: true,
        roles: ['admin', 'reviewer'],
      },
    },
  ],
}

export default SUBMIT_MATERIAL
