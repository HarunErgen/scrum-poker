package config

import (
	"net/http"
	"os"
	"strings"
)

type AppConfig struct {
	Env    string
	IsProd bool
	IsDev  bool
	Cookie CookieConfig
}

type CookieConfig struct {
	Secure   bool
	SameSite http.SameSite
}

var Cfg AppConfig

func init() {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	isProd := env == "prod"

	Cfg = AppConfig{
		Env:    env,
		IsProd: isProd,
		IsDev:  env == "dev",
		Cookie: CookieConfig{
			Secure:   resolveSecure(isProd),
			SameSite: resolveSameSite(isProd),
		},
	}
}

func resolveSameSite(isProd bool) http.SameSite {
	if isProd {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

func resolveSecure(isProd bool) bool {
	return isProd
}
