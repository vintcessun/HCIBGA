# 提交批次信息接口文档（/submit-material/submit-info）

> 更新时间：2025-11-03 23:17

## 模块概述

`/submit-material/submit-info` 页面用于展示用户可提交的批次信息，并提供提交操作。  
当前页面主要为静态演示数据，尚未接入真实接口，后续可对接提交批次 API。

---

## 页面组件

文件路径：`src/views/user/batch/submit-info.vue`

### 功能描述

- 展示可提交批次列表；
- 选择一个批次进行提交；
- 弹出提交确认提示。

---

## 组件结构

| 元素 | 类型 | 描述 |
|------|------|------|
| `<a-table />` | 组件 | 展示批次信息列表（标题、简介、状态） |
| `<a-select />` | 组件 | 选择批次下拉框 |
| `<a-button />` | 组件 | 点击触发提交事件 |
| `<a-card />` | 组件 | 承载数据展示与交互区块 |

---

## 模拟数据结构

```json
[
  {
    "title": "批次1",
    "description": "简介1",
    "status": "待提交"
  },
  {
    "title": "批次2",
    "description": "简介2",
    "status": "已提交"
  }
]
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 批次标题 |
| description | string | 批次简介 |
| status | string | 批次提交状态（待提交 / 已提交） |

---

## 模拟操作逻辑

### 提交逻辑

```ts
function handleSubmit() {
  if (!selectedBatch.value) {
    alert('请选择批次')
    return false
  }
  alert(`已提交批次: ${selectedBatch.value}`)
  return true
}
```

- 校验：若未选择批次，则弹框提示；
- 动作：提交成功后展示确认弹框。

---

## 计划接入接口

后续可在 `src/api/material.ts` 中添加以下接口定义：

```ts
/**
 * @description 提交批次材料
 * @param {string} batchTitle 批次标题
 * @returns {Promise<ApiResponse>}
 */
export function submitMaterialBatch(batchTitle) {
  return axios.post('/api/submit-material/submit-info', { batchTitle })
}
```

---

## 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "submittedBatch": "批次1",
    "timestamp": "2025-11-03T23:17:00Z"
  }
}
```

### 字段说明

| 字段 | 类型 | 描述 |
|------|------|------|
| code | number | 状态码（0 表示成功） |
| message | string | 返回消息 |
| data.submittedBatch | string | 提交的批次标题 |
| data.timestamp | string | 提交时间戳 |

---

## 错误码表

| 错误码 | 描述 | 处理建议 |
|--------|------|----------|
| 4001 | 未选择批次 | 请选择批次后重试 |
| 5001 | 提交失败 | 请检查网络或批次信息 |
| 9999 | 系统异常 | 联系系统管理员 |

---

## 附录

文档自动生成时间：`2025-11-03 23:17:59`  
维护人：系统自动生成任务（Cline）
