package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	adminhandler "github.com/parada3456/wat_project-backend/internal/admin/adapter/http"
	authhandler "github.com/parada3456/wat_project-backend/internal/auth/adapter/http"
	expensehandler "github.com/parada3456/wat_project-backend/internal/expense/adapter/http"
	friendhandler "github.com/parada3456/wat_project-backend/internal/friend/adapter/http"
	journeyhandler "github.com/parada3456/wat_project-backend/internal/gamification/adapter/http"
	jobhandler "github.com/parada3456/wat_project-backend/internal/job/adapter/http"
	"github.com/parada3456/wat_project-backend/internal/media"
	missionhandler "github.com/parada3456/wat_project-backend/internal/mission/adapter/http"
	notifhandler "github.com/parada3456/wat_project-backend/internal/notification/adapter/http"
	transporthttp "github.com/parada3456/wat_project-backend/internal/transport/http"
	userhandler "github.com/parada3456/wat_project-backend/internal/user/adapter/http"
	"github.com/stretchr/testify/assert"
)

func TestRouter_NewRouter(t *testing.T) {
	issuer := new(MockTokenIssuer)
	authH := authhandler.NewAuthHandler(nil, nil)
	userH := userhandler.NewUserHandler(nil)
	missionH := missionhandler.NewMissionHandler(nil, nil)
	journeyH := journeyhandler.NewJourneyHandler(nil, nil, nil)
	friendH := friendhandler.NewFriendHandler(nil, nil)
	expenseH := expensehandler.NewExpenseHandler(nil)
	notifH := notifhandler.NewNotificationHandler(nil)
	jobH := jobhandler.NewJobHandler(nil)
	adminH := adminhandler.NewAdminHandler(nil)
	mediaH := media.NewMediaHandler(nil)

	r := transporthttp.NewRouter(issuer, authH, userH, missionH, journeyH, friendH, expenseH, notifH, jobH, adminH, mediaH)
	assert.NotNil(t, r)

	// test public route health
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
