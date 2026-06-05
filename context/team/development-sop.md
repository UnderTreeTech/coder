# 团队开发标准 SOP (Standard Operating Procedure)

> **核心原则**：各服务代码实现必须遵循本标准套路，严禁 AI 或开发者在没有达成共识的情况下自行发挥。
>
> **特别强调：绝不瞎篡改、自创代码规范，具体服务的架构约束完全以其目录下的 `README.md` 为准！**

## 1. 架构分层与实现路径

业务功能的实现必须在设计门禁时确认，编码时严格遵守以下各层级的规范套路（以 `api/user` 架构为范本）：

### 1.1 数据访问层 (Model / IFace / DAO)
底层数据访问层采用接口隔离与依赖注入的设计，**一般情况下通过 `xo` 等组件自动生成**：
- **Model 层 (`internal/model`)**：定义与数据库表映射的结构体实体。例如 `internal/model/tuser.xo.go`。
- **IFace 层 (`internal/iface`)**：定义 DAO 层必须实现的接口。例如 `internal/iface/tuser.xo.go`。
- **DAO 注册与组合 (`internal/dao.go` 等)**：在 DAO 基础文件中，将 `iface` 中定义的各个接口组合进主 `Dao` 结构体/接口中。
- **DAO 实现层 (`internal/dao`)**：实现 `iface` 定义的接口逻辑，负责具体的 SQL 拼接与数据库交互。例如 `internal/dao/tuser.xo.go`。
> **AI 约束**：AI 在手写补充 DAO 方法时，必须同步维护 `iface` 的接口定义并在 `dao` 目录下实现，保持模式的一致性。

### 1.2 gRPC 接口服务 (gRPC Service)
- **场景**：服务间的高效内部 RPC 调用。
- **套路规范**：
  - **IDL 定义与生成**：在 `api/idl/{service_name}/` 找相应的 proto 文件定义接口，并**必须主动执行** `waterdrop protoc --grpc --swagger *.proto` 生成新的 stub 文件。
  - **业务实现 (`internal/service`)**：在 `service` 目录下实现 proto 中定义的 RPC 接口。
    - **实现样板**：强烈建议参考团队统一固化的 [`code-implementation-examples.md`](code-implementation-examples.md) 中的 gRPC 样板代码，严格遵循“入参校验 -> 核心业务逻辑 -> 调用 DAO”的流程。
  - **服务注册**：在 `internal/server/grpc/server.go` 中将实现类注册到 gRPC Server。

### 1.3 HTTP 接口服务 (HTTP Controller)
- **场景**：直接暴露给前端、第三方或内部基于 HTTP 的通信。
- **套路规范**：
  - **1. 路由注册**：在 `internal/server/http/router.go` 中注册 HTTP 路由。
  - **2. 定义入参/出参 (`http/model`)**：在 `internal/server/http/model/` 目录下（或指定的 http dto 目录），专门定义该接口的 Request 和 Response 结构体，隔离底层模型与外部表现。
  - **3. Controller 实现 (`http` 目录)**：在 `internal/server/http/{module}.go` 中实现控制器逻辑（如 `http/user.go`）。
    - 职责：解析 Gin 上下文参数 -> 转换为业务参数 -> 调用 `internal/service` 的逻辑层方法 -> 封装响应返回。
  - **4. 业务下沉**：Controller 中**严禁写重度业务逻辑或直接调 DAO**，必须将核心业务下沉到 `internal/service` 中的普通 Service 方法中去处理。

### 1.4 内部普通 Service 方法
- **场景**：服务内部的公共业务逻辑抽取，供 HTTP Controller、gRPC Service 或定时任务、MQ 消费者调用。
- **套路规范**：直接在 `internal/service/` 下建立对应的业务文件实现即可，纯 Go 方法，无对外网络协议绑定。

## 2. 核心红线 (AI 必须遵守)

1. **以服务级文档为尊**：`api/{service_name}/README.md` 是该服务的绝对规范源，任何生成代码的结构、命名、分层均需照抄该文档的范式，**不得发明所谓的 biz 层或凭空改造目录结构**。
2. **规范 IDL 编写**：对于 gRPC 接口，AI 应主动基于 `context/team/protobuf-style-guide.md` 规范给出 proto 变更，必须提示执行 `protoc` 生成代码命令。
## 3. 代码收尾要求

1. **自动格式化**：在完成一个功能模块的代码编写后，必须主动执行 `go fmt ./...` 对当前修改的代码进行格式化。
2. **存根更新**：所有 `.proto` 的变更，必须在代码编写后、提交前，确保执行过生成命令，保证存根与定义一致。
