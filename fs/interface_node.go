package fs

import (
	"context"
	"time"

	"bazil.org/fuse"
)

// A Node is the interface required of a file or directory.
// See the documentation for type FS for general information
// pertaining to all methods.
//
// A Node must be usable as a map key, that is, it cannot be a
// function, map or slice.
//
// Other FUSE requests can be handled by implementing methods from the
// Node* interfaces, for example NodeOpener.
//
// Methods returning Node should take care to return the same Node
// when the result is logically the same instance. Without this, each
// Node will get a new NodeID, causing spurious cache invalidations,
// extra lookups and aliasing anomalies. This may not matter for a
// simple, read-only filesystem.
type Node interface {
	// Attr fills attr with the standard metadata for the node.
	//
	// Fields with reasonable defaults are prepopulated. For example,
	// all times are set to a fixed moment when the program started.
	//
	// If Inode is left as 0, a dynamic inode number is chosen.
	//
	// The result may be cached for the duration set in Valid.
	// Attr(ctx context.Context, attr *fuse.Attr) error
}

type NodeGetAttrer interface {
	// GetAttr obtains the standard metadata for the receiver.
	// It should store that metadata in resp.
	//
	// If this method is not implemented, the attributes will be
	// generated based on Attr(), with zero values filled in.
	GetAttr(ctx context.Context, req *fuse.GetAttrRequest, resp *fuse.GetAttrResponse) error
}

type NodeSetAttrer interface {
	// SetAttr sets the standard metadata for the receiver.
	//
	// Note, this is also used to communicate changes in the size of
	// the file, outside of Writes.
	//
	// req.Valid is a bitmask of what fields are actually being set.
	// For example, the method should not change the mode of the file
	// unless req.Valid.Mode() is true.
	SetAttr(ctx context.Context, req *fuse.SetAttrRequest, resp *fuse.SetAttrResponse) error
}

type NodeSymlinker interface {
	// Symlink creates a new symbolic link in the receiver, which must be a directory.
	//
	// TODO is the above true about directories?
	Symlink(ctx context.Context, req *fuse.SymlinkRequest) (Node, error)
}

// This optional request will be called only for symbolic link nodes.
type NodeReadlinker interface {
	// Readlink reads a symbolic link.
	Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error)
}

type NodeLinker interface {
	// Link creates a new directory entry in the receiver based on an
	// existing Node. Receiver must be a directory.
	Link(ctx context.Context, req *fuse.LinkRequest, old Node) (Node, error)
}

type NodeRmdirer interface {
	// Rmdir removes the entry with the given name from
	// the receiver, which must be a directory.  The entry to be removed
	// may correspond to a file (unlink) or to a directory (rmdir).
	Rmdir(ctx context.Context, req *fuse.RmdirRequest) error
}

type NodeUnlinker interface {
	// Unlink removes the entry with the given name from
	// the receiver, which must be a directory.  The entry to be removed
	// may correspond to a file (unlink) or to a directory (rmdir).
	Unlink(ctx context.Context, req *fuse.UnlinkRequest) error
}

type NodeAccesser interface {
	// Access checks whether the calling context has permission for
	// the given operations on the receiver. If so, Access should
	// return nil. If not, Access should return EPERM.
	//
	// Note that this call affects the result of the access(2) system
	// call but not the open(2) system call. If Access is not
	// implemented, the Node behaves as if it always returns nil
	// (permission granted), relying on checks in Open instead.
	Access(ctx context.Context, req *fuse.AccessRequest) error
}

type NodeStringLookuper interface {
	// Lookup looks up a specific entry in the receiver,
	// which must be a directory.  Lookup should return a Node
	// corresponding to the entry.  If the name does not exist in
	// the directory, Lookup should return ENOENT.
	//
	// Lookup need not to handle the names "." and "..".
	Lookup(ctx context.Context, name string) (Node, error)
}

type NodeFsyncer interface {
	// Fsync is a signal to ensure writes to the Inode are flushed
	// to stable storage.
	Fsync(ctx context.Context, req *fuse.FsyncRequest) error
}

type NodeRequestLookuper interface {
	// Lookup looks up a specific entry in the receiver.
	// See NodeStringLookuper for more.
	Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (Node, error)
}

type NodeMkdirer interface {
	Mkdir(ctx context.Context, req *fuse.MkdirRequest) (Node, error)
}

type NodeOpener interface {
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
}

type NodeCreater interface {
	// Create creates a new directory entry in the receiver, which
	// must be a directory.
	Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (Node, Handle, error)
}

type NodeForgetter interface {
	// Forget about this node. This node will not receive further
	// method calls.
	//
	// Forget is not necessarily seen on unmount, as all nodes are
	// implicitly forgotten as part of the unmount.
	Forget()
}

type NodeRenamer interface {
	Rename(ctx context.Context, req *fuse.RenameRequest, newDir Node) error
}

type NodeMknoder interface {
	Mknod(ctx context.Context, req *fuse.MknodRequest) (Node, error)
}

type NodeGetXAttrer interface {
	// GetXAttr gets an extended attribute by the given name from the
	// node.
	//
	// If there is no XAttr by that name, returns fuse.ErrNoXAttr.
	GetXAttr(ctx context.Context, req *fuse.GetXAttrRequest, resp *fuse.GetXAttrResponse) error
}

type NodeListXAttrer interface {
	// ListXAttr lists the extended attributes recorded for the node.
	ListXAttr(ctx context.Context, req *fuse.ListXAttrRequest, resp *fuse.ListXAttrResponse) error
}

type NodeSetXAttrer interface {
	// SetXAttr sets an extended attribute with the given name and
	// value for the node.
	SetXAttr(ctx context.Context, req *fuse.SetXAttrRequest) error
}

type NodeRemoveXAttrer interface {
	// RemoveXAttr removes an extended attribute for the name.
	//
	// If there is no XAttr by that name, returns fuse.ErrNoXAttr.
	RemoveXAttr(ctx context.Context, req *fuse.RemoveXAttrRequest) error
}

var startTime = time.Now()

func nodeAttr(ctx context.Context, n Node, attr *fuse.Attr) error {
	attr.Valid = attrValidTime
	attr.Nlink = 1
	attr.Atime = startTime
	attr.Mtime = startTime
	attr.Ctime = startTime
	if err := n.Attr(ctx, attr); err != nil {
		return err
	}
	return nil
}
