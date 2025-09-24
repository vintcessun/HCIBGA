from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

# 导入路由
from routers import (
    users,
    materials,
    messages,
    forms,
    lists,
    profile,
    user_center,
    upload,
)
from auth import app as auth_app

# 创建FastAPI应用
app = FastAPI(title="HCI BGA Server", version="1.0.0")

# 配置CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
        "http://127.0.0.1:3000",
    ],  # 前端开发服务器地址
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 包含所有路由
app.include_router(users.router)
app.include_router(materials.router)
app.include_router(messages.router)
app.include_router(forms.router)
app.include_router(lists.router)
app.include_router(profile.router)
app.include_router(user_center.router)
app.include_router(upload.router)

# 包含认证路由
app.include_router(auth_app)


# 健康检查
@app.get("/")
async def root():
    return {"message": "HCI BGA Server is running"}


@app.get("/health")
async def health_check():
    return {"status": "healthy"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
