import axios from 'axios'

export interface BaseInfoModel {
  activityName: string
  channelType: string
  promotionTime: string[]
  promoteLink: string
}
export interface ChannelInfoModel {
  advertisingSource: string
  advertisingMedia: string
  keyword: string[]
  pushNotify: boolean
  advertisingContent: string
}

export type UnitChannelModel = BaseInfoModel & ChannelInfoModel

export function submitChannelForm(data: UnitChannelModel) {
  return axios.post('/api/channel-form/submit', { data })
}

export function llmFormRecognition(data: { files: string[] }) {
  return axios.post('/api/llm/form', data)
}
