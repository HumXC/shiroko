package client

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrGrpcCallFailed = errors.New("grpc call falied")

// 解析错误，返回一个合适的错误。如果是 grpc 造成的错误则会 Wrap 一个 ErrGrpcCallFailed
func ParseError(err error) error {
	if s, ok := status.FromError(err); ok {
		if s.Code() != codes.Aborted {
			return fmt.Errorf("%w: %w", ErrGrpcCallFailed, err)
		}
	}
	return err
}

type readCloser struct {
	stream        grpc.ClientStream
	buf           []byte
	dataChunkType reflect.Type
}

func (r *readCloser) Close() error {
	return r.stream.CloseSend()
}

func (r *readCloser) Read(p []byte) (n int, err error) {
	if len(r.buf) == 0 { // 如果buffer中没有数据，尝试从流中获取
		dataChunk := reflect.New(r.dataChunkType).Interface()
		err := r.stream.RecvMsg(dataChunk)
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, err
		}
		dataChunkVal := reflect.ValueOf(dataChunk).Elem()
		dataField := dataChunkVal.FieldByName("Data")
		r.buf = dataField.Bytes()
	}
	// 从buffer中复制数据到p
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

// 创建一个 io.ReadCloser, dataChunkStruct 是 grpc 生成的 DataChunk 空结构体
func NewReadCloser(stream grpc.ClientStream, dataChunkStruct any) io.ReadCloser {
	dataChunkType := reflect.TypeOf(dataChunkStruct)
	if dataChunkType.Kind() == reflect.Ptr {
		dataChunkType = dataChunkType.Elem()
	}
	return &readCloser{
		stream:        stream,
		buf:           make([]byte, 0, 1024*1024),
		dataChunkType: dataChunkType,
	}
}
