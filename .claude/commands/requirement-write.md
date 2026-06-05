# 需求撰写执行命令

> **说明**：基于外部产品文档内容，按照标准化需求模板，由 AI 自动拆解并填写需求。
> **前置条件**：已通过 `/requirement:new` 创建了需求目录和初始状态。
> **触发**：用户输入 `/requirement:write`

---

## 执行步骤

### Step 1: 确定当前需求 ID 和文档路径
1. 加载当前活动的 requirement id（可从 `.harness/local.yaml` 或询问用户获得）。
2. 定位到对应的 `requirements/{requirement-id}/requirement.md`。如果该文件不存在，提示用户先执行 `/requirement:new`。

### Step 2: 收集产品需求输入
询问用户提供包含产品需求文档内容的 MCP 服务或网页 URL。若用户直接以文本方式输入部分产品设计内容，也予以接受。
例如：“请输入能读取产品 PRD 的 MCP 服务指令（例如 `mcp fetch_prd doc_id`）或是需求内容的详细描述”

### Step 3: 读取需求模板
读取框架级模板 `context/harness-framework/templates/requirement-template.md` 的规范结构。

### Step 4: AI 拆解与组装
根据用户输入的产品文档和读取到的模板，执行如下拆分与补齐动作：
1. **背景和目标**：从产品文档中提炼为什么要做，以及可度量的目标。
2. **非目标**：识别当前范围不包含的内容。
3. **用户故事**：将业务描述转化为“作为...希望...以便...”的格式。
4. **功能需求条目**：将产品模块拆分成原子级的需求条目（REQ-xxx），并提取/生成具体可测试的验收标准。
5. **非功能需求**：识别性能、可用性、安全等要求。
6. **影响面分析**：初步提取可能涉及的模块/服务，并留空待开发者进一步补充确认。

### Step 5: 更新需求文档
将生成并组装好的 Markdown 内容写回 `requirements/{requirement-id}/requirement.md` 文件。

### Step 6: 引导用户下一步
输出提示：
```
✅ 需求文档已生成并写入：requirements/{requirement-id}/requirement.md

📝 下一步：
1. 请人工走查并微调生成的需求文档，特别是"影响面分析"和"验收标准"部分。
2. 确认无误后，运行 /requirement:review 进行需求评审门禁检查。
```

---

## 示例

```
/requirement:write
```