package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/go-chi/chi/v5"

	server "github.com/painb0w/data_processor/internal/api/player_service_api"
	playersubscriptionconsumer "github.com/painb0w/data_processor/internal/consumer/player_subscription_consumer"
	pb "github.com/painb0w/data_processor/internal/pb/players_api"
	"github.com/painb0w/data_processor/internal/services/playersservice"
)

func AppRun(ctx context.Context, cancel context.CancelFunc, api *server.PlayersServiceAPI, consumer *playersubscriptionconsumer.PlayerSubscriptionConsumer, playerService *playersservice.PlayerService) {
	go consumer.Consume(context.Background())

	defer cancel()

	go func() {
		if err := runGRPCServer(api); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %w", err))
		}
	}()

	go func() {
		playerService.StartBackgroundUpdates(ctx, 6*time.Hour)
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %w", err))
	}
}

func runGRPCServer(api *server.PlayersServiceAPI) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlayersServiceServer(s, api)

	slog.Info("gRPC server listening on :50051")
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterPlayersServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	r := chi.NewRouter()
	r.Mount("/", mux)

	slog.Info("gRPC-Gateway (HTTP/JSON) listening on :8080")
	return http.ListenAndServe(":8080", r)
}
