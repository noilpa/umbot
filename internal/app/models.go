package app

type UserData struct {
	ID        int     `json:"id"`
	IsActive  bool    `json:"is_active"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}
