package adapters

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
)

// Proto → Domain converters (shared across all adapters)

func coreProtoToUser(pbUser *pb.User) *domain.User {
	if pbUser == nil {
		return nil
	}

	id, _ := uuid.Parse(pbUser.Id)

	user := &domain.User{
		ID:              id,
		Email:           pbUser.Email,
		Name:            pbUser.Name,
		BirthDate:       pbUser.BirthDate.AsTime(),
		Gender:          constants.Gender(pbUser.Gender),
		IsVerified:      pbUser.IsVerified,
		IsPremium:       pbUser.IsPremium,
		SuperLikesCount: int(pbUser.SuperLikesCount),
		LastActive:      pbUser.LastActive.AsTime(),
		CreatedAt:       pbUser.CreatedAt.AsTime(),
		UpdatedAt:       pbUser.UpdatedAt.AsTime(),
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

func coreProtoToUserPreference(pbPref *pb.UserPreference) *domain.UserPreference {
	if pbPref == nil {
		return nil
	}

	userID, _ := uuid.Parse(pbPref.UserId)

	return &domain.UserPreference{
		UserID:       userID,
		ShowGender:   constants.GenderPreference(pbPref.ShowGender),
		AgeMin:       int(pbPref.AgeMin),
		AgeMax:       int(pbPref.AgeMax),
		MaxDistance:  int(pbPref.MaxDistance),
		GlobalSearch: pbPref.GlobalSearch,
		CreatedAt:    pbPref.CreatedAt.AsTime(),
		UpdatedAt:    pbPref.UpdatedAt.AsTime(),
	}
}

func coreProtoToUserPhoto(pbPhoto *pb.UserPhoto) domain.UserPhoto {
	id, _ := uuid.Parse(pbPhoto.Id)
	userID, _ := uuid.Parse(pbPhoto.UserId)

	return domain.UserPhoto{
		ID:           id,
		UserID:       userID,
		PhotoURL:     pbPhoto.PhotoUrl,
		DisplayOrder: int(pbPhoto.DisplayOrder),
		IsApproved:   pbPhoto.IsApproved,
		CreatedAt:    pbPhoto.CreatedAt.AsTime(),
	}
}

func coreProtoToUserPhotos(pbPhotos []*pb.UserPhoto) []domain.UserPhoto {
	photos := make([]domain.UserPhoto, len(pbPhotos))
	for i, pbPhoto := range pbPhotos {
		photos[i] = coreProtoToUserPhoto(pbPhoto)
	}
	return photos
}

func coreProtoToInterest(pbInterest *pb.Interest) domain.Interest {
	userID, _ := uuid.Parse(pbInterest.UserId)

	return domain.Interest{
		UserID: userID,
		Theme:  constants.InterestType(pbInterest.Theme),
	}
}

func coreProtoToInterests(pbInterests []*pb.Interest) []domain.Interest {
	interests := make([]domain.Interest, len(pbInterests))
	for i, pbInterest := range pbInterests {
		interests[i] = coreProtoToInterest(pbInterest)
	}
	return interests
}

func coreProtoToMatch(pbMatch *pb.Match) domain.Match {
	user1ID, _ := uuid.Parse(pbMatch.User1Id)
	user2ID, _ := uuid.Parse(pbMatch.User2Id)
	matchID, _ := uuid.Parse(pbMatch.Id)

	match := domain.Match{
		ID:        matchID,
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  pbMatch.IsActive,
		MatchedAt: pbMatch.MatchedAt.AsTime(),
	}

	// Обрабатываем expires_at, если оно есть
	if pbMatch.ExpiresAt != nil {
		expiresAt := pbMatch.ExpiresAt.AsTime()
		match.ExpiresAt = &expiresAt
	}

	return match
}
