# 提交材料批次管理模块 API 文档 (`/submit-material/manage-batch`)

## 概述
本模块用于管理已提交材料的批次信息，包括批次列表查看、批次状态更新等功能。  
前端界面对应文件：`src/views/user/batch-list/manage-batch.vue`  
当前版本页面为静态数据展示，后续可接入后端接口实现动态数据加载。

---

## 基础信息

| 项目 | 内容 |
|------|------|
| 模块路径 | `/submit-material/manage-batch` |
| 主要功能 | 查看批次列表及状态 |
| 关联接口 | GET `/api/submit-material/manage-batch/list`、POST `/api/submit-material/manage-batch/update` |
| 数据来源 | 计划接入 `src/api/submit-material.ts` |
| UI组件 | ArcoDesign `a-table`, `a-card` |

---

## 接口列表

### 1. 获取批次列表

#### 请求
- **Method**: `GET`
- **URL**: `/api/submit-material/manage-batch/list`
- **参数**:

| 参数名 | 类型 | 是否必填 | 示例值 | 说明 |
|--------|------|----------|--------|------|
| page | number | 否 | 1 | 页码 |
| pageSize | number | 否 | 10 | 每页条数 |

#### 响应

| 字段名 | 类型 | 示例值 | 说明 |
|--------|------|--------|------|
| batch | string | `"批次1"` | 批次名称 |
| total | number | 10 | 批次总数 |
| status | string | `"进行中"` | 批次状态（`进行中`, `已完成`） |

#### 返回示例
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "batch": "批次1",
      "total": 10,
      "status": "进行中"
    },
    {
      "batch": "批次2",
      "total": 15,
      "status": "已完成"
    }
  ],
  "timestamp": 1730656000000
}
```

---

### 2. 更新批次状态

#### 请求
- **Method**: `POST`
- **URL**: `/api/submit-material/manage-batch/update`
- **参数**:

| 参数名 | 类型 | 是否必填 | 示例值 | 说明 |
|--------|------|----------|--------|------|
| batch | string | 是 | `"批次1"` | 批次名称 |
| status | string | 是 | `"已完成"` | 更新后的批次状态 |

#### 响应
| 字段名 | 类型 | 示例值 | 说明 |
|--------|------|--------|------|
| success | boolean | true | 操作是否成功 |
| message | string | `"状态更新成功"` | 系统提示信息 |
| updatedAt | string | `"2025-11-03T23:30:00Z"` | 更新时间 |

#### 返回示例
```json
{
  "code": 0,
  "message": "状态更新成功",
  "data": {
    "success": true,
    "updatedAt": "2025-11-03T23:30:00Z"
  },
  "timestamp": 1730656200000
}
```

---

## 状态说明

| 状态值 | 中文说明 |
|---------|-----------|
| `进行中` | 批次正在处理 |
| `已完成` | 批次处理完成 |

---

## 错误码定义
| 错误码 | 说明 |
|--------|------|
| 4001 | 批次不存在 |
| 4002 | 状态值无效 |
| 5000 | 服务端内部错误 |

---

## 时间戳规范
- 服务端响应统一包含 `timestamp` 字段，单位：毫秒
- 格式参考 UTC+8，示例：`1730656200000`
- 前端通过 `new Date(timestamp).toLocaleString('zh-CN')` 格式化展示

---

## 接口调用流程

1. 页面加载时触发 `getManageBatchList` → 拉取批次列表。  
2. 用户在界面中选择批次并更新状态 → 调用 `updateManageBatchStatus` 完成状态更新。  
3. 更新完成后刷新数据源。

---

**最后更新时间**：2025-11-03 23:30 (UTC+8)
