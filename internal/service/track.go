package service

import (
	"spl-notification/internal/dto/request"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
	"spl-notification/internal/repository"
)

type trackServiceImpl struct {
	trackRepository     repository.TrackRepository
	accessService       AccessService
	notificationService NotificationService
}

func NewTrackServiceImpl(
	trackRepository repository.TrackRepository,
	accessService AccessService,
	notificationService NotificationService,
) TrackService {
	return &trackServiceImpl{
		trackRepository:     trackRepository,
		accessService:       accessService,
		notificationService: notificationService,
	}
}

func (t *trackServiceImpl) SendAllFollows(chatId string) *errors.AppError {
	followTracks, err := t.trackRepository.GetTracksByChatId(chatId)
	if err != nil {
		return err
	}

	err = t.notificationService.SendTracks(chatId, followTracks)
	if err != nil {
		return err
	}

	return nil

}

func (t *trackServiceImpl) GetFollowTracksByChatId(chatId string) ([]*model.Track, *errors.AppError) {
	followTracks, err := t.trackRepository.GetTracksByChatId(chatId)
	if err != nil {
		return nil, err
	}

	return followTracks, nil
}

func (t *trackServiceImpl) Create(trackDTO *request.CreateTrackDTO) *errors.AppError {
	accesses, err := t.accessService.GetRecentlyAccess()
	if err != nil {
		return err
	}

	var userAccess *model.Access
	for _, access := range accesses {
		if access.ExternalID == trackDTO.ExternalID {
			userAccess = access
			break
		}
	}

	if userAccess != nil {
		trackDTO.LastEntry = &userAccess.EntryAt
		trackDTO.LastExit = userAccess.ExitAt
	}

	err = t.trackRepository.Create(trackDTO)
	if err != nil {
		return err
	}

	err = t.notificationService.SendMessage(trackDTO.ChatID, "✅ Agregado")
	if err != nil {
		return err
	}

	return nil
}

func (t *trackServiceImpl) Delete(deleteDTO *request.DeleteTrackDTO) *errors.AppError {
	err := t.trackRepository.Delete(deleteDTO)
	if err != nil {
		return err
	}

	err = t.notificationService.SendMessage(deleteDTO.ChatID, "✅ Eliminado")
	if err != nil {
		return err
	}

	return nil
}

func (t *trackServiceImpl) error(err error) *errors.AppError {
	return errors.NewAppError("TrackService", err)
}
