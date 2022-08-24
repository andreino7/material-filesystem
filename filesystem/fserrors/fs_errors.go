package fserrors

import "errors"

var (
	ErrExist                   = &FileSystemError{err: errors.New("file already exists")}
	ErrNotExist                = &FileSystemError{err: errors.New("file does not exist")}
	ErrInvalid                 = &FileSystemError{err: errors.New("invalid argument")}
	ErrInvalidFileType         = &FileSystemError{err: errors.New("file is not a directory")}
	ErrOperationNotSupported   = &FileSystemError{err: errors.New("operation not supported")}
	ErrInvalidWorkingDirectory = &FileSystemError{err: errors.New("invalid working directory")}
	ErrSameFile                = &FileSystemError{err: errors.New("same file")}
)

type FileSystemError struct {
	err error
}

func (e FileSystemError) Error() string { return e.err.Error() }

func (e FileSystemError) Unwrap() error { return e.err }
