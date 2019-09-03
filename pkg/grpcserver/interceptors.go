package grpcserver

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"hermes/paygearauth"
	"time"
)

func UnaryAuthJWTInterceptor(ctx context.Context) (context.Context, error) {
	logrus.Info("In UnaryJWTInterceptor")
	ident, err := jwtCheck(ctx)
	if err != nil {
		logrus.Errorf("Authentication failed : %v", err)
		return ctx, errors.Wrap(err, "error in auth")
	}

	grpc_ctxtags.Extract(ctx).Set("identity", ident)
	newCtx := context.WithValue(ctx, "identity", ident)
	newCtx, _ = context.WithTimeout(newCtx, time.Hour)
	//newCtx, _ = context.WithTimeout(newCtx, time.Hour * 1)
	return newCtx, nil

}
func jwtCheck(ctx context.Context) (*paygearauth.Identity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("error in getting context meta data")
	}
	//check jwt
	bearer, exists := md["paygearauth"]
	if !exists {
		return nil, errors.New("no bearer token found")
	}
	_ = bearer
	// checkBearer(bearer)
	fmt.Println(bearer[0])
	ident, err := paygearauth.GetAuthentication(bearer[0], "")
	if err != nil {
		return nil, errors.Wrap(err, "error in verifying jwt")
	}
	return ident, nil
}