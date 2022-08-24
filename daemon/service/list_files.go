package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

func (daemon *FileSystemDaemon) ListFiles(ctx context.Context, request *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	files, err := daemon.fs.ListFiles(path, workingDir)
	if err != nil {
		msg, err := daemon.extractError(request.GetSessionId(), err)
		if err != nil {
			return nil, err
		}
		return &pb.ListFilesResponse{Error: proto.String(msg)}, nil
	}

	names := []string{}
	for _, info := range files {
		names = append(names, info.Name())
	}

	return &pb.ListFilesResponse{
		Names: names,
	}, nil
}
