# {ServiceName} 服务（{service-name}）知识库

> **使用说明**：复制本模板到 `context/project/{project-name}/{service-name}/INDEX.md`
> 由 `/service:onboard` 命令自动生成时预填已知字段，其余字段由开发者补充。

---

> **服务名**：`{service-name}`
> **代码路径**：`api/{service-name}/`
> **通信协议**：`{http|grpc}`
> **业务域**：{业务域描述}

---

## 架构概览

```
api/{service-name}/
├── api/v1/             # （无独立 IDL，使用全局 api/idl/{service-name}/v1/）
├── cmd/main.go
├── configs/application.toml
└── internal/
    ├── dao/
    ├── ecode/          # 错误码（{起始码}-{结束码}，在 context/team/error-code.md 中登记）
    ├── model/
    ├── server/
    │   ├── http/       # （HTTP 服务时使用）
    │   └── grpc/       # （gRPC 服务时使用）
    └── service/
```

---

## 外部依赖

| 依赖 | 类型 | 说明 |
|-----|------|------|
| {service-name} | {上游/下游/外部} | {说明} |
| MySQL | 数据库 | {数据库名} |

---

## 关键约束（AI 必读）

> ⚠️ 以下约束是从业务需求和踩坑经验中提炼的，AI 在处理本服务代码时**必须**遵守。

1. {约束 1}
2. {约束 2}

---

## 踩坑经验

详见 [`experience/`](experience/) 目录。

---

## API 说明

主要接口（详见 `api/idl/{service-name}/v1/{service-name}.proto`）：

| 接口 | 方法/RPC | 说明 |
|-----|---------|------|
| {接口路径或 RPC 名} | {GET/POST/rpc} | {说明} |
