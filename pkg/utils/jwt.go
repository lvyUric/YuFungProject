package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明
type Claims struct {
	UserID    string   `json:"user_id"`
	Username  string   `json:"username"`
	CompanyID string   `json:"company_id"`
	RoleIDs   []string `json:"role_ids"`
	jwt.RegisteredClaims
}

// JWTUtil JWT工具类
type JWTUtil struct {
	secretKey        []byte
	expiresIn        time.Duration
	refreshExpiresIn time.Duration
}

// NewJWTUtil 创建JWT工具实例
func NewJWTUtil(secretKey string, expiresIn, refreshExpiresIn time.Duration) *JWTUtil {
	return &JWTUtil{
		secretKey:        []byte(secretKey),
		expiresIn:        expiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}
}

// GenerateToken 生成访问令牌
func (j *JWTUtil) GenerateToken(userID, username, companyID string, roleIDs []string) (string, time.Time, error) {
	expiresAt := time.Now().Add(j.expiresIn)
	claims := Claims{
		UserID:    userID,
		Username:  username,
		CompanyID: companyID,
		RoleIDs:   roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "insurance-system",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWTUtil) GenerateRefreshToken(userID string) (string, error) {
	expiresAt := time.Now().Add(j.refreshExpiresIn)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "insurance-system",
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ParseToken 解析令牌
func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ParseRefreshToken 解析刷新令牌
func (j *JWTUtil) ParseRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("invalid refresh token")
}
