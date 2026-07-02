package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	adminhandler "github.com/j1hub/backend/internal/admin/adapter/http"
	authhandler "github.com/j1hub/backend/internal/auth/adapter/http"
	authport "github.com/j1hub/backend/internal/auth/port"
	expensehandler "github.com/j1hub/backend/internal/expense/adapter/http"
	friendhandler "github.com/j1hub/backend/internal/friend/adapter/http"
	journeyhandler "github.com/j1hub/backend/internal/gamification/adapter/http"
	jobhandler "github.com/j1hub/backend/internal/job/adapter/http"
	"github.com/j1hub/backend/internal/media"
	missionhandler "github.com/j1hub/backend/internal/mission/adapter/http"
	notifhandler "github.com/j1hub/backend/internal/notification/adapter/http"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	userhandler "github.com/j1hub/backend/internal/user/adapter/http"
)

func NewRouter(
	issuer authport.TokenIssuer,
	authH *authhandler.AuthHandler,
	userH *userhandler.UserHandler,
	missionH *missionhandler.MissionHandler,
	journeyH *journeyhandler.JourneyHandler,
	friendH *friendhandler.FriendHandler,
	expenseH *expensehandler.ExpenseHandler,
	notifH *notifhandler.NotificationHandler,
	jobH *jobhandler.JobHandler,
	adminH *adminhandler.AdminHandler,
	mediaH *media.MediaHandler,
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
			r.Post("/users/me/job", userH.AssignJob)
			r.Get("/users/{id}", userH.GetPublicProfile)
			r.Patch("/users/me/settings", userH.UpdateSettings)
			r.Delete("/users/me", userH.DeleteAccount)
			r.Patch("/profile/location", userH.UpdateLocation)
			r.Put("/users/me/password", userH.UpdatePassword)

			// Mission routes
			r.Get("/missions", missionH.ListMissions)
			r.Get("/user-missions", missionH.ListUserMissions)
			r.Get("/user-missions/{id}", missionH.GetMissionDetail)
			r.Post("/user-missions/{id}/proof", missionH.SubmitProof)
			r.Patch("/user-missions/{id}/tasks/{taskId}", missionH.ToggleTask)
			// New task routes
			r.Get("/tasks", missionH.ListTasks)
			r.Get("/user-tasks", missionH.ListUserTasks)

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
			r.Post("/jobs", jobH.CreateJob)
			r.Put("/jobs/{id}", jobH.UpdateJob)    // <-- ADDED: Update job details completely
			r.Patch("/jobs/{id}", jobH.PatchJob)   // <-- ADDED: Partially update job details (e.g., change status)
			r.Delete("/jobs/{id}", jobH.DeleteJob) // <-- ADDED: Delete or archive a job listing

			r.Get("/jobs/{id}/reviews", jobH.GetJobReviews)
			r.Post("/jobs/{id}/reviews", jobH.CreateReview)
			r.Get("/cart", jobH.ListCart)
			r.Post("/cart", jobH.AddToCart)
			r.Patch("/cart/{cartId}", jobH.UpdateCartStatus)
			r.Delete("/cart/{cartId}", jobH.RemoveFromCart)
			r.Post("/reviews", jobH.CreateReview)

			// Media & Upload routes
			r.Post("/media/upload", mediaH.UploadFile) // return url
			r.Delete("/media/{key}", mediaH.DeleteFile)
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
			r.Post("/admin/missions", adminH.CreateMission)
		})
	})

	return r
}
