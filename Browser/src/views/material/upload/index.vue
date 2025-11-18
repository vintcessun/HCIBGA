<template>
  <div class="material-upload">
    <a-card :title="$t('menu.material.upload')" class="upload-card">
      <!-- 文件上传区域 -->
      <div class="upload-section">
        <h3 class="section-title">{{ $t('material.upload.files') }}</h3>
        <a-upload
          :multiple="true"
          :file-list="fileList"
          :before-upload="beforeUpload"
          @change="handleFileChange"
          @remove="handleFileRemove"
          :show-upload-list="true"
          :disabled="uploading"
          :custom-request="handleFileUpload"
        >
          <template #upload-button>
            <div class="upload-area">
              <div class="upload-icon">
                <icon-upload />
              </div>
              <div class="upload-text">
                {{ $t('material.upload.clickToUpload') }}
              </div>
              <div class="upload-hint">
                支持格式: PDF, DOC, DOCX,
                <br />
                PPT, PPTX, XLS, XLSX, PNG, JPG等
              </div>
            </div>
          </template>
        </a-upload>

        <!-- 已上传文件列表 -->
        <div v-if="uploadedFiles.length > 0" class="uploaded-files">
          <h4 class="files-title">{{ $t('material.upload.uploadedFiles') }}</h4>
          <div v-for="(file, index) in uploadedFiles" :key="file.file_id" class="file-item">
            <icon-file />
            <a class="file-name" :href="fileList.find((f) => f.name === file.name)?.url" target="_blank">
              {{ file.name }}
            </a>
            <span class="file-status" :class="file.status">
              {{ file.status === 'uploaded' ? '已上传' : '已存在' }}
            </span>
            <!-- 上移按钮 -->
            <a-button size="mini" type="text" @click="moveFileUp(index)" :disabled="index === 0">↑</a-button>
            <!-- 下移按钮 -->
            <a-button size="mini" type="text" @click="moveFileDown(index)" :disabled="index === uploadedFiles.length - 1">↓</a-button>
            <!-- 删除按钮 -->
            <a-button size="mini" type="text" status="danger" @click="deleteUploadedFile(index)">删除</a-button>
          </div>
        </div>
      </div>

      <!-- 材料信息表单 -->
      <div class="form-section">
        <a-button type="primary" style="margin-bottom: 12px" @click="handleStartRecognition">开始识别</a-button>
        <h3 class="section-title">{{ $t('material.upload.materialInfo') }}</h3>
        <a-form
          :model="form"
          :label-col-props="{ span: 6 }"
          :wrapper-col-props="{ span: 18 }"
          class="upload-form"
          :rules="formRules"
          ref="formRef"
        >
          <a-form-item :label="$t('material.upload.title')" field="title">
            <a-input v-model="form.title" :placeholder="$t('material.upload.titlePlaceholder')" :max-length="100" />
          </a-form-item>

          <a-form-item :label="$t('material.upload.description')" field="description">
            <a-textarea
              v-model="form.description"
              :placeholder="$t('material.upload.descriptionPlaceholder')"
              :auto-size="{
                minRows: 3,
                maxRows: 5,
              }"
              :max-length="500"
            />
          </a-form-item>

          <a-form-item :label="$t('material.upload.category')" field="category">
            <a-select v-model="form.category" :placeholder="$t('material.upload.categoryPlaceholder')" allow-clear>
              <a-option value="学术专长成绩-科研成果">学术专长成绩-科研成果</a-option>
              <a-option value="学术专长成绩-学业竞赛">学术专长成绩-学业竞赛</a-option>
              <a-option value="学术专长成绩-创新创业训练">学术专长成绩-创新创业训练</a-option>
              <a-option value="综合表现加分-国际组织实习">综合表现加分-国际组织实习</a-option>
              <a-option value="综合表现加分-参军入伍服兵役">综合表现加分-参军入伍服兵役</a-option>
              <a-option value="综合表现加分-志愿服务">综合表现加分-志愿服务</a-option>
              <a-option value="综合表现加分-荣誉称号">综合表现加分-荣誉称号</a-option>
              <a-option value="综合表现加分-社会工作">综合表现加分-社会工作</a-option>
              <a-option value="综合表现加分-体育比赛">综合表现加分-体育比赛</a-option>
            </a-select>
          </a-form-item>

          <a-form-item :label="$t('material.upload.tags')" field="tags">
            <a-select v-model="form.tags" multiple :placeholder="$t('material.upload.tagsPlaceholder')" allow-clear>
              <a-option value="期刊论文发表（A 类）">期刊论文发表（A 类）</a-option>
              <a-option value="会议论文收录（B 类）">会议论文收录（B 类）</a-option>
              <a-option value="会议论文收录（C 类）">会议论文收录（C 类）</a-option>
              <a-option value="Nature/Science/Cell 主刊论文">Nature/Science/Cell 主刊论文</a-option>
              <a-option value="Cell 子刊论文（IF≥10）">Cell 子刊论文（IF≥10）</a-option>
              <a-option value="国家发明专利授权（第一作者）">国家发明专利授权（第一作者）</a-option>
              <a-option value="国家发明专利授权（独立作者）">国家发明专利授权（独立作者）</a-option>
              <a-option value="高水平中文学术期刊论文">高水平中文学术期刊论文</a-option>
              <a-option value="信息与通信工程国际期刊论文">信息与通信工程国际期刊论文</a-option>
              <a-option value="国家级 A + 类竞赛一等奖及以上">国家级 A + 类竞赛一等奖及以上</a-option>
              <a-option value="国家级 A 类竞赛二等奖">国家级 A 类竞赛二等奖</a-option>
              <a-option value="省级 A 类竞赛一等奖及以上">省级 A 类竞赛一等奖及以上</a-option>
              <a-option value="省级 A - 类竞赛二等奖">省级 A - 类竞赛二等奖</a-option>
              <a-option value="ICPC 亚洲区域赛获奖">ICPC 亚洲区域赛获奖</a-option>
              <a-option value="CCPC 竞赛获奖">CCPC 竞赛获奖</a-option>
              <a-option value="ICPC 全球总决赛获奖">ICPC 全球总决赛获奖</a-option>
              <a-option value="CCF CSP 认证前 0.2%（等同国一）">CCF CSP 认证前 0.2%（等同国一）</a-option>
              <a-option value="CCF CSP 认证前 1.5%（等同国二）">CCF CSP 认证前 1.5%（等同国二）</a-option>
              <a-option value="中国国际大学生创新大赛团体获奖">中国国际大学生创新大赛团体获奖</a-option>
              <a-option value="挑战杯学术科技作品竞赛个人获奖">挑战杯学术科技作品竞赛个人获奖</a-option>
              <a-option value="挑战杯创业计划大赛团队获奖">挑战杯创业计划大赛团队获奖</a-option>
              <a-option value="国家级创新实验计划项目（组长）">国家级创新实验计划项目（组长）</a-option>
              <a-option value="国家级创新实验计划项目（成员）">国家级创新实验计划项目（成员）</a-option>
              <a-option value="省级创新实验计划项目（组长）">省级创新实验计划项目（组长）</a-option>
              <a-option value="省级创新实验计划项目（成员）">省级创新实验计划项目（成员）</a-option>
              <a-option value="校级创新实验计划项目（组长）">校级创新实验计划项目（组长）</a-option>
              <a-option value="校级创新实验计划项目（成员）">校级创新实验计划项目（成员）</a-option>
              <a-option value="创新创业训练项目结项（教务处证明）">创新创业训练项目结项（教务处证明）</a-option>
              <a-option value="创新创业训练项目结项（创新网截图）">创新创业训练项目结项（创新网截图）</a-option>
              <a-option value="国际组织一学年实习">国际组织一学年实习</a-option>
              <a-option value="国际组织半年（超一学期）实习">国际组织半年（超一学期）实习</a-option>
              <a-option value="国际组织实习证明（满学年）">国际组织实习证明（满学年）</a-option>
              <a-option value="国际组织实习证明（半年）">国际组织实习证明（半年）</a-option>
              <a-option value="参军入伍满 1 年（含）服兵役">参军入伍满 1 年（含）服兵役</a-option>
              <a-option value="参军入伍满 2 年（含）服兵役">参军入伍满 2 年（含）服兵役</a-option>
              <a-option value="服兵役证明（1-2 年）">服兵役证明（1-2 年）</a-option>
              <a-option value="服兵役证明（2 年以上）">服兵役证明（2 年以上）</a-option>
              <a-option value="志愿服务时长">志愿服务时长</a-option>
              <a-option value="国家级志愿服务个人表彰">国家级志愿服务个人表彰</a-option>
              <a-option value="省级志愿服务团队表彰（队长）">省级志愿服务团队表彰（队长）</a-option>
              <a-option value="省级志愿服务团队表彰（队员）">省级志愿服务团队表彰（队员）</a-option>
              <a-option value="校级志愿服务个人表彰">校级志愿服务个人表彰</a-option>
              <a-option value="抢险救灾志愿服务突出表现">抢险救灾志愿服务突出表现</a-option>
              <a-option value="国家级优秀三好学生">国家级优秀三好学生</a-option>
              <a-option value="国家级优秀共产党员">国家级优秀共产党员</a-option>
              <a-option value="国家级‘自强之星’">国家级‘自强之星’</a-option>
              <a-option value="省级优秀学生干部">省级优秀学生干部</a-option>
              <a-option value="省级优秀团员">省级优秀团员</a-option>
              <a-option value="省级社会实践优秀个人">省级社会实践优秀个人</a-option>
              <a-option value="校级三好学生">校级三好学生</a-option>
              <a-option value="校级优秀学生干部">校级优秀学生干部</a-option>
              <a-option value="国家级五四红旗团支部（集体）">国家级五四红旗团支部（集体）</a-option>
              <a-option value="省级优秀班集体（集体）">省级优秀班集体（集体）</a-option>
              <a-option value="校级优秀团支部书记">校级优秀团支部书记</a-option>
              <a-option value="院学生会执行主席（任职 1 学年）">院学生会执行主席（任职 1 学年）</a-option>
              <a-option value="团总支书记（任职 1 学年）">团总支书记（任职 1 学年）</a-option>
              <a-option value="校学生会主席团成员（任职 1 学年）">校学生会主席团成员（任职 1 学年）</a-option>
              <a-option value="班长（任职 1 学年）">班长（任职 1 学年）</a-option>
              <a-option value="团支部书记（任职半学年）">团支部书记（任职半学年）</a-option>
              <a-option value="党支部书记（任职 1 学年）">党支部书记（任职 1 学年）</a-option>
              <a-option value="社团社长（任职 1 学年）">社团社长（任职 1 学年）</a-option>
              <a-option value="院学生会部长（辅导员打分 90+）">院学生会部长（辅导员打分 90+）</a-option>
              <a-option value="班委（任职 1 学年）">班委（任职 1 学年）</a-option>
              <a-option value="系团总支副书记（任职 1 学年）">系团总支副书记（任职 1 学年）</a-option>
              <a-option value="国际级体育团体冠军">国际级体育团体冠军</a-option>
              <a-option value="国际级体育团体亚军">国际级体育团体亚军</a-option>
              <a-option value="国家级体育团体冠军">国家级体育团体冠军</a-option>
              <a-option value="国家级体育团体季军">国家级体育团体季军</a-option>
              <a-option value="国际级体育比赛第四至八名">国际级体育比赛第四至八名</a-option>
              <a-option value="国家级体育个人亚军">国家级体育个人亚军</a-option>
              <a-option value="国家级体育个人季军">国家级体育个人季军</a-option>
              <a-option value="厦门大学代表国际体育比赛获奖">厦门大学代表国际体育比赛获奖</a-option>
              <a-option value="厦门大学代表国家级体育比赛获奖">厦门大学代表国家级体育比赛获奖</a-option>
              <a-option value="体育团体项目（5 人以内）获奖">体育团体项目（5 人以内）获奖</a-option>
            </a-select>
          </a-form-item>
        </a-form>
      </div>

      <!-- 操作按钮 -->
      <div class="actions-section">
        <a-form-item :wrapper-col-props="{ span: 24 }">
          <div class="upload-actions">
            <a-button type="primary" :loading="uploading" :disabled="uploadedFiles.length === 0 || !formValid" @click="handleSubmit">
              {{ $t('material.upload.submit') }}
            </a-button>
            <a-button @click="handleReset">
              {{ $t('material.upload.reset') }}
            </a-button>
          </div>
        </a-form-item>
      </div>
    </a-card>

    <!-- 成功模态框 -->
    <a-modal v-model:visible="showSuccessModal" :title="$t('material.upload.successTitle')" :footer="false">
      <div class="success-content">
        <icon-check-circle style="color: #00b42a; font-size: 48px; margin-bottom: 16px" />
        <p>{{ $t('material.upload.successMessage') }}</p>
        <a-button type="primary" @click="showSuccessModal = false">
          {{ $t('material.upload.confirm') }}
        </a-button>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { llmFormRecognition } from '@/api/form'
