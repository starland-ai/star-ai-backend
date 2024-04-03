package http

import (
	"fmt"
	v1 "starland-backend/api/http/v1"
	"starland-backend/configs"
	"starland-backend/internal/service"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewHTTPServer(config *configs.Config, us *service.Service) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		ReadTimeout:  config.HTTP.ReadTimeout * time.Second,
		WriteTimeout: config.HTTP.WriteTimeout * time.Second,
	})

	app.Use(recover.New(), pprof.New(), cors.New(), requestid.New())

	prometheus.MustRegister(v1.ChatCountMetric, v1.CreateCountMetric)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	prometheus := fiberprometheus.NewWithRegistry(prometheus.DefaultRegisterer, "starland-backend", "http", "", nil)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(logger.New(logger.Config{
		Format: fmt.Sprintf("${time} | ${ip} | ${status} | ${locals:%s} | ${latency} | ${method} | ${path} | RequestBody:${reqBody}"+
			"ResponseBody:${resBody} | Params:${} \n",
			requestid.ConfigDefault.ContextKey),
		Next: func(c *fiber.Ctx) bool {
			path := string(c.Request().URI().Path())
			if strings.Contains(path, "/v1/file") || strings.Contains(path, "/chat") || strings.Contains(path, "/v1/character") {
				return true
			} else {
				return false
			}
		},
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Shanghai",
	}))
	r := app.Group("")
	v1.InitAccountRouter(r, us.Account, config)
	v1.InitCharacterRouter(r, us.Character, config)
	v1.InitFileRouter(r, config)
	zap.S().Infof("addr:%s", config.HTTP.Addr)
	return app, nil
}
