package grpc

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
)

func (*Router) ValidateWorkspaceDomain(ctx context.Context, req *pb.ValidateWorkspaceDomainRequest) (*pb.ValidateWorkspaceDomainResponse, error) {
	service := service.CreateWorkspaceService()

	response := &pb.ValidateWorkspaceDomainResponse{
		Valid: true,
	}

	_, err := service.FindByDomain(slug.Make(req.GetDomain()))

	if err != nil {
		return response, nil
	}

	response.Valid = false

	return response, nil
}
