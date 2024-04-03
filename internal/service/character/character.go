package character

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	config "starland-backend/configs"
	"starland-backend/internal/biz"
	"starland-backend/internal/pkg/bizerr"
	"starland-backend/internal/pkg/util"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	Stage1 int = iota
	Stage2
	Stage3
	Stage4
	Stage5

	AdminName = "StarLand.AI"

	Unconfirmed int = -1
)

func (s *CharacterService) CreateCharacter(ctx context.Context,
	req *CreateCharacterRequest) (*CreateCharacterResponse, error) {
	zap.S().Infof("CreateCharacter: req: %+v", req)
	account := ctx.Value("account").(string)
	err := s.ativity.QueryActivityLimit(ctx, account, biz.CreateCharacter)
	if err != nil {
		return nil, err
	}

	if req.SessionID != "" {
		switch req.State {
		case Stage1, Stage3:
			ccRes, err := s.character.ChatCompletions(ctx, &biz.ChatCompletionsRequest{
				SessionID: req.SessionID,
				Message: &biz.ChatMessage{
					Role:    "human",
					Content: req.Message,
				},
			})
			if err != nil {
				return nil, fmt.Errorf("CreateCharacter: grpc exec ChatCompletions err: %w", err)
			}

			if ccRes.ImageMetas != nil {
				filePath := config.GetConfig().File.ImagesEndpoint
				for i := range ccRes.ImageMetas {
					ccRes.ImageMetas[i] = fmt.Sprintf("%s%s", filePath, ccRes.ImageMetas[i])
				}
			}

			return &CreateCharacterResponse{
				ChatMessage: makeChatMessage(ccRes.Message),
				ImageMeta:   ccRes.ImageMetas,
				ConfirmType: ccRes.ConfirmType,
				NeedConfirm: ccRes.NeedConfirm,
				SessionID:   req.SessionID,
			}, nil
		case Stage2:
			ccsRes, err := s.character.ConfirmCharacterSetting(ctx, &biz.ConfirmCharacterSettingRequest{
				SessionID: req.SessionID,
			})
			if err != nil {
				return nil, fmt.Errorf("CreateCharacter: grpc exec ConfirmCharacterSetting err: %w", err)
			}
			gender := 0
			if strings.Contains(strings.ToLower(ccsRes.CharacterSetting.Gender), "man") {
				gender = 1
			} else {
				gender = 2
			}

			accountInfo, err := s.ativity.QueryAccount(ctx, account)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacter: query account err: %w", err)
			}

			character := &biz.CharacterRequest{
				ID:          req.SessionID,
				AccountID:   account,
				Gender:      gender,
				Prompt:      ccsRes.CharacterSetting.Description,
				Name:        ccsRes.CharacterSetting.Name,
				AccountName: accountInfo.Name,
				AvatarURL:   accountInfo.AvatarURL,
				Tags:        ccsRes.CharacterSetting.Tags,
				State:       Unconfirmed,
			}

			_, err = s.character.SaveMyCharacter(ctx, character)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacter: save character err: %w", err)
			}

			return &CreateCharacterResponse{
				ChatMessage: makeChatMessage(ccsRes.Message),
				ConfirmType: ccsRes.ConfirmType,
				NeedConfirm: ccsRes.NeedConfirm,
				SessionID:   req.SessionID,
			}, nil
		case Stage4:
			character := &biz.CharacterRequest{
				ID:       req.SessionID,
				ImageURL: req.Message,
				State:    0,
			}
			_, err := s.character.SaveMyCharacter(ctx, character)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacter: save character err: %w", err)
			}
			err = s.ativity.PostActivity(ctx, account, biz.CreateCharacter)
			if err != nil {
				zap.S().Errorf("CreateCharacter: [Account: %s] add points err: %w ", account, err)
			}
		}
		return nil, nil
	} else {
		return &CreateCharacterResponse{
			SessionID: uuid.NewString(),
		}, nil
	}
}

