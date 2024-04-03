package v1

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"starland-backend/configs"
	"starland-backend/internal/pkg/middlewares"
	"starland-backend/internal/pkg/util"
	"starland-backend/internal/service/character"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	ImageCountLimit = 4
)

var (
	q chan struct{}
)

var (
	CreateCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "createCount",
			Help: "createCount custom metric",
		},
		[]string{"sum"},
	)
	ChatCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chatCount",
			Help: "chatCount custom metric",
		},
		[]string{"sum"})
)

type CharacterHTTPServer interface {
	CreateCharacter(context.Context, *character.CreateCharacterRequest) (*character.CreateCharacterResponse, error)
	QueryImageModel(context.Context) ([]*character.ImageModelResponse, error)
	CharacterLike(context.Context, string, string, bool) error
	Chat(context.Context, *character.ChatRequest) (*character.ChatResponse, error)
	QueryCharacters(context.Context, *character.QueryCharacterRequest) ([]*character.QueryCharacterResponse, int64, error)
	QueryCharacterInfo(context.Context, string) (*character.QueryCharacterResponse, error)
	CharacterMint(context.Context, *character.CharacterMintRequest) error
	QueryCharactersHistory(ctx context.Context, req *character.QueryCharactersHistoryRequest) ([]*character.QueryCharactersHistoryResponse,
		int64, error)
	CreateCharacterV2(context.Context, *character.CreateCharacterRequest) (*character.CreateCharacterResponse, error)
	ChatV2(context.Context, *character.ChatRequestV2) error
	MessageToVoice(context.Context, string, string) (string, error)
	QueryVoice(ctx context.Context) ([]*character.CharacterVoiceResponse, error)
	UpdateCharacter(context.Context, *character.UpdateCharacterRequest) error
	DeleteCharacter(context.Context, *character.DeleteCharacterRequest) error
}

func InitCharacterRouter(app fiber.Router, service CharacterHTTPServer, conf *configs.Config) {
	router := app.Group("/v1")
	q = make(chan struct{}, conf.Chat.ChatLimit)

	router.Get("/character/image_model", queryImageModels(service))
	router.Get("/character/voice", queryVoice(service))

	router.Get("/character", middlewares.JwtParse(), queryCharacter(service))
	router.Get("/character/my", middlewares.JwtParse(), queryMyCharacter(service))
	router.Post("/character", middlewares.JwtParse(), func(ctx *fiber.Ctx) error {
		CreateCountMetric.WithLabelValues("count").Inc()
		return ctx.Next()
	}, createV2(service))

	router.Put("/character/:id", middlewares.JwtParse(), updateCharacter(service))

	router.Get("/character/history", middlewares.JwtParse(), history(service))
	router.Post("/character/:id/like", middlewares.JwtParse(), like(service))
	router.Post("/character/:id/mint", middlewares.JwtParse(), mint(service))
	router.Delete("/character/:id", middlewares.JwtParse(), deleteCharacter(service))
	/* router.Post("/character/:id/chat", middlewares.JwtParse(), func(ctx *fiber.Ctx) error {
		ChatCountMetric.WithLabelValues("count").Inc()
		return ctx.Next()
	}, chat(service)) */
	router.Post("/character/:id/chat", middlewares.JwtParse(), func(ctx *fiber.Ctx) error {
		ChatCountMetric.WithLabelValues("count").Inc()
		return ctx.Next()
	}, ChatV2(service))

	router.Get("/character/:id", middlewares.JwtParse(), info(service))
	router.Get("/stream", adaptor.HTTPHandlerFunc(stream()))
}

