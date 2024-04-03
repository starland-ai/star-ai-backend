package data

import (
	"context"
	"errors"
	"starland-backend/configs"
	"starland-backend/internal/biz"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CharacterVoice struct {
	gorm.Model
	UUID     string `json:"uuid" gorm:"primary_key;size:255"`
	NameZH   string
	NameEN   string
	Gender   int
	ZHUrl    string
	ENUrl    string
	ZHRoleID string
	ENRoleID string
}

type characterVoiceRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewCharacterVoiceRepo(c *configs.Config, data *Data) biz.CharacterVoiceRepo {
	return &characterVoiceRepo{
		cfg:  c,
		data: data,
	}
}

func (r *characterVoiceRepo) CreateCharacterVoice(ctx context.Context, req *biz.CharacterVoice) error {
	var cv = CharacterVoice{
		UUID:   uuid.New().String(),
		NameZH: req.NameZH,
		NameEN: req.NameEN,
		ENUrl:  req.ENUrl,
		ZHUrl:  req.ENUrl,
		Gender: req.Gender,
	}
	return r.data.db.WithContext(ctx).Model(&CharacterVoice{}).Create(&cv).Error
}

func (r *characterVoiceRepo) QueryAllCharacterVoice(ctx context.Context, gender int) ([]*biz.CharacterVoice, error) {
	var cv []*CharacterVoice
	err := r.data.db.WithContext(ctx).Model(&CharacterVoice{}).Model(&CharacterVoice{}).Find(&cv).Error
	if err != nil {
		return nil, err
	}
	return makeBizCharacterVoiceList(cv), nil
}

func (r *characterVoiceRepo) QueryCharacterVoiceByID(ctx context.Context, id string) (*biz.CharacterVoice, error) {
	var cv *CharacterVoice
	err := r.data.db.WithContext(ctx).Model(&CharacterVoice{}).Model(&CharacterVoice{}).Where("uuid = ?", id).First(&cv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return makeBizCharacterVoice(cv), nil
}

func makeBizCharacterVoice(req *CharacterVoice) *biz.CharacterVoice {
	return &biz.CharacterVoice{
		UUID:     req.UUID,
		NameZH:   req.NameZH,
		NameEN:   req.NameEN,
		ENUrl:    req.ENUrl,
		ZHUrl:    req.ENUrl,
		Gender:   req.Gender,
		ZHRoleID: req.ZHRoleID,
		ENRoleID: req.ENRoleID,
	}
}

func makeBizCharacterVoiceList(req []*CharacterVoice) []*biz.CharacterVoice {
	res := make([]*biz.CharacterVoice, len(req))
	for i := range req {
		res[i] = makeBizCharacterVoice(req[i])
	}
	return res
}
