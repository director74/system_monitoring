package internalgrpc

import (
	"fmt"
	"log"
	"net"

	"github.com/director74/system_monitoring/internal/app"
	"github.com/director74/system_monitoring/pkg/grpc/protostat"
	"google.golang.org/grpc"
)

type Server struct {
	srv   *grpc.Server
	agent app.Application
	host  string
	port  string
}

func NewServer(port string, agent app.Application) *Server {
	if port == "" {
		port = agent.GetConfig().GetGRPCServerConf().Port
	}
	return &Server{
		host:  agent.GetConfig().GetGRPCServerConf().Host,
		port:  port,
		agent: agent,
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer()
	protostat.RegisterAgentServer(s.srv, NewService(s.agent))
	log.Printf("starting grpc server on %s", lsn.Addr().String())
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
