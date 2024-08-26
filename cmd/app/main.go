package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/oliverbenns/replicate-floor-plan/internal/services/app"
	"github.com/replicate/replicate-go"
)

func main() {
	ctx := context.Background()

	replicateApiToken := os.Getenv("REPLICATE_API_TOKEN")

	replicateClient, err := replicate.NewClient(replicate.WithToken(replicateApiToken))
	if err != nil {
		panic(err)
	}

	imagesDir := os.Getenv("IMAGES_DIR")

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	svc := app.NewService(replicateClient, imagesDir, logger)

	err = svc.Run(ctx)
	if err != nil {
		panic(err)
	}
}
