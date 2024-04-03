package v1

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"starland-backend/configs"
	"starland-backend/internal/pkg/middlewares"
	"starland-backend/internal/pkg/util"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileHTTPServer struct {
}

func InitFileRouter(app fiber.Router, conf *configs.Config) {
	router := app.Group("/v1")
	router.Get("/file/*", file())

	router.Post("/upload", middlewares.JwtParse(), upload())
}

func upload() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeErrResponse(err))
		}
		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		os.Mkdir(accountID, fs.ModeDir)
		if err = ctx.SaveFile(file, fmt.Sprintf("./%s/%s_%s", accountID, uuid.NewString(), file.Filename)); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponseWithMsg("ok"))
	}
}

func file() func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		zap.S().Info(c.Request().URI().String())
		path := strings.Split(c.Request().URI().String(), "/file")[1]
		if strings.Contains(path, "?") {
			path = strings.Split(path, "?")[0]
		}
		// Serve the file using SendFile function
		zap.S().Info(path)
		if strings.Contains(path, ".png") {
			c.Response().Header.Add("Content-Type", "image/png")
		}
		fileName := filepath.Base(path)
		c.Response().Header.Add("Content-Disposition", "attachment; filename="+fileName)
		err := filesystem.SendFile(c, http.Dir("."), path)
		if err != nil {
			// Handle the error, e.g., return a 404 Not Found response
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}
		return nil
	}
}
