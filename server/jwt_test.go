package main

import (
	"testing"
	"time"
)

var jwtM = NewJWTManagement()

func Test_Generate(t *testing.T) {
	tn := time.Now()
	user := User{
		UserName:  "pdkhang",
		Password:  "1234",
		CreatedAt: &tn,
		Role:      "admin",
	}
	token := jwtM.Generate(user)
	if token == "" {
		t.Fail()
	}

	userClaim, err := jwtM.ValidToken(token)
	if err != nil {
		t.Fail()
		t.Fatal(err)
	}
	if user.UserName != userClaim.Name {
		t.Fail()
	}
	if user.Role != userClaim.Role {
		t.Fail()
	}
}