import { uploadMaterial } from '@/api/material'
import { type UploadFileResponse, checkFileExists, uploadFile } from '@/api/upload'
import { Message } from '@arco-design/web-vue'
import type { UploadRequestOption } from '@arco-design/web-vue/es/upload'
import SparkMD5 from 'spark-md5'
import { computed, reactive, ref } from 'vue'

interface FileItem {
  uid: string
  name: string
  url?: string
  status?: 'uploading' | 'done' | 'error'
  percent?: number
}

interface UploadedFile {
  file_id: string
  name: string
  size: number
  status: 'uploaded' | 'exists'
}

const fileList = ref<FileItem[]>([])
const uploadedFiles = ref<UploadedFile[]>([])
const uploading = ref(false)
const showSuccessModal = ref(false)
const formRef = ref()

const form = reactive({
  title: '',
  description: '',
  category: '',
  tags: [] as string[],
})

const formRules = {
  title: [{ required: true, message: '请输入材料标题' }],
  category: [{ required: true, message: '请选择材料类别' }],
}

const formValid = computed(() => {
  return form.title.trim() && form.category
})

// 计算文件MD5
const calculateFileMD5 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const spark = new SparkMD5.ArrayBuffer()
    const reader = new FileReader()

    reader.onload = (e) => {
      spark.append(e.target?.result as ArrayBuffer)
      resolve(spark.end())
    }

    reader.onerror = reject
    reader.readAsArrayBuffer(file)
  })
}

