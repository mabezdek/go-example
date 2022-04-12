package service

import (
	"github.com/rubumo/core/internal/aggregate"
	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/repository/workspace"
)

type WorkspaceService struct {
	workspaceRepository workspace.WorkspaceRepository
}

func CreateWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository: &workspace.WorkspaceRepositoryMongo{},
	}
}

func (s *WorkspaceService) Create(w entity.Workspace) (entity.Workspace, error) {
	return s.workspaceRepository.Create(w)
}

func (s *WorkspaceService) FindByDomain(d string) (entity.Workspace, error) {
	return s.workspaceRepository.FindByDomain(d)
}

func (s *WorkspaceService) FindByDomainWithApplications(d string) (aggregate.WorkspaceWithApplications, error) {
	return s.workspaceRepository.FindByDomainWithApplications(d)
}

func (s *WorkspaceService) FindById(id string) (entity.Workspace, error) {
	return s.workspaceRepository.FindById(id)
}

func (s *WorkspaceService) FindBy(f map[string]interface{}) ([]entity.Workspace, error) {
	return s.workspaceRepository.FindBy(f)
}

func (s *WorkspaceService) StructToPb(w entity.Workspace) *pb.Workspace {
	return &pb.Workspace{
		Name:   w.Name,
		Domain: w.Domain,
	}
}
