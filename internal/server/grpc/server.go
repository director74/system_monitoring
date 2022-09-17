package internalgrpc

import (
	"fmt"
	"github.com/director74/system_monitoring/internal/cfg"
	"log"
	"net"

	"github.com/director74/system_monitoring/pkg/grpc/protostat"
	"google.golang.org/grpc"
)

type Server struct {
	srv  *grpc.Server
	host string
	port string
}

func NewServer(conf cfg.Configurable) *Server {
	return &Server{
		host: conf.GetGRPCServerConf().Host,
		port: conf.GetGRPCServerConf().Port,
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer()
	protostat.RegisterAgentServer(s.srv, NewService())
	log.Printf("starting grpc server on %s", lsn.Addr().String())
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
