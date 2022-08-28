package daemon

import (
	"context"
	"fmt"
	"log"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/fsservice"
)

// TODO: add mocking and unit tests
func (daemon *FileSystemDaemon) Mkdir(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - mkdir request recevied: {%+v}", request.GetSessionId(), request)
	mkdirReq := request.GetMkdir()
	if mkdirReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return mkdirReq.GetPath() })
	if err != nil {
		log.Printf("%s - mkdir daemon error: %s", request.GetSessionId(), err.Error())
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
		log.Printf("%s - mkdir fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Mkdir{
			Mkdir: &pb.MkdirResponse{Name: file.Info().Name()},
		},
	}, nil
}
