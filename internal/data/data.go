package data

import (
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	slog "log"
	"starland-backend/configs"

	"go.uber.org/zap"

	"os"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewCharacterRepo,
	NewImageModelRepo, NewCharacterAccountLikesRepo,
	NewConversationRepo, NewAccountRepo,
	NewCharacterVoiceRepo)

type Data struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewData(c *configs.Config) *Data {
	return &Data{
		db:  NewDB(c),
		rdb: NewRedis(c),
	}
}

func NewDB(c *configs.Config) *gorm.DB {
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, 
			Colorful:      true,        
			//IgnoreRecordNotFoundError: false,
			LogLevel: logger.Info, // Log lever
		},
	)

	db, err := gorm.Open(mysql.Open(c.Data.DB.Source), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{
		},
	})

	if err != nil {
		zap.S().Errorf("failed opening connection to sqlite: %v", err)
		panic("failed to connect database")
	}

	if err = db.AutoMigrate(&Character{}, &ImageModel{},
		&CharacterAccountLike{}, &Conversation{}, &CharacterVoice{}); err != nil {
		zap.S().Errorf("failed to migrate db: %v", err)
		panic("failed to connect database")
	}

	return db
}

func NewRedis(cfg *configs.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Data.Redis.Host,
		Password: cfg.Data.Redis.Password, // no password set
		DB:       0,
	})
	return rdb
}