func create(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Message   string `json:"message"`             
				SessionID string `json:"session_id,omitempty"` 
				State     int    `json:"state"`               
			}
		)

		if len(q) >= configs.GetConfig().Chat.ChatLimit {
			return ctx.Status(http.StatusOK).
				JSON(util.MakeResponseWithMsg("You're in queue at position 100. Please retry in 10 minutes.").SetCode("1"))
		} else {
			q <- struct{}{}
			defer func() {
				<-q
			}()
		}

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		res, err := service.CreateCharacter(ctx.Context(), &character.CreateCharacterRequest{
			Message:   req.Message,
			SessionID: req.SessionID,
			State:     req.State,
		})

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}

		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func createV2(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Message   string `json:"message"`    
				SessionID string `json:"session_id"` 
				State     int    `json:"state"`      
				Is3D      bool   `json:"is_3d"`
			}
		)
		accountID := ctx.Locals(middlewares.LocalsAccount).(string)

		if len(q) >= configs.GetConfig().Chat.ChatLimit {
			return ctx.Status(http.StatusOK).
				JSON(util.MakeResponseWithMsg("You're in queue at position 100. Please retry in 10 minutes.").SetCode("1"))
		} else {
			q <- struct{}{}
			defer func() {
				<-q
			}()
		}

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		resCh := make(chan character.ChatCompletionStreamResponseChunk)
		errCh := make(chan error)
		go func() {
			defer close(errCh)
			_, err := service.CreateCharacterV2(ctx.Context(), &character.CreateCharacterRequest{
				Message:   req.Message,
				SessionID: req.SessionID,
				State:     req.State,
				ResCh:     resCh,
				Is3D:      req.Is3D,
				AccountID: accountID,
			})
			if err != nil {
				zap.S().Errorf("createV2: create err: %v", err)
				errCh <- err
			}
		}()

		ctx.Context().SetContentType("text/event-stream")
		ctx.Set("Content-Type", "text/event-stream")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Set("Connection", "keep-alive")
		ctx.Set("Transfer-Encoding", "chunked")
		ctx.Set("X-Accel-Buffering", "no")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			zap.S().Info("WRITER")
			var (
				i    int
				data character.ChatCompletionStreamResponseChunk
				h    string
			)
			for {
				i++
				if resMessage, ok := <-resCh; ok {
					if req.SessionID == "" {
						resMessage.Message = resMessage.SessionID
					}
					if resMessage.Message == "" || h == resMessage.Message {
						continue
					}

					m := strings.ReplaceAll(resMessage.Message, "\n", "\\n")
					h = resMessage.Message
					msg := fmt.Sprintf("data: %s\n\n", m)
					fmt.Fprint(w, msg)
					err := w.Flush()
					if err != nil {
						fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
						return
					}
					data = resMessage
					time.Sleep(time.Millisecond * 50)
				} else {
					zap.S().Infof("%+v", data)

					jsonData, err := json.Marshal(util.MakeResponse(data))
					if err != nil {
						zap.S().Errorf("resMessage: %v", err)
					}
					fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
					err = w.Flush()
					if err != nil {
						fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
						return
					}
					return
				}
			}
		}))
		return nil
	}
}

func queryImageModels(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		res, err := service.QueryImageModel(ctx.Context())
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func like(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID   string `params:"id"`
				Flag bool   `json:"flag"`
			}
		)
		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		zap.S().Info("accountID:", accountID)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		zap.S().Infof("like: req: %+v", req)

		err := service.CharacterLike(ctx.Context(), req.ID, accountID, req.Flag)

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func queryCharacter(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Search string `query:"Search"`
				Page   int    `query:"page"`
				Limit  int    `query:"limit"`
			}
			res struct {
				Data  []*character.QueryCharacterResponse `json:"data"`
				Count int64                               `json:"count"`
			}
		)
		if err := ctx.QueryParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		characters, count, err := service.QueryCharacters(ctx.Context(), &character.QueryCharacterRequest{
			Search: req.Search,
			Page:   req.Page,
			Limit:  req.Limit,
		})
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.Data = characters
		res.Count = count
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func queryMyCharacter(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Search string `query:"Search"`
				Page   int    `query:"page"`
				Limit  int    `query:"limit"`
			}
			res struct {
				Data  []*character.QueryCharacterResponse `json:"data"`
				Count int64                               `json:"count"`
			}
			accountID string
		)

		if err := ctx.QueryParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		account := ctx.Locals(middlewares.LocalsAccount)
		if account == nil {
			accountID = ""
		} else {
			accountID = account.(string)
		}

		characters, count, err := service.QueryCharacters(ctx.Context(), &character.QueryCharacterRequest{
			Page:    req.Page,
			Limit:   req.Limit,
			Account: accountID,
			Search:  req.Search,
		})
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.Data = characters
		res.Count = count
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func chat(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID      string `params:"id"`
				Message string `json:"Message"`
			}

			res struct {
				ChatMessage []string `json:"chat_message"`
			}
		)

		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		zap.S().Infof("chat: req: %+v", req)

		r, err := service.Chat(ctx.Context(), &character.ChatRequest{
			Message:     req.Message,
			CharacterID: req.ID,
			AccountID:   accountID,
		})
		if r == nil {
			res.ChatMessage = []string{}
		} else {
			res.ChatMessage = []string{r.Message}
		}
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func ChatV2(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			reqData struct {
				ID      string `params:"id"`
				Message string `json:"Message"`
			}

			resData struct {
				ChatMessage string `json:"chat_message"`
				Voice       string `json:"voice"`
			}
			message string
		)

		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		if err := ctx.ParamsParser(&reqData); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		if err := ctx.BodyParser(&reqData); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		resCh := make(chan interface{})
		err := service.ChatV2(ctx.Context(), &character.ChatRequestV2{
			Message:     reqData.Message,
			CharacterID: reqData.ID,
			AccountID:   accountID,
			ResCh:       resCh,
		})
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}
		ctx.Context().SetContentType("text/event-stream")
		ctx.Set("Content-Type", "text/event-stream")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Set("Connection", "keep-alive")
		ctx.Set("Transfer-Encoding", "chunked")
		ctx.Set("X-Accel-Buffering", "no")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			zap.S().Info("WRITER")
			var (
				i       int
				content string
			)
			for {
				i++
				if resMessage, ok := <-resCh; ok {
					message = resMessage.(string)
					msg := strings.ReplaceAll(message, "\n", "\\n")

					fmt.Fprintf(w, "data: %s\n\n", msg)
					err = w.Flush()
					if err != nil {
						fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
						return
					}
					content = fmt.Sprintf("%s%s", content, message)
				} else {
					zap.S().Infof("chat res: %s", content)
					resData.ChatMessage = content
					voice, err := service.MessageToVoice(context.Background(), reqData.ID, content)
					if err != nil {
						zap.S().Errorf("MessageToVoice: err: %v", err)
						break
					}
					resData.Voice = voice
					zap.S().Errorf("resData: %+v", resData)
					res, err := json.Marshal(util.MakeResponse(resData))
					if err != nil {
						zap.S().Errorf("MessageToVoice: err: %v", err)
						break
					}
					fmt.Fprintf(w, "data: %s\n\n", string(res))
					break
				}
			}
		}))
		return nil
	}
}

