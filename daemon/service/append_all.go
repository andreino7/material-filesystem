package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) ReadAll(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	readReq := request.GetReadAll()
	if readReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return readReq.GetPath() })
	if err != nil {
		return nil, err
	}

	workDir := path.WorkingDir()
	content, err := daemon.fs.ReadAll(path)
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_ReadAll{
			ReadAll: &pb.ReadAllResponse{Content: content},
		},
	}, nil
}
