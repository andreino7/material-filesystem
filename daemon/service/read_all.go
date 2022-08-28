package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) AppendAll(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	appendReq := request.GetAppendAll()
	if appendReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return appendReq.GetPath() })
	if err != nil {
		return nil, err
	}

	err = daemon.fs.AppendAll(path, appendReq.GetContent())
	if err != nil {
		return daemon.extractError(request.GetSessionId(), err)
	}

	return &pb.Response{
		Response: &pb.Response_AppendAll{
			AppendAll: &pb.AppendAllResponse{},
		},
	}, nil
}
