package common

import (
	"io"

	"google.golang.org/grpc"
)

type ServerStreamer interface {
	Send(*DataChunk) error
	grpc.ClientStream
}
type writer struct {
	stream grpc.ServerStream
}

func (w *writer) Write(p []byte) (n int, err error) {
	chunk := &DataChunk{Data: p}
	if err := w.stream.SendMsg(chunk); err != nil {
		return 0, err
	}
	return len(p), nil
}
func NewWriter(stream grpc.ServerStream) io.Writer {
	return &writer{stream}
}

type ClientStreamer interface {
	Recv() (*DataChunk, error)
	grpc.ClientStream
}
type readCloser struct {
	stream ClientStreamer
	buf    []byte
}

func (r *readCloser) Close() error {
	return r.stream.CloseSend()
}

func (r *readCloser) Read(p []byte) (n int, err error) {
	if len(r.buf) == 0 { // 如果buffer中没有数据，尝试从流中获取
		chunk, err := r.stream.Recv()
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, err
		}
		r.buf = chunk.Data
	}
	// 从buffer中复制数据到p
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

// 创建一个 io.ReadCloser
func NewReadCloser(stream ClientStreamer) io.ReadCloser {
	return &readCloser{
		stream: stream,
		buf:    make([]byte, 0),
	}
}
