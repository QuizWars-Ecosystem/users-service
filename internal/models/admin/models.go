package admin

import (
	"time"

	usersv1 "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
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

const (
	User  Role = "user"
	Admin Role = "admin"
	Super Role = "super"
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

func (o Order) String() string {
	return string(o)
}

type Sort string

func (s Sort) String() string {
	return string(s)
}

func (o Order) ToGRPCEnum() usersv1.Order {
	switch o {
	case ID:
		return usersv1.Order_ORDER_ID
	case Username:
		return usersv1.Order_ORDER_USERNAME
	case Email:
		return usersv1.Order_ORDER_EMAIL
	case Rating:
		return usersv1.Order_ORDER_RATING
	case Coins:
		return usersv1.Order_ORDER_COINS
	case CreatedAt:
		return usersv1.Order_ORDER_CREATED_AT
	case DeletedAt:
		return usersv1.Order_ORDER_DELETED_AT
	default:
		return usersv1.Order_ORDER_USERNAME
	}
}

func (s Sort) ToGRPCEnum() usersv1.Sort {
	switch s {
	case ASC:
		return usersv1.Sort_SORT_ASC
	case DESC:
		return usersv1.Sort_SORT_DESC
	default:
		return usersv1.Sort_SORT_DESC
	}
}

func orderFromGRPCEnum(status usersv1.Order) Order {
	switch status {
	case usersv1.Order_ORDER_ID:
		return ID
	case usersv1.Order_ORDER_USERNAME:
		return Username
	case usersv1.Order_ORDER_EMAIL:
		return Email
	case usersv1.Order_ORDER_RATING:
		return Rating
	case usersv1.Order_ORDER_COINS:
		return Coins
	case usersv1.Order_ORDER_CREATED_AT:
		return CreatedAt
	case usersv1.Order_ORDER_DELETED_AT:
		return DeletedAt
	default:
		return Username
	}
}

func sortFromGRPCEnum(status usersv1.Sort) Sort {
	switch status {
	case usersv1.Sort_SORT_ASC:
		return ASC
	case usersv1.Sort_SORT_DESC:
		return DESC
	default:
		return DESC
	}
}

type Role string

func (r Role) String() string {
	return string(r)
}

func (r Role) ToGRPCEnum() usersv1.Role {
	switch r {
	case User:
		return usersv1.Role_ROLE_USER
	case Admin:
		return usersv1.Role_ROLE_ADMIN
	case Super:
		return usersv1.Role_ROLE_SUPER
	default:
		return usersv1.Role_ROLE_USER
	}
}

func RoleFromGRPCEnum(status usersv1.Role) Role {
	switch status {
	case usersv1.Role_ROLE_USER:
		return User
	case usersv1.Role_ROLE_ADMIN:
		return Admin
	case usersv1.Role_ROLE_SUPER:
		return Super
	default:
		return User
	}
}
