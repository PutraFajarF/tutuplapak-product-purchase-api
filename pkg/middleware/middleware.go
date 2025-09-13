package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CustomClaims sesuai token Python
type CustomClaims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

const ctxUserKey = "user_claims"

var ErrNoToken = errors.New("missing authorization token")

func extractToken(c echo.Context) (string, error) {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return "", ErrNoToken
	}
	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrNoToken
	}
	return parts[1], nil
}

// JWTMiddlewareHS256 : validasi token HS256 yang dibuat oleh service Python
func JWTMiddlewareHS256(cfg *config.Config, leeway time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := extractToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			claims := &CustomClaims{}
			parser := jwt.NewParser(jwt.WithLeeway(leeway))
			token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return cfg.JWT.Secret, nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set(ctxUserKey, claims)
			return next(c)
		}
	}
}

// Helper ambil user_id (sub)
func GetUserID(c echo.Context) (string, bool) {
	v := c.Get(ctxUserKey)
	if v == nil {
		return "", false
	}
	claims, ok := v.(*CustomClaims)
	if !ok || claims.Sub == "" {
		return "", false
	}
	return claims.Sub, true
}
