<template>
  <div class="login-bg">
    <div class="container">
      <div class="scan-login-btn" @click="handleScanLogin">
        <svg class="scan-icon" viewBox="0 0 24 24" width="20" height="20">
          <rect x="3" y="3" width="7" height="7" rx="2" fill="none" stroke="#00308f" stroke-width="2" />
          <rect x="14" y="3" width="7" height="7" rx="2" fill="none" stroke="#00308f" stroke-width="2" />
          <rect x="14" y="14" width="7" height="7" rx="2" fill="none" stroke="#00308f" stroke-width="2" />
          <rect x="3" y="14" width="7" height="7" rx="2" fill="none" stroke="#00308f" stroke-width="2" />
        </svg>
        <span>{{ $t('login.scan.title') }}</span>
      </div>
      <a-modal v-model:visible="showQr" :title="$t('login.scan.title')" :footer="false" width="320px">
        <div class="qr-modal-content">
          <img :src="qrUrl" alt="二维码" style="width: 200px; height: 200px; display: block; margin: 0 auto" />
          <div style="text-align: center; margin-top: 12px; color: #888">{{ $t('login.scan.tip') }}</div>
        </div>
      </a-modal>
      <div class="logo">
        <img alt="logo" src="//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/dfdba5317c0c20ce20e64fac803d52bc.svg~tplv-49unhts6dw-image.image" />
        <div class="logo-text">{{ $t('app.title') }}</div>
      </div>
      <LoginBanner />
      <div class="content">
        <div class="content-inner">
          <LoginForm />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { getQrAuthResult, getQrCode, pollQrStatus } from '@/api/user'
import useUserStore from '@/store/modules/user'
import { defineComponent, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import LoginBanner from './components/banner.vue'
import LoginForm from './components/login-form.vue'

export default defineComponent({
  components: {
    LoginBanner,
    LoginForm,
  },
  setup() {
    const { t } = useI18n()
    const showQr = ref(false)
    const qrUrl = ref('')
    const qrId = ref('')
    let pollingTimer: ReturnType<typeof setInterval> | null = null
    const userStore = useUserStore()
    const router = useRouter()

    const startPolling = () => {
      if (pollingTimer) {
        clearInterval(pollingTimer)
      }
      pollingTimer = setInterval(async () => {
        try {
          const statusRes = await pollQrStatus(qrId.value)
          if (statusRes.data.status === 'done') {
            clearInterval(pollingTimer as ReturnType<typeof setInterval>)
            const authRes = await getQrAuthResult(qrId.value)
            console.log('[ScanLogin] authRes', authRes)
            try {
              await userStore.loginToken(authRes.data.token, authRes.data.role)
              console.log('[ScanLogin] userStore.login 完成，开始跳转')
              console.log('[ScanLogin] 获取 router 实例')
              console.log('[ScanLogin] 当前路由 query:', router?.currentRoute?.value?.query)
              const { redirect, ...othersQuery } = router?.currentRoute?.value?.query || {}
              console.log('[ScanLogin] 准备跳转到:', redirect || 'MaterialList', 'query:', othersQuery)
              router
                .push({
                  name: (redirect as string) || 'MaterialList',
                  query: {
                    ...othersQuery,
                  },
                })
                .then(() => {
                  console.log('[ScanLogin] 跳转完成')
                })
                .catch((err: any) => {
                  console.error('[ScanLogin] 跳转出错', err)
                })
            } catch (err) {
              console.error('[ScanLogin] 登录错误', err)
            }
            showQr.value = false
          } else if (statusRes.data.status === 'expired') {
            clearInterval(pollingTimer as ReturnType<typeof setInterval>)
            console.warn(t('login.scan.expired'))
            qrUrl.value = ''
            qrId.value = ''
            showQr.value = false
          }
        } catch (err) {
          console.error('轮询二维码状态失败', err)
        }
      }, 3000)
    }

    const handleScanLogin = async () => {
      try {
        const res = await getQrCode()
        qrUrl.value = res.data.qrUrl
        qrId.value = res.data.qrId
        showQr.value = true
        startPolling()
      } catch (err) {
        console.error('获取二维码失败', err)
      }
    }

    return { showQr, qrUrl, handleScanLogin }
  },
})
</script>

<style lang="less" scoped>
.login-bg {
  min-height: 100vh;
  width: 100vw;
  background: linear-gradient(135deg, #e0e7ff 0%, #f8fafc 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}
.container {
  display: flex;
  min-height: 500px;
  width: 900px;
  max-width: 96vw;
  border-radius: 18px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.15);
  background: #fff;
  position: relative;
  overflow: hidden;
  margin: 40px 0;
  @media (max-width: 900px) {
    flex-direction: column;
    width: 98vw;
    min-height: unset;
  }
  .banner {
    width: 400px;
    background: linear-gradient(163.85deg, #1d2129 0%, #00308f 100%);
    @media (max-width: 900px) {
      width: 100%;
      min-height: 220px;
    }
  }
  .content {
    border-radius: 18px;
    position: relative;
    display: flex;
    flex: 1;
    align-items: center;
    justify-content: center;
    padding: 40px 0;
    background: #fff;
    @media (max-width: 900px) {
      padding: 24px 0;
    }
  }
}
.logo {
  position: absolute;
  top: 32px;
  left: 32px;
  z-index: 2;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.85);
  padding: 6px 18px 6px 10px;
  border-radius: 24px;
  box-shadow: 0 2px 8px 0 rgba(31, 38, 135, 0.08);
  &-text {
    margin-left: 8px;
    color: #00308f;
    font-size: 22px;
    font-weight: bold;
    letter-spacing: 1px;
  }
  img {
    height: 32px;
    width: 32px;
  }
  @media (max-width: 900px) {
    position: static;
    margin: 24px auto 0 auto;
    justify-content: center;
    box-shadow: none;
    background: transparent;
  }
}
.scan-login-btn {
  position: absolute;
  top: 24px;
  right: 32px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 6px;
  background: #f4f8ff;
  color: #00308f;
  border-radius: 18px;
  padding: 6px 16px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  box-shadow: 0 2px 8px 0 rgba(31, 38, 135, 0.08);
  transition: background 0.2s;
  &:hover {
    background: #e0e7ff;
  }
  .scan-icon {
    display: inline-block;
    vertical-align: middle;
  }
}
.qr-modal-content {
  padding: 12px 0 0 0;
}
</style>
