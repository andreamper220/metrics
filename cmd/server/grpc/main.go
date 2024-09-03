package grpc

import (
	"context"
	"google.golang.org/grpc"
	"net"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/application"
	"github.com/andreamper220/metrics.git/internal/server/domain/metrics"
	"github.com/andreamper220/metrics.git/internal/server/infrastructure/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
	pb "github.com/andreamper220/metrics.git/proto"
)

var (
	ErrIncorrectMetricID   = "Incorrect metric ID"
	ErrIncorrectMetricType = "Incorrect metric TYPE"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func (m MetricsServer) GetMetric(ctx context.Context, request *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse

	reqMetric := request.GetMetric()
	metric := shared.Metric{
		ID:    reqMetric.GetId(),
		MType: reqMetric.GetType().String(),
	}

	switch reqMetric.GetType() {
	case pb.Metric_COUNTER:
		isExisted := false
		counters, err := storages.Storage.GetCounters()
		if err != nil {
			response.Error = err.Error()
			return &response, nil
		}
		for _, counter := range counters {
			if counter.Name == shared.CounterMetricName(metric.ID) {
				reqMetric.Delta = &counter.Value
				isExisted = true
				break
			}
		}
		if !isExisted {
			response.Error = ErrIncorrectMetricID
			return &response, nil
		}
	case pb.Metric_GAUGE:
		isExisted := false
		gauges, err := storages.Storage.GetGauges()
		if err != nil {
			response.Error = err.Error()
			return &response, nil
		}
		for _, gauge := range gauges {
			if gauge.Name == shared.GaugeMetricName(metric.ID) {
				reqMetric.Value = &gauge.Value
				isExisted = true
				break
			}
		}
		if !isExisted {
			response.Error = ErrIncorrectMetricID
			return &response, nil
		}
	default:
		response.Error = ErrIncorrectMetricType
		return &response, nil
	}
	response.Metric = reqMetric
	return &response, nil
}

func (m MetricsServer) UpdateCounter(ctx context.Context, request *pb.UpdateCounterRequest) (*pb.UpdateCounterResponse, error) {
	var response pb.UpdateCounterResponse

	reqMetric := request.GetMetric()
	reqMetricInner := reqMetric.GetMetric()
	reqMetricDelta := reqMetric.GetDelta()
	metric := shared.Metric{
		ID:    reqMetricInner.GetId(),
		MType: reqMetricInner.GetType().String(),
		Delta: &reqMetricDelta,
	}

	if err := metrics.ProcessMetric(&metric); err != nil {
		response.Error = err.Error()
	} else {
		response.Metric = reqMetric
	}
	return &response, nil
}

func (m MetricsServer) UpdateGauge(ctx context.Context, request *pb.UpdateGaugeRequest) (*pb.UpdateGaugeResponse, error) {
	var response pb.UpdateGaugeResponse

	reqMetric := request.GetMetric()
	reqMetricInner := reqMetric.GetMetric()
	reqMetricValue := reqMetric.GetValue()
	metric := shared.Metric{
		ID:    reqMetricInner.GetId(),
		MType: reqMetricInner.GetType().String(),
		Value: &reqMetricValue,
	}

	if err := metrics.ProcessMetric(&metric); err != nil {
		response.Error = err.Error()
	} else {
		response.Metric = reqMetric
	}
	return &response, nil
}

func (m MetricsServer) UpdateMetrics(ctx context.Context, request *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse

	for _, reqMetric := range request.GetMetrics() {
		metric := shared.Metric{
			ID:    reqMetric.GetId(),
			MType: reqMetric.GetType().String(),
		}

		if reqMetricDelta := reqMetric.GetDelta(); reqMetricDelta != 0 {
			metric.Delta = &reqMetricDelta
		}
		if reqMetricValue := reqMetric.GetValue(); reqMetricValue != 0 {
			metric.Value = &reqMetricValue
		}

		if err := metrics.ProcessMetric(&metric); err != nil {
			response.Error = err.Error()
			return &response, nil
		}
	}

	response.Metrics = request.GetMetrics()
	return &response, nil
}

func main() {
	if err := application.ParseFlags(); err != nil {
		logger.Log.Fatal(err)
	}

	listen, err := net.Listen("tcp", application.Config.ServerAddress.String())
	if err != nil {
		logger.Log.Fatal(err)
	}

	if err := initGRPCServer().Serve(listen); err != nil {
		logger.Log.Fatal("gRPC server Serve: %v", err)
	}
}

func initGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterMetricsServer(s, &MetricsServer{})

	logger.Log.Info("gRPC server started")
	return s
}
