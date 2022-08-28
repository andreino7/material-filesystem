package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Close(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	closeReq := request.GetClose()
	if closeReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	daemon.fs.Close(closeReq.GetFileDescriptor())

	return &pb.Response{
		Response: &pb.Response_Close{
			Close: &pb.CloseResponse{},
		},
	}, nil
}
