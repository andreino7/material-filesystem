package daemon

import (
	"context"

	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

func (daemon *FileSystemDaemon) CreateRegularFile(ctx context.Context, request *pb.CreateRegularFileRequest) (*pb.CreateRegularFileResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.CreateRegularFile(path, workingDir)
	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.CreateRegularFileResponse{Error: proto.String(msg)}, nil
	}

	return &pb.CreateRegularFileResponse{
		Name: proto.String(file.Info().Name()),
	}, nil
}
