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
	ErrTooManyLinks            = &FileSystemError{err: errors.New("too many links")}
	ErrNotOpen                 = &FileSystemError{err: errors.New("file is not open")}
)

type FileSystemError struct {
	err error
}

func (e FileSystemError) Error() string { return e.err.Error() }

func (e FileSystemError) Unwrap() error { return e.err }
