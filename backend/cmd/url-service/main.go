package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	urlgrpc "zipit/internal/url/grpc"
	"zipit/internal/url/repository"
	"zipit/internal/url/service"
	"zipit/pkg/config"
	"zipit/pkg/database"
	"zipit/pkg/logger"
	"zipit/pkg/shortener"

	"github.com/joho/godotenv"
	grpc "google.golang.org/grpc"

	pb "zipit/gen/url"
)

func main() {
	// 1. Initialize Infrastructure: Logger, Env vars
	logger.SetLogger()
	_ = godotenv.Load() // Optional: only used in local dev, Railway injects env vars directly

	// 2. Database Setup
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		slog.Error("failed to load database config file", "error", err)
		os.Exit(1)
	}
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		slog.Error("failed to establish connection to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Initialize Layers: Repository -> Service -> gRPC Handler
	repo := repository.NewPostgresRepository(db)
	shortener := shortener.NewBase62Shortener()
	urlSvc := service.NewUrlSvc(repo, shortener)
	handler := urlgrpc.NewURLHandler(urlSvc)

	// 4. Start Network Listener
	port := os.Getenv("URL_SERVICE_PORT")
	if port == "" {
		port = "5051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// sig chan to catch interrupts
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 5. Setup and Start gRPC Server
	server := grpc.NewServer()
	pb.RegisterURLServiceServer(server, handler)
	errChan := make(chan error, 1)
	go func() {
		if err := server.Serve(lis); (err != nil) && (err != grpc.ErrServerStopped) {
			errChan <- err
		}
	}()
	slog.Info("url service up & ready!", "port", port)

	select {
	case <-stopChan:
		slog.Info("shutting down gracefully...")
	case err := <-errChan:
		slog.Error("failed to serve", "error", err)
	}
	server.GracefulStop()
}
