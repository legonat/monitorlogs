package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"monitorlogs/internal/db"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"os"
	"time"
)

func CreateAccessToken(user string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.AuthClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Login: user,
	})
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", erx.New(err)
	}
	return token, nil
}

func CreateRefreshToken() (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, 63)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b), nil
}

func ParseToken(tokenStr string) (string, error) {
	signKey := []byte(os.Getenv("ACCESS_SECRET"))
	tools.LogInfo("parse token")
	token, err := jwt.ParseWithClaims(tokenStr, &models.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, erx.NewError(609, fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return signKey, nil
	})
	if err != nil {
		return "", erx.New(err)
	}

	if claims, ok := token.Claims.(*models.AuthClaims); ok && token.Valid {
		return claims.Login, nil
	}
	return "", erx.NewError(610, "Invalid Access Token")
}

func RefreshSession(login string, ua string, fingerprint string, ip string, daysUntilExpire int) (string, string, error) {

	//login, err := db.CheckRefreshToken(token, fingerprint, ip)
	//if err != nil {
	//	return "", "", erx.New(err)
	//}

	newRefToken, err := CreateRefreshToken()

	err = db.WriteRefreshToken(login, newRefToken, ua, fingerprint, ip, daysUntilExpire)
	if err != nil {
		return "", "", erx.New(err)
	}

	newAccessToken, err := CreateAccessToken(login)
	if err != nil {
		return "", "", erx.New(err)
	}

	return newAccessToken, newRefToken, nil
}
