import type { RouteRecordNormalized } from 'vue-router'

import INFO_MANAGE from './modules/info'
import SUBMIT_MATERIAL from './modules/submit-material'
import USER from './modules/user'

const modules = import.meta.glob('./modules/*.ts', { eager: true })
const externalModules = import.meta.glob('./externalModules/*.ts', {
  eager: true,
})

function formatModules(_modules: any, result: RouteRecordNormalized[]) {
  Object.keys(_modules).forEach((key) => {
    const defaultModule = _modules[key].default
    if (!defaultModule) return
    const moduleList = Array.isArray(defaultModule) ? [...defaultModule] : [defaultModule]
    result.push(...moduleList)
  })
  return result
}

export const appRoutes: RouteRecordNormalized[] = [
  ...new Set([...formatModules(modules, []), USER as any, SUBMIT_MATERIAL as any, INFO_MANAGE as any]),
]

export const appExternalRoutes: RouteRecordNormalized[] = formatModules(externalModules, [])
