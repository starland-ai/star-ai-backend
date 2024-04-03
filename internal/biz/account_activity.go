package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"starland-backend/configs"
	"starland-backend/internal/pkg/bizerr"
	"starland-backend/internal/pkg/httpclientutil"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	queryAccountURL       = "/v1/account/"
	authAccountURL        = "/v1/account/"
	activityURL           = "/v1/activity"
	queryActivityLogURL   = "/v1/activity/log/"
	queryActivityLimitURL = "/v1/activity/Limit"

	claimPointsURL    = "/v1/account/claim_points"
	savePointsAddrURL = "/v1/account/%s/save_points_addr"
)

type ActivityCode int

const (
	Chat            ActivityCode = 0
	Like            ActivityCode = 10001
	CreateCharacter ActivityCode = 10002
	Login           ActivityCode = 10003
	LimitCode                    = "100"
)

type AccountRequest struct {
	AccountID   string `json:"account_id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	AvatarURL   string `json:"avatar_url"`
	ClaimPoints string `json:"claim_points"`
}

type ClaimPointsRequest struct {
	AccountID string `json:"account_id"`
	Points    int    `json:"points"`
	IsOk      bool   `json:"is_ok"`
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

type ActivityLogRequest struct {
	AccountID    string
	ActivityCode int
	ActivityName string
	Integral     int
}

type ActivityLogResponse struct {
	AccountID    string
	ActivityCode int
	ActivityName string
	Integral     int
	CreateAt     time.Time
}

type AccountAndActivitySerClientUsecase struct {
	conf *configs.AccountServiceConfig
	repo AccountRepo
}

func NewAccountUsecase(conf *configs.Config, repo AccountRepo) *AccountAndActivitySerClientUsecase {
	return &AccountAndActivitySerClientUsecase{conf: conf.Account, repo: repo}
}

type AccountRepo interface {
	SetLoginCode(context.Context, string, string, time.Duration) error
	GetLoginCode(context.Context, string) (string, error)
	DelLoginCode(context.Context, string) error
}

type AccountUsecase struct {
	repo AccountRepo
}

func (uc *AccountAndActivitySerClientUsecase) QueryAccount(ctx context.Context, account string) (*AccountResponse, error) {
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s%s", endpoint, queryAccountURL, account)
	req := &fasthttp.Request{}
	req.SetRequestURI(host)
	req.Header.SetContentType("application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	resp := &fasthttp.Response{}
	cl := &fasthttp.Client{}
	if err := cl.Do(req, resp); err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLimit: cl.do err: %w", err))
	}

	defer func() {
		resp.ConnectionClose()
	}()

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLimit: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Data struct {
			AccountID  string `json:"account_id"`
			AvatarURL  string `json:"avatar_url"`
			Email      string `json:"email"`
			Integral   int    `json:"integral"`
			Name       string `json:"name"`
			Provider   string `json:"provider"`
			Received   int    `json:"received"`
			SolanaAddr string `json:"solana_addr"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryAccount: response json decode err: %w", err))
	}
	if res.Code != "0" {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryAccount: response is failed : %+v", res))
	}

	accountResponse := &AccountResponse{
		AccountID:  res.Data.AccountID,
		AvatarURL:  res.Data.AvatarURL,
		Email:      res.Data.Email,
		Integral:   res.Data.Integral,
		Name:       res.Data.Name,
		Provider:   res.Data.Provider,
		Received:   res.Data.Received,
		SolanaAddr: res.Data.SolanaAddr,
	}
	return accountResponse, nil
}

func (uc *AccountAndActivitySerClientUsecase) Auth(ctx context.Context, reqData *AccountRequest) error {
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s", endpoint, authAccountURL)

	zap.S().Infof("reqData:%+v", reqData)
	reqBuf := new(bytes.Buffer)
	err := json.NewEncoder(reqBuf).Encode(reqData)
	if err != nil {
		return fmt.Errorf("req encode: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", host, reqBuf)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: post authAccountURL err: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: cl.do err: %w", err))
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != 200 {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: response json decode err: %w", err))
	}
	if res.Code != "0" {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: response is failed : %+v", res))
	}

	return nil
}

func (uc *AccountAndActivitySerClientUsecase) QueryActivityLog(ctx context.Context, account string, page, limit int) ([]*ActivityLogResponse, int64, error) {
	var count int64
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s%s?page=%d&limit=%d", endpoint, queryActivityLogURL, account, page, limit)
	var req *http.Request
	req, err := http.NewRequestWithContext(ctx, "GET", host, nil)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLog: get queryActivityLogURL: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLog: cl.do err: %w", err))
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != 200 {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLog: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Data struct {
			Count int64 `json:"count"`
			Data  []struct {
				Account      string    `json:"account,omitempty"`
				ActivityName string    `json:"activity_name,omitempty"`
				CreateAt     time.Time `json:"create_at,omitempty"`
				Integral     int       `json:"integral,omitempty"`
			} `json:"data"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLog: response json decode err: %w", err))
	}
	if res.Code != "0" {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLog: response is failed : %+v", res))
	}

	response := make([]*ActivityLogResponse, len(res.Data.Data))
	for i := range res.Data.Data {
		response[i] = &ActivityLogResponse{
			AccountID:    res.Data.Data[i].Account,
			ActivityName: res.Data.Data[i].ActivityName,
			Integral:     res.Data.Data[i].Integral,
			CreateAt:     res.Data.Data[i].CreateAt,
		}
	}
	count = res.Data.Count
	return response, count, nil
}

func (uc *AccountAndActivitySerClientUsecase) PostActivity(ctx context.Context, account string, activityCode ActivityCode) error {
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s", endpoint, activityURL)

	var (
		reqData struct {
			ActivityCode ActivityCode `json:"activity_code"`
			Account      string       `json:"account"`
		}
	)
	reqData.ActivityCode = activityCode
	reqData.Account = account
	zap.S().Infof("reqData:%+v", reqData)
	reqBuf := new(bytes.Buffer)
	err := json.NewEncoder(reqBuf).Encode(reqData)
	if err != nil {
		return fmt.Errorf("req encode: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", host, reqBuf)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("PostActivity: post activityURL err: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: cl.do err: %w", err))
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != 200 {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("PostActivity: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Data string `json:"data"`
		Msg  string `json:"msg"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("PostActivity: response json decode err: %w", err))
	}
	zap.S().Infof("res: %+v", res)

	if res.Code != "0" {
		if res.Code == LimitCode {
			return parseActivityCodeToErr(int(activityCode))
		}
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("PostActivity: response is failed : %+v", res))
	}

	return nil
}

