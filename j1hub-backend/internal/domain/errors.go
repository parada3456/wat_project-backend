package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrConflict              = errors.New("conflict")
	ErrInvalidInput          = errors.New("invalid input")
	ErrAlreadyCompleted      = errors.New("mission already completed")
	ErrSelfSplit             = errors.New("cannot split expense with yourself")
	ErrDuplicateFriend       = errors.New("friendship already exists")
	ErrPhaseNotComplete      = errors.New("current phase missions not all completed")
	ErrProofAlreadySubmitted = errors.New("proof already submitted for this mission")

	ErrInvalidCredentials      = errors.New("Invalid credentials")
	ErrPasswordHashFailed      = errors.New("Password Hashing Failed")
	ErrTransactionBeginFailed  = errors.New("Transaction Error")
	ErrUserCreationFailed      = errors.New("User Creation Failed")
	ErrProfileCreationFailed   = errors.New("Profile Creation Failed")
	ErrCreditCreationFailed     = errors.New("Credit Creation Failed")
	ErrTransactionCommitFailed = errors.New("Transaction Commit Failed")
	ErrTokenIssueFailed        = errors.New("Token Issue Failed")
)
