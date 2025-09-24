#!/usr/bin/env python3
"""
HCI BGA Server 启动脚本
使用: python run.py
"""

import uvicorn
from main import app

if __name__ == "__main__":
    print("正在启动 HCI BGA Server...")
    print("服务器地址: http://localhost:8000")
    print("API文档: http://localhost:8000/docs")
    print("按 Ctrl+C 停止服务器")

    uvicorn.run(
        "main:app", host="0.0.0.0", port=8000, reload=True
    )  # 开发模式下自动重载
