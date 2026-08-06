# grpc

## 消息大小限制

服务端创建时通过 `grpc.maxReceiveMessageBytes` 和 `grpc.maxSendMessageBytes` 设置接收、发送上限. `protoc-gen-go-grpc-x` 生成的客户端连接也会读取同一组配置, 分别设置客户端接收、发送上限.

通信双方都必须允许目标消息大小. 例如 A 向 B 发送 64MiB 以内的消息时, A 的发送上限和 B 的接收上限都必须至少为 67108864 bytes. 建议同一套服务统一配置这两个字段, 避免不同调用方向出现不对称限制.

# 路由功能
## 使用片键 shardkey, 通过不同的 loadblance policy, 路由到 service.
### 1. serverkey → IClientConn (1:1)
- **key**: serverkey (etcd 中存储的服务器信息)
- **value**: IClientConn (单个 gRPC 客户端连接)
- **关系**: 一对一映射，每个服务器对应一个连接

### 2. method → IClientConn (1:n)
- **key**: method (proto 文件中定义的方法)
- **value**: []IClientConn (方法对应的多个连接数组)
- **关系**: 一对多映射，一个方法可以对应多个连接

### 3. method&&shardkey → IClientConn (1:1)
- **key**: method&&shardkey (方法名和片键的组合结构体)
- **value**: IClientConn (单个 gRPC 客户端连接)
- **关系**: 一对一映射，特定方法和片键组合对应特定连接

```
                    ┌─────────────────────────────────────┐
                    │           IClientConn               │
                    │        (gRPC 客户端连接)              │
                    └─────────────────┬───────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    │                 │                 │
         ┌──────────▼─────────┐ ┌─────▼─────┐ ┌─────────▼─────────┐
         │     serverkey      │ │   method  │ │ method&&shardkey  │
         │ (etcd 中 server    │  │(proto 中的│  │ (方法+片键组合)    │
         │  信息)             │  │方法)      │  │                  │
         └────────────────────┘ └───────────┘ └───────────────────┘
                    │                 │                 │
                    │                 │                 │
                    │ 1:1 关系         │ 1:n 关系        │ 1:1 关系
                    │                 │                 │
                    │ 每个服务器对应     │ 一个方法对应      │ 特定方法和片键
                    │ 一个连接          │ 多个连接         │ 组合对应特定连接
                    │                 │                 │
```

## shardKey 片键
- 配置在 proto 文件中

## traceID 追踪ID
- 由服务端生成
- 在协议传递中, 存放于 context 中

## timeout 超时
- 配置在 proto 文件中


# 备注
## 字段 method: e.g.:"/packageName.serviceName/methodName"
