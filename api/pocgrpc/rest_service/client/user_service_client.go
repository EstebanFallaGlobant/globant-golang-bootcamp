package client

import (
	"context"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/transform"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
)

type userServiceClient struct {
	gRPCClient pb.UserDetailServiceClient
}

func (usc userServiceClient) GetUser(ctx context.Context, ID int64) (entities.User, error) {
	authToken, ok := ctx.Value(authTokenCtxKey).(string)

	if !ok {
		return entities.User{}, svcerr.NewInvalidArgumentError(paramAuthTokenName, ruleAuthTokenInvalidType)
	}

	request := pb.GetUserRequest{
		Id:        ID,
		AuthToken: authToken,
	}
	response, err := usc.gRPCClient.GetUser(ctx, &request, nil)

	if err != nil {
		return entities.User{}, err
	}

	return transform.FromPbUserToUser(response.User), nil
}
