package grpc

import (
	"context"
	"net/url"

	"github.com/gosimple/slug"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"google.golang.org/grpc/metadata"
)

func (*Router) ValidateApplicationDomain(ctx context.Context, req *pb.ValidateApplicationDomainRequest) (*pb.ValidateApplicationDomainResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	response := &pb.ValidateApplicationDomainResponse{
		Valid: true,
	}

	if ok {
		u, err := url.Parse(md.Get("origin")[0])

		if err != nil {
			return nil, err
		}

		service := service.CreateApplicationService()

		_, err = service.FindByDomain(u.Hostname(), slug.Make(req.GetDomain()))

		if err != nil {
			return response, nil
		}
	}

	response.Valid = false

	return response, nil
}
