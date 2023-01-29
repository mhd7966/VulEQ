package controllers

import (
	"strconv"
	"strings"

	gopkgs "github.com/abr-ooo/go-pkgs"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/VulEQ/log"
	"github.com/mhd7966/VulEQ/models"
	"github.com/mhd7966/VulEQ/repositories"
	"github.com/mhd7966/VulEQ/services/sonarqube"
	"github.com/sirupsen/logrus"
)

// GetScans godoc
// @Summary get all scans
// @Description return all scans info of a project(by del_project_id)
// @ID get_all_scans_by_delProjectID
// @Param del_project_id path int true "del_project_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Scan
// @Failure 400 json httputil.HTTPError
// @Router /project/{del_project_id}/scans [get]
func GetAllScans(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("GetAllScans. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("GetAllScans. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	delProjectID, err := c.ParamsInt("del_project_id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"delProjectID": delProjectID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("GetAllScans. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	exist, err := repositories.ProjectExist(strconv.Itoa(delProjectID), "del_project_id")
	if err != nil {
		response.Message = "Check Project Failed!"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
			"error":          err.Error(),
		}).Error("GetAllScans. Check exist del_project_id in project DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !exist {
		response.Message = "This project does not exist!"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
		}).Info("GetAllScans. This del_project_id doesn't exist!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserProject, err := VerifyUserProject(userID, "del_project_id", delProjectID)
	if err != nil {
		response.Message = "Check Access User To Project Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":        userID,
			"del_project_id": delProjectID,
			"verify":         verifyUserProject,
			"error":          err.Error(),
		}).Error("GetAllScans. Verify access user to project failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserProject {
		response.Message = "Access User To Project Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":        userID,
			"del_project_id": delProjectID,
			"verify":         verifyUserProject,
		}).Info("GetAllScans. User doesn't have access to project!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	scans, err := repositories.GetAllScan(delProjectID)
	if err != nil {
		response.Message = "Get Scans Failed!"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
			"error":          err.Error(),
		}).Error("GetAllScans. Get all scans from scan DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = scans
	log.Log.WithFields(logrus.Fields{
		"del_project_id": delProjectID,
		"scans":          scans,
	}).Info("GetAllScans. Get all scans succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// GetIssues godoc
// @Summary get all issues
// @Description return all issues of a project(by del_project_id)
// @ID get_all_issues_by_delProjectID
// @Param del_project_id path int true "del_project_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Issue
// @Failure 400 json httputil.HTTPError
// @Router /project/{del_project_id}/issues [get]
func GetIssues(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("GetIssues. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("GetIssues. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	delProjectID, err := c.ParamsInt("del_project_id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"delProjectID": delProjectID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("GetIssues. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	exist, err := repositories.ProjectExist(strconv.Itoa(delProjectID), "del_project_id")
	if err != nil {
		response.Message = "Check Project Failed"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
			"error":          err.Error(),
		}).Error("GetIssues. Check exist del_project_id in project DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !exist {
		response.Message = "This project does not exist"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
		}).Info("GetIssues. This del_project_id doesn't exist!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	verifyUserProject, err := VerifyUserProject(userID, "del_project_id", delProjectID)
	if err != nil {
		response.Message = "Check Access User To Project Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id":        userID,
			"del_project_id": delProjectID,
			"verify":         verifyUserProject,
			"error":          err.Error(),
		}).Error("GetIssues. Verify access user to project failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserProject {
		response.Message = "Access User To Project Doesn't Verify"
		log.Log.WithFields(logrus.Fields{
			"user_id":        userID,
			"del_project_id": delProjectID,
			"verify":         verifyUserProject,
		}).Info("GetIssues. User doesn't have access to project!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	project, err := repositories.GetProject(strconv.Itoa(delProjectID), "del_project_id")
	if err != nil {
		response.Message = "Check Project Failed"
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
			"error":          err.Error(),
		}).Error("GetIssues. Get project from project DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var issues []models.Issue

	issues, err = sonarqube.GetIssues(project)
	var issuesResp = make([]models.ReturnIssueResponse, len(issues))
	if err != nil {
		response.Message = err.Error()
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
			"project":        project,
			"error":          err.Error(),
		}).Error("GetIssues. Get issues from sonarqube DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	} else {
		for i := 0; i < len(issues); i++ {
			issuesResp[i].Component = issues[i].Component
			issuesResp[i].Line = issues[i].Line
			issuesResp[i].Message = issues[i].Message
			issuesResp[i].Severity = issues[i].Severity
			issuesResp[i].Type = issues[i].Type

		}
	}

	for i := 0; i < len(issuesResp); i++ {
		issuesResp[i].Component = strings.ReplaceAll(issuesResp[i].Component, project.Key, "source_code")
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = issuesResp
	log.Log.WithFields(logrus.Fields{
		"del_project_id": delProjectID,
		"issues":         issuesResp,
	}).Info("GetIssues. Get all issues succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

func GetUser(c *fiber.Ctx) string {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	user_id := (claims["sub"]).(string)
	return user_id

}

func VerifyUser(user_id int) (bool, int, error) {
	exist, err := repositories.ProjectExist(strconv.FormatUint(uint64(user_id), 10), "user_id")
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"user_id": user_id,
			"error":   err.Error(),
		}).Info("Verify User. check exist user failed!")
		return false, 0, err
	}
	log.Log.WithFields(logrus.Fields{
		"user_id": user_id,
		"exist":   exist,
	}).Info("VerifyUser. success!")
	return exist, user_id, nil
}

// func VerifyUserScan(user_id string, del_project_id string, scan_id string) (bool, error) {
// 	verify, err := repositories.VerifyUserDomain(user_id, del_project_id, scan_id)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"user_id": user_id,
// 			"scan_id": scan_id,
// 			"error":   err.Error(),
// 		}).Info("VerifyUserScan. check access user to scan failed!")
// 		return false, err
// 	}
// 	log.WithFields(log.Fields{
// 		"user_id": user_id,
// 		"scan_id": scan_id,
// 		"verify":  verify,
// 	}).Info("VerifyUserScan. success!")
// 	return verify, nil
// }

func VerifyUserProject(user_id int, field string, project_id int) (bool, error) {
	verify, err := repositories.VerifyUserProject(user_id, field, project_id)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"user_id":    user_id,
			"field":      field,
			"project_id": project_id,
			"error":      err.Error(),
		}).Info("VerifyUserProject. check access user to project failed!")
		return false, err
	}

	log.Log.WithFields(logrus.Fields{
		"user_id":    user_id,
		"field":      field,
		"project_id": project_id,
		"verify":     verify,
	}).Info("VerifyUserProject. success!")
	return verify, nil
}