func (s *CharacterService) CreateCharacterV2(ctx context.Context,
	req *CreateCharacterRequest) (*CreateCharacterResponse, error) {
	zap.S().Infof("CreateCharacterV2: req: %+v", req)
	account := req.AccountID
	err := s.ativity.QueryActivityLimit(ctx, account, biz.CreateCharacter)
	if err != nil {
		return nil, err
	}
	defer func() {
		close(req.ResCh)
	}()
	var content string

	ctx = context.Background()
	if req.SessionID != "" {
		switch req.State {
		case Stage1, Stage4:
			resCh := make(chan biz.ChatCompletionStreamResponseChunk)
			go func() {
				err = s.character.ChatCompletionsStream(context.Background(), &biz.ChatCompletionsRequest{
					SessionID: req.SessionID,
					ResCh:     resCh,
					Message: &biz.ChatMessage{
						Role:    "human",
						Content: req.Message,
					},
				})
				if err != nil {
					zap.S().Errorf("CreateCharacterV2: grpc exec ChatCompletions err: %v", err)
				}
			}()
			var voice string
			ch, _ := s.character.QueryCharacterByID(ctx, req.SessionID)
			if ch != nil {
				voice = ch.Voice
			}
			historyRes := &ChatCompletionStreamResponseChunk{
				SessionID:   req.SessionID,
				ConfirmType: "character_setting",
				ChatChunk:   &ChatMessage{},
			}

			for {
				if res, ok := <-resCh; ok {
					var is3D bool
					if res.ChunkType == biz.ImageChunk {
						is3D = res.Is3D
						filePath := config.GetConfig().File.ImagesEndpoint
						for i := range res.ImageChunk {
							res.ImageChunk[i] = fmt.Sprintf("%s%s.png", filePath, res.ImageChunk[i])
						}
						if ims, ok := s.imageCache.Get(req.SessionID); ok {
							list := ims.([]string)
							list = append(list, res.ImageChunk...)
							err = s.imageCache.Add(req.SessionID, list, time.Hour)
							if err != nil {
								zap.S().Errorf("CreateCharacterV2: gen image to cache")
							}
						} else {
							err = s.imageCache.Add(req.SessionID, res.ImageChunk, time.Hour)
							if err != nil {
								zap.S().Errorf("CreateCharacterV2: gen image to cache")
							}
						}
					}
					data := makeChatCompletionStreamResponseChunk(historyRes, req.SessionID, voice, is3D, res)
					req.ResCh <- data
					if res.ChatChunk != nil {
						content = fmt.Sprintf("%s%s", content, data.Message)
					}
				} else {
					historyRes.Message = content
					req.ResCh <- *historyRes
					return &CreateCharacterResponse{
						SessionID: req.SessionID,
					}, nil
				}
			}
		case Stage2:
			ccsRes, err := s.character.ConfirmCharacterSetting(context.Background(), &biz.ConfirmCharacterSettingRequest{
				SessionID: req.SessionID,
			})
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: grpc exec ConfirmCharacterSetting err: %w", err)
			}

			if ccsRes.CharacterSetting == nil {
				return nil, bizerr.ErrInternalError.Wrap(errors.New("Frequent operation, please try again"))
			}

			gender := 0
			if strings.Contains(strings.ToLower(ccsRes.CharacterSetting.Gender), "man") {
				gender = 1
			} else {
				gender = 2
			}
			accountInfo, err := s.ativity.QueryAccount(ctx, account)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: query account err: %w", err)
			}
			character := &biz.CharacterRequest{
				ID:           req.SessionID,
				Name:         ccsRes.CharacterSetting.Name,
				AccountID:    account,
				Gender:       gender,
				AccountName:  accountInfo.Name,
				AvatarURL:    accountInfo.AvatarURL,
				Tags:         ccsRes.CharacterSetting.Tags,
				Introduction: ccsRes.CharacterSetting.Introduction,
				Is3D:         req.Is3D,
				State:        Unconfirmed,
			}

			_, err = s.character.SaveMyCharacter(ctx, character)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: save character err: %w", err)
			}
			res := &CreateCharacterResponse{
				ChatMessage: makeChatMessage(ccsRes.Message),
				ConfirmType: ccsRes.ConfirmType,
				NeedConfirm: ccsRes.NeedConfirm,
				SessionID:   req.SessionID,
			}
			data := ChatCompletionStreamResponseChunk{
				CreateResponse: res,
			}
			req.ResCh <- data
			return res, nil
		case Stage3:
			_, err := s.voice.QueryCharacterVoice(ctx, req.Message)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: check voice err: %w", err)
			}
			character := &biz.CharacterRequest{
				ID:    req.SessionID,
				Voice: req.Message,
				State: Unconfirmed,
			}
			_, err = s.character.SaveMyCharacter(ctx, character)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: save character voice err: %w", err)
			}
			res := ChatCompletionStreamResponseChunk{
				NeedConfirmChunk: true,
			}
			req.ResCh <- res
			return nil, nil
		case Stage5:
			info, err := s.character.QueryCharacterByID(ctx, req.SessionID)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: check character err: %w", err)
			}
			if info == nil {
				return nil, bizerr.ErrCharacterNotExist
			}

			character := &biz.CharacterRequest{
				ID:       req.SessionID,
				ImageURL: req.Message,
				Is3D:     req.Is3D,
				State:    0,
			}
			_, err = s.character.SaveMyCharacter(ctx, character)
			if err != nil {
				return nil, fmt.Errorf("CreateCharacterV2: save character err: %w", err)
			}
			if ims, ok := s.imageCache.Get(req.SessionID); ok {
				imList := ims.([]string)
				for i := range imList {
					if !strings.Contains(req.Message, imList[i]) {
						file := fmt.Sprintf("%s%s", s.cfg.File.ImagePath, imList[i])
						os.Remove(file)
					}
				}
			}
			err = s.ativity.PostActivity(ctx, account, biz.CreateCharacter)
			if err != nil {
				zap.S().Errorf("CreateCharacterV2: [Account: %s] add points err: %w ", account, err)
			}
			data := ChatCompletionStreamResponseChunk{
				SessionID: req.SessionID,
			}
			req.ResCh <- data
		}
		return nil, nil
	} else {
		req.ResCh <- ChatCompletionStreamResponseChunk{
			SessionID: uuid.NewString(),
		}
		return nil, nil
	}
}

