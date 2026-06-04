# 服务级知识库 — INDEX

> **范围**：特定项目/服务的知识，高频演进，量最大。
> **层级**：服务级（AI 按"团队 → 项目 → 模块 → 服务"路径逐层缩小范围）

---

## 目录结构

```
context/project/
├── INDEX.md                              # 本文件 - 服务知识库入口
└── {project-name}/                       # 逻辑项目（对应一个业务域）
    ├── INDEX.md                          # 项目入口
    └── {module-name}/                    # 业务模块
        ├── INDEX.md                      # 模块入口
        ├── {service-name}/               # 具体服务
        │   ├── INDEX.md                  # 服务入口
        │   ├── architecture.md           # 架构图、技术选型
        │   ├── api-guide.md              # 接口说明（非 IDL 的补充说明）
        │   ├── operations-runbook.md     # 运维手册（部署、监控、告警）
        │   ├── sop/                      # 标准操作规程
        │   │   └── {操作名}.md
        │   └── experience/               # 踩坑经验（按日期命名）
        │       └── {YYYY-MM-DD}-{title}.md
        └── experience/                   # 模块级踩坑经验（跨服务通用）
            └── {YYYY-MM-DD}-{title}.md
```

---

## 示例结构（layout 项目）

```
context/project/
└── layout/
    ├── INDEX.md
    └── internal/
        ├── INDEX.md
        ├── dao/
        │   ├── INDEX.md
        │   ├── architecture.md
        │   └── experience/
        └── service/
            ├── INDEX.md
            ├── architecture.md
            └── experience/
```

---

## AI 使用说明

当处理特定服务的任务时：
1. 先读 `context/project/{project-name}/INDEX.md`
2. 再读对应模块的 `INDEX.md`
3. 最后读服务级的 `architecture.md` 和 `experience/`

**优先级**：服务级经验 > 模块级经验 > 团队级经验
