package internalgrpc

import (
	"log"
	"time"

	"github.com/director74/system_monitoring/pkg/grpc/protostat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	protostat.UnimplementedAgentServer
	activeClients int
}

func NewService() *Service {
	return &Service{
		activeClients: 0,
	}
}

func (s *Service) GetStats(timings *protostat.Timings, statStream protostat.Agent_GetStatsServer) error {
	s.activeClients++
	log.Printf("active clients: %d", s.activeClients)
	for {
		select {
		case <-statStream.Context().Done():
			s.activeClients--
			log.Printf("active clients: %d", s.activeClients)
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(time.Duration(timings.GetM()) * time.Second)

			err := statStream.SendMsg(&protostat.SystemStats{CpuLoad: &protostat.CpuLoad{UserMode: 1, SystemMode: 2, Idle: 3}})
			if err != nil {
				s.activeClients--
				log.Printf("active clients: %d", s.activeClients)
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}
