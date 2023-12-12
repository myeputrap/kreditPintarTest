package http

import (
	"goKreditPintar/domain"
	"goKreditPintar/kp/delivery/http/handler"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func RouterAPI(app *fiber.App, auth domain.AuthUsecase, authAction domain.ActionUsecase) {
	handlerAuthorization := &handler.AuthorizationHandler{AuthUsecase: auth}
	handlerAction := &handler.ActionHandler{ActionUsecase: authAction}
	handlerAuth := &handler.AuthHandler{AuthUsecase: auth}

	basePath := viper.GetString("server.base_path")
	aut := app.Group(basePath)

	aut.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))
	//auth
	//aut.Get("/auth", handlerAuthorization.MiddlewareAdmin(), handlerAuthorization.GetAuth)
	aut.Get("/consumer", handlerAuthorization.MiddlewareAdmin(), handlerAction.GetConsumer)
	aut.Post("/consumer", handlerAuthorization.MiddlewareAdmin(), handlerAction.PostConsumer)
	aut.Get("/consumer/:id", handlerAuthorization.MiddlewareAdmin(), handlerAction.GetConsumerDetail)
	aut.Delete("/logout", handlerAuth.DeleteLogout)

	aut.Post("/credit-card", handlerAuthorization.MiddlewareAdmin(), handlerAction.PostCreditCards)
	aut.Get("/credit-card", handlerAuthorization.MiddlewareAdmin(), handlerAction.GetCreditCards)
	aut.Get("/credit-card/:id", handlerAuthorization.MiddlewareAdmin(), handlerAction.GetCreditCardDetail)
	aut.Post("/transaction/credit", handlerAuthorization.MiddlewareConsumer(), handlerAction.PostTransactionDetail)
	aut.Patch("/billing/:id", handlerAuthorization.MiddlewareConsumer(), handlerAction.PatchBilling)
	aut.Get("/billing", handlerAuthorization.MiddlewareConsumer(), handlerAction.GetBilling)
	log.Debug(viper.GetInt("rate_limit.limiter_max"))
	log.Debug(time.Duration(viper.GetInt("rate_limit.limiter_expiration")) * time.Second)
	limiterConfig := limiter.Config{
		Max:        viper.GetInt("rate_limit.limiter_max"),
		Expiration: time.Duration(viper.GetInt("rate_limit.limiter_expiration")) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fasthttp.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}
	limiterFiber := limiter.New(limiterConfig)
	limiterLogin := app.Group(basePath + "")
	limiterLogin.Use(limiterFiber)
	aut.Post("/login", handlerAuth.PostLogin)
}
