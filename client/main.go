package main

import (
	"context"
	"google.golang.org/grpc/grpclog"
	"log"

	"github.com/PhamDuyKhang/demo-grpc/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	AuthInterceptor struct {
		token string
	}
)

func NewAI() AuthInterceptor {
	return AuthInterceptor{}
}

//Login use user name and password to get the JWT token
func (ai *AuthInterceptor) Login(userName, password string) error {

	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}

	userS := proto.NewUserClient(conn)

	userLogin, err := userS.Login(context.Background(), &proto.LoginMessage{UserName: userName, Password: password})
	if err != nil {
		log.Fatal("can't login cause ", err)
	}

	if userLogin.GetJWTToken() == "" {
		log.Fatal("the token is empty ")
	}
	log.Print("The token is: ", userLogin.JWTToken)
	ai.token = userLogin.JWTToken
	return nil
}

//SetToken set token for every time call service
func (ai *AuthInterceptor) SetToken(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", ai.token)
	err := invoker(authCtx, method, req, reply, cc, opts...)
	return err
}

func (ai *AuthInterceptor) SetTokenStream(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	grpclog.Info("% is sending....", method)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", ai.token)
	return streamer(authCtx, desc, cc, method, opts...)

}

func main() {

	ai := NewAI()
	err := ai.Login("pdkhang", "12345")
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(ai.SetToken), grpc.WithStreamInterceptor(ai.SetTokenStream))
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}

	pokemonGo := proto.NewExampleServiceClient(conn)

	res, err := pokemonGo.Do(context.Background(), &proto.CalculatorRequest{
		Name:        "Hello",
		A:           3,
		B:           5,
		O:           proto.Operator_PLUS,
		CallingTime: timestamppb.Now(),
	})

	if err != nil {
		log.Printf("can't call server cause %v", err)
		return
	}

	log.Printf("the response from server %+v", res)

	log.Print("=========================")

}
