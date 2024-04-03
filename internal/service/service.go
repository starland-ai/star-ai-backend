package service

import (
	account_service "starland-backend/internal/service/account"
	character_service "starland-backend/internal/service/character"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewService)

type Service struct {
	Account   *account_service.AccountService
	Character *character_service.CharacterService
}

func NewService(account *account_service.AccountService,
	character *character_service.CharacterService) *Service {
	return &Service{Account: account, Character: character}
}
