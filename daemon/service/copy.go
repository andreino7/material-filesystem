package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

func (daemon *FileSystemDaemon) Copy(ctx context.Context, request *pb.CopyRequest) (*pb.CopyResponse, error) {
	srcPath, destPath, workingDir, err := daemon.getSrcPathDestPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.Copy(srcPath, destPath, workingDir)

	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.CopyResponse{Error: proto.String(msg)}, nil
	}

	return &pb.CopyResponse{Name: proto.String(file.Name())}, nil
}
