import { DEFAULT_LAYOUT } from '../base'
import { AppRouteRecordRaw } from '../types'

const BONUS: AppRouteRecordRaw = {
  path: '/bonus',
  name: 'bonus',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.bonus',
    icon: 'icon-trophy',
    requiresAuth: true,
    order: 10,
  },
  children: [
    {
      path: 'academic',
      name: 'BonusAcademic',
      component: () => import('@/views/bonus/AcademicList.vue'),
      meta: {
        locale: 'menu.bonus.academic',
        requiresAuth: true,
        roles: ['user'],
      },
    },
    {
      path: 'comprehensive',
      name: 'BonusComprehensive',
      component: () => import('@/views/bonus/ComprehensiveList.vue'),
      meta: {
        locale: 'menu.bonus.comprehensive',
        requiresAuth: true,
        roles: ['user'],
      },
    },
    {
      path: 'summary',
      name: 'BonusSummary',
      component: () => import('@/views/bonus/Summary.vue'),
      meta: {
        locale: 'menu.bonus.summary',
        requiresAuth: true,
        roles: ['user'],
      },
    },
  ],
}

export default BONUS
