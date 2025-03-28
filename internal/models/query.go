package models

type CommonGetALL struct {
	Page   uint32 `json:"page"`
	Limit  uint32 `json:"limit"`
	Search string `json:"search"`
}

type CommonGetByIDRequest struct {
	ID string `json:"id"`
}

type CommonGetByIDResponse struct {
	ID string `json:"id"`
}
