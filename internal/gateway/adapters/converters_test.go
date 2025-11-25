package adapters

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCoreProtoToUser(t *testing.T) {
	tests := []struct {
		name     string
		pbUser   *pb.User
		expected *domain.User
	}{
		{
			name:     "nil_user",
			pbUser:   nil,
			expected: nil,
		},
		{
			name: "full_user_with_all_fields",
			pbUser: &pb.User{
				Id:         uuid.New().String(),
				Email:      "test@example.com",
				Name:       "John Doe",
				BirthDate:  timestamppb.New(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
				Gender:     "male",
				Phone:      strPtr("1234567890"),
				Bio:        strPtr("Test bio"),
				City:       strPtr("Moscow"),
				Artist:     strPtr("The Beatles"),
				Quote:      strPtr("To be or not to be"),
				IsVerified: true,
				LastActive: timestamppb.New(time.Now()),
				CreatedAt:  timestamppb.New(time.Now()),
				UpdatedAt:  timestamppb.New(time.Now()),
			},
			expected: &domain.User{
				Email:      "test@example.com",
				Name:       "John Doe",
				Gender:     constants.Gender("male"),
				IsVerified: true,
			},
		},
		{
			name: "minimal_user_without_optional_fields",
			pbUser: &pb.User{
				Id:         uuid.New().String(),
				Email:      "minimal@example.com",
				Name:       "Jane",
				BirthDate:  timestamppb.New(time.Now()),
				Gender:     "female",
				IsVerified: false,
				LastActive: timestamppb.New(time.Now()),
				CreatedAt:  timestamppb.New(time.Now()),
				UpdatedAt:  timestamppb.New(time.Now()),
			},
			expected: &domain.User{
				Email:      "minimal@example.com",
				Name:       "Jane",
				Gender:     constants.Gender("female"),
				IsVerified: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := coreProtoToUser(tt.pbUser)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.Email, result.Email)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Gender, result.Gender)
			assert.Equal(t, tt.expected.IsVerified, result.IsVerified)

			// Check optional fields
			if tt.pbUser.Phone != nil {
				assert.NotNil(t, result.Phone)
				assert.Equal(t, *tt.pbUser.Phone, *result.Phone)
			}
			if tt.pbUser.Bio != nil {
				assert.NotNil(t, result.Bio)
			}
			if tt.pbUser.City != nil {
				assert.NotNil(t, result.City)
			}
		})
	}
}

func TestCoreProtoToUserPreference(t *testing.T) {
	tests := []struct {
		name     string
		pbPref   *pb.UserPreference
		expected *domain.UserPreference
	}{
		{
			name:     "nil_preference",
			pbPref:   nil,
			expected: nil,
		},
		{
			name: "valid_preference",
			pbPref: &pb.UserPreference{
				UserId:       uuid.New().String(),
				ShowGender:   "male",
				AgeMin:       18,
				AgeMax:       30,
				MaxDistance:  50,
				GlobalSearch: false,
				CreatedAt:    timestamppb.New(time.Now()),
				UpdatedAt:    timestamppb.New(time.Now()),
			},
			expected: &domain.UserPreference{
				ShowGender:   constants.GenderPreference("male"),
				AgeMin:       18,
				AgeMax:       30,
				MaxDistance:  50,
				GlobalSearch: false,
			},
		},
		{
			name: "global_search_enabled",
			pbPref: &pb.UserPreference{
				UserId:       uuid.New().String(),
				ShowGender:   "female",
				AgeMin:       25,
				AgeMax:       35,
				MaxDistance:  100,
				GlobalSearch: true,
				CreatedAt:    timestamppb.New(time.Now()),
				UpdatedAt:    timestamppb.New(time.Now()),
			},
			expected: &domain.UserPreference{
				ShowGender:   constants.GenderPreference("female"),
				AgeMin:       25,
				AgeMax:       35,
				MaxDistance:  100,
				GlobalSearch: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := coreProtoToUserPreference(tt.pbPref)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.ShowGender, result.ShowGender)
			assert.Equal(t, tt.expected.AgeMin, result.AgeMin)
			assert.Equal(t, tt.expected.AgeMax, result.AgeMax)
			assert.Equal(t, tt.expected.MaxDistance, result.MaxDistance)
			assert.Equal(t, tt.expected.GlobalSearch, result.GlobalSearch)
		})
	}
}

