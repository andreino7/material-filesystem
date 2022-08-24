package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/session"
)

func (daemon *FileSystemDaemon) DeleteSession(ctx context.Context, deleteSessionRequest *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	return daemon.sessionStore.DeleteSession(deleteSessionRequest)
}
