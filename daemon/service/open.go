package daemon

import (
	"context"
	"fmt"
	"log"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Open(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - open request recevied: {%+v}", request.GetSessionId(), request)
	openReq := request.GetOpen()
	if openReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return openReq.GetPath() })
	if err != nil {
		log.Printf("%s - open path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := path.WorkingDir()
	fd, err := daemon.fs.Open(path)
	if err != nil {
		log.Printf("%s - open fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Open{
			Open: &pb.OpenResponse{FileDescriptor: fd},
		},
	}, nil
}
