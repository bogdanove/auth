package converter

import (
	"github.com/bogdanove/auth/internal/model"
	pb "github.com/bogdanove/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToPBUserFromService - конвертер информации о пользвателе в структуру протобаф
func ToPBUserFromService(user *model.User) *pb.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      pb.Role(pb.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserInfoFromPB - конвертер информации о пользователе из протобаф в сервисный слой
func ToUserInfoFromPB(info *pb.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role.String(),
	}
}

// ToUpdateUserInfoFromPB - конвертер обновления информации о пользователе из протобаф в сервисный слой
func ToUpdateUserInfoFromPB(info *pb.UpdateRequest) *model.UpdateUserInfo {

	return &model.UpdateUserInfo{
		ID:   info.Id,
		Name: info.UpdateUserInfo.Name.GetValue(),
		Role: info.UpdateUserInfo.Role.String(),
	}
}
