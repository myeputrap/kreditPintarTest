package http

import (
	"goKreditPintar/domain"
	"goKreditPintar/kp/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

// RouterAPI is main router for this Service Insurance Auth REST API
func RouterAPI(app *fiber.App, auth domain.AuthUsecase) {
	handlerAuthorization := &handler.AuthorizationHandler{AuthUsecase: auth}
	handlerAuth := &handler.AuthHandler{AuthUsecase: auth}

	// limiterFiber := limiter.New(limiter.Config{
	// 	Max:        viper.GetInt("limiter_max"),
	// 	Expiration: time.Duration(viper.GetInt("limiter_expiration")) * time.Second,
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		clientIP, err := helper.GetClientIP(c)
	// 		if err != nil {
	// 			log.Errorf("Error getting client IP: %s", err)
	// 		}
	// 		log.Infof("| Client IP: %s |", clientIP)
	// 		return clientIP
	// 	},
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.SendStatus(fasthttp.StatusTooManyRequests)
	// 	},
	// })

	basePath := viper.GetString("server.base_path")
	aut := app.Group(basePath)

	aut.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// Auth
	aut.Post("/login", handlerAuth.PostLogin)
	aut.Post("/validate", handlerAuth.PostValidate)
	aut.Get("/auth", handlerAuthorization.MiddlewareAuthorization(), handlerAuthorization.GetAuth)
}
