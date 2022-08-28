package daemon

import (
	"context"
	"fmt"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) ListFiles(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	lsReq := request.GetList()
	if lsReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return lsReq.GetPath() })
	if err != nil {
		return nil, err
	}

	files, err := daemon.fs.ListFiles(path)
	workDir := path.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	names := []string{}
	for _, info := range files {
		names = append(names, info.Name())
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_List{
			List: &pb.ListFilesResponse{Names: names},
		},
	}, nil
}
