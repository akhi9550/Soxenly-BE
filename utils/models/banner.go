package models

type BannerRequest struct {
	Title1    string `json:"title1"`
	Title2    string `json:"title2"`
	Subtitle1 string `json:"subtitle1"`
	Subtitle2 string `json:"subtitle2"`
	Link      string `json:"link"`
}

type BannerResponse struct {
	ID        uint   `json:"id"`
	Title1    string `json:"title1"`
	Title2    string `json:"title2"`
	Subtitle1 string `json:"subtitle1"`
	Subtitle2 string `json:"subtitle2"`
	Image     string `json:"image"`
	Link      string `json:"link"`
	IsActive  bool   `json:"is_active"`
}
