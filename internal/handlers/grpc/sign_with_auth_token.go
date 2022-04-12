package grpc

import (
	"context"

	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"github.com/rubumo/core/internal/utils"
)

func (*Router) SignWithAuthToken(ctx context.Context, req *pb.SignWithAuthTokenRequest) (*pb.SignInResponse, error) {
	_, claims, err := utils.VerifyJWT(req.GetToken())

	if err != nil {
		return nil, err
	}

	service := service.CreateUserService()
	user, err := service.FindById(claims.Issuer)

	if err != nil {
		return nil, err
	}

	err = utils.SaveJWT(ctx, user.ID.Hex())

	if err != nil {
		return nil, err
	}

	return &pb.SignInResponse{}, nil
}
