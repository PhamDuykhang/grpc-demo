package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/PhamDuyKhang/demo-grpc/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	ExampleService struct {
		d *InMemStore
	}
)

func NewExampleService(d *InMemStore) ExampleService {
	return ExampleService{d: d}
}

func (e ExampleService) Do(ctx context.Context, request *proto.CalculatorRequest) (*proto.CalculatorResponse, error) {
	log.Printf("stating do with request %s", request.String())
	result := &proto.CalculatorResponse{}
	switch request.O {
	case proto.Operator_DIV:
		a := request.A
		b := request.B
		if b == 0 {
			return &proto.CalculatorResponse{}, errors.New("div by zero")
		}
		r := a / b

		result.B = a
		result.B = b
		result.O = request.O
		result.R = r
		result.Name = request.Name
		return result, nil
	case proto.Operator_MULTI:
		a := request.A
		b := request.B
		r := a * b
		result.B = a
		result.B = b
		result.O = request.O
		result.R = r
		result.Name = request.Name
		return result, nil
	case proto.Operator_PLUS:
		a := request.A
		b := request.B
		r := a + b
		result.B = a
		result.B = b
		result.O = request.O
		result.R = r
		result.Name = request.Name
		return result, nil
	case proto.Operator_SUB:
		a := request.A
		b := request.B
		r := a - b
		result.B = a
		result.B = b
		result.O = request.O
		result.R = r
		result.Name = request.Name
		return result, nil
	default:
		return &proto.CalculatorResponse{}, fmt.Errorf("can't recognize operation %v", request.O)
	}
}

func (e ExampleService) AskProfessorOak(server proto.ExampleService_AskProfessorOakServer) error {
	log.Print("AskProfessorOak")
	defer log.Print("END AskProfessorOak")
	var listPokemon []PokemonDB
	metadata.FromIncomingContext(server.Context())
	result := proto.PokemonStatResponse{}
	var g int64
	var ex int64
	var n int64
	var b int64
	for {
		pokemon, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		l, mgs := CheckLevel(pokemon.GetStrength(), pokemon.GetHP(), pokemon.GetArmor())
		listPokemon = append(listPokemon, PokemonDB{
			ID:           pokemon.Id,
			Name:         pokemon.GetName(),
			Type:         pokemon.GetType(),
			Strength:     pokemon.GetStrength(),
			HP:           pokemon.GetHP(),
			Level:        l.String(),
			Comment:      mgs,
			Armor:        pokemon.GetArmor(),
			ValidateTime: timestamppb.Now(),
			CatchingTime: pokemon.GetCatchingTime(),
		})

		switch l {
		case proto.Level_Excellent:
			ex++
		case proto.Level_Good:
			g++
		case proto.Level_Normal:
			n++
		case proto.Level_Bad:
			b++
		}
	}
	err := e.d.InsertMany(listPokemon)
	if err != nil {
		log.Print("ERROR: can't save list pokemon into database ", err)
		result.Message = err.Error()
		err := server.SendAndClose(&result)
		if err != nil {
			log.Print("ERROR: can't send the response to client ", err)
			return err
		}
	}
	result.Total = int64(len(listPokemon))
	result.Message = "nice catch"
	result.Bad = b
	result.Excellent = ex
	result.Normal = n
	result.Good = g

	err = server.SendAndClose(&result)
	if err != nil {
		log.Print("ERROR: can't send the response to client")
		return err
	}
	return nil
}

func (e ExampleService) GetPokemonDetail(pokemon *proto.GetListPokemon, server proto.ExampleService_GetPokemonDetailServer) error {
	log.Print("GetPokemonDetail")
	defer log.Print("END GetPokemonDetail")
	size := pokemon.Size

	pk := e.d.GetAllPokemon(int(size))
	for _, k := range pk {
		time.Sleep(500 * time.Millisecond)
		pokemon := &proto.CheckedPokemon{
			Name:         k.Name,
			Type:         k.Type,
			Strength:     k.Strength,
			HP:           k.HP,
			Armor:        k.Armor,
			CatchingTime: k.CatchingTime,
			Level:        parseLv(k.Level),
			Comments:     k.Comment,
		}
		err := server.Send(pokemon)
		if err != nil {
			return err
		}
	}
	return nil

}

func (e ExampleService) Talk(server proto.ExampleService_TalkServer) error {
	log.Print("TALK")
	defer log.Print("END TALK")
	t := time.Tick(1 * time.Minute)

	for {
		select {
		case <-t:
			log.Print("END OF TIME")
			return nil
		case <-server.Context().Done():
			log.Print("END REQUEST ", server.Context().Err())
			return nil
		default:
			r, err := server.Recv()
			if err != nil {
				if err == io.EOF {
					log.Print("ERROR EOF")
					return nil
				}
				log.Print("ERROR when get data", err)
			} else {
				l, m := CheckLevel(r.Strength, r.HP, r.Armor)
				log.Print("checked level ", l)
				checked := &proto.CheckedPokemon{
					Name:         r.GetName(),
					Type:         r.GetType(),
					Strength:     r.GetStrength(),
					HP:           r.GetHP(),
					Armor:        r.GetArmor(),
					Level:        l,
					Comments:     m,
					CatchingTime: r.GetCatchingTime(),
					CheckTime:    timestamppb.Now(),
				}
				log.Print("sending data ", checked)
				err = server.Send(checked)
				if err != nil {
					log.Print("ERROR when send data", err)
				}
			}
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatal(err)
	}

	inMem := NewInMemStore()
	inMem.SeedData()
	jswM := NewJWTManagement()

	interceptor := NewAuthMiddleware(jswM)

	loginService := NewLoginService(inMem, jswM)

	service := NewExampleService(inMem)

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.Logger, interceptor.Authentication), grpc.StreamInterceptor(interceptor.AuthenticationStream))

	proto.RegisterExampleServiceServer(s, service)
	proto.RegisterUserServer(s, loginService)

	reflection.Register(s)

	log.Print("starting server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func CheckLevel(s, hp, armor int64) (proto.Level, string) {
	if inRange(s, 100, 200) && inRange(hp, 100, 200) && inRange(armor, 100, 200) {
		return proto.Level_Excellent, RandomMGS()
	}
	if inRange(s, 70, 99) && inRange(hp, 70, 99) && inRange(armor, 70, 99) {
		return proto.Level_Good, RandomMGS()

	}
	if inRange(s, 50, 69) && inRange(hp, 50, 69) && inRange(armor, 50, 69) {
		return proto.Level_Normal, RandomMGS()

	}
	if inRange(s, 0, 49) && inRange(hp, 0, 49) && inRange(armor, 0, 49) {
		return proto.Level_Bad, RandomMGS()

	}
	return proto.Level_Bad, RandomMGS()

}

func RandomMGS() string {
	a := []string{"I never see it before!", "It so cute", "Oh you can give me this pokemon please", "you should jettison that pokemon!"}
	r := rand.Intn(4)
	return a[r]
}

func parseLv(l string) proto.Level {
	switch l {
	case "Excellent":
		return proto.Level_Excellent
	case "Good":
		return proto.Level_Good
	case "Normal":
		return proto.Level_Normal
	case "Bad":
		return proto.Level_Bad
	default:
		panic("not found")
	}
}

func inRange(a, s, e int64) bool {
	if a >= s && a <= e {
		return true
	} else {
		return false
	}
}