func (s *CharacterService) QueryImageModel(ctx context.Context) ([]*ImageModelResponse, error) {
	ims, err := s.imageModel.QueryImageModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("QueryImageModel: qeury all image model err: %w", err)
	}
	return makeImageModelResponse(ims), nil
}

func (s *CharacterService) CharacterLike(ctx context.Context, characterID, accountID string, flag bool) error {

	err := s.ativity.QueryActivityLimit(ctx, accountID, biz.Like)
	if err != nil {
		return err
	}

	ch, err := s.character.QueryCharacterByID(ctx, characterID)
	if err != nil {
		return fmt.Errorf("CharacterLike: query err: %w", err)
	}
	go func() {
		err := s.ativity.PostActivity(ctx, accountID, biz.Like)
		if err != nil {
			zap.S().Infof("CharacterLike: exec ativity like err: %w", err)
		}
	}()
	if err = s.character.UpdateCharacterLikeAccountLike(ctx, characterID, accountID, flag); err != nil {
		return fmt.Errorf("CharacterLike: update err: %w", err)
	}

	account, err := s.character.QueryCharacterLikeCountByAccount(ctx, ch.ID)
	if err != nil {
		return fmt.Errorf("CharacterLike: (character_id:%s  account_id:%s)query like count err: %w", ch.ID, ch.AccountID, err)
	}

	req := &biz.CharacterRequest{
		ID:        ch.ID,
		AccountID: ch.AccountID,
		LikeCount: int(account),
	}

	_, err = s.character.SaveMyCharacter(ctx, req)
	if err != nil {
		return fmt.Errorf("refreshCharacterExec: (character_id:%s  account_id:%s) save like count err: %w", ch.ID, ch.AccountID, err)
	}
	return nil
}

func (s *CharacterService) QueryCharacterByID(ctx context.Context, characterID string) (*CharacterResponse, error) {
	res, err := s.character.QueryCharacterByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("QueryCharacterByID: update err: %w", err)
	}
	return makeCharacterResponse(res), nil
}