func (uc *AccountAndActivitySerClientUsecase) QueryActivityLimit(ctx context.Context, account string, activityCode ActivityCode) error {
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s?activity_code=%d&account=%s", endpoint, queryActivityLimitURL, activityCode, account)

	req := &fasthttp.Request{}
	req.SetRequestURI(host)
	req.Header.SetContentType("application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	resp := &fasthttp.Response{}
	cl := &fasthttp.Client{}
	if err := cl.Do(req, resp); err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLimit: cl.do err: %w", err))
	}

	defer func() {
		resp.ConnectionClose()
	}()

	if resp.StatusCode() != fasthttp.StatusOK {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLimit: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Data struct {
			IsLimit bool `json:"is_limit"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLimit: response json decode err: %w", err))
	}
	zap.S().Infof("res: %+v", res)

	if res.Code != "0" {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryActivityLimit: response is failed : %+v", res))
	}

	if res.Data.IsLimit {
		return parseActivityCodeToErr(int(activityCode))
	}
	return nil
}

func (uc *AccountAndActivitySerClientUsecase) ClaimPoints(ctx context.Context, reqData *ClaimPointsRequest) (string, error) {
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s", endpoint, claimPointsURL)

	zap.S().Infof("reqData:%+v", reqData)
	reqBuf := new(bytes.Buffer)
	err := json.NewEncoder(reqBuf).Encode(reqData)
	if err != nil {
		return "", fmt.Errorf("req encode: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", host, reqBuf)
	if err != nil {
		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("ClaimPoints: post savePointsAddrURL err: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("ClaimPoints: cl.do err: %w", err))
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != 200 {
		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("ClaimPoints: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("ClaimPoints: response json decode err: %w", err))
	}
	if res.Code != "0" {
		if res.Msg == "Not enough points" {
			return "", bizerr.ErrNotEnoughPoints
		}

		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("ClaimPoints: response is failed : %+v", res))
	}

	return res.Data, nil
}

func (uc *AccountAndActivitySerClientUsecase) SavePointsAddr(ctx context.Context, account, addr string) error {
	endpoint := uc.conf.Endpoint
	host := fmt.Sprintf("%s%s", endpoint, fmt.Sprintf(savePointsAddrURL, account))
	zap.S().Infof("SavePointsAddr: host:%s", host)
	var reqData struct {
		Addr    string `json:"addr"`
		Account string `json:"account"`
	}
	reqData.Addr = addr
	reqData.Account = account
	zap.S().Infof("reqData:%+v", reqData)
	reqBuf := new(bytes.Buffer)
	err := json.NewEncoder(reqBuf).Encode(reqData)
	if err != nil {
		return fmt.Errorf("req encode: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", host, reqBuf)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: post savePointsAddrURL err: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", uc.conf.Token)
	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: cl.do err: %w", err))
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != 200 {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: response status code is %d ", resp.StatusCode))
	}

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: response json decode err: %w", err))
	}
	if res.Code != "0" {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Auth: response is failed : %+v", res))
	}

	return nil
}

func parseActivityCodeToErr(activityCode int) error {
	switch activityCode {
	case int(Chat):
		return bizerr.NewBizError("You've reached the limit for chatting. Please come back tomorrow.", bizerr.Limit)
	case int(Like):
		return bizerr.NewBizError("You've reached the limit for liking. Please come back tomorrow.", bizerr.Limit)
	case int(CreateCharacter):
		return bizerr.NewBizError("You've reached the limit for creating Avatars. Please come back tomorrow.", bizerr.Limit)
	default:
		return bizerr.ErrLimit
	}
}

func (uc *AccountAndActivitySerClientUsecase) CheckLoginCode(ctx context.Context, email, code string) error {
	res, err := uc.repo.GetLoginCode(ctx, email)
	if err != nil {
		return fmt.Errorf("CheckLoginCode: get login code err: %w", err)
	}
	if res == code {
		var retry int
		for retry < 3 {
			err = uc.repo.DelLoginCode(ctx, email)
			if err != nil {
				retry++
				continue
			}
			break
		}
		if err != nil {
			return fmt.Errorf("CheckLoginCode: del login code err: %w", err)
		}

		return nil
	}
	return bizerr.ErrVerificationCodeFailed
}

func (uc *AccountAndActivitySerClientUsecase) SetLoginCode(ctx context.Context, email, code string, tokenExpires time.Duration) error {
	err := uc.repo.SetLoginCode(ctx, email, code, tokenExpires)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("SetLoginCode: set login code err: %w", err))
	}
	return nil
}
