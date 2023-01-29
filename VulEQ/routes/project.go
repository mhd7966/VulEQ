package routes

import (
	"github.com/abr-ooo/VulEQ/controllers"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)


func SonarProjectRouter(app fiber.Router) {

	project := app.Group("/project")
	project.Get("/:del_project_id/scans", controllers.GetAllScans)
	project.Get("/:del_project_id/issues", controllers.GetIssues)

	log.Info("Project routes created successfully :)")
}
