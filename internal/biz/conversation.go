package biz

import (
	"context"
	"fmt"
	"starland-backend/configs"
	"starland-backend/internal/pkg/bizerr"
	"time"
)

type ConversationRepo interface {
	SaveConversation(context.Context, *ConversationRequest) error
	QueryConversationByID(context.Context, string, string) (*ConversationResponse, error)
	QueryConversationsByAccountID(context.Context, string, int, int) ([]*ConversationResponse, int64, error)
	QueryConversationsCountByAccountID(context.Context, string) (int64, error)
}

type ConversationUsecase struct {
	repo ConversationRepo
	conf *configs.Config
}

type ConversationRequest struct {
	ConversationID string
	AccountID      string
	CharacterID    string
}

type ConversationResponse struct {
	ConversationID string
	AccountID      string
	CharacterID    string
	UpdateTime     time.Time
}

func NewConversationUsecase(conf *configs.Config, repo ConversationRepo) *ConversationUsecase {
	return &ConversationUsecase{repo: repo, conf: conf}
}

func (uc *ConversationUsecase) QueryConversation(ctx context.Context, account, characterID string) (*ConversationResponse, error) {
	cr, err := uc.repo.QueryConversationByID(ctx, account, characterID)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryConversation: query conversation by id err: %w", err))
	}
	return cr, nil
}

func (uc *ConversationUsecase) SaveConversation(ctx context.Context, account, characterID, conversationID string) (string, error) {
	err := uc.repo.SaveConversation(ctx, &ConversationRequest{
		ConversationID: conversationID,
		AccountID:      account,
		CharacterID:    characterID,
	})

	if err != nil {
		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("SaveConversation: save conversation to db err: %w", err))
	}
	return conversationID, nil
}

func (uc *ConversationUsecase) QueryConversations(ctx context.Context, accountID string, page, limit int) ([]*ConversationResponse, int64, error) {
	res, count, err := uc.repo.QueryConversationsByAccountID(ctx, accountID, page, limit)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryConversations: query conversations by account err: %w ", err))
	}
	return res, count, nil
}

func (uc *ConversationUsecase) QueryConversationsCount(ctx context.Context, accountID string) (int64, error) {
	count, err := uc.repo.QueryConversationsCountByAccountID(ctx, accountID)
	if err != nil {
		return count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryConversations: query conversations by account err: %w ", err))
	}
	return count, nil
}
