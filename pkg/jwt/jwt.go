package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"cinexus/config"
)

// 自定义错误
var (
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	ErrTokenMalformed   = errors.New("令牌格式错误")
	ErrTokenInvalid     = errors.New("无效的令牌")
)

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username, role string) (string, error) {
	// 设置JWT声明
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.JWT.ExpireTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Conf.JWT.Issuer,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString([]byte(config.Conf.JWT.Secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JWT.Secret), nil
	})

	if err != nil {
		// v5版本的错误处理方式有所不同
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else {
			return nil, ErrTokenInvalid
		}
	}

	// 验证令牌
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新JWT令牌
func RefreshToken(tokenString string) (string, error) {
	// 解析令牌
	claims, err := ParseToken(tokenString)
	if err != nil {
		// 如果令牌过期但其他部分有效，我们仍然可以刷新
		if errors.Is(err, ErrTokenExpired) {
			// 解析过期的令牌，忽略过期错误
			token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
				return []byte(config.Conf.JWT.Secret), nil
			}, jwt.WithoutClaimsValidation())

			if claims, ok := token.Claims.(*CustomClaims); ok {
				return GenerateToken(claims.UserID, claims.Username, claims.Role)
			}
		}
		return "", err
	}

	// 生成新令牌
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}
