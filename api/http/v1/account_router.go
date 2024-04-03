package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"starland-backend/configs"
	mycookie "starland-backend/internal/pkg/cookie"
	"starland-backend/internal/pkg/middlewares"
	"starland-backend/internal/pkg/util"
	"sync"

	"starland-backend/internal/service/account"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitterv2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	mailOnce sync.Once
	storage  *memory.Storage
	re       *regexp.Regexp
)

type AccountHTTPServer interface {
	LoginSendMail(ctx context.Context, email string, tokenExpires time.Duration) error
	Auth(context.Context, *account.AccountRequest) error
	QueryAccount(context.Context, string) (*account.AccountResponse, error)
	Activity(context.Context, *account.ActivityRequest) error
	ClaimPoints(context.Context, *account.ClaimPointsRequest) (string, error)
	QueryActivityLogs(context.Context, string, int, int) ([]*account.ActivityLogResponse, int64, error)
	SavePointsAddr(context.Context, string, string) error
}

func InitOAuth2(cfg *configs.Config) {
	v := viper.GetViper()
	key := v.GetString("oauth.secret") // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30               // 30 days

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = false

	gothic.Store = store
	goth.UseProviders(
		google.New(v.GetString("oauth.google.client_id"),
			v.GetString("oauth.google.key"),
			v.GetString("oauth.google.callback_url"),
			"email", "profile"),
		twitterv2.New(v.GetString("oauth.twitterv.key"),
			v.GetString("oauth.twitterv.secret"),
			v.GetString("oauth.twitterv.callback_url")),
	)
}

func InitAccountRouter(app fiber.Router, service AccountHTTPServer, conf *configs.Config) {
	router := app.Group("/v1")
	router.Post("/account", auth(service))

	authRouter := router.Group("/account", middlewares.JwtParse())
	authRouter.Get("/:id", queryAccount(service))

	authRouter.Put("/:id", queryAccount(service))
	authRouter.Get("/:id/activity_log", queryAccountLog(service))
	authRouter.Post("/:id/activity", activity(service))
	authRouter.Post("/:id/points", savePointsAddr(service))

	authRouter.Post("/claim_points", claimPoints(service))

	router.Get("/auth/:provider/callback", adaptor.HTTPHandlerFunc(authLoginCallback(service)))
	router.Get("/auth/:provider", adaptor.HTTPHandlerFunc(authLogin()))
	router.Post("/auth/email", sendSigninMail(service))
	InitOAuth2(conf)
}