func (s *CharacterService) refreshCharacterTask() {
	defer func() {
		if p := recover(); p != nil {
			zap.S().Errorf("refreshCharacterTask: recover err: %w", p)
		}
		s.refreshCharacterTask()
	}()

	t := time.NewTicker(5 * time.Minute)
	for range t.C {
		ctx := context.Background()
		characters, count, err := s.character.QueryCharactersByAccount(ctx, "", "", 1, 10)
		if err != nil {
			zap.S().Errorf("refreshCharacterTask: query characters err: %v", err)
			continue
		}
		err = s.refreshCharacterExec(ctx, characters)
		if err != nil {
			zap.S().Errorf("refreshCharacterTask: first refreshCharacterExec err: %v", err)
			continue
		}
		count -= 10
		for i := 2; count > 0; i++ {
			characters, _, err = s.character.QueryCharactersByAccount(ctx, "", "", i, 10)
			if err != nil {
				zap.S().Errorf("refreshCharacterTask: query characters err: %v", err)
				continue
			}
			err = s.refreshCharacterExec(ctx, characters)
			if err != nil {
				zap.S().Errorf("refreshCharacterTask: refreshCharacterExec err: %v", err)
				continue
			}
			count -= 10
		}
	}
}

func (s *CharacterService) refreshCharacterExec(ctx context.Context, characters []*biz.CharacterResponse) error {
	zap.S().Infof("refreshCharacterExec: len: %d", len(characters))
	for i := range characters {
		if characters[i].AccountID == "" {
			continue
		}
		accountRes, err := s.ativity.QueryAccount(ctx, characters[i].AccountID)
		if err != nil {
			return fmt.Errorf("refreshCharacterExec: (character_id:%s  account_id:%s)query accountInfo err: %w",
				characters[i].ID, characters[i].AccountID, err)
		}
		req := &biz.CharacterRequest{
			ID:          characters[i].ID,
			AccountID:   characters[i].AccountID,
			AccountName: accountRes.Name,
			AvatarURL:   accountRes.AvatarURL,
		}

		_, err = s.character.SaveMyCharacter(ctx, req)
		if err != nil {
			return fmt.Errorf("refreshCharacterExec: (character_id:%s  account_id:%s) save like count err: %w",
				characters[i].ID, characters[i].AccountID, err)
		}
	}
	return nil
}

func (s *CharacterService) QueryCharacters(ctx context.Context, req *QueryCharacterRequest) ([]*QueryCharacterResponse, int64, error) {
	var (
		characters []*biz.CharacterResponse
		err        error
		count      int64
		accountID  string
	)

	zap.S().Infof("QueryCharacters: req: %+v", *req)
	account := ctx.Value("account")
	if account != nil {
		zap.S().Info("accountID:", account)
		accountID = account.(string)
	}

	zap.S().Infof("QueryCharacters: query req: %+v", *req)
	if req.Account != "" {
		characters, count, err = s.character.QueryCharactersByAccount(ctx, req.Account, req.Search, req.Page, req.Limit)
	} else {
		characters, count, err = s.character.QueryCharactersByNameOrPrompt(ctx, req.Search, req.Page, req.Limit)
	}
	if err != nil {
		return nil, count, fmt.Errorf("QueryCharacters: query err: %w", err)
	}
	res := s.makeQueryCharacterResponse(accountID, characters)
	sort.Slice(res, func(i, j int) bool {
		if res[i].AccountName == AdminName && res[j].AccountName == AdminName {
			return res[i].LikeCount+res[i].ChatCount >= res[j].LikeCount+res[j].ChatCount
		} else if res[j].AccountName == AdminName {
			return false
		} else {
			return res[i].LikeCount+res[i].ChatCount >= res[j].LikeCount+res[j].ChatCount
		}
	})
	return res, count, nil
}

