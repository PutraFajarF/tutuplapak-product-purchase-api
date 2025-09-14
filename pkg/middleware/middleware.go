package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CustomClaims sesuai token Python
type CustomClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

const ctxUserKey = "user_claims"

var ErrNoToken = errors.New("missing authorization token")

func extractToken(c echo.Context) (string, error) {
	auth := c.Request().Header.Get("Authorization")
	logrus.Debugf("Authorization header: %s", auth)
	if auth == "" {
		logrus.Debug("Authorization header is empty")
		return "", ErrNoToken
	}
	parts := strings.Fields(auth)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		logrus.Debugf("Extracted token with Bearer: %s", parts[1])
		return parts[1], nil
	} else if len(parts) == 1 && strings.HasPrefix(auth, "eyJ") {
		// Assume it's a JWT token without Bearer prefix
		logrus.Debugf("Extracted token without Bearer: %s", auth)
		return auth, nil
	} else {
		logrus.Debugf("Invalid authorization header format: %v", parts)
		return "", ErrNoToken
	}
}

// JWTMiddlewareHS256 : validasi token HS256 yang dibuat oleh service Python
func JWTMiddlewareHS256(cfg *config.Config, leeway time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := extractToken(c)
			if err != nil {
				logrus.Debug("Failed to extract token")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			claims := &CustomClaims{}
			parser := jwt.NewParser(jwt.WithLeeway(leeway))
			token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					logrus.Debugf("Unexpected signing method: %v", t.Method.Alg())
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return []byte(cfg.JWT.Secret), nil
			})
			if err != nil {
				logrus.Debugf("Token parsing error: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}
			if !token.Valid {
				logrus.Debug("Token is not valid")
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			logrus.Debugf("Claims - Issuer: %s, Audience: %v", claims.Issuer, claims.Audience)

			// Validate issuer and audience
			if claims.Issuer != "tutuplapak-api" {
				logrus.Debugf("Invalid issuer: %s", claims.Issuer)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid issuer")
			}
			found := false
			for _, aud := range claims.Audience {
				if aud == "tutuplapak-services" {
					found = true
					break
				}
			}
			if !found {
				logrus.Debugf("Invalid audience: %v", claims.Audience)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid audience")
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

// Helper ambil email
func GetUserEmail(c echo.Context) (string, bool) {
	v := c.Get(ctxUserKey)
	if v == nil {
		return "", false
	}
	claims, ok := v.(*CustomClaims)
	if !ok || claims.Email == "" {
		return "", false
	}
	return claims.Email, true
}

// Helper ambil phone
func GetUserPhone(c echo.Context) (string, bool) {
	v := c.Get(ctxUserKey)
	if v == nil {
		return "", false
	}
	claims, ok := v.(*CustomClaims)
	if !ok || claims.Phone == "" {
		return "", false
	}
	return claims.Phone, true
}
