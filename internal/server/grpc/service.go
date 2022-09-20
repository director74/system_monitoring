package internalgrpc

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/director74/system_monitoring/pkg/grpc/protostat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	protostat.UnimplementedAgentServer
	activeClients int32
}

func NewService() *Service {
	return &Service{
		activeClients: 0,
	}
}

func (s *Service) GetStats(timings *protostat.Timings, statStream protostat.Agent_GetStatsServer) error {
	atomic.AddInt32(&s.activeClients, 1)
	log.Printf("active clients: %d", s.activeClients)
	for {
		select {
		case <-statStream.Context().Done():
			atomic.AddInt32(&s.activeClients, -1)
			log.Printf("active clients: %d", s.activeClients)
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(time.Duration(timings.GetM()) * time.Second)

			err := statStream.SendMsg(&protostat.SystemStats{CpuLoad: &protostat.CpuLoad{UserMode: 1, SystemMode: 2, Idle: 3}})
			if err != nil {
				atomic.AddInt32(&s.activeClients, -1)
				log.Printf("active clients: %d", s.activeClients)
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}
