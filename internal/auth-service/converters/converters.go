package converters

import (
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserToProto converts domain.User to protobuf User
func UserToProto(user *domain.User) *pb.User {
	if user == nil {
		return nil
	}

	pbUser := &pb.User{
		Id:         user.ID.String(),
		Email:      user.Email,
		Name:       user.Name,
		BirthDate:  timestamppb.New(user.BirthDate),
		Gender:     string(user.Gender),
		IsVerified: user.IsVerified,
		LastActive: timestamppb.New(user.LastActive),
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}

	// Handle optional fields
	if user.Phone != nil {
		pbUser.Phone = user.Phone
	}
	if user.Bio != nil {
		pbUser.Bio = user.Bio
	}
	if user.City != nil {
		pbUser.City = user.City
	}
	if user.Artist != nil {
		pbUser.Artist = user.Artist
	}
	if user.Quote != nil {
		pbUser.Quote = user.Quote
	}

	return pbUser
}

// SessionToProto converts domain.Session to protobuf Session
func SessionToProto(session *domain.Session) *pb.Session {
	if session == nil {
		return nil
	}

	return &pb.Session{
		Token:     session.Token,
		UserEmail: session.UserEmail,
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
