package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/rubumo/core/internal/entity"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/service"
	"github.com/rubumo/core/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type userInfo struct {
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	Locale     string `json:"locale"`
}

func (*Router) SignWithGoogle(ctx context.Context, req *pb.SignWithGoogleRequest) (*pb.SignInResponse, error) {
	ga := &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_SSO_CALLBACK"),
		ClientID:     os.Getenv("GOOGLE_SSO_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SSO_SECRET"),
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	token, err := ga.Exchange(
		context.Background(),
		req.GetRefreshToken(),
	)

	if err != nil {
		return nil, err
	}

	userInfo, err := getUserInfo(token.AccessToken)

	if err != nil {
		return nil, err
	}

	service := service.CreateUserService()
	user, err := service.FindByEmail(userInfo.Email)

	if err != nil {
		user, err = service.Create(entity.User{
			Name:              userInfo.Name,
			GivenName:         userInfo.GivenName,
			FamilyName:        userInfo.FamilyName,
			Email:             userInfo.Email,
			Picture:           userInfo.Picture,
			PreferredLanguage: userInfo.Locale,
		})

		if err != nil {
			return nil, err
		}
	}

	err = utils.SaveJWT(ctx, user.ID.Hex())

	if err != nil {
		return nil, err
	}

	return &pb.SignInResponse{}, nil
}

func getUserInfo(token string) (userInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + url.QueryEscape(token))

	if err != nil {
		return userInfo{}, err
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return userInfo{}, err
	}

	var userInfo userInfo

	json.Unmarshal(response, &userInfo)

	if userInfo.Email == "" {
		return userInfo, errors.New("Failed to associate your account with an e-mail.")
	}

	return userInfo, nil
}
