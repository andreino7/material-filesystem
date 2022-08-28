package daemon

import (
	"context"
	"fmt"
	"log"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) Remove(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - remove request recevied: {%+v}", request.GetSessionId(), request)

	rmReq := request.GetRemove()
	if rmReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	path, err := daemon.getPath(request, func() string { return rmReq.GetPath() })
	if err != nil {
		log.Printf("%s - remove path error: %s", request.GetSessionId(), err.Error())
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
		log.Printf("%s - remove fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	workDir, err = daemon.updateWorkingDirectory(request.GetSessionId(), file)
	if err != nil {
		log.Printf("%s - remove update working directory error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_Remove{
			Remove: &pb.RemoveResponse{},
		},
	}, nil
}
