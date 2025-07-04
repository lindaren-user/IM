package utils

import (
	"go.uber.org/zap"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var jwtKey []byte

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func InitJWTKey() {
	key := viper.GetString("token.key")
	if key == "" {
		GetLogger().Error("JWT 密钥未配置，请检查 config.yaml 中的 token.key")
	}
	jwtKey = []byte(key)
}

// GenerateJWT 生成 Token
func GenerateJWT(userID uint64) (string, error) {
	expiration := viper.GetInt("token.expiration")

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseJWT 解析 Token
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		GetLogger().Error("解析 token 失败", zap.Error(err))
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
