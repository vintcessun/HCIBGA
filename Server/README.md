# HCI BGA Server

基于 Python + FastAPI 的后端服务器，为人机交互材料审核平台提供API支持。

## 功能特性

- ✅ 用户认证和授权系统
- ✅ 材料管理（上传、审核、统计）
- ✅ 消息和通知系统
- ✅ 表单处理
- ✅ 列表数据管理
- ✅ 个人资料管理
- ✅ 用户中心功能
- ✅ CORS支持
- ✅ JWT令牌认证

## 技术栈

- **Python 3.11+**
- **FastAPI** - 高性能Web框架
- **Uvicorn** - ASGI服务器
- **SQLAlchemy** - ORM（预留）
- **Pydantic** - 数据验证
- **JWT** - 身份认证
- **uv** - Python包管理

## 安装和运行

### 1. 创建虚拟环境

```bash
uv venv .venv
```

### 2. 激活虚拟环境

Windows:
```bash
.venv\Scripts\activate
```

Linux/Mac:
```bash
source .venv/bin/activate
```

### 3. 安装依赖

```bash
uv sync
```

### 4. 运行服务器

```bash
python run.py
```

或者直接运行：

```bash
uvicorn main:app --reload --host 0.0.0.0 --port 8000
```

## API文档

服务器启动后，访问以下地址查看API文档：

- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## 主要API端点

### 用户认证
- `POST /api/user/login` - 用户登录
- `POST /api/user/logout` - 用户登出
- `POST /api/user/info` - 获取用户信息
- `POST /api/user/menu` - 获取用户菜单

### 材料管理
- `POST /api/material/upload` - 上传材料
- `POST /api/material/list` - 获取材料列表
- `POST /api/material/pending` - 获取待审核材料
- `POST /api/material/review` - 审核材料
- `POST /api/material/statistics` - 获取统计信息
- `POST /api/material/ai-review` - AI自动审核
- `POST /api/material/delete` - 删除材料
- `POST /api/material/batch-review` - 批量审核

### 消息系统
- `POST /api/message/list` - 获取消息列表
- `POST /api/message/read` - 标记消息已读
- `POST /api/chat/list` - 获取聊天记录

### 其他API
- `POST /api/channel-form/submit` - 提交渠道表单
- `GET /api/list/policy` - 获取策略列表
- `GET /api/list/quality-inspection` - 获取质检列表
- `GET /api/list/the-service` - 获取服务列表
- `GET /api/list/rules-preset` - 获取预设规则
- `GET /api/profile/basic` - 获取基本信息
- `GET /api/operation/log` - 获取操作日志
- `POST /api/user/my-project/list` - 获取我的项目
- `POST /api/user/my-team/list` - 获取我的团队
- `POST /api/user/latest-activity` - 获取最新活动
- `POST /api/user/save-info` - 保存用户信息
- `POST /api/user/certification` - 获取认证信息
- `POST /api/user/upload` - 用户文件上传

## 默认用户

- 用户名: `admin`，密码: `admin`
- 用户名: `user`，密码: `user`  
- 用户名: `reviewer`，密码: `reviewer`

## 项目结构

```
Server/
├── main.py              # 主应用文件
├── run.py              # 启动脚本
├── pyproject.toml      # 项目配置和依赖
├── uv.lock            # 依赖锁文件
├── routers/           # 路由模块
│   ├── users.py       # 用户相关路由
│   ├── materials.py   # 材料管理路由
│   ├── messages.py    # 消息相关路由
│   ├── forms.py       # 表单处理路由
│   ├── lists.py       # 列表数据路由
│   ├── profile.py     # 个人资料路由
│   └── user_center.py # 用户中心路由
└── README.md          # 项目说明
```

## 开发说明

1. 所有API都支持JWT认证
2. 使用Pydantic进行数据验证
3. 支持CORS，前端可运行在localhost:3000
4. 目前使用内存数据模拟，生产环境应连接数据库

## 下一步计划

- [ ] 集成真实数据库（SQLite/PostgreSQL）
- [ ] 实现文件上传存储
- [ ] 添加单元测试
- [ ] 实现真正的AI审核集成
- [ ] 添加日志系统
- [ ] 配置环境变量管理
- [ ] 实现用户权限管理

## 许可证

MIT License