func (s *CharacterService) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	account := ctx.Value("account").(string)
	err := s.ativity.QueryActivityLimit(ctx, account, biz.Chat)
	if err != nil {
		return nil, err
	}

	_, err = s.character.QueryCharacterByID(ctx, req.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("Chat: [CharacterId: %s Message: %s] query characterId err: %w ", req.CharacterID, req.Message, err)
	}

	cr, err := s.conversation.QueryConversation(ctx, req.AccountID, req.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("Chat: [CharacterID: %s Message: %s] query conversationId err: %w ", req.CharacterID, req.Message, err)
	}
	if cr == nil {
		cr = &biz.ConversationResponse{
			ConversationID: "",
		}
		var count int64
		count, err = s.conversation.QueryConversationsCount(ctx, req.CharacterID)
		if err != nil {
			return nil, fmt.Errorf("Chat: (character_id:%s  account_id:%s)query chat count err: %w", req.CharacterID, req.AccountID, err)
		}

		chReq := &biz.CharacterRequest{
			ID:        req.CharacterID,
			AccountID: req.AccountID,
			ChatCount: int(count) + 1,
		}

		_, err = s.character.SaveMyCharacter(ctx, chReq)
		if err != nil {
			return nil, fmt.Errorf("Chat: (character_id:%s  account_id:%s) save like count err: %w", req.CharacterID, req.AccountID, err)
		}
	}

	conversationID, err := s.conversation.SaveConversation(ctx, req.AccountID, req.CharacterID, cr.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("Chat: save conversation err: %w ", err)
	}

	cr = &biz.ConversationResponse{
		CharacterID:    req.CharacterID,
		ConversationID: conversationID,
	}
	res, err := s.character.Chat(ctx, &biz.ChatRequest{
		ConversationID: cr.ConversationID,
		CharacterId:    cr.CharacterID,
		Message:        req.Message,
	})
	if err != nil {
		return nil, fmt.Errorf("Chat: [ConversationId: %s Message: %s] err: %w ", cr.ConversationID, req.Message, err)
	}

	err = s.ativity.PostActivity(ctx, req.AccountID, biz.Chat)
	if err != nil {
		zap.S().Errorf("Chat: [Account: %s] add points err: %w ", req.AccountID, err)
	}

	return &ChatResponse{
		Message: res.ResponseText,
	}, nil
}

func (s *CharacterService) ChatV2(ctx context.Context, req *ChatRequestV2) error {
	account := req.AccountID
	err := s.ativity.QueryActivityLimit(ctx, account, biz.Chat)
	if err != nil {
		return err
	}

	ch, err := s.character.QueryCharacterByID(ctx, req.CharacterID)
	if err != nil {
		return fmt.Errorf("ChatV2: [CharacterId: %s Message: %s] query characterId err: %w ", req.CharacterID, req.Message, err)
	}

	cr, err := s.conversation.QueryConversation(ctx, req.AccountID, req.CharacterID)
	if err != nil {
		return fmt.Errorf("ChatV2: [CharacterID: %s Message: %s] query conversationId err: %w ", req.CharacterID, req.Message, err)
	}
	if cr == nil {
		cr = &biz.ConversationResponse{
			ConversationID: "",
		}
		var count int64
		count, err = s.conversation.QueryConversationsCount(ctx, req.CharacterID)
		if err != nil {
			return fmt.Errorf("ChatV2: (character_id:%s  account_id:%s)query chat count err: %w", req.CharacterID, req.AccountID, err)
		}

		chReq := &biz.CharacterRequest{
			ID:        req.CharacterID,
			ChatCount: int(count) + 1,
		}

		_, err = s.character.SaveMyCharacter(ctx, chReq)
		if err != nil {
			return fmt.Errorf("ChatV2: (character_id:%s  account_id:%s) save like count err: %w", req.CharacterID, req.AccountID, err)
		}
	}

	conversationID, err := s.conversation.SaveConversation(ctx, req.AccountID, req.CharacterID, cr.ConversationID)
	if err != nil {
		return fmt.Errorf("ChatV2: save conversation err: %w ", err)
	}

	cr = &biz.ConversationResponse{
		CharacterID:    req.CharacterID,
		ConversationID: conversationID,
	}
	go func() {
		err = s.character.ChatStream(context.Background(), &biz.ChatRequest{
			ConversationID: cr.ConversationID,
			CharacterId:    cr.CharacterID,
			Message:        req.Message,
			ResCh:          req.ResCh,
			CharacterName:  ch.Name,
		})
		if err != nil {
			zap.S().Errorf("ChatV2: [ConversationId: %s Message: %s] err: %v ", cr.ConversationID, req.Message, err)
		}
	}()

	err = s.ativity.PostActivity(ctx, req.AccountID, biz.Chat)
	if err != nil {
		zap.S().Errorf("Chat: [Account: %s] add points err: %w ", req.AccountID, err)
	}

	return nil
}
func (s *CharacterService) QueryCharactersHistory(ctx context.Context,
	req *QueryCharactersHistoryRequest) ([]*QueryCharactersHistoryResponse, int64, error) {
	cr, count, err := s.conversation.QueryConversations(ctx, req.Account, req.Page, req.Limit)
	if err != nil {
		return nil, count, fmt.Errorf("QueryCharactersHistory: query conversations err: %w ", err)
	}

	res := make([]*QueryCharactersHistoryResponse, 0)
	for i := range cr {
		character, err := s.character.QueryCharacterByID(ctx, cr[i].CharacterID)
		if err != nil {
			zap.S().Errorf("QueryCharactersHistory: query character err: %w ", err)
			continue
		}
		res = append(res, &QueryCharactersHistoryResponse{
			AccountName: character.AccountName,
			AvatarURL:   character.AvatarURL,
			Name:        character.Name,
			Gender:      character.Gender,
			ImageURL:    character.ImageURL,
			IsMint:      character.IsMint,
			LatestTime:  cr[i].UpdateTime,
			CharacterID: character.ID,
		})
	}

	return res, count, nil
}

