package cookie

import "net/http"

const (
	StarlandAIRedirect string = "starland-redirect"
	StarlandAIToken    string = "starland-token"
	Domain             string = "starland.ai"
)

func NewStarlandAICookie(name, value string) http.Cookie {
	return http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		Domain: Domain,
	}
}
