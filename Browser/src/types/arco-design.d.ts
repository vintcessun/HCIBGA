declare module '@arco-design/web-vue/es/form' {
  import type { FormInstance } from '@arco-design/web-vue'
  export { FormInstance }
}

declare module '@arco-design/web-vue/es/upload' {
  import type { RequestOption } from '@arco-design/web-vue'
  
  interface CustomUploadRequestOption extends RequestOption {
    file: File
    onProgress?: (progress: number) => void
    onSuccess?: (response: any) => void
    onError?: (error: Error) => void
  }
  
  export { CustomUploadRequestOption as UploadRequestOption }
}