func (s *CharacterService) QueryCharacterInfo(ctx context.Context, id string) (*QueryCharacterResponse, error) {
	cr, err := s.character.QueryCharacterByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("QueryCharacterInfo: query character info err: %w", err)
	}

	tags := make([]Tag, len(cr.Tag))
	for i := range cr.Tag {
		tags[i] = Tag{
			Key:   cr.Tag[i].Key,
			Value: cr.Tag[i].Value,
		}
	}
	var accountID string
	account := ctx.Value("account")
	if account != nil {
		zap.S().Info("accountID:", account)
		accountID = account.(string)
	}

	isLike, err := s.character.QueryCharacterLikeByAccount(context.Background(), id, accountID)
	if err != nil {
		zap.S().Errorf("makeQueryCharacterResponse: err: %v", err)
	}

	res := &QueryCharacterResponse{
		ID:          cr.ID,
		AccountID:   cr.AccountID,
		AccountName: cr.AccountName,
		AvatarURL:   cr.AvatarURL,
		Name:        cr.Name,
		Gender:      cr.Gender,
		ImageURL:    cr.ImageURL,
		IsMint:      cr.IsMint,
		Prompt:      cr.Introduction,
		Tags:        tags,
		Mint:        cr.Mint,
		IsLike:      isLike,
		LikeCount:   cr.LikeCount,
		ChatCount:   cr.ChatCount,
		ImageURLs:   cr.ImageURLs,
		Is3D:        cr.Is3D,
		Voice:       cr.Voice,
	}
	if cr.Is3D {
		res.ObjURL = strings.Replace(res.ImageURL, ".png", ".obj", 1)
		res.GlbURL = strings.Replace(res.ImageURL, ".png", ".glb", 1)
	}
	zap.S().Infof("res: ", account, cr.AccountID)
	return res, nil
}

func (s *CharacterService) CharacterMint(ctx context.Context, req *CharacterMintRequest) error {
	zap.S().Infof("CharacterMint: req: %+v", *req)
	if _, err := s.character.QueryCharacterByID(ctx, req.ID); err != nil {
		return fmt.Errorf("CharacterMint: query character err: %w", err)
	}

	err := s.character.Mint(ctx, req.ID, req.Mint)
	if err != nil {
		return fmt.Errorf("CharacterMint: save mint err: %w", err)
	}
	return nil
}

