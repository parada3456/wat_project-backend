package handler

import (
	"encoding/json"
	"log"
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
	adminH *AdminHandler,
) http.Handler {
	log.Println("debugprint: entering NewRouter")
	r := chi.NewRouter()
	r.Use(middleware.CORS)
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
			r.Get("/users/me", userH.GetProfile)
			r.Patch("/users/me", userH.UpdateProfile)
			r.Get("/users/{id}", userH.GetPublicProfile)
			r.Patch("/users/me/settings", userH.UpdateSettings)
			r.Delete("/users/me", userH.DeleteAccount)
			r.Patch("/profile/location", userH.UpdateLocation)

			// Mission routes
			r.Get("/missions", missionH.ListMissions)
			r.Get("/user-missions", missionH.ListUserMissions)
			r.Get("/user-missions/{id}", missionH.GetMissionDetail)
			r.Post("/user-missions/{id}/proof", missionH.SubmitProof)
			r.Patch("/user-missions/{id}/tasks/{taskId}", missionH.ToggleTask)

			// Journey routes
			r.Get("/journey/phases", journeyH.ListPhases)
			r.Post("/journey/phase-transitions", journeyH.AdvancePhase)
			r.Get("/journey/history", journeyH.GetHistory)
			r.Get("/leaderboard", journeyH.GetLeaderboard)
			r.Get("/user/badges", journeyH.ListBadges)
			r.Get("/user/points/ledger", journeyH.GetPointsLedger)
			r.Get("/user/credit-score/history", journeyH.GetCreditHistory)
			r.Get("/user/credit-history", journeyH.GetCreditHistory)

			// Friend routes
			r.Post("/friend-requests", friendH.SendRequest)
			r.Get("/friend-requests", friendH.ListPendingRequests)
			r.Patch("/friend-requests/{id}", friendH.RespondToRequest)
			r.Get("/friends", friendH.ListFriends)
			r.Delete("/friends/{friendshipId}", friendH.RemoveFriend)
			r.Get("/radar", friendH.GetRadar)

			// Expense routes
			r.Get("/expenses", expenseH.ListExpenses)
			r.Post("/expenses", expenseH.CreateExpense)
			r.Get("/expenses/{id}", expenseH.GetExpenseDetail)
			r.Delete("/expenses/{id}", expenseH.DeleteExpense)
			r.Get("/expense-splits", expenseH.ListPending)
			r.Patch("/expenses/{id}/splits/{splitId}", expenseH.UpdateSplit)

			// Notification routes
			r.Get("/notifications", notifH.ListNotifications)
			r.Patch("/notifications/{id}", notifH.MarkRead)
			r.Patch("/notifications", notifH.MarkAllRead)
			r.Delete("/notifications/{id}", notifH.DeleteNotification)

			// Job routes
			r.Get("/jobs", jobH.ListJobs)
			r.Get("/jobs/{id}", jobH.GetJobDetail)
			r.Get("/jobs/{id}/reviews", jobH.GetJobReviews)
			r.Post("/jobs/{id}/reviews", jobH.CreateReview)
			r.Get("/cart", jobH.ListCart)
			r.Post("/cart", jobH.AddToCart)
			r.Patch("/cart/{cartId}", jobH.UpdateCartStatus)
			r.Delete("/cart/{cartId}", jobH.RemoveFromCart)
			r.Post("/reviews", jobH.CreateReview)
		})

		// Public review routes for testing
		r.Get("/reviews", jobH.GetAllReviews)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(issuer))
			r.Use(middleware.RequireAdmin)

			// Admin routes
			r.Get("/admin/dashboard/stats", adminH.GetStats)
			r.Get("/admin/user-missions", adminH.ListPendingVerifications)
			r.Patch("/admin/user-missions/{id}/verify", adminH.VerifyMission)
			r.Get("/admin/users", adminH.ListUsers)
			r.Get("/admin/users/{id}", adminH.GetUserDetail)
			r.Post("/admin/users/{id}/adjust-points", adminH.AdjustPoints)
		})
	})

	return r
}
