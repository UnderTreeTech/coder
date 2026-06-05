# Skill: managing-requirement-lifecycle

> **定位**：整个框架的"中央调度"。需求工作通过它推进，以保证阶段、门禁和上下文口径一致。
> **职责**：意图识别、阶段检查、门禁验证、债务检查和计划更新。

---

## 触发条件

- 用户执行任何 `/requirement:*`, `/design:*`, `/coding:*` 命令时
- Agent 需要推进需求阶段时
- 门禁检查前后

---

## 执行流程

### Step 1: 加载当前状态

```
读取 requirements/{requirement-id}/status.json
确认当前阶段（current_stage）和步骤（current_step）
```

### Step 2: 意图识别

根据用户输入，识别以下意图之一：
- `NEW_REQ`：新建需求 (`/requirement:new`)
- `REQ_GATE`：需求门禁检查 (`/requirement:review`)
- `NEW_DESIGN`：新建设计 (`/design:new`)
- `DESIGN_GATE`：设计门禁检查 (`/design:review`)
- `START_CODING`：开始编码 (`/coding:start`)
- `CODE_REVIEW`：代码审查门禁 (`/coding:review`)
- `STATUS`：查看当前状态

### Step 3: 阶段合法性检查

```
检查当前意图是否符合阶段状态机规则：
- 不允许跳过门禁直接进入下一阶段
- 不允许在没有通过前置门禁的情况下开始编码
```

### Step 4: 执行对应操作

根据意图和当前阶段，调用对应的 Skill 或 Agent：

| 阶段 | 意图 | 调用 |
|-----|------|------|
| INIT | NEW_REQ | `requirement-bootstrapper` Agent |
| REQUIREMENT_DEFINING | REQ_GATE | `requirement-quality-reviewer` Agent |
| DESIGNING | NEW_DESIGN | `design-doc-writer` Skill |
| DESIGNING | DESIGN_GATE | `detail-design-quality-reviewer` Agent |
| CODING | START_CODING | `coding-session-initializer` Skill |
| CODING | CODE_REVIEW | `code-review-report` Skill |

### Step 5: 更新状态

操作完成后，更新 `requirements/{requirement-id}/status.json` 中的阶段信息。

---

## 状态机规则

```
INIT
  ↓ NEW_REQ
REQUIREMENT_DEFINING
  ↓ [GATE_1 需求门禁通过]
DESIGNING
  ↓ [GATE_2 设计门禁通过]
CODING
  ↓ [GATE_3 代码审查门禁通过]
DONE
```

**任何跳步都被禁止。**

---

## 输出格式

每次执行完成后，输出标准状态报告：

```markdown
## 📊 需求状态

**需求**：{title}（{requirement-id}）
**当前阶段**：{current_stage} - {current_step}
**已通过门禁**：{gates_passed}

**下一步**：{next_action}
```
