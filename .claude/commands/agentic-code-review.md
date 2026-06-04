# Command: /agentic:code-review

> **说明**：触发多维度代码审查，收集 diff 和上下文，生成审查报告。
> **触发**：用户输入 `/agentic:code-review`

---

## 执行步骤

### Step 1: 确认当前需求

读取 `requirements/` 目录，找到当前处于 CODING 阶段的需求。
如果有多个，询问用户选择哪个。

### Step 2: 委派给 code-review-preparer Agent

调用 `.claude/agents/Implementation/code-review-preparer.md` 执行：
1. 收集 git diff
2. 加载上下文
3. 执行 8 维度审查
4. 生成审查报告

### Step 3: 输出结果

显示审查摘要，并提示审查报告的路径。

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `req` | string | ❌ | 需求 ID，不填则自动检测当前需求 |

---

## 示例

```
/agentic:code-review
/agentic:code-review req=T12345678
```
