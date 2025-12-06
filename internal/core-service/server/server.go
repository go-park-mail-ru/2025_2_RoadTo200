package server

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/converters"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreServer) UploadPhoto(ctx context.Context, req *pb.UploadPhotoRequest) (*pb.UploadPhotoResponse, error) {
	s.logger.Infof("UploadPhoto called for user_id: %s", req.UserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	photo, err := s.profileService.UploadPhoto(ctx, userID, req.Content, req.ContentType)
	if err != nil {
		s.logger.Errorf("UploadPhoto error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to upload photo: %v", err)
	}

	return &pb.UploadPhotoResponse{
		Photo: converters.UserPhotoToProto(*photo),
	}, nil
}

type CoreServer struct {
	pb.UnimplementedCoreServiceServer
	profileService service.ProfileService
	feedService    service.FeedService
	swipeService   service.SwipeService
	matchService   service.MatchService
	strikeService  service.StrikeService
	logger         logger.Log
}

func NewCoreServer(
	profileService service.ProfileService,
	feedService service.FeedService,
	swipeService service.SwipeService,
	matchService service.MatchService,
	strikeServie service.StrikeService,
	logger logger.Log,
) *CoreServer {
	return &CoreServer{
		profileService: profileService,
		feedService:    feedService,
		swipeService:   swipeService,
		matchService:   matchService,
		strikeService:  strikeServie,
		logger:         logger,
	}
}

// ============= Profile Methods =============

