package daemon

import (
	"context"
	"fmt"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Copy(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	cpReq := request.GetCopy()
	if cpReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return cpReq.GetSrcPath() })
	if err != nil {
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return cpReq.GetDestPath() })
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.Copy(srcPath, destPath)

	if err != nil {
		return daemon.extractError(request.GetSessionId(), err)
	}

	return &pb.Response{
		Response: &pb.Response_Copy{
			Copy: &pb.CopyResponse{Name: file.Name()},
		},
	}, nil
}
