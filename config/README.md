# config 配置

## gRPC 消息大小

`grpc.maxReceiveMessageBytes` 和 `grpc.maxSendMessageBytes` 分别控制单条 gRPC 接收、发送消息的最大字节数, 配置值必须大于 0.

未配置时保持 gRPC 原有默认行为:

- 接收上限为 4194304 bytes, 即 4MiB.
- 发送上限为 2147483647 bytes.

示例:

```yaml
grpc:
  maxReceiveMessageBytes: 67108864
  maxSendMessageBytes: 67108864
```
