package fs

import (
	"context"
	"syscall"

	"bazil.org/fuse"
)

type FileSystemImpl struct{}

func (f FileSystemImpl) String() string {
	return "not implemented"
}

func (f FileSystemImpl) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (Node, error) {
	return nil, syscall.ENOSYS
}

func (f FileSystemImpl) Forget(nodeid, nlookup uint64) {

}

func (f FileSystemImpl) GetAttr(ctx context.Context, req *fuse.GetAttrRequest, resp *fuse.GetAttrResponse) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) SetAttr(ctx context.Context, req *fuse.SetAttrRequest, resp *fuse.SetAttrResponse) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Mknod(ctx context.Context, req *fuse.MknodRequest) (Node, error) {
	return nil, syscall.ENOSYS
}

func (f FileSystemImpl) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (Node, error) {
	return nil, syscall.ENOSYS
}

func (f FileSystemImpl) Unlink(ctx context.Context, req *fuse.UnlinkRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Rmdir(ctx context.Context, req *fuse.RmdirRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Rename(ctx context.Context, req *fuse.RenameRequest, new Node) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Link(ctx context.Context, req *fuse.LinkRequest, old Node) (Node, error) {
	return nil, syscall.ENOSYS
}

func (f FileSystemImpl) Symlink(ctx context.Context, req *fuse.SymlinkRequest) (Node, error) {
	return nil, syscall.ENOSYS
}

func (f FileSystemImpl) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
	return "", syscall.ENOSYS
}

func (f FileSystemImpl) Access(ctx context.Context, req *fuse.AccessRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) GetXAttr(ctx context.Context, req *fuse.GetXAttrRequest, resp *fuse.GetXAttrResponse) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) ListXAttr(ctx context.Context, req *fuse.ListXAttrRequest, resp *fuse.ListXAttrResponse) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) SetXAttr(ctx context.Context, req *fuse.SetXAttrRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) RemoveXAttr(ctx context.Context, req *fuse.RemoveXAttrRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (Node, Handle, error) {
	return nil, nil, syscall.ENOSYS
}

func (f FileSystemImpl) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (Handle, error) {
	return nil, syscall.ENOSYS
}

func (f FileSystemImpl) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	return syscall.ENOSYS
}

func (f FileSystemImpl) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	return syscall.ENOSYS
}

var _ FileSystem = (*FileSystemImpl)(nil)
