package daemon

import (
	"context"
	pb "material/filesystem/pb/proto/session"
)

func (daemon *FileSystemDaemon) NewSession(ctx context.Context, newSessionRequest *pb.NewSessionRequest) (*pb.NewSessionResponse, error) {
	return daemon.sessionStore.AddSession(newSessionRequest, daemon.fs.DefaultWorkingDirectory())
}
