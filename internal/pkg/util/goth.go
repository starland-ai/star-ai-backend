package util

import (
	"context"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"go.uber.org/zap"
)

var providerName = "provider"

func FiberGothAdapter(req *http.Request) *http.Request {
	pathParts := strings.Split(req.URL.Path, "/")
	index := 0
	for i := 0; i < len(pathParts); i++ {
		if pathParts[i] == "auth" {
			index = i + 1
			break
		}
	}

	provider := pathParts[index]
	zap.S().Infof("provider: %s",provider)
	if _, ok := goth.GetProviders()[provider]; ok {
		req = req.WithContext(context.WithValue(req.Context(), providerName, provider))
	}
	return req
}
