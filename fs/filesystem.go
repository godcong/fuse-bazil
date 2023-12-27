package fs

import (
	"context"

	"bazil.org/fuse"
)

type FileSystem interface {
	String() string
	// Lookup is called by the kernel when the VFS wants to know
	// about a file inside a directory. Many lookup calls can
	// occur in parallel, but only one call happens for each (dir,
	// name) pair.
	Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (Node, error)

	// Forget is called when the kernel discards entries from its
	// dentry cache. This happens on unmount, and when the kernel
	// is short on memory. Since it is not guaranteed to occur at
	// any moment, and since there is no return value, Forget
	// should not do I/O, as there is no channel to report back
	// I/O errors.
	Forget(nodeid, nlookup uint64)

	// GetAttr obtains the standard metadata for the receiver.
	// It should store that metadata in resp.
	//
	// If this method is not implemented, the attributes will be
	// generated based on Attr(), with zero values filled in.
	GetAttr(ctx context.Context, req *fuse.GetAttrRequest, resp *fuse.GetAttrResponse) error
	// SetAttr sets the standard metadata for the receiver.
	//
	// Note, this is also used to communicate changes in the size of
	// the file, outside Writes.
	//
	// req.Valid is a bitmask of what fields are actually being set.
	// For example, the method should not change the mode of the file
	// unless req.Valid.Mode() is true.
	SetAttr(ctx context.Context, req *fuse.SetAttrRequest, resp *fuse.SetAttrResponse) error
	// Modifying structure.
	Mknod(ctx context.Context, req *fuse.MknodRequest) (Node, error)
	Mkdir(ctx context.Context, req *fuse.MkdirRequest) (Node, error)
	Unlink(ctx context.Context, req *fuse.UnlinkRequest) error
	Rmdir(ctx context.Context, req *fuse.RmdirRequest) error
	Rename(ctx context.Context, req *fuse.RenameRequest, new Node) error
	Link(ctx context.Context, req *fuse.LinkRequest, old Node) (Node, error)

	Symlink(ctx context.Context, req *fuse.SymlinkRequest) (Node, error)
	Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error)
	Access(ctx context.Context, req *fuse.AccessRequest) error

	// Extended attributes.

	// GetXAttr reads an extended attribute, and should return the
	// number of bytes. If the buffer is too small, return ERANGE,
	// with the required buffer size.
	GetXAttr(ctx context.Context, req *fuse.GetXAttrRequest, resp *fuse.GetXAttrResponse) error

	// ListXAttr lists extended attributes as '\0' delimited byte
	// slice, and return the number of bytes. If the buffer is too
	// small, return ERANGE, with the required buffer size.
	ListXAttr(ctx context.Context, req *fuse.ListXAttrRequest, resp *fuse.ListXAttrResponse) error

	// SetXAttr sets an extended attribute with the given name and
	// value for the node.
	SetXAttr(ctx context.Context, req *fuse.SetXAttrRequest) error
	// RemoveXAttr removes an extended attribute.
	// RemoveXAttr removes an extended attribute for the name.
	//
	// If there is no XAttr by that name, returns fuse.ErrNoXAttr.
	RemoveXAttr(ctx context.Context, req *fuse.RemoveXAttrRequest) error

	// Create creates a new directory entry in the receiver, which
	// must be a directory.
	Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (Node, Handle, error)
	// Open opens the receiver. After a successful open, a client
	// process has a file descriptor referring to this Handle.
	//
	// Open can be also called on non-files. For example,
	// directories are Opened for ReadDir or fchdir(2).
	//
	// If this method is not implemented, the open will always
	// succeed, and the Node itself will be used as the Handle.
	//
	// XXX note about access.  XXX OpenFlags.
	Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (Handle, error)
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

	//Lseek(ctx context.Context, in *LseekIn, out *LseekOut) Status //todo: implement
	//GetLk(ctx context.Context, input *LkIn, out *LkOut) (code Status) //todo: implement
	//SetLk(ctx context.Context, input *LkIn) (code Status) //todo: implement
	//SetLkw(ctx context.Context, input *LkIn) (code Status) //todo: implement

	Release(ctx context.Context, req *fuse.ReleaseRequest) error
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
	//CopyFileRange(ctx context.Context, input *CopyFileRangeIn) (written uint32, code Status) //todo: implement

	// Flush is called each time the file or directory is closed.
	// Because there can be multiple file descriptors referring to a
	// single opened file, Flush can be called multiple times.
	Flush(ctx context.Context, req *fuse.FlushRequest) error
	// Fsync is a signal to ensure writes to the Inode are flushed
	// to stable storage.
	Fsync(ctx context.Context, req *fuse.FsyncRequest) error
	//Fallocate(ctx context.Context, input *FallocateIn) (code Status)  //todo: implement

	// Directory handling
	//OpenDir(ctx context.Context, input *OpenIn, out *OpenOut) (status Status)  //todo: implement
	//ReadDir(ctx context.Context, input *ReadIn, out *DirEntryList) Status  //todo: implement
	//ReadDirPlus(ctx context.Context, input *ReadIn, out *DirEntryList) Status  //todo: implement
	//ReleaseDir(input *ReleaseIn)  //todo: implement
	//FsyncDir(ctx context.Context, input *FsyncIn) (code Status)  //todo: implement

	//StatFs(ctx context.Context, input *InHeader, out *StatfsOut) (code Status)  //todo: implement

	// This is called on processing the first request. The
	// filesystem implementation can use the server argument to
	// talk back to the kernel (through notify methods).
	//Init(*Server)  //todo: implement
}
