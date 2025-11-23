package adapters

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ service.StrikeService = (*StrikeServiceAdapter)(nil)

type StrikeServiceAdapter struct {
	client pb.CoreServiceClient
}

func NewStrikeServiceAdapter(client pb.CoreServiceClient) service.StrikeService {
	return &StrikeServiceAdapter{
		client: client,
	}
}

func (a *StrikeServiceAdapter) CreateStrike(ctx context.Context, strikeData *dto.StrikeCreateRequest) (*domain.Strike, error) {
	resp, err := a.client.CreateStrike(ctx, &pb.CreateStrikeRequest{
		ReporterId:   strikeData.ReporterID.String(),
		TargetUserId: strikeData.TargetUserID.String(),
		Type:         string(strikeData.Type),
		Reason:       strikeData.Reason,
	})
	if err != nil {
		return nil, err
	}

	return protoToStrike(resp.Strike), nil
}

func (a *StrikeServiceAdapter) GetStrikeByID(ctx context.Context, strikeID string) (*domain.Strike, error) {
	resp, err := a.client.GetStrike(ctx, &pb.GetStrikeRequest{
		StrikeId: strikeID,
	})
	if err != nil {
		return nil, err
	}

	return protoToStrike(resp.Strike), nil
}

func (a *StrikeServiceAdapter) GetStrikesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Strike, error) {
	resp, err := a.client.GetStrikesByUserID(ctx, &pb.GetStrikesByUserIDRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	strikes := make([]*domain.Strike, len(resp.Strikes))
	for i, strikeProto := range resp.Strikes {
		strikes[i] = protoToStrike(strikeProto)
	}

	return strikes, nil
}

func (a *StrikeServiceAdapter) GetStrikesByType(ctx context.Context, strikeType constants.StrikeType, limit, offset int) ([]*domain.Strike, error) {
	resp, err := a.client.GetStrikesByType(ctx, &pb.GetStrikesByTypeRequest{
		Type:   string(strikeType),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	strikes := make([]*domain.Strike, len(resp.Strikes))
	for i, strikeProto := range resp.Strikes {
		strikes[i] = protoToStrike(strikeProto)
	}

	return strikes, nil
}

func (a *StrikeServiceAdapter) GetStrikesByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.Strike, error) {
	resp, err := a.client.GetStrikesByDateRange(ctx, &pb.GetStrikesByDateRangeRequest{
		From:   timestamppb.New(from),
		To:     timestamppb.New(to),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	strikes := make([]*domain.Strike, len(resp.Strikes))
	for i, strikeProto := range resp.Strikes {
		strikes[i] = protoToStrike(strikeProto)
	}

	return strikes, nil
}

func (a *StrikeServiceAdapter) UpdateStrikeStatus(ctx context.Context, strikeID string, status constants.StrikeStatus, moderatorID *uuid.UUID, note *string) error {
	req := &pb.UpdateStrikeStatusRequest{
		StrikeId: strikeID,
		Status:   string(status),
	}

	if moderatorID != nil {
		moderatorIDStr := moderatorID.String()
		req.ModeratorId = moderatorIDStr
	}

	if note != nil {
		req.Note = *note
	}

	_, err := a.client.UpdateStrikeStatus(ctx, req)
	return err
}

func (a *StrikeServiceAdapter) DeleteStrike(ctx context.Context, strikeID string) error {
	_, err := a.client.DeleteStrike(ctx, &pb.DeleteStrikeRequest{
		StrikeId: strikeID,
	})
	return err
}

func (a *StrikeServiceAdapter) GetUserStrikeStats(ctx context.Context, userID string) (*dto.StrikeStats, error) {
	resp, err := a.client.GetUserStrikeStats(ctx, &pb.GetUserStrikeStatsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return protoToStrikeStats(resp.Stats), nil
}

// Вспомогательные функции для преобразования proto <-> domain

func protoToStrike(strikeProto *pb.Strike) *domain.Strike {
	if strikeProto == nil {
		return nil
	}

	strike := &domain.Strike{
		Type:   constants.StrikeType(strikeProto.Type),
		Reason: strikeProto.Reason,
		Status: constants.StrikeStatus(strikeProto.Status),
	}

	if id, err := uuid.Parse(strikeProto.Id); err == nil {
		strike.ID = id
	}
	if reporterID, err := uuid.Parse(strikeProto.ReporterId); err == nil {
		strike.ReporterID = reporterID
	}
	if targetUserID, err := uuid.Parse(strikeProto.TargetUserId); err == nil {
		strike.TargetUserID = targetUserID
	}

	if strikeProto.CreatedAt != nil {
		strike.CreatedAt = strikeProto.CreatedAt.AsTime()
	}

	if strikeProto.UpdatedAt != nil {
		updatedAt := strikeProto.UpdatedAt.AsTime()
		strike.UpdatedAt = &updatedAt
	}

	if strikeProto.ModeratorId != "" {
		if moderatorID, err := uuid.Parse(strikeProto.ModeratorId); err == nil {
			strike.ModeratorID = &moderatorID
		}
	}

	if strikeProto.ModeratorNote != "" {
		note := strikeProto.ModeratorNote
		strike.ModeratorNote = &note
	}

	return strike
}

func protoToStrikeStats(statsProto *pb.StrikeStats) *dto.StrikeStats {
	if statsProto == nil {
		return nil
	}

	stats := &dto.StrikeStats{
		UserID:       statsProto.UserId,
		TotalStrikes: int(statsProto.TotalStrikes),
		StrikeTypes:  make(map[constants.StrikeType]int),
	}

	for strikeType, count := range statsProto.StrikeTypes {
		stats.StrikeTypes[constants.StrikeType(strikeType)] = int(count)
	}

	if statsProto.LastStrikeAt != nil {
		lastStrikeAt := statsProto.LastStrikeAt.AsTime()
		stats.LastStrikeAt = &lastStrikeAt
	}

	return stats
}
