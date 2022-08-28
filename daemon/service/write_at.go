package daemon

import (
	"context"
	"fmt"
	"log"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) WriteAt(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	log.Printf("%s - writeAt request recevied: {%+v}", request.GetSessionId(), request)

	writeReq := request.GetWriteAt()
	if writeReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	workDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(request.GetSessionId())
	if err != nil {
		log.Printf("%s - writeAt path error: %s", request.GetSessionId(), err.Error())
		return nil, err
	}

	size, err := daemon.fs.WriteAt(writeReq.GetFileDescriptor(), writeReq.GetContent(), int(writeReq.GetPos()))
	if err != nil {
		log.Printf("%s - writeAt fs error: %s", request.GetSessionId(), err.Error())
		return daemon.extractError(request.GetSessionId(), workDir, err)
	}

	return &pb.Response{
		WorkingDirPath: workDir.Info().AbsolutePath(),
		Response: &pb.Response_WriteAt{
			WriteAt: &pb.WriteAtResponse{
				NBytes: int32(size),
			},
		},
	}, nil
}
