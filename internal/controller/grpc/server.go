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
	"main/internal/state/ctx"
)

type Server struct {
	pb.UnimplementedLeaderboardServiceServer
	AppCtx *ctx.AppContext
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
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(ensureValidToken),
		// Enable TLS for all incoming connections.
		//grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	server := grpc.NewServer(opts...)

	pb.RegisterLeaderboardServiceServer(server, &Server{
		AppCtx: appCtx,
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

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
