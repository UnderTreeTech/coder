# CLAUDE.md — Harness Engineering 全局协作规范（Claude Code）

> **核心理念**：AI 参与问题分析、方案设计、编码实现、审查和验证，但**最终判断权始终留在工程师手中**。
>
> 代码产出 = AI 能力 × 上下文质量 —— 上下文质量的提升，完全掌握在团队自己手中。

---

## 1. 大仓（Monorepo）结构说明

本仓库是一个**大仓（Monorepo）**，包含多个微服务：

```
monorepo-layout/                    # 仓库根目录
├── CLAUDE.md                       # Harness 全局规范（本文件）
├── .claude/                        # Claude Code 配置（commands/skills/agents）
├── .service-matrix/                # 服务拓扑（单一真相源）
├── context/                        # 三层知识体系
├── requirements/                   # 需求生命周期产物
├── releases/                       # 版本管理
├── scripts/                        # 工具脚本
│
└── api/                            # 业务代码（按服务分目录）
    ├── order/                      # 订单服务
    │   ├── api/                    # Protobuf IDL
    │   ├── cmd/                    # 程序入口
    │   ├── configs/                # 配置文件
    │   └── internal/               # 业务逻辑
    ├── pay/                        # 支付服务
    │   ├── api/
    │   ├── cmd/
    │   ├── configs/
    │   └── internal/
    └── user/                       # 用户服务
        ├── api/
        ├── cmd/
        ├── configs/
        └── internal/
```

**关键约束**：
- Harness 工程制品（`CLAUDE.md`、`.claude/`、`context/`、`requirements/`、`releases/`）**只在根目录**，不在各服务目录中重复
- 各服务目录（`api/order/`、`api/pay/`、`api/user/`）只包含业务代码，不含 Harness 相关内容
- 所有需求、设计、门禁结论统一在根目录的 `requirements/` 下管理

---

## 2. 认知模式（必读）

### 2.1 Harness Engineering 的本质

Harness Engineering 不是让 AI "看起来更聪明"，而是让 AI 在真实业务系统里"长期更可靠"。

- **vibe coding 的底层逻辑**：让 AI 尽量自由地生成
- **Harness Engineering 的底层逻辑**：让 AI 在正确的轨道上尽量高效地生成

### 2.2 五大上下文缺口（AI 的盲区）

| 缺口类型 | 典型问题 | 解决方案 |
|---------|---------|---------|
| 隐性规范 | 团队约定的锁机制、埋点规则、错误码空间 | `context/team/` |
| 历史决策 | "为什么当时选了 A 方案不选 B" | `context/project/{service}/experience/` |
| 服务契约 | IDL 字段的冻结状态、下游是否强依赖 | `.service-matrix/dependencies.yaml` |
| 跨服务依赖 | 同一个需求要改哪几个服务、谁调谁 | 服务矩阵自动解析 |
| 演进轨迹 | 某个模块上次大改的坑、灰度策略 | Self-Refinement + `experience/*.md` |

---

## 3. 三阶段 + 三门禁流程

### 3.1 三阶段主流程

```
阶段 1: 需求定义 ⭐ → 撰写需求 → 评审（需求门禁）
阶段 2: 拆解设计 ⭐ → 方案与任务拆解 → 评审（设计门禁）
阶段 3: 编码实现 ⭐ → Vibe Coding 循环 → 审查（代码审查门禁）
```

⭐ = 强制门禁（共 3 个，不可跳过）

**核心理念：错误越早被拦住，代价越低。**

### 3.2 三道门禁

门禁口径收拢在 [`context/harness-framework/main-process-numbering.md`](context/harness-framework/main-process-numbering.md)。

| 门禁 | 触发指令 | 阻塞条件概览 |
|-----|---------|-------------|
| 需求门禁 | `/requirement:review` | 背景不清晰 / 验收标准缺失 / 影响面未分析 |
| 设计门禁 | `/design:review` | 漏了约束 / 任务未拆解 / `tasks/features.json` 不合法 |
| 代码审查门禁 | `/coding:review` | 不符合规范 / 逻辑不符合设计 / 单测未覆盖 |

---

## 4. 三层知识体系

| 层级 | 位置 | 范围 | 典型内容 |
|-----|------|------|---------|
| 团队级 | `context/team/` | 所有服务必须遵循 | Git 规范、错误码空间、日志规范 |
| 框架工程级 | `context/harness-framework/` | 所有需求研发必须遵循 | 三阶段流程、门禁规则、文档模板 |
| 服务级 | `context/project/api/{service}/` | 特定服务 | 架构图、API、踩坑经验 |

---

## 5. 大仓分支策略

每个需求，在**同一个仓库**中使用统一的分支名：

```
feature/{devops-name}
```

由于是大仓，所有服务的代码变更都在同一个分支上，通过目录路径区分服务。

---

## 6. 占位符词典

| 占位符 | 语义 | 举例 |
|-------|------|------|
| `{repo-root}` | 仓库根目录的磁盘路径（绝对） | `/data/workspace/api` |
| `{project-name}` | 逻辑项目名 | `api` |
| `{service-name}` | 服务名 | `order` / `pay` / `user` |
| `{requirement-id}` | 需求 ID | `T12345678` |

---

## 7. Self-Refinement 闭环

LLM 没有跨会话记忆。但团队的每一个"纠正"，都是一次宝贵的信号。

**闭环流程**：
1. 用户纠正 AI 某个错误
2. AI 识别：这是"模式性教训"还是"一次性 diff"？
3. AI 主动提议沉淀层级（团队级 / 框架级 / 服务级）
4. 用户确认 → 生成 experience 文档 / 更新规范
5. 下次同类场景，AI 主动引用

