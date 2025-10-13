package model

import "time"

type Track struct {
	ID         int        `json:"id"`
	ChatID     string     `json:"chatId"`
	ExternalID int32      `json:"externalId"`
	Run        string     `json:"run"`
	FullName   string     `json:"fullName"`
	Alias      *string    `json:"alias"`
	LastEntry  *time.Time `json:"lastEntry"`
	LastExit   *time.Time `json:"lastExit"`
}
