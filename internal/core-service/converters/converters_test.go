package converters

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ============= Domain → Proto Tests =============

func TestUserToProto(t *testing.T) {
	tests := []struct {
		name string
		user *domain.User
	}{
		{
			name: "nil_user",
			user: nil,
		},
		{
			name: "full_user_with_all_fields",
			user: &domain.User{
				ID:         uuid.New(),
				Email:      "test@example.com",
				Name:       "John Doe",
				BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Gender:     constants.GenderMale,
				Phone:      strPtr("1234567890"),
				Bio:        strPtr("Test bio"),
				City:       strPtr("Moscow"),
				Artist:     strPtr("The Beatles"),
				Quote:      strPtr("To be or not to be"),
				IsVerified: true,
				LastActive: time.Now(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			name: "minimal_user_without_optional_fields",
			user: &domain.User{
				ID:         uuid.New(),
				Email:      "minimal@example.com",
				Name:       "Jane",
				BirthDate:  time.Now(),
				Gender:     constants.GenderFemale,
				IsVerified: false,
				LastActive: time.Now(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UserToProto(tt.user)

			if tt.user == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.user.ID.String(), result.Id)
			assert.Equal(t, tt.user.Email, result.Email)
			assert.Equal(t, tt.user.Name, result.Name)
			assert.Equal(t, string(tt.user.Gender), result.Gender)
			assert.Equal(t, tt.user.IsVerified, result.IsVerified)

			// Check optional fields
			if tt.user.Phone != nil {
				assert.Equal(t, tt.user.Phone, result.Phone)
			}
			if tt.user.Bio != nil {
				assert.Equal(t, tt.user.Bio, result.Bio)
			}
		})
	}
}

func TestUserPhotoToProto(t *testing.T) {
	photo := domain.UserPhoto{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		PhotoURL:     "https://example.com/photo.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
		CreatedAt:    time.Now(),
	}

	result := UserPhotoToProto(photo)

	assert.NotNil(t, result)
	assert.Equal(t, photo.ID.String(), result.Id)
	assert.Equal(t, photo.UserID.String(), result.UserId)
	assert.Equal(t, photo.PhotoURL, result.PhotoUrl)
	assert.Equal(t, int32(photo.DisplayOrder), result.DisplayOrder)
	assert.Equal(t, photo.IsApproved, result.IsApproved)
}

func TestPhotosToProto(t *testing.T) {
	photos := []domain.UserPhoto{
		{ID: uuid.New(), UserID: uuid.New(), PhotoURL: "url1", DisplayOrder: 0},
		{ID: uuid.New(), UserID: uuid.New(), PhotoURL: "url2", DisplayOrder: 1},
	}

	result := PhotosToProto(photos)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "url1", result[0].PhotoUrl)
	assert.Equal(t, "url2", result[1].PhotoUrl)
}

func TestPreferenceToProto(t *testing.T) {
	tests := []struct {
		name string
		pref *domain.UserPreference
	}{
		{
			name: "nil_preference",
			pref: nil,
		},
		{
			name: "valid_preference",
			pref: &domain.UserPreference{
				UserID:       uuid.New(),
				ShowGender:   constants.GenderPrefMale,
				AgeMin:       18,
				AgeMax:       30,
				MaxDistance:  50,
				GlobalSearch: false,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreferenceToProto(tt.pref)

			if tt.pref == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.pref.UserID.String(), result.UserId)
			assert.Equal(t, string(tt.pref.ShowGender), result.ShowGender)
			assert.Equal(t, int32(tt.pref.AgeMin), result.AgeMin)
			assert.Equal(t, int32(tt.pref.AgeMax), result.AgeMax)
		})
	}
}

func TestInterestToProto(t *testing.T) {
	interest := domain.Interest{
		UserID: uuid.New(),
		Theme:  constants.InterestTypeWorkout,
	}

	result := InterestToProto(interest)

	assert.NotNil(t, result)
	assert.Equal(t, interest.UserID.String(), result.UserId)
	assert.Equal(t, string(interest.Theme), result.Theme)
}

