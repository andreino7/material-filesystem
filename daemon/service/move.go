package daemon

import (
	"context"
	"fmt"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Move(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	mvReq := request.GetCopy()
	if mvReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return mvReq.GetSrcPath() })
	if err != nil {
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return mvReq.GetDestPath() })
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.Move(srcPath, destPath)

	workDir := srcPath.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Move{
			Move: &pb.MoveResponse{Name: file.Name()},
		},
	}, nil
}
