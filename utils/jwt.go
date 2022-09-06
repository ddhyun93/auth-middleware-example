package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type DeviceType string

// SecretKey secret key being used to sign tokens
var (
	SecretKey = []byte("imverystrongsecretkey")
)

func GenerateToken(id string, hour time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * hour).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken : Middleware에서 parsing해서 valid한지 체크 + Refreshtoken요청할때 ParseToken으로 사용됨
func ParseToken(tokenStr string) (string, error) {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(string)
		return id, nil
	} else {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			// claim 에서 id 를 정상적으로 parsing 할수 있는경우는 토큰 만료
			if claims["id"] != nil {
				return "", ErrExpiredToken
			}
			// 시간이 지났는데 안에가 변조된 경우
			return "", errors.New("invalid token")
		} else {
			// 시간 안지났는데 안에가 변조된 경우
			return "", errors.New("invalid token")
		}
	}
}
