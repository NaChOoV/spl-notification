package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spl-notification/internal/config"
	"spl-notification/internal/dto/response"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
	"spl-notification/internal/repository"
	"strconv"
	"time"
)

type accessServiceImpl struct {
	trackRepository     repository.TrackRepository
	notificationService NotificationService
	enviromentConfig    *config.EnvironmentConfig
}

func NewAccessServiceImpl(
	trackRepository repository.TrackRepository,
	notificationService NotificationService,
	enviromentConfig *config.EnvironmentConfig,
) AccessService {
	return &accessServiceImpl{
		trackRepository:     trackRepository,
		notificationService: notificationService,
		enviromentConfig:    enviromentConfig,
	}
}

func (a *accessServiceImpl) CheckAccess(accessArray []*model.Access) *errors.AppError {
	allTracks, err := a.trackRepository.GetAll()
	if err != nil {
		return err
	}

	matchEntryAtTracks, matchExitAtTracks, err := a.compareTrackAndAccess(accessArray, allTracks)
	if err != nil {
		return err
	}

	notificationRequests := make([]*model.NotificationRequest, 0)
	if len(matchEntryAtTracks) > 0 {
		notificationRequests = append(
			notificationRequests,
			a.createNotificationRequest(model.NotificationTypeEntry, accessArray, matchEntryAtTracks)...,
		)
	}

	if len(matchExitAtTracks) > 0 {
		notificationRequests = append(
			notificationRequests,
			a.createNotificationRequest(model.NotificationTypeExit, accessArray, matchExitAtTracks)...,
		)
	}

	if len(notificationRequests) == 0 {
		return nil
	}

	a.notificationService.SendNotification(notificationRequests)

	return nil
}

func (a *accessServiceImpl) createNotificationRequest(
	notificationType model.NotificationType,
	accesses []*model.Access,
	tracks []*model.Track) []*model.NotificationRequest {
	notificationRequests := make([]*model.NotificationRequest, 0, len(tracks))
	for _, track := range tracks {
		var access *model.Access
		for _, a := range accesses {
			if a.ExternalID == track.ExternalID {
				access = a
				break
			}
		}
		var date time.Time
		if notificationType == model.NotificationTypeEntry {
			date = access.EntryAt
		} else {
			date = *access.ExitAt
		}

		notificationRequests = append(notificationRequests, &model.NotificationRequest{
			Type:     notificationType,
			Date:     date,
			ChatID:   track.ChatID,
			Run:      track.Run,
			FullName: track.FullName,
			Alias:    track.Alias,
			Location: access.Location,
		})
	}

	return notificationRequests
}

func (a *accessServiceImpl) compareTrackAndAccess(accessArray []*model.Access, tracks []*model.Track) ([]*model.Track, []*model.Track, *errors.AppError) {
	matchEntryAtTracks := make([]*model.Track, 0)
	matchExitAtTracks := make([]*model.Track, 0)
	trackToUpdateEntry := make(map[int32]*model.Access)
	trackToUpdateExit := make(map[int32]*model.Access)

	for _, access := range accessArray {
		for _, track := range tracks {
			if access.ExternalID == track.ExternalID {
				// EntryAt comparison
				if track.LastEntry == nil || !access.EntryAt.Equal(*track.LastEntry) {
					matchEntryAtTracks = append(matchEntryAtTracks, track)

					_, exist := trackToUpdateEntry[access.ExternalID]
					if !exist {
						trackToUpdateEntry[access.ExternalID] = access
					}
				}

				// ExitAt comparison
				isDifferent := false
				bothNil := track.LastExit == nil && access.ExitAt == nil

				if !bothNil {
					isDifferent = !track.LastExit.Equal(*access.ExitAt)
				}

				if !isDifferent {
					isDifferent = track.LastExit == nil && access.ExitAt != nil
				}

				if isDifferent {
					matchExitAtTracks = append(matchExitAtTracks, track)

					_, exist := trackToUpdateExit[access.ExternalID]
					if !exist {
						trackToUpdateExit[access.ExternalID] = access
					}
				}

			}

		}
	}

	trackToUpdateEntryArray := make([]*model.Access, 0, len(trackToUpdateEntry))
	for _, access := range trackToUpdateEntry {
		trackToUpdateEntryArray = append(trackToUpdateEntryArray, access)
	}

	trackToUpdateExitArray := make([]*model.Access, 0, len(trackToUpdateExit))
	for _, access := range trackToUpdateExit {
		trackToUpdateExitArray = append(trackToUpdateExitArray, access)
	}

	if len(trackToUpdateEntryArray) > 0 {
		err := a.trackRepository.UpdateEntryAt(trackToUpdateEntryArray)
		if err != nil {
			return nil, nil, err
		}
	}
	if len(trackToUpdateExitArray) > 0 {
		err := a.trackRepository.UpdateExitAt(trackToUpdateExitArray)
		if err != nil {
			return nil, nil, err
		}
	}

	return matchEntryAtTracks, matchExitAtTracks, nil
}

func (a *accessServiceImpl) GetRecentlyAccess() ([]*model.Access, *errors.AppError) {
	url := a.enviromentConfig.AccessServiceBaseUrl + "/api/access/recently"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, a.error(err)
	}
	req.Header.Set("X-Auth-Token", a.enviromentConfig.AccessServiceAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, a.error(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, a.error(fmt.Errorf("error fetching recently access: %s", resp.Status))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, a.error(err)
	}

	var response struct {
		Data []*response.AccessDTO `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, a.error(err)
	}

	// Convert response.AccessDTO to model.Access
	accesses := make([]*model.Access, 0, len(response.Data))
	for _, dto := range response.Data {
		externalID, err := strconv.ParseInt(dto.ExternalID, 10, 32)
		if err != nil {
			return nil, a.error(err)
		}

		location, err := strconv.ParseInt(dto.Location, 10, 8)
		if err != nil {
			return nil, a.error(err)
		}

		// Parsear EntryAt de string a time.Time
		entryAt, err := time.Parse(time.RFC3339, dto.EntryAt)
		if err != nil {
			return nil, a.error(err)
		}

		// Parsear ExitAt si no es nil
		var exitAt *time.Time
		if dto.ExitAt != nil {
			parsed, err := time.Parse(time.RFC3339, *dto.ExitAt)
			if err != nil {
				return nil, a.error(err)
			}
			exitAt = &parsed
		}

		accesses = append(accesses, &model.Access{
			ExternalID: int32(externalID),
			Run:        dto.Run,
			FullName:   dto.FullName,
			Location:   int8(location),
			EntryAt:    entryAt,
			ExitAt:     exitAt,
		})
	}

	return accesses, nil
}

func (a *accessServiceImpl) error(err error) *errors.AppError {
	return errors.NewAppError("AccessService", err)
}
