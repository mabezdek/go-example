package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	GivenName         string             `bson:"givenName"`
	FamilyName        string             `bson:"familyName"`
	Email             string             `bson:"email"`
	Picture           string             `bson:"picture"`
	PreferredLanguage string             `bson:"preferredLanguage"`
}
