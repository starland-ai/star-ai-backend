package account

import (
	"context"
	"fmt"
	"starland-backend/internal/biz"
	"starland-backend/internal/pkg/util"
	"time"

	"go.uber.org/zap"
)

func (s *AccountService) LoginSendMail(ctx context.Context, mail string, tokenExpires time.Duration) error {
	code := util.GenValidateCode(6)
	s.mailPool.SendMail(ctx, mail, code)
	if err := s.account.SetLoginCode(ctx, mail, code, tokenExpires); err != nil {
		return fmt.Errorf("LoginSendMail: set login code err: %w", err)
	}

	return nil
}

func (s *AccountService) Auth(ctx context.Context, req *AccountRequest) error {
	bizReq := &biz.AccountRequest{
		AccountID:   req.AccountID,
		Email:       req.Email,
		Name:        req.Name,
		Provider:    req.Provider,
		AvatarURL:   req.AvatarURL,
		ClaimPoints: req.ClaimPoints,
	}
	zap.S().Info("Auth: req: %+v", *bizReq)
	if err := s.account.Auth(ctx, bizReq); err != nil {
		return fmt.Errorf("Auth: auth accout err: %w ", err)
	}
	return nil
}

func (s *AccountService) QueryAccount(ctx context.Context, id string) (*AccountResponse, error) {
	account, err := s.account.QueryAccount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("QueryAccount: query accout by addr(%s) err: %w", id, err)
	}
	return makeBizToAccountResponse(account), nil
}

func (s *AccountService) QueryActivityLogs(ctx context.Context, id string, page, limit int) ([]*ActivityLogResponse, int64, error) {
	account, count, err := s.account.QueryActivityLog(ctx, id, page, limit)
	if err != nil {
		return nil, count, fmt.Errorf("QueryActivityLogs: query activityLog by account(%s) err: %w", id, err)
	}
	return makeBizToActivityLogResponse(account), count, nil
}

func (s *AccountService) Activity(ctx context.Context, req *ActivityRequest) error {
	err := s.account.QueryActivityLimit(ctx, req.AccountID, biz.ActivityCode(req.ActivityCode))
	if err != nil {
		return err
	}

	err = s.account.PostActivity(ctx, req.AccountID, biz.ActivityCode(req.ActivityCode))
	if err != nil {
		return fmt.Errorf("Activity: post activity err: %w ", err)
	}
	return nil
}

func (s *AccountService) ClaimPoints(ctx context.Context, req *ClaimPointsRequest) (string, error) {
	res, err := s.account.ClaimPoints(ctx, &biz.ClaimPointsRequest{
		AccountID: req.AccountID,
		Points:    req.Points,
		IsOk:      req.IsOk,
	})
	if err != nil {
		return "", fmt.Errorf("ClaimPoints: claim points err: %w", err)
	}
	return res, nil
}

func (s *AccountService) SavePointsAddr(ctx context.Context, account, add string) error {
	if err := s.account.SavePointsAddr(ctx, account, add); err != nil {
		return fmt.Errorf("SavePointsAddr: save addr err: %w ", err)
	}
	return nil
}

func makeBizToAccountResponse(req *biz.AccountResponse) *AccountResponse {
	return &AccountResponse{
		AccountID:  req.AccountID,
		AvatarURL:  req.AvatarURL,
		Email:      req.Email,
		Integral:   req.Integral,
		Name:       req.Name,
		Provider:   req.Provider,
		Received:   req.Received,
		SolanaAddr: req.SolanaAddr,
	}
}

func makeBizToActivityLogResponse(req []*biz.ActivityLogResponse) []*ActivityLogResponse {
	res := make([]*ActivityLogResponse, len(req))
	for i := range req {
		res[i] = &ActivityLogResponse{
			CreateAt:     req[i].CreateAt,
			Account:      req[i].AccountID,
			ActivityName: req[i].ActivityName,
			Integral:     req[i].Integral,
		}
	}
	return res
}
