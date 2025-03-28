package models

type CommonGetALL struct {
	Page   uint32 `json:"page"`
	Limit  uint32 `json:"limit"`
	Search string `json:"search"`
}

type CommonGetByID struct {
	ID string `json:"id"`
}
