package intercept

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"strings"
)

type jwt struct {
	secretKey string
}

func JwtMiddleware() *jwt {
	return &jwt{secretKey: viper.GetString("auth.jwt.secret")}
}

func (m *jwt) JwtConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Skipper: func(ctx echo.Context) bool {
			apiV1URI := "/api/v1"
			if strings.EqualFold(ctx.Request().RequestURI, apiV1URI+"/users/token") {
				return true
			}
			return false
		},
		SigningKey:    []byte(m.secretKey),
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "header:Authorization",
		AuthScheme:    "Bearer",
	}
}
