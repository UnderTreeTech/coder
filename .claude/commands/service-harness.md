# Command: /service:harness

> **说明**：仅生成 Harness 工程制品（知识库目录 + 服务矩阵 + INDEX 更新），
> 适用于**用户已手动执行 `waterdrop new {service_name}` 后**，再调用本命令补全 Harness 接入。
>
> **触发**：用户输入 `/service:harness {service-name}` 或 `/service:harness`

---

## 与 /service:onboard 的区别

| 命令 | waterdrop new | Harness 制品 |
|-----|--------------|-------------|
| `/service:onboard` | ✅ 自动执行 | ✅ 自动生成 |
| `/service:harness` | ❌ 用户自己执行 | ✅ 自动生成 |

**使用场景**：
- 已经手动执行过 `waterdrop new chat`
- 或者业务代码目录是从其他方式创建的（非 waterdrop）
- 只需要补全 Harness 工程制品

---

## 执行步骤

### Step 1: 收集服务基本信息

交互式询问（同 `/service:onboard`）：
- 服务名、描述、通信协议（HTTP/gRPC）
- 上游/下游服务
- 错误码段起始值

### Step 2: 验证业务代码目录存在

```
检查 api/{service-name}/ 目录是否存在...
```

- 目录**存在**：继续
- 目录**不存在**：提示用户先执行 `cd api && waterdrop new {service_name}`

### Step 3: 生成知识库目录

> ⚠️ **关键路径约束**：知识库目录必须创建在 `context/project/api/{service-name}/` 下，
> **绝对不能**创建在 `api/{service-name}/` 业务代码目录下。
> 业务代码目录 `api/{service-name}/` 中只放代码，不放任何 Harness 制品。

在 **`context/project/api/{service-name}/`** 下创建：

```
context/project/api/{service-name}/
├── INDEX.md
├── architecture.md
├── api-guide.md
├── operations-runbook.md
├── sop/.gitkeep
└── experience/.gitkeep
```

### Step 4: 更新服务矩阵

在 `.service-matrix/dependencies.yaml` 追加服务配置。

### Step 5: 更新项目 INDEX.md

在 `context/project/api/INDEX.md` 的服务列表表格中追加新服务行，并更新服务间依赖关系图。

### Step 6: 预留错误码段

在 `context/team/error-code.md` 中预留码段。

### Step 7: 输出完成摘要

```
✅ {service-name} 服务 Harness 制品已生成！

📁 新建文件：
  - context/project/api/{service-name}/INDEX.md
  - context/project/api/{service-name}/architecture.md
  - context/project/api/{service-name}/api-guide.md
  - context/project/api/{service-name}/operations-runbook.md

📝 更新文件：
  - .service-matrix/dependencies.yaml
  - context/project/api/INDEX.md
  - context/team/error-code.md

🔜 下一步：
  1. 在 api/idl/{service-name}/v1/ 下定义 proto 文件
  2. 在 context/project/api/{service-name}/INDEX.md 中补充关键约束
  3. 运行 /service:deps 验证服务依赖关系
```

---

## 示例

```
# 用户先手动执行 waterdrop
cd api && waterdrop new chat

# 然后告诉 Claude 生成 Harness 制品
/service:harness chat
```
