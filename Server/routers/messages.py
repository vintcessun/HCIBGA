from fastapi import APIRouter, Depends
from typing import List, Optional
from pydantic import BaseModel
from auth import get_current_active_user, User

router = APIRouter(prefix="/api", tags=["messages"])


# 消息记录模型
class MessageRecord(BaseModel):
    id: int
    type: str
    title: str
    subTitle: str
    avatar: Optional[str] = None
    content: str
    time: str
    status: int  # 0 未读, 1 已读
    messageType: Optional[int] = None


class MessageStatus(BaseModel):
    ids: List[int]


# 聊天记录模型
class ChatRecord(BaseModel):
    id: int
    username: str
    content: str
    time: str
    isCollect: bool


# 模拟消息数据
fake_messages = [
    MessageRecord(
        id=1,
        type="notification",
        title="系统通知",
        subTitle="新版本发布",
        avatar="https://example.com/avatar.png",
        content="系统已升级到v2.0版本，新增多项功能",
        time="2024-01-15 10:30:00",
        status=0,
        messageType=1,
    ),
    MessageRecord(
        id=2,
        type="message",
        title="用户消息",
        subTitle="张三",
        avatar="https://example.com/avatar2.png",
        content="您好，我有一个问题需要咨询",
        time="2024-01-14 15:20:00",
        status=1,
        messageType=2,
    ),
    MessageRecord(
        id=3,
        type="alert",
        title="安全提醒",
        subTitle="重要通知",
        avatar="https://example.com/alert.png",
        content="请及时更新密码以确保账户安全",
        time="2024-01-13 09:00:00",
        status=0,
        messageType=3,
    ),
]

# 模拟聊天数据
fake_chats = [
    ChatRecord(
        id=1,
        username="张三",
        content="您好，项目进展如何？",
        time="2024-01-15 14:30:00",
        isCollect=False,
    ),
    ChatRecord(
        id=2,
        username="李四",
        content="需要帮忙审核一下材料吗？",
        time="2024-01-15 13:45:00",
        isCollect=True,
    ),
    ChatRecord(
        id=3,
        username="王五",
        content="会议安排在明天下午3点",
        time="2024-01-15 12:20:00",
        isCollect=False,
    ),
]


@router.post("/message/list")
async def get_message_list(current_user: User = Depends(get_current_active_user)):
    return fake_messages


@router.post("/message/read")
async def mark_message_read(
    data: MessageStatus, current_user: User = Depends(get_current_active_user)
):
    # 标记消息为已读
    for msg in fake_messages:
        if msg.id in data.ids:
            msg.status = 1
    return fake_messages


@router.post("/chat/list")
async def get_chat_list(current_user: User = Depends(get_current_active_user)):
    return fake_chats
