package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/rubumo/core/internal/database"
	"github.com/rubumo/core/internal/pb"
	"github.com/rubumo/core/internal/temporal"
	"github.com/rubumo/core/internal/utils"

	handlers "github.com/rubumo/core/internal/handlers/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if strings.Index(info.FullMethod, "/pb.AuthService/SignWith") == 0 {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		cookies := strings.Split(md.Get("cookie")[0], "; ")

		m := make(map[string]string)

		for _, c := range cookies {
			parts := strings.Split(c, "=")
			m[parts[0]] = parts[1]
		}

		if val, ok := m["jwt"]; ok {
			_, claims, err := utils.VerifyJWT(val)

			if err == nil {
				md.Set("jwt", val)
				md.Set("user", claims.Issuer)

				ctx = metadata.NewIncomingContext(ctx, md)

				return handler(ctx, req)
			}
		}
	}

	return nil, errors.New("Forbidden.")
}

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	database.Connect()
	temporal.Connect()

	lis, err := net.Listen("tcp", os.Getenv("TCP_SERVER"))

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		withServerUnaryInterceptor(),
	)

	grpcRouter := &handlers.Router{}

	pb.RegisterAuthServiceServer(s, grpcRouter)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	if err := database.Database.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	temporal.Client.Close()

	s.Stop()
}
