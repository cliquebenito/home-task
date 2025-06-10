package models

type CreateBannerRequest struct {
	Name string `json:"name"`
}

type RegisterClickRequest struct {
	To   string `json:"to"`
	From string `json:"from"`
}
