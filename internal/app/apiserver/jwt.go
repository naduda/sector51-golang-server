package apiserver

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/naduda/sector51-golang/internal/app/model"
)

func (s *Server) generateJWT(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"phone": user.Phone,
		"exp":   time.Now().Add(time.Hour * 20).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Server) parseJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(s.jwtSecret), nil
	})

	if token == nil {
		return nil, errors.New("Token is nil")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return map[string]interface{}{
			"id":    claims["id"],
			"phone": claims["phone"],
		}, nil
	}

	return nil, err
}
