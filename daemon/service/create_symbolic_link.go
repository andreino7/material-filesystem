package daemon

import (
	"context"
	"fmt"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) CreateSymbolicLink(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	createReq := request.GetSymLink()
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

	file, err := daemon.fs.CreateSymbolicLink(srcPath, destPath)

	workDir := srcPath.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_SymLink{
			SymLink: &pb.CreateSymbolicLinkResponse{Name: file.Name()},
		},
	}, nil
}
