import localeMessageBox from '@/components/message-box/locale/zh-CN'
import localeLogin from '@/views/login/locale/zh-CN'

import localeCardList from '@/views/list/card/locale/zh-CN'
import localeSearchTable from '@/views/list/search-table/locale/zh-CN'

import localeGroupForm from '@/views/form/group/locale/zh-CN'
import localeStepForm from '@/views/form/step/locale/zh-CN'

import localeBasicProfile from '@/views/profile/basic/locale/zh-CN'

import localeError from '@/views/result/error/locale/zh-CN'
import localeSuccess from '@/views/result/success/locale/zh-CN'

import locale403 from '@/views/exception/403/locale/zh-CN'
import locale404 from '@/views/exception/404/locale/zh-CN'
import locale500 from '@/views/exception/500/locale/zh-CN'

import localeUserInfo from '@/views/user/info/locale/zh-CN'
import localeUserSetting from '@/views/user/setting/locale/zh-CN'
/** simple end */
import localeSettings from './zh-CN/settings'

export default {
  'menu.dashboard': '仪表盘',
  'menu.server.dashboard': '仪表盘-服务端',
  'menu.server.workplace': '工作台-服务端',
  'menu.server.monitor': '实时监控-服务端',
  'menu.list': '列表页',
  'menu.result': '结果页',
  'menu.exception': '异常页',
  'menu.form': '表单页',
  'menu.profile': '详情页',
  'menu.visualization': '数据可视化',
  'menu.user': '个人中心',
  'menu.arcoWebsite': 'Arco Design',
  'menu.faq': '常见问题',
  'menu.material': '材料管理',
  'menu.material.upload': '上传材料',
  'menu.material.list': '材料列表',
  'menu.material.review': '材料审核',
  'menu.material.statistics': '材料统计',
  'material.upload.files': '上传文件',
  'material.upload.materialInfo': '材料信息',
  'material.upload.uploadedFiles': '已上传文件',
  'material.upload.clickToUpload': '点击上传材料',
  'material.upload.supportFormat': '支持格式: PDF, DOC, DOCX, PPT, PPTX, XLS, XLSX, PNG, JPG等',
  'material.upload.title': '标题',
  'material.upload.titlePlaceholder': '请输入材料标题',
  'material.upload.description': '描述',
  'material.upload.descriptionPlaceholder': '请输入材料描述',
  'material.upload.category': '分类',
  'material.upload.categoryPlaceholder': '请选择分类',
  'material.upload.tags': '标签',
  'material.upload.tagsPlaceholder': '请选择标签',
  'material.upload.submit': '提交',
  'material.upload.reset': '重置',
  'material.upload.successTitle': '上传成功',
  'material.upload.successMessage': '材料已成功上传，等待审核',
  'material.upload.confirm': '确定',
  'material.list.status': '状态',
  'material.list.statusPlaceholder': '请选择状态',
  'material.list.category': '分类',
  'material.list.categoryPlaceholder': '请选择分类',
  'material.list.uploader': '上传者',
  'material.list.uploaderPlaceholder': '请输入上传者',
  'material.list.dateRange': '上传时间',
  'material.list.search': '搜索',
  'material.list.reset': '重置',
  'material.list.batchDelete': '批量删除',
  'material.list.export': '导出',
  'navbar.docs': '文档中心',
  'navbar.action.locale': '切换为中文',
  ...localeSettings,
  ...localeMessageBox,
  ...localeLogin,
  ...localeSearchTable,
  ...localeCardList,
  ...localeStepForm,
  ...localeGroupForm,
  ...localeBasicProfile,
  ...localeSuccess,
  ...localeError,
  ...locale403,
  ...locale404,
  ...locale500,
  ...localeUserInfo,
  ...localeUserSetting,
  /** simple end */
}
