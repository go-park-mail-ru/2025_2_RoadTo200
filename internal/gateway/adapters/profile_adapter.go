package adapters

import (
	"context"
	"mime/multipart"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProfileServiceAdapter struct {
	client pb.CoreServiceClient
}

func NewProfileServiceAdapter(client pb.CoreServiceClient) service.ProfileService {
	return &ProfileServiceAdapter{
		client: client,
	}
}

func (a *ProfileServiceAdapter) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.ProfileResponse, error) {
	resp, err := a.client.GetProfile(ctx, &pb.GetProfileRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &domain.ProfileResponse{
		User:        coreProtoToUser(resp.User),
		Preferences: coreProtoToUserPreference(resp.Preferences),
		Photos:      coreProtoToUserPhotos(resp.Photos),
	}, nil
}

func (a *ProfileServiceAdapter) UpdateProfileInfo(ctx context.Context, userID uuid.UUID, updateData *domain.ProfileUpdateRequest) error {
	req := &pb.UpdateProfileInfoRequest{
		UserId: userID.String(),
	}

	if updateData.Name != "" {
		req.Name = &updateData.Name
	}
	if updateData.Phone != nil {
		req.Phone = updateData.Phone
	}
	if updateData.BirthDate != nil {
		req.BirthDate = timestamppb.New(*updateData.BirthDate)
	}
	if updateData.Gender != "" {
		gender := string(updateData.Gender)
		req.Gender = &gender
	}
	if updateData.Bio != nil {
		req.Bio = updateData.Bio
	}
	if updateData.Artist != nil {
		req.Artist = updateData.Artist
	}
	if updateData.Quote != nil {
		req.Quote = updateData.Quote
	}
	if updateData.Latitude != nil {
		req.Latitude = updateData.Latitude
	}
	if updateData.Longitude != nil {
		req.Longitude = updateData.Longitude
	}

	_, err := a.client.UpdateProfileInfo(ctx, req)
	return err
}

func (a *ProfileServiceAdapter) UpdatePreferences(ctx context.Context, userID uuid.UUID, updateData *domain.PreferencesUpdateRequest) error {
	req := &pb.UpdatePreferencesRequest{
		UserId: userID.String(),
	}

	if updateData.ShowGender != "" {
		showGender := string(updateData.ShowGender)
		req.ShowGender = &showGender
	}
	if updateData.AgeMin != 0 {
		ageMin := int32(updateData.AgeMin)
		req.AgeMin = &ageMin
	}
	if updateData.AgeMax != 0 {
		ageMax := int32(updateData.AgeMax)
		req.AgeMax = &ageMax
	}
	if updateData.MaxDistance != 0 {
		maxDist := int32(updateData.MaxDistance)
		req.MaxDistance = &maxDist
	}
	req.GlobalSearch = &updateData.GlobalSearch

	_, err := a.client.UpdatePreferences(ctx, req)
	return err
}

func (a *ProfileServiceAdapter) UpdateInterests(ctx context.Context, userID uuid.UUID, interests []domain.Interest) error {
	themes := make([]string, len(interests))
	for i, interest := range interests {
		themes[i] = string(interest.Theme)
	}

	_, err := a.client.UpdateInterests(ctx, &pb.UpdateInterestsRequest{
		UserId: userID.String(),
		Themes: themes,
	})
	return err
}

func (a *ProfileServiceAdapter) UploadPhotos(ctx context.Context, userID uuid.UUID, photos []*multipart.FileHeader) ([]domain.UserPhoto, error) {
	// Photo upload stays in Gateway (cannot send multipart via gRPC efficiently)
	// This method should not be called through adapter
	return nil, nil
}

func (a *ProfileServiceAdapter) DeletePhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error {
	_, err := a.client.DeletePhoto(ctx, &pb.DeletePhotoRequest{
		UserId:  userID.String(),
		PhotoId: photoID.String(),
	})
	return err
}

func (a *ProfileServiceAdapter) SetPrimaryPhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error {
	_, err := a.client.SetPrimaryPhoto(ctx, &pb.SetPrimaryPhotoRequest{
		UserId:  userID.String(),
		PhotoId: photoID.String(),
	})
	return err
}

func (a *ProfileServiceAdapter) ReorderPhotos(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error {
	photoIDStrs := make([]string, len(photoIDs))
	for i, id := range photoIDs {
		photoIDStrs[i] = id.String()
	}

	_, err := a.client.ReorderPhotos(ctx, &pb.ReorderPhotosRequest{
		UserId:   userID.String(),
		PhotoIds: photoIDStrs,
	})
	return err
}

func (a *ProfileServiceAdapter) ValidateProfileUpdate(ctx context.Context, updateData *domain.ProfileUpdateRequest) error {
	// Validation stays in Gateway
	return nil
}

func (a *ProfileServiceAdapter) ValidatePreferencesUpdate(ctx context.Context, updateData *domain.PreferencesUpdateRequest) error {
	// Validation stays in Gateway
	return nil
}
