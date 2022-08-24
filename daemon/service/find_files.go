package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

func (daemon *FileSystemDaemon) FindFiles(ctx context.Context, request *pb.FindFilesRequest) (*pb.FindFilesResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	// TODO: validate name
	files, err := daemon.fs.FindFiles(request.GetName(), path, workingDir)
	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.FindFilesResponse{Error: proto.String(msg)}, nil
	}

	paths := []string{}
	for _, info := range files {
		paths = append(paths, info.AbsolutePath())
	}

	return &pb.FindFilesResponse{
		Paths: paths,
	}, nil
}
