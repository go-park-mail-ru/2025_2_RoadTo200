package adapters

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"go.uber.org/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFeedServiceAdapter_GetFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewFeedServiceAdapter(mockClient)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Prepare mock protobuf response
	interest1 := &pb.Interest{
		UserId: uuid.New().String(),
		Theme:  string(constants.InterestTypeCinema),
	}
	interest2 := &pb.Interest{
		UserId: uuid.New().String(),
		Theme:  string(constants.InterestTypeWorkout),
	}

	pbUser := &pb.FeedUser{
		Id:          uuid.New().String(),
		Name:        "John Doe",
		Age:         25,
		Gender:      "male",
		Description: "Test user",
		Images:      []string{"image1.jpg", "image2.jpg"},
		PhotosCount: 2,
		Artist:      stringPtr("Artist Name"),
		Quote:       stringPtr("Test quote"),
		Interests:   []*pb.Interest{interest1, interest2},
	}

	expectedResponse := &pb.GetFeedResponse{
		Users:  []*pb.FeedUser{pbUser},
		Total:  1,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Setup mock expectation
	mockClient.EXPECT().
		GetFeed(gomock.Any(), &pb.GetFeedRequest{
			UserId: userID.String(),
			Limit:  int32(limit),
			Offset: int32(offset),
		}).
		Return(expectedResponse, nil)

	// Execute test
	result, err := adapter.GetFeed(context.Background(), userID, limit, offset)

	// Verify results
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	user := result[0]
	assert.Equal(t, pbUser.Id, user.ID)
	assert.Equal(t, pbUser.Name, user.Name)
	assert.Equal(t, int(pbUser.Age), user.Age)
	assert.Equal(t, pbUser.Gender, user.Gender)
	assert.Equal(t, pbUser.Description, user.Description)
	assert.Equal(t, pbUser.Images, user.Images)
	assert.Equal(t, int(pbUser.PhotosCount), user.PhotosCount)
	assert.Equal(t, stringPtr("Artist Name"), user.Artist)
	assert.Equal(t, stringPtr("Test quote"), user.Quote)
	assert.Len(t, user.Interests, 2)
	assert.Equal(t, constants.InterestTypeCinema, user.Interests[0].Theme)
	assert.Equal(t, constants.InterestTypeWorkout, user.Interests[1].Theme)
}

func TestFeedServiceAdapter_GetFeed_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewFeedServiceAdapter(mockClient)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Setup mock expectation for error
	mockClient.EXPECT().
		GetFeed(gomock.Any(), &pb.GetFeedRequest{
			UserId: userID.String(),
			Limit:  int32(limit),
			Offset: int32(offset),
		}).
		Return(nil, assert.AnError)

	// Execute test
	result, err := adapter.GetFeed(context.Background(), userID, limit, offset)

	// Verify results
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFeedServiceAdapter_GetFeed_EmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCoreServiceClient(ctrl)
	adapter := NewFeedServiceAdapter(mockClient)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Prepare empty response
	expectedResponse := &pb.GetFeedResponse{
		Users:  []*pb.FeedUser{},
		Total:  0,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Setup mock expectation
	mockClient.EXPECT().
		GetFeed(gomock.Any(), &pb.GetFeedRequest{
			UserId: userID.String(),
			Limit:  int32(limit),
			Offset: int32(offset),
		}).
		Return(expectedResponse, nil)

	// Execute test
	result, err := adapter.GetFeed(context.Background(), userID, limit, offset)

	// Verify results
	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

// Helper function for creating string pointers
func stringPtr(s string) *string {
	return &s
}
