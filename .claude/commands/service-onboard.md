# Command: /service:onboard

> **说明**：新增服务时，自动生成 Harness 工程制品（知识库目录、服务矩阵登记、IDL 骨架），
> **业务代码由 `waterdrop new {service_name}` 生成，不在本命令范围内**。
>
> **触发**：用户输入 `/service:onboard {service-name}` 或 `/service:onboard`

---

## 职责边界

本命令是**一条龙**自动化命令，用户只需提供服务名，自动完成：

| 步骤 | 操作 | 方式 |
|-----|------|------|
| ① 生成业务代码 | `cd api && waterdrop new {service_name}` | 自动执行 shell 命令 |
| ② 生成 Harness 制品 | 知识库目录 + 服务矩阵 + INDEX 更新 | Claude Code 自动创建文件 |

> **注意**：IDL（proto 文件）由用户后续手动定义，本命令不关心。

---

## 执行步骤

### Step 1: 收集服务基本信息

如果用户没有提供，交互式询问：

```
请提供新服务的基本信息：

1. 服务名（英文小写，如 chat）：
2. 服务描述（中文，如 即时通讯服务）：
3. 通信协议：
   [1] HTTP（对外暴露，使用 grpc-gateway）
   [2] gRPC（内部服务）
4. 所属项目（如 api）：
5. 上游服务（哪些服务会调用本服务，逗号分隔，如 order,user）：
6. 下游服务（本服务会调用哪些服务，逗号分隔，如 user）：
7. 错误码段起始值（参考 context/team/error-code.md，如 5000）：
```

### Step 2: 自动执行 waterdrop new 生成业务代码

**在 `api/` 目录下执行 `waterdrop new {service_name}`**：

```bash
cd api && waterdrop new {service-name}
```

- 如果执行**成功**：继续下一步
- 如果执行**失败**：输出错误信息，提示用户检查 waterdrop 是否已安装，终止流程

```
⚠️ waterdrop new chat 执行失败。
请确认 waterdrop 已正确安装：go install github.com/UnderTreeTech/waterdrop/cmd/waterdrop@latest
```

### Step 3: 生成服务知识库目录

> ⚠️ **关键路径约束**：知识库目录必须创建在 `context/project/api/{service-name}/` 下，
> **绝对不能**创建在 `api/{service-name}/` 业务代码目录下。
> 业务代码目录 `api/{service-name}/` 中只放代码，不放任何 Harness 制品。

在 **`context/project/api/{service-name}/`** 下创建以下文件：

```
context/project/api/chat/
├── INDEX.md                    # 从 service-knowledge-template.md 生成
├── architecture.md             # 架构图（空模板，待填写）
├── api-guide.md                # 接口补充说明（非 IDL 部分）
├── operations-runbook.md       # 运维手册（空模板，待填写）
├── sop/
│   └── .gitkeep
└── experience/
    └── .gitkeep
```

`INDEX.md` 从 `context/harness-framework/templates/service-knowledge-template.md` 模板生成，预填：
- 服务名、代码路径（指向 `api/{service-name}/`）、通信协议、业务域描述
- 错误码段
- 上下游依赖关系

### Step 4: 更新服务矩阵

在 `.service-matrix/dependencies.yaml` 的 `services:` 下**追加**新服务配置：

```yaml
  # ===== {ServiceName} 服务（{protocol}）=====
  {service-name}:
    module: {project-name}
    path: "{repo-root}/api/{service-name}"
    entrypoint: "cmd/main.go"
    protocol: "{http|grpc}"
    idl_path: "{repo-root}/api/idl/{service-name}/v1/{service-name}.proto"
    description: "{描述}"
    dependencies: []
    upstream: [{上游服务列表}]
    downstream: [{下游服务列表}]
```

同时在 `modules:` 下追加模块配置（如果所属模块不存在）。

### Step 5: 更新项目 INDEX.md

在 `context/project/api/INDEX.md` 的**服务列表表格**中追加新服务行：

```markdown
| 服务 | 说明 | 路径 |
|-----|------|------|
| `chat` | 即时通讯服务 | [`chat/INDEX.md`](chat/INDEX.md) |   
```

同时更新该文件中的**服务间依赖关系图**，将新服务的调用关系补充进去。

**注意**：此文件是 AI 按"团队 → 项目 → 模块"路径定位服务的入口，新服务不在此登记则 AI 无法感知其存在。

### Step 7: 在 error-code.md 中预留错误码段

在 `context/team/error-code.md` 的错误码登记表中，为新服务预留码段。

### Step 8: 输出完成摘要

```
✅ chat 服务 Harness 制品已生成！

📁 新建文件：
  - context/project/api/chat/INDEX.md
  - context/project/api/chat/architecture.md
  - context/project/api/chat/api-guide.md
  - context/project/api/chat/operations-runbook.md

📝 更新文件：
  - .service-matrix/dependencies.yaml（新增 chat 服务）
  - context/project/api/INDEX.md（新增 chat 服务条目）
  - context/team/error-code.md（预留 {起始码}-{结束码} 给 chat）

🔜 下一步（需要手动完成）：
  1. 在 api/idl/chat/v1/ 下定义 proto 文件并生成代码
  2. 在 context/project/api/chat/INDEX.md 中补充关键约束
  3. 如果有跨服务调用，更新 .service-matrix/dependencies.yaml 中的 upstream/downstream
  4. 运行 /service:deps 验证服务依赖关系
```

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `name` | string | ❌ | 服务名，不填则交互式询问 |
| `protocol` | string | ❌ | `http` 或 `grpc`，不填则询问 |
| `project` | string | ❌ | 所属项目，默认 `api` |

---

## 示例

```
# 交互式（推荐）
/service:onboard

# 直接指定服务名
/service:onboard chat

# 完整参数
/service:onboard name=chat protocol=grpc project=api
```

---

## 完整新增服务流程

只需一条命令：

```
/service:onboard chat
```

Claude Code 将自动完成：
1. 交互式询问协议/依赖等信息
2. 执行 `cd api && waterdrop new chat` 生成业务代码
3. 创建 `context/project/api/chat/` 知识库目录
4. 更新 `.service-matrix/dependencies.yaml`
5. 更新 `context/project/api/INDEX.md`
6. 预留错误码段

之后用户只需手动完成：
- 在 `api/idl/chat/v1/` 下定义 proto 文件
- 在 `context/project/api/chat/INDEX.md` 中补充关键约束
