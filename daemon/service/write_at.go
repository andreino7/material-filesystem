package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) WriteAt(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	writeReq := request.GetWriteAt()
	if writeReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	size, err := daemon.fs.WriteAt(writeReq.GetFileDescriptor(), writeReq.GetContent(), int(writeReq.GetPos()))
	if err != nil {
		return daemon.extractError(request.GetSessionId(), err)
	}

	return &pb.Response{
		Response: &pb.Response_WriteAt{
			WriteAt: &pb.WriteAtResponse{
				NBytes: int32(size),
			},
		},
	}, nil
}
