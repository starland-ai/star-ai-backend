package biz

import (
	"context"
	"fmt"
	"io"
	"starland-backend/configs"
	"starland-backend/internal/pkg/bizerr"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"

	"starland-backend/internal/proto/character_agent"
	"starland-backend/internal/proto/chat_agent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	ChatChunk        = 1
	ImageChunk       = 2
	NeedConfirmChunk = 3
	SettingChunk     = 4
)

type CharacterRepo interface {
	SaveCharacter(context.Context, *CharacterRequest) (string, error)
	QueryCharacterByID(context.Context, string) (*CharacterResponse, error)
	QueryCharactersByAccountID(context.Context, string, string, int, int) ([]*CharacterResponse, int64, error)
	QueryCharactersByNameOrPrompt(context.Context, string, int, int) ([]*CharacterResponse, int64, error)
	CharacterMintSave(context.Context, string, string) error
	UpdateCharacter(context.Context, *UpdateCharacterRequest) error
	DeleteCharacterByID(context.Context, string) error
}

type CharacterAccountLikesRepo interface {
	QueryCharacterAccountLike(context.Context, string, string) (bool, error)
	SaveCharacterAccountLike(context.Context, string, string, bool) error
	QueryCharacterLikeCount(context.Context, string) (int64, error)
}

type CharacterUsecase struct {
	conf               *configs.Config
	characterRepo      CharacterRepo
	characterLikesRepo CharacterAccountLikesRepo
	characterPoll      *grpcpool.Pool
	chatPoll           *grpcpool.Pool
}

func NewCharacterUsecase(conf *configs.Config, characterRepo CharacterRepo,
	characterLikesRepo CharacterAccountLikesRepo) *CharacterUsecase {

	character, err := grpcpool.New(func() (*grpc.ClientConn, error) {
		return grpc.Dial(conf.ChatCompletions.Endpoint, grpc.WithInsecure())
	}, 10, 30, time.Second*30)

	chat, err := grpcpool.New(func() (*grpc.ClientConn, error) {
		return grpc.Dial(conf.Chat.Endpoint, grpc.WithInsecure())
	}, 10, 30, time.Second*30)

	if err != nil {
		zap.S().Fatalf("NewChatCompletionsClient: failed to create client pool: %v", err)
	}
	return &CharacterUsecase{conf: conf,
		characterRepo:      characterRepo,
		characterPoll:      character,
		characterLikesRepo: characterLikesRepo,
		chatPoll:           chat}
}

type CharacterRequest struct {
	ID           string
	AccountID    string
	AccountName  string
	AvatarURL    string
	Name         string
	Gender       int
	Prompt       string
	ImageURL     string
	LikeCount    int
	ChatCount    int
	State        int
	Tags         map[string]string
	Voice        string
	Introduction string
	Is3D         bool
}

type CharacterResponse struct {
	ID           string
	AccountID    string
	AccountName  string
	AvatarURL    string
	Name         string
	Gender       int
	Prompt       string
	ImageURL     string
	ImageURLs    []string
	IsMint       bool
	UpdateTime   time.Time
	LikeCount    int
	ChatCount    int
	Tag          []Tag
	Mint         string
	Voice        string
	Is3D         bool
	Introduction string
}
type Tag struct {
	Key   string
	Value string
}
type ChatMessage struct {
	Role    string
	Content string
}

type ChatCompletionsRequest struct {
	SessionID string
	Message   *ChatMessage
	ResCh     chan ChatCompletionStreamResponseChunk
}

type ChatRequest struct {
	ConversationID string
	Message        string
	CharacterId    string
	ResCh          chan interface{}
	CharacterName  string
}

type ChatResponse struct {
	ResponseText string
}
type ChatCompletionResponse struct {
	Message     []*ChatMessage
	ImageMetas  []string
	NeedConfirm bool
	ConfirmType string
}

type ChatCompletionStreamResponse struct {
	Code   int
	ErrMsg string
	Chunk  *ChatCompletionStreamResponseChunk
}

type ChatCompletionStreamResponseChunk struct {
	Is3D              bool         `json:is_3d,omitempty`
	ChunkType         uint32       `json:"chunk_type,omitempty"`
	ChunkSessionIndex uint32       `json:"chunk_session_index,omitempty"`
	ChatChunk         *ChatMessage `json:"chat_chunk,omitempty"`
	ImageChunk        []string     `json:"image_chunk,omitempty"`
	NeedConfirmChunk  bool         `json:"need_confirm_chunk,omitempty"`
	SettingChunk      string       `json:"setting_chunk,omitempty"`
}

type ConfirmCharacterSettingRequest struct {
	SessionID string
}

