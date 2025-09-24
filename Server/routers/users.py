from fastapi import APIRouter, Depends, HTTPException
from typing import List
from pydantic import BaseModel
from datetime import datetime, timedelta
from auth import (
    get_current_active_user,
    User,
    authenticate_user,
    fake_users_db,
    create_access_token,
)
from models.response import success_response, invalid_credentials_response

router = APIRouter(prefix="/api/user", tags=["users"])


# 用户信息响应模型
class UserInfoResponse(BaseModel):
    name: str
    avatar: str
    email: str
    job: str
    jobName: str
    organization: str
    organizationName: str
    location: str
    locationName: str
    introduction: str
    personalWebsite: str
    phone: str
    registrationDate: str
    accountId: str
    certification: int
    role: str


# 登录请求模型
class LoginRequest(BaseModel):
    username: str
    password: str


class LoginResponse(BaseModel):
    token: str
    role: str
    username: str
    message: str = "登录成功"


# 菜单项模型
class MenuItem(BaseModel):
    path: str
    name: str
    meta: dict
    children: List["MenuItem"] = []


MenuItem.update_forward_refs()

# 模拟用户数据
fake_user_info = {
    "admin": UserInfoResponse(
        name="王立群",
        avatar="https://i.gtimg.cn/club/item/face/img/2/15922_100.gif",
        email="wangliqun@email.com",
        job="frontend",
        jobName="前端艺术家",
        organization="Frontend",
        organizationName="前端",
        location="beijing",
        locationName="北京",
        introduction="人潇洒，性温存",
        personalWebsite="https://www.arco.design",
        phone="150****0000",
        registrationDate="2013-05-10 12:10:00",
        accountId="15012312300",
        certification=1,
        role="admin",
    ),
    "user": UserInfoResponse(
        name="普通用户",
        avatar="https://i.gtimg.cn/club/item/face/img/2/15922_100.gif",
        email="user@example.com",
        job="user",
        jobName="普通用户",
        organization="User",
        organizationName="用户组",
        location="shanghai",
        locationName="上海",
        introduction="普通用户介绍",
        personalWebsite="https://example.com",
        phone="150****0001",
        registrationDate="2023-01-01 10:00:00",
        accountId="15012312301",
        certification=0,
        role="user",
    ),
    "reviewer": UserInfoResponse(
        name="审核员",
        avatar="https://i.gtimg.cn/club/item/face/img/2/15922_100.gif",
        email="reviewer@example.com",
        job="reviewer",
        jobName="内容审核员",
        organization="Review",
        organizationName="审核组",
        location="guangzhou",
        locationName="广州",
        introduction="专业内容审核",
        personalWebsite="https://review.example.com",
        phone="150****0002",
        registrationDate="2022-06-15 09:30:00",
        accountId="15012312302",
        certification=1,
        role="reviewer",
    ),
}

# 模拟菜单数据
fake_menus = [
    {
        "path": "/material",
        "name": "material",
        "meta": {
            "locale": "menu.material",
            "requiresAuth": True,
            "icon": "icon-file",
            "order": 1,
        },
        "children": [
            {
                "path": "upload",
                "name": "MaterialUpload",
                "meta": {"locale": "menu.material.upload", "requiresAuth": True},
            },
            {
                "path": "list",
                "name": "MaterialList",
                "meta": {"locale": "menu.material.list", "requiresAuth": True},
            },
            {
                "path": "review",
                "name": "MaterialReview",
                "meta": {
                    "locale": "menu.material.review",
                    "requiresAuth": True,
                    "roles": ["admin", "reviewer"],
                },
            },
            {
                "path": "statistics",
                "name": "MaterialStatistics",
                "meta": {
                    "locale": "menu.material.statistics",
                    "requiresAuth": True,
                    "roles": ["admin"],
                },
            },
        ],
    }
]


@router.post("/login")
async def login(request: LoginRequest):
    # 使用真实的认证机制
    user = authenticate_user(fake_users_db, request.username, request.password)
    if not user:
        return invalid_credentials_response()

    # 生成真实的JWT token
    access_token_expires = timedelta(minutes=30)
    access_token = create_access_token(
        data={"sub": user.username, "role": user.role},
        expires_delta=access_token_expires,
    )

    return success_response(
        data={
            "token": access_token,
            "role": user.role,
            "username": user.username,
            "message": "登录成功",
        },
        msg="登录成功",
    )


@router.post("/logout")
async def logout():
    return success_response(msg="登出成功")


@router.post("/info")
async def get_user_info(current_user: User = Depends(get_current_active_user)):
    if current_user.username in fake_user_info:
        return success_response(
            data=fake_user_info[current_user.username], msg="请求成功"
        )
    raise HTTPException(status_code=404, detail="User not found")


@router.post("/menu")
async def get_user_menu(current_user: User = Depends(get_current_active_user)):
    return success_response(data=fake_menus, msg="请求成功")
