package workspace

import (
	"context"
	"errors"
	"strings"

	"github.com/rubumo/core/internal/aggregate"
	"github.com/rubumo/core/internal/database"
	"github.com/rubumo/core/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkspaceRepositoryMongo struct{}

var collection *mongo.Collection = database.GetCollection("workspaces")

func (s *WorkspaceRepositoryMongo) Create(u entity.Workspace) (entity.Workspace, error) {
	res, err := collection.InsertOne(context.TODO(), u)

	if err != nil {
		return entity.Workspace{}, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return entity.Workspace{}, errors.New("Unable to retrieve last inserted id.")
	}

	u.ID = id

	return u, nil
}

func (s *WorkspaceRepositoryMongo) FindByDomain(d string) (entity.Workspace, error) {
	f := bson.M{
		"domain": bson.M{
			"$in": []string{
				d,
				strings.Replace(d, ".rubumo.io", "", 1),
			},
		},
	}

	return s.FindOne(f)
}

func (s *WorkspaceRepositoryMongo) FindByDomainWithApplications(d string) (aggregate.WorkspaceWithApplications, error) {
	p := mongo.Pipeline{
		bson.D{
			primitive.E{
				Key: "$match",
				Value: bson.M{
					"domain": bson.M{
						"$in": []string{
							d,
							strings.Replace(d, ".rubumo.io", "", 1),
						},
					},
				},
			},
		},
		bson.D{
			primitive.E{
				Key:   "$limit",
				Value: 1,
			},
		},
		bson.D{
			primitive.E{
				Key: "$lookup",
				Value: bson.D{
					primitive.E{
						Key:   "from",
						Value: "applications",
					},
					primitive.E{
						Key:   "localField",
						Value: "_id",
					},
					primitive.E{
						Key:   "foreignField",
						Value: "workspace._id",
					},
					primitive.E{
						Key:   "as",
						Value: "applications",
					},
				},
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), p)

	if err != nil {
		return aggregate.WorkspaceWithApplications{}, err
	}

	var l []aggregate.WorkspaceWithApplications

	if err = cursor.All(context.Background(), &l); err != nil || len(l) == 0 {
		return aggregate.WorkspaceWithApplications{}, nil
	}

	return l[0], nil
}

func (s *WorkspaceRepositoryMongo) FindById(id string) (entity.Workspace, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entity.Workspace{}, err
	}

	f := bson.M{"_id": objectId}

	return s.FindOne(f)
}

func (s *WorkspaceRepositoryMongo) FindOne(f bson.M) (entity.Workspace, error) {
	w := entity.Workspace{}

	res := collection.FindOne(context.TODO(), f)

	if err := res.Decode(&w); err != nil {
		return entity.Workspace{}, err
	}

	return w, nil
}

func (s *WorkspaceRepositoryMongo) FindBy(f map[string]interface{}) ([]entity.Workspace, error) {
	res, err := collection.Find(context.TODO(), f)

	if err != nil {
		return make([]entity.Workspace, 0), err
	}

	c := make([]entity.Workspace, 5)

	err = res.All(context.TODO(), &c)

	if err != nil {
		return make([]entity.Workspace, 0), err
	}

	return c, nil
}
