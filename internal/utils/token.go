package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWT(uid uint, expireTime int, secret string) (string, error) {
	// 创建一个新的 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
	})

	// 签名并获取完整的 token 字符串
	return token.SignedString([]byte(secret))
}

func parseJWT(key string, jwtStr string) (uint, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		// 检查算法是否为 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(key), nil
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

	uid, ok := claims["uid"].(float64)
	if !ok {
		return 0, errors.New("invalid uid")
	}

	return uint(uid), nil
}

func GenerateAccessToken(uid uint, secret string) (string, error) {
	return generateJWT(uid, 3600, secret)
}

func ValidateAccessToken(key string, tokenStr string) (uint, error) {
	return parseJWT(key, tokenStr)
}

func GenerateRefreshToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b)[:length], nil
}
