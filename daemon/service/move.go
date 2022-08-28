package daemon

import (
	"context"
	"fmt"
	"log"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Move(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - move request recevied: {%+v}", request.GetSessionId(), request)
	mvReq := request.GetMove()
	if mvReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return mvReq.GetSrcPath() })
	if err != nil {
		log.Printf("%s - move srcPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return mvReq.GetDestPath() })
	if err != nil {
		log.Printf("%s - move destPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := srcPath.WorkingDir()
	file, err := daemon.fs.Move(srcPath, destPath)
	if err != nil {
		log.Printf("%s - move fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Move{
			Move: &pb.MoveResponse{Name: file.Name()},
		},
	}, nil
}
