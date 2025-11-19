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
                {{ $t('material.upload.supportFormat') }}
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
              {{ file.status === 'uploaded' ? $t('material.upload.status.uploaded') : $t('material.upload.status.exists') }}
            </span>
            <!-- 上移按钮 -->
            <a-button size="mini" type="text" @click="moveFileUp(index)" :disabled="index === 0">↑</a-button>
            <!-- 下移按钮 -->
            <a-button size="mini" type="text" @click="moveFileDown(index)" :disabled="index === uploadedFiles.length - 1">↓</a-button>
            <!-- 删除按钮 -->
            <a-button size="mini" type="text" status="danger" @click="deleteUploadedFile(index)">
              {{ $t('material.upload.delete') }}
            </a-button>
          </div>
        </div>
      </div>

      <!-- 材料信息表单 -->
      <div class="form-section">
        <a-button type="primary" style="margin-bottom: 12px" @click="handleStartRecognition">
          {{ $t('material.upload.startRecognition') }}
        </a-button>
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

          <a-form-item :label="$t('material.upload.category')" field="category">
            <a-select v-model="form.category" :placeholder="$t('material.upload.categoryPlaceholder')" allow-clear>
              <a-option value="学术专长成绩-科研成果">{{ $t('material.category.academic.research') }}</a-option>
              <a-option value="学术专长成绩-学业竞赛">{{ $t('material.category.academic.competition') }}</a-option>
              <a-option value="学术专长成绩-创新创业训练">{{ $t('material.category.academic.innovation') }}</a-option>
              <a-option value="综合表现加分-国际组织实习">{{ $t('material.category.comprehensive.internship') }}</a-option>
              <a-option value="综合表现加分-参军入伍服兵役">{{ $t('material.category.comprehensive.military') }}</a-option>
              <a-option value="综合表现加分-志愿服务">{{ $t('material.category.comprehensive.volunteer') }}</a-option>
              <a-option value="综合表现加分-荣誉称号">{{ $t('material.category.comprehensive.honor') }}</a-option>
              <a-option value="综合表现加分-社会工作">{{ $t('material.category.comprehensive.social') }}</a-option>
              <a-option value="综合表现加分-体育比赛">{{ $t('material.category.comprehensive.sports') }}</a-option>
            </a-select>
          </a-form-item>

          <a-form-item v-if="isVolunteerCategory" label=" " class="volunteer-helper-item">
            <div class="volunteer-helper">
              <span class="volunteer-helper__desc">{{ $t('material.upload.fetchVolunteerHint') }}</span>
              <a-button type="outline" size="small" @click="openVolunteerModal">
                {{ $t('material.upload.fetchVolunteerInfo') }}
              </a-button>
            </div>
          </a-form-item>

          <a-form-item :label="$t('material.upload.description')" field="description">
            <template #extra>
              <div v-if="isVolunteerCategory" class="volunteer-extra">
                <span>{{ $t('material.upload.volunteerHoursRequired') }}</span>
                <a-tag v-if="volunteerHours" color="arcoblue">
                  {{ $t('material.upload.volunteerTotalHours', { hours: volunteerHours.total_hours }) }}h
                </a-tag>
              </div>
            </template>
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

          <a-form-item :label="$t('material.upload.tags')" field="tags">
            <a-select v-model="form.tags" multiple :placeholder="$t('material.upload.tagsPlaceholder')" allow-clear>
              <a-option value="期刊论文发表（A 类）">{{ $t('material.tag.paper.journalA') }}</a-option>
              <a-option value="会议论文收录（B 类）">{{ $t('material.tag.paper.conferenceB') }}</a-option>
              <a-option value="会议论文收录（C 类）">{{ $t('material.tag.paper.conferenceC') }}</a-option>
              <a-option value="Nature/Science/Cell 主刊论文">{{ $t('material.tag.paper.nsc') }}</a-option>
              <a-option value="Cell 子刊论文（IF≥10）">{{ $t('material.tag.paper.cellSub') }}</a-option>
              <a-option value="国家发明专利授权（第一作者）">{{ $t('material.tag.patent.nationalFirst') }}</a-option>
              <a-option value="国家发明专利授权（独立作者）">{{ $t('material.tag.patent.nationalIndep') }}</a-option>
              <a-option value="高水平中文学术期刊论文">{{ $t('material.tag.paper.chineseHigh') }}</a-option>
              <a-option value="信息与通信工程国际期刊论文">{{ $t('material.tag.paper.iceInt') }}</a-option>
              <a-option value="国家级 A + 类竞赛一等奖及以上">{{ $t('material.tag.competition.nationalAPlusFirst') }}</a-option>
              <a-option value="国家级 A 类竞赛二等奖">{{ $t('material.tag.competition.nationalASecond') }}</a-option>
              <a-option value="省级 A 类竞赛一等奖及以上">{{ $t('material.tag.competition.provincialAFirst') }}</a-option>
              <a-option value="省级 A - 类竞赛二等奖">{{ $t('material.tag.competition.provincialAMinusSecond') }}</a-option>
              <a-option value="ICPC 亚洲区域赛获奖">{{ $t('material.tag.competition.icpcAsia') }}</a-option>
              <a-option value="CCPC 竞赛获奖">{{ $t('material.tag.competition.ccpc') }}</a-option>
              <a-option value="ICPC 全球总决赛获奖">{{ $t('material.tag.competition.icpcWorld') }}</a-option>
              <a-option value="CCF CSP 认证前 0.2%（等同国一）">{{ $t('material.tag.csp.top02') }}</a-option>
              <a-option value="CCF CSP 认证前 1.5%（等同国二）">{{ $t('material.tag.csp.top15') }}</a-option>
              <a-option value="中国国际大学生创新大赛团体获奖">{{ $t('material.tag.innovation.intlGroup') }}</a-option>
              <a-option value="挑战杯学术科技作品竞赛个人获奖">{{ $t('material.tag.innovation.challengeCupIndividual') }}</a-option>
              <a-option value="挑战杯创业计划大赛团队获奖">{{ $t('material.tag.innovation.challengeCupGroup') }}</a-option>
              <a-option value="国家级创新实验计划项目（组长）">{{ $t('material.tag.project.nationalLeader') }}</a-option>
              <a-option value="国家级创新实验计划项目（成员）">{{ $t('material.tag.project.nationalMember') }}</a-option>
              <a-option value="省级创新实验计划项目（组长）">{{ $t('material.tag.project.provincialLeader') }}</a-option>
              <a-option value="省级创新实验计划项目（成员）">{{ $t('material.tag.project.provincialMember') }}</a-option>
              <a-option value="校级创新实验计划项目（组长）">{{ $t('material.tag.project.schoolLeader') }}</a-option>
              <a-option value="校级创新实验计划项目（成员）">{{ $t('material.tag.project.schoolMember') }}</a-option>
              <a-option value="创新创业训练项目结项（教务处证明）">{{ $t('material.tag.project.completionJiaowu') }}</a-option>
              <a-option value="创新创业训练项目结项（创新网截图）">{{ $t('material.tag.project.completionScreenshot') }}</a-option>
              <a-option value="国际组织一学年实习">{{ $t('material.tag.internship.intlYear') }}</a-option>
              <a-option value="国际组织半年（超一学期）实习">{{ $t('material.tag.internship.intlHalfYearPlus') }}</a-option>
              <a-option value="国际组织实习证明（满学年）">{{ $t('material.tag.internship.intlProofYear') }}</a-option>
              <a-option value="国际组织实习证明（半年）">{{ $t('material.tag.internship.intlProofHalfYear') }}</a-option>
              <a-option value="参军入伍满 1 年（含）服兵役">{{ $t('material.tag.military.oneYear') }}</a-option>
              <a-option value="参军入伍满 2 年（含）服兵役">{{ $t('material.tag.military.twoYears') }}</a-option>
              <a-option value="服兵役证明（1-2 年）">{{ $t('material.tag.military.proofOneTwo') }}</a-option>
              <a-option value="服兵役证明（2 年以上）">{{ $t('material.tag.military.proofTwoPlus') }}</a-option>
              <a-option value="志愿服务时长">{{ $t('material.tag.volunteer.hours') }}</a-option>
              <a-option value="国家级志愿服务个人表彰">{{ $t('material.tag.volunteer.nationalIndividual') }}</a-option>
              <a-option value="省级志愿服务团队表彰（队长）">{{ $t('material.tag.volunteer.provincialTeamLeader') }}</a-option>
              <a-option value="省级志愿服务团队表彰（队员）">{{ $t('material.tag.volunteer.provincialTeamMember') }}</a-option>
              <a-option value="校级志愿服务个人表彰">{{ $t('material.tag.volunteer.schoolIndividual') }}</a-option>
              <a-option value="抢险救灾志愿服务突出表现">{{ $t('material.tag.volunteer.disasterRelief') }}</a-option>
              <a-option value="国家级优秀三好学生">{{ $t('material.tag.honor.nationalMeritStudent') }}</a-option>
              <a-option value="国家级优秀共产党员">{{ $t('material.tag.honor.nationalPartyMember') }}</a-option>
              <a-option value="国家级‘自强之星’">{{ $t('material.tag.honor.nationalSelfImprovement') }}</a-option>
              <a-option value="省级优秀学生干部">{{ $t('material.tag.honor.provincialStudentCadre') }}</a-option>
              <a-option value="省级优秀团员">{{ $t('material.tag.honor.provincialLeagueMember') }}</a-option>
              <a-option value="省级社会实践优秀个人">{{ $t('material.tag.honor.provincialSocialPractice') }}</a-option>
              <a-option value="校级三好学生">{{ $t('material.tag.honor.schoolMeritStudent') }}</a-option>
              <a-option value="校级优秀学生干部">{{ $t('material.tag.honor.schoolStudentCadre') }}</a-option>
              <a-option value="国家级五四红旗团支部（集体）">{{ $t('material.tag.honor.nationalRedFlag') }}</a-option>
              <a-option value="省级优秀班集体（集体）">{{ $t('material.tag.honor.provincialClass') }}</a-option>
              <a-option value="校级优秀团支部书记">{{ $t('material.tag.honor.schoolLeagueSecretary') }}</a-option>
              <a-option value="院学生会执行主席（任职 1 学年）">{{ $t('material.tag.work.collegeChair') }}</a-option>
              <a-option value="团总支书记（任职 1 学年）">{{ $t('material.tag.work.leagueSecretary') }}</a-option>
              <a-option value="校学生会主席团成员（任职 1 学年）">{{ $t('material.tag.work.schoolPresidium') }}</a-option>
              <a-option value="班长（任职 1 学年）">{{ $t('material.tag.work.monitor') }}</a-option>
              <a-option value="团支部书记（任职半学年）">{{ $t('material.tag.work.branchSecretaryHalf') }}</a-option>
              <a-option value="党支部书记（任职 1 学年）">{{ $t('material.tag.work.partySecretary') }}</a-option>
              <a-option value="社团社长（任职 1 学年）">{{ $t('material.tag.work.clubPresident') }}</a-option>
              <a-option value="院学生会部长（辅导员打分 90+）">{{ $t('material.tag.work.collegeMinister') }}</a-option>
              <a-option value="班委（任职 1 学年）">{{ $t('material.tag.work.classCommittee') }}</a-option>
              <a-option value="系团总支副书记（任职 1 学年）">{{ $t('material.tag.work.deptDeputySecretary') }}</a-option>
              <a-option value="国际级体育团体冠军">{{ $t('material.tag.sports.intlTeamChamp') }}</a-option>
              <a-option value="国际级体育团体亚军">{{ $t('material.tag.sports.intlTeamRunnerUp') }}</a-option>
              <a-option value="国家级体育团体冠军">{{ $t('material.tag.sports.nationalTeamChamp') }}</a-option>
              <a-option value="国家级体育团体季军">{{ $t('material.tag.sports.nationalTeamThird') }}</a-option>
              <a-option value="国际级体育比赛第四至八名">{{ $t('material.tag.sports.intlTop8') }}</a-option>
              <a-option value="国家级体育个人亚军">{{ $t('material.tag.sports.nationalIndepRunnerUp') }}</a-option>
              <a-option value="国家级体育个人季军">{{ $t('material.tag.sports.nationalIndepThird') }}</a-option>
              <a-option value="厦门大学代表国际体育比赛获奖">{{ $t('material.tag.sports.xmuIntl') }}</a-option>
              <a-option value="厦门大学代表国家级体育比赛获奖">{{ $t('material.tag.sports.xmuNational') }}</a-option>
              <a-option value="体育团体项目（5 人以内）获奖">{{ $t('material.tag.sports.teamSmall') }}</a-option>
            </a-select>
          </a-form-item>
        </a-form>
      </div>

      <a-modal
        v-model:visible="showVolunteerModal"
        :title="$t('material.upload.fetchVolunteerHours')"
        @ok="handleFetchVolunteerHours"
        :confirm-loading="volunteerLoading"
      >
        <a-form :model="volunteerForm" layout="vertical">
          <a-form-item field="username" :label="$t('material.upload.volunteerUsername')">
            <a-input v-model="volunteerForm.username" :placeholder="$t('material.upload.volunteerUsernamePlaceholder')" />
          </a-form-item>
          <a-form-item field="password" :label="$t('material.upload.volunteerPassword')">
            <a-input-password v-model="volunteerForm.password" :placeholder="$t('material.upload.volunteerPasswordPlaceholder')" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="showVolunteerModal = false">{{ $t('material.upload.cancel') }}</a-button>
          <a-button type="primary" :loading="volunteerLoading" @click="handleFetchVolunteerHours">
            {{ $t('material.upload.confirm') }}
          </a-button>
        </template>
      </a-modal>

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
import { checkFileExists, uploadFile, type UploadFileResponse } from '@/api/upload'
import { fetchVolunteerHours, type VolunteerHoursResponse } from '@/api/volunteer'
import { Message } from '@arco-design/web-vue'
import type { UploadRequestOption } from '@arco-design/web-vue/es/upload'
import SparkMD5 from 'spark-md5'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

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

