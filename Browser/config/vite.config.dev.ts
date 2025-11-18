import { mergeConfig } from 'vite'
import eslint from 'vite-plugin-eslint'
import baseConfig from './vite.config.base'

export default mergeConfig(
  {
    mode: 'development',
    server: {
      open: true,
      host: '0.0.0.0',
      allowedHosts: ['vintces.icu'],
      fs: {
        strict: true,
      },
      proxy: {
        '/api': {
          target: process.env.VITE_API_BASE_URL || 'http://127.0.0.1:8088',
          changeOrigin: true,
          rewrite: (path: string) => path.replace(/^\/api/, '/api'),
        },
      },
    },
    plugins: [
      eslint({
        cache: false,
        include: ['src/**/*.ts', 'src/**/*.tsx', 'src/**/*.vue'],
        exclude: ['node_modules'],
      }),
    ],
    define: {
      'process.env.mock': false,
    },
  },
  baseConfig
)
