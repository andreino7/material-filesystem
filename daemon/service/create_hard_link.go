package daemon

import (
	"context"
	"fmt"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) CreateHardLink(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	createReq := request.GetHardLink()
	if createReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return createReq.GetSrcPath() })
	if err != nil {
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return createReq.GetDestPath() })
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.CreateHardLink(srcPath, destPath)

	workDir := srcPath.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_HardLink{
			HardLink: &pb.CreateHardLinkResponse{Name: file.Name()},
		},
	}, nil
}
