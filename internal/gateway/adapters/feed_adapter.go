package adapters

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
)

var _ service.FeedService = (*FeedServiceAdapter)(nil)

type FeedServiceAdapter struct {
	client pb.CoreServiceClient
}

func NewFeedServiceAdapter(client pb.CoreServiceClient) service.FeedService {
	return &FeedServiceAdapter{
		client: client,
	}
}

func (a *FeedServiceAdapter) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]dto.FeedUser, error) {
	resp, err := a.client.GetFeed(ctx, &pb.GetFeedRequest{
		UserId: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	feedUsers := make([]dto.FeedUser, len(resp.Users))
	for i, pbUser := range resp.Users {
		feedUsers[i] = dto.FeedUser{
			ID:          pbUser.Id,
			Name:        pbUser.Name,
			Age:         int(pbUser.Age),
			Gender:      pbUser.Gender,
			Description: pbUser.Description,
			Images:      pbUser.Images,
			PhotosCount: int(pbUser.PhotosCount),
			Artist:      pbUser.Artist,
			Quote:       pbUser.Quote,
			IsPremium:   pbUser.IsPremium,
			Interests:   coreProtoToInterests(pbUser.Interests),
		}
	}

	return feedUsers, nil
}
