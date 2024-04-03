//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"starland-backend/configs"
	"starland-backend/internal/biz"
	"starland-backend/internal/data"
	"starland-backend/internal/service"
	account_service "starland-backend/internal/service/account"
	character_service "starland-backend/internal/service/character"

	"github.com/google/wire"
)

// initApp
func initApp(cfg *configs.Config) (*service.Service, error) {
	panic(wire.Build(data.ProviderSet,
		biz.ProviderSet,
		account_service.ProviderSet,
		character_service.ProviderSet,
		service.ProviderSet))
}
