package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

func (daemon *FileSystemDaemon) Move(ctx context.Context, request *pb.MoveRequest) (*pb.MoveResponse, error) {
	srcPath, destPath, workingDir, err := daemon.getSrcPathDestPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.Move(srcPath, destPath, workingDir)

	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.MoveResponse{Error: proto.String(msg)}, nil
	}

	return &pb.MoveResponse{Name: proto.String(file.Name())}, nil
}
