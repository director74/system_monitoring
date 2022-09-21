package internalgrpc

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/director74/system_monitoring/internal/app"
	"github.com/director74/system_monitoring/internal/metrics"
	"github.com/director74/system_monitoring/pkg/grpc/protostat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	protostat.UnimplementedAgentServer
	agent         app.Application
	activeClients int32
}

func NewService(agent app.Application) *Service {
	return &Service{
		activeClients: 0,
		agent:         agent,
	}
}

func (s *Service) GetStats(timings *protostat.Timings, statStream protostat.Agent_GetStatsServer) error {
	atomic.AddInt32(&s.activeClients, 1)
	log.Printf("active clients: %d", s.activeClients)
	beginTime := time.Now().Unix()
	for {
		select {
		case <-statStream.Context().Done():
			atomic.AddInt32(&s.activeClients, -1)
			log.Printf("active clients: %d", s.activeClients)
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(time.Duration(timings.GetN()) * time.Second)

			currentTime := time.Now().Unix()
			if (beginTime + timings.GetM()) < currentTime {
				break
			}

			response, err := s.buildResponse(currentTime, timings.GetM())
			if err != nil {
				log.Println(err)
				break
			}
			err = statStream.SendMsg(response)
			if err != nil {
				atomic.AddInt32(&s.activeClients, -1)
				log.Printf("active clients: %d", s.activeClients)
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}

func (s *Service) buildResponse(currentTime int64, periodSeconds int64) (*protostat.SystemStats, error) {
	from := currentTime - periodSeconds
	resultStat := &protostat.SystemStats{}

	names := s.agent.GetAllMetricNames()
	for _, name := range names {
		switch name {
		case "LoadAverage":
			la, err := s.agent.GetMetricStat(name)
			if err != nil {
				return nil, err
			}
			laValues, err := la.GetAverageByPeriod(from, currentTime)
			if err != nil {
				return nil, err
			}
			typedValues, success := laValues.(metrics.LoadAverageResult)
			if !success {
				return nil, fmt.Errorf("cast LoadAverageResult problem")
			}
			resultStat.LoadAverage = &protostat.LoadAverage{Minute1: typedValues.Minute1, Minute5: typedValues.Minute5, Minute15: typedValues.Minute15}
		}
	}

	return resultStat, nil
}
