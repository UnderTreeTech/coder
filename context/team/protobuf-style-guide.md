# Protobuf 定义规范 (团队级)

> **范围**：所有服务间通信的 gRPC 及基于 proto 定义的 HTTP 接口。
> **层级**：团队级（最高优先级，不可被项目级覆盖）

本规范提炼自基础服务（如 user 服务）的长期最佳实践。为了保证 IDL 文件的可读性与可维护性，所有新增的 `.proto` 文件及接口定义必须遵循以下规则：

---

## 1. 文件及包定义

- 语法：必须使用 `proto3`。
- 包命名：`package service.{service_name}.v1;` (例：`service.user.v1`)
- Go 包选项：`option go_package = "./;{service_name}";` (例：`option go_package = "./;user";`)
- **文件拆分规范**：存放路径为 `api/idl/{service_name}/`。为了保证业务清晰度，**严禁将所有定义揉在一个文件中**。必须按实体（Entity，例如 `user`、`order`）拆分为三个独立文件：
  1. `{entity}_iface.proto`：仅用于定义 Service 及其包含的增删改查等 RPC 方法。
  2. `{entity}_model.proto`：仅用于定义 RPC 接口的出入参（Request / Reply）等 `message` 定义。
  3. `{entity}_common.proto`：仅用于定义公共的 `enum` 枚举常量或被多个 model 共享的底层基础结构。

> **引用关系（Import）**：`iface` 文件应 `import "{entity}_model.proto"`；`model` 文件视需要 `import "{entity}_common.proto"`。

**示例（以 User 模块为例）：**

`user_common.proto`
```protobuf
syntax = "proto3";
package service.user.v1;
option go_package = "./;user";

// UserStatus 用户状态枚举
enum UserStatus {
  USER_STATUS_UNSPECIFIED = 0;
  USER_STATUS_ACTIVE = 1;
  USER_STATUS_BANNED = 2;
}
```

`user_model.proto`
```protobuf
syntax = "proto3";
package service.user.v1;
option go_package = "./;user";

import "user_common.proto";

// CreateUserReq 创建用户请求
message CreateUserReq {
  string name = 1;
  UserStatus status = 2;
}

// CreateUserReply 创建用户返回结果
message CreateUserReply {
  uint64 uid = 1;
}
```

`user_iface.proto`
```protobuf
syntax = "proto3";
package service.user.v1;
option go_package = "./;user";

import "user_model.proto";

// User 用户服务接口
service User {
  // CreateUser 创建用户
  rpc CreateUser(CreateUserReq) returns(CreateUserReply);
}
```

---

## 2. 注释规范（强制要求）

### 2.1 Service 与 RPC 方法
- **Service 分区注释**：使用区块注释将不同业务域的方法进行分组。
  ```protobuf
  //******************************************************************************//
  //*****************************订单管理接口**************************************//
  //******************************************************************************//
  ```
- **RPC 方法注释**：每个 `rpc` 必须有注释，注释格式为：`// {方法名} {描述说明}`。
  ```protobuf
  // CreateOrder 创建订单业务逻辑：处理金额计算、库存扣减
  rpc CreateOrder(CreateOrderReq) returns(CreateOrderReply);
  ```

### 2.2 Message 与字段（重点）
- **Message 注释**：所有 `message` 定义都**必须**包含说明注释。禁止出现没有任何说明的结构体。
- **字段注释**：对于有业务含义的字段，应尽量加注释说明其用途、枚举值或数据来源。

```protobuf
// CreateOrderReq 创建订单请求
message CreateOrderReq {
  uint64 cid = 1; // 企业ID
  uint64 uid = 2; // 用户ID
  uint64 sku_id = 3; // 商品SKU
  uint32 count = 4; // 购买数量
}

// CreateOrderReply 创建订单返回结果
message CreateOrderReply {
  uint64 order_id = 1; // 生成的订单ID
}
```

---

## 3. 命名约定

- **Service 名**：大驼峰（CamelCase），通常为名词，如 `Order`。
- **RPC 方法名**：大驼峰（CamelCase），动宾结构，如 `CreateOrder`。
- **Message 名**：大驼峰（CamelCase）。
  - 请求参数统一以 `Req` 结尾，如 `CreateOrderReq`。
  - 响应参数统一以 `Reply`（或 `Resp`） 结尾，如 `CreateOrderReply`。
- **字段名**：小写下划线分隔（snake_case），如 `order_id`、`created_time`。

---

## 4. 类型使用建议

- 对于 ID、时间戳、数量等，推荐使用 `uint64` 或 `uint32`。
- 若无特殊返回值，引用 `import "google/protobuf/empty.proto";` 并返回 `google.protobuf.Empty`。
- 通用的调用元信息（如操作人、来源等），推荐封装为独立的公用 Message（如 `RevokeInfo`）。
