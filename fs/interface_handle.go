package fs

import (
	"context"

	"bazil.org/fuse"
)

type HandleFlusher interface {
	// Flush is called each time the file or directory is closed.
	// Because there can be multiple file descriptors referring to a
	// single opened file, Flush can be called multiple times.
	Flush(ctx context.Context, req *fuse.FlushRequest) error
}

type HandleReadAller interface {
	ReadAll(ctx context.Context) ([]byte, error)
}

type HandleReadDirAller interface {
	ReadDirAll(ctx context.Context) ([]fuse.Dirent, error)
}

type HandleReader interface {
	// Read requests to read data from the handle.
	//
	// Copy the response bytes to the byte slice resp.Data, slicing
	// it shorter when needed.
	//
	// There is a page cache in the kernel that normally submits only
	// page-aligned reads spanning one or more pages. However, you
	// should not rely on this. To see individual requests as
	// submitted by the file system clients, set OpenDirectIO.
	//
	// Note that reads beyond the size of the file as reported by Attr
	// are not even attempted (except in OpenDirectIO mode).
	Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error
}

type HandleWriter interface {
	// Write requests to write data into the handle at the given offset.
	// Store the amount of data written in resp.Size.
	//
	// There is a writeback page cache in the kernel that normally submits
	// only page-aligned writes spanning one or more pages. However,
	// you should not rely on this. To see individual requests as
	// submitted by the file system clients, set OpenDirectIO.
	//
	// Writes that grow the file are expected to update the file size
	// (as seen through Attr). Note that file size changes are
	// communicated also through SetAttr.
	Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error
}

type HandleReleaser interface {
	Release(ctx context.Context, req *fuse.ReleaseRequest) error
}

type HandlePoller interface {
	// Poll checks whether the handle is currently ready for I/O, and
	// may request a wakeup when it is.
	//
	// Poll should always return quickly. Clients waiting for
	// readiness can be woken up by passing the return value of
	// PollRequest.Wakeup to fs.Server.NotifyPollWakeup or
	// fuse.Conn.NotifyPollWakeup.
	//
	// To allow supporting poll for only some of your Nodes/Handles,
	// the default behavior is to report immediate readiness. If your
	// FS does not support polling and you want to minimize needless
	// requests and log noise, implement NodePoller and return
	// syscall.ENOSYS.
	//
	// The Go runtime uses epoll-based I/O whenever possible, even for
	// regular files.
	Poll(ctx context.Context, req *fuse.PollRequest, resp *fuse.PollResponse) error
}

type HandleFsyncer interface {
	Fsync(ctx context.Context, req *fuse.FsyncRequest) error
}