type ConfirmCharacterSettingResponse struct {
	Message          []*ChatMessage
	NeedConfirm      bool
	ConfirmType      string
	CharacterSetting *CharacterSetting
}

type CharacterSetting struct {
	Name         string
	Gender       string
	Description  string
	Tags         map[string]string
	Introduction string
}

type UpdateCharacterRequest struct {
	ID          string
	Images      []string
	Image       string
	Name        string
	Description string
	Voice       string
}

func (uc *CharacterUsecase) ChatCompletions(ctx context.Context, req *ChatCompletionsRequest) (*ChatCompletionResponse, error) {
	conn, err := uc.characterPoll.Get(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletions: get grpc conn err: %w", err))
	}
	defer conn.Close()

	cli := character_agent.NewAgentClient(conn)
	grpcReq := &character_agent.ChatCompletionRequest{
		SessionId: req.SessionID,
		Message: &character_agent.ChatMessage{
			Role:    req.Message.Role,
			Content: req.Message.Content,
		},
	}
	resp, err := cli.ChatCompletions(ctx, grpcReq)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			err = fmt.Errorf("ChatCompletions: invoke chat agent timeout ")
		}
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletions: grpc exec failed err: %w ", err))
	}

	if resp.Code != uint32(codes.OK) {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletions: grpc res failed err: %s ", resp.ErrMsg))
	}

	res := &ChatCompletionResponse{
		Message:     makeChatMessages(resp.Messages),
		ImageMetas:  makeImageMetas(resp.ImageMetas),
		ConfirmType: resp.ConfirmType,
		NeedConfirm: resp.NeedConfirm,
	}
	return res, nil
}

/* func (uc *CharacterUsecase) MockChatCompletions(ctx context.Context,
	req *ChatCompletionsRequest) (*ChatCompletionResponse, error) {

	res := &ChatCompletionResponse{
		Message: []*ChatMessage{{
			Role:    "AI",
			Content: "Test",
		}},
		NeedConfirm: false,
	}
	return res, nil
}

func (uc *CharacterUsecase) MockConfirmCharacterSetting(ctx context.Context,
	req *ConfirmCharacterSettingRequest) (*ConfirmCharacterSettingResponse, error) {
	res := &ConfirmCharacterSettingResponse{
		Message: []*ChatMessage{{
			Role:    "AI",
			Content: "Test",
		}},
		NeedConfirm: false,
		CharacterSetting: &CharacterSetting{
			Name:        "test",
			Gender:      "man",
			Description: "xxxx",
			Tags:        map[string]string{"hhh": "HHHH"},
		},
	}
	return res, nil
} */

func (uc *CharacterUsecase) ConfirmCharacterSetting(ctx context.Context,
	req *ConfirmCharacterSettingRequest) (*ConfirmCharacterSettingResponse, error) {
	conn, err := uc.characterPoll.Get(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("ConfirmCharacterSetting: get grpc conn err: %w", err))
	}
	defer conn.Close()

	cli := character_agent.NewAgentClient(conn)
	grpcReq := &character_agent.ConfirmCharacterSettingRequest{
		SessionId: req.SessionID,
	}
	resp, err := cli.ConfirmCharacterSetting(ctx, grpcReq)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			err = fmt.Errorf("ConfirmCharacterSetting: invoke chat agent timeout ")
		}
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("ConfirmCharacterSetting: grpc exec failed err: %w ", err))
	}
	if resp.Code != uint32(codes.OK) {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("ConfirmCharacterSetting: grpc res failed err: %s ", resp.ErrMsg))
	}

	res := &ConfirmCharacterSettingResponse{
		Message:     makeChatMessages(resp.Messages),
		ConfirmType: resp.ConfirmType,
		NeedConfirm: resp.NeedConfirm,
	}
	if resp.CharacterSetting != nil {
		res.CharacterSetting = &CharacterSetting{
			Name:         resp.CharacterSetting.Name,
			Gender:       resp.CharacterSetting.Gender,
			Description:  resp.CharacterSetting.Description,
			Tags:         resp.CharacterSetting.Tags,
			Introduction: resp.CharacterSetting.Introduction,
		}
	} else {
		zap.S().Infof("CharacterSetting is null")
	}
	return res, nil
}

