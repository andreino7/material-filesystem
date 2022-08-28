package daemon

import (
	"context"
	"fmt"
	"log"
	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) CreateSymbolicLink(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - createSymbolicLink request recevied: {%+v}", request.GetSessionId(), request)
	createReq := request.GetSymLink()
	if createReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	srcPath, err := daemon.getPath(request, func() string { return createReq.GetSrcPath() })
	if err != nil {
		log.Printf("%s - createSymbolicLink srcPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	destPath, err := daemon.getPath(request, func() string { return createReq.GetDestPath() })
	if err != nil {
		log.Printf("%s - createSymbolicLink destPath error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	workDir := srcPath.WorkingDir()
	file, err := daemon.fs.CreateSymbolicLink(srcPath, destPath)
	if err != nil {
		log.Printf("%s - createSymbolicLink fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_SymLink{
			SymLink: &pb.CreateSymbolicLinkResponse{Name: file.Name()},
		},
	}, nil
}
