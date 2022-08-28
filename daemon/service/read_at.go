package daemon

import (
	"context"
	"fmt"

	pb "material/filesystem/pb/proto/fsservice"
)

func (daemon *FileSystemDaemon) ReadAt(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	readReq := request.GetReadAt()
	if readReq == nil {
		return nil, fmt.Errorf("invalid request")
	}

	content, err := daemon.fs.ReadAt(readReq.GetFileDescriptor(), int(readReq.GetStartPos()), int(readReq.GetEndPos()))
	if err != nil {
		return daemon.extractError(request.GetSessionId(), err)
	}

	return &pb.Response{
		Response: &pb.Response_ReadAt{
			ReadAt: &pb.ReadAtResponse{
				Content: content,
			},
		},
	}, nil
}
