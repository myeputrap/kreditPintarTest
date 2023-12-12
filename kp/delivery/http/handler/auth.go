package handler

import (
	"goKreditPintar/domain"
	"goKreditPintar/helper"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// AuthHandler is handler for authentication
type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

// PostLogin is handler for login
func (ah *AuthHandler) PostLogin(c *fiber.Ctx) (err error) {
	var input domain.LoginRequest
	err = c.BodyParser(&input)
	if err != nil {
		log.Errorf("error bodyparser PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	input.Username = c.FormValue("user_name")
	input.PhoneNumber = c.FormValue("phone_number")
	isAdminStr := c.FormValue("is_admin")
	if isAdminStr == "true" {
		input.IsAdmin = true
	} else {
		input.IsAdmin = false
	}
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Errorf("error validator PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	res, err := ah.AuthUsecase.PostLogin(c.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(fasthttp.StatusBadRequest).SendString("The data you filled is incorrect.")
		}
		if strings.Contains(err.Error(), "phone number is not valid") {
			return c.Status(fasthttp.StatusBadRequest).SendString("phone number is not valid")
		}
		log.Errorf("error  PostLogin %s", err.Error())
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	err = c.JSON(res)

	return err
}

func (ah *AuthHandler) DeleteLogout(c *fiber.Ctx) (err error) {
	var token string
	auth := c.GetReqHeaders()
	authorization := auth["Authorization"]
	if len(authorization) != 0 {
		token = authorization[0][7:]
	}

	if token == "" {
		return c.Status(401).SendString("token not found")
	}
	err = ah.AuthUsecase.PostLogout(c.Context(), token)
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.JSON("success logout")

	return err
}
