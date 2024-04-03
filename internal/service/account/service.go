package account

import (
	"starland-backend/configs"
	"starland-backend/internal/biz"
	"time"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAccountService)

type AccountService struct {
	cfg      *configs.Config
	account  *biz.AccountAndActivitySerClientUsecase
	mailPool *biz.MailPool
}

func NewAccountService(cfg *configs.Config, account *biz.AccountAndActivitySerClientUsecase) *AccountService {
	s := &AccountService{cfg: cfg, account: account, mailPool: biz.NewMailPool(cfg)}
	return s
}

type AccountRequest struct {
	AccountID   string
	Email       string
	Name        string
	Provider    string
	AvatarURL   string
	ClaimPoints string
}

type ActivityRequest struct {
	AccountID    string
	ActivityCode int
}

type AccountResponse struct {
	AccountID  string `json:"account_id"`
	AvatarURL  string `json:"avatar_url"`
	Email      string `json:"email"`
	Integral   int    `json:"integral"`
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Received   int    `json:"received"`
	SolanaAddr string `json:"solana_addr"`
}

type ActivityLogResponse struct {
	CreateAt     time.Time `json:"create_at"`
	Account      string    `json:"account"`
	ActivityName string    `json:"activity_name"`
	Integral     int       `json:"integral"`
}

type ClaimPointsRequest struct {
	AccountID string `json:"account_id"`
	Points    int    `json:"points"`
	IsOk      bool   `json:"is_ok"`
}
