package data

import (
	"context"
	"errors"
	"starland-backend/configs"
	"starland-backend/internal/biz"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ConversationID string `json:"conversation_id" gorm:"primary_key;size:255"`
	AccountID      string
	CharacterID    string
}

type conversationRepo struct {
	cfg  *configs.Config
	data *Data
}

func NewConversationRepo(c *configs.Config, data *Data) biz.ConversationRepo {
	return &conversationRepo{
		cfg:  c,
		data: data,
	}
}

func (r *conversationRepo) QueryConversationByID(ctx context.Context, account, characterID string) (*biz.ConversationResponse, error) {
	var res *Conversation
	if err := r.data.db.WithContext(ctx).Model(&Conversation{}).Where("account_id = ? and character_id = ?", account, characterID).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &biz.ConversationResponse{
		ConversationID: res.ConversationID,
		AccountID:      res.AccountID,
		CharacterID:    res.CharacterID,
		UpdateTime:     res.UpdatedAt,
	}, nil
}

func (r *conversationRepo) QueryConversationsByAccountID(ctx context.Context, accountID string, page, limit int) ([]*biz.ConversationResponse, int64, error) {
	var (
		res   []*Conversation
		count int64
	)
	if err := r.data.db.WithContext(ctx).Model(&Conversation{}).Where("account_id = ?", accountID).Order("updated_at desc ").Find(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, count, nil
		}
		return nil, count, err
	}

	if err := r.data.db.WithContext(ctx).Model(&Conversation{}).Where("account_id = ?", accountID).Count(&count).Error; err != nil {
		return nil, count, err
	}

	return makeBizConversationResponse(res), count, nil
}

func (r *conversationRepo) SaveConversation(ctx context.Context, req *biz.ConversationRequest) error {
	var con *Conversation
	if req.ConversationID == "" {
		con = &Conversation{
			ConversationID: uuid.NewString(),
			AccountID:      req.AccountID,
			CharacterID:    req.CharacterID,
		}
		if err := r.data.db.WithContext(ctx).Model(&Conversation{}).Where("conversation_id = ?", req.ConversationID).Create(&con).Error; err != nil {
			return err
		}
	} else {
		if err := r.data.db.WithContext(ctx).Model(&Conversation{}).Where("conversation_id = ?", req.ConversationID).
			Updates(Conversation{ConversationID: req.ConversationID, CharacterID: req.CharacterID, AccountID: req.AccountID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func makeBizConversationResponse(req []*Conversation) []*biz.ConversationResponse {
	res := make([]*biz.ConversationResponse, len(req))
	for i := range res {
		res[i] = &biz.ConversationResponse{
			ConversationID: req[i].ConversationID,
			AccountID:      req[i].AccountID,
			CharacterID:    req[i].CharacterID,
			UpdateTime:     req[i].UpdatedAt,
		}
	}
	return res
}

func (r *conversationRepo) QueryConversationsCountByAccountID(ctx context.Context, id string) (int64, error) {
	var count int64

	if err := r.data.db.WithContext(ctx).Model(&Conversation{}).Where("character_id = ?", id).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
