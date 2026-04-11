package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	secretKeyPem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(secretKeyPem)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyPem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPem)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func generateJWT(uid uint, expireTime int, secretKeyPath, role string) (string, error) {
	privateKey, err := loadPrivateKey(secretKeyPath)
	if err != nil {
		return "", errors.New("无法加载密钥文件: " + err.Error())
	}
	// 创建一个新的 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
		"rol": role,
	})

	// 签名并获取完整的 token 字符串
	return token.SignedString(privateKey)
}

func parseJWT(publicKeyPath string, jwtStr string) (uint, string, error) {
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return 0, "", errors.New("无法加载公钥文件: " + err.Error())
	}

	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		// 检查算法是否为 RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return 0, "", errors.New("token expired")
	}

	uid, ok := claims["uid"].(float64)
	if !ok {
		return 0, "", errors.New("invalid uid")
	}

	role, ok := claims["rol"].(string)
	if !ok {
		return 0, "", errors.New("invalid role")
	}

	return uint(uid), role, nil
}

func GenerateUserAccessToken(uid uint, expireTime int, secretKeyPath string) (string, error) {
	return generateJWT(uid, expireTime, secretKeyPath, "user")
}

func GenerateAdminAccessToken(uid uint, expireTime int, secretKeyPath string) (string, error) {
	return generateJWT(uid, expireTime, secretKeyPath, "admin")
}

func ValidateAccessToken(publicKeyPath string, tokenStr string) (uint, string, error) {
	return parseJWT(publicKeyPath, tokenStr)
}

func GenerateRefreshToken(length int) (string, error) {
	return GenerateRandomString(length)
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b)[:length], nil
}

func GenerateTraceID(length int) (string, error) {
	return GenerateRandomString(length)
}
