package transform

import (
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	pb "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
)

func FromUserToPbUser(user entities.User) pb.User {
	return pb.User{
		Id:       user.ID,
		Name:     user.Name,
		PwdHash:  user.Password,
		Age:      uint32(user.Age),
		ParentId: user.ParentID,
	}
}

func FromPbUserToUser(user *pb.User) entities.User {
	return entities.User{
		ID:       user.Id,
		Name:     user.Name,
		Password: user.PwdHash,
		Age:      uint8(user.Age),
		ParentID: user.ParentId,
	}
}
