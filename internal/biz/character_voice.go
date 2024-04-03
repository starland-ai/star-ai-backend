package biz

import (
	"context"
	"fmt"
	"starland-backend/configs"
	"starland-backend/internal/pkg/bizerr"
)

type CharacterVoice struct {
	UUID   string
	NameZH string
	NameEN string
	Gender int
	ZHUrl  string
	ENUrl  string
	ZHRoleID string
	ENRoleID string
}

func NewCharacterVoiceUsecase(conf *configs.Config, repo CharacterVoiceRepo) *CharacterVoiceUsecase {
	return &CharacterVoiceUsecase{repo: repo}
}

type CharacterVoiceRepo interface {
	QueryAllCharacterVoice(context.Context, int) ([]*CharacterVoice, error)
	CreateCharacterVoice(context.Context, *CharacterVoice) error
	QueryCharacterVoiceByID(context.Context, string) (*CharacterVoice, error)
}

type CharacterVoiceUsecase struct {
	repo CharacterVoiceRepo
}

func (uc *CharacterVoiceUsecase) QueryAllCharacterVoice(ctx context.Context) ([]*CharacterVoice, error) {
	res, err := uc.repo.QueryAllCharacterVoice(ctx, 0)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryAllCharacterVoice: query all voice err: %w", err))
	}
	return res, nil
}

func (uc *CharacterVoiceUsecase) QueryCharacterVoice(ctx context.Context, id string) (*CharacterVoice, error) {
	res, err := uc.repo.QueryCharacterVoiceByID(ctx, id)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharacterVoice: query all voice err: %w", err))
	}
	if res == nil {
		return nil, bizerr.ErrVoiceNotExist
	}
	return res, nil
}
