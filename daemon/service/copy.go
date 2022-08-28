package daemon

import (
	"context"
	"fmt"
	"log"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Copy(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - copy request recevied: {%+v}", request.GetSessionId(), request)
	cpReq := request.GetCopy()
	if cpReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return cpReq.GetSrcPath() })
	if err != nil {
		log.Printf("%s - copy srcPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return cpReq.GetDestPath() })
	if err != nil {
		log.Printf("%s - copy destPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := srcPath.WorkingDir()
	file, err := daemon.fs.Copy(srcPath, destPath)
	if err != nil {
		log.Printf("%s - copy daemon error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Copy{
			Copy: &pb.CopyResponse{Name: file.Name()},
		},
	}, nil
}
