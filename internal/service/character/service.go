package character

import (
	"starland-backend/configs"
	"starland-backend/internal/biz"

	"time"

	"github.com/google/wire"
	"github.com/patrickmn/go-cache"
)

var ProviderSet = wire.NewSet(NewCharacterService)

type ConfirmType string

const (
	CharacterSetting ConfirmType = "character_setting"
	ImageSetting     ConfirmType = "image_setting"
)

type CharacterService struct {
	cfg          *configs.Config
	character    *biz.CharacterUsecase
	ativity      *biz.AccountAndActivitySerClientUsecase
	imageModel   *biz.ImageModelUsecase
	conversation *biz.ConversationUsecase
	voice        *biz.CharacterVoiceUsecase
	imageCache   *cache.Cache
}

func NewCharacterService(cfg *configs.Config,
	model *biz.ImageModelUsecase,
	character *biz.CharacterUsecase,
	ativity *biz.AccountAndActivitySerClientUsecase,
	conversation *biz.ConversationUsecase,
	voice *biz.CharacterVoiceUsecase) *CharacterService {
	c := cache.New(30*time.Minute, 30*time.Minute)
	s := &CharacterService{cfg: cfg,
		character:    character,
		imageModel:   model,
		ativity:      ativity,
		conversation: conversation,
		voice:        voice,
		imageCache:   c}
	go s.refreshCharacterTask()
	return s
}

type CreateCharacterRequest struct {
	AccountID string `json:"account_id"`
	Message   string `json:"message"`   
	SessionID string `json:"session_id"` 
	State     int    `json:"state"`
	Is3D      bool   `json:"is_3d"`
	ResCh     chan ChatCompletionStreamResponseChunk
}

type CreateCharacterResponse struct {
	ChatMessage []string `json:"chat_message"`
	ConfirmType string   `json:"confirm_type,omitempty"`
	ImageMeta   []string `json:"image_meta"`
	NeedConfirm bool     `json:"need_confirm"`
	SessionID   string   `json:"session_id"`
}

type ImageModelResponse struct {
	UUID                     string  `json:"uuid,omitempty"`
	NameEN                   string  `json:"nameEN,omitempty"`
	NameZH                   string  `json:"nameZH,omitempty"`
	URL                      string  `json:"URL,omitempty"`
	InferModelType           string  `json:"inferModelType,omitempty"`
	InferModelName           string  `json:"inferModelName,omitempty"`
	ComfyuiModelName         string  `json:"comfyuiModelName,omitempty"`
	InferModelDownloadURL    string  `json:"inferModelDownloadURL,omitempty"`
	InferDepModelName        string  `json:"inferDepModelName,omitempty"`
	InferDepModelDownloadURL string  `json:"inferDepModelDownloadURL,omitempty"`
	NegativePrompt           string  `json:"negativePrompt,omitempty"`
	SamplerName              string  `json:"samplerName,omitempty"`
	CfgScale                 float32 `json:"cfgScale,omitempty"`
	Steps                    int     `json:"steps,omitempty"`
	Width                    int     `json:"width,omitempty"`
	Height                   int     `json:"height,omitempty"`
	BatchSize                int     `json:"batchSize,omitempty"`
	ClipSkip                 int     `json:"clipSkip,omitempty"`
	DenoisingStrength        float32 `json:"denoisingStrength,omitempty"`
	Ensd                     int     `json:"ensd,omitempty"`
	HrUpscaler               string  `json:"hrUpscaler,omitempty"`
	EnableHr                 bool    `json:"enableHr,omitempty"`
	RestoreFaces             bool    `json:"restoreFaces,omitempty"`
	Trigger                  string  `json:"trigger,omitempty"`
	Gender                   int     `json:"gender,omitempty"`
}

type CharacterResponse struct {
	ID        string `json:"id,omitempty"`
	AccountID string `json:"account_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Gender    int    `json:"gender,omitempty"`
	Prompt    string `json:"prompt,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
}

type QueryCharacterRequest struct {
	Search  string
	Account string
	Page    int
	Limit   int
}
type QueryCharacterResponse struct {
	ID          string   `json:"id"`
	AccountID   string   `json:"account_id"`
	AccountName string   `json:"account_name"`
	AvatarURL   string   `json:"avatar_url"`
	Name        string   `json:"name"`
	Gender      int      `json:"gender"`
	Prompt      string   `json:"prompt"`
	ImageURL    string   `json:"image_url"`
	ImageURLs   []string `json:"image_urls"`
	IsMint      bool     `json:"is_mint"`
	Tags        []Tag    `json:"tags"`
	ChatCount   int      `json:"chat_count"`
	LikeCount   int      `json:"like_count"`
	Mint        string   `json:"mint"`
	IsLike      bool     `json:"is_like"`
	Is3D        bool     `json:"is_3d"`
	ObjURL      string   `json:"obj_url,omitempty"`
	GlbURL      string   `json:"glb_url,omitempty"`
	Voice       string   `json:"voice"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ChatRequest struct {
	Message     string `json:"message"`
	CharacterID string `json:"character_id"`
	AccountID   string `json:"account_id"`
}

type ChatRequestV2 struct {
	Message     string `json:"message"`
	CharacterID string `json:"character_id"`
	AccountID   string `json:"account_id"`
	ResCh       chan interface{}
}
type ChatResponse struct {
	Message string `json:"message"`
}

type QueryCharactersHistoryRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Account string `json:"account"`
}

type QueryCharactersHistoryResponse struct {
	AccountName string    `json:"account_name"`
	AvatarURL   string    `json:"avatar_url"`
	Name        string    `json:"name"`
	Gender      int       `json:"gender"`
	ImageURL    string    `json:"image_url"`
	IsMint      bool      `json:"is_mint"`
	LatestTime  time.Time `json:"latest_time"`
	CharacterID string    `json:"character_id"`
}

type CharacterMintRequest struct {
	Mint string
	ID   string
}

type ChatCompletionStreamResponseChunk struct {
	SessionID   string   `json:"session_id"`
	Message     string   `json:"message"`
	ChatMessage []string `json:"chat_message"`
	ConfirmType string   `json:"confirm_type,omitempty"`
	ImageMeta   []string `json:"image_meta"`
	NeedConfirm bool     `json:"need_confirm"`
	Is3D        bool     `json:"is_3d"`
	ObjURL      string   `json:"obj_url"`

	ChunkType         uint32                   `json:"chunk_type,omitempty"`
	ChunkSessionIndex uint32                   `json:"chunk_session_index,omitempty"`
	ChatChunk         *ChatMessage             `json:"chat_chunk,omitempty"`
	ImageChunk        []string                 `json:"image_chunk,omitempty"`
	NeedConfirmChunk  bool                     `json:"need_confirm_chunk,omitempty"`
	SettingChunk      string                   `json:"setting_chunk,omitempty"`
	CreateResponse    *CreateCharacterResponse `json:"data,omitempty"`
}

type ChatMessage struct {
	Role    string
	Content string
}

type CharacterVoiceResponse struct {
	UUID   string `json:"uuid"`
	NameZH string `json:"name_zh,omitempty"`
	NameEN string `json:"name_en,omitempty"`
	ZHUrl  string `json:"zh_url,omitempty"`
	ENUrl  string `json:"en_url,omitempty"`
}

type UpdateCharacterRequest struct {
	ID          string
	AccountID   string
	Images      []string
	Image       string
	Name        string
	Description string
	Voice       string
}

type DeleteCharacterRequest struct {
	ID        string
	AccountID string
}
