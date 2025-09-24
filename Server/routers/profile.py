from fastapi import APIRouter, Depends
from pydantic import BaseModel
from typing import List
from datetime import datetime
from auth import get_current_active_user, User

router = APIRouter(prefix="/api", tags=["profile"])


# 基本信息响应模型
class ProfileBasicRes(BaseModel):
    status: int
    video: dict
    audio: dict


# 操作日志记录
class OperationLogRecord(BaseModel):
    key: str
    contentNumber: str
    updateContent: str
    status: int
    updateTime: str


# 模拟基本信息
fake_profile_basic = ProfileBasicRes(
    status=200,
    video={
        "mode": "custom",
        "acquisition": {"resolution": "1280x720", "frameRate": 30},
        "encoding": {
            "resolution": "1280x720",
            "rate": {"min": 300, "max": 2000, "default": 1500},
            "frameRate": 30,
            "profile": "high",
        },
    },
    audio={
        "mode": "custom",
        "acquisition": {"channels": 2},
        "encoding": {"channels": 2, "rate": 44100, "profile": "aac_low"},
    },
)

# 模拟操作日志
fake_operation_logs = [
    OperationLogRecord(
        key="1",
        contentNumber="CONT001",
        updateContent="修改了个人资料信息",
        status=1,
        updateTime="2024-01-15 14:30:00",
    ),
    OperationLogRecord(
        key="2",
        contentNumber="CONT002",
        updateContent="上传了新的材料文件",
        status=1,
        updateTime="2024-01-14 10:15:00",
    ),
    OperationLogRecord(
        key="3",
        contentNumber="CONT003",
        updateContent="审核了用户提交的材料",
        status=1,
        updateTime="2024-01-13 16:45:00",
    ),
]


@router.get("/profile/basic")
async def get_profile_basic(current_user: User = Depends(get_current_active_user)):
    return fake_profile_basic


@router.get("/operation/log")
async def get_operation_log(current_user: User = Depends(get_current_active_user)):
    return fake_operation_logs
