package grpc

import (
	"context"
	"net/url"

	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"github.com/rubumo/core/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (*Router) InitWorkspace(ctx context.Context, req *pb.InitWorkspaceRequest) (*pb.Workspace, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		u, err := url.Parse(md.Get("origin")[0])

		if err != nil {
			return nil, err
		}

		service := service.CreateWorkspaceService()

		workspace, err := service.FindByDomainWithApplications(u.Hostname())

		if err != nil {
			return nil, err
		}

		userId := md.Get("user")[0]
		jwt, err := utils.GenerateJWT(userId)

		if err != nil {
			return nil, err
		}

		grpc.SetTrailer(ctx, metadata.Pairs("auth_token", jwt))

		return workspace.ToPb(), nil
	}

	return &pb.Workspace{}, nil
}
