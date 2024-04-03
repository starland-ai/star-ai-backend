package data

import (
	"context"
	config "starland-backend/configs"
	"starland-backend/internal/biz"
	"time"

	"go.uber.org/zap"
)

type accountRepo struct {
	cfg  *config.Config
	data *Data
}

func NewAccountRepo(c *config.Config, data *Data) biz.AccountRepo {
	return &accountRepo{
		cfg:  c,
		data: data,
	}
}

func (r *accountRepo) SetLoginCode(ctx context.Context, email, code string, expiration time.Duration) error {
	if res, err := r.data.rdb.WithContext(ctx).Set(email, code, expiration).Result(); err != nil {
		return err
	} else {
		zap.S().Infof("SetLoginCode: res: %s", res)
		return nil
	}
}

func (r *accountRepo) GetLoginCode(ctx context.Context, email string) (string, error) {
	if res, err := r.data.rdb.WithContext(ctx).Get(email).Result(); err != nil {
		return "", err
	} else {
		zap.S().Infof("GetLoginCode: res: %s", res)
		return res, nil
	}
}

func (r *accountRepo) DelLoginCode(ctx context.Context, mail string) error {
	if _, err := r.data.rdb.WithContext(ctx).Del(mail).Result(); err != nil {
		return err
	} else {
		zap.S().Infof("DelLoginCode: mail: %s", mail)
		return nil
	}
}
