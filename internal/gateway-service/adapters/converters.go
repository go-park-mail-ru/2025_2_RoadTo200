package adapters

import (
	auth_c "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/constants"
	auth_d "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	core_c "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	core_d "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
)

// Proto → Domain converters (shared across all adapters)

func coreProtoToUser(pbUser *pb.User) *auth_d.User {
	if pbUser == nil {
		return nil
	}

	id, _ := uuid.Parse(pbUser.Id)

	user := &auth_d.User{
		ID:         id,
		Email:      pbUser.Email,
		Name:       pbUser.Name,
		BirthDate:  pbUser.BirthDate.AsTime(),
		Gender:     auth_c.Gender(pbUser.Gender),
		IsVerified: pbUser.IsVerified,
		LastActive: pbUser.LastActive.AsTime(),
		CreatedAt:  pbUser.CreatedAt.AsTime(),
		UpdatedAt:  pbUser.UpdatedAt.AsTime(),
	}

	if pbUser.Phone != nil {
		user.Phone = pbUser.Phone
	}
	if pbUser.Bio != nil {
		user.Bio = pbUser.Bio
	}
	if pbUser.City != nil {
		user.City = pbUser.City
	}
	if pbUser.Artist != nil {
		user.Artist = pbUser.Artist
	}
	if pbUser.Quote != nil {
		user.Quote = pbUser.Quote
	}

	return user
}

func coreProtoToUserPreference(pbPref *pb.UserPreference) *core_d.UserPreference {
	if pbPref == nil {
		return nil
	}

	userID, _ := uuid.Parse(pbPref.UserId)

	return &core_d.UserPreference{
		UserID:       userID,
		ShowGender:   core_c.GenderPreference(pbPref.ShowGender),
		AgeMin:       int(pbPref.AgeMin),
		AgeMax:       int(pbPref.AgeMax),
		MaxDistance:  int(pbPref.MaxDistance),
		GlobalSearch: pbPref.GlobalSearch,
		CreatedAt:    pbPref.CreatedAt.AsTime(),
		UpdatedAt:    pbPref.UpdatedAt.AsTime(),
	}
}

func coreProtoToUserPhoto(pbPhoto *pb.UserPhoto) core_d.UserPhoto {
	id, _ := uuid.Parse(pbPhoto.Id)
	userID, _ := uuid.Parse(pbPhoto.UserId)

	return core_d.UserPhoto{
		ID:           id,
		UserID:       userID,
		PhotoURL:     pbPhoto.PhotoUrl,
		DisplayOrder: int(pbPhoto.DisplayOrder),
		IsApproved:   pbPhoto.IsApproved,
		CreatedAt:    pbPhoto.CreatedAt.AsTime(),
	}
}

func coreProtoToUserPhotos(pbPhotos []*pb.UserPhoto) []core_d.UserPhoto {
	photos := make([]core_d.UserPhoto, len(pbPhotos))
	for i, pbPhoto := range pbPhotos {
		photos[i] = coreProtoToUserPhoto(pbPhoto)
	}
	return photos
}

func coreProtoToInterest(pbInterest *pb.Interest) core_d.Interest {
	userID, _ := uuid.Parse(pbInterest.UserId)

	return core_d.Interest{
		UserID: userID,
		Theme:  core_c.InterestType(pbInterest.Theme),
	}
}

func coreProtoToInterests(pbInterests []*pb.Interest) []core_d.Interest {
	interests := make([]core_d.Interest, len(pbInterests))
	for i, pbInterest := range pbInterests {
		interests[i] = coreProtoToInterest(pbInterest)
	}
	return interests
}

func coreProtoToMatch(pbMatch *pb.Match) domain.Match {
	user1ID, _ := uuid.Parse(pbMatch.User1Id)
	user2ID, _ := uuid.Parse(pbMatch.User2Id)
	matchID, _ := uuid.Parse(pbMatch.Id)

	return domain.Match{
		ID:        matchID,
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  pbMatch.IsActive,
		MatchedAt: pbMatch.MatchedAt.AsTime(),
	}
}
