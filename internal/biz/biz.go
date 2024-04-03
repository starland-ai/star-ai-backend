package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAccountUsecase,
	NewCharacterUsecase,
	NewImageModelUsecase,
	NewConversationUsecase,
	NewCharacterVoiceUsecase)
