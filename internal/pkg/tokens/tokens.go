package tokens

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"user-service/internal/entity"
	"user-service/internal/pkg/config"
)

func NewToken(user *entity.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(config.Token()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
