from fastapi import APIRouter, Depends
from pydantic import BaseModel
from typing import List, Optional
from auth import get_current_active_user, User

router = APIRouter(prefix="/api/list", tags=["lists"])


# 策略记录模型
class PolicyRecord(BaseModel):
    id: str
    number: int
    name: str
    contentType: str  # img, horizontalVideo, verticalVideo
    filterType: str  # artificial, rules
    count: int
    status: str  # online, offline
    createdTime: str


class PolicyParams(BaseModel):
    current: int
    pageSize: int
    id: Optional[str] = None
    number: Optional[int] = None
    name: Optional[str] = None
    contentType: Optional[str] = None
    filterType: Optional[str] = None
    count: Optional[int] = None
    status: Optional[str] = None
    createdTime: Optional[str] = None


class PolicyListRes(BaseModel):
    list: List[PolicyRecord]
    total: int


# 服务记录模型
class ServiceRecord(BaseModel):
    id: int
    title: str
    description: str
    name: Optional[str] = None
    actionType: Optional[str] = None
    icon: Optional[str] = None
    data: Optional[List[dict]] = None
    enable: Optional[bool] = None
    expires: Optional[bool] = None


# 模拟数据
fake_policies = [
    PolicyRecord(
        id="1",
        number=1001,
        name="图片审核策略",
        contentType="img",
        filterType="artificial",
        count=150,
        status="online",
        createdTime="2024-01-10",
    ),
    PolicyRecord(
        id="2",
        number=1002,
        name="视频审核策略",
        contentType="horizontalVideo",
        filterType="rules",
        count=89,
        status="offline",
        createdTime="2024-01-12",
    ),
    PolicyRecord(
        id="3",
        number=1003,
        name="内容过滤策略",
        contentType="verticalVideo",
        filterType="artificial",
        count=203,
        status="online",
        createdTime="2024-01-15",
    ),
]

fake_inspection_list = [
    ServiceRecord(
        id=1,
        title="内容质量检查",
        description="检查内容质量和完整性",
        name="quality-inspection",
        actionType="view",
        icon="check-circle",
    ),
    ServiceRecord(
        id=2,
        title="安全审核",
        description="安全检查和安全策略",
        name="security-check",
        actionType="edit",
        icon="shield",
    ),
]

fake_service_list = [
    ServiceRecord(
        id=1,
        title="用户服务",
        description="用户相关的服务功能",
        name="user-service",
        actionType="manage",
        icon="user",
    ),
    ServiceRecord(
        id=2,
        title="系统服务",
        description="系统管理和维护服务",
        name="system-service",
        actionType="config",
        icon="setting",
    ),
]

fake_rules_preset = [
    ServiceRecord(
        id=1,
        title="基础规则集",
        description="基础的内容审核规则",
        name="basic-rules",
        actionType="apply",
        icon="file-text",
    ),
    ServiceRecord(
        id=2,
        title="高级规则集",
        description="高级的内容过滤规则",
        name="advanced-rules",
        actionType="customize",
        icon="filter",
    ),
]


@router.get("/policy")
async def get_policy_list(
    current: int = 1,
    pageSize: int = 10,
    id: Optional[str] = None,
    number: Optional[int] = None,
    name: Optional[str] = None,
    contentType: Optional[str] = None,
    filterType: Optional[str] = None,
    count: Optional[int] = None,
    status: Optional[str] = None,
    createdTime: Optional[str] = None,
):
    # 简单的分页和筛选实现
    start = (current - 1) * pageSize
    end = start + pageSize
    policies = fake_policies[start:end]

    # 应用筛选条件
    filtered_policies = fake_policies
    if id:
        filtered_policies = [p for p in filtered_policies if p.id == id]
    if number:
        filtered_policies = [p for p in filtered_policies if p.number == number]
    if name:
        filtered_policies = [
            p for p in filtered_policies if name.lower() in p.name.lower()
        ]
    if contentType:
        filtered_policies = [
            p for p in filtered_policies if p.contentType == contentType
        ]
    if filterType:
        filtered_policies = [p for p in filtered_policies if p.filterType == filterType]
    if count:
        filtered_policies = [p for p in filtered_policies if p.count == count]
    if status:
        filtered_policies = [p for p in filtered_policies if p.status == status]
    if createdTime:
        filtered_policies = [
            p for p in filtered_policies if p.createdTime == createdTime
        ]

    # 应用分页
    policies = filtered_policies[start:end]

    return PolicyListRes(list=policies, total=len(filtered_policies))


@router.get("/quality-inspection")
async def get_inspection_list():
    return fake_inspection_list


@router.get("/the-service")
async def get_service_list():
    return fake_service_list


@router.get("/rules-preset")
async def get_rules_preset_list():
    return fake_rules_preset
