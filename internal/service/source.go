package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spl-notification/internal/config"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
	"time"
)

type sourceServiceImpl struct {
	enviromentConfig *config.EnvironmentConfig
	httpClient       *http.Client
}

func NewSourceServiceImpl(
	enviromentConfig *config.EnvironmentConfig,
) SourceService {
	return &sourceServiceImpl{
		enviromentConfig: enviromentConfig,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (s *sourceServiceImpl) GetABMUserByRun(run string) (*model.ABMUser, *errors.AppError) {
	req, err := http.NewRequest("GET", s.enviromentConfig.SourceBaseUrl+"/user/abm/"+run, nil)
	if err != nil {
		return nil, s.error(err)
	}
	req.Header.Set("X-Auth-String", s.enviromentConfig.AuthString)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, s.error(err)
	}
	defer resp.Body.Close()

	var user model.ABMUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}
		return nil, s.error(err)
	}

	if user.ExternalID == 0 {
		return nil, nil
	}

	return &user, nil
}

func (s *sourceServiceImpl) GetUserByExternalId(externalId int32) (*model.User, *errors.AppError) {
	req, err := http.NewRequest("GET", s.enviromentConfig.SourceBaseUrl+"/user/"+fmt.Sprintf("%d", externalId), nil)
	if err != nil {
		return nil, s.error(err)
	}
	req.Header.Set("X-Auth-String", s.enviromentConfig.AuthString)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, s.error(err)
	}
	defer resp.Body.Close()

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}
		return nil, s.error(err)
	}

	if user.Run == "" {
		return nil, nil
	}

	return &user, nil
}

func (s *sourceServiceImpl) error(err error) *errors.AppError {
	return errors.NewAppError("SourceService", err)
}
