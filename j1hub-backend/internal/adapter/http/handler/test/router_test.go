package handler_test
import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/j1hub/backend/internal/adapter/http/handler"
	"github.com/stretchr/testify/assert"
)

func TestRouter_NewRouter(t *testing.T) {
	issuer := new(MockTokenIssuer)
	authH := handler.NewAuthHandler(nil, nil)
	userH := handler.NewUserHandler(nil)
	missionH := handler.NewMissionHandler(nil, nil)
	journeyH := handler.NewJourneyHandler(nil, nil, nil)
	friendH := handler.NewFriendHandler(nil, nil)
	expenseH := handler.NewExpenseHandler(nil)
	notifH := handler.NewNotificationHandler(nil)
	jobH := handler.NewJobHandler(nil)

	r := handler.NewRouter(issuer, authH, userH, missionH, journeyH, friendH, expenseH, notifH, jobH, nil)
	assert.NotNil(t, r)

	// test public route health
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
