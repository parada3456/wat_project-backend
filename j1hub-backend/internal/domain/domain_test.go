package domain

import (
	"testing"
	"time"
)

func TestPaymentStatusValid(t *testing.T) {
	tests := []struct {
		status PaymentStatus
		valid  bool
	}{
		{PaymentPending, true},
		{PaymentSubmitted, true},
		{PaymentApproved, true},
		{PaymentOverdue, true},
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
		status ApprovalStatus
		valid  bool
	}{
		{ApprovalPending, true},
		{ApprovalApproved, true},
		{ApprovalRejected, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestExpenseSplitIsSettled(t *testing.T) {
	split := &ExpenseSplit{PaymentStatus: PaymentApproved}
	if !split.IsSettled() {
		t.Error("expected IsSettled() to be true for PaymentApproved")
	}

	split.PaymentStatus = PaymentPending
	if split.IsSettled() {
		t.Error("expected IsSettled() to be false for PaymentPending")
	}
}

func TestFriendshipStatusValid(t *testing.T) {
	tests := []struct {
		status FriendshipStatus
		valid  bool
	}{
		{FriendshipPending, true},
		{FriendshipAccepted, true},
		{FriendshipBlocked, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestCanonicalOrder(t *testing.T) {
	r1, r2 := CanonicalOrder("alice", "bob")
	if r1 != "alice" || r2 != "bob" {
		t.Errorf("CanonicalOrder(\"alice\", \"bob\") = (%s, %s); want (alice, bob)", r1, r2)
	}

	r1, r2 = CanonicalOrder("bob", "alice")
	if r1 != "alice" || r2 != "bob" {
		t.Errorf("CanonicalOrder(\"bob\", \"alice\") = (%s, %s); want (alice, bob)", r1, r2)
	}
}

func TestSourceTypeValid(t *testing.T) {
	tests := []struct {
		st    SourceType
		valid bool
	}{
		{SourceMissionBase, true},
		{SourceSpeedBonus, true},
		{SourceStreakBonus, true},
		{SourceFirstCompleter, true},
		{SourceExpensePenalty, true},
		{SourceAdminAdjust, true},
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
		tt    TriggerType
		valid bool
	}{
		{TriggerSpeed, true},
		{TriggerStreak, true},
		{TriggerFirstCompleter, true},
		{TriggerPhaseComplete, true},
		{TriggerManual, true},
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
		status CartStatus
		valid  bool
	}{
		{CartSaved, true},
		{CartViewed, true},
		{CartApplied, true},
		{CartRemoved, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.status.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.status, tc.valid)
		}
	}
}

func TestJobReviewScoreMap(t *testing.T) {
	review := &JobReview{
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
		status UserMissionStatus
		valid  bool
	}{
		{StatusNotStarted, true},
		{StatusInProgress, true},
		{StatusPendingVerification, true},
		{StatusCompleted, true},
		{StatusOverdue, true},
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
		vt    VerificationType
		valid bool
	}{
		{VerificationNone, true},
		{VerificationUpload, true},
		{VerificationAdmin, true},
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
	m1 := &Mission{
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
	m2 := &Mission{
		DueDateType:  "Fixed",
		FixedDueDate: &fixed,
	}
	res2 := m2.CalculateDueDate(now)
	if !res2.Equal(fixed) {
		t.Errorf("Fixed: expected %v, got %v", fixed, res2)
	}
}

func TestCanAdvancePhase(t *testing.T) {
	missions := []UserMission{
		{Status: StatusCompleted},
		{Status: StatusCompleted},
	}
	if !CanAdvancePhase(missions) {
		t.Error("expected CanAdvancePhase to return true when all completed")
	}

	missions2 := []UserMission{
		{Status: StatusCompleted},
		{Status: StatusInProgress},
	}
	if CanAdvancePhase(missions2) {
		t.Error("expected CanAdvancePhase to return false when one is in progress")
	}
}

func TestRadarVisibilityValid(t *testing.T) {
	tests := []struct {
		v     RadarVisibility
		valid bool
	}{
		{VisibilityShowAnonymous, true},
		{VisibilityShowFriends, true},
		{VisibilityHidden, true},
		{"INVALID", false},
	}
	for _, tc := range tests {
		if tc.v.Valid() != tc.valid {
			t.Errorf("expected Valid() for %s to be %v", tc.v, tc.valid)
		}
	}
}
