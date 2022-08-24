package daemon

import (
	"context"
	"material/filesystem/filesystem/file"
	pb "material/filesystem/pb/proto/fsservice"

	"google.golang.org/protobuf/proto"
)

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
		msg, err := daemon.extractError(request.GetSessionId(), err)
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
