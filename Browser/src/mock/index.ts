import Mock from 'mockjs'

import './material'
import './message-box'
import './upload'
import './user'

import '@/views/list/card/mock'
import '@/views/list/search-table/mock'

import '@/views/form/step/mock'

import '@/views/profile/basic/mock'

import '@/views/user/info/mock'
import '@/views/user/setting/mock'
/** simple end */

Mock.setup({
  timeout: '600-1000',
})
