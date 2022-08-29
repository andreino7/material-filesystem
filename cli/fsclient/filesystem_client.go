package fsclient

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"material/filesystem/pb/proto/fsservice"
	"material/filesystem/pb/proto/session"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Session FileSystemSession

type FileSystemSession struct {
	sessionId      string
	workingDirPath string
	conn           *grpc.ClientConn
	fsservice.FileSystemServiceClient
	session.SessionServiceClient
}

type grpcFileSystemCall func(ctx context.Context, req *fsservice.Request, opts ...grpc.CallOption) (*fsservice.Response, error)
type onSuccessFn func(*fsservice.Response)

func (f *FileSystemSession) DoRequest(req *fsservice.Request, call grpcFileSystemCall, onSuonSuccessFn onSuccessFn) {
	// TODO: make timeout configurable
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req.SessionId = f.sessionId
	resp, err := call(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("unexpected error: %v", err))
		return
	}

	defer f.updateWokingDirectory(resp)
	if resp.GetError() != "" {
		fmt.Println(resp.GetError())
	} else {
		onSuonSuccessFn(resp)
	}
}

func (f *FileSystemSession) updateWokingDirectory(resp *fsservice.Response) {
	f.workingDirPath = resp.GetWorkingDirPath()
}

func (f *FileSystemSession) WorkingDirPath() string {
	return f.workingDirPath
}

func (f *FileSystemSession) WorkingDirName() string {
	return filepath.Base(f.workingDirPath)
}

func (f *FileSystemSession) SetWorkingDirPath(workingDirPath string) {
	f.workingDirPath = workingDirPath
}

func Initialize(port string) error {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial("localhost:"+port, opts...)
	if err != nil {
		return fmt.Errorf("fail to dial: %v", err)
	}
	fsClient := fsservice.NewFileSystemServiceClient(conn)
	sessionClient := session.NewSessionServiceClient(conn)

	ctx := context.Background()
	resp, err := sessionClient.NewSession(ctx, &session.NewSessionRequest{})
	if err != nil {
		conn.Close()
		return fmt.Errorf("error creating new session: %w", err)
	}

	Session = FileSystemSession{
		sessionId:               resp.GetSessionId(),
		workingDirPath:          resp.GetWorkingDirectoryPath(),
		conn:                    conn,
		FileSystemServiceClient: fsClient,
		SessionServiceClient:    sessionClient,
	}

	return nil
}

func Close() {
	ctx := context.Background()
	if Session.sessionId != "" {
		Session.DeleteSession(ctx, &session.DeleteSessionRequest{SessionId: Session.sessionId})
	}

	if Session.conn != nil {
		Session.conn.Close()
	}
}
