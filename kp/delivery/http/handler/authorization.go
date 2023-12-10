package handler

import (
	"errors"
	"goKreditPintar/domain"
	"goKreditPintar/helper"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// AuthorizationHandler is authorization handler for middleware
type AuthorizationHandler struct {
	AuthUsecase domain.AuthUsecase
}

func (ah *AuthorizationHandler) authorizationAuth(c *fiber.Ctx) (resAuth domain.Client, err error) {
	var token string
	auth := c.GetReqHeaders()
	authorization := auth["Authorization"]
	if len(authorization) != 0 {
		token = authorization[0][7:]
	}
	if token == "" {
		err = errors.New("token not found")
		return
	}
	c.Locals("token", token)

	resAuth, err = ah.AuthUsecase.Authorize(c.Context(), token)
	if err != nil {
		log.Error(err)
	}

	return
}

// MiddlewareAuthorization is authorization middleware
func (ah *AuthorizationHandler) MiddlewareAuthorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, err := ah.authorizationAuth(c)
		if err != nil {
			log.Error(err)
			return helper.HTTPSimpleResponse(c, fasthttp.StatusUnauthorized)
		}
		err = nil

		c.Locals("session", auth)

		return c.Next()
	}
}

// GetAUth is check and get authentication data
func (ah *AuthorizationHandler) GetAuth(c *fiber.Ctx) (err error) {
	return nil
}
