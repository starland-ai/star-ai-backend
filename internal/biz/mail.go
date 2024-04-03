package biz

import (
	"bytes"
	"context"
	"fmt"
	"starland-backend/configs"
	"starland-backend/internal/pkg/alert"
	"sync"
	"text/template"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var (
	signinURL string
)

type MailPool struct {
	Index       int
	mu          sync.Mutex
	config      *configs.Config
	MailClients []MailClient
}

type MailClientInfo struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	PassWord string `mapstructure:"password"`
}

type MailClient struct {
	Info   *MailClientInfo
	Dialer *gomail.Dialer
}

func NewMailPool(cfg *configs.Config) *MailPool {
	signinURL = cfg.Login.Mail.SigninURL
	return &MailPool{
		Index:       0,
		MailClients: NewMailClients(cfg),
		config:      cfg,
	}
}

func NewMailClients(cfg *configs.Config) []MailClient {
	mailClients := make([]MailClient, 0, len(cfg.Login.Mail.MailAccount))
	for _, item := range cfg.Login.Mail.MailAccount {
		info := MailClientInfo{
			Host:     item.Host,
			Port:     item.Port,
			UserName: item.Username,
			PassWord: item.Password,
		}
		mailClients = append(mailClients, MailClient{
			Info:   &info,
			Dialer: gomail.NewDialer(info.Host, info.Port, info.UserName, info.PassWord),
		})
	}
	return mailClients
}

func (mp *MailPool) SendMail(ctx context.Context, toMail, code string) {
	zap.S().Infof("send tomail(%s) start", toMail)
	index := mp.Index % len(mp.MailClients)
	mp.mu.Lock()
	defer mp.mu.Unlock()
	msg := mp.MailClients[index].makeMessage( toMail, code)
	go func() {
		if err := mp.MailClients[index].Dialer.DialAndSend(msg); err != nil {
			zap.S().Errorf(fmt.Sprintf("%s send email(%s) failed err: %s", mp.MailClients[index].Info.UserName, toMail, err.Error()))
			alert.SendAlertMsg(fmt.Sprintf("%s 需要重新验证", mp.MailClients[index].Info.UserName))
		} else {
			zap.S().Info(fmt.Sprintf("%s send email(%s) success", mp.MailClients[index].Info.UserName, toMail))
		}
	}()
	mp.Index++
}

func (mc *MailClient) makeMessage(toMail, code string) *gomail.Message {
	/*
		localizer := i18n.NewLocalizer(bundle, lang)

		 // Set title message.
		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HelloPerson",
				Other: "Hello {{.Name}}",
			},
			TemplateData: map[string]string{
				"Name": toMail,
			},
		})
		subject := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "Subject",
			},
		})
		description := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "Description",
			},
		})
	*/
	helloPerson := "Hello {{.Name}}"
	subject := "Your verification code"
	description := "It appears that you are attempting to log in using a new device. Here is the token verification code required for accessing your account:"
	m := gomail.NewMessage()
	m.SetHeader("From", mc.Info.UserName)
	m.SetHeader("To", toMail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", readTemplate(code, helloPerson, description))
	return m
}

func readTemplate(code, helloPerson, description string) string {
	body, err := template.ParseFiles("./conf/template/index.html")
	if err != nil {
		zap.S().Error(err)
		return ""
	}
	dataTemplate := struct {
		Code string
	}{
		Code: code,
	}
	buf := new(bytes.Buffer)
	err = body.Execute(buf, dataTemplate)
	if err != nil {
		zap.S().Error(err)
	}
	return buf.String()
}
