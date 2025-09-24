from fastapi import APIRouter, Depends, UploadFile, File
from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime
from auth import get_current_active_user, User
from models.response import success_response

router = APIRouter(prefix="/api/user", tags=["user-center"])


# 项目记录模型
class MyProjectRecord(BaseModel):
    id: int
    name: str
    description: str
    peopleNumber: int
    contributors: List[dict]


# 团队记录模型
class MyTeamRecord(BaseModel):
    id: int
    avatar: str
    name: str
    peopleNumber: int


# 最新活动模型
class LatestActivity(BaseModel):
    id: int
    title: str
    description: str
    avatar: str


# 基本信息模型
class BasicInfoModel(BaseModel):
    email: str
    nickname: str
    countryRegion: str
    area: str
    address: str
    profile: str


# 企业认证模型
class EnterpriseCertificationModel(BaseModel):
    accountType: int
    status: int
    time: str
    legalPerson: str
    certificateType: str
    authenticationNumber: str
    enterpriseName: str


# 认证记录
class CertificationRecord(BaseModel):
    certificationType: int
    certificationContent: str
    status: int
    time: str


class UnitCertification(BaseModel):
    enterpriseInfo: EnterpriseCertificationModel
    record: List[CertificationRecord]


# 模拟数据
fake_my_projects = [
    MyProjectRecord(
        id=1,
        name="HCI BGA项目",
        description="人机交互材料审核平台",
        peopleNumber=456,
        contributors=[
            {
                "name": "秦臻宇",
                "email": "qingzhenyu@arco.design",
                "avatar": "//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/a8c8cdb109cb051163646151a4a5083b.png~tplv-uwbnlip3yd-webp.webp",
            },
            {
                "name": "于涛",
                "email": "yuebao@arco.design",
                "avatar": "//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/a8c8cdb109cb051163646151a4a5083b.png~tplv-uwbnlip3yd-webp.webp",
            },
        ],
    ),
    MyProjectRecord(
        id=2,
        name="智能审核系统",
        description="基于AI的内容审核系统",
        peopleNumber=3,
        contributors=[
            {
                "name": "王五",
                "email": "wangwu@example.com",
                "avatar": "https://example.com/avatar3.png",
            }
        ],
    ),
]

fake_my_teams = [
    MyTeamRecord(
        id=1,
        avatar="https://example.com/team1.png",
        name="前端开发团队",
        peopleNumber=8,
    ),
    MyTeamRecord(
        id=2,
        avatar="https://example.com/team2.png",
        name="后端开发团队",
        peopleNumber=6,
    ),
]

fake_latest_activities = [
    LatestActivity(
        id=1,
        title="项目启动会议",
        description="召开了项目启动会议，讨论了项目规划和分工",
        avatar="https://example.com/activity1.png",
    ),
    LatestActivity(
        id=2,
        title="系统升级完成",
        description="完成了系统v2.0版本的升级部署",
        avatar="https://example.com/activity2.png",
    ),
]

fake_certification = UnitCertification(
    enterpriseInfo=EnterpriseCertificationModel(
        accountType=1,
        status=2,
        time="2024-01-10",
        legalPerson="张三",
        certificateType="营业执照",
        authenticationNumber="CERT123456",
        enterpriseName="示例科技有限公司",
    ),
    record=[
        CertificationRecord(
            certificationType=1,
            certificationContent="企业实名认证",
            status=2,
            time="2024-01-10",
        ),
        CertificationRecord(
            certificationType=2,
            certificationContent="法人身份验证",
            status=2,
            time="2024-01-11",
        ),
    ],
)


@router.post("/my-project/list")
async def get_my_project_list(current_user: User = Depends(get_current_active_user)):
    return success_response(data=fake_my_projects, msg="请求成功")


@router.post("/my-team/list")
async def get_my_team_list(current_user: User = Depends(get_current_active_user)):
    return success_response(data=fake_my_teams, msg="请求成功")


@router.post("/latest-activity")
async def get_latest_activity(current_user: User = Depends(get_current_active_user)):
    return success_response(data=fake_latest_activities, msg="请求成功")


@router.post("/save-info")
async def save_user_info(current_user: User = Depends(get_current_active_user)):
    return success_response(msg="用户信息保存成功")


@router.post("/certification")
async def get_certification(current_user: User = Depends(get_current_active_user)):
    return success_response(data=fake_certification, msg="请求成功")


@router.post("/upload")
async def user_upload(
    file: UploadFile = File(...), current_user: User = Depends(get_current_active_user)
):
    # 处理文件上传
    return success_response(
        data={
            "filename": file.filename,
            "size": 0,  # 实际应该计算文件大小
        },
        msg="文件上传成功",
    )
