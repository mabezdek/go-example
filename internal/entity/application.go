package entity

import (
	"errors"

	"github.com/rubumo/core/internal/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Application struct {
	ID           primitive.ObjectID       `bson:"_id,omitempty"`
	Workspace    Workspace                `bson:"workspace"`
	Environments []ApplicationEnvironment `bson:"environments"`
	Access       []ApplicationAccess      `bson:"access"`
}

type ApplicationAccess struct {
	UserID primitive.ObjectID `bson:"userId"`
	Role   string             `bson:"role"`
}

type ApplicationEnvironment struct {
	Code          string                    `bson:"code"`
	Name          string                    `bson:"name"`
	Domain        string                    `bson:"domain"`
	Localizations []ApplicationLocalization `bson:"localizations"`
	IsDevelopment bool                      `bson:"isDevelopment"`
}

type ApplicationLocalization struct {
	Name string `bson:"name"`
	Lang string `bson:"lang"`
}

func (a *Application) GetDevelopmentEnv() (ApplicationEnvironment, error) {
	for _, e := range a.Environments {
		if e.IsDevelopment {
			return e, nil
		}
	}

	return ApplicationEnvironment{}, errors.New("Not environment found")
}

func (a *Application) GetDomain() (string, error) {
	env, err := a.GetDevelopmentEnv()

	if err != nil {
		return "", err
	}

	return "http://" + env.Domain + "." + a.Workspace.Domain + ".rubumo.io:8080", nil
}

func (a *Application) ToPb() *pb.Application {
	environments := []*pb.Environment{}

	for _, e := range a.Environments {
		localizations := []*pb.Localization{}

		for _, l := range e.Localizations {
			localizations = append(localizations, &pb.Localization{
				Name: e.Name,
				Lang: l.Lang,
			})
		}

		environments = append(environments, &pb.Environment{
			Code:          e.Code,
			Name:          e.Name,
			Domain:        e.Domain,
			IsDevelopment: e.IsDevelopment,
			Localizations: localizations,
		})
	}

	return &pb.Application{
		Id:           a.ID.Hex(),
		Environments: environments,
	}
}
