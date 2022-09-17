package internalgrpc

import (
	"github.com/director74/system_monitoring/pkg/grpc/protostat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	protostat.UnimplementedAgentServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetStats(timings *protostat.Timings, statStream protostat.Agent_GetStatsServer) error {
	for {
		select {
		case <-statStream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(15 * time.Second)

			err := statStream.SendMsg(&protostat.SystemStats{NetStats: nil})
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}
