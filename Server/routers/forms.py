from fastapi import APIRouter, Depends
from pydantic import BaseModel
from typing import List
from auth import get_current_active_user, User

router = APIRouter(prefix="/api", tags=["forms"])


# 表单模型
class BaseInfoModel(BaseModel):
    activityName: str
    channelType: str
    promotionTime: List[str]
    promoteLink: str


class ChannelInfoModel(BaseModel):
    advertisingSource: str
    advertisingMedia: str
    keyword: List[str]
    pushNotify: bool
    advertisingContent: str


class UnitChannelModel(BaseInfoModel, ChannelInfoModel):
    pass


@router.post("/channel-form/submit")
async def submit_channel_form(
    data: UnitChannelModel, current_user: User = Depends(get_current_active_user)
):
    # 这里应该保存表单数据到数据库
    # 暂时返回成功消息
    return {"success": True, "message": "表单提交成功", "data": data.model_dump()}
