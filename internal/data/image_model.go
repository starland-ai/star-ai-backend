package data

import (
	"context"
	"errors"
	config "starland-backend/configs"
	"starland-backend/internal/biz"

	"gorm.io/gorm"
)

type ImageModel struct {
	gorm.Model
	UUID                     string `json:"uuid" gorm:"primary_key;size:255"`
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
	Gender                   int // 1: man,2:women,0:all
}

type imageModelRepo struct {
	cfg  *config.Config
	data *Data
}

func NewImageModelRepo(c *config.Config, data *Data) biz.ImageModelRepo {
	return &imageModelRepo{
		cfg:  c,
		data: data,
	}
}

func (r *imageModelRepo) QueryAllImageModelAbbreviate(ctx context.Context) ([]*biz.ImageModelAbbreviate, error) {
	var ims []*ImageModel
	err := r.data.db.Model(&ImageModel{}).Find(&ims).Error
	if err != nil {
		return nil, err
	}
	return makeBizImageModelAbbreviates(ims), nil
}

func (r *imageModelRepo) QueryAllImageModel(ctx context.Context) ([]*biz.ImageModel, error) {
	var ims []*ImageModel
	err := r.data.db.Model(&ImageModel{}).Find(&ims).Error
	if err != nil {
		return nil, err
	}
	return makeBizImageModels(ims), nil
}

func (r *imageModelRepo) QueryImageModelByID(ctx context.Context, id string) (*biz.ImageModel, error) {
	var im *ImageModel
	err := r.data.db.First(&im, "uuid = ?", id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return makeBizImageModel(im), nil
}

func makeBizImageModel(im *ImageModel) *biz.ImageModel {
	return &biz.ImageModel{
		UUID:                     im.UUID,
		NameEN:                   im.NameEN,
		NameZH:                   im.NameZH,
		URL:                      im.URL,
		InferModelType:           im.InferModelType,
		InferModelName:           im.InferModelName,
		ComfyuiModelName:         im.ComfyuiModelName,
		InferModelDownloadURL:    im.InferModelDownloadURL,
		InferDepModelName:        im.InferDepModelName,
		InferDepModelDownloadURL: im.InferDepModelDownloadURL,
		NegativePrompt:           im.NegativePrompt,
		SamplerName:              im.SamplerName,
		CfgScale:                 im.CfgScale,
		Steps:                    im.Steps,
		Width:                    im.Width,
		Height:                   im.Height,
		BatchSize:                im.BatchSize,
		ClipSkip:                 im.ClipSkip,
		DenoisingStrength:        im.DenoisingStrength,
		Ensd:                     im.Ensd,
		HrUpscaler:               im.HrUpscaler,
		EnableHr:                 im.EnableHr,
		RestoreFaces:             im.RestoreFaces,
		Trigger:                  im.Trigger,
		Gender:                   im.Gender,
	}
}

func makeBizImageModels(ims []*ImageModel) []*biz.ImageModel {
	res := make([]*biz.ImageModel, len(ims))
	for i := range ims {
		res[i] = makeBizImageModel(ims[i])
	}
	return res
}

func makeBizImageModelAbbreviates(ims []*ImageModel) []*biz.ImageModelAbbreviate {
	res := make([]*biz.ImageModelAbbreviate, len(ims))
	for i := range ims {
		res[i] = &biz.ImageModelAbbreviate{
			UUID:   ims[i].UUID,
			NameEN: ims[i].NameEN,
			NameZH: ims[i].NameZH,
			URL:    ims[i].URL,
			Gender: ims[i].Gender,
		}
	}
	return res
}
