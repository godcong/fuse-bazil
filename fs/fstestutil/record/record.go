package record // import "bazil.org/fuse/fs/fstestutil/record"

import (
	"context"
	"sync"
	"sync/atomic"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// Writes gathers data from FUSE Write calls.
type Writes struct {
	buf Buffer
}

var _ fs.HandleWriter = (*Writes)(nil)

func (w *Writes) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	n, err := w.buf.Write(req.Data)
	resp.Size = n
	if err != nil {
		return err
	}
	return nil
}

func (w *Writes) RecordedWriteData() []byte {
	return w.buf.Bytes()
}

// Counter records number of times a thing has occurred.
type Counter struct {
	count uint32
}

func (r *Counter) Inc() {
	atomic.AddUint32(&r.count, 1)
}

func (r *Counter) Count() uint32 {
	return atomic.LoadUint32(&r.count)
}

// MarkRecorder records whether a thing has occurred.
type MarkRecorder struct {
	count Counter
}

func (r *MarkRecorder) Mark() {
	r.count.Inc()
}

func (r *MarkRecorder) Recorded() bool {
	return r.count.Count() > 0
}

// Flushes notes whether a FUSE Flush call has been seen.
type Flushes struct {
	rec MarkRecorder
}

var _ fs.HandleFlusher = (*Flushes)(nil)

func (r *Flushes) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	r.rec.Mark()
	return nil
}

func (r *Flushes) RecordedFlush() bool {
	return r.rec.Recorded()
}

type Recorder struct {
	mu  sync.Mutex
	val interface{}
}

// Record that we've seen value. A nil value is indistinguishable from
// no value recorded.
func (r *Recorder) Record(value interface{}) {
	r.mu.Lock()
	r.val = value
	r.mu.Unlock()
}

func (r *Recorder) Recorded() interface{} {
	r.mu.Lock()
	val := r.val
	r.mu.Unlock()
	return val
}

type RequestRecorder struct {
	rec Recorder
}

// Record a fuse.Request, after zeroing header fields that are hard to
// reproduce.
//
// Make sure to record a copy, not the original request.
func (r *RequestRecorder) RecordRequest(req fuse.Request) {
	hdr := req.Hdr()
	*hdr = fuse.Header{}
	r.rec.Record(req)
}

func (r *RequestRecorder) Recorded() fuse.Request {
	val := r.rec.Recorded()
	if val == nil {
		return nil
	}
	return val.(fuse.Request)
}

// SetAttrs records a SetAttr request and its fields.
type SetAttrs struct {
	rec RequestRecorder
}

var _ fs.NodeSetAttrer = (*SetAttrs)(nil)

func (r *SetAttrs) SetAttr(ctx context.Context, req *fuse.SetAttrRequest, resp *fuse.SetAttrResponse) error {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil
}

func (r *SetAttrs) RecordedSetAttr() fuse.SetAttrRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.SetAttrRequest{}
	}
	return *(val.(*fuse.SetAttrRequest))
}

// Fsyncs records an Fsync request and its fields.
type Fsyncs struct {
	rec RequestRecorder
}

var _ fs.NodeFsyncer = (*Fsyncs)(nil)

func (r *Fsyncs) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil
}

func (r *Fsyncs) RecordedFsync() fuse.FsyncRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.FsyncRequest{}
	}
	return *(val.(*fuse.FsyncRequest))
}

// Mkdirs records a Mkdir request and its fields.
type Mkdirs struct {
	rec RequestRecorder
}

var _ fs.NodeMkdirer = (*Mkdirs)(nil)

// Mkdir records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *Mkdirs) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil, syscall.EIO
}

// RecordedMkdir returns information about the Mkdir request.
// If no request was seen, returns a zero value.
func (r *Mkdirs) RecordedMkdir() fuse.MkdirRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.MkdirRequest{}
	}
	return *(val.(*fuse.MkdirRequest))
}

// Symlinks records a Symlink request and its fields.
type Symlinks struct {
	rec RequestRecorder
}

var _ fs.NodeSymlinker = (*Symlinks)(nil)

// Symlink records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *Symlinks) Symlink(ctx context.Context, req *fuse.SymlinkRequest) (fs.Node, error) {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil, syscall.EIO
}

// RecordedSymlink returns information about the Symlink request.
// If no request was seen, returns a zero value.
func (r *Symlinks) RecordedSymlink() fuse.SymlinkRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.SymlinkRequest{}
	}
	return *(val.(*fuse.SymlinkRequest))
}

// Links records a Link request and its fields.
type Links struct {
	rec RequestRecorder
}

var _ fs.NodeLinker = (*Links)(nil)

// Link records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *Links) Link(ctx context.Context, req *fuse.LinkRequest, old fs.Node) (fs.Node, error) {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil, syscall.EIO
}

