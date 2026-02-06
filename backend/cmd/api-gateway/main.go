package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zipit/internal/gateway/handler"
	"zipit/internal/gateway/router"
	"zipit/pkg/logger"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "zipit/gen/url"
)

func main() {
	logger.SetLogger()
	_ = godotenv.Load() // Optional: only used in local dev, Railway injects env vars directly

	urlServiceHost := getEnvOrDefault("URL_SERVICE_HOST", "localhost")
	urlServicePort := getEnvOrDefault("URL_SERVICE_PORT", "5051")
	urlServiceAddr := urlServiceHost + ":" + urlServicePort

	conn, err := grpc.NewClient(urlServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to connect to url service", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	urlSvc := pb.NewURLServiceClient(conn)
	gatewayHandler := handler.NewGatewayHandler(urlSvc)
	apiRouter := router.New(gatewayHandler)

	port := getEnvOrDefault("GATEWAY_PORT", "8080")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           apiRouter,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); (err != nil) && (err != http.ErrServerClosed) {
			errChan <- err
		}
	}()

	slog.Info("api gateway up & ready", "port", port, "url_service", urlServiceAddr)

	select {
	case <-stopChan:
		slog.Info("shutting down gracefully...")
	case err := <-errChan:
		slog.Error("failed to serve http", "error", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown http server", "error", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
