package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"go.uber.org/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMatchServiceAdapter_GetUserMatches_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewMatchServiceAdapter(mockClient)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Prepare mock protobuf user
	pbUser := &pb.User{
		Id:              uuid.New().String(),
		Email:           "test@example.com",
		Name:            "Jane Doe",
		BirthDate:       timestamppb.New(time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC)),
		Gender:          "female",
		Bio:             stringPtr("Test bio"),
		City:            stringPtr("Moscow"),
		Artist:          stringPtr("Artist"),
		Quote:           stringPtr("Quote"),
		IsVerified:      true,
		IsPremium:       false,
		SuperLikesCount: 5,
		LastActive:      timestamppb.New(time.Now()),
		CreatedAt:       timestamppb.New(time.Now().AddDate(0, -1, 0)),
		UpdatedAt:       timestamppb.New(time.Now()),
	}

	// Prepare mock protobuf match
	pbMatch := &pb.Match{
		Id:        uuid.New().String(),
		User1Id:   userID.String(),
		User2Id:   pbUser.Id,
		IsActive:  true,
		MatchedAt: timestamppb.New(time.Now().AddDate(0, -1, 0)),
	}

	// Prepare mock match response
	pbMatchResponse := &pb.MatchResponse{
		Match:       pbMatch,
		User:        pbUser,
		Photos:      []string{"photo1.jpg", "photo2.jpg"},
		Age:         29,
		Description: "Match description",
		PhotosCount: 2,
	}

	expectedResponse := &pb.GetUserMatchesResponse{
		Matches: []*pb.MatchResponse{pbMatchResponse},
		Total:    1,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	// Setup mock expectation
	mockClient.EXPECT().
		GetUserMatches(gomock.Any(), &pb.GetUserMatchesRequest{
			UserId: userID.String(),
			Limit:  int32(limit),
			Offset: int32(offset),
		}).
		Return(expectedResponse, nil)

	// Execute test
	result, err := adapter.GetUserMatches(context.Background(), userID, limit, offset)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Matches, 1)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)

	match := result.Matches[0]
	assert.Equal(t, pbMatch.Id, match.Match.ID.String())
	assert.Equal(t, userID.String(), match.Match.User1ID.String())
	assert.Equal(t, pbMatch.User2Id, match.Match.User2ID.String())
	assert.True(t, match.Match.IsActive)
	assert.NotNil(t, match.Match.MatchedAt)

	// Verify user data
	assert.Equal(t, pbUser.Id, match.User.ID.String())
	assert.Equal(t, pbUser.Email, match.User.Email)
	assert.Equal(t, pbUser.Name, match.User.Name)
	assert.Equal(t, constants.Gender("female"), match.User.Gender)
	assert.Equal(t, pbUser.IsVerified, match.User.IsVerified)
	assert.Equal(t, pbUser.IsPremium, match.User.IsPremium)
	assert.Equal(t, int(pbUser.SuperLikesCount), match.User.SuperLikesCount)

	// Verify additional fields
	assert.Equal(t, pbMatchResponse.Photos, match.Photos)
	assert.Equal(t, int(pbMatchResponse.Age), match.Age)
	assert.Equal(t, pbMatchResponse.Description, match.Description)
	assert.Equal(t, int(pbMatchResponse.PhotosCount), match.PhotosCount)
}

func TestMatchServiceAdapter_GetUserMatches_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewMatchServiceAdapter(mockClient)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Setup mock expectation for error
	mockClient.EXPECT().
		GetUserMatches(gomock.Any(), &pb.GetUserMatchesRequest{
			UserId: userID.String(),
			Limit:  int32(limit),
			Offset: int32(offset),
		}).
		Return(nil, assert.AnError)

	// Execute test
	result, err := adapter.GetUserMatches(context.Background(), userID, limit, offset)

	// Verify results
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMatchServiceAdapter_GetUserMatches_EmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewMatchServiceAdapter(mockClient)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Prepare empty response
	expectedResponse := &pb.GetUserMatchesResponse{
		Matches: []*pb.MatchResponse{},
		Total:    0,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	// Setup mock expectation
	mockClient.EXPECT().
		GetUserMatches(gomock.Any(), &pb.GetUserMatchesRequest{
			UserId: userID.String(),
			Limit:  int32(limit),
			Offset: int32(offset),
		}).
		Return(expectedResponse, nil)

	// Execute test
	result, err := adapter.GetUserMatches(context.Background(), userID, limit, offset)

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Matches, 0)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
}

func TestMatchServiceAdapter_Unmatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewMatchServiceAdapter(mockClient)

	userID := uuid.New()
	targetUserID := uuid.New()

	expectedResponse := &pb.UnmatchResponse{
		Success: true,
	}

	// Setup mock expectation
	mockClient.EXPECT().
		Unmatch(gomock.Any(), &pb.UnmatchRequest{
			UserId:       userID.String(),
			TargetUserId: targetUserID.String(),
		}).
		Return(expectedResponse, nil)

	// Execute test
	err := adapter.Unmatch(context.Background(), userID, targetUserID)

	// Verify results
	assert.NoError(t, err)
}

func TestMatchServiceAdapter_Unmatch_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewMatchServiceAdapter(mockClient)

	userID := uuid.New()
	targetUserID := uuid.New()

	// Setup mock expectation for error
	mockClient.EXPECT().
		Unmatch(gomock.Any(), &pb.UnmatchRequest{
			UserId:       userID.String(),
			TargetUserId: targetUserID.String(),
		}).
		Return(nil, assert.AnError)

	// Execute test
	err := adapter.Unmatch(context.Background(), userID, targetUserID)

	// Verify results
	assert.Error(t, err)
}
