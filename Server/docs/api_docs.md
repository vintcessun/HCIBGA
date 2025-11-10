# API 功能与执行流程说明

本文档描述了现有所有 API 的功能、执行流程以及所需第三方库。

---

## `/api/user/save`
**功能**: 保存用户设置信息（使用 SQLite 持久化）  
**执行流程**:
1. 接收 JSON 格式的用户信息（姓名、邮箱、头像）。
2. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
3. 创建 `user_settings` 表（如不存在）。
4. 插入一条包含姓名、邮箱、头像的记录。
5. 返回保存成功的响应。  
**第三方库**:
- `github.com/mattn/go-sqlite3`: SQLite 驱动。

---

## `/api/user/auth`
**功能**: 用户认证并记录登录信息  
**执行流程**:
1. 接收包含用户名的 JSON 数据。
2. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
3. 创建 `login_records` 表（如不存在）。
4. 插入一条包含用户名、IP、当前时间的记录。
5. 返回认证成功信息。  
**第三方库**:
- `github.com/mattn/go-sqlite3`: SQLite 驱动。

---

## `/api/user/upload`
**功能**: 上传用户文件  
**执行流程**:
1. 解析 `multipart/form-data` 上传内容。
2. 将文件保存至临时目录 `hci_user_upload`。
3. 返回保存路径和成功状态。  
**第三方库**: 无。

---

## `/api/export/students/excel`
**功能**: 从数据库导出学生信息为 CSV 文件（Excel 可打开）  
**执行流程**:
1. 从 SQLite 数据库（`students` 表）中读取所有学生信息（包括学号、姓名、专业、班级、成绩）。  
2. 设置响应头为 `application/vnd.ms-excel`，并指定下载文件名（如 `students_时间戳.csv`）。  
3. 使用 `encoding/csv` 生成 CSV 文件流：  
   - 首行写入表头（`id, name, studentId, major, class, score`）。  
   - 随后按顺序写入数据库中每条学生记录。  
4. 处理可能的异常情况，例如数据库连接错误或文件写入异常，并返回 HTTP 状态码与错误信息。  
5. 浏览器将直接下载生成的 CSV 文件。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。  
- Go 标准库 `encoding/csv`: 用于写入 CSV 格式文件。

## `/api/export/students/txt`
**功能**: 从数据库导出学生信息为 TXT 文件（纯文本格式）  
**执行流程**:
1. 从 SQLite 数据库（`students` 表）中读取所有学生信息（包括学号、姓名、专业、班级、成绩）。  
2. 设置响应头为 `text/plain; charset=utf-8`，并指定下载文件名（如 `students_时间戳.txt`）。  
3. 按每位学生一行的格式输出详细信息，例如：  
   `ID: 1, 姓名: 张三, 学号: 2021001, 专业: 计算机科学, 班级: 计科1班, 成绩: 5.0`  
4. 处理可能的异常情况，例如数据库连接错误或 I/O 写入异常，并返回 HTTP 状态码与错误信息。  
5. 浏览器将直接下载生成的 TXT 文件。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。  
- Go 标准库：用于处理文本输出。

---

## `/api/info/import/excel`
**功能**: 接收并保存 Excel 文件，并将其内容导入到数据库  
**执行流程**:
1. 接收并解析上传的 Excel 文件。
2. 检查扩展名 `.xls` / `.xlsx` 是否有效，不符合则返回错误。
3. 保存至临时目录 `hci_info_import`。
4. 使用如 `excelize` 等库读取表格内容。
5. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
6. 将 Excel 中的每一行数据插入到目标数据表（例如 `students`），包括字段映射与类型转换。
7. 处理可能的异常情况，例如文件解析错误、数据库连接或插入异常，并返回 HTTP 状态码与错误信息。
8. 返回带有任务 ID 或导入结果的响应。  
**第三方库**:  
- `github.com/xuri/excelize/v2`: 用于读取 Excel 文件内容。  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。

## `/api/info/import/txt`
**功能**: 接收 TXT 文件或直接文本数据，并导入到数据库  
**执行流程**:
1. 接收 TXT 文件或 `text` 字段内容。
2. 验证文件扩展名是否为 `.txt`，不符合则返回错误。
3. 保存至临时目录 `hci_info_import`。
4. 打开 TXT 文件并按行读取，解析每行数据（例如以逗号或制表符分隔的字段）。
5. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
6. 将 TXT 中解析出的每行记录插入到目标数据表（例如 `students`），包括字段映射与类型转换。
7. 处理可能的异常情况，例如文件读取错误、解析异常、数据库连接或插入失败，并返回 HTTP 状态码与错误信息。
8. 返回成功信息及导入结果（例如导入的行数）。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。  
- Go 标准库：用于文件读取与字符串解析。

