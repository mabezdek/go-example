package service

import (
	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/repository/application"
)

type ApplicationService struct {
	applicationRepository application.ApplicationRepository
}

func CreateApplicationService() *ApplicationService {
	return &ApplicationService{
		applicationRepository: &application.ApplicationRepositoryMongo{},
	}
}

func (s *ApplicationService) Create(w entity.Application) (entity.Application, error) {
	return s.applicationRepository.Create(w)
}

func (s *ApplicationService) FindByDomain(w string, d string) (entity.Application, error) {
	return s.applicationRepository.FindByDomain(w, d)
}

func (s *ApplicationService) FindById(id string) (entity.Application, error) {
	return s.applicationRepository.FindById(id)
}

func (s *ApplicationService) FindBy(f map[string]interface{}) ([]entity.Application, error) {
	return s.applicationRepository.FindBy(f)
}
