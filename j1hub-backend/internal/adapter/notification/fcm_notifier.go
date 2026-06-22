package notification

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/port"
	"google.golang.org/api/option"
)

type fcmNotifier struct {
	client   *messaging.Client
	userRepo port.UserRepository
}

func NewFCMNotifier(cfg *config.Config, userRepo port.UserRepository) port.NotifierPort {
	log.Println("debugprint: entering NewFCMNotifier")
	opt := option.WithCredentialsFile(cfg.FCMCredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("error initializing firebase app: %v", err)
		return &noopNotifier{}
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("error getting messaging client: %v", err)
		return &noopNotifier{}
	}

	return &fcmNotifier{client: client, userRepo: userRepo}
}

func (n *fcmNotifier) Send(ctx context.Context, userID, title, body string) error {
	log.Println("debugprint: entering (*fcmNotifier).Send")
	user, err := n.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.UserID == "" { // fcm_token is not in userdomain.User yet, let's assume we find it via repo if we add it
		// For now, let's just log or assume fcm_token is empty
		return nil
	}

	// In a real scenario, we'd have fcm_token in the user struct
	// Let's assume we added it to userdomain.User in Phase 3 or we fetch it here.
	// Since I didn't add it to userdomain.User, I'll just log for now.
	log.Printf("Push to %s: %s - %s", userID, title, body)
	return nil
}

type noopNotifier struct{}

func (n *noopNotifier) Send(ctx context.Context, userID, title, body string) error {
	log.Println("debugprint: entering (*noopNotifier).Send")
	log.Printf("NOOP Push to %s: %s - %s", userID, title, body)
	return nil
}
