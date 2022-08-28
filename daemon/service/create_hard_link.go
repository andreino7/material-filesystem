package daemon

import (
	"context"
	"fmt"
	"log"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) CreateHardLink(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - createHardLink request recevied: {%+v}", request.GetSessionId(), request)
	createReq := request.GetHardLink()
	if createReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return createReq.GetSrcPath() })
	if err != nil {
		log.Printf("%s - createHardLink srcPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return createReq.GetDestPath() })
	if err != nil {
		log.Printf("%s - createHardLink destPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := srcPath.WorkingDir()
	file, err := daemon.fs.CreateHardLink(srcPath, destPath)
	if err != nil {
		log.Printf("%s - createHardLink fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_HardLink{
			HardLink: &pb.CreateHardLinkResponse{Name: file.Name()},
		},
	}, nil
}
