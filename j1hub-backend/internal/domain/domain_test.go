package domain

import (
	"testing"
	"time"

	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"
)

func TestPaymentStatusValid(t *testing.T) {
	tests := []struct {
		status expensedomain.PaymentStatus
		valid  bool
	}{
		{expensedomain.PaymentPending, true},
		{expensedomain.PaymentSubmitted, true},
		{expensedomain.PaymentApproved, true},
		{expensedomain.PaymentOverdue, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestApprovalStatusValid(t *testing.T) {
	tests := []struct {
		status expensedomain.ApprovalStatus
		valid  bool
	}{
		{expensedomain.ApprovalPending, true},
		{expensedomain.ApprovalApproved, true},
		{expensedomain.ApprovalRejected, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestExpenseSplitIsSettled(t *testing.T) {
	split := &expensedomain.ExpenseSplit{PaymentStatus: expensedomain.PaymentApproved}
	if !split.IsSettled() {
		t.Error("expected IsSettled() to be true for expensedomain.PaymentApproved")
	}

	split.PaymentStatus = expensedomain.PaymentPending
	if split.IsSettled() {
		t.Error("expected IsSettled() to be false for expensedomain.PaymentPending")
	}
}

func TestFriendshipStatusValid(t *testing.T) {
	tests := []struct {
		status frienddomain.FriendshipStatus
		valid  bool
	}{
		{frienddomain.FriendshipPending, true},
		{frienddomain.FriendshipAccepted, true},
		{frienddomain.FriendshipBlocked, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestCanonicalOrder(t *testing.T) {
	r1, r2 := frienddomain.CanonicalOrder("alice", "bob")
	if r1 != "alice" || r2 != "bob" {
		t.Errorf("frienddomain.CanonicalOrder(\"alice\", \"bob\") = (%s, %s); want (alice, bob)", r1, r2)
	}

	r1, r2 = frienddomain.CanonicalOrder("bob", "alice")
	if r1 != "alice" || r2 != "bob" {
		t.Errorf("frienddomain.CanonicalOrder(\"bob\", \"alice\") = (%s, %s); want (alice, bob)", r1, r2)
	}
}

func TestSourceTypeValid(t *testing.T) {
	tests := []struct {
		st    gamificationdomain.SourceType
		valid bool
	}{
		{gamificationdomain.SourceMissionBase, true},
		{gamificationdomain.SourceSpeedBonus, true},
		{gamificationdomain.SourceStreakBonus, true},
		{gamificationdomain.SourceFirstCompleter, true},
		{gamificationdomain.SourceExpensePenalty, true},
		{gamificationdomain.SourceAdminAdjust, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.st.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.st, tc.valid)
		}
	}
}

func TestTriggerTypeValid(t *testing.T) {
	tests := []struct {
		tt    gamificationdomain.TriggerType
		valid bool
	}{
		{gamificationdomain.TriggerSpeed, true},
		{gamificationdomain.TriggerStreak, true},
		{gamificationdomain.TriggerFirstCompleter, true},
		{gamificationdomain.TriggerPhaseComplete, true},
		{gamificationdomain.TriggerManual, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.tt.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.tt, tc.valid)
		}
	}
}

func TestCartStatusValid(t *testing.T) {
	tests := []struct {
		status jobdomain.CartStatus
		valid  bool
	}{
		{jobdomain.CartSaved, true},
		{jobdomain.CartViewed, true},
		{jobdomain.CartApplied, true},
		{jobdomain.CartRemoved, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestJobReviewScoreMap(t *testing.T) {
	review := &jobdomain.JobReview{
		ScoreAgency:               1,
		ScoreJob:                  2,
		ScoreCoworkers:            3,
		ScoreTown:                 4,
		ScoreHours:                5,
		ScoreHousing:              6,
		ScoreSecondJobFeasibility: 7,
		ScoreOvertimeAvailability: 8,
	}
	m := review.ScoreMap()
	if m["agency"] != 1 || m["job"] != 2 || m["coworkers"] != 3 || m["town"] != 4 ||
		m["hours"] != 5 || m["housing"] != 6 || m["second_job_feasibility"] != 7 || m["overtime_availability"] != 8 {
		t.Errorf("ScoreMap failed: got %+v", m)
	}
}

func TestUserMissionStatusValid(t *testing.T) {
	tests := []struct {
		status missiondomain.UserMissionStatus
		valid  bool
	}{
		{missiondomain.StatusNotStarted, true},
		{missiondomain.StatusInProgress, true},
		{missiondomain.StatusPendingVerification, true},
		{missiondomain.StatusCompleted, true},
		{missiondomain.StatusOverdue, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestVerificationTypeValid(t *testing.T) {
	tests := []struct {
		vt    missiondomain.VerificationType
		valid bool
	}{
		{missiondomain.VerificationNone, true},
		{missiondomain.VerificationUpload, true},
		{missiondomain.VerificationAdmin, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.vt.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.vt, tc.valid)
		}
	}
}

func TestMissionCalculateDueDate(t *testing.T) {
	now := time.Now()
	// Relative Type
	m1 := &missiondomain.Mission{
		DueDateType:        "Relative",
		RelativeDaysOffset: 5,
	}
	expected1 := now.AddDate(0, 0, 5)
	res1 := m1.CalculateDueDate(now)
	if !res1.Equal(expected1) {
		t.Errorf("Relative: expected %v, got %v", expected1, res1)
	}

	// Fixed Type
	fixed := now.Add(24 * time.Hour)
	m2 := &missiondomain.Mission{
		DueDateType:  "fixed",
		FixedDueDate: &fixed,
	}
	res2 := m2.CalculateDueDate(now)
	if !res2.Equal(fixed) {
		t.Errorf("Fixed: expected %v, got %v", fixed, res2)
	}
}

func TestCanAdvancePhase(t *testing.T) {
	missions := []missiondomain.UserMission{
		{Status: missiondomain.StatusCompleted},
		{Status: missiondomain.StatusCompleted},
	}
	if !missiondomain.CanAdvancePhase(missions) {
		t.Error("expected missiondomain.CanAdvancePhase to return true when all completed")
	}

	missions2 := []missiondomain.UserMission{
		{Status: missiondomain.StatusCompleted},
		{Status: missiondomain.StatusInProgress},
	}
	if missiondomain.CanAdvancePhase(missions2) {
		t.Error("expected missiondomain.CanAdvancePhase to return false when one is in progress")
	}
}

func TestRadarVisibilityValid(t *testing.T) {
	tests := []struct {
		v     userdomain.RadarVisibility
		valid bool
	}{
		{userdomain.VisibilityShowAnonymous, true},
		{userdomain.VisibilityShowFriends, true},
		{userdomain.VisibilityHidden, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.v.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.v, tc.valid)
		}
	}
}
