package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workspace struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name"`
	Domain string             `bson:"domain"`
	Access []WorkspaceAccess  `bson:"access"`
}

type WorkspaceAccess struct {
	UserID primitive.ObjectID `bson:"userId"`
	Role   string             `bson:"role"`
}

func (w *Workspace) GetDomain() (string, error) {
	return "http://" + w.Domain + ".rubumo.io:8080", nil
}
