package main

import (
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sync"
	"time"
)

var (
	ErrorNotFound     = errors.New("not found")
	ErrorDuplicateKey = errors.New("duplicate")
)

type (
	User struct {
		UserName  string
		Password  string
		CreatedAt *time.Time
		Role      string
	}

	InMemStore struct {
		rw       *sync.RWMutex
		UserData map[string]User
		Pokemon  map[string]PokemonDB
	}
)

func NewInMemStore() *InMemStore {
	rw := &sync.RWMutex{}
	userData := make(map[string]User)
	pokemon := make(map[string]PokemonDB)
	return &InMemStore{UserData: userData, Pokemon: pokemon, rw: rw}

}

func (im *InMemStore) GetUser(uID string) (User, error) {
	user, ok := im.UserData[uID]
	if !ok {
		return User{}, ErrorNotFound
	}
	return user, nil
}

func (im *InMemStore) InsertUser(us User) error {
	_, ok := im.UserData[us.UserName]
	if ok {
		return ErrorDuplicateKey
	}
	im.rw.Lock()
	defer im.rw.Unlock()

	im.UserData[us.UserName] = us
	return nil
}

func (im *InMemStore) InsertManyUser(us []User) error {
	log.Print("Data before insert ", im.UserData)
	im.rw.Lock()
	for i := range us {
		_, ok := im.UserData[us[i].UserName]
		if ok {
			return ErrorDuplicateKey
		}
		im.UserData[us[i].UserName] = us[i]
	}
	im.rw.Unlock()
	log.Print("inserting is successfully data now: ", im.UserData)
	return nil

}

func (im *InMemStore) GetPokemon(pokemonID string) (PokemonDB, error) {
	p, ok := im.Pokemon[pokemonID]
	if !ok {
		return PokemonDB{}, ErrorNotFound
	}
	return p, nil
}

func (im *InMemStore) GetAllPokemon(size int) []PokemonDB {
	var listP = []PokemonDB{}

	for _, v := range im.Pokemon {
		listP = append(listP, v)
		size--
		if size == 0 {
			break
		}
	}
	return listP
}

func (im *InMemStore) InsertOne(p PokemonDB) error {

	_, ok := im.Pokemon[p.ID]
	if ok {
		return ErrorDuplicateKey
	}
	im.rw.Lock()
	defer im.rw.Unlock()

	im.Pokemon[p.ID] = p
	log.Print("inserting is successfully data now: ", im.Pokemon)
	return nil
}

func (im *InMemStore) InsertMany(p []PokemonDB) error {
	log.Print("Data before insert ", im.Pokemon)
	im.rw.Lock()
	for i := range p {
		_, ok := im.Pokemon[p[i].ID]
		if ok {
			return ErrorDuplicateKey
		}
		im.Pokemon[p[i].ID] = p[i]
	}
	im.rw.Unlock()
	log.Print("inserting is successfully data now: ", im.Pokemon)
	return nil
}
func (im *InMemStore) SeedData() {
	tn := time.Now()
	user1 := User{
		UserName:  "pdkhang",
		Password:  "12345",
		CreatedAt: &tn,
		Role:      "admin",
	}
	user2 := User{
		UserName:  "user1",
		Password:  "user1",
		CreatedAt: &tn,
		Role:      "normal",
	}
	user3 := User{
		UserName:  "user2",
		Password:  "user2",
		CreatedAt: &tn,
		Role:      "admin",
	}

	pk1 := PokemonDB{
		ID:           "PK1",
		Name:         "Jil",
		Type:         "Fire",
		Strength:     130,
		HP:           120,
		Armor:        134,
		Level:        "Excellent",
		Comment:      "nice catch",
		CatchingTime: timestamppb.Now(),
		ValidateTime: timestamppb.Now(),
	}

	pk2 := PokemonDB{
		ID:           "PK2",
		Name:         "Jack",
		Type:         "Bug",
		Strength:     120,
		HP:           125,
		Armor:        124,
		Level:        "Excellent",
		Comment:      "nice catch",
		CatchingTime: timestamppb.Now(),
		ValidateTime: timestamppb.Now(),
	}

	pokemons := []PokemonDB{pk1, pk2}

	listUser := []User{user1, user2, user3}

	err := im.InsertManyUser(listUser)

	if err != nil {
		log.Fatal("can't init user list ", err)
	}

	err = im.InsertMany(pokemons)
	if err != nil {
		log.Fatal("can't init user list ", err)
	}
	return

}
