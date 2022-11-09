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
	DefaultWriter io.Writer = os.Stdout
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
func New(server *grpc.Server, network string, port int, l logger.Interface, opts ...Option) (*Server, error) {
	lis, err := net.Listen(network, fmt.Sprintf("localhost:%d", port))
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

	s.debugMessage()

	go s.start()

	return s, nil
}

func (s *Server) debugMessage() {
	fmt.Fprintln(DefaultWriter, "[GRPC]")
	for k, v := range s.server.GetServiceInfo() {
		for _, method := range v.Methods {
			fmt.Fprintf(
				DefaultWriter,
				" - %s:  -->  %s stream: server:%v/client:%v\n",
				k,
				method.Name,
				method.IsServerStream,
				method.IsClientStream,
			)
		}
		fmt.Fprintln(DefaultWriter)
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
	select {
	case <-s.notify:
		return nil
	default:
	}

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("grpc server - Server - Shutdown - s.listener.Close: %w", err)
	}

	s.server.GracefulStop()
	return nil
}
