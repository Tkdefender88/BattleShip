package repository

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"

	"github.com/Tkdefender88/BattleShip/battlestate"
)

const (
	modelsDir = "./models"
)

type ModelRepository interface {
	FindModel(name string) (*battlestate.BsState, error)
	ListModels() ([]string, error)
	DeleteModel(name string) error
	CreateModel(name string, model *battlestate.BsState) (primitive.ObjectID, error)
}

type ModelRepo struct {
	db *mongo.Database
}

func NewModelRepo(db *mongo.Database) *ModelRepo {
	return &ModelRepo{db}
}

func (mr *ModelRepo) FindModel(name string) (*battlestate.BsState, error) {
	col := mr.db.Collection("battleStates")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var bsState *battlestate.BsState

	if err := col.FindOne(ctx, bson.M{"name": name}).Decode(bsState); err != nil {
		return nil, err
	}

	return bsState, nil
}

func (mr *ModelRepo) ListModels() ([]string, error) {
	col := mr.db.Collection("battleStates")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var names []string
	for cursor.Next(ctx) {
		var bs battlestate.BsState
		if err := cursor.Decode(&bs); err != nil {
			log.Println(err)
			continue
		}
		names = append(names, bs.Name)
	}

	return names, nil
}

func (mr *ModelRepo) CreateModel(name string, model *battlestate.BsState) (primitive.ObjectID, error) {

	col := mr.db.Collection("battleStates")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	res, err := col.InsertOne(ctx, model)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, err
	}

	return id, err
}

func (mr *ModelRepo) DeleteModel(name string) error {
	collection := mr.db.Collection("battleStates")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}

	return nil
}
