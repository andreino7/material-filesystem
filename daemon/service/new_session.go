package daemon

import (
	"context"
	"log"
	pb "material/filesystem/pb/proto/session"
)

func (daemon *FileSystemDaemon) NewSession(ctx context.Context, newSessionRequest *pb.NewSessionRequest) (*pb.NewSessionResponse, error) {
	log.Printf("newSession request recevied: {%+v}", newSessionRequest)
	return daemon.sessionStore.AddSession(newSessionRequest, daemon.fs.DefaultWorkingDirectory())
}
