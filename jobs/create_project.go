package jobs

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/abr-ooo/VulEQ/connections"
	"github.com/abr-ooo/VulEQ/repositories"
	"github.com/abr-ooo/VulEQ/log"
	"github.com/abr-ooo/VulEQ/services/sonarqube"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeCreateProject = "project:create"
)

type CreateProjectPayload struct {
	ProjectID int
	ScanID    int
}

func NewCreateProjectTask(projectID int, scanID int) (*asynq.Task, error) {
	payload, err := json.Marshal(CreateProjectPayload{ProjectID: projectID, ScanID: scanID})
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project_id": projectID,
			"scan_id":    scanID,
			"error":      err.Error(),
		}).Error("create payload json for project failed!\n", err)
		return nil, err
	}
	return asynq.NewTask(TypeCreateProject, payload), nil
}

func HandleCreateProjectTask(ctx context.Context, t *asynq.Task) error {
	var p CreateProjectPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Log.Error("Worker. unmarshal failed: projectID=%s!", p.ProjectID)
		return err
	}
	log.Log.Info("Worker. Creating Project: projectID=%s!", p.ProjectID)

	project, err := repositories.GetProject(strconv.FormatUint(uint64(p.ProjectID), 10), "id")
	if err != nil {
		log.Log.Error("Worker. get project from projectID in CreateProject failed!\n", err)
		return err
	}

	err = sonarqube.CreateProject(project)
	if err != nil {
		log.Log.Error("Worker. create project from sonarqube in CreateProject failed!\n", err)
		return err
	}

	err = sonarqube.GenerateToken(project)
	if err != nil {
		log.Log.Error("Worker. generate token from sonarqube in CreateProject failed!\n", err)
		return err
	}

	err = repositories.UpdateProject(project.ID, project.SonarToken, "sonar_token")
	if err != nil {
		log.Log.Error("Worker. update project sonar_token in CreateProject failed!\n", err)
		return err
	}

	task, err := NewScanProjectTask(p.ScanID)
	if err != nil {
		log.Log.Error("Worker. Scan - could not create task: %v!", err)
		return err
	}
	info, err := connections.RedisClient.Enqueue(task)
	if err != nil {
		log.Log.Error("Worker. Scan - could not enqueue task: %v!", err)
		return err
	}
	log.Log.Info("Worker. Scan - enqueued task: id=%s queue=%s!", info.ID, info.Queue)

	return nil
}
