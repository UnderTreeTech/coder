# Command: /service:deps

> **说明**：查看服务依赖关系，分析影响面。
> **触发**：用户输入 `/service:deps` 或 `/service:deps {service-name}`

---

## 执行步骤

### Step 1: 加载服务矩阵

读取 `.service-matrix/dependencies.yaml`。

### Step 2: 调用 service-dependency-analyzer Skill

委派给 `.claude/skills/service-dependency-analyzer.md` 执行分析。

### Step 3: 输出依赖图

以可读格式输出服务依赖关系。

---

## 参数

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|-----|
| `service` | string | ❌ | 服务名，不填则显示所有服务 |

---

## 示例

```
/service:deps
/service:deps dao
/service:deps http-server
```

输出示例：
```
📊 服务依赖图

http-server
  ↓ 调用
  service
    ↓ 调用
    dao
      ↓ 依赖
      MySQL (主数据库)
      Redis (缓存)

grpc-server
  ↓ 调用
  dao
    ↓ 依赖
    MySQL (主数据库)
```
