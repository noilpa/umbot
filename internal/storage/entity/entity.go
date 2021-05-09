package entity

type ChatData struct {
	ID        int64   `json:"id"`
	IsActive  bool    `json:"is_active"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}
