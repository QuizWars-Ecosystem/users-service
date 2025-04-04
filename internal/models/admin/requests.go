package admin

import (
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
)

var _ abstractions.Requestable[SearchFilter, *userspb.SearchUsersRequest] = (*SearchFilter)(nil)

func (s SearchFilter) Request(req *userspb.SearchUsersRequest) (*SearchFilter, error) {
	s.Offset, s.Limit = offsetLimit(req.Page, req.Size)

	if req.Order != nil {
		s.Order = orderFromGRPCEnum(*req.Order)
	} else {
		s.Order = Username
	}

	if req.Sort != nil {
		s.Sort = sortFromGRPCEnum(*req.Sort)
	} else {
		s.Sort = DESC
	}

	if req.UserRating != nil {
		s.RatingFilter = &Filter[int32]{
			From: req.UserRating.From,
			To:   req.UserRating.To,
		}
	}

	if req.UserCoins != nil {
		s.CoinsFilter = &Filter[int64]{
			From: req.UserCoins.From,
			To:   req.UserCoins.To,
		}
	}

	if req.UserCreatedAt != nil {
		s.CreatedAtFilter = &Filter[time.Time]{
			From: req.UserCreatedAt.From.AsTime(),
			To:   req.UserCreatedAt.To.AsTime(),
		}
	}

	if req.UserDeletedAt != nil {
		s.DeletedAtFilter = &Filter[time.Time]{
			From: req.UserDeletedAt.From.AsTime(),
			To:   req.UserDeletedAt.To.AsTime(),
		}
	}

	return &s, nil
}

func offsetLimit(page, size uint64) (uint64, uint64) {
	offset := (page - 1) * size
	limit := size

	return offset, limit
}
