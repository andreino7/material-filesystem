package daemon

import (
	"fmt"
	"log"
	"material/filesystem/daemon/session"
	"material/filesystem/filesystem"
	pb "material/filesystem/pb/proto/session"

	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type FileSystemDaemon struct {
	fs           filesystem.FileSystem
	sessionStore session.SessionStore
	pb.UnimplementedSessionServiceServer
}

func NewFileSystemDaemon(fsType filesystem.FileSystemType) (*FileSystemDaemon, error) {
	fs, err := filesystem.NewFileSystem(fsType)
	if err != nil {
		return nil, fmt.Errorf("error creating filesystem: %w", err)
	}
	return &FileSystemDaemon{
		fs:                                fs,
		sessionStore:                      session.NewSessionStore(),
		UnimplementedSessionServiceServer: pb.UnimplementedSessionServiceServer{},
	}, nil
}

func (daemon *FileSystemDaemon) Run() error {
	// TODO: make port configurable
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 3333))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSessionServiceServer(grpcServer, daemon)
	reflection.Register(grpcServer)
	fmt.Println("Daemon listening on port: 3333")
	grpcServer.Serve(lis)
	return nil
}
