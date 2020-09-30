package main

import (
	"context"
	"fmt"
	"github.com/PhamDuyKhang/demo-grpc/proto/proto"
)

type (
	LoginService struct {
		d               *InMemStore
		tokenManagement *JWTManagement
	}
)

func NewLoginService(im *InMemStore, tokenManagement *JWTManagement) *LoginService {
	return &LoginService{tokenManagement: tokenManagement, d: im}
}

func (l LoginService) Login(ctx context.Context, message *proto.LoginMessage) (*proto.TokenPlaceHolder, error) {
	userName := message.UserName

	user, err := l.d.GetUser(userName)

	if err != nil {
		return nil, err
	}
	if user.Password != message.GetPassword() {
		return nil, fmt.Errorf("unauthorizre")
	}

	token := l.tokenManagement.Generate(user)

	return &proto.TokenPlaceHolder{JWTToken: token}, nil

}
