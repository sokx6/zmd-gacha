package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(uid uint, expireTime int, secret string) (string, error) {
	// 创建一个新的 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
	})

	// 签名并获取完整的 token 字符串
	return token.SignedString([]byte(secret))
}

func ParseJWT(key string, jwtStr string) (uint, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	uid, ok := claims["uid"].(uint)
	if !ok {
		return 0, errors.New("invalid uid")
	}

	return uid, nil
}

func GenerateAccessToken(uid uint, secret string) (string, error) {
	return GenerateJWT(uid, 3600, secret)
}

func ValidateAccessToken(tokenStr string, secret string) (uint, error) {
	return ParseJWT(secret, tokenStr)
}

func GenerateRefreshToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b)[:length], nil
}
