package model

import "time"

type NotificationType int8

const (
	NotificationTypeEntry = iota + 1
	NotificationTypeExit
)

func (t NotificationType) String() string {
	switch t {
	case NotificationTypeEntry:
		return "ENTRY"
	case NotificationTypeExit:
		return "EXIT"
	default:
		return "UNKNOWN"
	}
}

func (n NotificationRequest) LocationName() string {
	switch n.Location {
	case 102:
		return "Espacio Urbano"
	case 104:
		return "Calama"
	case 105:
		return "Pacífico"
	case 106:
		return "Arauco"
	case 107:
		return "Iquique"
	case 108:
		return "Angamos"
	default:
		return "Unknown Location"
	}
}

type NotificationRequest struct {
	Type     NotificationType `json:"type"`
	Date     time.Time        `json:"date"`
	ChatID   string           `json:"chatId"`
	Run      string           `json:"run"`
	FullName string           `json:"fullName"`
	Alias    *string          `json:"alias"`
	Location int8             `json:"location"`
}
