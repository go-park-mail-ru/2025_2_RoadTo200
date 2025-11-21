package adapters

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"github.com/google/uuid"
)

type SwipeServiceAdapter struct {
	client pb.CoreServiceClient
}

func NewSwipeServiceAdapter(client pb.CoreServiceClient) service.SwipeService {
	return &SwipeServiceAdapter{
		client: client,
	}
}

func (a *SwipeServiceAdapter) ProcessSwipe(swiperID uuid.UUID, request *dto.SwipeRequest) (*dto.SwipeResponse, error) {
	resp, err := a.client.ProcessSwipe(context.Background(), &pb.ProcessSwipeRequest{
		SwiperId: swiperID.String(),
		CardId:   request.CardID,
		Action:   request.Action,
	})
	if err != nil {
		return nil, err
	}

	swipeResp := &dto.SwipeResponse{
		IsMatch: resp.IsMatch,
	}

	if resp.MatchId != nil {
		swipeResp.MatchID = *resp.MatchId
	}
	if resp.Message != nil {
		swipeResp.Message = *resp.Message
	}

	return swipeResp, nil
}