func TestCoreProtoToUserPhoto(t *testing.T) {
	photoID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now()

	pbPhoto := &pb.UserPhoto{
		Id:           photoID.String(),
		UserId:       userID.String(),
		PhotoUrl:     "https://example.com/photo.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
		CreatedAt:    timestamppb.New(createdAt),
	}

	result := coreProtoToUserPhoto(pbPhoto)

	assert.Equal(t, photoID, result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "https://example.com/photo.jpg", result.PhotoURL)
	assert.Equal(t, 1, result.DisplayOrder)
	assert.True(t, result.IsApproved)
}

func TestCoreProtoToUserPhotos(t *testing.T) {
	tests := []struct {
		name     string
		pbPhotos []*pb.UserPhoto
		expected int
	}{
		{
			name:     "empty_list",
			pbPhotos: []*pb.UserPhoto{},
			expected: 0,
		},
		{
			name: "single_photo",
			pbPhotos: []*pb.UserPhoto{
				{
					Id:           uuid.New().String(),
					UserId:       uuid.New().String(),
					PhotoUrl:     "https://example.com/1.jpg",
					DisplayOrder: 0,
					IsApproved:   true,
					CreatedAt:    timestamppb.New(time.Now()),
				},
			},
			expected: 1,
		},
		{
			name: "multiple_photos",
			pbPhotos: []*pb.UserPhoto{
				{
					Id:           uuid.New().String(),
					UserId:       uuid.New().String(),
					PhotoUrl:     "https://example.com/1.jpg",
					DisplayOrder: 0,
					IsApproved:   true,
					CreatedAt:    timestamppb.New(time.Now()),
				},
				{
					Id:           uuid.New().String(),
					UserId:       uuid.New().String(),
					PhotoUrl:     "https://example.com/2.jpg",
					DisplayOrder: 1,
					IsApproved:   false,
					CreatedAt:    timestamppb.New(time.Now()),
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := coreProtoToUserPhotos(tt.pbPhotos)

			assert.Equal(t, tt.expected, len(result))

			for i, photo := range result {
				assert.Equal(t, tt.pbPhotos[i].PhotoUrl, photo.PhotoURL)
				assert.Equal(t, int(tt.pbPhotos[i].DisplayOrder), photo.DisplayOrder)
			}
		})
	}
}

func TestCoreProtoToInterest(t *testing.T) {
	userID := uuid.New()

	pbInterest := &pb.Interest{
		UserId: userID.String(),
		Theme:  "music",
	}

	result := coreProtoToInterest(pbInterest)

	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, constants.InterestType("music"), result.Theme)
}

func TestCoreProtoToInterests(t *testing.T) {
	userID := uuid.New()

	pbInterests := []*pb.Interest{
		{UserId: userID.String(), Theme: "music"},
		{UserId: userID.String(), Theme: "sports"},
		{UserId: userID.String(), Theme: "travel"},
	}

	result := coreProtoToInterests(pbInterests)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, constants.InterestType("music"), result[0].Theme)
	assert.Equal(t, constants.InterestType("sports"), result[1].Theme)
	assert.Equal(t, constants.InterestType("travel"), result[2].Theme)
}

func TestCoreProtoToMatch(t *testing.T) {
	matchID := uuid.New()
	user1ID := uuid.New()
	user2ID := uuid.New()
	matchedAt := time.Now()

	pbMatch := &pb.Match{
		Id:        matchID.String(),
		User1Id:   user1ID.String(),
		User2Id:   user2ID.String(),
		IsActive:  true,
		MatchedAt: timestamppb.New(matchedAt),
	}

	result := coreProtoToMatch(pbMatch)

	assert.Equal(t, matchID, result.ID)
	assert.Equal(t, user1ID, result.User1ID)
	assert.Equal(t, user2ID, result.User2ID)
	assert.True(t, result.IsActive)
}

// Helper function
func strPtr(s string) *string {
	return &s
}
