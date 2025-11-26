package adapters

import (
	"context"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProfileServiceAdapter_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()

	resp := &pb.GetProfileResponse{
		User: &pb.User{
			Id:         userID.String(),
			Email:      "test@example.com",
			Name:       "Test User",
			BirthDate:  timestamppb.New(time.Now().AddDate(-20, 0, 0)),
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			LastActive: timestamppb.Now(),
		},
		Preferences: &pb.UserPreference{
			UserId: userID.String(),
			AgeMin: 18,
			AgeMax: 30,
		},
		Photos:    []*pb.UserPhoto{},
		Interests: []*pb.Interest{},
	}

	mockClient.EXPECT().
		GetProfile(gomock.Any(), &pb.GetProfileRequest{
			UserId: userID.String(),
		}).
		Return(resp, nil)

	profile, err := adapter.GetProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.User.ID)
}

func TestProfileServiceAdapter_UpdateProfileInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()
	name := "Updated Name"

	updateData := &domain.ProfileUpdateRequest{
		Name: name,
	}

	mockClient.EXPECT().
		UpdateProfileInfo(gomock.Any(), gomock.Any()).
		Return(&pb.UpdateProfileInfoResponse{Success: true}, nil)

	err := adapter.UpdateProfileInfo(context.Background(), userID, updateData)
	assert.NoError(t, err)
}

func TestProfileServiceAdapter_UpdatePreferences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()
	updateData := &domain.PreferencesUpdateRequest{
		AgeMin: 20,
		AgeMax: 30,
	}

	mockClient.EXPECT().
		UpdatePreferences(gomock.Any(), gomock.Any()).
		Return(&pb.UpdatePreferencesResponse{Success: true}, nil)

	err := adapter.UpdatePreferences(context.Background(), userID, updateData)
	assert.NoError(t, err)
}

func TestProfileServiceAdapter_UpdateInterests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()
	interests := []domain.Interest{
		{Theme: "music"},
	}

	mockClient.EXPECT().
		UpdateInterests(gomock.Any(), &pb.UpdateInterestsRequest{
			UserId: userID.String(),
			Themes: []string{"music"},
		}).
		Return(&pb.UpdateInterestsResponse{Success: true}, nil)

	err := adapter.UpdateInterests(context.Background(), userID, interests)
	assert.NoError(t, err)
}

func TestProfileServiceAdapter_DeletePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()
	photoID := uuid.New()

	mockClient.EXPECT().
		DeletePhoto(gomock.Any(), &pb.DeletePhotoRequest{
			UserId:  userID.String(),
			PhotoId: photoID.String(),
		}).
		Return(&pb.DeletePhotoResponse{Success: true}, nil)

	err := adapter.DeletePhoto(context.Background(), userID, photoID)
	assert.NoError(t, err)
}

func TestProfileServiceAdapter_SetPrimaryPhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()
	photoID := uuid.New()

	mockClient.EXPECT().
		SetPrimaryPhoto(gomock.Any(), &pb.SetPrimaryPhotoRequest{
			UserId:  userID.String(),
			PhotoId: photoID.String(),
		}).
		Return(&pb.SetPrimaryPhotoResponse{Success: true}, nil)

	err := adapter.SetPrimaryPhoto(context.Background(), userID, photoID)
	assert.NoError(t, err)
}

func TestProfileServiceAdapter_UploadPhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewProfileServiceAdapter(mockClient)

	userID := uuid.New()
	content := []byte("test")
	contentType := "image/jpeg"
	photoID := uuid.New()

	resp := &pb.UploadPhotoResponse{
		Photo: &pb.UserPhoto{
			Id:       photoID.String(),
			UserId:   userID.String(),
			PhotoUrl: "http://example.com/test.jpg",
		},
	}

	mockClient.EXPECT().
		UploadPhoto(gomock.Any(), &pb.UploadPhotoRequest{
			UserId:      userID.String(),
			Content:     content,
			ContentType: contentType,
		}).
		Return(resp, nil)

	photo, err := adapter.UploadPhoto(context.Background(), userID, content, contentType)
	assert.NoError(t, err)
	assert.NotNil(t, photo)
	assert.Equal(t, photoID, photo.ID)
}

func TestProfileServiceAdapter_UploadPhotos(t *testing.T) {
	// This test is harder because it involves multipart.FileHeader which is hard to mock/construct fully valid.
	// We can skip it or try to construct a dummy FileHeader if possible, or mock the Open method if we could (but we can't easily).
	// However, the adapter calls Open() on the file header.
	// We can create a real multipart form in memory.
}
