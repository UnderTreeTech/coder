# Command: /requirement:new

> **说明**：新建需求，创建标准目录结构和需求骨架。
> **触发**：用户输入 `/requirement:new` 或 `/requirement:new {需求标题}`

---

## 执行步骤

### Step 1: 收集基本信息

如果用户没有提供需求标题，询问：
- 需求标题
- 需求 ID（如有 Geelib/TAPD/Jira ID）
- 优先级（P0/P1/P2/P3）

### Step 2: 加载上下文

按顺序加载：
1. `context/team/INDEX.md` - 团队规范
2. `context/harness-framework/main-process-numbering.md` - 流程规范
3. `.service-matrix/dependencies.yaml` - 服务矩阵
4. `context/project/{project-name}/INDEX.md` - 项目知识

### Step 3: 创建目录结构

```
requirements/{requirement-id}/
├── requirement.md          # 从模板创建，预填已知信息
├── status.json             # 初始化为 INIT 阶段
└── tasks/
    └── .gitkeep
```

### Step 4: 预填需求文档

使用 `context/harness-framework/templates/requirement-template.md` 模板，预填：
- 需求 ID
- 创建时间
- 从服务矩阵中提取可能涉及的服务列表

### Step 5: 初始化 status.json

```json
{
    "requirement_id": "{requirement-id}",
    "title": "{需求标题}",
    "current_stage": "REQUIREMENT_DEFINING",
    "current_step": "2.1",
    "gates_passed": [],
    "assignee": "",
    "created_at": "{当前时间}",
    "last_updated": "{当前时间}",
    "branch": "feature/{devops-name}/{tapd-id}"
}
```

### Step 6: 引导用户

输出提示：
```
✅ 需求目录已创建：requirements/{requirement-id}/

📝 下一步：
1. 填写 requirements/{requirement-id}/requirement.md 中的需求内容
2. 完成后运行 /requirement:gate-check 进行需求评审门禁检查
3. 通过后运行 /requirement:next 进入设计阶段

💡 提示：
- 影响面分析请参考 .service-matrix/dependencies.yaml
- 需求模板说明请参考 context/harness-framework/templates/requirement-template.md
```

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `title` | string | ❌ | 需求标题，不填则交互式询问 |
| `id` | string | ❌ | 需求 ID，不填则自动生成 |

---

## 示例

```
/requirement:new
/requirement:new 用户续费功能优化
/requirement:new id=T12345678 用户续费功能优化
```
