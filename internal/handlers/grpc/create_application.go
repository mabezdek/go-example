package grpc

import (
	"context"
	"net/url"
	"strings"

	"github.com/gosimple/slug"
	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
)

func (*Router) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.Application, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		u, err := url.Parse(md.Get("origin")[0])

		if err != nil {
			return nil, err
		}

		workspace, err := service.CreateWorkspaceService().FindByDomain(u.Hostname())

		if err != nil {
			return nil, err
		}

		service := service.CreateApplicationService()

		_, err = service.FindByDomain(u.Hostname(), slug.Make(req.GetDomain()))

		if err == nil {
			return nil, err
		}

		userId := strings.Join(md.Get("user"), "")

		objectId, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			return nil, err
		}

		application, err := service.Create(entity.Application{
			Workspace: workspace,
			Environments: []entity.ApplicationEnvironment{
				{
					Code:          "default",
					IsDevelopment: true,
					Domain:        req.GetDomain(),
					Name:          "Default",
					Localizations: []entity.ApplicationLocalization{
						{
							Name: req.GetName(),
							Lang: req.GetPrimaryLanguage(),
						},
					},
				},
			},
			Access: []entity.ApplicationAccess{
				{
					UserID: objectId,
					Role:   "owner",
				},
			},
		})

		if err != nil {
			return nil, err
		}

		return application.ToPb(), nil
	}

	return &pb.Application{}, nil
}
