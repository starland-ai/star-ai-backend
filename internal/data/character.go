package data

import (
	"context"
	"errors"
	"starland-backend/configs"
	"starland-backend/internal/biz"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	Unconfirmed int = -1
)

type Character struct {
	gorm.Model
	ID           string `json:"id" gorm:"primary_key;size:255"`
	Name         string
	Gender       int // 0:all 1:man 2:wowem
	Prompt       string
	Introduction string
	AccountID    string
	AccountName  string
	AvatarURL    string
	ImageURL     string
	ImageURLs    datatypes.JSONSlice[string] `gorm:"type:text"`
	LikeCount    int
	ChatCount    int
	IsMint       bool
	Mint         string
	Tag          datatypes.JSONSlice[Tag] `gorm:"type:text"`
	State        int
	VoiceID      string
	IsCustomized bool
	Is3D         bool `json:"is_3d" gorm:"column:is_3d"`
}
type Tag struct {
	Key   string
	Value string
}

type characterRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewCharacterRepo(c *configs.Config, data *Data) biz.CharacterRepo {
	return &characterRepo{
		cfg:  c,
		data: data,
	}
}

func (r *characterRepo) SaveCharacter(ctx context.Context, req *biz.CharacterRequest) (string, error) {
	zap.S().Infof("SaveCharacter: %+v", *req)
	var c *Character
	if err := r.data.db.WithContext(ctx).Model(&Character{}).Where("id = ?", req.ID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag := make([]Tag, 0, len(req.Tags))
			for k, v := range req.Tags {
				tag = append(tag, Tag{Key: k, Value: v})
			}

			c = &Character{
				ID:           req.ID,
				AccountID:    req.AccountID,
				Name:         req.Name,
				Gender:       req.Gender,
				Prompt:       req.Prompt,
				ImageURL:     req.ImageURL,
				AccountName:  req.AccountName,
				AvatarURL:    req.AvatarURL,
				State:        req.State,
				Tag:          tag,
				Is3D:         req.Is3D,
				Introduction: req.Introduction,
				VoiceID:      req.Voice,
				ChatCount:    req.ChatCount,
			}
			if qErr := r.data.db.WithContext(ctx).Model(&Character{}).Create(&c).Error; qErr != nil {
				return "", qErr
			}
			return req.ID, nil
		} else {
			return "", err
		}
	}

	c = &Character{
		ID:           req.ID,
		AccountID:    req.AccountID,
		Name:         req.Name,
		Gender:       req.Gender,
		Prompt:       req.Prompt,
		ImageURL:     req.ImageURL,
		AccountName:  req.AccountName,
		AvatarURL:    req.AvatarURL,
		State:        req.State,
		Is3D:         req.Is3D,
		Introduction: req.Introduction,
		VoiceID:      req.Voice,
		ChatCount:    req.ChatCount,
	}

	zap.S().Infof("save to db req: %+v", *c)
	if err := r.data.db.WithContext(ctx).Model(&Character{}).Where("id = ?", c.ID).Updates(&c).Error; err != nil {
		return "", err
	}
	return req.ID, nil
}

func (r *characterRepo) QueryCharacterByID(ctx context.Context, id string) (*biz.CharacterResponse, error) {
	var c *Character
	if err := r.data.db.WithContext(ctx).Model(&Character{}).Where("id = ?", id).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return makeBizCharacterResponse(c), nil
}

func (r *characterRepo) QueryCharactersByAccountID(ctx context.Context, accountID string, query string,
	page, limit int) ([]*biz.CharacterResponse, int64, error) {
	var (
		res   []*Character
		count int64
	)
	if accountID == "" {
		if err := r.data.db.Model(&Character{}).Where("state != ?", Unconfirmed).Offset((page - 1) * limit).Limit(limit).
			Order("is_customized desc,like_count+chat_count desc,Created_at desc").Find(&res).Error; err != nil {
			return nil, count, err
		}
		if err := r.data.db.Model(&Character{}).Count(&count).Error; err != nil {
			return nil, count, err
		}
	} else {
		queryWhere := "%" + query + "%"
		if err := r.data.db.Model(&Character{}).Where("account_id = ? and (prompt like ? or account_name like ? or  name like ? )",
			accountID, queryWhere, queryWhere, queryWhere).Offset((page - 1) * limit).
			Limit(limit).Order("is_customized desc,like_count+chat_count desc,Created_at desc").Find(&res).Error; err != nil {
			return nil, count, err
		}

		if err := r.data.db.Model(&Character{}).Where("account_id = ? and (prompt like ? or account_name like ? or  name like ? )",
			accountID, queryWhere, queryWhere, queryWhere).Count(&count).Error; err != nil {
			return nil, count, err
		}
	}
	return makeBizCharacterResponses(res), count, nil
}

