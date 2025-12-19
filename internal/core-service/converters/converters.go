package converters

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ============= Domain → Proto =============

func UserToProto(user *domain.User) *pb.User {
	if user == nil {
		return nil
	}

	pbUser := &pb.User{
		Id:              user.ID.String(),
		Email:           user.Email,
		Name:            user.Name,
		BirthDate:       timestamppb.New(user.BirthDate),
		Gender:          string(user.Gender),
		IsVerified:      user.IsVerified,
		IsPremium:       user.IsPremium,
		SuperLikesCount: int32(user.SuperLikesCount),
		LastActive:      timestamppb.New(user.LastActive),
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}

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

func UserPhotoToProto(photo domain.UserPhoto) *pb.UserPhoto {
	return &pb.UserPhoto{
		Id:           photo.ID.String(),
		UserId:       photo.UserID.String(),
		PhotoUrl:     photo.PhotoURL,
		DisplayOrder: int32(photo.DisplayOrder),
		IsApproved:   photo.IsApproved,
		CreatedAt:    timestamppb.New(photo.CreatedAt),
	}
}

func PhotosToProto(photos []domain.UserPhoto) []*pb.UserPhoto {
	pbPhotos := make([]*pb.UserPhoto, len(photos))
	for i, photo := range photos {
		pbPhotos[i] = UserPhotoToProto(photo)
	}
	return pbPhotos
}

func PreferenceToProto(pref *domain.UserPreference) *pb.UserPreference {
	if pref == nil {
		return nil
	}

	return &pb.UserPreference{
		UserId:       pref.UserID.String(),
		ShowGender:   string(pref.ShowGender),
		AgeMin:       int32(pref.AgeMin),
		AgeMax:       int32(pref.AgeMax),
		MaxDistance:  int32(pref.MaxDistance),
		GlobalSearch: pref.GlobalSearch,
		CreatedAt:    timestamppb.New(pref.CreatedAt),
		UpdatedAt:    timestamppb.New(pref.UpdatedAt),
	}
}

func InterestToProto(interest domain.Interest) *pb.Interest {
	return &pb.Interest{
		UserId: interest.UserID.String(),
		Theme:  string(interest.Theme),
	}
}

func InterestsToProto(interests []domain.Interest) []*pb.Interest {
	pbInterests := make([]*pb.Interest, len(interests))
	for i, interest := range interests {
		pbInterests[i] = InterestToProto(interest)
	}
	return pbInterests
}

func FeedUserToProto(feedUser dto.FeedUser) *pb.FeedUser {
	pbFeedUser := &pb.FeedUser{
		Id:          feedUser.ID,
		Name:        feedUser.Name,
		Age:         int32(feedUser.Age),
		Gender:      feedUser.Gender,
		Description: feedUser.Description,
		Images:      feedUser.Images,
		PhotosCount: int32(feedUser.PhotosCount),
		IsPremium:   feedUser.IsPremium,
		Interests:   InterestsToProto(feedUser.Interests),
	}

	if feedUser.Artist != nil {
		pbFeedUser.Artist = feedUser.Artist
	}
	if feedUser.Quote != nil {
		pbFeedUser.Quote = feedUser.Quote
	}

	return pbFeedUser
}

func FeedUsersToProto(feedUsers []dto.FeedUser) []*pb.FeedUser {
	pbFeedUsers := make([]*pb.FeedUser, len(feedUsers))
	for i, feedUser := range feedUsers {
		pbFeedUsers[i] = FeedUserToProto(feedUser)
	}
	return pbFeedUsers
}

func MatchToProto(match domain.Match) *pb.Match {
	pbMatch := &pb.Match{
		Id:        match.ID.String(),
		User1Id:   match.User1ID.String(),
		User2Id:   match.User2ID.String(),
		IsActive:  match.IsActive,
		MatchedAt: timestamppb.New(match.MatchedAt),
	}

	// Добавляем expires_at, если оно не nil
	if match.ExpiresAt != nil {
		pbMatch.ExpiresAt = timestamppb.New(*match.ExpiresAt)
	}

	return pbMatch
}

func MatchResponseToProto(matchResp dto.MatchResponse) *pb.MatchResponse {
	return &pb.MatchResponse{
		Match:       MatchToProto(matchResp.Match),
		User:        UserToProto(&matchResp.User),
		Photos:      matchResp.Photos,
		Age:         int32(matchResp.Age),
		Description: matchResp.Description,
		PhotosCount: int32(matchResp.PhotosCount),
	}
}

func MatchResponsesToProto(matchResps []dto.MatchResponse) []*pb.MatchResponse {
	pbMatchResps := make([]*pb.MatchResponse, len(matchResps))
	for i, matchResp := range matchResps {
		pbMatchResps[i] = MatchResponseToProto(matchResp)
	}
	return pbMatchResps
}

// ============= Proto → Domain =============

func ProtoToProfileUpdateRequest(req *pb.UpdateProfileInfoRequest) *domain.ProfileUpdateRequest {
	updateData := &domain.ProfileUpdateRequest{}

	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Phone != nil {
		updateData.Phone = req.Phone
	}
	if req.BirthDate != nil {
		birthDate := req.BirthDate.AsTime()
		updateData.BirthDate = &birthDate
	}
	if req.Gender != nil {
		gender := constants.Gender(*req.Gender)
		updateData.Gender = gender
	}
	if req.Bio != nil {
		updateData.Bio = req.Bio
	}
	if req.Artist != nil {
		updateData.Artist = req.Artist
	}
	if req.Quote != nil {
		updateData.Quote = req.Quote
	}
	if req.City != nil {
		updateData.City = req.City
	}
	if req.Latitude != nil {
		updateData.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		updateData.Longitude = req.Longitude
	}

	return updateData
}

func ProtoToPreferencesUpdateRequest(req *pb.UpdatePreferencesRequest) *domain.PreferencesUpdateRequest {
	updateData := &domain.PreferencesUpdateRequest{}

	if req.ShowGender != nil {
		showGender := constants.GenderPreference(*req.ShowGender)
		updateData.ShowGender = showGender
	}
	if req.AgeMin != nil {
		updateData.AgeMin = int(*req.AgeMin)
	}
	if req.AgeMax != nil {
		updateData.AgeMax = int(*req.AgeMax)
	}

	return updateData
}

func ProtoToUser(pbUser *pb.User) *domain.User {
	if pbUser == nil {
		return nil
	}

	id, _ := uuid.Parse(pbUser.Id)

	user := &domain.User{
		ID:         id,
		Email:      pbUser.Email,
		Name:       pbUser.Name,
		BirthDate:  pbUser.BirthDate.AsTime(),
		Gender:     constants.Gender(pbUser.Gender),
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

// StrikeToProto преобразует доменную сущность Strike в proto сообщение
func StrikeToProto(strike *domain.Strike) *pb.Strike {
	if strike == nil {
		return nil
	}

	strikeProto := &pb.Strike{
		Id:           strike.ID.String(),
		ReporterId:   strike.ReporterID.String(),
		TargetUserId: strike.TargetUserID.String(),
		Type:         string(strike.Type),
		Reason:       strike.Reason,
		Status:       string(strike.Status),
		CreatedAt:    timestamppb.New(strike.CreatedAt),
	}

	if strike.UpdatedAt != nil {
		strikeProto.UpdatedAt = timestamppb.New(*strike.UpdatedAt)
	}

	if strike.ModeratorID != nil {
		moderatorID := strike.ModeratorID.String()
		strikeProto.ModeratorId = moderatorID
	}

	if strike.ModeratorNote != nil {
		strikeProto.ModeratorNote = *strike.ModeratorNote
	}

	return strikeProto
}

// StrikesToProto преобразует список доменных сущностей Strike в proto сообщения
func StrikesToProto(strikes []*domain.Strike) []*pb.Strike {
	if strikes == nil {
		return nil
	}

	strikeProtos := make([]*pb.Strike, len(strikes))
	for i, strike := range strikes {
		strikeProtos[i] = StrikeToProto(strike)
	}

	return strikeProtos
}

// StrikeStatsToProto преобразует статистику жалоб в proto сообщение
func StrikeStatsToProto(stats *dto.StrikeStats) *pb.StrikeStats {
	if stats == nil {
		return nil
	}

	statsProto := &pb.StrikeStats{
		UserId:       stats.UserID,
		TotalStrikes: int32(stats.TotalStrikes),
		StrikeTypes: &pb.StrikeTypeStat{
			Pending:  int32(stats.StrikeTypes.Pending),
			Approved: int32(stats.StrikeTypes.Approved),
			Rejected: int32(stats.StrikeTypes.Rejected),
			Resolved: int32(stats.StrikeTypes.Resolved),
		},
	}

	if stats.LastStrikeAt != nil {
		statsProto.LastStrikeAt = timestamppb.New(*stats.LastStrikeAt)
	}

	return statsProto
}

// ProtoToStrikeCreateRequest преобразует proto сообщение в DTO для создания жалобы
func ProtoToStrikeCreateRequest(req *pb.CreateStrikeRequest) *dto.StrikeCreateRequest {
	reporterID, _ := uuid.Parse(req.ReporterId)
	targetUserID, _ := uuid.Parse(req.TargetUserId)

	return &dto.StrikeCreateRequest{
		ReporterID:   reporterID,
		TargetUserID: targetUserID,
		Type:         constants.StrikeType(req.Type),
		Reason:       req.Reason,
	}
}

// ProtoToStrikeStatusUpdateRequest преобразует proto сообщение в DTO для обновления статуса жалобы
func ProtoToStrikeStatusUpdateRequest(req *pb.UpdateStrikeStatusRequest) *dto.StrikeStatusUpdateRequest {
	var moderatorID *uuid.UUID
	if req.ModeratorId != "" {
		id, _ := uuid.Parse(req.ModeratorId)
		moderatorID = &id
	}

	var note *string
	if req.Note != "" {
		note = &req.Note
	}

	return &dto.StrikeStatusUpdateRequest{
		Status:      constants.StrikeStatus(req.Status),
		ModeratorID: moderatorID,
		Note:        note,
	}
}

// ============= Report/Support Converters =============

func SupportTicketResponseToProto(ticket *dto.SupportTicketResponse) *pb.Report {
	return &pb.Report{
		Id:        ticket.ID,
		UserId:    "", // Not available in SupportTicketResponse
		Category:  ticket.Category,
		Status:    ticket.Status,
		CreatedAt: timestamppb.New(ticket.CreatedAt),
	}
}

func SupportTicketDetailResponseToProto(ticket *dto.SupportTicketDetailResponse) *pb.ReportDetail {
	return &pb.ReportDetail{
		Id:        ticket.ID,
		Category:  ticket.Category,
		Text:      ticket.Text,
		Email:     ticket.Email,
		Status:    ticket.Status,
		CreatedAt: timestamppb.New(ticket.CreatedAt),
	}
}
