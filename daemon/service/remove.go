package daemon

import (
	"context"
	"fmt"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Remove(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	rmReq := request.GetRemove()
	if rmReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return rmReq.GetPath() })
	if err != nil {
		return nil, err
	}

	var file file.FileInfo
	if rmReq.GetRecursive() {
		file, err = daemon.fs.RemoveAll(path)
	} else {
		file, err = daemon.fs.Remove(path)
	}

	workDir := path.WorkingDir()
	if err != nil {
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	workDir, err = daemon.updateWorkingDirectory(request.GetSessionId(), file)
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Remove{
			Remove: &pb.RemoveResponse{},
		},
	}, nil
}