const VOLUNTEER_CATEGORY = '综合表现加分-志愿服务'

const fileList = ref<FileItem[]>([])
const uploadedFiles = ref<UploadedFile[]>([])
const uploading = ref(false)
const showSuccessModal = ref(false)
const showVolunteerModal = ref(false)
const volunteerLoading = ref(false)
const formRef = ref()
const { t } = useI18n()

const form = reactive({
  title: '',
  description: '',
  category: '',
  tags: [] as string[],
})

const volunteerForm = reactive({
  username: '',
  password: '',
})

const volunteerHours = ref<VolunteerHoursResponse | null>(null)

const formRules = {
  title: [{ required: true, message: t('material.upload.rules.title') }],
  category: [{ required: true, message: t('material.upload.rules.category') }],
  description: [
    {
      validator: (value: string, cb: (message?: string) => void) => {
        if (form.category === VOLUNTEER_CATEGORY) {
          if (!volunteerHours.value) {
            cb(t('material.upload.validator.volunteerInfo'))
            return
          }
          if (!value?.trim()) {
            cb(t('material.upload.validator.description'))
            return
          }
        }
        cb()
      },
    },
  ],
}

const isVolunteerCategory = computed(() => form.category === VOLUNTEER_CATEGORY)

const formValid = computed(() => {
  if (!form.title.trim() || !form.category) return false
  if (isVolunteerCategory.value) {
    return Boolean(volunteerHours.value) && Boolean(form.description.trim())
  }
  return true
})

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
    Message.error(t('material.upload.error.fileSize'))
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

  const uploadedIndex = uploadedFiles.value.findIndex((item) => item.name === file.name)
  if (uploadedIndex !== -1) {
    uploadedFiles.value.splice(uploadedIndex, 1)
  }
}

