package utils

import (
	"encoding/json"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	jwt.StandardClaims
}

func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	claims := &Claims{
		Username: user.Username,
		ID:       user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractJWTFromRequest(e echo.Context, key string) (*Claims, error) {
	tokenString, err := extractBearerToken(e)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	userClaims, err := convertMapToClaims(claims)
	if err != nil {
		return nil, err
	}

	return userClaims, err
}

func convertMapToClaims(token map[string]interface{}) (*Claims, error) {
	var claims Claims
	tokenByte, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tokenByte, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func extractBearerToken(e echo.Context) (string, error) {
	headerAuthorization := e.Request().Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")

	if len(bearerToken) < 2 {
		return "", errors.New("invalid format token")
	}

	return html.EscapeString(bearerToken[1]), nil
}
