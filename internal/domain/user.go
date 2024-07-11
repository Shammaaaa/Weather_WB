package domain

type User struct {
	ID       int64  `json:"ID"`
	Username string `json:"username"`
}

type FavoriteCity struct {
	ID       int64  `json:"ID"`
	UserID   int64  `json:"userID"`
	CityID   int64  `json:"cityID"`
	CityName string `json:"cityName"`
}
