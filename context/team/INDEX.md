# 团队级知识库 — INDEX

> **范围**：所有项目必须遵循的团队规范。
> **层级**：团队级（最高优先级，不可被项目级覆盖）

---

## 目录结构

```
context/team/
├── INDEX.md                    # 本文件 - 团队知识库入口
├── git-convention.md           # Git 提交规范、分支策略
├── error-code.md               # 错误码空间分配规范
├── logging.md                  # 日志规范（级别、字段、格式）
├── security.md                 # 安全规范（鉴权、加密、敏感数据）
├── protobuf-style-guide.md     # pb定义规范（Go 语言）
├── testing.md                  # 测试规范（单测、集成测试覆盖率要求）
├── development-sop.md          # 服务功能开发标准 SOP
└── experience/                 # 团队级踩坑经验（跨项目通用）
    └── .gitkeep
```

---

## 快速导航

| 文件 | 内容 | 适用场景 |
|-----|------|---------|
| [`protobuf-style-guide.md`](protobuf-style-guide.md) | gRPC 和 HTTP 接口 Protobuf 定义规范（包名、注释强制要求等） | 设计 API、新增或修改 proto 接口文件 |
| [`code-implementation-examples.md`](code-implementation-examples.md) | 涵盖 gRPC, HTTP, DAO, Model 的标准代码实现样板 | 新项目初始化、AI 代码生成参考 |
| [`git-convention.md`](git-convention.md) | 分支命名、提交信息格式、PR 规范 | 创建分支、提交代码、发起 CR |
| [`error-code.md`](error-code.md) | 错误码空间分配、错误码格式 | 新增错误码、错误处理 |
| [`logging.md`](logging.md) | 日志级别定义、必填字段、禁止字段 | 添加日志、日志排查 |
| [`security.md`](security.md) | 鉴权方式、敏感数据处理、加密规范 | 接口鉴权、数据加密 |
| [`code-style.md`](code-style.md) | Go 代码风格、命名规范、注释规范 | 代码审查、新功能开发 |
| [`testing.md`](testing.md) | 单测覆盖率要求、测试命名规范 | 编写测试、CI 门禁 |
| [`development-sop.md`](development-sop.md) | 服务功能开发标准 SOP（分层、HTTP/gRPC） | 需求开发、vibe coding |

---

## AI 使用说明

当需要了解团队规范时，按以下顺序读取：
1. 读本文件了解整体结构
2. 根据任务类型，读对应的规范文件
3. 如果规范文件中有 `experience/` 引用，读取对应经验文档

**不要猜测团队规范，始终以本目录为准。**