func (uc *CharacterUsecase) ChatCompletionsStream(ctx context.Context, req *ChatCompletionsRequest) error {
	defer func() {
		close(req.ResCh)
		zap.S().Info("ChatCompletionsStream end")
	}()
	zap.S().Infof("ChatCompletionsStream: %+v", *req)
	conn, err := grpc.Dial(uc.conf.ChatCompletions.Endpoint, grpc.WithInsecure())
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletionsStream: get grpc conn err: %w", err))
	}
	defer conn.Close()

	cli := character_agent.NewAgentClient(conn)
	grpcReq := &character_agent.ChatCompletionRequest{
		SessionId: req.SessionID,
		Message: &character_agent.ChatMessage{
			Role:    req.Message.Role,
			Content: req.Message.Content,
		},
	}
	stream, err := cli.ChatCompletionsStream(ctx, grpcReq)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			err = fmt.Errorf("ChatCompletionsStream: invoke chat agent timeout ")
		}
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletionsStream: get stream cli err: %w ", err))
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				zap.S().Info("ChatCompletionsStream EOF")
				break
			}
			return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletionsStream: recv failed err: %w ", err))
		}

		if resp.Code != uint32(codes.OK) {
			return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatCompletionsStream: resp err code:%d err:%s ", resp.Code, resp.ErrMsg))
		}

		var data ChatCompletionStreamResponseChunk
		switch resp.Chunk.ChunkType {
		case ChatChunk:
			data = ChatCompletionStreamResponseChunk{
				ChunkType:         resp.Chunk.ChunkType,
				ChunkSessionIndex: resp.Chunk.ChunkSessionIndex,
				ChatChunk: &ChatMessage{
					Role:    resp.Chunk.ChatChunk.Role,
					Content: resp.Chunk.ChatChunk.Content,
				},
			}
		case ImageChunk:
			var is3D bool
			images := make([]string, len(resp.Chunk.ImageChunk))
			for i := range resp.Chunk.ImageChunk {
				images[i] = resp.Chunk.ImageChunk[i].Id
				is3D = resp.Chunk.ImageChunk[i].Enable3D
			}

			zap.S().Infof("images: %+v %+v", images, resp.Chunk.ImageChunk[0])
			data = ChatCompletionStreamResponseChunk{
				Is3D:              is3D,
				ChunkType:         resp.Chunk.ChunkType,
				ChunkSessionIndex: resp.Chunk.ChunkSessionIndex,
				ImageChunk:        images,
			}
		case NeedConfirmChunk:
			data = ChatCompletionStreamResponseChunk{
				ChunkType:         resp.Chunk.ChunkType,
				ChunkSessionIndex: resp.Chunk.ChunkSessionIndex,
				NeedConfirmChunk:  resp.Chunk.NeedConfirmChunk,
			}
		case SettingChunk:
			data = ChatCompletionStreamResponseChunk{
				ChunkType:         resp.Chunk.ChunkType,
				ChunkSessionIndex: resp.Chunk.ChunkSessionIndex,
				SettingChunk:      resp.Chunk.SettingChunk,
			}
		default:
			return bizerr.ErrChunkNotExist.Wrap(fmt.Errorf("ChatCompletionsStream: chunk not exists"))
		}
		req.ResCh <- data
	}
	return nil
}

func makeChatMessages(req []*character_agent.ChatMessage) []*ChatMessage {
	res := make([]*ChatMessage, len(req))
	for i := range req {
		res[i] = &ChatMessage{
			Role:    req[i].Role,
			Content: req[i].Content,
		}
	}
	return res
}

func makeImageMetas(req []*character_agent.ImageMeta) []string {
	res := make([]string, len(req))
	for i := range req {
		res[i] = req[i].Name
	}
	return res
}

func (uc *CharacterUsecase) SaveMyCharacter(ctx context.Context, req *CharacterRequest) (string, error) {
	zap.S().Infof("SaveMyCharacter: req: %+v", req)
	res, err := uc.characterRepo.SaveCharacter(ctx, req)
	if err != nil {
		return "", bizerr.ErrInternalError.Wrap(fmt.Errorf("SaveMyCharacter: save character to db err: %w", err))
	}
	return res, nil
}

func (uc *CharacterUsecase) QueryCharacterByID(ctx context.Context, id string) (*CharacterResponse, error) {
	res, err := uc.characterRepo.QueryCharacterByID(ctx, id)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharacterByID: query character to db err: %w", err))
	}
	if res == nil {
		return nil, bizerr.ErrCharacterNotExist
	}
	return res, nil
}

func (uc *CharacterUsecase) QueryCharactersByAccount(ctx context.Context, id, query string, page, limit int) ([]*CharacterResponse, int64, error) {
	res, count, err := uc.characterRepo.QueryCharactersByAccountID(ctx, id, query, page, limit)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharacterByID: query character to db err: %w", err))
	}

	return res, count, nil
}

