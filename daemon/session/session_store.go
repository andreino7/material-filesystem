package session

import (
	"fmt"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/session"

	"github.com/google/uuid"
)

type session struct {
	sessionId        string
	workingDirectory file.File
}

// TODO: implement cleanup for inactive sessions
type SessionStore map[string]*session

func NewSessionStore() SessionStore {
	return SessionStore{}
}

func (store SessionStore) AddSession(request *pb.NewSessionRequest, workingDirectory file.File) (*pb.NewSessionResponse, error) {
	session := &session{
		sessionId:        uuid.NewString(),
		workingDirectory: workingDirectory,
	}
	// This should never happen
	if _, found := store[session.sessionId]; found {
		return nil, fmt.Errorf("invalid session id")
	}
	store[session.sessionId] = session
	return &pb.NewSessionResponse{
		SessionId:            session.sessionId,
		WorkingDirectoryName: session.workingDirectory.Info().Name(),
	}, nil
}

func (store SessionStore) DeleteSession(request *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	if _, found := store[request.SessionId]; !found {
		return nil, fmt.Errorf("session not found")
	}
	delete(store, request.SessionId)
	return &pb.DeleteSessionResponse{
		SessionId: request.SessionId,
	}, nil
}
