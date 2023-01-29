package repositories

import (

	"github.com/abr-ooo/VulEQ/connections"
	"github.com/abr-ooo/VulEQ/log"
	"github.com/abr-ooo/VulEQ/models"
	"github.com/sirupsen/logrus"
)

func CreateProject(p *models.Project) error {

	query := "INSERT INTO project(del_project_id, git_branch, git_url, git_token, user_id, date, name, key, token, sonar_token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"

	result, err := connections.DB.Query(query, p.DelProjectID, p.GitBranch, p.GitURL, p.GitToken, p.UserID, p.Date, p.Name, p.Key, p.Token, p.SonarToken)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project": p,
			"error":   err.Error(),
		}).Debug("Repo. Execution *CreateProject* query in DB have error!")
		return err
	}

	var id int
	for result.Next() {
		err = result.Scan(&id)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project": p,
				"id":      id,
				"error":   err.Error(),
			}).Debug("Repo. Scan result of *CreateProject* to id query have error!")
			return err
		}
	}

	p.ID = id

	log.Log.WithField("project_id", p.ID).Debug("Repo. Create Project Finish :))")
	return nil
}

func GetProject(value string, field string) (*models.Project, error) {

	project := new(models.Project)

	query := "SELECT * FROM project WHERE " + field + "=$1"
	result, err := connections.DB.Query(query, value)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"filed":       field,
			"field_value": value,
			"error":       err.Error(),
		}).Debug("Repo. Execution *GetProject* query in DB have error!")
		return nil, err
	}

	for result.Next() {
		err = result.Scan(&project.ID, &project.DelProjectID, &project.GitBranch,
			&project.GitURL, &project.GitToken, &project.UserID, &project.Date, &project.Name, &project.Key,
			&project.Token, &project.SonarToken, &project.ScanCounter)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project":     project,
				"filed_value": value,
				"field":       field,
				"error":       err.Error(),
			}).Debug("Repo. Scan result of *GetProject* to project model query have error!")
			return nil, err
		}
	}

	log.Log.WithField("project", project).Debug("Repo. Get Project Finish :))")
	return project, nil
}

func UpdateProject(projectID int, value interface{}, field string) error {
	query := "UPDATE project SET " + field + "=$1 WHERE id=$2"

	_, err := connections.DB.Query(query, value, projectID)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project_id":   projectID,
			"update_field": field,
			"new_value":    value,
			"error":        err.Error(),
		}).Debug("Repo. Execution *UpdateProject* query in DB have error!")
		return err
	}
	log.Log.Debug("Repo. Update Project Finish :))")

	return nil
}

func ProjectExist(value string, field string) (bool, error) {

	query := "SELECT * FROM project WHERE " + field + "=$1"
	result, err := connections.DB.Query(query, value)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"field": field,
			"value": value,
			"error": err.Error(),
		}).Debug("Repo. Execution *ExistProject* query in DB have error!")
		return false, err
	}
	exist := result.Next()

	log.Log.WithFields(logrus.Fields{
		"exist_project": exist,
		"field":         field,
		"value":         value,
	}).Debug("Repo. Exist Project Finish :))")
	return exist, nil

}

func VerifyUserProject(userID int, field string, projectID int) (bool, error) {
	query := "SELECT * FROM project WHERE user_id=$1 and " + field + "=$2"
	result, err := connections.DB.Query(query, userID, projectID)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"user_id":    userID,
			"field":      field,
			"project_id": projectID,
			"error":      err.Error(),
		}).Debug("Repo. Execution *VerifyUserProject* query in DB have error!")
		return false, err
	}
	verify := result.Next()

	log.Log.WithFields(logrus.Fields{
		"verify":     verify,
		"user_id":    userID,
		"field":      field,
		"project_id": projectID,
		}).Debug("Repo.Verify User Project Finish :))")
	return verify, nil
}
