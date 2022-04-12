package application

import "github.com/rubumo/core/internal/entity"

type ApplicationRepository interface {
	Create(entity.Application) (entity.Application, error)
	FindById(string) (entity.Application, error)
	FindByDomain(string, string) (entity.Application, error)
	FindBy(map[string]interface{}) ([]entity.Application, error)
}
