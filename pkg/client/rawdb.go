package main

import (
	"fmt"
	"grpc/pkg/grpctask"
	proto "grpc/pkg/proto"
	"grpc/pkg/protobuf"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

)

const (
	port = "41005"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var grpcsize = 1024 * 1024 * 128 // 128MB
	var kaep = keepalive.EnforcementPolicy{
		// MinTime:             time.Second * 60,
		PermitWithoutStream: true,
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle: 3600 * time.Second,
		Time:              10 * time.Second,
		Timeout:           1 * time.Second,
	}
	var grpcmessage chan grpctask.GRPCTask
	grpcmessage = make(chan grpctask.GRPCTask, 1024)
	sigchan := make(chan os.Signal, 1)
	go GRPCQueuetask(sigchan, grpcmessage)
	//起服务
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
	grpc.KeepaliveParams(kasp),
	grpc.KeepaliveEnforcementPolicy(kaep),
	grpc.MaxRecvMsgSize(grpcsize),
	grpc.MaxSendMsgSize(grpcsize),
	)

	proto.RegisterFrameDataServer(s, &protobuf.FrameData{Grpcmessage: grpcmessage})
	s.Serve(lis)

	fmt.Println("grpc server in: ", port)
}


// GRPCQueuetask -
func GRPCQueuetask(sigchan chan os.Signal, message chan grpctask.GRPCTask) {
	isStop := false
	count := 0
	for {
		if isStop {
			break
		}

		select {
		case sig := <-sigchan:
			fmt.Println("Caught signal : terminating", sig)
			return

		case mes := <-message:
			isStop = mes.IsStop
			if len(mes.Key) > 0 {
				fmt.Println("**********Task gRPC message******************", count)
			}
		}
		count++
	}
}
