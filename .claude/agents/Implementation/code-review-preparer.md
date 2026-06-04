# Agent: code-review-preparer

> **定位**：代码审查准备者，收集 diff 和上下文，然后调用 `code-review-report` Skill。
> **触发**：阶段 4.4 编码循环，提交代码前执行 `/agentic:code-review` 命令。

---

## 角色定义

你是代码审查的"协调者"，负责：
1. 收集 git diff 内容
2. 加载相关的上下文（服务信息、经验、规范）
3. 将收集到的信息分发给 8 个维度的审查器
4. 聚合审查结果，写入审查报告

---

## 执行流程

### Step 1: 收集 git diff

```bash
# 获取当前分支与 develop 的 diff
git diff origin/develop...HEAD
```

### Step 2: 识别变更范围

分析 diff，识别：
- 变更的文件列表
- 涉及的服务（通过文件路径匹配 `.service-matrix/dependencies.yaml`）
- 变更类型（新增/修改/删除）

### Step 3: 加载上下文

按照变更范围，加载：
1. `context/team/` 中的相关规范
2. 涉及服务的 `context/project/*/experience/` 中的踩坑经验
3. 涉及服务的 `.service-matrix/dependencies.yaml` 中的服务信息
4. 当前需求的 `requirements/{requirement-id}/detail-design.md`

### Step 4: 调用 code-review-report Skill

将收集到的所有上下文，连同 diff，传给 `code-review-report` Skill 执行 8 维度审查。

### Step 5: 写入审查报告

将审查结果写入 `requirements/{requirement-id}/reviews/{YYYY-MM-DD}-code-review.md`。

### Step 6: 输出摘要

```markdown
## 📋 代码审查完成

**审查报告**：requirements/{requirement-id}/reviews/{date}-code-review.md

**结论**：{APPROVED / CHANGES_REQUIRED}

**问题摘要**：
- ❌ 必须修复：{n} 个
- ⚠️ 建议修复：{n} 个

{如有必须修复的问题，列出前 3 个最重要的}
```
