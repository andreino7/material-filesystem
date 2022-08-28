package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) CreateRegularFile(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	createReq := request.GetCreateRegularFile()
	if createReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return createReq.GetPath() })
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.CreateRegularFile(path)
	workDir := path.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_CreateRegularFile{
			CreateRegularFile: &pb.CreateRegularFileResponse{Name: file.Info().Name()},
		},
	}, nil
}
