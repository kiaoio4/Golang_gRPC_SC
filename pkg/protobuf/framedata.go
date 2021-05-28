package protobuf

import (
	pb "grpc/pkg/proto"
	"grpc/pkg/grpctask"
	"fmt"

)

// FrameData -
type FrameData struct {
	Grpcmessage chan grpctask.GRPCTask
}

// FrameDataCallback -
func (t *FrameData) FrameDataCallback1(request pb.FrameData_FrameDataCallbackServer)  (err error) {
	tem, err := request.Recv()
	if err == nil {
            fmt.Println(tem)
	} else {
		fmt.Println("break, err :", err)
		return err
	}

	message := &grpctask.GRPCTask{
		Key:    tem.Key,
		Value:  tem.Value,
		IsStop: false,
	}

	t.Grpcmessage <- *message
	return  nil
}

func (t *FrameData) FrameDataCallback(request pb.FrameData_FrameDataCallbackServer)  (err error) {
	fmt.Println("start new server")
	ctx := request.Context()
	for{
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		tem, err := request.Recv()
		if err !=nil {
			return request.SendAndClose(&pb.FrameDataResponse{Successed:false})
		} 

		if err != nil {
          	fmt.Println(err)
        }
		message := &grpctask.GRPCTask{
			Key:    tem.Key,
			Value:  tem.Value,
			IsStop: false,
		}
		t.Grpcmessage <- *message
	}
	return nil
}
