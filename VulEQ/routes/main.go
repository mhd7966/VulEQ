package routes

import (

	"github.com/gofiber/fiber/v2"
	"github.com/abr-ooo/go-pkgs"
	"github.com/abr-ooo/VulEQ/log"
)

func MainRouter(app fiber.Router) {
	api := app.Group("/v0", gopkgs.Auth)

	SonarProjectRouter(api)
	SonarScanRouter(api)

	log.Log.Info("All routes created successfully :)")

}