// RecordedLink returns information about the Link request.
// If no request was seen, returns a zero value.
func (r *Links) RecordedLink() fuse.LinkRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.LinkRequest{}
	}
	return *(val.(*fuse.LinkRequest))
}

// Mknods records a Mknod request and its fields.
type Mknods struct {
	rec RequestRecorder
}

var _ fs.NodeMknoder = (*Mknods)(nil)

// Mknod records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *Mknods) Mknod(ctx context.Context, req *fuse.MknodRequest) (fs.Node, error) {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil, syscall.EIO
}

// RecordedMknod returns information about the Mknod request.
// If no request was seen, returns a zero value.
func (r *Mknods) RecordedMknod() fuse.MknodRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.MknodRequest{}
	}
	return *(val.(*fuse.MknodRequest))
}

// Opens records a Open request and its fields.
type Opens struct {
	rec RequestRecorder
}

var _ fs.NodeOpener = (*Opens)(nil)

// Open records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *Opens) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil, syscall.EIO
}

// RecordedOpen returns information about the Open request.
// If no request was seen, returns a zero value.
func (r *Opens) RecordedOpen() fuse.OpenRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.OpenRequest{}
	}
	return *(val.(*fuse.OpenRequest))
}

// GetXAttrs records a GetXAttr request and its fields.
type GetXAttrs struct {
	rec RequestRecorder
}

var _ fs.NodeGetXAttrer = (*GetXAttrs)(nil)

// GetXAttr records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *GetXAttrs) GetXAttr(ctx context.Context, req *fuse.GetXAttrRequest, resp *fuse.GetXAttrResponse) error {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return fuse.ErrNoXAttr
}

// RecordedGetXAttr returns information about the GetXAttr request.
// If no request was seen, returns a zero value.
func (r *GetXAttrs) RecordedGetXAttr() fuse.GetXAttrRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.GetXAttrRequest{}
	}
	return *(val.(*fuse.GetXAttrRequest))
}

// ListXAttrs records a ListXAttr request and its fields.
type ListXAttrs struct {
	rec RequestRecorder
}

var _ fs.NodeListXAttrer = (*ListXAttrs)(nil)

// ListXAttr records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *ListXAttrs) ListXAttr(ctx context.Context, req *fuse.ListXAttrRequest, resp *fuse.ListXAttrResponse) error {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return fuse.ErrNoXAttr
}

// RecordedListXAttr returns information about the ListXAttr request.
// If no request was seen, returns a zero value.
func (r *ListXAttrs) RecordedListXAttr() fuse.ListXAttrRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.ListXAttrRequest{}
	}
	return *(val.(*fuse.ListXAttrRequest))
}

// SetXAttrs records a SetXAttr request and its fields.
type SetXAttrs struct {
	rec RequestRecorder
}

var _ fs.NodeSetXAttrer = (*SetXAttrs)(nil)

// SetXAttr records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *SetXAttrs) SetXAttr(ctx context.Context, req *fuse.SetXAttrRequest) error {
	tmp := *req
	// The byte slice points to memory that will be reused, so make a
	// deep copy.
	tmp.Xattr = append([]byte(nil), req.Xattr...)
	r.rec.RecordRequest(&tmp)
	return nil
}

// RecordedSetXAttr returns information about the SetXAttr request.
// If no request was seen, returns a zero value.
func (r *SetXAttrs) RecordedSetXAttr() fuse.SetXAttrRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.SetXAttrRequest{}
	}
	return *(val.(*fuse.SetXAttrRequest))
}

// RemoveXAttrs records a RemoveXAttr request and its fields.
type RemoveXAttrs struct {
	rec RequestRecorder
}

var _ fs.NodeRemoveXAttrer = (*RemoveXAttrs)(nil)

// RemoveXAttr records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *RemoveXAttrs) RemoveXAttr(ctx context.Context, req *fuse.RemoveXAttrRequest) error {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil
}

// RecordedRemoveXAttr returns information about the RemoveXAttr request.
// If no request was seen, returns a zero value.
func (r *RemoveXAttrs) RecordedRemoveXAttr() fuse.RemoveXAttrRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.RemoveXAttrRequest{}
	}
	return *(val.(*fuse.RemoveXAttrRequest))
}

// Creates records a Create request and its fields.
type Creates struct {
	rec RequestRecorder
}

var _ fs.NodeCreater = (*Creates)(nil)

// Create records the request and returns an error. Most callers should
// wrap this call in a function that returns a more useful result.
func (r *Creates) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	tmp := *req
	r.rec.RecordRequest(&tmp)
	return nil, nil, syscall.EIO
}

// RecordedCreate returns information about the Create request.
// If no request was seen, returns a zero value.
func (r *Creates) RecordedCreate() fuse.CreateRequest {
	val := r.rec.Recorded()
	if val == nil {
		return fuse.CreateRequest{}
	}
	return *(val.(*fuse.CreateRequest))
}
