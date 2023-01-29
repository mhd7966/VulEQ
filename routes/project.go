package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/VulEQ/controllers"
	log "github.com/sirupsen/logrus"
)

func SonarProjectRouter(app fiber.Router) {

	project := app.Group("/project")
	project.Get("/:del_project_id/scans", controllers.GetAllScans)
	project.Get("/:del_project_id/issues", controllers.GetIssues)

	log.Info("Project routes created successfully :)")
}
