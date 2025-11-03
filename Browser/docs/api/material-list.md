# 材料列表模块接口文档

## 模块简介
`src/views/material/list/index.vue` 用于展示、筛选、查看和审核材料列表，支持按状态、类别、上传者、日期范围进行过滤，并提供删除、批量删除、导出及审核功能。

---

## 接口定义

| 接口名称 | 方法 | 路径 | 描述 |
|-----------|--------|----------------|----------------|
| 获取材料列表 | `GET` | `/api/material/list` | 按条件获取材料列表 |
| 删除材料 | `DELETE` | `/api/material/{id}` | 删除指定材料 |
| 审核材料 | `POST` | `/api/material/review` | 审核指定材料（通过/拒绝） |

---

## 请求参数

### 1. 获取材料列表 `/api/material/list`
**Query 参数:**
| 字段 | 类型 | 示例值 | 说明 |
|------|------|------|------|
| status | string | approved | 材料状态（pending/approved/rejected） |
| category | string | document | 材料类别 |
| uploader | string | 张三 | 上传者姓名 |
| startDate | string | 2025-11-01 | 起始日期 |
| endDate | string | 2025-11-03 | 结束日期 |

**响应:**
```json
[
  {
    "id": "uuid",
    "title": "项目申报材料",
    "category": "document",
    "status": "approved",
    "fileSize": 102400,
    "uploader": "张三",
    "uploadTime": "2025-11-03T10:00:00Z",
    "reviewer": "李四",
    "reviewTime": "2025-11-03T12:00:00Z"
  }
]
```

---

### 2. 删除材料 `/api/material/{id}`
**路径参数:**
| 字段 | 类型 | 示例值 | 说明 |
|------|------|------|------|
| id | string | uuid | 材料唯一标识 |

**响应:**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 3. 审核材料 `/api/material/review`
**请求体:**
```json
{
  "materialId": "uuid",
  "status": "approved",
  "comment": "材料符合要求"
}
```

**响应:**
```json
{
  "code": 200,
  "message": "审核成功"
}
```

---

## 错误示例

| 状态码 | 说明 |
|---------|--------|
| 400 | 参数错误 |
| 404 | 材料不存在 |
| 500 | 服务器错误 |

---

## 操作流程图

```mermaid
flowchart TD
    A[进入材料列表页] --> B[选择筛选条件]
    B --> C[点击搜索按钮]
    C --> D[调用 /api/material/list 获取数据]
    D --> E[渲染表格]
    E --> F{选择操作}
    F -- 查看 --> G[打开详情模态框]
    F -- 删除 --> H[调用 /api/material/{id} 删除]
    F -- 审核 --> I[打开审核模态框]
    I --> J[提交 /api/material/review]
    J --> K[刷新列表]
```

---

## 接口调用依赖关系
| 文件 | 引用函数 | 来源 |
|------|-----------|------|
| `getMaterialList()` | 获取材料列表 | `@/api/material.ts` |
| `deleteMaterial()` | 删除材料 | `@/api/material.ts` |
| `reviewMaterial()` | 审核材料 | `@/api/material.ts` |

---

## 最佳实践
- 使用分页与筛选条件减少数据量。
- 审核操作应记录审核意见。
- 删除操作需二次确认以防误删。

---
文档生成时间：2025-11-03T22:26:00+08:00
