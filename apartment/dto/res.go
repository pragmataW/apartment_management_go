package dto

type ApartmentResponse struct {
	FlatNo       int    `json:"flat_no"`
	OwnerName    string `json:"owner_name"`
	OwnerSurname string `json:"owner_surname"`
	Mail         string `json:"mail"`
	Password     string `json:"password"`
	DuesCount    int    `json:"dues_count"`
}

type DuesPriceResponse struct{
	DuesPrice float64 `json:"dues_price"`
}

type PayDayResponse struct {
	PayDay int `json:"payday"`
}

type AnnouncementResponse struct {
	AnnouncementID int	`json:"announcement_id"`
	Title          string `json:"title"`
	Content        string `json:"content"`
}