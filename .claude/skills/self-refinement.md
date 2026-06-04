# Skill: self-refinement

> **定位**：让 AI 从错误中沉淀经验，实现团队知识的持续积累。
> **触发**：用户纠正 AI 错误后，或 `/knowledge:extract-experience` 命令。

---

## 核心理念

LLM 没有跨会话记忆。但团队的每一个"纠正"，都是一次宝贵的信号。
**错误不再是"走一次算一次"，而是成为团队资产。**

---

## 执行流程

### Step 1: 识别纠正类型

当用户纠正 AI 时，判断：

```
这是"模式性教训"还是"一次性 diff"？

模式性教训：
- 同类错误在不同场景下会重复出现
- 涉及团队规范、业务约束、架构决策
- 新人/新模型也可能犯同样的错误

一次性 diff：
- 特定上下文下的临时修正
- 不具有普遍性
```

### Step 2: 确定沉淀层级

如果是模式性教训，确定沉淀到哪一层：

| 层级 | 条件 | 目标路径 |
|-----|------|---------|
| 团队级 | 所有项目都适用 | `context/team/experience/` |
| 框架工程级 | 所有需求研发都适用 | `context/harness-framework/` |
| 服务级 | 特定服务/模块适用 | `context/project/{project}/{module}/experience/` |

### Step 3: 向用户确认

```
AI 提议：
"这个纠正看起来是一个模式性教训，建议沉淀到 [层级]。
具体内容：[约束描述]
是否确认沉淀？"
```

### Step 4: 生成经验文档

用户确认后，使用 `context/harness-framework/templates/experience-template.md` 模板生成经验文档。

**文件命名**：`{YYYY-MM-DD}-{简短标题}.md`

### Step 5: 更新约束

如果经验涉及代码约束，同时更新：
- 对应服务的 `context/project/*/experience/` 中的 `🔒 约束` 部分
- 如果是通用约束，更新 `context/team/` 中的相关规范

---

## 输出示例

```markdown
## 🧠 Self-Refinement 建议

**识别到模式性教训**：分页查询未设置 LIMIT 上限

**建议沉淀到**：`context/project/layout/internal/dao/experience/`

**约束内容**：
> 所有分页查询必须设置 LIMIT 上限，不得超过 1000 条。
> 超过 1000 条的场景必须使用游标分页（cursor-based pagination）。

**是否确认沉淀？** [是/否/修改内容]
```

---

## 与 AGENTS.md 的关联

AGENTS.md 第 5 条认知模式：
> 当遇到模式性错误时，主动触发 `self-refinement` Skill，提议沉淀经验。
