package errors

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Component *string // (ej: "TrackRepository", "TrackService")
	Type      *string // Tipo/código del error
	Err       error   // Error original para wrapping
}

func NewAppError(component string, err error) *AppError {
	return &AppError{
		Component: &component,
		Err:       err,
	}
}

func NewAppErrorWithType(component string, errType string, err error) *AppError {
	return &AppError{
		Component: &component,
		Type:      &errType,
		Err:       err,
	}
}

func (e *AppError) Error() string {
	errStr := ""
	if e.Component != nil {
		errStr += fmt.Sprintf("| %s", *e.Component)
	}

	if e.Type != nil {
		errStr += fmt.Sprintf(":%s |", *e.Type)
	} else {
		errStr += " |"
	}

	return fmt.Sprintf("%s %v", errStr, e.Err)
}

func InternalError(c *fiber.Ctx, err error) error {
	if se, ok := err.(*AppError); ok {
		log.Printf("%s\n", se.Error())
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Internal Server Error"})
	}

	log.Printf("Internal error occurred: %s\n", err.Error())
	return c.
		Status(fiber.StatusInternalServerError).
		JSON(fiber.Map{"error": "Internal Server Error"})
}
