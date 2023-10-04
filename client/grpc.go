package client

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 解析错误，返回一个合适的错误。
// 如果错误并不是由 grpc 等外部原因造成的，例如网络断开等错误
// 便会返回一个由服务端函数返回的原始错误，否则返回原错误
func ParseError(err error) error {
	s, ok := status.FromError(err)
	if ok && s.Code() == codes.Unknown {
		return errors.New(s.Message())
	}
	return err
}
