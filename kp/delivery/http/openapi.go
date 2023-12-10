package http

import (
	"embed"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/utils"
	"github.com/spf13/viper"
)

//go:embed openapi/*
var openAPI embed.FS

// RouterOpenAPI is Open API Specification UI for Service Structure System
func RouterOpenAPI(app *fiber.App) {
	app.Use(viper.GetString("server.base_path")+"/openapi", func(c *fiber.Ctx) error {
		originalURL := utils.ImmutableString(c.OriginalURL())

		// Check if the client is requesting a file extension
		extMatch, _ := regexp.MatchString("\\.[a-zA-Z0-9]+$", originalURL)

		if !strings.HasSuffix(originalURL, "/") && !extMatch {
			return c.Redirect(originalURL + "/")
		}
		return c.Next()
	}, filesystem.New(filesystem.Config{
		Root:       http.FS(openAPI),
		PathPrefix: "openapi",
	}))
}
