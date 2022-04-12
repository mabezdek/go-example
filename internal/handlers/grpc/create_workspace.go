package grpc

import (
	"context"
	"strings"

	"github.com/gosimple/slug"
	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
)

func (*Router) CreateWorkspace(ctx context.Context, req *pb.CreateWorkspaceRequest) (*pb.Workspace, error) {
	service := service.CreateWorkspaceService()

	exists, err := service.FindByDomain(slug.Make(req.GetDomain()))

	if exists.Domain != "" {
		return nil, err
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		userId := strings.Join(md.Get("user"), "")

		objectId, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			return nil, err
		}

		workspace, err := service.Create(entity.Workspace{
			Name:   req.GetName(),
			Domain: slug.Make(req.GetDomain()),
			Access: []entity.WorkspaceAccess{
				{
					UserID: objectId,
					Role:   "owner",
				},
			},
		})

		if err != nil {
			return nil, err
		}

		return service.StructToPb(workspace), nil
	}

	return &pb.Workspace{}, nil
}
