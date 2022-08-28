# material-filesystem

material-filesystem is an in-memory file system

## Structure
---------------
### filesystem
A file system library. 

At this time the only file system type supported is an in-memory file system.


### daemon
A grpc service acting as a "very basic OS" and using the filesystem's library.

Basically, it just keeps track of every CLI connection and translates requests coming from the CLI into calls to the filesystem library and relay the result back to the CLI that made the request.

### cli
A grpc client implementic a very basic command line interface to interact with the filesystem daemon.

## Getting started
---------------

The latest release is [1.0](). You can find precompiled binaries ready to use. 

Please note that they have been tested only on a Mac Intel, there's no guarantee they would work on other OSs or architecutres. 

If you need to build the binaries please go to the [build](#build) section.

### Daemon

```console
foo@bar:~$ ./fs-daemon

2022/08/28 12:21:09 Initializing file system daemon
2022/08/28 12:21:09 Starting daemon
2022/08/28 12:21:09 Daemon listening on port: 2222
```

### CLI

```console
foo@bar:~$ ./fs-cli

Welcome to material filesystem client!
Type help for help :)
material-filesystem / $
```




## Build
---------------
### Prerequisites
* [Go](https://go.dev/doc/install)
* [Protoc and go plugins](https://grpc.io/docs/languages/go/quickstart/#prerequisites)

### Daemon
```console
foo@bar:~$ make build-daemon
```
This will create a binary in `build/fs-daemon`

### CLI
```console
foo@bar:~$ make build-cli
```
This will create a binary in `build/fs-cli`

## What's supported?
---------------
* Absolute and relative paths
* Changing and getting the current working directory
* Creating and new directories and any parent directory
* Listing a directory content
* Creating empty files
* Removing files or directories
* Appending content to files or write at specific positions.
* Reading all content of a file or only a portion of it.
* Moving (rename) a file or directory to a new location and automatically resolves any name conflict.
* Copying a file or directory to a new location and automatically resolves any name conflict
* Hard links to regular files
* Symbolic links to files and directories
* Files can have multiple readers at the same time


See [filesystem.go](https://github.com/andreino7/material-filesystem/blob/main/filesystem/filesystem.go) for more details or type help in `fs-cli`:
```console
foo@bar:~$ ./fs-cli

Welcome to material filesystem client!
Type help for help :)
material-filesystem / $ help
```

## Future improvements
---------------
* Optimizing file writes at random locations
* Merge [user permissions](https://github.com/andreino7/material-filesystem/tree/users-v2) once stable
* Better logging framework
