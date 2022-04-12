package application

import (
	"context"
	"errors"
	"strings"

	"github.com/rubumo/core/internal/database"
	"github.com/rubumo/core/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationRepositoryMongo struct{}

var collection *mongo.Collection = database.GetCollection("applications")

func (s *ApplicationRepositoryMongo) Create(a entity.Application) (entity.Application, error) {
	res, err := collection.InsertOne(context.TODO(), a)

	if err != nil {
		return entity.Application{}, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return entity.Application{}, errors.New("Unable to retrieve last inserted id.")
	}

	a.ID = id

	return a, nil
}

func (s *ApplicationRepositoryMongo) FindByDomain(w string, d string) (entity.Application, error) {
	f := bson.M{
		"environments.domain": d,
		"workspace.domain": bson.M{
			"$in": []string{
				w,
				strings.Replace(w, ".rubumo.io", "", 1),
			},
		},
	}

	return s.FindOne(f)
}

func (s *ApplicationRepositoryMongo) FindById(id string) (entity.Application, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entity.Application{}, err
	}

	f := bson.M{"_id": objectId}

	return s.FindOne(f)
}

func (s *ApplicationRepositoryMongo) FindOne(f bson.M) (entity.Application, error) {
	a := entity.Application{}

	res := collection.FindOne(context.TODO(), f)

	if err := res.Decode(&a); err != nil {
		return entity.Application{}, err
	}

	return a, nil
}

func (s *ApplicationRepositoryMongo) FindBy(f map[string]interface{}) ([]entity.Application, error) {
	res, err := collection.Find(context.TODO(), f)

	if err != nil {
		return make([]entity.Application, 0), err
	}

	c := make([]entity.Application, 15)

	err = res.All(context.TODO(), &c)

	if err != nil {
		return make([]entity.Application, 0), err
	}

	return c, nil
}
