package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"service-url-shortener/pkg/logger"
)

const (
	_defaultTimeout = 2 * time.Second
)

var (
	defaultWriter io.Writer = os.Stdout
	debugEnv                = "GRPC_DEBUG"
)

// Server -.
type Server struct {
	listener net.Listener
	server   *grpc.Server
	notify   chan error

	timeout time.Duration

	logger logger.Interface
}

// New -.
func New(server *grpc.Server, network, port string, l logger.Interface, opts ...Option) (*Server, error) {
	lis, err := net.Listen(network, net.JoinHostPort("", port))
	if err != nil {
		return nil, fmt.Errorf("grpc server - NewServer - net.Listen: %w", err)
	}

	s := &Server{
		listener: lis,
		server:   server,
		notify:   make(chan error),
		timeout:  _defaultTimeout,
		logger:   l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	go s.start()

	if os.Getenv(debugEnv) == "true" {
		s.debugMessage()
	}

	return s, nil
}

func (s *Server) debugMessage() {
	fmt.Fprintln(defaultWriter, "[GRPC]")

	for k, v := range s.server.GetServiceInfo() {
		for _, method := range v.Methods {
			fmt.Fprintf(
				defaultWriter,
				" - %s:  -->  %s stream(server:%v client:%v)\n",
				k,
				method.Name,
				method.IsServerStream,
				method.IsClientStream,
			)
		}

		fmt.Fprintln(defaultWriter)
	}
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.Serve(s.listener)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	if s.notify == nil {
		return nil
	}

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("grpc server - Server - Shutdown - s.listener.Close: %w", err)
	}

	s.server.GracefulStop()
	s.notify = nil

	return nil
}
