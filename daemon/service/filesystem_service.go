package daemon

import (
	"context"
	"errors"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

type singlePathRequest interface {
	GetPath() string
	GetSessionId() string
}

// TODO: add mocking and unit tests
func (daemon *FileSystemDaemon) Mkdir(ctx context.Context, request *pb.MkdirRequest) (*pb.MkdirResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	var file file.File
	if request.GetRecursive() {
		file, err = daemon.fs.MkdirAll(path, workingDir)
	} else {
		file, err = daemon.fs.Mkdir(path, workingDir)
	}

	if err != nil {
		msg, err := extractError(err)
		if err != nil {
			return nil, err
		}
		return &pb.MkdirResponse{Error: proto.String(msg)}, nil
	}

	return &pb.MkdirResponse{
		Name: proto.String(file.Info().Name()),
	}, nil
}

func (daemon *FileSystemDaemon) CreateRegularFile(ctx context.Context, request *pb.CreateRegularFileRequest) (*pb.CreateRegularFileResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.CreateRegularFile(path, workingDir)
	if err != nil {
		msg, err := extractError(err)
		if err != nil {
			return nil, err
		}
		return &pb.CreateRegularFileResponse{Error: proto.String(msg)}, nil
	}

	return &pb.CreateRegularFileResponse{
		Name: proto.String(file.Info().Name()),
	}, nil
}

func (daemon *FileSystemDaemon) ChangeWorkingDirectory(ctx context.Context, request *pb.ChangeWorkingDirectoryRequest) (*pb.ChangeWorkingDirectoryResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	file, err := daemon.fs.GetDirectory(path, workingDir)
	if err != nil {
		msg, err := extractError(err)
		if err != nil {
			return nil, err
		}
		return &pb.ChangeWorkingDirectoryResponse{Error: proto.String(msg)}, nil
	}

	err = daemon.sessionStore.ChangeWorkingDirectory(request.GetSessionId(), file)
	if err != nil {
		return nil, err
	}

	return &pb.ChangeWorkingDirectoryResponse{
		Name: proto.String(file.Info().Name()),
	}, nil
}

func (daemon *FileSystemDaemon) Remove(ctx context.Context, request *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	var file file.FileInfo
	if request.GetRecursive() {
		file, err = daemon.fs.RemoveAll(path, workingDir)
	} else {
		file, err = daemon.fs.Remove(path, workingDir)
	}

	if err != nil {
		msg, err := extractError(err)
		if err != nil {
			return nil, err
		}
		return &pb.RemoveResponse{Error: proto.String(msg)}, nil
	}

	if workingDir.Info().AbsolutePath() == file.AbsolutePath() {
		err = daemon.sessionStore.ChangeWorkingDirectory(request.GetSessionId(), daemon.fs.DefaultWorkingDirectory())
		if err != nil {
			return nil, err
		}
	}

	// TODO: add working dir to resp
	return &pb.RemoveResponse{}, nil
}

func (daemon *FileSystemDaemon) FindFiles(ctx context.Context, request *pb.FindFilesRequest) (*pb.FindFilesResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	// TODO: validate name
	files, err := daemon.fs.FindFiles(request.GetName(), path, workingDir)
	if err != nil {
		msg, err := extractError(err)
		if err != nil {
			return nil, err
		}
		return &pb.FindFilesResponse{Error: proto.String(msg)}, nil
	}

	paths := []string{}
	for _, info := range files {
		paths = append(paths, info.AbsolutePath())
	}

	return &pb.FindFilesResponse{
		Paths: paths,
	}, nil
}

func (daemon *FileSystemDaemon) ListFiles(ctx context.Context, request *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	path, workingDir, err := daemon.getPathAndWorkDir(request)
	if err != nil {
		return nil, err
	}

	files, err := daemon.fs.ListFiles(path, workingDir)
	if err != nil {
		msg, err := extractError(err)
		if err != nil {
			return nil, err
		}
		return &pb.ListFilesResponse{Error: proto.String(msg)}, nil
	}

	names := []string{}
	for _, info := range files {
		names = append(names, info.Name())
	}

	return &pb.ListFilesResponse{
		Names: names,
	}, nil
}

func extractError(err error) (string, error) {
	target := &fserrors.FileSystemError{}
	if errors.As(err, &target) {
		return err.Error(), nil
	}
	return "", err
}

func (daemon *FileSystemDaemon) getPathAndWorkDir(req singlePathRequest) (*fspath.FileSystemPath, file.File, error) {
	workingDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(req)
	if err != nil {
		return nil, nil, err
	}

	path := fspath.NewFileSystemPath(req.GetPath())
	return path, workingDir, nil
}