func auth(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID string `json:"id"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		accountInfo := &account.AccountRequest{
			AccountID: req.ID,
			Provider:  "Blockchain",
		}

		if err := service.Auth(ctx.Context(), accountInfo); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		t, err := middlewares.NewJwtToken(accountInfo.AccountID, middlewares.Expires*time.Hour)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		cookie := new(fiber.Cookie)
		cookie.Name = mycookie.StarlandAIToken
		cookie.Value = t
		cookie.Path = "/"
		cookie.Domain = mycookie.Domain
		ctx.Cookie(cookie)
		zap.S().Info(t)
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(t))
	}
}

func queryAccount(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID string `params:"id"`
			}
		)
		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		response, err := service.QueryAccount(ctx.Context(), req.ID)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(response))
	}
}

func queryAccountLog(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Account string `params:"id"`
				Page    int    `query:"page"`
				Limit   int    `query:"limit"`
			}
			res struct {
				Data  []*account.ActivityLogResponse `json:"data"`
				Count int64                          `json:"count"`
			}
		)

		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := ctx.QueryParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		zap.S().Infof("queryActivityLogs: req: %+v", req)
		response, count, err := service.QueryActivityLogs(ctx.Context(), req.Account, req.Page, req.Limit)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.Count = count
		res.Data = response
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func authLoginCallback(server AccountHTTPServer) func(w http.ResponseWriter, r *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		req = util.FiberGothAdapter(req)

		zap.S().Info("authLoginCallback cookie", req.Cookies())

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			res.WriteHeader(http.StatusOK)
			bytes, e := json.Marshal(util.MakeResponseWithMsg(err.Error()))
			if e != nil {
				zap.S().Error(fmt.Errorf("json Marshal is failed: %w", e))
			}
			_, err = res.Write(bytes)
			if err != nil {
				zap.S().Error(fmt.Errorf("request write CompleteUserAuth err is failed: %w", err))
			}
			return
		}

		accountInfo := &account.AccountRequest{
			Name:      user.Name,
			Provider:  user.Provider,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		}

		err = server.Auth(context.Background(), accountInfo)
		if err != nil {
			res.WriteHeader(http.StatusOK)
			bytes, e := json.Marshal(util.MakeResponseWithMsg(err.Error()))
			if e != nil {
				zap.S().Error(fmt.Errorf("json Marshal is failed: %w", e))
			}
			_, err = res.Write(bytes)
			if err != nil {
				zap.S().Error(fmt.Errorf("request write CompleteUserAuth err is failed: %w", err))
			}
			return
		}

		t, err := middlewares.NewJwtToken(accountInfo.AccountID, middlewares.Expires*time.Hour)
		if err != nil {
			res.WriteHeader(http.StatusOK)
			bytes, e := json.Marshal(util.MakeResponseWithMsg(err.Error()))
			if e != nil {
				zap.S().Error(fmt.Errorf("json Marshal is failed: %w", e))
			}
			_, err = res.Write(bytes)
			if err != nil {
				zap.S().Error(fmt.Errorf("request write NewJwtToken err is failed: %w", err))
			}
		}

		redirect := fmt.Sprintf("%s/home", configs.GetConfig().RedirectURL)

		zap.S().Info(t)
		cookie := mycookie.NewStarlandAICookie(mycookie.StarlandAIToken, t)

		http.SetCookie(res, &cookie)

		zap.S().Info("redirect  ", redirect)
		http.Redirect(res, req, redirect, http.StatusFound)
	}
}

func authLogin() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		req = util.FiberGothAdapter(req)

		l, err := gothic.GetAuthURL(res, req)
		if err != nil {
			res.WriteHeader(http.StatusOK)
			bytes, e := json.Marshal(util.MakeResponseWithMsg(err.Error()))
			if e != nil {
				zap.S().Error(fmt.Errorf("json Marshal is failed: %w", e))
			}
			_, err = res.Write(bytes)
			if err != nil {
				zap.S().Error(fmt.Errorf("request write GetAuthURL err is failed: %w", err))
			}
			return
		}
		zap.S().Infof("gothic: url :%s", l)
		zap.S().Info("authLogin cookie", req.Cookies())
		redirect := fmt.Sprintf("%s/home", configs.GetConfig().RedirectURL)
		zap.S().Infof("redirect: %s", redirect)
		cookie := mycookie.NewStarlandAICookie(mycookie.StarlandAIRedirect, redirect)

		http.SetCookie(res, &cookie)

		escape := url.QueryEscape(l)
		response := util.MakeResponse(fiber.Map{"url": escape})
		marshal, err := json.Marshal(response)
		if err != nil {
			res.WriteHeader(http.StatusOK)
			bytes, e := json.Marshal(util.MakeResponseWithMsg(err.Error()))
			if e != nil {
				zap.S().Error(fmt.Errorf("json Marshal is failed: %w", e))
			}
			_, err = res.Write(bytes)
			if err != nil {
				zap.S().Error(fmt.Errorf("request write json marshal err is failed: %w", err))
			}
			return
		}
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(marshal)
		if err != nil {
			zap.S().Error(fmt.Errorf("request write marshal err is failed: %w", err))
		}
	}
}

func activity(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Account      string `params:"id"`
				ActivityCode int    `json:"activity_code"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		aid := ctx.Locals(middlewares.LocalsAccount)
		accountID := aid.(string)

		if req.Account != accountID {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg("not auth"))
		}

		activityInfo := &account.ActivityRequest{
			AccountID:    req.Account,
			ActivityCode: req.ActivityCode,
		}

		if err := service.Activity(ctx.Context(), activityInfo); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}

		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func sendSigninMail(server AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req struct {
			Mail string `form:"mail"`
		}
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		mailOnce.Do(func() {
			storage = memory.New()
			re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]" +
				"{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		})

		if !re.MatchString(req.Mail) {
			return ctx.Status(http.StatusBadRequest).
				JSON(util.MakeErrResponse(errors.New(" Email format is incorrect")))
		}
		tokenExpires := middlewares.EmailExpires * time.Minute

		if v, _ := storage.Get(req.Mail); v != nil {
			return ctx.Status(http.StatusOK).
				JSON(util.MakeErrResponse(errors.New(" Too many retry attempts, please wait")))
		}

		err := server.LoginSendMail(ctx.Context(), req.Mail, tokenExpires)
		if err != nil {
			zap.S().Errorf("sendSigninMail: LoginSendMail err: %v", err)
			return ctx.Status(http.StatusInternalServerError).
				JSON(util.MakeErrResponse(err))
		}

		if err = storage.Set(req.Mail, []byte("ok"), middlewares.EmailExpires*time.Second); err != nil {
			return ctx.Status(http.StatusInternalServerError).
				JSON(util.MakeErrResponse(err))
		}
		return ctx.Status(http.StatusOK).
			JSON(util.MakeResponse(fmt.Sprintf("email send successfully to %s", req.Mail)).SetCode("0"))
	}
}

func claimPoints(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ClaimPoints int  `json:"claim_points"`
				IsOK        bool `json:"is_ok"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		aid := ctx.Locals(middlewares.LocalsAccount)
		accountID := aid.(string)
		claimPointsRequest := &account.ClaimPointsRequest{
			AccountID: accountID,
			Points:    req.ClaimPoints,
			IsOk:      req.IsOK,
		}

		res, err := service.ClaimPoints(ctx.Context(), claimPointsRequest)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}

		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func savePointsAddr(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Account    string `params:"id"`
				SolanaAddr string `json:"solana_addr"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		aid := ctx.Locals(middlewares.LocalsAccount)
		accountID := aid.(string)
		zap.S().Info(accountID)
		if req.Account != accountID {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg("not auth"))
		}

		err := service.SavePointsAddr(ctx.Context(), accountID, req.SolanaAddr)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}
