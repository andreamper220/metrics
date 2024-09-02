package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreamper220/metrics.git/internal/agent"
	"github.com/andreamper220/metrics.git/internal/logger"
	pb "github.com/andreamper220/metrics.git/proto"
)

type requestStruct struct {
	url        string
	bodyStruct interface{}
	client     *pb.MetricsClient
}

func main() {
	if err := agent.ParseFlags(); err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(agent.Config.ServerAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewMetricsClient(conn)

	go agent.UpdateMetrics()
	go agent.UpdatePsUtilsMetrics()

	requestCh := make(chan requestStruct)
	errCh := make(chan error)
	for s := 1; s <= agent.Config.RateLimit; s++ {
		go Sender(requestCh, errCh, c)
	}

	ctxSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	stopCh := make(chan struct{})
	go sendMetrics(ctxSignal, requestCh, stopCh)

	for {
		select {
		case err := <-errCh:
			logger.Log.Error(err.Error())
		case <-stopCh:
			return
		}
	}
}

func Sender(requestCh <-chan requestStruct, errCh chan<- error, c pb.MetricsClient) {
	defer close(errCh)

	for request := range requestCh {
		_, err := c.UpdateMetrics(context.Background(), &pb.UpdateMetricsRequest{
			Metrics: request.bodyStruct.([]*pb.Metric),
		})
		if err != nil {
			errCh <- err
			continue
		}
	}
}

func sendMetrics(context context.Context, requestCh chan<- requestStruct, stopCh chan<- struct{}) {
	defer close(requestCh)

	reportTicker := time.NewTicker(time.Duration(agent.Config.ReportInterval) * time.Second)
	for {
		select {
		case <-context.Done():
			reportTicker.Stop()
			stopCh <- struct{}{}
		case <-reportTicker.C:
			go func() {
				requestCh <- requestStruct{
					bodyStruct: buildMetrics(),
				}
			}()
		}
	}
}

func buildMetrics() []*pb.Metric {
	currentMetrics := agent.ReadMetrics()

	metricsSlice := make([]*pb.Metric, len(currentMetrics.Gauges)+len(currentMetrics.Counters))
	metricsIndex := 0
	for name, value := range currentMetrics.Gauges {
		metricsSlice[metricsIndex] = &pb.Metric{
			Id:    string(name),
			Type:  pb.Metric_GAUGE,
			Value: &value,
		}
		metricsIndex++
	}
	for name, delta := range currentMetrics.Counters {
		metricsSlice[metricsIndex] = &pb.Metric{
			Id:    string(name),
			Type:  pb.Metric_COUNTER,
			Delta: &delta,
		}
		metricsIndex++
	}

	return metricsSlice
}
