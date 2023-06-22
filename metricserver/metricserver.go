package metricserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

const (
	DEFAULT_PORT = 2112
)

type MetricServerConfiguration func(*MetricServer) error

type MetricServer struct {
	s http.Server

	port int

	meterProvider *sdk.MeterProvider
}

func New(cfgs ...MetricServerConfiguration) (*MetricServer, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	provider := sdk.NewMeterProvider(sdk.WithReader(exporter))

	ms := &MetricServer{
		port:          DEFAULT_PORT,
		meterProvider: provider,
	}

	for _, cfg := range cfgs {
		err := cfg(ms)

		if err != nil {
			return nil, err
		}
	}

	return ms, nil
}

func WithPort(Port int) MetricServerConfiguration {
	return func(ms *MetricServer) error {
		ms.port = Port

		return nil
	}
}

func (m *MetricServer) Start() {
	router := http.NewServeMux()

	router.Handle("/metrics", promhttp.Handler())

	m.s.Handler = router
	m.s.Addr = fmt.Sprintf(":%v", m.port)

	logrus.Infof("Started metric server on port %v", m.port)

	err := m.s.ListenAndServe()

	// expected when Shutdown is invoked
	if errors.Is(err, http.ErrServerClosed) {
		logrus.Info(err)
		return
	}

	// Any other error would fatal
	logrus.Fatal(err)
}

func (m *MetricServer) Shutdown() error {
	logrus.Info("shutting down metrics server")
	if err := m.s.Shutdown(context.Background()); err != nil {
		logrus.Error("shutdown metrics server failed:", err)
		return err
	}
	logrus.Info("shutdown metrics server complete")
	return nil
}

func (m *MetricServer) GetMeter(scope string) metric.Meter {
	meter := m.meterProvider.Meter(scope)
	return meter
}
