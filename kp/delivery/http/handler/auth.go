package handler

import (
	"goKreditPintar/domain"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler is handler for authentication
type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

// PostLogin is handler for login
func (ah *AuthHandler) PostLogin(c *fiber.Ctx) (err error) {
	return nil
}

// PostValidate is handler for login validation
func (ah *AuthHandler) PostValidate(c *fiber.Ctx) (err error) {
	return nil
}
