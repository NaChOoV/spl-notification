package repository

import (
	"spl-notification/internal/dto/request"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
)

type TrackRepository interface {
	GetAll() ([]*model.Track, *errors.AppError)
	GetTracksByChatId(chatId string) ([]*model.Track, *errors.AppError)
	UpdateEntryAt(accessArray []*model.Access) *errors.AppError
	UpdateExitAt(accessArray []*model.Access) *errors.AppError
	Create(trackDTO *request.CreateTrackDTO) *errors.AppError
	Delete(trackDTO *request.DeleteTrackDTO) *errors.AppError
}