func TestInterestsToProto(t *testing.T) {
	userID := uuid.New()
	interests := []domain.Interest{
		{UserID: userID, Theme: constants.InterestTypeWorkout},
		{UserID: userID, Theme: constants.InterestTypeFun},
	}

	result := InterestsToProto(interests)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, string(constants.InterestTypeWorkout), result[0].Theme)
	assert.Equal(t, string(constants.InterestTypeFun), result[1].Theme)
}

func TestFeedUserToProto(t *testing.T) {
	userID := uuid.New()
	feedUser := dto.FeedUser{
		ID:          userID.String(),
		Name:        "Test User",
		Age:         25,
		Gender:      "male",
		Description: "Test description",
		Images:      []string{"img1.jpg", "img2.jpg"},
		PhotosCount: 2,
		Artist:      strPtr("Beatles"),
		Quote:       strPtr("Be yourself"),
		IsPremium:   true,
		Interests: []domain.Interest{
			{UserID: userID, Theme: constants.InterestTypeWorkout},
		},
	}

	result := FeedUserToProto(feedUser)

	assert.NotNil(t, result)
	assert.Equal(t, feedUser.ID, result.Id)
	assert.Equal(t, feedUser.Name, result.Name)
	assert.Equal(t, int32(feedUser.Age), result.Age)
	assert.Equal(t, feedUser.Gender, result.Gender)
	assert.Equal(t, feedUser.IsPremium, result.IsPremium)
	assert.Equal(t, 2, len(result.Images))
	assert.NotNil(t, result.Artist)
	assert.NotNil(t, result.Quote)
}

func TestFeedUsersToProto(t *testing.T) {
	users := []dto.FeedUser{
		{ID: uuid.New().String(), Name: "User1", Age: 20},
		{ID: uuid.New().String(), Name: "User2", Age: 25},
	}

	result := FeedUsersToProto(users)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "User1", result[0].Name)
	assert.Equal(t, "User2", result[1].Name)
}

func TestMatchToProto(t *testing.T) {
	matchedAt := time.Now()
	expiresAt := matchedAt.Add(24 * time.Hour)

	match := domain.Match{
		ID:        uuid.New(),
		User1ID:   uuid.New(),
		User2ID:   uuid.New(),
		IsActive:  true,
		MatchedAt: matchedAt,
		ExpiresAt: &expiresAt,
	}

	result := MatchToProto(match)

	assert.NotNil(t, result)
	assert.Equal(t, match.ID.String(), result.Id)
	assert.Equal(t, match.User1ID.String(), result.User1Id)
	assert.Equal(t, match.User2ID.String(), result.User2Id)
	assert.Equal(t, match.IsActive, result.IsActive)
	assert.NotNil(t, result.ExpiresAt)
	assert.Equal(t, expiresAt.Unix(), result.ExpiresAt.AsTime().Unix())
}

func TestMatchToProto_WithNilExpiresAt(t *testing.T) {
	matchedAt := time.Now()

	match := domain.Match{
		ID:        uuid.New(),
		User1ID:   uuid.New(),
		User2ID:   uuid.New(),
		IsActive:  true,
		MatchedAt: matchedAt,
		ExpiresAt: nil, // Матч активен навсегда
	}

	result := MatchToProto(match)

	assert.NotNil(t, result)
	assert.Nil(t, result.ExpiresAt) // Должно быть nil
}

func TestMatchResponseToProto(t *testing.T) {
	matchResp := dto.MatchResponse{
		Match: domain.Match{
			ID:        uuid.New(),
			User1ID:   uuid.New(),
			User2ID:   uuid.New(),
			IsActive:  true,
			MatchedAt: time.Now(),
		},
		User: domain.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Name:  "Test",
		},
		Photos:      []string{"photo1.jpg"},
		Age:         25,
		Description: "Description",
		PhotosCount: 1,
	}

	result := MatchResponseToProto(matchResp)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Match)
	assert.NotNil(t, result.User)
	assert.Equal(t, 1, len(result.Photos))
	assert.Equal(t, int32(25), result.Age)
}

