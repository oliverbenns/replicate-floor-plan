package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/replicate/replicate-go"
)

type Service struct {
	replicateClient *replicate.Client
	imagesDir       string
	Logger          *slog.Logger
}

func NewService(replicateClient *replicate.Client, imagesDir string, logger *slog.Logger) *Service {
	return &Service{
		replicateClient: replicateClient,
		imagesDir:       imagesDir,
		Logger:          logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	imageFileNames, err := s.getImageFileNames()
	if err != nil {
		return fmt.Errorf("failed to get image file names: %w", err)
	}

	floorPlans := make([]*FloorPlan, len(imageFileNames))

	for i, fileName := range imageFileNames {
		floorPlan, err := s.getFloorPlan(ctx, fileName)
		if err != nil {
			return fmt.Errorf("failed to get floor plan: %w", err)
		}

		floorPlans[i] = floorPlan

	}

	s.Logger.Info("floor plans", "data", floorPlans)

	return nil
}

type FloorPlan struct {
	SqFt      int `json:"sq_ft"`
	NumFloors int `json:"num_floors"`
}

const modelOwner = "meta"
const modelName = "meta-llama-3-8b"

func (s *Service) getFloorPlan(ctx context.Context, _ string) (*FloorPlan, error) {
	//imageData, err := os.ReadFile(fileName)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	//}

	input := replicate.PredictionInput{
		"prompt": "How many planets are in our solar system?",
	}

	prediction, err := s.replicateClient.CreatePredictionWithModel(ctx, modelOwner, modelName, input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to run model: %w", err)
	}

	err = s.replicateClient.Wait(ctx, prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for prediction: %w", err)
	}

	log.Print("output", prediction.Output)

	return nil, nil

}
