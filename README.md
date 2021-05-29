# Golang_gRPC_SC

使用基于HTTP/2的gRPC流式传输，使用三种类型流中的客户端流式请求

- [ ] 服务端流式响应
- [x] 客户端流式请求
- [ ] 两端双向流式

服务端for循环调用stream.Rece()，接收客户端消息并阻塞，等客户端调用stream.CloseAndRecv()
关闭流的发送后进入阻塞监听，服务端调用stream.SendAndClose()，返回响应体并关闭流。
此方式客户端只负责发送流的结束，服务端可以在中途结束整个流处理。

## proto文件编译

windows
```bash
protoc --proto_path=pkg/proto pkg/proto/*.proto --go_out=plugins=grpc:pkg/proto -Ipkg pkg/proto/*.proto
```

Linux
```bash
.PHONY: generate proto file

gen:
	protoc -I ./proto ./proto/frame.proto --go_out=plugins=grpc:./pb

clean:
	rm pb/*.go

```