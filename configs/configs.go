package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/viper"
)

var cfg *Config

func InitConfig() {
	confPath := "./conf"
	value := os.Getenv("CONF_PATH")
	if value != "" {
		confPath = value
	}

	cf := viper.GetViper()
	cf.AddConfigPath(confPath)
	cf.SetConfigName("config.yaml")
	cf.SetConfigType("yaml")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("InitLog: config has change: %v\n", e.Name)
		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Printf("InitLog: unmarshal config failed: %v\n", err)
			return
		}
	})

	if err := cf.ReadInConfig(); err != nil {
		fmt.Printf("InitLog:  reading cf file: %v\n", err)
		os.Exit(-1)
	} else {
		fmt.Println("InitLog: using cf file:", cf.ConfigFileUsed())
	}

	if err := cf.Unmarshal(&cfg); err != nil {
		fmt.Printf("InitLog:  unmarshaling cf file: %v\n", err)
		os.Exit(1)
	} else if GetConfig().Debug {
		// un
		_, err := pp.Println(cfg)
		if err != nil {
			fmt.Printf("InitLog: pp err %s\n", err.Error())
		}
	}
}

func GetConfig() *Config {
	if cfg == nil {
		InitConfig()
	}
	return cfg
}

type Config struct {
	Debug           bool                  `mapstructure:"debug"`
	Env             string                `mapstructure:"env"`
	RedirectURL     string                `mapstructure:"redirectUrl"`
	HTTP            *HTTPConfig           `mapstructure:"http"`
	HTTPS           *HTTPSConfig          `mapstructure:"https"`
	Log             *LogConfig            `mapstructure:"log"`
	Data            *DataConfig           `mapstructure:"data"`
	Account         *AccountServiceConfig `mapstructure:"account"`
	ChatCompletions *AgentEndpointConfig  `mapstructure:"chatCompletions"`
	Voice           *AgentEndpointConfig  `mapstructure:"voice"`
	Chat            *AgentEndpointConfig  `mapstructure:"chat"`
	FeiShuAlertURL  string                `mapstructure:"feiShuAlertUrl"`
	File            FileConfig            `mapstructure:"file"`
	Login           *LoginConfig          `mapstructure:"login"`
}

type HTTPConfig struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	LogFormat    string        `mapstructure:"log_format"`
}

type HTTPSConfig struct {
	Addr         string        `mapstructure:"addr"`
	PemPtah      string        `mapstructure:"pem"`
	KeyPath      string        `mapstructure:"key"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	ColorLevel bool   `mapstructure:"color_level"`
}

type DataConfig struct {
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
}

type DBConfig struct {
	Source string `mapstructure:"source"`
}

type RedisConfig struct {
	Host       string        `mapstructure:"host"`
	Password   string        `mapstructure:"password"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type AgentConfig struct {
	MaxChatHistoryContextLength int32   `mapstructure:"maxChatHistoryContextLength"`
	Temperature                 float32 `mapstructure:"temperature"`
}

type AccountServiceConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Token    string `mapstructure:"token"`
}

type AgentEndpointConfig struct {
	ChatLimit int    `mapstructure:"chatLimit"`
	Endpoint  string `mapstructure:"endpoint"`
}

type FileConfig struct {
	ImagesEndpoint string `mapstructure:"imagesEndpoint"`
	UploadSaveDir  string `mapstructure:"uploadSaveDir"`
	ImagePath      string `mapstructure:"imagePath"`
	VoiceEndpoint  string `mapstructure:"voiceEndpoint"`
	VoicePath      string `mapstructure:"voicePath"`
}

type LoginConfig struct {
	RedirectURL string      `mapstructure:"redirect_url"`
	Mail        *MailConfig `mapstructure:"mail"`
}

type MailConfig struct {
	SigninURL   string                  `mapstructure:"signin_url"`
	PactURL     string                  `mapstructure:"pact_url"`
	MailAccount map[string]*MailAccount `mapstructure:"mail_map"`
}

type MailAccount struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
