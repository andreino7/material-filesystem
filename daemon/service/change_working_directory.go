package daemon

import (
	"context"
	"fmt"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) ChangeWorkingDirectory(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	cdReq := request.GetChangeWorkingDirectory()
	if cdReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return cdReq.GetPath() })
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.GetDirectory(path)
	if err != nil {
		return daemon.extractError(request.GetSessionId(), path.WorkingDir(), err)
	}

	err = daemon.sessionStore.ChangeWorkingDirectory(request.GetSessionId(), file)
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		WorkingDirPath: file.Info().AbsolutePath(),
		Response: &pb.Response_ChangeWorkingDirectory{
			ChangeWorkingDirectory: &pb.ChangeWorkingDirectoryResponse{Path: file.Info().AbsolutePath()},
		},
	}, nil
}
