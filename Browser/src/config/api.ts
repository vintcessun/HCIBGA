// API基础配置 - 仅包含login相关配置
export const API_CONFIG = {
  // 基础URL - 从环境变量获取
  BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8088',

  // API路径配置 - 仅包含login相关
  PATHS: {
    // 用户登录相关API
    USER: {
      LOGIN: '/api/user/auth',
    },
  },
} as const

// 获取完整的API URL
export const getApiUrl = (path: string): string => {
  return `${API_CONFIG.BASE_URL}${path}`
}
