package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
	"unsafe"
	proto "grpc/pkg/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

)

var (
	wg sync.WaitGroup
)

const (
	networkType = "tcp"
	server      = "localhost"
	port        = "41005"
	parallel    = 1     //连接并行度
	times       = 10000 //每连接请求次数
)

func mai1n() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	currTime := time.Now()

	//并行请求
	for i := 0; i < int(parallel); i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			// exe(num)
		}(i)
	}
	wg.Wait()

	log.Printf("time taken: %.2f ", time.Now().Sub(currTime).Seconds())
}
var (
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "echo.Echo","method":"UnaryEcho"}],
		  "retryPolicy": {
			  "MaxAttempts": 100, 
			  "InitialBackoff": "0.1s",
			  "MaxBackoff": "1s",
			  "BackoffMultiplier": 2,
			  "RetryableStatusCodes": [ "UNKNOWN" ]
		  }
		}]}`
)
/*
MaxAttempts：最大尝试次数
InitialBackoff：默认退避时间
MaxBackoff：最大退避时间
BackoffMultiplier：退避时间增加倍率
RetryableStatusCodes：服务端返回什么错误码才重试
*/
var kacp = keepalive.ClientParameters{
    Time:                10 * time.Second, // 客户端没有回调服务，每10秒ping一次服务端
    Timeout:             time.Second,      // 如果服务端已经断开了，等待ping后回应的ack多少时间
    PermitWithoutStream: true,             // 没有数据流时是否允许ping
}
func main() {
	go func() {
Retry:
		var grpcsize = 1024 * 1024 * 128 // 128MB
		conn, err := grpc.Dial(
			server+":"+port,
			grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcsize)),
			grpc.WithKeepaliveParams(kacp),
			grpc.WithDefaultServiceConfig(retryPolicy),
		)
		if err != nil {
			log.Fatal("error....", err)
			time.Sleep(time.Second * 2)
			goto Retry
		}

		client := proto.NewFrameDataClient(conn)
		stream, err := client.FrameDataCallback(context.Background())
		if err != nil {
			
			fmt.Println("openn stream error ", err)
			time.Sleep(time.Second * 2)
			goto Retry
		}
		for{  
			if conn != nil {
				now := time.Now().UTC()
				request := proto.FrameDataRequest{
					Key:   str2bytes(now.String()),
					Value: []byte{0x01, 0x02, 0x03, 0x04},
				}
				if err := stream.Send(&request); err != nil {
					fmt.Println("can not send ", err)
					goto Retry
				}

				time.Sleep(time.Millisecond * 1)

			}
		}
	}()

	select {}
}


func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
