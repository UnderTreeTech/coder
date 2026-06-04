# 项目整体说明 (Project README)

## ⚠️ 重要提示：HTTP/GRPC 服务的开启与关闭

在本项目中，默认会在 `cmd/main.go` 里同时启动 **HTTP** 和 **GRPC** 两种服务。

**实际生产环境中，一般建议只开启一种服务**（或者根据具体微服务职责划分）。
如何开启或关闭：
打开 `cmd/main.go` 文件，定位到以下代码段：
```go
http := http.New(s)
rpc := grpc.New(s)

etcd.Register(context.Background(), rpc.ServiceInfo)
etcd.Register(context.Background(), http.ServiceInfo)
```
- **仅保留 HTTP 服务**：注释或删除 `rpc := grpc.New(s)` 及其注册和 Stop 代码。
- **仅保留 GRPC 服务**：注释或删除 `http := http.New(s)` 及其注册和 Stop 代码。

---

## 1. HTTP 服务说明

### 路由入口在哪？
HTTP 服务的入口及路由定义主要在 `internal/server/http/server.go`（或者 `router.go`）中。在这里会初始化 HTTP Server 并挂载对应的路由。

### 如何添加路由及代码实现？
以下是添加一个 HTTP 接口的完整代码实现流程：

**1. 路由注册与 Controller 层 (`internal/server/http/server.go`)**
```go
package http

import (
	"github.com/gin-gonic/gin"
	"github.com/UnderTreeTech/layout/internal/service"
)

// New 初始化 HTTP 服务
func New(s *service.Service) *Server {
	engine := gin.Default()

	// 添加路由，绑定 Controller 方法
	engine.GET("/api/v1/user", getUser(s))

	return &Server{engine: engine}
}

// Controller 层实现：解析参数并调用 Service
func getUser(s *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		// 调用 Service 层
		user, err := s.GetUser(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, user)
	}
}
```

**2. Service 层 (`internal/service/user.go`)**
```go
package service

import "context"

// GetUser 核心业务逻辑处理
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	// 可以在这里组合多个 DAO 调用或其他业务逻辑
	return s.dao.GetUserByID(ctx, id)
}
```

**3. DAO 层 (`internal/dao/user.go`)**
```go
package dao

import "context"

// GetUserByID 数据访问层，直接与数据库交互
func (d *Dao) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	// 伪代码：执行 SQL 查询
	// err := d.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, nil
}
```

### 代码分层拆解介绍
本项目采用标准的洋葱模型/MVC架构分层设计：
- **Controller 层 (`internal/server/http/`)**：负责接收 HTTP 请求，解析参数（Header, Query, Body），进行基础的参数校验，并封装统一的 HTTP Response。
- **Service 层 (`internal/service/`)**：核心业务逻辑的承载者。组合调用多个 DAO 层方法或其他 RPC 接口来完成复杂的业务场景。不包含任何与协议（HTTP/gRPC）相关的代码。
- **Iface 层 (Interfaces / 接口层)**：定义 Service 和 DAO 之间、或内部各模块之间的契约，方便 Mock 测试和解耦。
- **DAO 层 (`internal/dao/`)**：Data Access Object，数据访问层。负责与底层存储（如 MySQL, PostgreSQL, MongoDB, Redis 等）直接交互，屏蔽底层细节。

---

## 2. GRPC 服务说明

### Service 中如何实现 RPC 定义及代码实现？
作为 GRPC 服务时，服务的定义与实现流程如下：

**1. 编写 Protobuf 文件 (`api/v1/user.proto`)**
```protobuf
syntax = "proto3";
package api.v1;
option go_package = "layout/api/v1;v1";

// 定义 RPC 服务
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
*编写完成后，使用 `protoc` 工具生成对应的 Go 代码。*

**2. 在 Service 层实现接口 (`internal/service/user.go`)**
生成的 Go 代码中会包含一个 `UserServiceServer` 接口。我们需要在 Service 层实现它：
```go
package service

import (
	"context"
	pb "github.com/UnderTreeTech/layout/api/v1"
)

// GetUser 实现 Protobuf 中定义的 RPC 方法
func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	// 调用 DAO 层获取数据
	user, err := s.dao.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 组装并返回 Protobuf 定义的 Reply 结构
	return &pb.GetUserReply{
		Name: user.Name,
	}, nil
}
```

**3. 注册 GRPC 服务 (`internal/server/grpc/server.go`)**
在 GRPC Server 初始化时，将实现了接口的 Service 实例注册进去：
```go
package grpc

import (
	"google.golang.org/grpc"
	pb "github.com/UnderTreeTech/layout/api/v1"
	"github.com/UnderTreeTech/layout/internal/service"
)

// New 初始化 GRPC 服务
func New(s *service.Service) *Server {
	srv := grpc.NewServer()

	// 将 Service 实例注册到 GRPC Server 中
	pb.RegisterUserServiceServer(srv, s)

	return &Server{Server: srv}
}