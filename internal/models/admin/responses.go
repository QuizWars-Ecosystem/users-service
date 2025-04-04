package admin

import (
	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"
	userspb "github.com/QuizWars-Ecosystem/users-service/gen/external/users/v1"
	"github.com/QuizWars-Ecosystem/users-service/internal/models/profile"
)

var _ abstractions.Responseable[userspb.SearchUsersResponse] = (*SearchUsersResponse)(nil)

type SearchUsersResponse struct {
	Users  []*profile.UserAdmin
	Page   uint64
	Size   uint64
	Order  Order
	Sort   Sort
	Amount int64
}

func (s *SearchUsersResponse) Response() (*userspb.SearchUsersResponse, error) {
	var res userspb.SearchUsersResponse

	users := make([]*userspb.UserAdmin, len(s.Users))
	for i, user := range s.Users {
		u, err := user.Response()
		if err != nil {
			return nil, err
		}

		users[i] = u
	}

	res.Users = users
	res.Page = s.Page
	res.Size = s.Size
	res.Order = s.Order.ToGRPCEnum()
	res.Sort = s.Sort.ToGRPCEnum()
	res.Amount = s.Amount

	return &res, nil
}
