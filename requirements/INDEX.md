# 需求生命周期产物目录 — INDEX

> **说明**：每个需求在此目录下有独立的子目录，包含从需求定义到交付的全部产物。
> 需求、设计、任务、门禁结论和代码形成完整追溯链。

---

## 目录结构

```
requirements/
├── INDEX.md                              # 本文件
├── {requirement-id}/                     # 单个需求目录
│   ├── requirement.md                    # 需求文档（来自模板）
│   ├── status.json                       # 当前阶段状态
│   ├── outline-design.md                 # 概要设计
│   ├── detail-design.md                  # 详细设计
│   ├── tasks/
│   │   └── features.json                 # 任务拆分（Dev门禁检查）
│   ├── gate-1-requirement-review.md      # 需求评审门禁结论
│   ├── gate-2-design-review.md           # 设计门禁结论
│   ├── gate-3-dev-entry.md               # Dev进入门禁结论
│   ├── gate-4-service-repo-check.md      # 服务仓库检查门禁结论
│   └── reviews/                          # 代码审查报告
│       └── {YYYY-MM-DD}-code-review.md
└── _template/                            # 需求目录模板
    └── status.json
```

---

## 需求状态说明

`status.json` 记录当前需求所处阶段：

```json
{
    "requirement_id": "{requirement-id}",
    "title": "{需求标题}",
    "current_stage": "INIT",
    "current_step": "1.1",
    "gates_passed": [],
    "assignee": "{@姓名}",
    "created_at": "{YYYY-MM-DDTHH:mm:ss+08:00}",
    "last_updated": "{YYYY-MM-DDTHH:mm:ss+08:00}",
    "branch": "feature/{devops-name}/{tapd-id}"
}
```

**阶段枚举**：
- `INIT` - 初始化
- `REQUIREMENT_DEFINING` - 需求定义中
- `DESIGNING` - 设计中
- `DEV_PREPARING` - 开发准备
- `CODING` - 编码中
- `DELIVERING` - 交付中
- `DONE` - 完成

---

## 快速命令

```bash
# 新建需求（自动创建目录结构）
/requirement:new

# 恢复上下文（读取 status.json 恢复到上次状态）
/requirement:continue

# 进入下一阶段
/requirement:next

# 门禁自检
/requirement:review
```
