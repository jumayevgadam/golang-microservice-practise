package server

import (
	"bytes"
	grpcV1 "cart/internal/controller/grpc/v1"
	"cart/internal/metrics"
	"cart/internal/repository/postgres"
	"cart/internal/service/stockms"
	"cart/internal/usecase/carts"
	pb "cart/pkg/api/cart"
	"cart/pkg/log"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

func (s *Server) registerGRPCServices() error {
	// repos.
	cartRepo := postgres.NewCartItemRepository(s.psqlDB)

	// services.
	stockService, err := stockms.NewGRPCStockService(s.cfg.StockServiceGRPCAddress())
	if err != nil {
		return fmt.Errorf("failed to create new gRPC stock service: %w", err)
	}

	// usecases.
	cartUseCase := carts.NewCartServiceUseCase(stockService, cartRepo, s.kafkaProducer)

	cartGRPCHandler := grpcV1.NewCartGRPCHandler(cartUseCase, s.logger)

	pb.RegisterCartServiceServer(s.grpcServer, cartGRPCHandler)

	return nil
}

func grpcMiddleware(logger log.Logger, metrics metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		requestID := uuid.New().String()

		ctx, span := otel.Tracer(os.Getenv("SERVICE_NAME")).Start(ctx, info.FullMethod)
		defer span.End()

		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		traceID := span.SpanContext().TraceID().String()
		logFields := map[string]interface{}{
			"method":     info.FullMethod,
			"trace_id":   traceID,
			"request_id": requestID,
			"duration":   duration,
		}

		var msg string
		if err != nil {
			logFields["level"] = "error"
			logFields["msg"] = "gRPC request failed"
			logFields["error"] = err.Error()
			metrics.IncError(info.FullMethod)
			span.SetAttributes(attribute.String("error.message", err.Error()))
		} else {
			logFields["level"] = "info"
			logFields["msg"] = "gRPC request processed"
			msg = "gRPC request processed"
		}

		fields := make([]log.Field, 0, len(logFields))
		for k, v := range logFields {
			fields = append(fields, log.Any(k, v))
		}

		logger.Info(msg, fields...)
		metrics.ObserveLatency(info.FullMethod, duration)

		return resp, err
	}
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
