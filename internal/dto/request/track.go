package request

import "time"

type DeleteTrackDTO struct {
	ChatID string `json:"chatId" validate:"required"`
	Run    string `json:"run" validate:"required"`
}

type CreateTrackDTO struct {
	ChatID     string     `json:"chatId" validate:"required"`
	ExternalID int32      `json:"externalId"`
	Run        string     `json:"run" validate:"required"`
	Alias      *string    `json:"alias" validate:"omitempty,min=1,max=100"`
	FullName   string     `json:"fullName"`
	LastEntry  *time.Time `json:"lastEntry"`
	LastExit   *time.Time `json:"lastExit"`
}
