# 服务功能开发标准 SOP (Standard Operating Procedure)

本文档是基于项目架构和分层模型制定的**服务功能开发标准 SOP**。在进行需求开发（vibe coding）时，AI 和开发者都必须严格按照本规范进行编码实现，确保代码架构一致性、职责分明以及服务正确集成。

## 0. 架构分层规范（洋葱模型 / MVC）

本项目严格遵循标准的洋葱模型，分为以下几层，各层职责界限清晰：

- **Controller 层 (`internal/server/http/`)**
  - **职责**：接收 HTTP 请求，解析请求参数（Header, Query, Body），进行基础参数校验，并封装统一的 HTTP Response 返回。
  - **限制**：不得包含任何核心业务逻辑，仅作数据组装和路由处理。
- **Service 层 (`internal/service/`)**
  - **职责**：核心业务逻辑的承载者。负责组合调用多个 DAO 层方法或其他外部 RPC 接口来完成复杂的业务场景。
  - **限制**：绝对不包含任何与底层协议（HTTP/gRPC）直接相关的代码。
- **Iface 层 (Interfaces / 接口层)**
  - **职责**：定义 Service 和 DAO 之间、或内部各模块之间的契约，方便 Mock 测试、依赖注入和层与层之间的解耦。
- **DAO 层 (`internal/dao/`)**
  - **职责**：Data Access Object，数据访问层。负责与底层存储（如 MySQL, PostgreSQL, MongoDB, Redis 等）直接交互。
  - **限制**：必须屏蔽底层存储细节，对上层提供业务友好的数据结构和接口。

---

## 1. 服务的开启与关闭 (HTTP vs gRPC)

项目中默认会在服务入口 `cmd/main.go` 中同时注册 **HTTP** 和 **gRPC** 两种服务，但在**实际生产环境中，通常建议根据微服务职责只开启一种服务**。

**操作指南**（编辑 `cmd/main.go`）：
```go
http := http.New(s)
rpc := grpc.New(s)

etcd.Register(context.Background(), rpc.ServiceInfo)
etcd.Register(context.Background(), http.ServiceInfo)
```
- **若仅作为 HTTP 服务**：请将与 `grpc.New(s)` 及其相关的注册、Stop 等代码进行注释或删除。
- **若仅作为 gRPC 服务**：请将与 `http.New(s)` 及其相关的注册、Stop 等代码进行注释或删除。

---

## 2. HTTP 接口开发 SOP

当你需要为服务增加一个 HTTP 接口时，请严格按照以下顺序自下而上（或自上而下）进行实现：

### Step 2.1: 实现 DAO 层 (`internal/dao/{entity}.go`)
直接与数据库交互，实现数据增删改查。
```go
package dao

import "context"

func (d *Dao) GetUserByID(ctx context.Context, id string) (*User, error) {
    var user User
    // 执行 SQL 查询并将结果映射到实体
    return &user, nil
}
```

### Step 2.2: 实现 Service 层 (`internal/service/{entity}.go`)
组合 DAO 方法，处理核心业务逻辑，不涉及 HTTP 语义。
```go
package service

import "context"

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // 业务校验，组合调用 dao
    return s.dao.GetUserByID(ctx, id)
}
```

### Step 2.3: 实现 Controller 及路由挂载 (`internal/server/http/server.go` 或独立 `router.go`)
解析参数，调用 Service 并响应。
```go
// 1. Controller 方法定义
func getUser(s *service.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Query("id")
        
        // 调用 Service
        user, err := s.GetUser(c.Request.Context(), id)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()}) // 依据团队错误码规范处理
            return
        }
        c.JSON(200, user)
    }
}

// 2. 路由注册
func New(s *service.Service) *Server {
    engine := gin.Default()
    engine.GET("/api/v1/user", getUser(s)) // 挂载路由
    return &Server{engine: engine}
}
```

---

## 3. gRPC 接口开发 SOP

当你需要为服务增加一个 gRPC 接口时，请严格按照以下步骤进行：

### Step 3.1: 编写 Protobuf 文件 (`api/{version}/{entity}.proto`)
定义 RPC 服务契约、Request 及 Reply。
```protobuf
syntax = "proto3";
package api.v1;
option go_package = "{repo}/api/v1;v1"; // 替换为真实包路径

service UserService {
    rpc GetUser (GetUserRequest) returns (GetUserReply);
}

message GetUserRequest {
    string id = 1;
}

message GetUserReply {
    string name = 1;
}
```
> **要求**：编写完成后，必须执行 `protoc` (或项目对应的生成命令) 生成对应的 Go 桩代码。

### Step 3.2: Service 层实现 gRPC 接口 (`internal/service/{entity}.go`)
在 Service 层实现 Protobuf 生成的 `Server` 接口。
```go
package service

import (
    "context"
    pb "github.com/UnderTreeTech/layout/api/v1" // 引入 pb 包
)

// 必须实现 pb 定义的 RPC 接口
func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
    // 调用 DAO 层核心逻辑
    user, err := s.dao.GetUserByID(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    
    // 将业务层结构转换为 PB Reply 结构返回
    return &pb.GetUserReply{
        Name: user.Name,
    }, nil
}
```

### Step 3.3: gRPC Server 注册 (`internal/server/grpc/server.go`)
在 gRPC Server 启动入口处，将 Service 实例绑定到 gRPC Server 上。
```go
package grpc

import (
    "google.golang.org/grpc"
    pb "github.com/UnderTreeTech/layout/api/v1"
    "github.com/UnderTreeTech/layout/internal/service"
)

func New(s *service.Service) *Server {
    srv := grpc.NewServer()
    // 注册 gRPC 服务
    pb.RegisterUserServiceServer(srv, s)
    return &Server{Server: srv}
}
```

---

## 总结

AI（Claude）在接收到针对本仓库的具体开发任务时：
1. **优先确认**：当前需求是需要开发 HTTP 接口还是 gRPC 接口，或者仅仅是内部逻辑。
2. **遵守分层**：任何逻辑变更必须放置在正确的层（Router, Controller, Service, DAO）。不要跨层污染（例如在 DAO 层中解析 gin.Context）。
3. **闭环验证**：确保新增代码涉及到的 `main.go`、路由注册或 gRPC Server 注册都已经正确就绪。