const beforeUpload = (file: File) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    Message.error('文件大小不能超过10MB')
    return false
  }
  return true
}

const handleFileChange = (files: FileItem[]) => {
  fileList.value = files
}

const handleFileRemove = (file: FileItem) => {
  const index = fileList.value.findIndex((item) => item.uid === file.uid)
  if (index !== -1) {
    fileList.value.splice(index, 1)
  }

  // 同时从已上传文件列表中移除
  const uploadedIndex = uploadedFiles.value.findIndex((item) => item.name === file.name)
  if (uploadedIndex !== -1) {
    uploadedFiles.value.splice(uploadedIndex, 1)
  }
}

const handleFileUpload = async (options: UploadRequestOption & { fileItem?: any }) => {
  console.log('[handleFileUpload] start', options)
  const { file, fileItem, onProgress, onSuccess, onError } = options
  let fileObj: File

  // 根据 Arco 官方文档，优先取 option.file
  if (file) {
    fileObj = file as File
    console.log('[handleFileUpload] using option.file', fileObj.name)
  } else if (fileItem?.file) {
    fileObj = fileItem.file as File
    console.log('[handleFileUpload] using fileItem.file', fileObj.name)
  } else if ((fileItem as any)?.originFile) {
    fileObj = (fileItem as any).originFile as File
    console.log('[handleFileUpload] using fileItem.originFile', fileObj.name)
  } else {
    console.error('[handleFileUpload] No valid file found', options)
    onError?.(new Error('No valid file provided for upload'))
    return
  }

  try {
    console.log('[handleFileUpload] calculating MD5 for', fileObj.name)
    // 计算文件MD5
    const md5 = await calculateFileMD5(fileObj)
    console.log('[handleFileUpload] MD5 calculated:', md5)

    // 先检查文件是否存在
    console.log('[handleFileUpload] checking if file exists')
    const checkResponse = await checkFileExists({ md5, filename: fileObj.name })
    console.log('[handleFileUpload] checkFileExists response:', checkResponse)
    const { exists: fileExists, file_id: existingFileId, url: existingFileUrl } = checkResponse.data

    if (fileExists) {
      // 文件已存在，直接返回成功
      const response: UploadFileResponse = {
        file_id: existingFileId!,
        url: existingFileUrl!,
        md5,
      }

      // 添加到已上传文件列表
      uploadedFiles.value.push({
        file_id: response.file_id,
        name: fileObj.name,
        size: fileObj.size,
        status: 'exists',
      })

      // 同步更新 fileList 中对应文件状态为 done，以显示 √
      const target = fileList.value.find((f) => f.name === fileObj.name)
      if (target) {
        target.status = 'done'
      }

      console.log('[handleFileUpload] file exists, skipping upload')
      return
    }

    // 文件不存在，调用真正的上传API
    console.log('[handleFileUpload] file not exists, starting upload')
    const uploadResponse = await uploadFile(fileObj, (progress: number) => {
      console.log('[handleFileUpload] upload progress:', progress)
      // 上传完成后将 fileList 中对应文件状态改为 done，以显示 √
      if (progress === 100) {
        const target = fileList.value.find((f) => f.name === fileObj.name)
        if (target) {
          target.status = 'done'
        }
      }
    })

    // 添加到已上传文件列表
    uploadedFiles.value.push({
      file_id: uploadResponse.data.file_id,
      name: fileObj.name,
      size: fileObj.size,
      status: 'uploaded', // 状态保持类型兼容，并可在界面通过已上传样式显示√
    })

    console.log('[handleFileUpload] upload finished:', uploadResponse.data)
  } catch (error) {
    console.error('[handleFileUpload] error:', error)
  }
}

