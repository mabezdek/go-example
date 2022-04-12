package grpc

import (
	"context"
	"os"
	"strings"

	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"github.com/rubumo/core/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
)

func (*Router) Me(ctx context.Context, req *pb.MeRequest) (*pb.MeResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	response := &pb.MeResponse{}

	if ok {
		userId := strings.Join(md.Get("user"), "")

		userService := service.CreateUserService()
		user, err := userService.FindById(userId)

		if err != nil {
			return nil, err
		}

		origin := md.Get("origin")[0]

		if origin == os.Getenv("APP_BASE_URI") {
			objectId, err := primitive.ObjectIDFromHex(userId)

			if err != nil {
				return nil, err
			}

			jwt, err := utils.GenerateJWT(user.ID.Hex())

			if err != nil {
				return nil, err
			}

			applicationService := service.CreateApplicationService()

			applications, err := applicationService.FindBy(bson.M{
				"access.userId": objectId,
			})

			if err != nil || len(applications) > 0 {
				domain, err := applications[0].GetDomain()

				if err != nil {
					return nil, err
				}

				response.Redirect = domain + "/login?auth_token=" + jwt
				return response, nil
			}

			workspaceService := service.CreateWorkspaceService()

			workspaces, err := workspaceService.FindBy(bson.M{
				"access.userId": objectId,
			})

			if err != nil || len(workspaces) > 0 {
				domain, err := workspaces[0].GetDomain()

				if err != nil {
					return nil, err
				}

				response.Redirect = domain + "/login?auth_token=" + jwt
				return response, nil
			}

			response.Redirect = "/create-workspace"
			return response, nil
		}
	}

	return response, nil
}
