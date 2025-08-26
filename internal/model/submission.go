package model

import "time"

type SubmitReq struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type FeedItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}
