# 五阶段流程 + 四门禁 — 唯一真相源

> ⚠️ **本文件是整条流程语义的唯一真相源。**
> AGENTS.md、每一个 Skill、每一个 Command 都围绕它保持一致。
> 一次更新，全仓生效，避免规范口径散落多处并逐渐漂移。

---

## 流程总览

```
阶段 1: 初始化
    └── 1.1 目录骨架创建
    └── 1.2 上下文加载（team + harness-framework + project）
    └── 1.3 服务矩阵确认

         ↓

阶段 2: 需求定义 ⭐ [需求评审门禁]
    └── 2.1 需求录入（/requirement:new）
    └── 2.2 需求撰写（AI 补齐背景、目标、非目标、验收标准）
    └── 2.3 ⭐ 需求评审门禁（requirement-quality-reviewer Agent）

         ↓

阶段 3: 设计 ⭐ [设计门禁]
    └── 3.1 技术预研（tech-feasibility-assessor Agent）
    └── 3.2 概要设计（outline-design-doc-writer Skill）
    └── 3.3 ⭐ 设计门禁（detail-design-quality-reviewer Agent）
    └── 3.4 详细设计（detail-design-doc-writer Skill）
    └── 3.5 追溯链建立（需求条目 → 设计决策 → 开发任务）

         ↓

阶段 4: 开发 ⭐⭐ [Dev门禁 + 服务仓库检查门禁]
    └── 4.1 任务拆分（tasks/features.json 生成）
    └── 4.2 ⭐ Dev 进入门禁（tasks/features.json 合法性检查）
    └── 4.3 ⭐ 服务仓库检查门禁（三仓分支一致性 + 服务仓库就位）
    └── 4.4 编码循环（选择上下文 → 编码 → 审查 → 提交）
        └── code-review-preparer Agent（收集 diff + 上下文）
        └── 8 个维度并行审查（code-review-report Skill）

         ↓

阶段 5: 交付
    └── 5.1 验收测试（test-runner Agent）
    └── 5.2 经验沉淀（/knowledge:extract-experience）
    └── 5.3 版本发布（releases/ 目录更新）
    └── 5.4 Done
```

---

## 四道门禁详细定义

### ⭐ 门禁 1：需求评审门禁（阶段 2.3）

**触发时机**：需求文档初稿完成后

**阻塞条件**（任一不满足则阻塞）：
- [ ] 需求背景描述清晰，有业务价值说明
- [ ] 目标明确，可量化
- [ ] 非目标明确列出（防止范围蔓延）
- [ ] 验收标准具体、可测试
- [ ] 影响面分析完成（涉及哪些服务/模块）
- [ ] 与历史需求无冲突

**执行 Agent**：`requirement-quality-reviewer`
**结论写入**：`requirements/{requirement-id}/gate-1-requirement-review.md`

---

### ⭐ 门禁 2：设计门禁（阶段 3.3）

**触发时机**：详细设计完成后

**阻塞条件**（任一不满足则阻塞）：
- [ ] 方案覆盖所有需求条目（追溯链完整）
- [ ] 服务边界清晰，无跨域耦合
- [ ] IDL 变更风险评估完成（字段冻结状态确认）
- [ ] 数据库变更方案（DDL/DML）已列出
- [ ] 配置变更已列出
- [ ] 回滚方案已设计
- [ ] 性能影响评估完成

**执行 Agent**：`detail-design-quality-reviewer`
**结论写入**：`requirements/{requirement-id}/gate-2-design-review.md`

---

### ⭐ 门禁 3：Dev 进入门禁（阶段 4.2）

**触发时机**：开始编码前

**阻塞条件**（任一不满足则阻塞）：
- [ ] `tasks/features.json` 存在且格式合法
- [ ] 所有任务已关联需求条目
- [ ] 任务粒度合理（单个任务不超过 1 天工作量）
- [ ] 分支已按规范创建（三仓同名）

**执行 Skill**：`dev-entry-gate-checker`
**结论写入**：`requirements/{requirement-id}/gate-3-dev-entry.md`

---

### ⭐ 门禁 4：服务仓库检查门禁（阶段 4.3）

**触发时机**：编码开始前（每次编码循环）

**阻塞条件**（任一不满足则阻塞）：
- [ ] 三仓分支名完全一致
- [ ] 业务代码仓分支已从最新 develop 切出
- [ ] IDL 仓分支已从最新 develop 切出（如涉及 IDL 变更）
- [ ] 服务仓库路径在 `.service-matrix/dependencies.yaml` 中已登记

**执行规则**：`context/harness-framework/gates/service-repo-check.md`
**结论写入**：`requirements/{requirement-id}/gate-4-service-repo-check.md`

---

## 阶段状态机

```
INIT → REQUIREMENT_DEFINING → [GATE_1] → DESIGNING → [GATE_2]
     → DEV_PREPARING → [GATE_3] → [GATE_4] → CODING → DELIVERING → DONE
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
 |                                              /
 |                                         /
 |                              /
 |                    /
 |          /
 |    /
 |___________________________________________→ 阶段
   阶段2    阶段3    阶段4.1  阶段4.3  阶段4.4
   (改几行  (改设计  (改任务  (重切分支  (回滚代码+
    文档)    文档)    拆分)    改环境)   回滚IDL+
                                        回滚数据迁移)
```

⭐ 门禁正是设在"代价最低的拐点上"。
