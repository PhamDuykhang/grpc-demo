package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	AuthMiddleware struct {
		jwtM *JWTManagement
	}
)

func NewAuthMiddleware(jwtM *JWTManagement) *AuthMiddleware {
	return &AuthMiddleware{jwtM: jwtM}
}

func (ai AuthMiddleware) Authentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if info.FullMethod != "/proto.ExampleService/Do" {
		return handler(ctx, req)
	}
	if !ok {
		log.Print("ERROR: get metadata err")
		return nil, fmt.Errorf("[unauthorize]:missing authentication token")
	}
	tokens := metaData.Get("authorization")
	if len(tokens) != 1 {
		log.Print("ERROR: token lens isn't valid")
		return nil, fmt.Errorf("[unauthorize]:missing authentication token")
	}
	if tokens[0] == "" {
		log.Print("ERROR: token lens is nill")
		return nil, fmt.Errorf("[unauthorize]:missing authentication token")
	}

	user, err := ai.jwtM.ValidToken(tokens[0])

	if err != nil {
		log.Printf("the token is not valid %s", tokens[0])
		return nil, fmt.Errorf("[unauthorize] the token is not valid")
	}
	log.Print("the user name ", user)
	if user.Role != "admin" {
		log.Print("ERROR:the user role is not admin")
		return nil, fmt.Errorf("unauthorize")
	}
	return handler(ctx, req)

}

func (ai AuthMiddleware) AuthenticationStream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Print(info)
	if info.FullMethod == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" {
		return handler(srv, ss)
	}
	metaData, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return fmt.Errorf("[unauthorize]:missing authentication token")
	}
	tokens := metaData.Get("authorization")
	if len(tokens) != 1 {
		return fmt.Errorf("[unauthorize]:missing authentication token")
	}
	if tokens[0] == "" {
		return fmt.Errorf("[unauthorize]:missing authentication token")
	}
	user, err := ai.jwtM.ValidToken(tokens[0])

	if err != nil {
		log.Printf("the token is not valid %s", tokens[0])
		return fmt.Errorf("[unauthorize] the token is not valid")
	}
	log.Print("the user name ", user)
	return handler(srv, ss)

}

func (ai AuthMiddleware) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tn := time.Now()
	log.Printf("Starting Func %s at time", info.FullMethod)
	h, err := handler(ctx, req)
	log.Printf("Function: %s done with Duration %s", info.FullMethod, time.Since(tn))
	return h, err

}