---

## `/api/material/list`
**功能**: 从数据库获取真实材料列表（包含 ID 列表）  
**执行流程**:
1. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
2. 查询 `materials` 表，获取所有材料记录（包括 `id`、`name`、`type`、`status` 等字段）。
3. 将结果转换为 JSON 数组返回，其中必须包含材料 `id` 列表，便于后续操作。
4. 处理可能的异常情况，例如数据库连接错误或查询失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。

## `/api/material/review`
**功能**: 审核指定材料并保存审核意见  
**执行流程**:
1. 接收包含材料 `id` 和审核意见（`comment`）的 JSON 请求。
2. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
3. 更新 `materials` 表中对应 `id` 的记录，写入审核意见并修改状态（如 `status` = "reviewed"）。
4. 返回审核成功信息及更新后的记录。
5. 处理可能的异常情况，例如数据库连接错误、记录不存在或更新失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。

## `/api/material/{id}`
**功能**: 删除指定材料记录  
**执行流程**:
1. 从 URL 路径参数中提取材料 `id`。
2. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
3. 在 `materials` 表中删除对应 `id` 的记录。
4. 返回删除成功信息及被删除的 `id`。
5. 处理可能的异常情况，例如数据库连接错误、记录不存在或删除失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。

---

## `/api/upload/check`
**功能**: 检查上传服务状态  
**执行流程**:
1. 返回 “ready” 状态信息。  
**第三方库**: 无。

## `/api/upload/file`
**功能**: 上传材料文件  
**执行流程**:
1. 解析上传的文件并保存至 `hci_upload_material` 临时目录。
2. 返回文件信息和保存路径。  
**第三方库**: 无。

## `/api/material/llm-fill`
**功能**: 材料自动填充（模拟 LLM）  
**执行流程**:
1. 接收材料 ID。
2. 返回自动填充数据。  
**第三方库**: 无。

## `/api/material/upload`
**功能**: 材料上传完成记录（引用已上传文件 ID）  
**执行流程**:
1. 接收包含 `file_id`（来自 `/api/upload/file` 返回的唯一 ID）及材料相关信息的 JSON 数据。
2. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
3. 查询 `uploaded_files` 表，验证 `file_id` 是否存在。
4. 在 `materials` 表中插入一条新记录，关联该 `file_id` 及材料信息（如 `name`、`type`、`status`）。
5. 返回材料记录的 ID 及状态信息。
6. 处理可能的异常情况，例如 `file_id` 不存在、数据库连接错误或插入失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  

---

## `/api/submit-material/batch-list`
**功能**: 获取批量材料列表（真实数据）  
**执行流程**:
1. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。  
2. 查询与用户信息同库的 `submit_material_batches` 表，获取所有批次记录（包括 `id`、`name`、`status`、`reviewer` 等字段）。  
3. 将结果转换为 JSON 数组返回。  
4. 处理可能的异常情况，例如数据库连接错误、查询失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。  

---

## `/api/user/projects`
**功能**: 获取用户项目列表（真实数据）  
**执行流程**:
1. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。  
2. 查询与用户信息同库的 `user_projects` 表，获取所有项目记录（包括 `id`、`name` 等字段）。  
3. 将结果转换为 JSON 数组返回。  
4. 处理可能的异常情况，例如数据库连接错误、查询失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。  

## `/api/user/activities`
**功能**: 获取用户活动记录（真实数据）  
**执行流程**:
1. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。  
2. 查询与用户信息同库的 `user_activities` 表，获取所有活动记录（包括 `id`、`content`、`time` 等字段）。  
3. 将结果转换为 JSON 数组返回。  
4. 处理可能的异常情况，例如数据库连接错误、查询失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:  
- `github.com/mattn/go-sqlite3`: SQLite 驱动。  

## `/api/user/teams`
**功能**: 获取用户所在团队信息（真实数据）  
**执行流程**:
1. 使用 `database/sql` 连接 SQLite 数据库（`github.com/mattn/go-sqlite3`）。
2. 查询与用户信息同库的 `user_teams` 表，获取所有团队记录（包括 `id`、`name`、`role`、`created_at` 等字段）。
3. 将结果转换为 JSON 数组返回。
4. 处理可能的异常情况，例如数据库连接错误、查询失败，并返回 HTTP 状态码与错误信息。  
**第三方库**:
