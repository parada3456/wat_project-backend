package usecase

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/apperror"
	"net/http"
)

type LoginUseCase struct {
	userRepo    port.UserRepository
	hasher      port.PasswordHasher
	tokenIssuer port.TokenIssuer
}

func NewLoginUseCase(
	userRepo port.UserRepository,
	hasher port.PasswordHasher,
	tokenIssuer port.TokenIssuer,
) *LoginUseCase {
	log.Println("debugprint: entering NewLoginUseCase")
	return &LoginUseCase{
		userRepo:    userRepo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
	}
}

type LoginCommand struct {
	Email    string
	Password string
}

func (uc *LoginUseCase) Login(ctx context.Context, cmd LoginCommand) (*domain.User, *port.TokenPair, error) {
	log.Println("debugprint: entering (*LoginUseCase).Login")
	user, err := uc.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/invalid-credentials",
			Title:  "Invalid credentials",
			Status: http.StatusUnauthorized,
			Detail: "The email or password provided is incorrect.",
		}
	}

	if !uc.hasher.Verify(cmd.Password, user.PasswordHash) {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/invalid-credentials",
			Title:  "Invalid credentials",
			Status: http.StatusUnauthorized,
			Detail: "The email or password provided is incorrect.",
		}
	}

	tokens, err := uc.tokenIssuer.Issue(user.UserID, false)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (uc *LoginUseCase) Refresh(ctx context.Context, refreshToken string) (*port.TokenPair, error) {
	log.Println("debugprint: entering (*LoginUseCase).Refresh")
	return uc.tokenIssuer.Refresh(refreshToken)
}
