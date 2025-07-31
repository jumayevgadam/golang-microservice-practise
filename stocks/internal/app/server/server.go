package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"stocks/internal/config"
	"stocks/internal/kafka"
	"stocks/internal/metrics"
	pb "stocks/pkg/api/stocks"
	"stocks/pkg/connection"
	"stocks/pkg/constants"
	"stocks/pkg/log"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Server represent server configurations for this stocks service.
type Server struct {
	server        *http.Server
	grpcServer    *grpc.Server
	metricsServer *http.Server
	cfg           config.Config
	psqlDB        connection.DB
	kafkaProducer kafka.StocksEventProducer
	logger        log.Logger
	metrics       metrics.Metrics
}

// NewServer creates and returns a new instance of Server.
func NewServer(
	cfg config.Config,
	psqlDB connection.DB,
	kafkaProducer kafka.StocksEventProducer,
	logger log.Logger,
) *Server {
	return &Server{
		server:        nil,
		grpcServer:    nil,
		metricsServer: nil,
		cfg:           cfg,
		psqlDB:        psqlDB,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		metrics:       metrics.RegisterMetrics(),
	}
}

// RunServer starts http and grpc server in goroutines and gracefully shutdown if signal catches.
func (s *Server) RunServer() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// start grpc server.
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.runGRPCServer(); err != nil {
			errChan <- fmt.Errorf("runGRPCServer: %w", err)
		}
	}()

	// start grpc-gateway server.
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.runGatewayServer(); err != nil {
			errChan <- fmt.Errorf("runGatewayServer: %w", err)
		}
	}()

	// start metrics server.
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.runMetricsServer(); err != nil {
			errChan <- fmt.Errorf("runMetricsServer: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// wait for a signal or an error from the servers.
	select {
	case <-quit:
		s.logger.Info("shutting down server...")
	case err := <-errChan:
		s.logger.Errorf("Server error: %v", err.Error())
	}

	// Create context for shutdown
	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.SrvTimeOut*time.Second)
	defer cancel()

	// shutdown http server.
	if s.server != nil {
		if err := s.server.Shutdown(ctxTimeOut); err != nil {
			s.logger.Errorf("s.server.Shutdown: %v", err.Error())
		}
	}

	// Shutdown gRPC server.
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	// Shutdown metrics server.
	if s.metricsServer != nil {
		if err := s.metricsServer.Shutdown(ctxTimeOut); err != nil {
			s.logger.Errorf("s.metricsServer.Shutdown: %v", err.Error())
		}
	}

	wg.Wait()

	s.logger.Info("stock service successfully shut down...")

	return nil
}

func (s *Server) runGRPCServer() error {
	lis, err := net.Listen("tcp", s.cfg.GRPCAddress())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.GRPCAddress(), err)
	}
	defer lis.Close()
	// create a grpc server.
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware(s.logger, s.metrics)),
	)
	// enable reflection for grpcui.
	s.registerGRPCServices()
	reflection.Register(s.grpcServer)

	s.logger.Infof("grpc Server starting on %s", s.cfg.GRPCAddress())

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve stock service gRPC: %w", err)
	}

	return nil
}

func (s *Server) runGatewayServer() error {
	// create a context for a gateway.
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create grpc-gateway mux.
	gatewayMux := runtime.NewServeMux()

	handler := observalityMiddleware(s.logger, s.metrics)(gatewayMux)

	// create a new serve mux for api and metrics endpoint.
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	grpcEndpoint := s.cfg.GRPCAddress()

	// register gRPC-gateway handlers.
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterStocksServiceHandlerFromEndpoint(ctx, gatewayMux, grpcEndpoint, opts)
	if err != nil {
		s.logger.Errorf("error register stocks service handler: %v", err.Error())
		return fmt.Errorf("failed to register gateway handler: %w", err)
	}

	s.server = &http.Server{
		Addr:         s.cfg.Address(),
		Handler:      mux,
		ReadTimeout:  s.cfg.SrvConfig().ReadTimeOut,
		WriteTimeout: s.cfg.SrvConfig().WriteTimeOut,
	}

	s.logger.Infof("Gateway server starting on %s", s.cfg.Address())

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve gateway: %w", err)
	}

	return nil
}

func (s *Server) runMetricsServer() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	s.metricsServer = &http.Server{
		Addr:         s.cfg.MetricsAddress(),
		Handler:      mux,
		ReadTimeout:  s.cfg.SrvConfig().ReadTimeOut,
		WriteTimeout: s.cfg.SrvConfig().WriteTimeOut,
	}

	s.logger.Infof("Metrics server starting on: %s", s.cfg.MetricsAddress())

	if err := s.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve metrics server: %w", err)
	}

	return nil
}
