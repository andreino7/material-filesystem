package daemon

import (
	"context"
	"fmt"
	"log"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) ListFiles(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - listFiles request recevied: {%+v}", request.GetSessionId(), request)
	lsReq := request.GetList()
	if lsReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return lsReq.GetPath() })
	if err != nil {
		log.Printf("%s - listFiles path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := path.WorkingDir()
	files, err := daemon.fs.ListFiles(path)
	if err != nil {
		log.Printf("%s - listFiles fs error: %s", request.GetSessionId(), err.Error())
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
