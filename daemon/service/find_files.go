package daemon

import (
	"context"
	"fmt"
	"log"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) FindFiles(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - findFiles request recevied: {%+v}", request.GetSessionId(), request)
	findReq := request.GetFind()
	if findReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return findReq.GetPath() })
	if err != nil {
		log.Printf("%s - findFiles path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	files, err := daemon.fs.FindFiles(findReq.GetName(), path)
	workDir := path.WorkingDir()
	if err != nil {
		log.Printf("%s - findFiles fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	paths := []string{}
	for _, info := range files {
		paths = append(paths, info.AbsolutePath())
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Find{
			Find: &pb.FindFilesResponse{Paths: paths},
		},
	}, nil
}
