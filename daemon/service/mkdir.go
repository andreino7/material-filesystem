package daemon

import (
	"context"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

// TODO: add mocking and unit tests
func (daemon *FileSystemDaemon) Mkdir(ctx context.Context, request *pb.MkdirRequest) (*pb.MkdirResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	var file file.File
	if request.GetRecursive() {
		file, err = daemon.fs.MkdirAll(path, workingDir)
	} else {
		file, err = daemon.fs.Mkdir(path, workingDir)
	}

	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.MkdirResponse{Error: proto.String(msg)}, nil
	}

	return &pb.MkdirResponse{
		Name: proto.String(file.Info().Name()),
	}, nil
}
