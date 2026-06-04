# Command: /knowledge:extract-experience

> **说明**：从本次开发过程中提取踩坑经验，沉淀为团队资产。
> **触发**：用户输入 `/knowledge:extract-experience`，通常在需求交付阶段或遇到问题修复后。

---

## 执行步骤

### Step 1: 收集本次开发信息

读取：
- 当前需求的代码审查报告（`requirements/*/reviews/`）
- 本次开发过程中修复的问题
- 用户在对话中提到的"踩坑"或"注意事项"

### Step 2: 调用 self-refinement Skill

委派给 `.claude/skills/self-refinement.md` 执行：
1. 识别哪些是模式性教训
2. 确定沉淀层级
3. 向用户确认

### Step 3: 生成经验文档

使用 `context/harness-framework/templates/experience-template.md` 生成经验文档。

### Step 4: 更新相关约束

如果经验涉及代码约束，同步更新对应目录下的约束列表。

---

## 示例

```
/knowledge:extract-experience
```

输出示例：
```
🧠 从本次开发中识别到 2 个模式性教训：

1. 分页查询未设置 LIMIT 上限 → 建议沉淀到 context/project/layout/internal/dao/experience/
2. goroutine 泄漏处理 → 建议沉淀到 context/team/experience/

是否确认沉淀？[全部确认 / 逐个确认 / 跳过]
```
