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

	workDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(request.GetSessionId())
	if err != nil {
		return nil, err
	}

	daemon.fs.Close(closeReq.GetFileDescriptor())

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Close{
			Close: &pb.CloseResponse{},
		},
	}, nil
}