const moveFileUp = (index: number) => {
  if (index > 0) {
    const temp = uploadedFiles.value[index]
    uploadedFiles.value.splice(index, 1)
    uploadedFiles.value.splice(index - 1, 0, temp)
  }
}

const moveFileDown = (index: number) => {
  if (index < uploadedFiles.value.length - 1) {
    const temp = uploadedFiles.value[index]
    uploadedFiles.value.splice(index, 1)
    uploadedFiles.value.splice(index + 1, 0, temp)
  }
}

const deleteUploadedFile = (index: number) => {
  uploadedFiles.value.splice(index, 1)
}

const handleReset = () => {
  fileList.value = []
  uploadedFiles.value = []
  form.title = ''
  form.description = ''
  form.category = ''
  form.tags = []
  formRef.value?.resetFields()
}

const handleStartRecognition = async () => {
  if (uploadedFiles.value.length === 0) {
    Message.error('请先上传文件')
    return
  }
  try {
    uploading.value = true
    const fileIds = uploadedFiles.value.map((f) => f.file_id)
    const res = await llmFormRecognition({ files: fileIds })
    if (res && typeof res.data === 'object') {
      const { title, category, tags, description } = res.data
      form.title = title || ''
      form.category = category || ''
      form.tags = Array.isArray(tags) ? tags : []
      form.description = description || ''
      Message.success('识别完成，已填充表单')
    } else {
      Message.error('识别结果格式不正确')
    }
  } catch (err) {
    console.error(err)
    Message.error('识别失败，请重试')
  } finally {
    uploading.value = false
  }
}

