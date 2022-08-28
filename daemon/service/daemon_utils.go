package daemon

import (
	"errors"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

type pathExtractorFn func() string

func (daemon *FileSystemDaemon) extractError(sessionId string, workDir file.File, err error) (*pb.Response, error) {
	target := &fserrors.FileSystemError{}
	if errors.As(err, &target) {
		workDir := daemon.maybeChangeWorkDirectory(sessionId, workDir, target)
		return &pb.Response{
			Error:          proto.String(err.Error()),
			WorkingDirPath: workDir.Info().AbsolutePath(),
		}, nil
	}
	return nil, err
}

func (daemon *FileSystemDaemon) maybeChangeWorkDirectory(sessionId string, workDir file.File, err *fserrors.FileSystemError) file.File {
	if err == fserrors.ErrInvalidWorkingDirectory {
		daemon.sessionStore.ChangeWorkingDirectory(sessionId, daemon.fs.DefaultWorkingDirectory())
		return daemon.fs.DefaultWorkingDirectory()
	}
	return workDir
}

func (daemon *FileSystemDaemon) updateWorkingDirectory(sessionId string, deletedFile file.FileInfo) (file.File, error) {
	workingDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(sessionId)
	if err != nil {
		return nil, err
	}

	if workingDir.Info().AbsolutePath() == deletedFile.AbsolutePath() {
		daemon.sessionStore.ChangeWorkingDirectory(sessionId, daemon.fs.DefaultWorkingDirectory())
		return daemon.fs.DefaultWorkingDirectory(), nil
	}

	return workingDir, nil
}

func (daemon *FileSystemDaemon) getPath(req *pb.Request, pathpathExtractorFn pathExtractorFn) (*fspath.FileSystemPath, error) {
	workingDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(req.SessionId)
	if err != nil {
		return nil, err
	}

	return fspath.NewFileSystemPath(pathpathExtractorFn(), workingDir)
}
