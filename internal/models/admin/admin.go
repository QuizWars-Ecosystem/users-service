package admin

import (
	"time"

	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
)

const (
	ID        Order = "id"
	Username  Order = "username"
	Email     Order = "email"
	Rating    Order = "rating"
	Coins     Order = "coins"
	CreatedAt Order = "created_at"
	DeletedAt Order = "deleted_at"
)

const (
	ASC  Sort = "ASC"
	DESC Sort = "DESC"
)

type Filter[T any] struct {
	From T
	To   T
}

type SearchFilter struct {
	Offset          uint64
	Limit           uint64
	Order           Order
	Sort            Sort
	RatingFilter    *Filter[int32]
	CoinsFilter     *Filter[int64]
	CreatedAtFilter *Filter[time.Time]
	DeletedAtFilter *Filter[time.Time]
}

type Order string

func (s Order) String() string {
	return string(s)
}

type Sort string

func (s Sort) String() string {
	return string(s)
}

func orderFromGRPCEnum(status userspb.Order) Order {
	switch status {
	case userspb.Order_ORDER_ID:
		return ID
	case userspb.Order_ORDER_USERNAME:
		return Username
	case userspb.Order_ORDER_EMAIL:
		return Email
	case userspb.Order_ORDER_RATING:
		return Rating
	case userspb.Order_ORDER_COINS:
		return Coins
	case userspb.Order_ORDER_CREATED_AT:
		return CreatedAt
	case userspb.Order_ORDER_DELETED_AT:
		return DeletedAt
	default:
		return Username
	}
}

func sortFromGRPCEnum(status userspb.Sort) Sort {
	switch status {
	case userspb.Sort_SORT_ASC:
		return ASC
	case userspb.Sort_SORT_DESC:
		return DESC
	default:
		return DESC
	}
}
