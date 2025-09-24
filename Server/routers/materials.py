from fastapi import APIRouter, Depends, HTTPException, Form
from typing import List, Optional
from pydantic import BaseModel
from datetime import datetime
import uuid
from auth import get_current_active_user, User
from models.response import success_response

router = APIRouter(prefix="/api/material", tags=["materials"])


# 材料模型
class Material(BaseModel):
    id: str
    title: str
    description: str
    category: str
    tags: List[str]
    fileUrl: str
    fileName: str
    fileSize: int
    status: str  # pending, approved, rejected
    uploader: str
    uploadTime: datetime
    reviewer: Optional[str] = None
    reviewTime: Optional[datetime] = None
    reviewComment: Optional[str] = None
    aiReviewResult: Optional[dict] = None


class AIReviewResult(BaseModel):
    score: int
    confidence: int
    suggestions: List[str]
    riskLevel: str  # low, medium, high


class UploadMaterialRequest(BaseModel):
    title: str
    description: str
    category: str
    tags: List[str]
    files: List[str]  # 文件UUID列表


class ReviewMaterialRequest(BaseModel):
    materialId: str
    status: str  # approved, rejected
    comment: Optional[str] = None


class MaterialFilter(BaseModel):
    status: Optional[str] = None
    category: Optional[str] = None
    uploader: Optional[str] = None
    startDate: Optional[str] = None
    endDate: Optional[str] = None


# 模拟材料数据
fake_materials = [
    Material(
        id="1",
        title="项目文档",
        description="项目详细说明文档",
        category="文档",
        tags=["项目", "文档", "说明"],
        fileUrl="/files/project-doc.pdf",
        fileName="project-doc.pdf",
        fileSize=1024000,
        status="approved",
        uploader="admin",
        uploadTime=datetime(2024, 1, 15, 10, 30),
        reviewer="reviewer1",
        reviewTime=datetime(2024, 1, 16, 14, 20),
        reviewComment="文档内容完整，符合要求",
        aiReviewResult={
            "score": 85,
            "confidence": 90,
            "suggestions": ["格式可以进一步优化", "建议添加目录"],
            "riskLevel": "low",
        },
    ),
    Material(
        id="2",
        title="设计稿",
        description="UI设计初稿",
        category="设计",
        tags=["UI", "设计", "初稿"],
        fileUrl="/files/design-draft.png",
        fileName="design-draft.png",
        fileSize=2048000,
        status="pending",
        uploader="user1",
        uploadTime=datetime(2024, 1, 20, 9, 15),
    ),
]


@router.post("/upload")
async def upload_material(
    request: UploadMaterialRequest,
    current_user: User = Depends(get_current_active_user),
):
    # 创建新材料
    new_material = Material(
        id=str(uuid.uuid4()),
        title=request.title,
        description=request.description,
        category=request.category,
        tags=request.tags,
        fileUrl=f"/files/{request.files[0] if request.files else 'default'}",  # 使用第一个文件的URL
        fileName=f"material-{str(uuid.uuid4())[:8]}",  # 生成随机文件名
        fileSize=0,  # 文件大小需要从文件存储系统获取
        status="pending",
        uploader=current_user.username,
        uploadTime=datetime.now(),
    )

    fake_materials.append(new_material)
    return success_response(data=new_material, msg="材料上传成功")


@router.post("/list")
async def get_material_list(
    filter: Optional[MaterialFilter] = None,
    current_user: User = Depends(get_current_active_user),
):
    materials = fake_materials

    # 应用筛选条件
    if filter:
        if filter.status:
            materials = [m for m in materials if m.status == filter.status]
        if filter.category:
            materials = [m for m in materials if m.category == filter.category]
        if filter.uploader:
            materials = [m for m in materials if m.uploader == filter.uploader]
        # 日期筛选需要实现

    return success_response(data=materials, msg="请求成功")


@router.post("/pending")
async def get_pending_materials(current_user: User = Depends(get_current_active_user)):
    pending_materials = [m for m in fake_materials if m.status == "pending"]
    return success_response(data=pending_materials, msg="请求成功")


@router.post("/review")
async def review_material(
    request: ReviewMaterialRequest,
    current_user: User = Depends(get_current_active_user),
):
    material = next((m for m in fake_materials if m.id == request.materialId), None)
    if not material:
        raise HTTPException(status_code=404, detail="Material not found")

    material.status = request.status
    material.reviewer = current_user.username
    material.reviewTime = datetime.now()
    material.reviewComment = request.comment

    return success_response(data=material, msg="审核成功")


@router.post("/statistics")
async def get_material_statistics(
    current_user: User = Depends(get_current_active_user),
):
    total = len(fake_materials)
    pending = len([m for m in fake_materials if m.status == "pending"])
    approved = len([m for m in fake_materials if m.status == "approved"])
    rejected = len([m for m in fake_materials if m.status == "rejected"])

    return success_response(
        data={
            "total": total,
            "pending": pending,
            "approved": approved,
            "rejected": rejected,
            "byCategory": {
                "文档": len([m for m in fake_materials if m.category == "文档"]),
                "设计": len([m for m in fake_materials if m.category == "设计"]),
                "代码": len([m for m in fake_materials if m.category == "代码"]),
            },
        },
        msg="统计信息获取成功",
    )


@router.post("/ai-review")
async def ai_review_material(
    materialId: str, current_user: User = Depends(get_current_active_user)
):
    material = next((m for m in fake_materials if m.id == materialId), None)
    if not material:
        raise HTTPException(status_code=404, detail="Material not found")

    # 模拟AI审核结果
    ai_result = {
        "score": 78,
        "confidence": 85,
        "suggestions": ["内容质量良好", "建议优化格式"],
        "riskLevel": "low",
    }

    material.aiReviewResult = ai_result
    return success_response(data=ai_result, msg="AI审核完成")


@router.post("/delete")
async def delete_material(
    materialId: str, current_user: User = Depends(get_current_active_user)
):
    global fake_materials
    material = next((m for m in fake_materials if m.id == materialId), None)
    if not material:
        raise HTTPException(status_code=404, detail="Material not found")

    fake_materials = [m for m in fake_materials if m.id != materialId]
    return success_response(msg="材料删除成功")


@router.post("/batch-review")
async def batch_review_materials(
    materialIds: List[str],
    status: str,
    comment: Optional[str] = None,
    current_user: User = Depends(get_current_active_user),
):
    results = []
    for material_id in materialIds:
        material = next((m for m in fake_materials if m.id == material_id), None)
        if material:
            material.status = status
            material.reviewer = current_user.username
            material.reviewTime = datetime.now()
            material.reviewComment = comment
            results.append(material)

    return success_response(data=results, msg="批量审核成功")
