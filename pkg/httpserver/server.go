// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"finstar-test-task/internal/usecase"
	"finstar-test-task/proto"
	"log"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	proto.UnimplementedUserServiceServer
	mux *runtime.ServeMux
	// could be set trough options to use ServerValidationUnaryInterceptor
	// grpcServer *grpc.Server
	rl usecase.Controller
}

func (s *Server) IncreaseBalance(_ context.Context, request *proto.IncreaseBalanceRequest) (*proto.IncreaseBalanceResponse, error) {
	return s.rl.IncreaseBalance(*request)
}

func (s *Server) TransferBalance(_ context.Context, request *proto.TransferBalanceRequest) (*proto.TransferBalanceResponse, error) {
	return s.rl.TransferBalance(*request)
}

// New -.
func New(rl usecase.Controller, opts ...Option) *Server {
	httpServer := &http.Server{
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		rl:              rl,
		mux: runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{EmitUnpopulated: true},
			}),
		),
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := proto.RegisterUserServiceHandlerServer(ctx, s.mux, s); err != nil {
		log.Panic(err)
	}

	s.start()

	log.Printf("Listening gateway %s", s.server.Addr)

	return s
}

func startgRpcServer() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			ServerValidationUnaryInterceptor,
		)),
	)
	proto.RegisterUserServiceServer(grpcServer, &Server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Validator interface {
	Validate() error
}

// ServerValidationUnaryInterceptor to validate all requests with validate options
func ServerValidationUnaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if r, ok := req.(Validator); ok {
		if err := r.Validate(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}

func (s *Server) start() {
	go func() {
		s.notify <- http.ListenAndServe(s.server.Addr, s.mux)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
