package controller

import (
	"fmt"
	"spl-notification/internal/dto/request"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
	"spl-notification/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TrackController struct {
	trackService        service.TrackService
	sourceService       service.SourceService
	notificationService service.NotificationService
	validation          *validator.Validate
}

func NewTrackController(
	trackService service.TrackService,
	sourceService service.SourceService,
	notificationService service.NotificationService,
	validation *validator.Validate,
) *TrackController {
	return &TrackController{
		trackService:        trackService,
		sourceService:       sourceService,
		notificationService: notificationService,
		validation:          validation,
	}
}

func (t *TrackController) GetAllFollowTracks(c *fiber.Ctx) error {
	chatId := c.Params("chatId")
	if chatId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "chatId parameter is required",
		})
	}

	tracks, err := t.trackService.GetFollowTracksByChatId(chatId)
	if err != nil {
		return errors.InternalError(c, err)
	}

	if len(tracks) == 0 {
		tracks = []*model.Track{}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tracks,
	})
}

func (t *TrackController) SendAllFollowTracks(c *fiber.Ctx) error {
	chatId := c.Params("chatId")
	if chatId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "chatId parameter is required",
		})
	}

	err := t.trackService.SendAllFollows(chatId)
	if err != nil {
		return errors.InternalError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (t *TrackController) CreateTrack(c *fiber.Ctx) error {
	var createTrackDto request.CreateTrackDTO
	if err := c.BodyParser(&createTrackDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := t.validation.Struct(createTrackDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	abmUser, err := t.sourceService.GetABMUserByRun(createTrackDto.Run)
	if err != nil {
		return errors.InternalError(c, err)
	}

	if abmUser == nil {
		err := t.notificationService.SendMessage(createTrackDto.ChatID, "Usuario no existente")
		if err != nil {
			return errors.InternalError(c, err)
		}
		return c.SendStatus(fiber.StatusNotFound)
	}

	createTrackDto.FullName = fmt.Sprintf("%s %s", abmUser.FirstName, abmUser.LastName)
	createTrackDto.ExternalID = abmUser.ExternalID

	err = t.trackService.Create(&createTrackDto)
	if err != nil {
		return errors.InternalError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (t *TrackController) DeleteTrack(c *fiber.Ctx) error {
	var deleteTrackDto request.DeleteTrackDTO
	if err := c.BodyParser(&deleteTrackDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := t.validation.Struct(deleteTrackDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := t.trackService.Delete(&deleteTrackDto)
	if err != nil {
		return errors.InternalError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
