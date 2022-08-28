package daemon

import (
	"context"
	"fmt"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/fsservice"
)

// TODO: add mocking and unit tests
func (daemon *FileSystemDaemon) Mkdir(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	mkdirReq := request.GetMkdir()
	if mkdirReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return mkdirReq.GetPath() })
	if err != nil {
		return nil, err
	}

	var file file.File
	if mkdirReq.GetRecursive() {
		file, err = daemon.fs.MkdirAll(path)
	} else {
		file, err = daemon.fs.Mkdir(path)
	}

	workDir := path.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Mkdir{
			Mkdir: &pb.MkdirResponse{Name: file.Info().Name()},
		},
	}, nil
}