func TestMatchResponsesToProto(t *testing.T) {
	matches := []dto.MatchResponse{
		{
			Match: domain.Match{ID: uuid.New()},
			User:  domain.User{Name: "User1"},
			Age:   20,
		},
		{
			Match: domain.Match{ID: uuid.New()},
			User:  domain.User{Name: "User2"},
			Age:   25,
		},
	}

	result := MatchResponsesToProto(matches)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "User1", result[0].User.Name)
	assert.Equal(t, "User2", result[1].User.Name)
}

// ============= Proto → Domain Tests =============

func TestProtoToProfileUpdateRequest(t *testing.T) {
	email := "newemail@example.com"
	name := "Updated Name"
	phone := "9876543210"
	bio := "New bio"

	req := &pb.UpdateProfileInfoRequest{
		Email:     &email,
		Name:      &name,
		Phone:     &phone,
		BirthDate: timestamppb.New(time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC)),
		Gender:    strPtr("female"),
		Bio:       &bio,
	}

	result := ProtoToProfileUpdateRequest(req)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Email)
	assert.Equal(t, email, *result.Email)
	assert.Equal(t, name, result.Name)
	assert.NotNil(t, result.Phone)
	assert.Equal(t, phone, *result.Phone)
	assert.NotNil(t, result.Bio)
	assert.NotNil(t, result.BirthDate)
}

func TestProtoToPreferencesUpdateRequest(t *testing.T) {
	showGender := "male"
	ageMin := int32(20)
	ageMax := int32(30)

	req := &pb.UpdatePreferencesRequest{
		ShowGender: &showGender,
		AgeMin:     &ageMin,
		AgeMax:     &ageMax,
	}

	result := ProtoToPreferencesUpdateRequest(req)

	assert.NotNil(t, result)
	assert.Equal(t, constants.GenderPreference(showGender), result.ShowGender)
	assert.Equal(t, int(ageMin), result.AgeMin)
	assert.Equal(t, int(ageMax), result.AgeMax)
}

func TestProtoToUser(t *testing.T) {
	pbUser := &pb.User{
		Id:         uuid.New().String(),
		Email:      "proto@example.com",
		Name:       "Proto User",
		BirthDate:  timestamppb.New(time.Now()),
		Gender:     "male",
		IsVerified: true,
		LastActive: timestamppb.New(time.Now()),
		CreatedAt:  timestamppb.New(time.Now()),
		UpdatedAt:  timestamppb.New(time.Now()),
	}

	result := ProtoToUser(pbUser)

	assert.NotNil(t, result)
	assert.Equal(t, pbUser.Email, result.Email)
	assert.Equal(t, pbUser.Name, result.Name)
	assert.Equal(t, constants.Gender(pbUser.Gender), result.Gender)
}

// ============= Strike Converter Tests =============

func TestStrikeToProto(t *testing.T) {
	tests := []struct {
		name   string
		strike *domain.Strike
	}{
		{
			name:   "nil_strike",
			strike: nil,
		},
		{
			name: "full_strike_with_moderator",
			strike: &domain.Strike{
				ID:            uuid.New(),
				ReporterID:    uuid.New(),
				TargetUserID:  uuid.New(),
				ModeratorID:   uuidPtr(uuid.New()),
				Type:          constants.StrikeTypeSpam,
				Status:        constants.StrikeStatusPending,
				Reason:        "Test reason",
				ModeratorNote: strPtr("Moderator note"),
				CreatedAt:     time.Now(),
				UpdatedAt:     timePtr(time.Now()),
			},
		},
		{
			name: "minimal_strike_without_moderator",
			strike: &domain.Strike{
				ID:           uuid.New(),
				ReporterID:   uuid.New(),
				TargetUserID: uuid.New(),
				Type:         constants.StrikeTypeInappropriate,
				Status:       constants.StrikeStatusPending,
				Reason:       "Test reason",
				CreatedAt:    time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrikeToProto(tt.strike)

			if tt.strike == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.strike.ID.String(), result.Id)
			assert.Equal(t, tt.strike.ReporterID.String(), result.ReporterId)
			assert.Equal(t, tt.strike.TargetUserID.String(), result.TargetUserId)
			assert.Equal(t, string(tt.strike.Type), result.Type)
			assert.Equal(t, string(tt.strike.Status), result.Status)
			assert.Equal(t, tt.strike.Reason, result.Reason)

			if tt.strike.ModeratorID != nil {
				assert.NotEmpty(t, result.ModeratorId)
			}
			if tt.strike.ModeratorNote != nil {
				assert.NotEmpty(t, result.ModeratorNote)
			}
		})
	}
}

