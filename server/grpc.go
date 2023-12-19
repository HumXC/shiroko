package server

import (
	pInput "github.com/HumXC/shiroko/protos/input"
	pManager "github.com/HumXC/shiroko/protos/manager"
	pMinicap "github.com/HumXC/shiroko/protos/minicap"
	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"google.golang.org/grpc"
)

func newGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pManager.RegisterManagerServer(grpcServer, NewManagerServer())
	pScreencap.RegisterScreencapServer(grpcServer, NewScreencapServer())
	pMinicap.RegisterMinicapServer(grpcServer, NewMinicapServer())
	pInput.RegisterInputServer(grpcServer, NewInputServer())
	return grpcServer
}
