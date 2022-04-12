package service

import (
	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/repository/user"
)

type UserService struct {
	userRepository user.UserRepository
}

func CreateUserService() *UserService {
	return &UserService{
		userRepository: &user.UserRepositoryMongo{},
	}
}

func (s *UserService) Create(u entity.User) (entity.User, error) {
	return s.userRepository.Create(u)
}

func (s *UserService) FindByEmail(e string) (entity.User, error) {
	return s.userRepository.FindByEmail(e)
}

func (s *UserService) FindById(id string) (entity.User, error) {
	return s.userRepository.FindById(id)
}
