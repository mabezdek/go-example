package aggregate

import (
	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/pb"
)

type WorkspaceWithApplications struct {
	entity.Workspace `bson:",inline"`
	Applications     []entity.Application `bson:"applications"`
}

func (w *WorkspaceWithApplications) ToPb() *pb.Workspace {
	applications := []*pb.Application{}

	for _, a := range w.Applications {
		applications = append(applications, a.ToPb())
	}

	return &pb.Workspace{
		Id:           w.ID.Hex(),
		Domain:       w.Domain,
		Name:         w.Name,
		Applications: applications,
	}
}
