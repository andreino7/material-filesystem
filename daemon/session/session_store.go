package session

import (
	"fmt"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/session"
	"sync"

	"github.com/google/uuid"
)

type session struct {
	sessionId        string
	workingDirectory file.File
}

// TODO: implement cleanup for inactive sessions
type SessionStore struct {
	sessions map[string]*session
	mutex    sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: map[string]*session{},
	}
}

func (store *SessionStore) AddSession(request *pb.NewSessionRequest, workingDirectory file.File) (*pb.NewSessionResponse, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	session := &session{
		sessionId:        uuid.NewString(),
		workingDirectory: workingDirectory,
	}
	// This should never happen
	if _, found := store.sessions[session.sessionId]; found {
		return nil, fmt.Errorf("invalid session id")
	}
	store.sessions[session.sessionId] = session
	return &pb.NewSessionResponse{
		SessionId:            session.sessionId,
		WorkingDirectoryName: session.workingDirectory.Info().Name(),
	}, nil
}

func (store *SessionStore) DeleteSession(request *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, found := store.sessions[request.SessionId]; !found {
		return nil, fmt.Errorf("session not found")
	}
	delete(store.sessions, request.SessionId)
	return &pb.DeleteSessionResponse{
		SessionId: request.SessionId,
	}, nil
}

func (store *SessionStore) GetWorkingDirectoryForSession(sessionId string) (file.File, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if sessionId == "" {
		return nil, fmt.Errorf("invalid session id")
	}

	session, found := store.sessions[sessionId]
	if !found {
		return nil, fmt.Errorf("session not found")
	}

	return session.workingDirectory, nil
}

func (store *SessionStore) ChangeWorkingDirectory(sessionId string, workingDir file.File) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	session, found := store.sessions[sessionId]
	if !found {
		return fmt.Errorf("session not found")
	}
	session.workingDirectory = workingDir
	return nil
}
