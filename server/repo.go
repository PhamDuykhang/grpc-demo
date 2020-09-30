package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const DBName = "Axx"
const Pokemon = "POKEMON"

type (
	PokemonDB struct {
		ID           string               `json:"id" bson:"id"`
		Name         string               `json:"name" bson:"name"`
		Type         string               `json:"type"bson:"type"`
		Strength     int64                `json:"strength"bson:"strength"`
		HP           int64                `json:"hp"bson:"hp"`
		Armor        int64                `json:"armor"bson:"armor"`
		Level        string               `json:"level" bson:"level"`
		Comment      string               `json:"comment" bson:"comment"`
		CatchingTime *timestamp.Timestamp `json:"catching_time",bson:"catching_time"`
		ValidateTime *timestamp.Timestamp `json:"validate_time",bson:"validate_time"`
	}

	Database struct {
		DBClient *mongo.Client
	}
)

func NewDB() Database {
	clO := options.Client()
	clO.SetAppName("Pokemon")
	clO.SetHosts([]string{"localhost:27017"})

	cl, err := mongo.NewClient(clO)
	if err != nil {
		log.Fatalf("can't init mongodb client cause %s", err)
	}
	err = cl.Connect(context.Background())
	if err != nil {
		log.Fatalf("can't connect mongodb client cause %s", err)
	}
	err = cl.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("can't ping mongodb client cause %s", err)
	}
	log.Print("every things are OK")
	return Database{DBClient: cl}
}

func (d Database) SavePokemon(p PokemonDB) error {

	_, err := d.DBClient.Database(DBName).Collection(Pokemon).InsertOne(context.Background(), p)
	if err != nil {
		log.Print("ERROR: can't insert the pokemon list ", err)
		return err
	}
	return nil
}

func (d Database) GetAll(numDoc int) ([]PokemonDB, error) {

	s := bson.M{}
	fo := options.FindOptions{}

	fo.SetBatchSize(int32(numDoc))
	var result []PokemonDB
	c, err := d.DBClient.Database(DBName).Collection(Pokemon).Find(context.Background(), s, &fo)
	if err != nil {
		log.Print("ERROR: can't get the pokemon list")
		return nil, err
	}

	err = c.All(context.Background(), &result)
	if err != nil {
		log.Print("ERROR: can't decode result pokemon list")
		return nil, err
	}
	return result, nil
}