const handleFileUpload = async (options: UploadRequestOption & { fileItem?: any }) => {
  console.log('[handleFileUpload] start', options)
  const { file, fileItem } = options
  let fileObj: File

  if (file) {
    fileObj = file as File
  } else if (fileItem?.file) {
    fileObj = fileItem.file as File
  } else if ((fileItem as any)?.originFile) {
    fileObj = (fileItem as any).originFile as File
  } else {
    Message.error(t('material.upload.error.missingFile'))
    return
  }

  try {
    const md5 = await calculateFileMD5(fileObj)
    const checkResponse = await checkFileExists({ md5, filename: fileObj.name })
    const { exists: fileExists, file_id: existingFileId, url: existingFileUrl } = checkResponse.data

    if (fileExists) {
      const response: UploadFileResponse = {
        file_id: existingFileId!,
        url: existingFileUrl!,
        md5,
      }

      uploadedFiles.value.push({
        file_id: response.file_id,
        name: fileObj.name,
        size: fileObj.size,
        status: 'exists',
      })

      const target = fileList.value.find((f) => f.name === fileObj.name)
      if (target) {
        target.status = 'done'
      }
      return
    }

    const uploadResponse = await uploadFile(fileObj, (progress: number) => {
      if (progress === 100) {
        const target = fileList.value.find((f) => f.name === fileObj.name)
        if (target) {
          target.status = 'done'
        }
      }
    })

    uploadedFiles.value.push({
      file_id: uploadResponse.data.file_id,
      name: fileObj.name,
      size: fileObj.size,
      status: 'uploaded',
    })
  } catch (error) {
    console.error('[handleFileUpload] error:', error)
    Message.error(t('material.upload.error.uploadFailed'))
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

const openVolunteerModal = () => {
  volunteerForm.username = ''
  volunteerForm.password = ''
  showVolunteerModal.value = true
}

const prependVolunteerDescription = () => {
  if (!volunteerHours.value) return ''
  return t('material.upload.volunteerSummary', {
    total: volunteerHours.value.total_hours,
    credit: volunteerHours.value.credit_hours,
    honor: volunteerHours.value.honor_hours,
  })
}

const handleFetchVolunteerHours = async () => {
  if (!volunteerForm.username || !volunteerForm.password) {
    Message.error(t('material.upload.error.volunteerAuth'))
    return
  }
  volunteerLoading.value = true
  try {
    const { data } = await fetchVolunteerHours({ ...volunteerForm })
    volunteerHours.value = data
    showVolunteerModal.value = false
    console.log('[Zyh365] ', prependVolunteerDescription())
    Message.success(t('material.upload.success.volunteerFetch'))
  } catch (error) {
    Message.error(t('material.upload.error.volunteerFetch'))
  } finally {
    volunteerLoading.value = false
  }
}

const handleReset = () => {
  fileList.value = []
  uploadedFiles.value = []
  form.title = ''
  form.description = ''
  form.category = ''
  form.tags = []
  volunteerHours.value = null
  showVolunteerModal.value = false
  volunteerLoading.value = false
  formRef.value?.resetFields()
}

const handleStartRecognition = async () => {
  if (uploadedFiles.value.length === 0) {
    Message.error(t('material.upload.error.noFile'))
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
      Message.success(t('material.upload.success.recognition'))
    } else {
      Message.error(t('material.upload.error.recognitionFormat'))
    }
  } catch (err) {
    console.error(err)
    Message.error(t('material.upload.error.recognitionFailed'))
  } finally {
    uploading.value = false
  }
}

watch(
  () => form.category,
  (val) => {
    if (val !== VOLUNTEER_CATEGORY) {
      volunteerHours.value = null
    }
  }
)

const handleSubmit = async () => {
  if (uploadedFiles.value.length === 0) {
    Message.error(t('material.upload.error.noFile'))
    return
  }

  if (!formValid.value) {
    Message.error(t('material.upload.error.incompleteForm'))
    return
  }

  let descriptionToSubmit = form.description

  if (isVolunteerCategory.value) {
    const volunteerInfo = prependVolunteerDescription()
    if (!volunteerInfo) {
      Message.error(t('material.upload.error.volunteerSync'))
      return
    }
    const userDescription = form.description?.trim() ?? ''
    descriptionToSubmit = [volunteerInfo, userDescription].filter(Boolean).join('\n')
  }

  uploading.value = true
  try {
    await uploadMaterial({
      title: form.title,
      description: descriptionToSubmit,
      category: form.category,
      tags: form.tags,
      files: uploadedFiles.value.map((file) => file.file_id),
    })

    Message.success(t('material.upload.success.upload'))
    showSuccessModal.value = true
    handleReset()
  } catch (error) {
    Message.error(t('material.upload.error.uploadFailed'))
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
      max-width: 100%;
      padding: 0 20px;
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

  .volunteer-helper {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;

    &__desc {
      color: #86909c;
      font-size: 12px;
    }
  }

  .volunteer-extra {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #86909c;
    font-size: 12px;
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
