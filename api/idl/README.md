# 全局 IDL 目录 — README

> **说明**：本目录是所有服务 Protobuf IDL 的统一管理目录。
> 各服务的 proto 文件集中在此，便于跨服务引用和版本管理。

---

## 目录结构

```
api/idl/
├── README.md                    # 本文件
├── order/
│   └── v1/
│       └── order.proto          # 订单服务 gRPC 接口定义
├── pay/
│   └── v1/
│       └── pay.proto            # 支付服务 gRPC 接口定义
└── user/
    └── v1/
        └── user.proto           # 用户服务 HTTP 接口定义（使用 grpc-gateway）
```

---

## 协议说明

| 服务 | 协议 | 说明 |
|-----|------|------|
| `order` | **gRPC** | 内部服务间调用，高性能 |
| `pay` | **gRPC** | 内部服务间调用，高性能 |
| `user` | **HTTP** | 对外暴露，使用 grpc-gateway 转换 |

---

## IDL 变更规范

1. **新增字段**：可以直接添加，不影响已有调用方
2. **修改字段类型**：**禁止**，会破坏二进制兼容性
3. **删除字段**：**禁止**，改为标记 `deprecated`
4. **修改字段编号**：**严格禁止**，会导致数据解析错误
5. **新增 RPC 方法**：可以直接添加

---

## 生成代码

```bash
# 生成 Go 代码（在仓库根目录执行）
cd idl/chat/v1
waterdrop protoc --grpc --swagger *.proto
```
---

## 冻结字段登记

冻结字段在 `.service-matrix/dependencies.yaml` 的 `idl_frozen_fields` 中登记。