---

## 8. 工程制品清单

| 工程制品 | 作用 |
|---------|------|
| `CLAUDE.md` | 全局协作规范（本文件，Claude Code 自动读取） |
| `.claude/commands/` | Slash Commands（`/requirement:new` 等） |
| `.claude/skills/` | 可复用工作流规范 |
| `.claude/agents/` | 专家角色（门禁执行者） |
| `context/team/` | 团队级规范 |
| `context/harness-framework/` | 框架工程规范 |
| `context/project/api/{service}/` | 服务级知识 |
| `.service-matrix/dependencies.yaml` | 服务拓扑（单一真相源） |
| `requirements/` | 需求生命周期产物 |
| `releases/` | 版本管理（SQL/配置/脚本变更） |

---

## 9. 硬规则（不可违反）

1. **不跳过门禁**：三道门禁是强制的
2. **不硬编码路径**：所有路径使用占位符
3. **不口头通过**：门禁结论必须写入文件
4. **不散落聊天**：需求、设计、经验必须沉淀到对应目录
5. **知识及时沉淀**：每次踩坑后，必须更新 `experience/` 目录

---

## 10. Claude Code 使用说明

### 10.1 可用命令一览

所有命令定义在 `.claude/commands/` 目录下，Claude Code 自动识别。

#### 需求管理命令

| 命令 | 说明 | 示例 |
|-----|------|------|
| `/requirement:new` | 新建需求，创建标准目录骨架和需求文档 | `/requirement:new 用户续费功能优化` |
| `/requirement:write` | 使用外部产品文档 MCP 填充拆解需求 | `/requirement:write` |
| `/requirement:review` | 触发需求阶段门禁自检，生成检查结论 | `/requirement:review` |

#### 设计管理命令

| 命令 | 说明 | 示例 |
|-----|------|------|
| `/design:new` | 新建设计方案 | `/design:new` |
| `/design:review` | 触发设计阶段门禁检查 | `/design:review` |

#### 服务管理命令

| 命令 | 说明 | 示例 |
|-----|------|------|
| `/service:onboard` | **一条龙**：自动执行 `waterdrop new` + 生成 Harness 制品 | `/service:onboard chat` |
| `/service:harness` | **仅 Harness 制品**：用户已手动执行 waterdrop 后调用 | `/service:harness chat` |
| `/service:offboard` | **下线服务**：清理 Harness 制品（知识库/服务矩阵/依赖图） | `/service:offboard chat` |
| `/service:deps` | 查看服务依赖关系和影响面 | `/service:deps order` |

#### 编码与审查命令

| 命令 | 说明 | 示例 |
|-----|------|------|
| `/coding:start` | 开始编码循环 | `/coding:start` |
| `/coding:review` | 触发代码审查门禁 | `/coding:review` |
| `/agentic:code-review` | 触发 8 维度并行代码审查 | `/agentic:code-review` |

#### 知识沉淀命令

| 命令 | 说明 | 示例 |
|-----|------|------|
| `/knowledge:extract-experience` | 从本次开发中提取踩坑经验 | `/knowledge:extract-experience` |

---

### 10.2 新项目初始化流程

```bash
# 1. 填写服务矩阵（必须）
vim .service-matrix/dependencies.yaml

# 2. 运行初始化验证
bash scripts/install.sh

# 3. 在 Claude Code 中开始第一个需求
/requirement:new
```

---

### 10.3 新增服务流程

#### 方式一：一条龙（推荐）

```
/service:onboard chat
```

Claude Code 自动执行 `waterdrop new chat` + 生成全部 Harness 制品。

#### 方式二：分步执行

```bash
# 用户先手动执行 waterdrop
cd api && waterdrop new chat

# 然后告诉 Claude 生成 Harness 制品
/service:harness chat
```

两种方式最终效果相同，区别只在于 `waterdrop new` 是 Claude 执行还是用户手动执行。

---

### 10.4 日常需求研发流程

```
# 1. 新建需求
/requirement:new 某某功能需求

# 2. 撰写需求文档（执行 /requirement:write 或手动补全）
# 生成并编辑 requirements/{req-id}/requirement.md

# 3. 需求门禁
/requirement:review

# 4. 新建设计方案
/design:new

# 5. 撰写设计方案和任务拆解
# 生成 tasks/features.json 等

# 6. 设计门禁
/design:review

# 7. 开始编码循环
/coding:start

# 8. 编码完成后代码审查门禁
/coding:review

# 9. 交付后沉淀经验
/knowledge:extract-experience
```

---

### 10.5 版本发布流程

```bash
# 1. 创建版本目录
bash scripts/new-release.sh v1.0.0

# 2. 添加变更内容
# 编辑 releases/v1.0.0/sql/ddl/   （DDL 变更）
# 编辑 releases/v1.0.0/sql/dml/   （DML 变更）
# 编辑 releases/v1.0.0/configs/   （配置变更）
# 编辑 releases/v1.0.0/scripts/   （运维脚本）

# 3. 填写发布说明和检查清单
# 编辑 releases/v1.0.0/RELEASE_NOTES.md
# 编辑 releases/v1.0.0/checklist.md

# 4. 运维按 checklist.md 执行发布
```

---

> **最后一句话**：Context Engineering + Spec-First + Knowledge as Code，构成了可验证、可演进的 AI 协作工程基线。
>
> 工程规范驱动 Claude Code，让 AI 在正确的轨道上高效生成，而不是让 AI 成为"绕过约束的捷径"。
