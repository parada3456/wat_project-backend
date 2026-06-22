package dto

type FriendRequestReq struct {
	TargetUserID string `json:"target_user_id" validate:"required"`
}

type RespondFriendReq struct {
	Accept *bool  `json:"accept"`
	Status string `json:"status"`
}
