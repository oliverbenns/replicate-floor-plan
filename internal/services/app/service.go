package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

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

const modelOwner = "yorickvp"
const modelName = "llava-13b"
const modelVersion = "80537f9eead1a5bfa72d5ac6ea6414379be41d4d4f6679fd776e9535d1eb58bb"

func (s *Service) getFloorPlan(ctx context.Context, fileName string) (*FloorPlan, error) {
	imageData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	encImageData := base64.StdEncoding.EncodeToString(imageData)

	input := replicate.PredictionInput{
		"image":  "data:image/jpeg;base64," + encImageData,
		"prompt": `In this image, extract the information about the number of floors and the square footage of the building. The output should solely be a valid json object that is the following schema: {"sq_ft": number, "num_floors": number}. Do not escape the text`,
	}

	prediction, err := s.replicateClient.CreatePrediction(ctx, modelVersion, input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create prediction: %w", err)
	}

	err = s.replicateClient.Wait(ctx, prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for prediction: %w", err)
	}

	data := ""
	for _, v := range prediction.Output.([]interface{}) {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert prediction to []byte")
		}
		data += str
	}

	floorPlan := FloorPlan{}
	err = json.Unmarshal([]byte(data), &floorPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w, %v", err, prediction.Output)
	}

	return &floorPlan, nil
}
