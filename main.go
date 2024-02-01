package main

import (
	dao "cloudbees/dao"
	postsGrpc "cloudbees/genproto/posts"
	"cloudbees/services"
	svc "cloudbees/services"
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

var postsService *svc.PostsService
var logger *zap.Logger

type server struct {
	server *grpc.Server
}

func createServer() *server {
	return &server{
		server: grpc.NewServer(
			grpc.UnaryInterceptor(loggingInterceptor),
		),
	}
}

func initPostsService(postsDao *dao.PostDAO) {
	postsService = services.NewPostsService(postsDao)
}

func init() {
	postsDao := dao.NewPostDAO()
	initPostsService(postsDao)

	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func (s *server) registerService(service grpc.ServiceRegistrar) {
	postsGrpc.RegisterBlogServiceServer(s.server, postsService)
}

func (s *server) serve(listener net.Listener) error {
	return s.server.Serve(listener)
}

func (s *server) start(listener net.Listener) {
	err := s.server.Serve(listener)
	if err != nil {
		logger.Sugar().Fatalf("cannot serve grpc server: %s", err)
	}
}

// loggingInterceptor  logs the request and response.
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Info("gRPC method", zap.String("method", info.FullMethod), zap.Any("request", req))
	resp, err := handler(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		logger.Error("gRPC method", zap.String("method", info.FullMethod), zap.Error(err), zap.String("code", st.Code().String()))
	} else {
		logger.Info("gRPC method", zap.String("method", info.FullMethod), zap.Any("response", resp))
	}
	return resp, err
}

func main() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		logger.Sugar().Fatalf("cannot create Listener: %s", err)
	}
	defer listener.Close()

	logger.Info("Server started on port 80")

	s := createServer()
	s.registerService(s.server)
	s.start(listener)
}
