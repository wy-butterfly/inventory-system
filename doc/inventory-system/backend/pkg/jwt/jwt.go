package jwt

import (
	"errors"
	"time"

	"inventory-system/config"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

// Claims 自定义的JWT载荷
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwtLib.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
	cfg := config.GlobalConfig.JWT

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(
				time.Duration(cfg.ExpireHours) * time.Hour,
			)),
			IssuedAt: jwtLib.NewNumericDate(time.Now()),
			Issuer:   "inventory-system",
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析并验证JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GlobalConfig.JWT

	token, err := jwtLib.ParseWithClaims(tokenString, &Claims{},
		func(token *jwtLib.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的Token")
}
