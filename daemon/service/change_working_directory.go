package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

func (daemon *FileSystemDaemon) ChangeWorkingDirectory(ctx context.Context, request *pb.ChangeWorkingDirectoryRequest) (*pb.ChangeWorkingDirectoryResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.GetDirectory(path, workingDir)
	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.ChangeWorkingDirectoryResponse{Error: proto.String(msg)}, nil
	}

	err = daemon.sessionStore.ChangeWorkingDirectory(request.GetSessionId(), file)
	if err != nil {
		return nil, err
	}

	return &pb.ChangeWorkingDirectoryResponse{
		Name: proto.String(file.Info().Name()),
	}, nil
}
