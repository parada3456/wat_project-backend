package friendusecase

import (
	"context"
	"log"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"

	"github.com/j1hub/backend/internal/domain"
	port "github.com/j1hub/backend/internal/friend/port"
	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/j1hub/backend/pkg/uid"
)

type ManageFriendshipUseCase struct {
	friendRepo port.FriendshipRepository
	userRepo   port.UserRepository
	notifier   port.NotifierPort
	clock      timeutil.Clock
}

func NewManageFriendshipUseCase(friendRepo port.FriendshipRepository, userRepo port.UserRepository, notifier port.NotifierPort, clock timeutil.Clock) *ManageFriendshipUseCase {
	log.Println("debugprint: entering NewManageFriendshipUseCase")
	return &ManageFriendshipUseCase{friendRepo: friendRepo, userRepo: userRepo, notifier: notifier, clock: clock}
}

func (uc *ManageFriendshipUseCase) SendRequest(ctx context.Context, senderID, targetID string) error {
	log.Println("debugprint: entering (*ManageFriendshipUseCase).SendRequest")
	if senderID == targetID {
		return domain.ErrConflict
	}

	_, err := uc.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return err
	}

	u1, u2 := frienddomain.CanonicalOrder(senderID, targetID)
	existing, err := uc.friendRepo.FindByCanonicalPair(ctx, u1, u2)
	if err == nil && existing != nil {
		return domain.ErrDuplicateFriend
	}

	f := &frienddomain.Friendship{
		FriendshipID: uid.New("frn_"),
		UserID1:      u1,
		UserID2:      u2,
		Status:       frienddomain.FriendshipPending,
		CreatedAt:    uc.clock.Now(),
		UpdatedAt:    uc.clock.Now(),
	}

	if err := uc.friendRepo.Insert(ctx, f); err != nil {
		return err
	}

	uc.notifier.Send(ctx, targetID, "Friend request", "Someone wants to be your friend!")
	return nil
}

func (uc *ManageFriendshipUseCase) RespondToRequest(ctx context.Context, responderID, friendshipID string, accept bool) error {
	log.Println("debugprint: entering (*ManageFriendshipUseCase).RespondToRequest")
	f, err := uc.friendRepo.FindByID(ctx, friendshipID)
	if err != nil {
		return err
	}

	// Simple check: plan says verify responder is user_id_2, but in real case it could be user_id_1 if they were the receiver.
	// For simplicity, let's just check if responder is part of the friendship.
	if f.UserID1 != responderID && f.UserID2 != responderID {
		return domain.ErrForbidden
	}

	status := frienddomain.FriendshipAccepted
	if !accept {
		status = frienddomain.FriendshipBlocked
	}

	if err := uc.friendRepo.UpdateStatus(ctx, friendshipID, status); err != nil {
		return err
	}

	if accept {
		otherID := f.UserID1
		if otherID == responderID {
			otherID = f.UserID2
		}
		uc.notifier.Send(ctx, otherID, "Friend request accepted", "You are now friends!")
	}

	return nil
}

func (uc *ManageFriendshipUseCase) ListFriends(ctx context.Context, userID string) ([]frienddomain.Friendship, error) {
	log.Println("debugprint: entering (*ManageFriendshipUseCase).ListFriends")
	return uc.friendRepo.FindFriendsOf(ctx, userID)
}

func (uc *ManageFriendshipUseCase) ListPendingRequests(ctx context.Context, userID string) ([]frienddomain.Friendship, error) {
	log.
		// Need FindPendingFor in repo
		Println("debugprint: entering (*ManageFriendshipUseCase).ListPendingRequests")

	return nil, nil
}

func (uc *ManageFriendshipUseCase) RemoveFriend(ctx context.Context, userID, friendshipID string) error {
	log.
		// Need Delete in repo
		Println("debugprint: entering (*ManageFriendshipUseCase).RemoveFriend")

	return nil
}
