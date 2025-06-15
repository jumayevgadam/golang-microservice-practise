package server

// Server represent server configurations for this stocks service.
type Server struct {
}

// NewServer creates and returns a new instance of Server.
func NewServer() *Server {
	return &Server{}
}

// RunHTTPServer starts http server in goroutines and gracefully shutdown if signal catches.
func (s *Server) RunHTTPServer() error {
	return nil
}
