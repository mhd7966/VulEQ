package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/VulEQ/controllers"
	log "github.com/sirupsen/logrus"
)

func SonarScanRouter(app fiber.Router) {

	scan := app.Group("/scan")

	scan.Post("", controllers.SonarScan)
	scan.Get("/:scan_id", controllers.ScanInfo)

	log.Info("Scan routes created successfully :)")
}
