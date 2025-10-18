package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"spl-notification/internal/config"
	"spl-notification/internal/dto/request"
	"spl-notification/internal/dto/response"
	apperrors "spl-notification/internal/errors"
	"spl-notification/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTrackRepository struct {
	mock.Mock
}

func (m *MockTrackRepository) GetAll() ([]*model.Track, *apperrors.AppError) {
	args := m.Called()
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*apperrors.AppError)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*model.Track), nil
	}
	return args.Get(0).([]*model.Track), args.Get(1).(*apperrors.AppError)
}

func (m *MockTrackRepository) GetTracksByChatId(chatId string) ([]*model.Track, *apperrors.AppError) {
	args := m.Called(chatId)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*apperrors.AppError)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*model.Track), nil
	}
	return args.Get(0).([]*model.Track), args.Get(1).(*apperrors.AppError)
}

func (m *MockTrackRepository) UpdateEntryAt(accesses []*model.Access) *apperrors.AppError {
	args := m.Called(accesses)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

func (m *MockTrackRepository) UpdateExitAt(accesses []*model.Access) *apperrors.AppError {
	args := m.Called(accesses)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

func (m *MockTrackRepository) Create(trackDTO *request.CreateTrackDTO) *apperrors.AppError {
	args := m.Called(trackDTO)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

func (m *MockTrackRepository) Delete(trackDTO *request.DeleteTrackDTO) *apperrors.AppError {
	args := m.Called(trackDTO)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(tracks []*model.NotificationRequest) *apperrors.AppError {
	args := m.Called(tracks)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

func (m *MockNotificationService) HandleNotification() {
	m.Called()
}

func (m *MockNotificationService) SendTracks(chatID string, tracks []*model.Track) *apperrors.AppError {
	args := m.Called(chatID, tracks)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

func (m *MockNotificationService) SendMessage(chatID string, message string) *apperrors.AppError {
	args := m.Called(chatID, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

func (m *MockNotificationService) Close() error {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*apperrors.AppError)
}

// Tests for CheckAccess

func TestCheckAccess_Success_WithEntryMatches(t *testing.T) {
	now := time.Now()
	oldEntry := now.Add(-2 * time.Hour)

	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	expectedTracks := []*model.Track{
		{
			ID:         1,
			ChatID:     "chat123",
			ExternalID: 12345,
			Run:        "12345678-9",
			FullName:   "John Doe",
			Alias:      stringPtr("Johnny"),
			LastEntry:  &oldEntry,
		},
	}

	expectedNotifications := []*model.NotificationRequest{
		{
			Type:     model.NotificationTypeEntry,
			ChatID:   "chat123",
			Run:      "12345678-9",
			FullName: "John Doe",
			Alias:    stringPtr("Johnny"),
			Location: 1,
			Date:     now,
		},
	}

	accesses := []*model.Access{
		{
			ExternalID: 12345,
			Run:        "12345678-9",
			FullName:   "John Doe",
			Location:   1,
			EntryAt:    now,
			ExitAt:     nil,
		},
	}

	mockRepo.On("GetAll").Return(expectedTracks, nil)
	mockRepo.On("UpdateEntryAt", mock.Anything).Return(nil)
	mockNotifyService.On("SendNotification", mock.Anything).Return(nil)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Test: CheckAccess debe completarse sin error
	err := service.CheckAccess(accesses)

	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertCalled(t, "UpdateEntryAt", accesses)
	mockRepo.AssertNotCalled(t, "UpdateExitAt")
	mockNotifyService.AssertCalled(t, "SendNotification", expectedNotifications)
}

func TestCheckAccess_Success_WithExitMatches(t *testing.T) {
	now := time.Now()
	oldEntry := now.Add(-3 * time.Hour)
	oldExit := now.Add(-2 * time.Hour)
	newExit := now

	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	expectedTracks := []*model.Track{
		{
			ID:         1,
			ChatID:     "chat456",
			ExternalID: 67890,
			Run:        "98765432-1",
			FullName:   "Jane Smith",
			Alias:      nil,
			LastEntry:  &oldEntry,
			LastExit:   &oldExit,
		},
	}

	accesses := []*model.Access{
		{
			ExternalID: 67890,
			Run:        "98765432-1",
			FullName:   "Jane Smith",
			Location:   2,
			EntryAt:    oldEntry,
			ExitAt:     &newExit,
		},
	}

	// Setup expectations - solo lo mínimo necesario
	mockRepo.On("GetAll").Return(expectedTracks, nil)
	mockRepo.On("UpdateEntryAt", mock.Anything).Return(nil)
	mockRepo.On("UpdateExitAt", mock.Anything).Return(nil)
	mockNotifyService.On("SendNotification", mock.Anything).Return(nil)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Test: CheckAccess debe completarse sin error
	err := service.CheckAccess(accesses)

	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertNotCalled(t, "UpdateEntryAt")
	mockRepo.AssertCalled(t, "UpdateExitAt", accesses)
	mockNotifyService.AssertCalled(t, "SendNotification", mock.Anything)
}

func TestCheckAccess_Success_WithBothEntryAndExit(t *testing.T) {
	now := time.Now()
	oldEntry1 := now.Add(-3 * time.Hour)
	oldEntry2 := now.Add(-4 * time.Hour)
	oldExit := now.Add(-2 * time.Hour)
	newExit := now

	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	expectedTracks := []*model.Track{
		{
			ID:         1,
			ChatID:     "chat111",
			ExternalID: 11111,
			Run:        "11111111-1",
			FullName:   "User One",
			LastEntry:  &oldEntry1,
		},
		{
			ID:         2,
			ChatID:     "chat222",
			ExternalID: 22222,
			Run:        "22222222-2",
			FullName:   "User Two",
			LastEntry:  &oldEntry2,
			LastExit:   &oldExit,
		},
	}

	// Setup expectations - solo lo mínimo necesario
	mockRepo.On("GetAll").Return(expectedTracks, nil)
	mockRepo.On("UpdateEntryAt", mock.Anything).Return(nil)
	mockRepo.On("UpdateExitAt", mock.Anything).Return(nil)
	mockNotifyService.On("SendNotification", mock.Anything).Return(nil)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	accesses := []*model.Access{
		{
			ExternalID: 11111,
			Run:        "11111111-1",
			FullName:   "User One",
			Location:   1,
			EntryAt:    now,
			ExitAt:     nil,
		},
		{
			ExternalID: 22222,
			Run:        "22222222-2",
			FullName:   "User Two",
			Location:   2,
			EntryAt:    oldEntry2,
			ExitAt:     &newExit,
		},
	}

	// Test: CheckAccess debe completarse sin error
	err := service.CheckAccess(accesses)

	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertCalled(t, "UpdateEntryAt", mock.Anything)
	mockRepo.AssertCalled(t, "UpdateExitAt", mock.Anything)
	mockNotifyService.AssertCalled(t, "SendNotification", mock.Anything)
}

func TestCheckAccess_NoMatches(t *testing.T) {
	now := time.Now()

	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	expectedTracks := []*model.Track{
		{
			ID:         1,
			ExternalID: 99999,
			LastEntry:  &now,
		},
	}

	// Setup expectations - solo lo mínimo necesario
	mockRepo.On("GetAll").Return(expectedTracks, nil)
	mockRepo.On("UpdateEntryAt", mock.Anything).Return(nil)
	mockRepo.On("UpdateExitAt", mock.Anything).Return(nil)
	mockNotifyService.On("SendNotification", mock.Anything).Return(nil)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	accesses := []*model.Access{
		{
			ExternalID: 12345,
			EntryAt:    now,
		},
	}

	// Test: CheckAccess debe completarse sin error
	err := service.CheckAccess(accesses)

	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertNotCalled(t, "UpdateEntryAt")
	mockRepo.AssertNotCalled(t, "UpdateExitAt")
	mockNotifyService.AssertNotCalled(t, "SendNotification")
}

func TestCheckAccess_GetAllError(t *testing.T) {
	expectedError := apperrors.NewAppError("TestError", errors.New("database connection error"))

	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	// Setup expectations
	mockRepo.On("GetAll").Return(nil, expectedError)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	accesses := []*model.Access{
		{
			ExternalID: 12345,
			Run:        "12345678-9",
			FullName:   "John Doe",
			Location:   1,
			EntryAt:    time.Now(),
		},
	}

	err := service.CheckAccess(accesses)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertNotCalled(t, "UpdateEntryAt")
	mockRepo.AssertNotCalled(t, "UpdateExitAt")
	mockNotifyService.AssertNotCalled(t, "SendNotification")
}

func TestCheckAccess_SyncTrackAndAccessError(t *testing.T) {
	now := time.Now()
	oldEntry := now.Add(-1 * time.Hour)
	expectedError := apperrors.NewAppError("TestError", errors.New("sync error"))

	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	expectedTracks := []*model.Track{
		{
			ID:         1,
			ExternalID: 12345,
			LastEntry:  &oldEntry,
		},
	}

	// Setup expectations
	mockRepo.On("GetAll").Return(expectedTracks, nil)
	mockRepo.On("UpdateEntryAt", mock.AnythingOfType("[]*model.Access")).Return(expectedError)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	accesses := []*model.Access{
		{
			ExternalID: 12345,
			EntryAt:    now,
		},
	}

	err := service.CheckAccess(accesses)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertCalled(t, "UpdateEntryAt", mock.AnythingOfType("[]*model.Access"))
	mockRepo.AssertNotCalled(t, "UpdateExitAt")
	mockNotifyService.AssertNotCalled(t, "SendNotification")
}

func TestCheckAccess_EmptyAccessArray(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)

	// Setup expectations - solo lo mínimo necesario
	mockRepo.On("GetAll").Return([]*model.Track{}, nil)
	mockRepo.On("UpdateEntryAt", mock.Anything).Return(nil)
	mockRepo.On("UpdateExitAt", mock.Anything).Return(nil)
	mockNotifyService.On("SendNotification", mock.Anything).Return(nil)

	envConfig := &config.EnvironmentConfig{}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Test: CheckAccess debe completarse sin error con array vacío
	err := service.CheckAccess([]*model.Access{})

	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "GetAll")
	mockRepo.AssertNotCalled(t, "UpdateEntryAt")
	mockRepo.AssertNotCalled(t, "UpdateExitAt")
	mockNotifyService.AssertNotCalled(t, "SendNotification")
}

// Tests for GetCompleteAccess

func TestGetCompleteAccess_Success(t *testing.T) {
	// Setup mock HTTP server
	exitAtValue := "2025-10-05T16:58:54Z"
	mockResponse := struct {
		Data []*response.AccessDTO `json:"data"`
	}{
		Data: []*response.AccessDTO{
			{
				ExternalID: "95729",
				Run:        "21480585-5",
				FullName:   "William Gabriel Araya Pino",
				Location:   "104",
				EntryAt:    "2025-10-05T16:59:48Z",
				ExitAt:     nil,
			},
			{
				ExternalID: "80352",
				Run:        "21209061-1",
				FullName:   "Constanza Paz Saavedra Suarez",
				Location:   "102",
				EntryAt:    "2025-10-05T16:58:54Z",
				ExitAt:     &exitAtValue,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/access/complete", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, accesses)
	assert.Len(t, accesses, 2)

	// Verify first access
	assert.Equal(t, int32(95729), accesses[0].ExternalID)
	assert.Equal(t, "21480585-5", accesses[0].Run)
	assert.Equal(t, "William Gabriel Araya Pino", accesses[0].FullName)
	assert.Equal(t, int8(104), accesses[0].Location)
	expectedEntryAt1, _ := time.Parse(time.RFC3339, "2025-10-05T16:59:48Z")
	assert.Equal(t, expectedEntryAt1, accesses[0].EntryAt)
	assert.Nil(t, accesses[0].ExitAt)

	// Verify second access
	assert.Equal(t, int32(80352), accesses[1].ExternalID)
	assert.Equal(t, "21209061-1", accesses[1].Run)
	assert.Equal(t, "Constanza Paz Saavedra Suarez", accesses[1].FullName)
	assert.Equal(t, int8(102), accesses[1].Location)
	expectedEntryAt2, _ := time.Parse(time.RFC3339, "2025-10-05T16:58:54Z")
	assert.Equal(t, expectedEntryAt2, accesses[1].EntryAt)
	assert.NotNil(t, accesses[1].ExitAt)
	expectedExitAt, _ := time.Parse(time.RFC3339, "2025-10-05T16:58:54Z")
	assert.Equal(t, expectedExitAt, *accesses[1].ExitAt)
}

func TestGetCompleteAccess_ParseIntError_ExternalID(t *testing.T) {
	mockResponse := struct {
		Data []*response.AccessDTO `json:"data"`
	}{
		Data: []*response.AccessDTO{
			{
				ExternalID: "invalid",
				Run:        "21480585-5",
				FullName:   "William Gabriel Araya Pino",
				Location:   "104",
				EntryAt:    "2025-10-05T16:59:48Z",
				ExitAt:     nil,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, accesses)
	assert.Contains(t, err.Error(), "invalid syntax")
}

func TestGetCompleteAccess_ParseIntError_Location(t *testing.T) {
	mockResponse := struct {
		Data []*response.AccessDTO `json:"data"`
	}{
		Data: []*response.AccessDTO{
			{
				ExternalID: "95729",
				Run:        "21480585-5",
				FullName:   "William Gabriel Araya Pino",
				Location:   "invalid",
				EntryAt:    "2025-10-05T16:59:48Z",
				ExitAt:     nil,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, accesses)
	assert.Contains(t, err.Error(), "invalid syntax")
}

func TestGetCompleteAccess_ParseTimeError_EntryAt(t *testing.T) {
	mockResponse := struct {
		Data []*response.AccessDTO `json:"data"`
	}{
		Data: []*response.AccessDTO{
			{
				ExternalID: "95729",
				Run:        "21480585-5",
				FullName:   "William Gabriel Araya Pino",
				Location:   "104",
				EntryAt:    "invalid-time",
				ExitAt:     nil,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, accesses)
}

func TestGetCompleteAccess_ParseTimeError_ExitAt(t *testing.T) {
	exitAtValue := "invalid-time"
	mockResponse := struct {
		Data []*response.AccessDTO `json:"data"`
	}{
		Data: []*response.AccessDTO{
			{
				ExternalID: "95729",
				Run:        "21480585-5",
				FullName:   "William Gabriel Araya Pino",
				Location:   "104",
				EntryAt:    "2025-10-05T16:59:48Z",
				ExitAt:     &exitAtValue,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, accesses)
}

func TestGetCompleteAccess_UnmarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, accesses)
}

func TestGetCompleteAccess_EmptyResponse(t *testing.T) {
	mockResponse := struct {
		Data []*response.AccessDTO `json:"data"`
	}{
		Data: []*response.AccessDTO{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Setup service with server URL
	mockRepo := new(MockTrackRepository)
	mockNotifyService := new(MockNotificationService)
	envConfig := &config.EnvironmentConfig{
		AccessServiceBaseUrl: server.URL,
	}
	service := NewAccessServiceImpl(mockRepo, mockNotifyService, envConfig)

	// Execute
	accesses, err := service.GetCompleteAccess()

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, accesses)
	assert.Len(t, accesses, 0)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
