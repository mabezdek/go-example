package user

import "github.com/rubumo/core/internal/entity"

type UserRepository interface {
	Create(entity.User) (entity.User, error)
	FindById(string) (entity.User, error)
	FindByEmail(string) (entity.User, error)
}
