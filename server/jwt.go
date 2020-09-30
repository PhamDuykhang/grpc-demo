package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	JWTManagement struct {
		PrivateKey string
	}
	UserClaims struct {
		jwt.StandardClaims
		Name string
		Role string
	}
)

func NewJWTManagement() *JWTManagement {
	return &JWTManagement{PrivateKey: "Ktech"}
}

func (jw *JWTManagement) Generate(user User) string {
	jwtClaims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Ktech",
			Subject:   "For trainer",
		},
		Name: user.UserName,
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	tkString, err := token.SignedString([]byte(jw.PrivateKey))
	if err != nil {
		log.Fatal("can't sing the token ", err)
	}
	return tkString
}

func (jw *JWTManagement) ValidToken(tk string) (UserClaims, error) {

	token, err := jwt.ParseWithClaims(tk, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jw.PrivateKey), nil
	})
	if err != nil {
		log.Print("ERROR: can't valid ", err)
		return UserClaims{}, err
	}
	claims := token.Claims
	err = claims.Valid()
	if err != nil {
		log.Print("ERROR: can't valid ", err)
		return UserClaims{}, err
	}
	c, _ := claims.(*UserClaims)
	return *c, nil
}
