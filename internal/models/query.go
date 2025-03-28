package models

type GetALLRequest struct {
	Page   uint32 `json:"page"`
	Limit  uint32 `json:"limit"`
	Search string `json:"search"`
}
