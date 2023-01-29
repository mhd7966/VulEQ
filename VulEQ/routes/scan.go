package routes

import (
	"github.com/abr-ooo/VulEQ/controllers"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func SonarScanRouter(app fiber.Router) {

	scan := app.Group("/scan")

	scan.Post("", controllers.SonarScan)
	scan.Get("/:scan_id", controllers.ScanInfo)

	log.Info("Scan routes created successfully :)")
}
