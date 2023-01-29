package routes

import (
	gopkgs "github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/VulEQ/log"
)

func MainRouter(app fiber.Router) {
	api := app.Group("/v0", gopkgs.Auth)

	SonarProjectRouter(api)
	SonarScanRouter(api)

	log.Log.Info("All routes created successfully :)")

}
