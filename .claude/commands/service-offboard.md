# Command: /service:offboard

> **说明**：下线/删除某个服务时，清理所有相关的 Harness 工程制品，保持制品与实际服务一致。
> **触发**：用户输入 `/service:offboard {service-name}`

---

## 执行步骤

### Step 1: 确认下线意图

```
⚠️ 即将下线服务：{service-name}

将执行以下操作：
  - 删除 context/project/api/{service-name}/ 知识库目录
  - 从 .service-matrix/dependencies.yaml 中移除该服务
  - 从 context/project/api/INDEX.md 中移除该服务行和依赖关系
  - 从 context/team/error-code.md 中标记该错误码段为已废弃

⚠️ 注意：本命令不会删除业务代码目录 api/{service-name}/，请手动处理。

确认下线？[是/否]
```

### Step 2: 检查依赖关系

在 `.service-matrix/dependencies.yaml` 中检查是否有其他服务依赖该服务：

- 如果**有上游服务依赖**：⚠️ 警告并列出依赖方，要求用户确认是否继续
- 如果**无依赖**：直接继续

```
⚠️ 以下服务仍在依赖 {service-name}：
  - order（在 downstream 中引用了 {service-name}）

请先处理这些依赖关系，或确认强制下线。
[继续强制下线 / 取消]
```

### Step 3: 删除知识库目录

```
删除 context/project/api/{service-name}/ 整个目录
```

### Step 4: 更新服务矩阵

从 `.service-matrix/dependencies.yaml` 中：
- 移除 `services:` 下该服务的配置块
- 移除 `modules:` 下该模块（如果该模块下无其他服务）
- 移除 `idl_frozen_fields:` 中该服务相关的冻结字段记录
- 清理其他服务 `upstream/downstream` 中对该服务的引用

### Step 5: 更新项目 INDEX.md

从 `context/project/api/INDEX.md` 中：
- 移除服务列表表格中该服务的行
- 更新服务间依赖关系图（移除该服务相关的连线）

### Step 6: 标记错误码段废弃

在 `context/team/error-code.md` 中，将该服务的错误码段标记为 `[已废弃]`，但不删除（保留审计记录）。

### Step 7: 输出完成摘要

```
✅ {service-name} 服务已从 Harness 制品中下线！

🗑️ 已删除：
  - context/project/api/{service-name}/（知识库目录）

📝 已更新：
  - .service-matrix/dependencies.yaml（移除服务配置）
  - context/project/api/INDEX.md（移除服务行和依赖图）
  - context/team/error-code.md（错误码段标记为已废弃）

⚠️ 需要手动处理：
  - 删除业务代码目录 api/{service-name}/（如确认不再需要）
  - 删除 IDL 文件 api/idl/{service-name}/（如确认不再需要）
  - 检查是否有未关闭的需求引用了该服务
```

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `name` | string | ✅ | 要下线的服务名，或 `all` 一键卸载所有服务 |
| `force` | bool | ❌ | 强制下线（跳过依赖检查确认），默认 false |

---

## 示例

```
# 下线单个服务
/service:offboard chat

# 强制下线（跳过依赖检查）
/service:offboard name=chat force=true

# 一键卸载所有服务（重置为空仓库状态）
/service:offboard all
```

---

## `/service:offboard all` 特殊逻辑

当参数为 `all` 时，执行以下操作：

1. 从 `.service-matrix/dependencies.yaml` 读取所有已注册服务
2. 逐个清理每个服务的 Harness 制品
3. 清空 `context/project/api/` 下所有服务知识库目录（保留 INDEX.md）
4. 重置 `.service-matrix/dependencies.yaml` 中 `services:` 和 `modules:` 为空
5. 重置 `context/project/api/INDEX.md` 服务列表和依赖图为空
6. 在 `context/team/error-code.md` 中将所有服务错误码段标记为已废弃

```
⚠️ 即将下线所有服务（共 {n} 个）：
  - order
  - pay
  - user

这将清空所有服务的 Harness 制品，是否确认？[是/否]
```

> **注意**：`/service:offboard all` 不会删除业务代码目录 `api/*/` 和 IDL 文件 `api/idl/*/`，这些需要用户手动处理。
