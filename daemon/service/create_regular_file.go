package daemon

import (
	"context"
	"fmt"
	"log"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) CreateRegularFile(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - createRegularFile request recevied: {%+v}", request.GetSessionId(), request)
	createReq := request.GetCreateRegularFile()
	if createReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return createReq.GetPath() })
	if err != nil {
		log.Printf("%s - createRegularFile path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	file, err := daemon.fs.CreateRegularFile(path)
	workDir := path.WorkingDir()
	if err != nil {
		log.Printf("%s - createRegularFile fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_CreateRegularFile{
			CreateRegularFile: &pb.CreateRegularFileResponse{Name: file.Info().Name()},
		},
	}, nil
}
