package ui

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

const APIPathPrefix = "/api"

func NewAPIProxy(grpcGW *url.URL, accessToken string) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = grpcGW.Scheme
			r.URL.Host = grpcGW.Host

			// NOTE This should always be the case if the proxy handler is correct
			// configured using APIPathPrefix. If not, do nothing :shrug: and let the
			// proxy target handle it and presumably fail.
			if strings.HasPrefix(APIPathPrefix+"/", r.URL.Path) {
				newPath := strings.TrimPrefix(APIPathPrefix, r.URL.Path)
				r.URL.Path, r.URL.RawPath = newPath, url.PathEscape(r.URL.Path)
			}

			if ah := r.Header.Get("Authorization"); ah != "" {
				log.Warn().
					Str("request_url", r.URL.String()).
					Str("authorization_header", ah).
					Msg("recevied request with existing authorization header")
			}
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		},
	}
}
