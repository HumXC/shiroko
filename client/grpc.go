package client

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrGrpcCallFailed = errors.New("grpc call falied")

// 解析错误，返回一个合适的错误。如果是 grpc 造成的错误则会 Wrap 一个 ErrGrpcCallFailed
func ParseError(err error) error {
	if s, ok := status.FromError(err); ok {
		code := s.Code()
		if code != codes.OK && code != codes.Aborted {
			return fmt.Errorf("%w: %w", ErrGrpcCallFailed, err)
		}
	}
	return err
}
