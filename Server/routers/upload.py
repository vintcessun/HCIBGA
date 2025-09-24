from fastapi import APIRouter, Depends, UploadFile, File
from typing import List, Optional
import uuid
import hashlib
from datetime import datetime
from pydantic import BaseModel
from auth import get_current_active_user, User
from models.response import success_response

router = APIRouter(prefix="/api/upload", tags=["upload"])

# 模拟文件存储
fake_files_db = {}


# 文件信息模型
class FileInfo(BaseModel):
    id: str  # UUID
    md5: str  # 文件MD5哈希
    filename: str
    original_name: str
    size: int
    uploader: str
    upload_time: datetime
    url: str


# 检查文件是否已存在请求
class CheckFileRequest(BaseModel):
    md5: str
    filename: str


# 检查文件响应
class CheckFileResponse(BaseModel):
    exists: bool
    file_id: Optional[str] = None
    url: Optional[str] = None


def calculate_md5(file_content: bytes) -> str:
    """计算文件的MD5哈希"""
    return hashlib.md5(file_content).hexdigest()


@router.post("/check")
async def check_file_exists(
    request: CheckFileRequest,
    current_user: User = Depends(get_current_active_user),
):
    """检查文件是否已存在"""
    # 查找相同MD5和文件名的文件
    for file_id, file_info in fake_files_db.items():
        if file_info.md5 == request.md5 and file_info.filename == request.filename:
            return success_response(
                data={"exists": True, "file_id": file_id, "url": file_info.url},
                msg="文件已存在",
            )

    return success_response(data={"exists": False}, msg="文件不存在")


@router.post("/file")
async def upload_file(
    file: UploadFile = File(...),
    current_user: User = Depends(get_current_active_user),
):
    """上传单个文件"""
    # 读取文件内容
    file_content = await file.read()

    # 计算MD5
    file_md5 = calculate_md5(file_content)

    # 检查是否已存在相同文件
    for file_id, file_info in fake_files_db.items():
        if file_info.md5 == file_md5 and file_info.filename == file.filename:
            return success_response(
                data={"file_id": file_id, "url": file_info.url, "md5": file_md5},
                msg="文件已存在，直接返回文件信息",
            )

    # 生成文件ID
    file_id = str(uuid.uuid4())

    # 创建文件信息
    file_info = FileInfo(
        id=file_id,
        md5=file_md5,
        filename=file.filename,
        original_name=file.filename,
        size=len(file_content),
        uploader=current_user.username,
        upload_time=datetime.now(),
        url=f"/files/{file_id}/{file.filename}",
    )

    # 存储文件信息
    fake_files_db[file_id] = file_info

    return success_response(
        data={"file_id": file_id, "url": file_info.url, "md5": file_md5},
        msg="文件上传成功",
    )


@router.post("/batch")
async def upload_batch_files(
    files: List[UploadFile] = File(...),
    current_user: User = Depends(get_current_active_user),
):
    """批量上传文件"""
    results = []

    for file in files:
        file_content = await file.read()
        file_md5 = calculate_md5(file_content)

        # 检查是否已存在
        existing_file_id = None
        existing_url = None
        for file_id, file_info in fake_files_db.items():
            if file_info.md5 == file_md5 and file_info.filename == file.filename:
                existing_file_id = file_id
                existing_url = file_info.url
                break

        if existing_file_id:
            results.append(
                {
                    "file_id": existing_file_id,
                    "url": existing_url,
                    "md5": file_md5,
                    "status": "exists",
                }
            )
        else:
            # 创建新文件
            file_id = str(uuid.uuid4())
            file_info = FileInfo(
                id=file_id,
                md5=file_md5,
                filename=file.filename,
                original_name=file.filename,
                size=len(file_content),
                uploader=current_user.username,
                upload_time=datetime.now(),
                url=f"/files/{file_id}/{file.filename}",
            )
            fake_files_db[file_id] = file_info
            results.append(
                {
                    "file_id": file_id,
                    "url": file_info.url,
                    "md5": file_md5,
                    "status": "uploaded",
                }
            )

    return success_response(data=results, msg="批量文件上传完成")


@router.get("/files")
async def get_uploaded_files(
    current_user: User = Depends(get_current_active_user),
):
    """获取用户上传的文件列表"""
    user_files = [
        file_info
        for file_info in fake_files_db.values()
        if file_info.uploader == current_user.username
    ]
    return success_response(data=user_files, msg="获取文件列表成功")
