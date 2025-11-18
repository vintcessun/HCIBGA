import { LoginData, getUserInfo, login as userLogin, logout as userLogout } from '@/api/user'
import { clearToken } from '@/utils/auth'
import { removeRouteListener } from '@/utils/route-listener'
import { defineStore } from 'pinia'
import useAppStore from '../app'
import { RoleType, UserState } from './types'

const useUserStore = defineStore('user', {
  state: (): UserState => ({
    name: undefined,
    avatar: undefined,
    job: undefined,
    organization: undefined,
    location: undefined,
    email: undefined,
    introduction: undefined,
    personalWebsite: undefined,
    jobName: undefined,
    organizationName: undefined,
    locationName: undefined,
    phone: undefined,
    registrationDate: undefined,
    accountId: undefined,
    certification: undefined,
    role: '',
  }),

  getters: {
    userInfo(state: UserState): UserState {
      return { ...state }
    },
  },

  actions: {
    switchRoles() {
      return new Promise((resolve) => {
        this.role = this.role === 'user' ? 'admin' : 'user'
        resolve(this.role)
      })
    },
    // Set user's information
    setInfo(partial: Partial<UserState>) {
      this.$patch(partial)
    },

    // Reset user's information
    resetInfo() {
      this.$reset()
    },

    // Get user's information
    async info() {
      const { getToken } = await import('@/utils/auth')
      const res = await getUserInfo()

      this.setInfo(res.data)
    },

    async loginToken(token: string, role: RoleType) {
      try {
        // 后端已改为 string 类型，不必再用 String() 转换
        this.accountId = token
        this.role = role || ('user' as any as typeof this.role)
        console.info('[login] 设置角色为:', this.role)

        // 将 accountId 作为 token 存储，确保 isLogin() 返回 true
        const { setToken } = await import('@/utils/auth')
        setToken(token)
        console.info('[login] 已将 accountId 作为 token 存储:', token)
      } catch (err) {
        clearToken()
        throw err
      }
    },
    // Login
    async login(loginForm: LoginData) {
      try {
        const res = await userLogin(loginForm)
        console.info('[login] 完整响应对象:', JSON.stringify(res, null, 2))

        // 后端已改为 string 类型，不必再用 String() 转换
        this.accountId = res.data.id
        this.name = res.data.username
        this.location = res.data.ip
        this.registrationDate = res.data.login_time
        this.role = res.data.role || ('user' as any as typeof this.role)
        console.info('[login] 设置角色为:', this.role)
        if (!res.data.role) {
          console.warn('[login] API未返回角色，已设置默认角色为 user')
        }

        // 将 accountId 作为 token 存储，确保 isLogin() 返回 true
        const { setToken } = await import('@/utils/auth')
        setToken(res.data.id)
        console.info('[login] 已将 accountId 作为 token 存储:', res.data.id)
      } catch (err) {
        clearToken()
        throw err
      }
    },
    logoutCallBack() {
      const appStore = useAppStore()
      this.resetInfo()
      clearToken()
      removeRouteListener()
      appStore.clearServerMenu()
    },
    // Logout
    async logout() {
      try {
        await userLogout()
      } finally {
        this.logoutCallBack()
      }
    },
  },
})

export default useUserStore
