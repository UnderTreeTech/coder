# Git 提交规范

## 1. 分支命名

### 功能分支（三仓必须相同）
```
feature/{devops-name}/{geelib-id}
```

例如：`feature/vip-refactor/T12345678`

### 修复分支
```
fix/{devops-name}/{geelib-id}
```

### 热修复分支
```
hotfix/{version}/{geelib-id}
```

### 发布分支
```
release/v{major}.{minor}.{patch}
```

---

## 2. 提交信息格式（Conventional Commits）

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### type 类型

| type | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档变更 |
| `style` | 代码格式（不影响逻辑） |
| `refactor` | 重构（不是新功能也不是 Bug 修复） |
| `perf` | 性能优化 |
| `test` | 测试相关 |
| `chore` | 构建/工具链变更 |
| `ci` | CI/CD 配置变更 |
| `revert` | 回滚提交 |

### 示例
```
feat(vip): add subscription renewal API

- 新增 VIP 续费接口 /vip/renew
- 支持月卡、季卡、年卡三种规格
- 接入埋点 vip_renew_submit

Closes #T12345678
```

---

## 3. 分支策略

### 主分支保护规则
- `main` / `master`：需要 2 人 review + CI 全通过
- `test`：需要 1 人 review + CI 全通过

### PR / MR 规范
- 标题格式：`[{geelib-id}] {简要描述}`
- 必须关联需求 ID
- 代码变更量超过 500 行时，需要分拆 PR
- IDL 变更必须在独立 PR 中处理

---

## 4. 三仓联动提交规范

当需求涉及多个仓库（业务仓 + IDL 仓 + Harness 仓）时：

1. 三个仓使用**完全相同的分支名**
2. IDL 变更先合并，业务代码后合并
3. Harness 仓（需求文档、设计文档、经验沉淀）随业务代码同步合并

---

## 5. 回滚规范

```bash
# 单仓回滚
git revert {commit-hash}

# 三仓同步回滚（按顺序）
# 1. 先回滚业务代码仓
# 2. 再回滚 IDL 仓（如有变更）
# 3. 最后更新 Harness 仓的 experience/ 记录本次回滚原因
```
