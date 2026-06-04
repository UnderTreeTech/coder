# 版本发布管理 — INDEX

> **说明**：每个版本发布在此目录下有独立的子目录，包含该版本的所有变更清单。
> 运维可按此版本目录自动化做升级和出包。

---

## 目录结构

```
releases/
├── INDEX.md                          # 本文件
├── v{major}.{minor}.{patch}/         # 版本目录（语义化版本）
│   ├── RELEASE_NOTES.md              # 版本说明（变更摘要、影响面）
│   ├── checklist.md                  # 发布前检查清单
│   ├── sql/                          # 数据库变更
│   │   ├── ddl/                      # 结构变更（建表/改表/加索引）
│   │   │   └── {序号}-{描述}.sql
│   │   └── dml/                      # 数据变更（初始化数据/迁移数据）
│   │       └── {序号}-{描述}.sql
│   ├── configs/                      # 配置变更
│   │   ├── {env}/                    # 按环境分组（dev/test/staging/prod）
│   │   │   └── {配置项}.yaml
│   │   └── diff/                     # 配置变更 diff 说明
│   │       └── {配置项}-diff.md
│   └── scripts/                      # 运维脚本
│       ├── pre-deploy/               # 部署前执行
│       │   └── {序号}-{描述}.sh
│       ├── post-deploy/              # 部署后执行
│       │   └── {序号}-{描述}.sh
│       └── rollback/                 # 回滚脚本
│           └── {序号}-{描述}.sh
└── _template/                        # 版本目录模板
    ├── RELEASE_NOTES.md
    ├── checklist.md
    ├── sql/
    │   ├── ddl/.gitkeep
    │   └── dml/.gitkeep
    ├── configs/
    │   └── diff/.gitkeep
    └── scripts/
        ├── pre-deploy/.gitkeep
        ├── post-deploy/.gitkeep
        └── rollback/.gitkeep
```

---

## SQL 变更规范

### 文件命名
```
{序号}-{描述}.sql
```
例如：`001-add-user-vip-level-column.sql`、`002-init-vip-config-data.sql`

### 执行顺序
SQL 文件按**文件名序号升序**执行，必须保证幂等性。

### DDL 变更约束

```sql
-- ✅ 正确：添加字段使用 IF NOT EXISTS（MySQL 5.7 以上）
ALTER TABLE `users` ADD COLUMN IF NOT EXISTS `vip_level` TINYINT NOT NULL DEFAULT 0 COMMENT 'VIP等级';

-- ✅ 正确：创建索引使用 IF NOT EXISTS
CREATE INDEX IF NOT EXISTS `idx_users_vip_level` ON `users` (`vip_level`);

-- ❌ 禁止：直接 DROP 列（需要确认无下游依赖）
-- ALTER TABLE `users` DROP COLUMN `old_column`;

-- ❌ 禁止：修改已有字段类型（可能破坏数据兼容性）
-- ALTER TABLE `users` MODIFY COLUMN `id` VARCHAR(64);
```

### DML 变更约束

```sql
-- ✅ 正确：INSERT 使用 INSERT IGNORE 或 ON DUPLICATE KEY UPDATE
INSERT IGNORE INTO `config` (`key`, `value`) VALUES ('vip_price_monthly', '30');

-- ✅ 正确：批量 UPDATE 限制影响行数，分批执行
UPDATE `users` SET `vip_level` = 1 WHERE `is_vip` = 1 LIMIT 1000;

-- ❌ 禁止：无 WHERE 条件的 UPDATE/DELETE
-- UPDATE `users` SET `status` = 0;
```

---

## 配置变更规范

### 目录结构

```
configs/
├── dev/        # 开发环境
├── test/       # 测试环境
├── staging/    # 预发布环境
└── prod/       # 生产环境（敏感配置不提交，使用密钥管理）
```

### Diff 说明格式

每个配置变更必须有对应的 diff 说明文件：

```markdown
# 配置变更说明：{配置项名称}

## 变更原因
{为什么需要修改这个配置}

## 变更内容
| 配置项 | 原值 | 新值 | 说明 |
|-------|------|------|-----|
| `xxx.yyy` | `old_value` | `new_value` | {说明} |

## 影响范围
{哪些服务/功能受影响}

## 生效时间
- 是否需要重启服务：{是/否}
- 建议生效时间：{低峰期/任意时间}
```

---

## 运维脚本规范

### 脚本命名
```
{序号}-{描述}.sh
```
例如：`001-check-db-connection.sh`、`002-warm-up-cache.sh`

### 脚本要求

```bash
#!/usr/bin/env bash
# =============================================================================
# 脚本说明：{描述}
# 执行阶段：{pre-deploy/post-deploy/rollback}
# 预计耗时：{X 分钟}
# 是否幂等：{是/否}
# 依赖：{无/依赖列表}
# =============================================================================

set -euo pipefail

# 错误处理
trap 'echo "Error on line $LINENO"' ERR

# 脚本内容
echo "执行：{操作描述}"

# ... 脚本逻辑 ...

echo "✅ 完成：{操作描述}"
```

---

## 发布流程

```
1. 在需求研发完成后（阶段 5.3），创建版本目录
2. 收集本次需求涉及的 SQL/配置/脚本变更
3. 填写 RELEASE_NOTES.md 和 checklist.md
4. 提交 code review
5. 运维按 checklist.md 执行发布
```

---

## 快速创建版本目录

```bash
# 创建新版本目录
VERSION="v1.0.0"
cp -r releases/_template releases/${VERSION}

# 或使用脚本
bash scripts/new-release.sh v1.0.0
```
