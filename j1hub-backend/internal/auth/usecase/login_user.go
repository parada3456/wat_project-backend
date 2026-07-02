package authusecase

import (
	"context"
	"fmt"
	"log"

	"github.com/parada3456/wat_project-backend/internal/domain"
	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"

	port "github.com/parada3456/wat_project-backend/internal/auth/port"
)

type LoginUseCase struct {
	userRepo    port.UserRepository
	profileRepo port.ProfileRepository
	hasher      port.PasswordHasher
	tokenIssuer port.TokenIssuer
}

func NewLoginUseCase(
	userRepo port.UserRepository,
	profileRepo port.ProfileRepository,
	hasher port.PasswordHasher,
	tokenIssuer port.TokenIssuer,
) *LoginUseCase {
	log.Println("debugprint: entering NewLoginUseCase")
	return &LoginUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
	}
}

type LoginCommand struct {
	Email    string
	Password string
}

func (uc *LoginUseCase) Login(ctx context.Context, cmd LoginCommand) (*userdomain.User, *userdomain.Profile, *port.TokenPair, error) {
	log.Println("debugprint: entering (*LoginUseCase).Login")
	user, err := uc.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: The email or password provided is incorrect.", domain.ErrInvalidCredentials)
	}

	if !uc.hasher.Verify(cmd.Password, user.PasswordHash) {
		return nil, nil, nil, fmt.Errorf("%w: The email or password provided is incorrect.", domain.ErrInvalidCredentials)
	}

	tokens, err := uc.tokenIssuer.Issue(user.UserID, false)
	if err != nil {
		return nil, nil, nil, err
	}

	profile, _ := uc.profileRepo.FindByUserID(ctx, user.UserID)

	return user, profile, tokens, nil
}

func (uc *LoginUseCase) Refresh(ctx context.Context, refreshToken string) (*port.TokenPair, error) {
	log.Println("debugprint: entering (*LoginUseCase).Refresh")
	return uc.tokenIssuer.Refresh(refreshToken)
}