func (s *CoreServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	s.logger.Infof("GetProfile called for user_id: %s", req.UserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	profileResp, err := s.profileService.GetProfile(ctx, userID)
	if err != nil {
		s.logger.Errorf("GetProfile error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get profile: %v", err)
	}

	return &pb.GetProfileResponse{
		User:        converters.UserToProto(profileResp.User),
		Preferences: converters.PreferenceToProto(profileResp.Preferences),
		Photos:      converters.PhotosToProto(profileResp.Photos),
		Interests:   converters.InterestsToProto(profileResp.Interests),
	}, nil
}

func (s *CoreServer) UpdateProfileInfo(ctx context.Context, req *pb.UpdateProfileInfoRequest) (*pb.UpdateProfileInfoResponse, error) {
	s.logger.Infof("UpdateProfileInfo called for user_id: %s", req.UserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	updateData := converters.ProtoToProfileUpdateRequest(req)

	if err := s.profileService.UpdateProfileInfo(ctx, userID, updateData); err != nil {
		s.logger.Errorf("UpdateProfileInfo error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &pb.UpdateProfileInfoResponse{Success: true}, nil
}

func (s *CoreServer) UpdatePreferences(ctx context.Context, req *pb.UpdatePreferencesRequest) (*pb.UpdatePreferencesResponse, error) {
	s.logger.Infof("UpdatePreferences called for user_id: %s", req.UserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	updateData := converters.ProtoToPreferencesUpdateRequest(req)

	if err := s.profileService.UpdatePreferences(ctx, userID, updateData); err != nil {
		s.logger.Errorf("UpdatePreferences error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update preferences: %v", err)
	}

	return &pb.UpdatePreferencesResponse{Success: true}, nil
}

func (s *CoreServer) UpdateInterests(ctx context.Context, req *pb.UpdateInterestsRequest) (*pb.UpdateInterestsResponse, error) {
	s.logger.Infof("UpdateInterests called for user_id: %s", req.UserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	interests := make([]domain.Interest, len(req.Themes))
	for i, theme := range req.Themes {
		interests[i] = domain.Interest{
			UserID: userID,
			Theme:  constants.InterestType(theme),
		}
	}

	if err := s.profileService.UpdateInterests(ctx, userID, interests); err != nil {
		s.logger.Errorf("UpdateInterests error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update interests: %v", err)
	}

	return &pb.UpdateInterestsResponse{Success: true}, nil
}

func (s *CoreServer) DeletePhoto(ctx context.Context, req *pb.DeletePhotoRequest) (*pb.DeletePhotoResponse, error) {
	s.logger.Infof("DeletePhoto called for user_id: %s, photo_id: %s", req.UserId, req.PhotoId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	photoID, err := uuid.Parse(req.PhotoId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid photo_id: %v", err)
	}

	if err := s.profileService.DeletePhoto(ctx, userID, photoID); err != nil {
		s.logger.Errorf("DeletePhoto error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete photo: %v", err)
	}

	return &pb.DeletePhotoResponse{Success: true}, nil
}

func (s *CoreServer) SetPrimaryPhoto(ctx context.Context, req *pb.SetPrimaryPhotoRequest) (*pb.SetPrimaryPhotoResponse, error) {
	s.logger.Infof("SetPrimaryPhoto called for user_id: %s, photo_id: %s", req.UserId, req.PhotoId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	photoID, err := uuid.Parse(req.PhotoId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid photo_id: %v", err)
	}

	if err := s.profileService.SetPrimaryPhoto(ctx, userID, photoID); err != nil {
		s.logger.Errorf("SetPrimaryPhoto error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to set primary photo: %v", err)
	}

	return &pb.SetPrimaryPhotoResponse{Success: true}, nil
}

func (s *CoreServer) ReorderPhotos(ctx context.Context, req *pb.ReorderPhotosRequest) (*pb.ReorderPhotosResponse, error) {
	s.logger.Infof("ReorderPhotos called for user_id: %s", req.UserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	photoIDs := make([]uuid.UUID, len(req.PhotoIds))
	for i, id := range req.PhotoIds {
		photoID, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid photo_id at index %d: %v", i, err)
		}
		photoIDs[i] = photoID
	}

	if err := s.profileService.ReorderPhotos(ctx, userID, photoIDs); err != nil {
		s.logger.Errorf("ReorderPhotos error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to reorder photos: %v", err)
	}

	return &pb.ReorderPhotosResponse{Success: true}, nil
}

// ============= Feed Methods =============

func (s *CoreServer) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (*pb.GetFeedResponse, error) {
	s.logger.Infof("GetFeed called for user_id: %s, limit: %d, offset: %d", req.UserId, req.Limit, req.Offset)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	feedUsers, err := s.feedService.GetFeed(ctx, userID, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Errorf("GetFeed error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get feed: %v", err)
	}

	return &pb.GetFeedResponse{
		Users:  converters.FeedUsersToProto(feedUsers),
		Total:  int32(len(feedUsers)),
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// ============= Swipe Methods =============

func (s *CoreServer) ProcessSwipe(ctx context.Context, req *pb.ProcessSwipeRequest) (*pb.ProcessSwipeResponse, error) {
	s.logger.Infof("ProcessSwipe called for swiper_id: %s, card_id: %s, action: %s", req.SwiperId, req.CardId, req.Action)

	swiperID, err := uuid.Parse(req.SwiperId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid swiper_id: %v", err)
	}

	swipeReq := &dto.SwipeRequest{
		CardID: req.CardId,
		Action: req.Action,
	}

	swipeResp, err := s.swipeService.ProcessSwipe(ctx, swiperID, swipeReq)
	if err != nil {
		s.logger.Errorf("ProcessSwipe error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to process swipe: %v", err)
	}

	resp := &pb.ProcessSwipeResponse{
		IsMatch: swipeResp.IsMatch,
	}
	if swipeResp.MatchID != "" {
		resp.MatchId = &swipeResp.MatchID
	}
	if swipeResp.Message != "" {
		resp.Message = &swipeResp.Message
	}

	return resp, nil
}

// ============= Match Methods =============

func (s *CoreServer) GetUserMatches(ctx context.Context, req *pb.GetUserMatchesRequest) (*pb.GetUserMatchesResponse, error) {
	s.logger.Infof("GetUserMatches called for user_id: %s, limit: %d, offset: %d", req.UserId, req.Limit, req.Offset)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	matchesResp, err := s.matchService.GetUserMatches(ctx, userID, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Errorf("GetUserMatches error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get matches: %v", err)
	}

	return &pb.GetUserMatchesResponse{
		Matches: converters.MatchResponsesToProto(matchesResp.Matches),
		Total:   int32(matchesResp.Total),
		Limit:   req.Limit,
		Offset:  req.Offset,
	}, nil
}

func (s *CoreServer) Unmatch(ctx context.Context, req *pb.UnmatchRequest) (*pb.UnmatchResponse, error) {
	s.logger.Infof("Unmatch called for user_id: %s, target_user_id: %s", req.UserId, req.TargetUserId)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	targetUserID, err := uuid.Parse(req.TargetUserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid target_user_id: %v", err)
	}

	if err := s.matchService.Unmatch(ctx, userID, targetUserID); err != nil {
		s.logger.Errorf("Unmatch error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to unmatch: %v", err)
	}

	return &pb.UnmatchResponse{Success: true}, nil
}

// ============= Strike Methods =============

func (s *CoreServer) CreateStrike(ctx context.Context, req *pb.CreateStrikeRequest) (*pb.CreateStrikeResponse, error) {
	s.logger.Infof("CreateStrike called for reporter_id: %s, target_user_id: %s, type: %s",
		req.ReporterId, req.TargetUserId, req.Type)

	reporterID, err := uuid.Parse(req.ReporterId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid reporter_id: %v", err)
	}

	targetUserID, err := uuid.Parse(req.TargetUserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid target_user_id: %v", err)
	}

	strikeData := &dto.StrikeCreateRequest{
		ReporterID:   reporterID,
		TargetUserID: targetUserID,
		Type:         constants.StrikeType(req.Type),
		Reason:       req.Reason,
	}

	strike, err := s.strikeService.CreateStrike(ctx, strikeData)
	if err != nil {
		s.logger.Errorf("CreateStrike error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create strike: %v", err)
	}

	return &pb.CreateStrikeResponse{
		Strike: converters.StrikeToProto(strike),
	}, nil
}

func (s *CoreServer) GetStrike(ctx context.Context, req *pb.GetStrikeRequest) (*pb.GetStrikeResponse, error) {
	s.logger.Infof("GetStrike called for strike_id: %s", req.StrikeId)

	strike, err := s.strikeService.GetStrikeByID(ctx, req.StrikeId)
	if err != nil {
		s.logger.Errorf("GetStrike error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get strike: %v", err)
	}

	return &pb.GetStrikeResponse{
		Strike: converters.StrikeToProto(strike),
	}, nil
}

func (s *CoreServer) GetStrikesByUserID(ctx context.Context, req *pb.GetStrikesByUserIDRequest) (*pb.GetStrikesByUserIDResponse, error) {
	s.logger.Infof("GetStrikesByUserID called for user_id: %s, limit: %d, offset: %d",
		req.UserId, req.Limit, req.Offset)

	strikes, err := s.strikeService.GetStrikesByUserID(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Errorf("GetStrikesByUserID error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get strikes by user ID: %v", err)
	}

	return &pb.GetStrikesByUserIDResponse{
		Strikes: converters.StrikesToProto(strikes),
	}, nil
}

func (s *CoreServer) GetStrikesByType(ctx context.Context, req *pb.GetStrikesByTypeRequest) (*pb.GetStrikesByTypeResponse, error) {
	s.logger.Infof("GetStrikesByType called for type: %s, limit: %d, offset: %d",
		req.Type, req.Limit, req.Offset)

	strikes, err := s.strikeService.GetStrikesByType(ctx, constants.StrikeType(req.Type), int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Errorf("GetStrikesByType error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get strikes by type: %v", err)
	}

	return &pb.GetStrikesByTypeResponse{
		Strikes: converters.StrikesToProto(strikes),
	}, nil
}

func (s *CoreServer) GetStrikesByDateRange(ctx context.Context, req *pb.GetStrikesByDateRangeRequest) (*pb.GetStrikesByDateRangeResponse, error) {
	s.logger.Infof("GetStrikesByDateRange called from: %v, to: %v, limit: %d, offset: %d",
		req.From, req.To, req.Limit, req.Offset)

	var from, to time.Time

	if req.From != nil {
		from = req.From.AsTime()
	}

	if req.To != nil {
		to = req.To.AsTime()
	}

	strikes, err := s.strikeService.GetStrikesByDateRange(ctx, from, to, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Errorf("GetStrikesByDateRange error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get strikes by date range: %v", err)
	}

	return &pb.GetStrikesByDateRangeResponse{
		Strikes: converters.StrikesToProto(strikes),
	}, nil
}

func (s *CoreServer) UpdateStrikeStatus(ctx context.Context, req *pb.UpdateStrikeStatusRequest) (*pb.UpdateStrikeStatusResponse, error) {
	s.logger.Infof("UpdateStrikeStatus called for strike_id: %s, status: %s",
		req.StrikeId, req.Status)

	var moderatorID *uuid.UUID
	if req.ModeratorId != "" {
		id, err := uuid.Parse(req.ModeratorId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid moderator_id: %v", err)
		}
		moderatorID = &id
	}

	var note *string
	if req.Note != "" {
		note = &req.Note
	}

	err := s.strikeService.UpdateStrikeStatus(ctx, req.StrikeId, constants.StrikeStatus(req.Status), moderatorID, note)
	if err != nil {
		s.logger.Errorf("UpdateStrikeStatus error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update strike status: %v", err)
	}

	// Получаем обновленную жалобу для ответа
	strike, err := s.strikeService.GetStrikeByID(ctx, req.StrikeId)
	if err != nil {
		s.logger.Errorf("Failed to get updated strike: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get updated strike: %v", err)
	}

	return &pb.UpdateStrikeStatusResponse{
		Strike: converters.StrikeToProto(strike),
	}, nil
}

func (s *CoreServer) DeleteStrike(ctx context.Context, req *pb.DeleteStrikeRequest) (*pb.DeleteStrikeResponse, error) {
	s.logger.Infof("DeleteStrike called for strike_id: %s", req.StrikeId)

	err := s.strikeService.DeleteStrike(ctx, req.StrikeId)
	if err != nil {
		s.logger.Errorf("DeleteStrike error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete strike: %v", err)
	}

	return &pb.DeleteStrikeResponse{}, nil
}

func (s *CoreServer) GetUserStrikeStats(ctx context.Context, req *pb.GetUserStrikeStatsRequest) (*pb.GetUserStrikeStatsResponse, error) {
	s.logger.Infof("GetUserStrikeStats called for user_id: %s", req.UserId)

	stats, err := s.strikeService.GetUserStrikeStats(ctx, req.UserId)
	if err != nil {
		s.logger.Errorf("GetUserStrikeStats error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user strike stats: %v", err)
	}

	return &pb.GetUserStrikeStatsResponse{
		Stats: converters.StrikeStatsToProto(stats),
	}, nil
}
