package server

import (
	"cart/internal/config"
	"cart/internal/kafka"
	pb "cart/pkg/api/cart"
	"cart/pkg/connection"
	"cart/pkg/constants"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server        *http.Server
	grpcServer    *grpc.Server
	cfg           config.Config
	psqlDB        connection.DB
	kafkaProducer kafka.CartEventProducer
}

func NewServer(
	cfg config.Config,
	psqlDB connection.DB,
	kafkaProducer kafka.CartEventProducer,
) *Server {
	return &Server{
		server:        nil,
		cfg:           cfg,
		psqlDB:        psqlDB,
		kafkaProducer: kafkaProducer,
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
		log.Println("shutting down server...")
	case err := <-errChan:
		log.Printf("Server error: %v", err.Error())
	}

	// Create context for shutdown
	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.SrvTimeOut*time.Second)
	defer cancel()

	// shutdown http server.
	if s.server != nil {
		if err := s.server.Shutdown(ctxTimeOut); err != nil {
			log.Printf("s.server.Shutdown: %v", err.Error())
		}
	}

	// Shutdown gRPC server.
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	wg.Wait()

	log.Println("cart service successfully shut down...")

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

	log.Printf("grpc Server starting on %s", s.cfg.GRPCAddress())
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve cart service gRPC: %w", err)
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

	grpcEndpoint := s.cfg.GRPCAddress()

	// register gRPC-gateway handlers.
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterCartServiceHandlerFromEndpoint(ctx, gatewayMux, grpcEndpoint, opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway handler: %w", err)
	}

	s.server = &http.Server{
		Addr:         s.cfg.Address(),
		Handler:      gatewayMux,
		ReadTimeout:  s.cfg.SrvConfig().ReadTimeOut,
		WriteTimeout: s.cfg.SrvConfig().WriteTimeOut,
	}

	log.Printf("Gateway server starting on %s", s.cfg.Address())
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve gateway: %w", err)
	}

	return nil
}
