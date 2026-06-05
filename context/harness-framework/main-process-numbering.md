# 五阶段流程 + 三门禁 — 唯一真相源

> ⚠️ **本文件是整条流程语义的唯一真相源。**
> AGENTS.md、每一个 Skill、每一个 Command 都围绕它保持一致。
> 一次更新，全仓生效，避免规范口径散落多处并逐渐漂移。

---

## 流程总览（简化版）

```
阶段 1: 需求定义 ⭐ [需求门禁]
    └── 1.1 需求录入（/requirement:new）
    └── 1.2 需求撰写（AI 补齐背景、目标、验收标准）
    └── 1.3 ⭐ 需求门禁（/requirement:gate-check）

         ↓

阶段 2: 拆解设计 ⭐ [设计门禁]
    └── 2.1 方案设计（/design:new）
    └── 2.2 任务拆解（生成 tasks/features.json）
    └── 2.3 ⭐ 设计门禁（/design:gate-check）

         ↓

阶段 3: Vibe Coding 实现 ⭐ [代码审查门禁]
    └── 3.1 编码循环（/coding:start）
    └── 3.2 ⭐ 代码审查门禁（/coding:review）
    └── 3.3 经验沉淀（/knowledge:extract-experience）
```

---

## 三道门禁详细定义

### ⭐ 门禁 1：需求门禁（阶段 1.3）

**触发指令**：`/requirement:gate-check`
**触发时机**：需求文档初稿完成后

**阻塞条件**（任一不满足则阻塞）：
- [ ] 需求背景描述清晰，有业务价值说明
- [ ] 验收标准具体、可测试
- [ ] 影响面分析完成（涉及哪些服务/模块）

**结论写入**：`requirements/{requirement-id}/gate-1-requirement-review.md`

---

### ⭐ 门禁 2：设计门禁（阶段 2.3）

**触发指令**：`/design:gate-check`
**触发时机**：方案设计与任务拆解完成后

**阻塞条件**（任一不满足则阻塞）：
- [ ] 方案覆盖所有需求条目
- [ ] 设计文档已包含核心要素：改动服务范围、整体流程、架构方案、性能瓶颈点、风险点
- [ ] 数据库变更方案（DDL/DML）已列出
- [ ] `tasks/features.json` 存在且格式合法，任务粒度合理

**结论写入**：`requirements/{requirement-id}/gate-2-design-review.md`

---

### ⭐ 门禁 3：代码审查门禁（阶段 3.2）

**触发指令**：`/coding:review`
**触发时机**：编码完成，准备提交前

**阻塞条件**（任一不满足则阻塞）：
- [ ] 代码符合团队规范（`context/team/`）
- [ ] 业务逻辑符合设计文档预期
- [ ] 单元测试覆盖核心逻辑

**结论写入**：`requirements/{requirement-id}/gate-3-code-review.md`

---

## 阶段状态机

```
REQUIREMENT_DEFINING → [GATE_1] → DESIGNING → [GATE_2] → CODING → [GATE_3] → DONE
```

当前阶段记录在：`requirements/{requirement-id}/status.json`

```json
{
    "requirement_id": "{requirement-id}",
    "current_stage": "DESIGNING",
    "current_step": "3.2",
    "gates_passed": ["GATE_1"],
    "last_updated": "2026-06-04T10:00:00+08:00"
}
```

---

## 错误代价递增曲线

```
代价
 ▲
 |                                   /
 |                              /
 |                         /
 |                    /
 |               /
 |          /
 |    /
 |___________________________________________→ 阶段
   阶段1      阶段2      阶段3
   (改需求)   (改设计)   (改代码/回滚)
```

⭐ 门禁正是设在"代价最低的拐点上"。