func (s *CharacterService) MessageToVoice(ctx context.Context, id, msg string) (string, error) {

	if 1 == 1 {
		ch, err := s.character.QueryCharacterByID(ctx, id)
		if err != nil {
			zap.S().Errorf("MessageToVoice: query voice by id : %v", err)
			return "", nil
		}
		voice, err := s.voice.QueryCharacterVoice(ctx, ch.Voice)
		if err != nil {
			zap.S().Errorf("MessageToVoice: query voice by id : %v", err)
			return "", nil
		}
		var roleID string
		if util.ContainsChinese(msg) {
			roleID = voice.ZHRoleID
		} else {
			roleID = voice.ENRoleID
		}
		res, err := util.GenerateVoice(msg, roleID)
		if err != nil {
			return "", fmt.Errorf("MessageToVoice: gen voice : %w", err)
		}
		return config.GetConfig().File.VoiceEndpoint + res, nil
	}
	return "", nil
}

func (s *CharacterService) QueryVoice(ctx context.Context) ([]*CharacterVoiceResponse, error) {

	voices, err := s.voice.QueryAllCharacterVoice(ctx)
	if err != nil {
		return nil, fmt.Errorf("QueryVoice: query voice err: %w", err)
	}

	return makeCharacterVoiceResponse(voices), nil
}

func (s *CharacterService) UpdateCharacter(ctx context.Context, req *UpdateCharacterRequest) error {
	zap.S().Infof("UpdateCharacter: req: %+v", *req)
	ch, err := s.character.QueryCharacterByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("UpdateCharacter: query character err: %w", err)
	}
	if ch.AccountID != req.AccountID {
		return bizerr.ErrNoPermissionToModify
	}
	zap.S().Info(req.Images[0])

	sort.Slice(req.Images, func(i, j int) bool {
		if req.Images[i] == req.Image {
			return true
		} else {
			return false
		}
	})
	zap.S().Info(req.Images[0])
	if err := s.character.UpdateCharacter(ctx, &biz.UpdateCharacterRequest{
		ID:          req.ID,
		Images:      req.Images,
		Description: req.Description,
		Name:        req.Name,
		Image:       req.Image,
		Voice:       req.Voice,
	}); err != nil {
		return fmt.Errorf("UpdateCharacter: update err: %w", err)
	}
	return nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, req *DeleteCharacterRequest) error {

	ch, err := s.character.QueryCharacterByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("DeleteCharacter: query character err: %w", err)
	}
	if ch.AccountID != req.AccountID {
		return bizerr.ErrNoPermissionToModify
	}

	err = s.character.DeleteCharacter(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("DeleteCharacter: delete character err: %w", err)
	}
	return nil
}

func makeChatMessage(req []*biz.ChatMessage) []string {
	res := make([]string, len(req))
	for i := range req {
		res[i] = req[i].Content
	}
	return res
}

func makeImageModelResponse(req []*biz.ImageModel) []*ImageModelResponse {
	res := make([]*ImageModelResponse, len(req))

	for i := range req {
		res[i] = &ImageModelResponse{
			UUID:                     req[i].UUID,
			NameEN:                   req[i].NameEN,
			NameZH:                   req[i].NameZH,
			URL:                      req[i].URL,
			InferModelType:           req[i].InferModelType,
			InferModelName:           req[i].InferModelName,
			ComfyuiModelName:         req[i].ComfyuiModelName,
			InferModelDownloadURL:    req[i].InferModelDownloadURL,
			InferDepModelName:        req[i].InferDepModelName,
			InferDepModelDownloadURL: req[i].InferDepModelDownloadURL,
			NegativePrompt:           req[i].NegativePrompt,
			SamplerName:              req[i].SamplerName,
			CfgScale:                 req[i].CfgScale,
			Steps:                    req[i].Steps,
			Width:                    req[i].Width,
			Height:                   req[i].Height,
			BatchSize:                req[i].BatchSize,
			ClipSkip:                 req[i].ClipSkip,
			DenoisingStrength:        req[i].DenoisingStrength,
			Ensd:                     req[i].Ensd,
			HrUpscaler:               req[i].HrUpscaler,
			EnableHr:                 req[i].EnableHr,
			RestoreFaces:             req[i].RestoreFaces,
			Trigger:                  req[i].Trigger,
			Gender:                   req[i].Gender,
		}
	}
	return res
}