func (r *characterRepo) QueryCharactersByNameOrPrompt(ctx context.Context, query string,
	page, limit int) ([]*biz.CharacterResponse, int64, error) {
	var (
		res   []*Character
		count int64
	)
	if query == "" {
		if err := r.data.db.Model(&Character{}).Where("state != ?", Unconfirmed).Offset((page - 1) * limit).Limit(limit).
			Order("is_customized desc,like_count+chat_count desc,Created_at desc").Find(&res).Error; err != nil {
			return nil, count, err
		}

		if err := r.data.db.Model(&Character{}).Where("state != ?", Unconfirmed).Count(&count).Error; err != nil {
			return nil, count, err
		}
	} else {
		queryWhere := "%" + query + "%"
		if err := r.data.db.Model(&Character{}).Where("state != ? and (prompt like ? or account_name like ? or  name like ? )",
			Unconfirmed, queryWhere, queryWhere, queryWhere).Offset((page - 1) * limit).
			Limit(limit).Order("is_customized desc,like_count+chat_count desc,Created_at desc").Find(&res).Error; err != nil {
			return nil, count, err
		}

		if err := r.data.db.Model(&Character{}).Where("state != ? and (prompt like ? or account_name like ? or  name like ? )",
			Unconfirmed, queryWhere, queryWhere, queryWhere).Count(&count).Error; err != nil {
			return nil, count, err
		}
	}

	return makeBizCharacterResponses(res), count, nil
}

func (r *characterRepo) CharacterMintSave(ctx context.Context, id, mint string) error {
	return r.data.db.Model(Character{}).WithContext(ctx).Where("id = ?", id).
		Updates(Character{IsMint: true, Mint: mint}).Error
}

func (r *characterRepo) UpdateCharacter(ctx context.Context, req *biz.UpdateCharacterRequest) error {
	var c *Character
	c = &Character{
		Name:         req.Name,
		Introduction: req.Description,
		ImageURL:     req.Image,
		ImageURLs:    req.Images,
		VoiceID:      req.Voice,
	}
	if err := r.data.db.WithContext(ctx).Model(&Character{}).Where("id = ?", req.ID).
		Updates(&c).Error; err != nil {
		return err
	}
	return nil
}

func (r *characterRepo) DeleteCharacterByID(ctx context.Context, id string) error {
	if err := r.data.db.WithContext(ctx).Model(&Character{}).Where("id = ?", id).Delete(&Character{}).Error; err != nil {
		return err
	}
	return nil
}

func makeBizCharacterResponse(c *Character) *biz.CharacterResponse {
	tags := make([]biz.Tag, len(c.Tag))

	for i := range c.Tag {
		tags[i] = biz.Tag{
			Key:   c.Tag[i].Key,
			Value: c.Tag[i].Value,
		}
	}
	return &biz.CharacterResponse{
		ID:           c.ID,
		AccountID:    c.AccountID,
		Name:         c.Name,
		Gender:       c.Gender,
		Prompt:       c.Prompt,
		ImageURL:     c.ImageURL,
		AccountName:  c.AccountName,
		AvatarURL:    c.AvatarURL,
		IsMint:       c.IsMint,
		UpdateTime:   c.UpdatedAt,
		Tag:          tags,
		LikeCount:    c.LikeCount,
		ChatCount:    c.ChatCount,
		Introduction: c.Introduction,
		Mint:         c.Mint,
		Voice:        c.VoiceID,
		ImageURLs:    c.ImageURLs,
		Is3D:         c.Is3D,
	}
}

func makeBizCharacterResponses(req []*Character) []*biz.CharacterResponse {
	res := make([]*biz.CharacterResponse, len(req))
	for i := range req {
		res[i] = makeBizCharacterResponse(req[i])
	}
	return res
}
