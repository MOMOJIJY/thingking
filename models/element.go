package models

type ElementResponse struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}
