package daemon

import (
	"context"
	"fmt"
	"log"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) ReadAt(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - readAt request recevied: {%+v}", request.GetSessionId(), request)
	readReq := request.GetReadAt()
	if readReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	workDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(request.GetSessionId())
	if err != nil {
		log.Printf("%s - readAt path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	size := readReq.GetEndPos() - readReq.GetStartPos()
	buff := make([]byte, size)
	_, err = daemon.fs.ReadAt(readReq.GetFileDescriptor(), buff, int(readReq.GetStartPos()))
	if err != nil {
		log.Printf("%s - readAt fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_ReadAt{
			ReadAt: &pb.ReadAtResponse{
				Content: buff,
			},
		},
	}, nil
}
