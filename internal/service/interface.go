package service

import (
	"spl-notification/internal/dto/request"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
)

type AccessService interface {
	CheckAccess(access []*model.Access) *errors.AppError
	GetRecentlyAccess() ([]*model.Access, *errors.AppError)
}

type NotificationService interface {
	SendNotification(tracks []*model.NotificationRequest) *errors.AppError
	HandleNotification()
	SendTracks(chatId string, tracks []*model.Track) *errors.AppError
	SendMessage(chatID string, message string) *errors.AppError
	Close() error
}

type TrackService interface {
	SendAllFollows(chatId string) *errors.AppError
	GetFollowTracksByChatId(chatId string) ([]*model.Track, *errors.AppError)
	Create(trackDTO *request.CreateTrackDTO) *errors.AppError
	Delete(deleteDTO *request.DeleteTrackDTO) *errors.AppError
}

type SourceService interface {
	GetABMUserByRun(run string) (*model.ABMUser, *errors.AppError)
	GetUserByExternalId(externalId int32) (*model.User, *errors.AppError)
}
