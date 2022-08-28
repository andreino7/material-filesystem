package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Open(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	openReq := request.GetOpen()
	if openReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return openReq.GetPath() })
	if err != nil {
		return nil, err
	}

	fd, err := daemon.fs.Open(path)
	if err != nil {
		return daemon.extractError(request.GetSessionId(), err)
	}

	return &pb.Response{
		Response: &pb.Response_Open{
			Open: &pb.OpenResponse{FileDescriptor: fd},
		},
	}, nil
}