func history(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Page  int `query:"page"`
				Limit int `query:"limit"`
			}
			res struct {
				Data  []*character.QueryCharactersHistoryResponse `json:"data"`
				Count int64                                       `json:"count"`
			}
		)
		if err := ctx.QueryParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		characters, count, err := service.QueryCharactersHistory(ctx.Context(), &character.QueryCharactersHistoryRequest{
			Page:    req.Page,
			Limit:   req.Limit,
			Account: accountID,
		})

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.Data = characters
		res.Count = count
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func info(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID string `params:"id"`
			}
		)
		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		characters, err := service.QueryCharacterInfo(ctx.Context(), req.ID)

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(characters))
	}
}

func mint(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID   string `params:"id"`
				Mint string `json:"mint"`
			}
		)
		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		err := service.CharacterMint(ctx.Context(), &character.CharacterMintRequest{
			ID:   req.ID,
			Mint: req.Mint,
		})

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func updateCharacter(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Description string   `form:"description"`
				ID          string   `params:"id"`
				Images      []string `form:"images"`
				Image       string   `form:"image"`
				Name        string   `form:"name"`
				Voice       string   `form:"voice"`
			}
		)

		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		num, err := strconv.Atoi(req.Image)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		files, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeErrResponse(err))
		}

		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		zap.S().Infof("updateCharacter: file size: %d", len(files.File["files"]))
		dir := fmt.Sprintf("%s%s", configs.GetConfig().File.ImagePath, accountID)
		if err = util.Mkdir(dir); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeErrResponse(err))
		}

		imageMap := map[string]struct{}{"1": {}, "2": {}, "3": {}, "4": {}}

		for i := range req.Images {
			fileName := util.URL2FileName(req.Images[i])
			fileKey := strings.Split(fileName, ".")[0]
			delete(imageMap, fileKey)
		}
		if len(req.Images)+len(files.File["files"]) > ImageCountLimit {
			return ctx.Status(http.StatusBadRequest).JSON(util.MakeErrResponse(errors.New(" You've reached the limit for upload image")))
		}
		urlLen := len(req.Images)
		for _, file := range files.File["files"] {
			for key := range imageMap {
				fileExt := filepath.Ext(file.Filename)
				filePath := fmt.Sprintf("%s/%s%s", dir, key, fileExt)
				if err = ctx.SaveFile(file, filePath); err != nil {
					return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
				}
				req.Images = append(req.Images, fmt.Sprintf("%s?t=%d", configs.GetConfig().File.ImagesEndpoint+accountID+"/"+key+fileExt, time.Now().Nanosecond()))
				zap.S().Infof("images :%s", req.Images[len(req.Images)-1])
				delete(imageMap, key)
				break
			}
		}

		if num > len(req.Images) && num <= 9+len(files.File["files"]) {
			num = num%10 + urlLen
			req.Image = req.Images[num]
		} else if num < len(req.Images) && num >= 0 {
			req.Image = req.Images[num]
		}

		err = service.UpdateCharacter(ctx.Context(), &character.UpdateCharacterRequest{
			ID:          req.ID,
			AccountID:   accountID,
			Description: req.Description,
			Name:        req.Name,
			Images:      req.Images,
			Image:       req.Image,
			Voice:       req.Voice,
		})
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeErrResponse(err))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func queryVoice(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		res, err := service.QueryVoice(ctx.Context())
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func deleteCharacter(service CharacterHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID string `params:"id"`
			}
		)
		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		accountID := ctx.Locals(middlewares.LocalsAccount).(string)
		err := service.DeleteCharacter(ctx.Context(), &character.DeleteCharacterRequest{
			ID:        req.ID,
			AccountID: accountID,
		})

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func stream() func(w http.ResponseWriter, r *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain")

		for i := 0; i < 10; i++ {
			data := fmt.Sprintf("Data %d\n", i)
			_, err := res.Write([]byte(data))
			if err != nil {
				break
			}
			res.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
	}
}