func makeCharacterResponse(req *biz.CharacterResponse) *CharacterResponse {
	return &CharacterResponse{
		ID:        req.ID,
		AccountID: req.AccountID,
		Name:      req.Name,
		Gender:    req.Gender,
		Prompt:    req.Introduction,
		ImageURL:  req.ImageURL,
	}
}

func (s *CharacterService) makeQueryCharacterResponse(accountID string, req []*biz.CharacterResponse) []*QueryCharacterResponse {
	res := make([]*QueryCharacterResponse, len(req))
	for i := range req {
		isLike, err := s.character.QueryCharacterLikeByAccount(context.Background(), req[i].ID, accountID)
		if err != nil {
			zap.S().Errorf("makeQueryCharacterResponse: err: %w", err)
		}
		res[i] = &QueryCharacterResponse{
			ID:          req[i].ID,
			AccountID:   req[i].AccountID,
			AccountName: req[i].AccountName,
			AvatarURL:   req[i].AvatarURL,
			Name:        req[i].Name,
			Gender:      req[i].Gender,
			Prompt:      req[i].Introduction,
			ImageURL:    req[i].ImageURL,
			IsMint:      req[i].IsMint,
			ChatCount:   req[i].ChatCount,
			LikeCount:   req[i].LikeCount,
			IsLike:      isLike,
			Is3D:        req[i].Is3D,
			ImageURLs:   req[i].ImageURLs,
			Voice:       req[i].Voice,
		}
	}
	return res
}

func makeChatCompletionStreamResponseChunk(historyRes *ChatCompletionStreamResponseChunk, session_id, voice string, Is3D bool, req biz.ChatCompletionStreamResponseChunk) ChatCompletionStreamResponseChunk {
	var settingChunk struct {
		Description string `json:"description"`
		Gender      string `json:"gender"`
		Name        string `json:"name"`
		Voice       string `json:"voice,omitempty"`
	}

	switch req.ChunkType {
	case biz.ChatChunk:
		historyRes.Message = req.ChatChunk.Content
	case biz.ImageChunk:
		historyRes.Is3D = Is3D
		historyRes.ImageMeta = req.ImageChunk
		if historyRes.Is3D {
			historyRes.ObjURL = strings.Replace(req.ImageChunk[0], ".png", ".obj", 1)
		}

		historyRes.ConfirmType = "image_setting"
		historyRes.NeedConfirm = true

		zap.S().Info("ImageMeta:", historyRes.ImageMeta)
	case biz.NeedConfirmChunk:
		if !historyRes.NeedConfirm {
			historyRes.NeedConfirm = req.NeedConfirmChunk
		}
	case biz.SettingChunk:
		if req.SettingChunk != "" {
			err := json.Unmarshal([]byte(req.SettingChunk), &settingChunk)
			if err != nil {
				zap.S().Errorf("makeChatCompletionStreamResponseChunk: json Unmarshal err: %v", err)
			}
			settingChunk.Voice = voice
			jsonData, err := json.Marshal(settingChunk)
			if err != nil {
				zap.S().Errorf("makeChatCompletionStreamResponseChunk: json Unmarshal err: %v", err)
			}
			req.SettingChunk = string(jsonData)
		}
		historyRes.SettingChunk = req.SettingChunk

	}
	return *historyRes
}

func makeCharacterVoiceResponse(req []*biz.CharacterVoice) []*CharacterVoiceResponse {
	res := make([]*CharacterVoiceResponse, len(req))

	for i := range req {
		res[i] = &CharacterVoiceResponse{
			UUID:   req[i].UUID,
			NameZH: req[i].NameZH,
			NameEN: req[i].NameEN,
			ZHUrl:  req[i].ZHUrl,
			ENUrl:  req[i].ENUrl,
		}
	}
	return res
}
