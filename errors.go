package stdserver

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func newError(statusCode int, err error, messages ...string) *fiber.Error {
	if err == nil && len(messages) == 0 {
		return fiber.NewError(statusCode)
	}
	message := strings.Join(messages, "; ")
	if err != nil {
		message += ": " + err.Error()
	}
	return fiber.NewError(statusCode, message)
}

func ErrBadRequest(err error, messages ...string) *fiber.Error {
	return newError(fiber.StatusBadRequest, err, messages...)
}

func ErrInternalServerError(err error, messages ...string) *fiber.Error {
	return newError(fiber.StatusInternalServerError, err, messages...)
}

func ErrForbidden(err error, messages ...string) *fiber.Error {
	return newError(fiber.StatusForbidden, err, messages...)
}

func ErrUnauthorized(err error, messages ...string) *fiber.Error {
	return newError(fiber.StatusUnauthorized, err, messages...)
}

func ErrConflict(err error, message ...string) *fiber.Error {
	return newError(fiber.StatusConflict, err, message...)
}
