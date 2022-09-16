package internalgrpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	srv  *grpc.Server
	host string
	port string
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer()
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
