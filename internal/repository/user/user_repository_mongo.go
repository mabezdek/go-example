package user

import (
	"context"
	"errors"

	"github.com/rubumo/core/internal/database"
	"github.com/rubumo/core/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct{}

var collection *mongo.Collection = database.GetCollection("users")

func (s *UserRepositoryMongo) Create(u entity.User) (entity.User, error) {
	res, err := collection.InsertOne(context.TODO(), u)

	if err != nil {
		return entity.User{}, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return entity.User{}, errors.New("Unable to retrieve last inserted id.")
	}

	u.ID = id

	return u, nil
}

func (s *UserRepositoryMongo) FindByEmail(e string) (entity.User, error) {
	f := bson.M{"email": e}

	return s.FindOne(f)
}

func (s *UserRepositoryMongo) FindById(id string) (entity.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entity.User{}, err
	}

	f := bson.M{"_id": objectId}

	return s.FindOne(f)
}

func (s *UserRepositoryMongo) FindOne(f bson.M) (entity.User, error) {
	u := entity.User{}

	res := collection.FindOne(context.TODO(), f)

	if err := res.Decode(&u); err != nil {
		return entity.User{}, err
	}

	return u, nil
}
