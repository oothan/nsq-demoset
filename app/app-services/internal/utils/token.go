package utils

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/internal/model"
	"time"
)

type accessTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

func GenerateAccessToken(user *model.User, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + 60*15 //15 minutes

	claims := accessTokenCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		_applib.Sugar.Error("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

type refreshTokenCustomClaims struct {
	UserId uint64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateRefreshToken(userId uint64, key string) (*model.RefreshTokenDate, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(60*60*24*3) * time.Second) // 3 days
	tokenId, err := uuid.NewRandom()

	if err != nil {
		_applib.Sugar.Error("Failed to generate refresh token Id")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenId.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		_applib.Sugar.Error("Failed to sign refresh token string")
		return nil, err
	}

	return &model.RefreshTokenDate{
		SS:        ss,
		ID:        tokenId,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

func ValidateAccessToken(tokenString string, key *rsa.PublicKey) (*accessTokenCustomClaims, error) {
	claims := &accessTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Access token is invalid ")
	}

	claims, ok := token.Claims.(*accessTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Access token valid but couldn't parse claims ")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Reresh token is invalid ")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims ")
	}

	return claims, nil
}