const handleSubmit = async () => {
  if (uploadedFiles.value.length === 0) {
    Message.error('请先上传文件')
    return
  }

  if (!formValid.value) {
    Message.error('请填写完整的材料信息')
    return
  }

  uploading.value = true
  try {
    // 调用材料上传API
    await uploadMaterial({
      title: form.title,
      description: form.description,
      category: form.category,
      tags: form.tags,
      files: uploadedFiles.value.map((file) => file.file_id),
    })

    Message.success('材料上传成功！')
    showSuccessModal.value = true
    handleReset()
  } catch (error) {
    Message.error('上传失败，请重试')
  } finally {
    uploading.value = false
  }
}
</script>

<style lang="less" scoped>
.material-upload {
  padding: 20px;

  .upload-card {
    min-width: 800px;
    margin: 0 auto;
  }

  .section-title {
    margin-bottom: 16px;
    color: #1d2129;
    font-size: 16px;
    font-weight: 600;
  }

  .upload-section {
    text-align: center;
    margin-bottom: 24px;
  }

  .upload-area {
    padding: 40px 0;
    text-align: center;
    border: 2px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    transition: border-color 0.3s;
    margin: 0 auto;
    max-width: 400px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    &:hover {
      border-color: #165dff;
    }

    .upload-icon {
      font-size: 48px;
      color: #165dff;
      margin-bottom: 16px;
    }

    .upload-text {
      font-size: 16px;
      color: #1d2129;
      margin-bottom: 8px;
    }

    .upload-hint {
      width: 500px;
      font-size: 14px;
      color: #86909c;
    }
  }

  .uploaded-files {
    margin-top: 16px;
    padding: 16px;
    border: 1px solid #e5e6eb;
    border-radius: 6px;
    background: #f7f8fa;

    .files-title {
      margin-bottom: 12px;
      color: #1d2129;
      font-size: 14px;
      font-weight: 600;
    }

    .file-item {
      display: flex;
      align-items: center;
      padding: 8px;
      margin-bottom: 8px;
      background: white;
      border-radius: 4px;
      border: 1px solid #e5e6eb;

      &:last-child {
        margin-bottom: 0;
      }

      .file-name {
        flex: 1;
        margin-left: 8px;
        color: #1d2129;
      }

      .file-status {
        padding: 2px 8px;
        border-radius: 12px;
        font-size: 12px;

        &.uploaded {
          background: #e8ffea;
          color: #00b42a;
        }

        &.exists {
          background: #e8f4ff;
          color: #165dff;
        }
      }
    }
  }

  .form-section {
    text-align: center;
    margin-bottom: 24px;
  }

  .upload-form {
    max-width: 1000px;
    margin-top: 16px;
  }

  .actions-section {
    margin-top: 24px;
    display: flex;
    justify-content: center;
  }

  .upload-actions {
    text-align: center;
    min-width: 1200px;
    display: flex;
    justify-content: center;
    gap: 12px;

    .arco-btn {
      min-width: 100px;
    }
  }

  .success-content {
    text-align: center;
    padding: 20px;

    p {
      margin: 16px 0 24px;
      color: #1d2129;
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .material-upload {
    padding: 12px;

    .upload-card {
      margin: 0;
    }

    .upload-area {
      padding: 24px 0;
    }

    .upload-form {
      :deep(.arco-form-item) {
        margin-bottom: 16px;
      }
    }

    .upload-actions {
      .arco-btn {
        width: 100%;
        margin: 8px 0;
      }
    }
  }
}
</style>
