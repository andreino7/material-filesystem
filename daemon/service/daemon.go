package daemon

import (
	"fmt"
	"log"
	"material/filesystem/daemon/session"
	"material/filesystem/filesystem"
	pbFs "material/filesystem/pb/proto/fsservice"
	pbSession "material/filesystem/pb/proto/session"

	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type FileSystemDaemon struct {
	fs           filesystem.FileSystem
	sessionStore *session.SessionStore
	pbSession.UnimplementedSessionServiceServer
	pbFs.UnimplementedFileSystemServiceServer
}

func NewFileSystemDaemon(fsType filesystem.FileSystemType) (*FileSystemDaemon, error) {
	fs, err := filesystem.NewFileSystem(fsType)
	if err != nil {
		return nil, fmt.Errorf("error creating filesystem: %w", err)
	}
	return &FileSystemDaemon{
		fs:                                   fs,
		sessionStore:                         session.NewSessionStore(),
		UnimplementedSessionServiceServer:    pbSession.UnimplementedSessionServiceServer{},
		UnimplementedFileSystemServiceServer: pbFs.UnimplementedFileSystemServiceServer{},
	}, nil
}

func (daemon *FileSystemDaemon) Run(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pbSession.RegisterSessionServiceServer(grpcServer, daemon)
	pbFs.RegisterFileSystemServiceServer(grpcServer, daemon)
	reflection.Register(grpcServer)
	log.Printf("Daemon listening on port: %s", port)
	grpcServer.Serve(lis)
	return nil
}
