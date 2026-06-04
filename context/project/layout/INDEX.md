# layout 项目知识库

> **项目名**：`layout`
> **业务域**：Go 微服务标准布局模板
> **维护团队**：{team-name}

---

## 模块列表

| 模块 | 说明 | 路径 |
|-----|------|------|
| `internal` | 核心业务逻辑 | [`internal/INDEX.md`](internal/INDEX.md) |

---

## 项目架构概览

```
layout/
├── api/          # Protobuf IDL 定义
├── cmd/          # 程序入口
├── configs/      # 配置文件
├── internal/     # 内部业务逻辑（禁止外部引用）
│   ├── dao/      # 数据访问层
│   ├── ecode/    # 错误码定义
│   ├── i18n/     # 国际化
│   ├── model/    # 数据模型
│   ├── server/   # HTTP/gRPC 服务器
│   ├── service/  # 业务逻辑层
│   └── utils/    # 工具函数
└── templates/    # 代码生成模板
```

---

## 技术栈

| 组件 | 技术选型 | 版本 |
|-----|---------|------|
| Web 框架 | gin | v1.12.0 |
| RPC 框架 | gRPC | v1.80.0 |
| ORM/SQL | squirrel | v1.4.0 |
| 数据库 | MySQL / OpenGauss | - |
| 缓存 | Redis | go-redis v9 |
| 消息队列 | - | - |
| 服务注册 | etcd | v3.5.21 |
| 链路追踪 | Jaeger | - |
| 监控 | Prometheus | - |

---

## 关键约束（AI 必读）

1. `internal/` 目录下的包**禁止**被外部项目直接引用
2. 数据库操作必须通过 `internal/dao/` 层，禁止在 service 层直接操作 DB
3. 错误码必须使用 `internal/ecode/` 中定义的常量
4. HTTP 响应格式统一使用 `internal/utils/reply/` 中的工具函数
