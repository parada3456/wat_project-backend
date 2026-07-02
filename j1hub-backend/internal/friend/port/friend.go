package port

import (
	"context"

	frienddomain "github.com/parada3456/wat_project-backend/internal/friend/domain"
	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"
)

type FriendshipRepository interface {
	Insert(ctx context.Context, f *frienddomain.Friendship) error
	FindByCanonicalPair(ctx context.Context, u1, u2 string) (*frienddomain.Friendship, error)
	FindByID(ctx context.Context, id string) (*frienddomain.Friendship, error)
	UpdateStatus(ctx context.Context, id string, status frienddomain.FriendshipStatus) error
	FindFriendsOf(ctx context.Context, userID string, limit, offset int) ([]frienddomain.Friendship, int, error)
	FindPendingFor(ctx context.Context, userID string, limit, offset int) ([]frienddomain.Friendship, int, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}
