# 材料审核模块 API 文档 (`/material/review`)

## 概述
本模块用于管理与审核待审核材料，包括单条快速审核、详细审核及批量操作等功能。  
前端界面对应文件：`src/views/material/review/index.vue`  
调用接口：`getPendingMaterials`, `reviewMaterial`, `batchReviewMaterials`

---

## 基础信息

| 项目 | 内容 |
|------|------|
| 模块路径 | `/material/review` |
| 主要功能 | 审核待提交的材料，支持查看、通过、拒绝、批量操作 |
| 关联接口 | GET `/api/material/pending`、POST `/api/material/review`、POST `/api/material/batch` |
| 数据来源 | `src/api/material.ts` |
| UI组件 | ArcoDesign `a-table`, `a-modal`, `a-form`, `a-tag`, `a-select`, `a-button` |

---

## 接口列表

### 1. 获取待审核材料列表

#### 请求
- **Method**: `GET`
- **URL**: `/api/material/pending`
- **参数**:

| 参数名 | 类型 | 是否必填 | 示例值 | 说明 |
|--------|------|----------|--------|------|
| page | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

#### 响应

| 字段名 | 类型 | 示例值 | 说明 |
|--------|------|--------|------|
| id | string | `"mat_001"` | 材料唯一ID |
| title | string | `"关于材料管理规范"` | 材料标题 |
| category | string | `"document"` | 材料分类（`document`, `image`, `video`, …） |
| fileSize | number | 102400 | 文件大小（单位字节） |
| uploader | string | `"张三"` | 上传者姓名 |
| uploadTime | string | `"2025-09-18T14:00:00Z"` | 上传时间（ISO格式） |
| status | string | `"pending"` | 状态（`pending`, `approved`, `rejected`） |

#### 返回示例
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "mat_001",
      "title": "关于材料管理规范",
      "category": "document",
      "fileSize": 102400,
      "uploader": "张三",
      "uploadTime": "2025-09-18T14:00:00Z",
      "status": "pending"
    }
  ],
  "timestamp": 1730652000000
}
```

---

### 2. 单条材料审核

#### 请求
- **Method**: `POST`
- **URL**: `/api/material/review`
- **参数**:

| 参数名 | 类型 | 是否必填 | 示例值 | 说明 |
|--------|------|----------|--------|------|
| materialId | string | 是 | `"mat_001"` | 材料唯一ID |
| status | string | 是 | `"approved"` | 审核状态（`approved`通过、`rejected`拒绝） |
| comment | string | 否 | `"快速审核通过"` | 审核意见/拒绝原因 |

#### 响应
| 字段名 | 类型 | 示例值 | 说明 |
|--------|------|--------|------|
| success | boolean | true | 操作是否成功 |
| message | string | `"审核成功"` | 系统提示信息 |
| reviewer | string | `"李四"` | 当前操作用户 |
| reviewedAt | string | `"2025-11-03T23:21:00Z"` | 审核时间 |
| status | string | `"approved"` | 最新审核状态 |

#### 返回示例
```json
{
  "code": 0,
  "message": "审核通过",
  "data": {
    "success": true,
    "reviewer": "李四",
    "reviewedAt": "2025-11-03T23:21:01Z",
    "status": "approved"
  },
  "timestamp": 1730652061000
}
```

---

### 3. 批量审核材料

#### 请求
- **Method**: `POST`
- **URL**: `/api/material/batch`
- **参数**:

| 参数名 | 类型 | 是否必填 | 示例值 | 说明 |
|--------|------|----------|--------|------|
| materialIds | array<string> | 是 | `["mat_001", "mat_002"]` | 材料ID数组 |
| status | string | 是 | `"rejected"` | 审核状态 |
| comment | string | 否 | `"批量审核完成"` | 统一审核意见 |

#### 响应
| 字段名 | 类型 | 示例值 | 说明 |
|--------|------|--------|------|
| successCount | number | 2 | 成功审核数 |
| failedCount | number | 0 | 失败审核数 |
| status | string | `"rejected"` | 执行的审核状态 |

#### 返回示例
```json
{
  "code": 0,
  "message": "批量审核完成",
  "data": {
    "successCount": 2,
    "failedCount": 0,
    "status": "rejected"
  },
  "timestamp": 1730652107000
}
```

---

## 状态说明

| 状态值 | 中文说明 |
|---------|-----------|
| `pending` | 待审核 |
| `approved` | 审核通过 |
| `rejected` | 审核拒绝 |

---

## 错误码定义
| 错误码 | 说明 |
|--------|------|
| 4001 | 材料不存在 |
| 4002 | 审核状态无效 |
| 4003 | 批量操作数量超限 |
| 5000 | 服务端内部错误 |

---

## 时间戳规范
- 服务端响应统一包含 `timestamp` 字段，单位：毫秒
- 格式参考 UTC+8，示例：`1730652107000`
- 前端通过 `new Date(timestamp).toLocaleString('zh-CN')` 格式化展示

---

## 接口调用流程

1. 页面加载时触发 `getPendingMaterials` → 拉取待审核列表。  
2. 用户点击“快速通过/拒绝” → 调用 `reviewMaterial` 完成单条审核。  
3. 若选中多条材料 → 调用 `batchReviewMaterials` 批量处理。  
4. 审核完成后刷新数据源。

---

## 示例前端交互截图
组件位置：`src/views/material/review/index.vue`  
操作区：`快速通过 / 拒绝 / 批量审 / 查看详情` 触发对应接口调用。  

---

**最后更新时间**：2025-11-03 23:23 (UTC+8)
