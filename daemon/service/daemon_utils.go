package daemon

import (
	"errors"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

type singlePathRequest interface {
	GetPath() string
	GetSessionId() string
}

type srcAndDestPathRequest interface {
	GetSrcPath() string
	GetDestPath() string
	GetSessionId() string
}

func (daemon *FileSystemDaemon) extractError(sessionId string, err error) (string, error) {
	target := &fserrors.FileSystemError{}
	if errors.As(err, &target) {
		daemon.maybeChangeWorkDirectory(sessionId, target)
		return err.Error(), nil
	}
	return "", err
}

func (daemon *FileSystemDaemon) maybeChangeWorkDirectory(sessionId string, err *fserrors.FileSystemError) {
	if err == fserrors.ErrInvalidWorkingDirectory {
		daemon.sessionStore.ChangeWorkingDirectory(sessionId, daemon.fs.DefaultWorkingDirectory())
	}
}

func (daemon *FileSystemDaemon) getPathAndWorkDir(req singlePathRequest) (*fspath.FileSystemPath, file.File, error) {
	workingDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(req)
	if err != nil {
		return nil, nil, err
	}

	path := fspath.NewFileSystemPath(req.GetPath())
	return path, workingDir, nil
}

func (daemon *FileSystemDaemon) getSrcPathDestPathAndWorkDir(req srcAndDestPathRequest) (*fspath.FileSystemPath, *fspath.FileSystemPath, file.File, error) {
	workingDir, err := daemon.sessionStore.GetWorkingDirectoryForSession(req)
	if err != nil {
		return nil, nil, nil, err
	}

	srcPath := fspath.NewFileSystemPath(req.GetSrcPath())
	destPath := fspath.NewFileSystemPath(req.GetDestPath())
	return srcPath, destPath, workingDir, nil
}