func (uc *CharacterUsecase) QueryCharactersByNameOrPrompt(ctx context.Context, query string,
	page, limit int) ([]*CharacterResponse, int64, error) {
	res, count, err := uc.characterRepo.QueryCharactersByNameOrPrompt(ctx, query, page, limit)
	if err != nil {
		return nil, count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharactersByNameOrPrompt: query character to db err: %w", err))
	}
	return res, count, nil
}

func (uc *CharacterUsecase) QueryCharacterLikeCountByAccount(ctx context.Context, id string) (int64, error) {
	count, err := uc.characterLikesRepo.QueryCharacterLikeCount(ctx, id)
	if err != nil {
		return count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharacterLikeCountByAccount: query in db err: %w ", err))
	}
	return count, nil
}

func (uc *CharacterUsecase) QueryCharacterLikeByAccount(ctx context.Context, characterID, accountID string) (bool, error) {
	count, err := uc.characterLikesRepo.QueryCharacterAccountLike(ctx, characterID, accountID)
	if err != nil {
		return count, bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharacterLikeCountByAccount: query in db err: %w ", err))
	}
	return count, nil
}

func (uc *CharacterUsecase) UpdateCharacterLikeAccountLike(ctx context.Context, characterID, accountID string, flag bool) error {
	err := uc.characterLikesRepo.SaveCharacterAccountLike(ctx, characterID, accountID, flag)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("QueryCharacterLikeCountByAccount: query in db err: %w ", err))
	}
	return nil
}

func (uc *CharacterUsecase) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	conn, err := uc.chatPoll.Get(ctx)
	if err != nil {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("Chat: get grpc conn err: %w", err))
	}
	defer conn.Close()

	cli := chat_agent.NewAgentClient(conn)
	grpcReq := &chat_agent.ChatRequest{
		CharacterId:    req.CharacterId,
		ConversationId: req.ConversationID,
		Message:        []byte(req.Message),
	}
	resp, err := cli.Chat(ctx, grpcReq)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			err = fmt.Errorf("Chat: invoke chat agent timeout ")
		}
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("Chat: grpc exec failed err: %w ", err))
	}

	if resp.Code != uint32(codes.OK) {
		return nil, bizerr.ErrInternalError.Wrap(fmt.Errorf("Chat: grpc res failed err: %s ", resp.ErrMsg))
	}

	res := &ChatResponse{
		ResponseText: resp.ResponseText,
	}
	zap.S().Info(resp.ResponseText)
	return res, nil
}

func (uc *CharacterUsecase) ChatStream(ctx context.Context, req *ChatRequest) error {
	defer func() {
		close(req.ResCh)
		zap.S().Info("ChatStream end")
	}()
	conn, err := uc.chatPoll.Get(ctx)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatStream: get grpc conn err: %w", err))
	}
	defer conn.Close()

	cli := chat_agent.NewAgentClient(conn)
	grpcReq := &chat_agent.ChatRequest{
		CharacterId:    req.CharacterId,
		ConversationId: req.ConversationID,
		Message:        []byte(req.Message),
	}
	zap.S().Infof("ChatStream: req: %+v", grpcReq)
	stream, err := cli.ChatStream(ctx, grpcReq)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			err = fmt.Errorf("ChatStream: invoke chat agent timeout ")
		}
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatStream: get stream cli err: %w ", err))
	}
	msg := ""
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				zap.S().Info("stream.Recv eof")
				break
			}
			return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatStream: recv failed err: %w ", err))
		}

		if resp.Code != 0 {
			return bizerr.ErrInternalError.Wrap(fmt.Errorf("ChatStream: resp err code:%d err:%s ", resp.Code, resp.ErrMsg))
		}

		msg = fmt.Sprintf("%s%s", msg, resp.Chunk)
		req.ResCh <- resp.Chunk
	}
	zap.S().Info("stream.Recv end")
	return nil
}

func (s *CharacterUsecase) UpdateCharacter(ctx context.Context, req *UpdateCharacterRequest) error {
	zap.S().Infof("UpdateCharacter: req: %+v", *req)
	if err := s.characterRepo.UpdateCharacter(ctx, req); err != nil {
		return bizerr.ErrInternalError.Wrap(err)
	}
	return nil
}

func (s *CharacterUsecase) DeleteCharacter(ctx context.Context, id string) error {
	if err := s.characterRepo.DeleteCharacterByID(ctx, id); err != nil {
		return bizerr.ErrInternalError.Wrap(err)
	}
	return nil
}

func (uc *CharacterUsecase) Mint(ctx context.Context, id string, mint string) error {
	err := uc.characterRepo.CharacterMintSave(ctx, id, mint)
	if err != nil {
		return bizerr.ErrInternalError.Wrap(fmt.Errorf("Mint: save to db err: %w", err))
	}
	return nil
}
