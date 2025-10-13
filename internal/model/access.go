package model

import "time"

type Access struct {
	ExternalID int32      `json:"externalId"`
	Run        string     `json:"run"`
	FullName   string     `json:"fullName"`
	Location   int8       `json:"location"`
	EntryAt    time.Time  `json:"entryAt"`
	ExitAt     *time.Time `json:"exitAt"`
}
