package server

import (
	"net"
	"strings"

	"github.com/HumXC/shiroko/android"
	"github.com/grandcat/zeroconf"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	zeroConf   *zeroconf.Server
	lis        net.Listener
	name       string
}

func (s *Server) Serve(lis net.Listener) error {
	s.lis = lis
	text := []string{"model=", "port="}
	z, err := zeroconf.Register(s.name, "_shiroko._tcp", "local.", 15600, text, nil)
	if err != nil {
		return err
	}
	s.zeroConf = z
	model := android.Model()
	text[0] += model
	ports := strings.Split(s.lis.Addr().String(), ":")
	text[1] += ports[len(ports)-1]
	return s.grpcServer.Serve(lis)
}
func (s *Server) Stop() {
	if s.zeroConf != nil {
		s.zeroConf.Shutdown()
	}
	s.grpcServer.Stop()
}
func New(name string) *Server {
	s := &Server{
		name:       name,
		grpcServer: newGrpcServer(),
	}
	return s
}
