package workspace

import (
	"github.com/rubumo/core/internal/aggregate"
	"github.com/rubumo/core/internal/entity"
)

type WorkspaceRepository interface {
	Create(entity.Workspace) (entity.Workspace, error)
	FindById(string) (entity.Workspace, error)
	FindByDomain(string) (entity.Workspace, error)
	FindBy(map[string]interface{}) ([]entity.Workspace, error)
	FindByDomainWithApplications(string) (aggregate.WorkspaceWithApplications, error)
}
