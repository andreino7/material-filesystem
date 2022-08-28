package daemon

import (
	"context"
	"fmt"
	"log"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) AppendAll(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - appendAll request recevied: {%+v}", request.GetSessionId(), request)
	appendReq := request.GetAppendAll()
	if appendReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return appendReq.GetPath() })
	if err != nil {
		log.Printf("%s - appendAll path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := path.WorkingDir()
	err = daemon.fs.AppendAll(path, appendReq.GetContent())
	if err != nil {
		log.Printf("%s - appendAll fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_AppendAll{
			AppendAll: &pb.AppendAllResponse{},
		},
	}, nil
}
