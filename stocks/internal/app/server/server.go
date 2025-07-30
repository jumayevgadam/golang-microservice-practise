package server

import (
	"bytes"
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

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Server represent server configurations for this stocks service.
type Server struct {
	server        *http.Server
	grpcServer    *grpc.Server
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

	errChan := make(chan error, 2)

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
	s.grpcServer = grpc.NewServer()
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
	mux.Handle("/metrics", promhttp.Handler())

	grpcEndpoint := s.cfg.GRPCAddress()

	// register gRPC-gateway handlers.
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterStocksServiceHandlerFromEndpoint(ctx, gatewayMux, grpcEndpoint, opts)
	if err != nil {
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

func observalityMiddleware(logger log.Logger, metrics metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// extract or generate request ID.
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// extract trace context and start a new span.
			ctx, span := otel.Tracer(os.Getenv("SERVICE_NAME")).Start(r.Context(), r.URL.Path)
			defer span.End()
			// we need to update request context with trace.
			r = r.WithContext(ctx)
			// wrap response writer to capture status code and body.
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           new(bytes.Buffer),
			}

			// call the next handler.
			next.ServeHTTP(rw, r)

			// calculate duration.
			duration := time.Since(start).Seconds()

			traceID := span.SpanContext().TraceID().String()
			logFields := map[string]interface{}{
				"level":      "info",
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     rw.statusCode,
				"trace_id":   traceID,
				"request_id": requestID,
				"duration":   duration,
			}

			// we need to handle errors.
			var msg string
			if rw.statusCode >= 400 {
				logFields["level"] = "error"
				logFields["msg"] = "Request failed"

				if rw.body.Len() > 0 {
					logFields["error"] = rw.body.String()
				} else {
					logFields["error"] = http.StatusText(rw.statusCode)
				}

				metrics.IncError(r.URL.Path)

				// record error in span.
				span.SetAttributes(attribute.Int("http.status_code", rw.statusCode))
				span.SetAttributes(attribute.String("error.message", logFields["error"].(string)))
			} else {
				logFields["level"] = "info"
				logFields["msg"] = "HTTP request processed"
				msg = "HTTP request processed"
			}

			fields := make([]log.Field, 0, len(logFields))
			for k, v := range logFields {
				fields = append(fields, log.Any(k, v))
			}

			logger.Info(msg, fields...)
			// at last record latency.
			metrics.ObserveLatency(r.URL.Path, duration)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code and response body
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
