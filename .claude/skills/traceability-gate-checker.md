# Skill: traceability-gate-checker

> **定位**：追溯链校验器，确保需求条目 → 设计决策 → 开发任务 → 代码 diff 的追溯链完整。
> **触发**：设计门禁检查时，或代码审查时。

---

## 追溯链定义

```
需求条目（REQ-XXX）
    ↓
设计决策（detail-design.md 中的对应章节）
    ↓
开发任务（tasks/features.json 中的 task）
    ↓
代码变更（git commit message 中的 REQ-XXX 引用）
```

---

## 执行流程

### Step 1: 读取需求条目

从 `requirements/{requirement-id}/requirement.md` 中提取所有 `REQ-XXX` 条目。

### Step 2: 检查设计覆盖

对每个 REQ-XXX，在 `requirements/{requirement-id}/detail-design.md` 中查找对应的设计决策。

### Step 3: 检查任务覆盖

对每个 REQ-XXX，在 `requirements/{requirement-id}/tasks/features.json` 中查找对应的任务。

### Step 4: 检查代码覆盖（可选）

如果已有代码提交，检查 git log 中是否有引用 REQ-XXX 的 commit。

### Step 5: 生成追溯链报告

---

## 输出格式

```markdown
## 追溯链检查报告

| 需求条目 | 设计覆盖 | 任务覆盖 | 代码覆盖 |
|---------|---------|---------|---------|
| REQ-001 | ✅ | ✅ | ✅ |
| REQ-002 | ✅ | ❌ | - |
| REQ-003 | ❌ | - | - |

**结论**：{COMPLETE / INCOMPLETE}

**缺失项**：
- REQ-002：缺少对应的开发任务
- REQ-003：缺少设计决策
```
