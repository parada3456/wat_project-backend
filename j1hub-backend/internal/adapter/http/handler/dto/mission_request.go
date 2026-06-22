package dto

type ToggleTaskReq struct {
	Completed   *bool `json:"completed"`
	IsCompleted *bool `json:"isCompleted"`
}

