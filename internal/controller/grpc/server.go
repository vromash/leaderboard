package grpc_controller

import (
	"context"
	"fmt"
	"net"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	config "github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "main/internal/controller/protos"
	"main/internal/score"
	"main/internal/state/ctx"
)

type Server struct {
	pb.UnimplementedLeaderboardServiceServer
	appCtx   *ctx.AppContext
	scoreSvc *score.ScoreService
}

func RunServer(appCtx *ctx.AppContext) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GetString("GRPC_PORT")))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen to GRPC port")
	}

	srv := setup(appCtx)

	log.Info().Msgf("GRPC server listening at: %s", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve GRPC")
	}
}

func setup(appCtx *ctx.AppContext) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureUnaryValidToken),
		grpc.StreamInterceptor(ensureStreamValidToken),
	}

	server := grpc.NewServer(opts...)
	scoreSvc := score.NewScoreService(appCtx)

	pb.RegisterLeaderboardServiceServer(server, &Server{
		appCtx:   appCtx,
		scoreSvc: scoreSvc,
	})
	reflection.Register(server)

	return server
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "secret-token"
}

// ensureUnaryValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureUnaryValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return handler(ctx, req)
}

func ensureStreamValidToken(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return handler(srv, stream)
}
