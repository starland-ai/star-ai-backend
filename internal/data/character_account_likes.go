package data

import (
	"context"
	"errors"
	"starland-backend/configs"
	"starland-backend/internal/biz"

	"gorm.io/gorm"
)

type CharacterAccountLike struct {
	gorm.Model
	CharacterID string `gorm:"primary_key"`
	AccountID   string `gorm:"primary_key"`
	Flag        bool
}

type characterAccountLikesRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewCharacterAccountLikesRepo(c *configs.Config, data *Data) biz.CharacterAccountLikesRepo {
	return &characterAccountLikesRepo{
		cfg:  c,
		data: data,
	}
}

func (r *characterAccountLikesRepo) SaveCharacterAccountLike(ctx context.Context,
	characterID, accountID string, flag bool) error {
	var characterAccountLike *CharacterAccountLike

	if err := r.data.db.WithContext(ctx).Model(&CharacterAccountLike{}).
		Where("character_id=? and account_id=? ", characterID, accountID).First(&characterAccountLike).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	characterAccountLike.Flag = flag
	characterAccountLike.AccountID = accountID
	characterAccountLike.CharacterID = characterID
	if err := r.data.db.WithContext(ctx).Model(&CharacterAccountLike{}).
		Where("character_id=? and account_id=? ", characterID, accountID).Save(&characterAccountLike).Error; err != nil {
		return nil
	}
	return nil
}

func (r *characterAccountLikesRepo) QueryCharacterAccountLike(ctx context.Context,
	characterID, accountID string) (bool, error) {
	var characterAccountLike *CharacterAccountLike

	if err := r.data.db.WithContext(ctx).Model(&CharacterAccountLike{}).
		Where("character_id=? and account_id=? ", characterID, accountID).First(&characterAccountLike).Error; err != nil {
		return false, err
	}
	return characterAccountLike.Flag, nil
}

func (r *characterAccountLikesRepo) QueryCharacterLikeCount(ctx context.Context, characterID string) (int64, error) {
	var count int64
	if err := r.data.db.WithContext(ctx).Model(&CharacterAccountLike{}).
		Where("character_id = ? and flag = true", characterID).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
