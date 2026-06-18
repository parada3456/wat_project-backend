package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/port"
)

func NewRouter(
	issuer port.TokenIssuer,
	authH *AuthHandler,
	userH *UserHandler,
	missionH *MissionHandler,
	journeyH *JourneyHandler,
	friendH *FriendHandler,
	expenseH *ExpenseHandler,
	notifH *NotificationHandler,
	jobH *JobHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)
		r.Post("/auth/refresh", authH.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(issuer))

			r.Post("/auth/logout", authH.Logout)

			// User routes
			r.Get("/user/profile", userH.GetProfile)
			r.Patch("/user/profile", userH.UpdateProfile)

			// Mission routes
			r.Get("/missions", missionH.ListMissions)
			r.Get("/missions/{id}", missionH.GetMissionDetail)
			r.Post("/missions/{id}/verify", missionH.SubmitProof)
			r.Patch("/tasks/{id}", missionH.ToggleTask)

			// Journey routes
			r.Get("/journey/phases", journeyH.ListPhases)
			r.Post("/journey/phase/transition", journeyH.AdvancePhase)
			r.Get("/journey/history", journeyH.GetHistory)
			r.Get("/leaderboard", journeyH.GetLeaderboard)
			r.Get("/user/badges", journeyH.ListBadges)
			r.Get("/user/credit-score/history", journeyH.GetCreditHistory)

			// Friend routes
			r.Post("/friends/request", friendH.SendRequest)
			r.Get("/friends/requests/pending", friendH.ListPendingRequests)
			r.Patch("/friends/respond", friendH.RespondToRequest)
			r.Get("/friends", friendH.ListFriends)
			r.Delete("/friends/{id}", friendH.RemoveFriend)
			r.Get("/radar", friendH.GetRadar)

			// Expense routes
			r.Get("/expenses", expenseH.ListExpenses)
			r.Post("/expenses", expenseH.CreateExpense)
			r.Get("/expenses/{id}", expenseH.GetExpenseDetail)
			r.Delete("/expenses/{id}", expenseH.DeleteExpense)
			r.Get("/expenses/pending", expenseH.ListPending)
			r.Post("/expenses/splits/{id}/pay", expenseH.PaySplit)
			r.Patch("/expenses/splits/{id}/approve", expenseH.ApproveSplit)

			// Notification routes
			r.Get("/notifications", notifH.ListNotifications)
			r.Patch("/notifications/{id}/read", notifH.MarkRead)
			r.Patch("/notifications/read-all", notifH.MarkAllRead)
			r.Delete("/notifications/{id}", notifH.DeleteNotification)

			// Job routes
			r.Get("/jobs", jobH.ListJobs)
			r.Get("/jobs/{id}", jobH.GetJobDetail)
			r.Post("/cart", jobH.AddToCart)
			r.Get("/cart", jobH.ListCart)
			r.Delete("/cart/{id}", jobH.RemoveFromCart)
			r.Post("/reviews", jobH.CreateReview)
		})

		// Public review routes for testing
		r.Get("/reviews", jobH.GetAllReviews)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(issuer))
			r.Use(middleware.RequireAdmin)

			// Admin routes
			// ...
		})
	})

	return r
}
