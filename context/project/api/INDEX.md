# api 项目知识库

> **项目名**：`api`
> **仓库类型**：大仓（Monorepo）
> **业务域**：核心业务 API
> **维护团队**：api team

---

## 服务列表

| 服务 | 协议 | 说明 | 路径 |
|-----|------|------|------|

> **新增服务时**，通过 `/service:onboard` 或 `/service:harness` 命令自动在此表格追加新行。

---

## 大仓架构概览

```
monorepo-layout/                # 仓库根目录
├── CLAUDE.md                   # Harness 全局规范
├── .claude/                    # Claude Code 配置
├── context/                    # 三层知识体系
├── .service-matrix/            # 服务拓扑
├── requirements/               # 需求生命周期产物
├── releases/                   # 版本管理
├── scripts/                    # 工具脚本
│
└── api/                        # 业务代码目录
    ├── idl/                    # 全局 IDL（集中管理 proto 文件）
```

---

## 服务间依赖关系

```
```

**调用链说明**：
- （暂无服务依赖关系）

> **新增服务时**，通过 `/service:onboard` 或 `/service:harness` 命令自动更新此依赖关系图。

---

## 关键约束（AI 必读）

1. 各服务目录下**不包含**任何 Harness 相关文件（`CLAUDE.md`、`.claude/`、`context/` 等）
2. 所有需求/设计/经验文档统一在**根目录的**对应目录下管理
3. 跨服务调用必须通过 HTTP/gRPC 接口，禁止直接跨 module 引用
4. 各服务的 `internal/` 目录下的包**禁止**被其他服务直接引用
5. 数据库**禁止**跨服务共享，每个服务有独立的数据库