func TestStrikesToProto(t *testing.T) {
	strikes := []*domain.Strike{
		{
			ID:           uuid.New(),
			ReporterID:   uuid.New(),
			TargetUserID: uuid.New(),
			Type:         constants.StrikeTypeSpam,
			Status:       constants.StrikeStatusPending,
			Reason:       "Strike 1",
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			ReporterID:   uuid.New(),
			TargetUserID: uuid.New(),
			Type:         constants.StrikeTypeInappropriate,
			Status:       constants.StrikeStatusApproved,
			Reason:       "Strike 2",
			CreatedAt:    time.Now(),
		},
	}

	result := StrikesToProto(strikes)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Strike 1", result[0].Reason)
	assert.Equal(t, "Strike 2", result[1].Reason)
}

func TestStrikeStatsToProto(t *testing.T) {
	lastStrike := time.Now()
	stats := &dto.StrikeStats{
		UserID:       uuid.New().String(),
		TotalStrikes: 5,
		StrikeTypes: dto.StrikeTypeStat{
			Pending:  2,
			Approved: 3,
			Rejected: 0,
			Resolved: 0,
		},
		LastStrikeAt: &lastStrike,
	}

	result := StrikeStatsToProto(stats)

	assert.NotNil(t, result)
	assert.Equal(t, stats.UserID, result.UserId)
	assert.Equal(t, int32(5), result.TotalStrikes)
	assert.NotNil(t, result.StrikeTypes)
	assert.NotNil(t, result.LastStrikeAt)
}

func TestProtoToStrikeCreateRequest(t *testing.T) {
	req := &pb.CreateStrikeRequest{
		ReporterId:   uuid.New().String(),
		TargetUserId: uuid.New().String(),
		Type:         "spam",
		Reason:       "Spam content",
	}

	result := ProtoToStrikeCreateRequest(req)

	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ReporterID)
	assert.NotEqual(t, uuid.Nil, result.TargetUserID)
	assert.Equal(t, constants.StrikeType("spam"), result.Type)
	assert.Equal(t, "Spam content", result.Reason)
}

func TestProtoToStrikeStatusUpdateRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.UpdateStrikeStatusRequest
	}{
		{
			name: "with_moderator_and_note",
			req: &pb.UpdateStrikeStatusRequest{
				StrikeId:    uuid.New().String(),
				Status:      "approved",
				ModeratorId: uuid.New().String(),
				Note:        "Approved by moderator",
			},
		},
		{
			name: "without_moderator_and_note",
			req: &pb.UpdateStrikeStatusRequest{
				StrikeId: uuid.New().String(),
				Status:   "rejected",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProtoToStrikeStatusUpdateRequest(tt.req)

			require.NotNil(t, result)
			assert.Equal(t, constants.StrikeStatus(tt.req.Status), result.Status)

			if tt.req.ModeratorId != "" {
				assert.NotNil(t, result.ModeratorID)
			}
			if tt.req.Note != "" {
				assert.NotNil(t, result.Note)
			}
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func timePtr(t time.Time) *time.Time {
	return &t
}
