package controllers

import (
	"strconv"

	"github.com/abr-ooo/VulEQ/connections"
	"github.com/abr-ooo/VulEQ/jobs"
	"github.com/abr-ooo/VulEQ/models"
	"github.com/abr-ooo/VulEQ/repositories"
	gopkgs "github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/abr-ooo/VulEQ/log"
)

// ScanProject godoc
// @Summary scan project
// @Description if project doesn't exist, first create that then scan it else just clone and scan
// @ID scan_by_del_project_id
// @Accept  json
// @Produce  json
// @Param projectBody body models.ProjectBody true "Project info"
// @Security ApiKeyAuth
// @Success 200 {object} models.Scan
// @Failure 400 json httputil.HTTPError
// @Router /scan [post]
func SonarScan(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	projectModel, scan, err := GetProjectModel(c)
	if err != nil {
		response.Message = "Parse body Failed"
		log.Log.WithFields(logrus.Fields{
			"project_model": projectModel,
			"scan":          scan,
			"error":         err.Error(),
		}).Error("Scan. Parse request body to project and scan failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectModel.UserID = int(gopkgs.UID(c))
	exist, err := repositories.ProjectExist(strconv.Itoa(projectModel.DelProjectID), "del_project_id")
	if err != nil {
		response.Message = "Check Project Failed"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": projectModel.DelProjectID,
			"error":          err.Error(),
		}).Error("Scan. Check exist this del_project_id failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if exist {

		project, err := repositories.GetProject(strconv.Itoa(projectModel.DelProjectID), "del_project_id")
		if err != nil {
			response.Message = "Get Project Failed"
			log.Log.WithFields(logrus.Fields{
				"del_project_id": projectModel.DelProjectID,
				"project":        project,
				"error":          err.Error(),
			}).Error("Scan1. Get project from project DB failed!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		err = repositories.CreateScan(scan, project.ID)
		if err != nil {
			response.Message = "Create Scan Record Failed"
			log.Log.WithFields(logrus.Fields{
				"project_id": projectModel.ID,
				"scan":       scan,
				"error":      err.Error(),
			}).Error("Scan1. Create scan record in DB failed!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		task, err := jobs.NewScanProjectTask(scan.ID)
		if err != nil {
			response.Message = "Scan - could not create task"
			log.Log.WithFields(logrus.Fields{
				"scan_id": scan.ID,
				"error":   err.Error(),
			}).Error("Scan1. Couldn't create new scan task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		_, err = connections.RedisClient.Enqueue(task)
		if err != nil {
			response.Message = "Scan - could not enqueue task"
			log.Log.WithFields(logrus.Fields{
				"scan_id": scan.ID,
				"task":    task,
				"error":   err.Error(),
			}).Error("Scan1. Couldn't enqueue scan task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response.Message = "OK!"
		response.Status = "succes"
		response.Data = scan.ID
		log.Log.WithField("scan_id", scan.ID).Info("Scan1. Create new scan task worker succuesfully :)")
		return c.Status(fiber.StatusOK).JSON(response)

	} else {

		jobs.SetProjectInfo(projectModel)
		log.Log.Info("Projct info set to projectModel succesfully :)")

		err = repositories.CreateProject(projectModel)
		if err != nil {
			response.Message = "Create Project Failed"
			log.Log.WithFields(logrus.Fields{
				"project_model": projectModel,
				"error":         err.Error(),
			}).Error("Scan2. Create project record in DB failed!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		err = repositories.CreateScan(scan, projectModel.ID)
		if err != nil {
			response.Message = "Create Scan Failed!"
			log.Log.WithFields(logrus.Fields{
				"project_id": projectModel.ID,
				"scan":       scan,
				"error":      err.Error(),
			}).Error("Scan2. Create scan record in DB failed!")
			return c.Status(fiber.StatusBadRequest).JSON(response)

		}

		task, err := jobs.NewCreateProjectTask(projectModel.ID, scan.ID)
		if err != nil {
			response.Message = "Project - could not create project task"
			log.Log.WithFields(logrus.Fields{
				"project_id": projectModel.ID,
				"scan_id":    scan.ID,
				"error":      err.Error(),
			}).Error("Scan2. Couldn't create new project task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		_, err = connections.RedisClient.Enqueue(task)
		if err != nil {
			response.Message = "Project - could not enqueue task"
			log.Log.WithFields(logrus.Fields{
				"project_id": projectModel.ID,
				"scan_id":    scan.ID,
				"task":       task,
				"error":      err.Error(),
			}).Error("Scan2. Couldn't enqueue project task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response.Message = "OK!"
		response.Status = "succes"
		response.Data = scan.ID
		log.Log.WithField("scan_id", scan.ID).Info("Scan2. Create new project task worker succuesfully :)")
		return c.Status(fiber.StatusOK).JSON(response)

	}

}

// GetScan godoc
// @Summary get scan information
// @Description get all information about a specific scan of project
// @ID get_info_scan_by_scan_id
// @Param scan_id path int true "ID that return after scan a project"
// @Security ApiKeyAuth
// @Success 200 {object} models.Scan
// @Failure 400 json httputil.HTTPError
// @Router /scan/{scan_id} [get]
func ScanInfo(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("ScanInfo. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("ScanInfo. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	scanID, err := c.ParamsInt("scan_id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"scan_id":  scanID,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("ScanInfo. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	exist, err := repositories.ScanExist(scanID)
	if err != nil {
		response.Message = "Check Scan Failed"
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
			"error":   err.Error(),
		}).Error("ScanInfo. Check exist this scan_id in scan DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !exist {
		response.Message = "This scan does not exist"
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
		}).Info("ScanInfo. This scan_id doesn't exist!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	scan, err := repositories.GetScan(scanID)
	if err != nil {
		response.Message = "Get Scan Failed"
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
			"error":   err.Error(),
		}).Error("ScanInfo. Get scan info from scan DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserProject, err := VerifyUserProject(userID, "id", scan.ProjectID)
	if err != nil {
		response.Message = "Check Access User To Scan Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"scan_id": scanID,
			"verify":  verifyUserProject,
			"error":   err.Error(),
		}).Error("ScanInfo. Verify access user to scan failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserProject {
		response.Message = "Access User To Scan Doesn't Verify"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"scan_id": scanID,
			"verify":  verifyUserProject,
		}).Info("ScanInfo. User doesn't have access to scan!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = scan
	log.Log.WithFields(logrus.Fields{
		"scan_id": scanID,
		"scan":    scan,
	}).Info("GetIssues. Get scan info succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// func CreateProject(c *fiber.Ctx) error {

// 	projectBody := new(models.Project)

// 	err := c.BodyParser(projectBody)
// 	if err != nil {
// 		return c.JSON(err)
// 	}

// 	projectBody.Date = time.Now().Format("2006-01-02|15:04:05")
// 	lastIndex := strings.LastIndex(projectBody.GitURL, "/")
// 	projectBody.Name = projectBody.GitURL[lastIndex+1 : len(projectBody.GitURL)-4]

// 	repositories.CreateProject(projectBody)

// 	err = git.CloneCode(*projectBody)
// 	if err != nil {
// 		return c.JSON(err)
// 	}

// 	return c.JSON("success")
// }

func GetProjectModel(c *fiber.Ctx) (*models.Project, *models.Scan, error) {

	projectBody := new(models.ProjectBody)

	err := c.BodyParser(projectBody)
	if err != nil {
		log.Log.Info("Parse request body to models.ProjectBody failed !")
		return nil, nil, err
	}

	projectModel := new(models.Project)
	projectModel.DelProjectID = projectBody.DelProjectID
	projectModel.GitBranch = projectBody.GitBranch
	projectModel.GitToken = projectBody.GitToken
	projectModel.GitURL = projectBody.GitURL
	log.Log.Debug("Set project model fields :)")

	scan := new(models.Scan)
	scan.GitCommitHash = projectBody.GitCommitHash
	scan.PipelineID = projectBody.PipelineID
	log.Log.Debug("Set scan model fields :)")

	return projectModel, scan, nil
}
