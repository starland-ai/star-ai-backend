package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	LocalsAccount = "account"
	Expires       = 3000
	EmailExpires  = 30
	TokenSecret   = "starland-ai"
	SigningMethod = "HS256"
	/*	TokenClaimsName     = "name"
		TokenClaimsID       = "UserID"
		TokenClaimsExpires  = "exp"*/
	TokenClaimsEmail    = "Email"
	TokenClaimsProvider = "Provider"
	StatusUnauthorized  = "40001"
)

type MyClaims struct {
	UID string `json:"UID"`
	jwt.RegisteredClaims
}

func JwtParse() fiber.Handler {
	return jwtware.New(jwtware.Config{
		TokenLookup: fmt.Sprintf("header:%s,query:token", fiber.HeaderAuthorization),
		SuccessHandler: func(ctx *fiber.Ctx) error {
			zap.S().Info("HeaderAuthorization:", ctx.Get(fiber.HeaderAuthorization))
			token := ctx.Locals(LocalsAccount).(*jwt.Token)
			ctx.Locals(LocalsAccount, token.Claims.(jwt.MapClaims)["UID"])
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if strings.HasSuffix(ctx.Path(), "/character") && string(ctx.Request().Header.Method()) == "GET" {
				return ctx.Next()
			}
			zap.S().Info(ctx.Get(fiber.HeaderAuthorization))
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code": StatusUnauthorized,
				"msg":  "token verification failed",
			})
		},
		SigningKey:    []byte(TokenSecret),
		SigningMethod: SigningMethod,
		Claims:        jwt.MapClaims{},
		ContextKey:    LocalsAccount,
	})
}

func JwtParseRedirect(url string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		TokenLookup: fmt.Sprintf("header:%s,query:token", fiber.HeaderAuthorization),
		SuccessHandler: func(ctx *fiber.Ctx) error {
			token := ctx.Locals(LocalsAccount).(*jwt.Token)
			ctx.Locals(LocalsAccount, token.Claims.(jwt.MapClaims)["UID"])
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			zap.S().Error(err)
			return ctx.Redirect(url, http.StatusFound)
		},
		SigningKey:    []byte(TokenSecret),
		SigningMethod: SigningMethod,
		Claims:        jwt.MapClaims{},
		ContextKey:    LocalsAccount,
	})
}

func NewJwtToken(uid string, expiresAt time.Duration) (string, error) {
	claims := &MyClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(TokenSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}
