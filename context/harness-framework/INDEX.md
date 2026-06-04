# Harness Framework 工程规范 — INDEX

> **范围**：所有需求研发必须遵循的框架工程规范。
> **层级**：框架工程级（中频更新，被所有需求研发继承）

---

## 目录结构

```
context/harness-framework/
├── INDEX.md                          # 本文件 - 框架知识库入口
├── main-process-numbering.md         # 五阶段流程 + 四门禁（唯一真相源）
├── gates/                            # 门禁规则详细定义
│   ├── requirement-quality-gate.md   # 需求评审门禁规则
│   ├── design-quality-gate.md        # 设计门禁规则
│   ├── dev-entry-gate.md             # Dev 进入门禁规则
│   └── service-repo-check.md         # 服务仓库检查门禁规则
├── templates/                        # 文档模板
│   ├── requirement-template.md       # 需求文档模板
│   ├── outline-design-template.md    # 概要设计模板
│   ├── detail-design-template.md     # 详细设计模板
│   └── experience-template.md        # 踩坑经验模板
└── context-collection-rules.md       # 上下文收集规范
```

---

## 快速导航

| 文件 | 内容 | 适用场景 |
|-----|------|---------|
| [`main-process-numbering.md`](main-process-numbering.md) | 五阶段流程、四门禁定义 | 了解整体流程、门禁自检 |
| [`gates/requirement-quality-gate.md`](gates/requirement-quality-gate.md) | 需求评审门禁检查项 | 阶段 2.2 门禁 |
| [`gates/design-quality-gate.md`](gates/design-quality-gate.md) | 设计门禁检查项 | 阶段 3.3 门禁 |
| [`gates/dev-entry-gate.md`](gates/dev-entry-gate.md) | Dev 进入门禁检查项 | 阶段 4.2 门禁 |
| [`gates/service-repo-check.md`](gates/service-repo-check.md) | 服务仓库检查门禁 | 阶段 4.3 门禁 |
| [`templates/requirement-template.md`](templates/requirement-template.md) | 需求文档模板 | 新建需求 |
| [`templates/detail-design-template.md`](templates/detail-design-template.md) | 详细设计模板 | 设计阶段 |
| [`templates/experience-template.md`](templates/experience-template.md) | 踩坑经验模板 | 经验沉淀 |

---

## AI 使用说明

1. 开始任何需求研发前，先读 `main-process-numbering.md` 了解当前所处阶段
2. 到达门禁节点时，读对应的 `gates/` 文件执行检查
3. 创建文档时，使用 `templates/` 中的模板
4. 遇到问题时，先查 `context/project/` 中的服务级经验
