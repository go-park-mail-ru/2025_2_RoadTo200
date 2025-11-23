package adapters

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
)

var _ service.MatchService = (*MatchServiceAdapter)(nil)

type MatchServiceAdapter struct {
	client pb.CoreServiceClient
}

func NewMatchServiceAdapter(client pb.CoreServiceClient) service.MatchService {
	return &MatchServiceAdapter{
		client: client,
	}
}

func (a *MatchServiceAdapter) GetUserMatches(userID uuid.UUID, limit, offset int) (*dto.MatchesResponse, error) {
	resp, err := a.client.GetUserMatches(context.Background(), &pb.GetUserMatchesRequest{
		UserId: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	matches := make([]dto.MatchResponse, len(resp.Matches))
	for i, pbMatch := range resp.Matches {
		matches[i] = dto.MatchResponse{
			Match:       coreProtoToMatch(pbMatch.Match),
			User:        *coreProtoToUser(pbMatch.User),
			Photos:      pbMatch.Photos,
			Age:         int(pbMatch.Age),
			Description: pbMatch.Description,
			PhotosCount: int(pbMatch.PhotosCount),
		}
	}

	return &dto.MatchesResponse{
		Matches: matches,
		Total:   int(resp.Total),
		Limit:   int(resp.Limit),
		Offset:  int(resp.Offset),
	}, nil
}

func (a *MatchServiceAdapter) Unmatch(userID, targetUserID uuid.UUID) error {
	_, err := a.client.Unmatch(context.Background(), &pb.UnmatchRequest{
		UserId:       userID.String(),
		TargetUserId: targetUserID.String(),
	})
	return err
}
