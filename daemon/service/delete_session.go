package daemon

import (
	"context"
	"log"
	pb "material/filesystem/pb/proto/session"
)

func (daemon *FileSystemDaemon) DeleteSession(ctx context.Context, deleteSessionRequest *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	log.Printf("%s - deleteSession request recevied: {%+v}", deleteSessionRequest.GetSessionId(), deleteSessionRequest)
	return daemon.sessionStore.DeleteSession(deleteSessionRequest)
}
