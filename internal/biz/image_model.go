package biz

import (
	"context"
	"fmt"
	"starland-backend/configs"
	"starland-backend/internal/pkg/bizerr"
)

type ImageModel struct {
	UUID                     string
	NameEN                   string
	NameZH                   string
	URL                      string
	InferModelType           string
	InferModelName           string
	ComfyuiModelName         string
	InferModelDownloadURL    string
	InferDepModelName        string
	InferDepModelDownloadURL string
	NegativePrompt           string
	SamplerName              string
	CfgScale                 float32
	Steps                    int
	Width                    int
	Height                   int
	BatchSize                int
	ClipSkip                 int
	DenoisingStrength        float32
	Ensd                     int
	HrUpscaler               string
	EnableHr                 bool
	RestoreFaces             bool
	Trigger                  string
	Gender                   int
}

type ImageModelAbbreviate struct {
	UUID   string
	NameEN string
	NameZH string
	URL    string
	Gender int
}

type ImageModelRepo interface {
	QueryAllImageModel(context.Context) ([]*ImageModel, error)
	QueryAllImageModelAbbreviate(context.Context) ([]*ImageModelAbbreviate, error)
	QueryImageModelByID(context.Context, string) (*ImageModel, error)
}

type ImageModelUsecase struct {
	repo ImageModelRepo
	conf *configs.Config
}

func NewImageModelUsecase(conf *configs.Config, repo ImageModelRepo) *ImageModelUsecase {
	return &ImageModelUsecase{repo: repo, conf: conf}
}

func (uc *ImageModelUsecase) QueryAllImageModel(ctx context.Context) ([]*ImageModelAbbreviate, error) {
	ims, err := uc.repo.QueryAllImageModelAbbreviate(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryAllImageModel: query all model err: %w", err))
	}
	return ims, nil
}

func (uc *ImageModelUsecase) QueryImageModelByID(ctx context.Context, id string) (*ImageModel, error) {
	im, err := uc.repo.QueryImageModelByID(ctx, id)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryImageModelByID: query model by id(%s) err: %w", id, err))
	}
	if im == nil {
		return nil, bizerr.ErrModelNotExist
	}
	return im, nil
}

func (uc *ImageModelUsecase) QueryImageModels(ctx context.Context) ([]*ImageModel, error) {
	im, err := uc.repo.QueryAllImageModel(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryImageModels: query model err: %w", err))
	}
	if im == nil {
		return nil, bizerr.ErrModelNotExist
	}
	return im, nil
}
