package converter

import (
	"github.com/bogdanove/auth/internal/model"
	pb "github.com/bogdanove/auth/pkg/user_v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToPBUserFromService - конвертер информации о пользвателе в структуру протобаф
func ToPBUserFromService(user *model.User) (*pb.User, error) {
	if user == nil {
		return nil, errors.New("user for convert is nil")
	}

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
	}, nil
}

// ToUserInfoFromPB - конвертер информации о пользователе из протобаф в сервисный слой
func ToUserInfoFromPB(info *pb.UserInfo) (*model.UserInfo, error) {
	if info == nil {
		return nil, errors.New("info for convert is nil")
	}

	return &model.UserInfo{
		Name:            info.Name,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Role:            info.Role.String(),
	}, nil
}

// ToUpdateUserInfoFromPB - конвертер обновления информации о пользователе из протобаф в сервисный слой
func ToUpdateUserInfoFromPB(info *pb.UpdateRequest) (*model.UpdateUserInfo, error) {
	if info == nil {
		return nil, errors.New("info for convert is nil")
	}

	var user = &model.UpdateUserInfo{
		ID: info.Id,
	}

	if info.UpdateUserInfo.Name != nil {
		name := info.UpdateUserInfo.Name.GetValue()
		user.Name = &name
	}

	if info.UpdateUserInfo.Role != nil {
		role := info.UpdateUserInfo.Role.String()
		user.Role = &role
	}

	return user, nil
}
