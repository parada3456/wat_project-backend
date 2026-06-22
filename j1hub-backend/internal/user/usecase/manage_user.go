package userusecase

import (
	"context"
	"log"

	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type UserUseCase struct {
	userRepo    port.UserRepository
	profileRepo port.ProfileRepository
	creditRepo  port.CreditScoreRepository
	friendRepo  port.FriendshipRepository
	hasher      port.PasswordHasher
}

func NewUserUseCase(
	userRepo port.UserRepository,
	profileRepo port.ProfileRepository,
	creditRepo port.CreditScoreRepository,
	friendRepo port.FriendshipRepository,
	hasher port.PasswordHasher,
) *UserUseCase {
	log.Println("debugprint: entering NewUserUseCase")
	return &UserUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		creditRepo:  creditRepo,
		friendRepo:  friendRepo,
		hasher:      hasher,
	}
}

type UserProfileResponse struct {
	User        *userdomain.User        `json:"user"`
	Profile     *userdomain.Profile     `json:"profile"`
	CreditScore *userdomain.CreditScore `json:"credit_score_detail,omitempty"`
}

func (uc *UserUseCase) GetProfile(ctx context.Context, userID string) (*UserProfileResponse, error) {
	log.Println("debugprint: entering (*UserUseCase).GetProfile")
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	credit, err := uc.creditRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserProfileResponse{
		User:        user,
		Profile:     profile,
		CreditScore: credit,
	}, nil
}

type UpdateProfileCommand struct {
	FirstName string
	LastName  string
	Bio       string
	AvatarURL string
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, userID string, cmd UpdateProfileCommand) error {
	log.Println("debugprint: entering (*UserUseCase).UpdateProfile")
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.FirstName = cmd.FirstName
	user.LastName = cmd.LastName

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	profile.Bio = cmd.Bio
	profile.AvatarURL = cmd.AvatarURL

	return uc.profileRepo.Update(ctx, profile)
}

func (uc *UserUseCase) UpdateLocation(ctx context.Context, userID string, lat, lng float64) error {
	log.Println("debugprint: entering (*UserUseCase).UpdateLocation")
	return uc.profileRepo.UpdateLocation(ctx, userID, lat, lng)
}

func (uc *UserUseCase) UpdateSettings(ctx context.Context, userID string, settings map[string]interface{}) error {
	log.Println("debugprint: entering (*UserUseCase).UpdateSettings")
	visibilityRaw, exists := settings["radar_visibility"]
	if !exists {
		return nil
	}
	visibilityStr, ok := visibilityRaw.(string)
	if !ok {
		return domain.ErrInvalidInput
	}
	visibility := domain.RadarVisibility(visibilityStr)
	if !visibility.Valid() {
		return domain.ErrInvalidInput
	}
	return uc.profileRepo.UpdateVisibility(ctx, userID, visibility)
}

func (uc *UserUseCase) DeleteAccount(ctx context.Context, userID string, password string) error {
	log.Println("debugprint: entering (*UserUseCase).DeleteAccount")
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if !uc.hasher.Verify(password, user.PasswordHash) {
		return domain.ErrUnauthorized
	}

	return uc.userRepo.Delete(ctx, userID)
}

func (uc *UserUseCase) GetPublicProfile(ctx context.Context, currentUserID, targetUserID string) (*userdomain.User, *userdomain.Profile, error) {
	log.Println("debugprint: entering (*UserUseCase).GetPublicProfile")

	targetUser, err := uc.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return nil, nil, err
	}

	targetProfile, err := uc.profileRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return nil, nil, err
	}

	if currentUserID == targetUserID {
		return targetUser, targetProfile, nil
	}

	switch targetProfile.RadarVisibility {
	case domain.VisibilityHidden:
		return nil, nil, domain.ErrNotFound
	case domain.VisibilityShowFriends:
		u1, u2 := domain.CanonicalOrder(currentUserID, targetUserID)
		f, err := uc.friendRepo.FindByCanonicalPair(ctx, u1, u2)
		if err != nil || f.Status != domain.FriendshipAccepted {
			return nil, nil, domain.ErrNotFound
		}
	case domain.VisibilityShowAnonymous:
		// allowed
	}

	return targetUser, targetProfile, nil
}
