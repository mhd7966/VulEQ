package jobs

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/mhd7966/VulEQ/configs"
	"github.com/mhd7966/VulEQ/log"
	"github.com/mhd7966/VulEQ/models"
	"github.com/mhd7966/VulEQ/repositories"
	"github.com/mhd7966/VulEQ/services/git"
	"github.com/mhd7966/VulEQ/services/sonarqube"
	"github.com/sirupsen/logrus"
)

const (
	TypeScanProject = "project:scan"
)

type ScanProjectPayload struct {
	ScanID int
}

func NewScanProjectTask(scanID int) (*asynq.Task, error) {
	payload, err := json.Marshal(ScanProjectPayload{ScanID: scanID})
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
			"error":   err.Error(),
		}).Error("create payload json for scan failed!\n", err)
		return nil, err
	}
	return asynq.NewTask(TypeScanProject, payload), nil
}

func HandleScanProjectTask(ctx context.Context, t *asynq.Task) error {
	var p ScanProjectPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Log.Error("Worker. json.Unmarshal failed: %v:", err)
		return err
	}
	log.Log.Info("Worker. Scan Project: scan_id=%s", p.ScanID)

	scan, err := repositories.GetScan(p.ScanID)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan_id": p.ScanID,
			"error":   err.Error(),
		}).Error("Worker. get scan from scan_id in ScanProject failed!")
		return err
	}

	project, err := repositories.GetProject(strconv.FormatUint(uint64(scan.ProjectID), 10), "id")
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scna_project_id": scan.ProjectID,
			"error":           err.Error(),
		}).Error("Worker. get project from projectID in ScanProject failed!")
		return err
	}

	scanResult, err := Scan(*scan, *project)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan":    scan,
			"project": project,
			"error":   err.Error(),
		}).Error("Worker. Scan in ScanProject failed!\n", err)
		return err
	}

	log.Log.WithField("scan_result", scanResult).Info("Jobs. Worker job finish successfully :)")
	return nil
}

func Scan(scan models.Scan, project models.Project) (*models.Scan, error) {

	err := git.CloneCode(project, scan.GitCommitHash)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project": project,
			"error":   err.Error(),
		}).Debug("Jobs. Clone code from git failed!")

		err = RemoveDir(project)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project": project,
				"error":   err.Error(),
			}).Debug("Jobs. Remove dir after clone code failed!")
			return nil, err
		}
		return nil, err
	}

	log.Log.Debug("Jobs. Clone code from git successful :)")

	err = RunCommand(project)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project": project,
			"error":   err.Error(),
		}).Debug("Jobs. Run command failed!")

		// err = RemoveDir(project)

		// if err != nil {
		// 	log.WithFields(log.Fields{
		// 		"project": project,
		// 		"error":   err.Error(),
		// 	}).Debug("Jobs. Remove dir after run command failed!")
		// 	return nil, err
		// }
		// return nil, err
	}
	log.Log.Debug("Jobs. Run Command successful :)")

	time.Sleep(3 * time.Second)

	var measures []models.Measure
	measures, err = sonarqube.GetMeasures(&project)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project": project,
			"error":   err.Error(),
		}).Debug("Jobs. get measures from sonar failed!")

		err = RemoveDir(project)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project": project,
				"error":   err.Error(),
			}).Debug("Jobs. Remove dir after get measures form sonar failed!")
			return nil, err
		}

		return nil, err
	}
	log.Log.Debug("Jobs. Get measure from sonar successful :)")

	err = SetScanInfo(&scan, measures, &project)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project":  project,
			"scan":     scan,
			"measures": measures,
			"error":    err.Error(),
		}).Debug("Jobs. Set scan info failed!")

		err = RemoveDir(project)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project": project,
				"error":   err.Error(),
			}).Debug("Jobs. Remove dir after set scan info failed!")
			return nil, err
		}
		return nil, err
	}
	log.Log.Debug("Jobs. Set scan info successful :)")

	//update scan
	err = repositories.UpdateScan(&scan)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan":  scan,
			"error": err.Error(),
		}).Debug("Jobs. Update scan failed!")

		err = RemoveDir(project)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project": project,
				"error":   err.Error(),
			}).Debug("Jobs. Remove dir after update scan failed!")
			return nil, err
		}
		return nil, err
	}
	log.Log.Debug("Jobs. Update scan successful :)")

	project.ScanCounter++
	err = repositories.UpdateProject(project.ID, project.ScanCounter, "scan_counter")
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project": project,
			"error":   err.Error(),
		}).Debug("Jobs. Update scan counter from project failed!")

		err = RemoveDir(project)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"project": project,
				"error":   err.Error(),
			}).Debug("Jobs. Remove dir after update project scan_counter failed!")
			return nil, err
		}
		return nil, err
	}
	log.Log.Debug("Jobs. Update scan_counter from project successful :)")

	err = RemoveDir(project)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"project": project,
			"error":   err.Error(),
		}).Debug("Jobs. Remove dir in last of scan process failed!")

		return nil, err
		// err = repositories.UpdateProject(project.MainID, project.ScanCounter-1, "scan_counter")
		// if err != nil {
		// 	return nil, err
		// }
	}
	log.Log.Debug("Jobs. Remove directory successful :)")
	log.Log.Debug("Jobs. Scan successfully :)")
	return &scan, nil
}

func SetProjectInfo(project *models.Project) error {

	config := configs.Cfg.SonarQube

	project.Date = time.Now().Format("2006-01-02|15:04:05")
	project.Name = uuid.New().String()
	project.Key = config.KeyPrefix + project.Name
	project.Token = config.TokenPrefix + project.Name

	log.Log.Debug("Jobs. Set Project Info successfully :)")
	return nil
}

func SetScanInfo(scan *models.Scan, measures []models.Measure, project *models.Project) error {
	scan.ProjectID = project.ID
	var err error

	for i := 0; i < len(measures); i++ {
		switch measures[i].Metric {
		case "bugs":
			scan.NumBug = measures[i].Value
		case "vulnerabilities":
			scan.NumVulnerability = measures[i].Value
		case "sqale_index":
			scan.NumDebt = measures[i].Value
		case "code_smells":
			scan.NumCodeSmell = measures[i].Value
		case "files":
			scan.NumFile = measures[i].Value
		case "duplicated_lines":
			scan.NumDuplicateLine = measures[i].Value
		case "ncloc":
			scan.LineCode = measures[i].Value
		case "comment_lines":
			scan.LineComment = measures[i].Value

			if err != nil {
				// error nadarim mesle inke :)
				return err
			}
		}
	}

	log.Log.Debug("Jobs. Set Scan Info successfully :)")
	return nil
}

func RunCommand(project models.Project) error {

	config := configs.Cfg

	key := "-Dsonar.projectKey=" + project.Key
	source := "-Dsonar.sources=."
	url := "-Dsonar.host.url=http://" + config.SonarQube.Host
	login := "-Dsonar.login=" + project.SonarToken

	cmd := exec.Command(config.SonarQube.ScannerPath, key, source, url, login)
	cmd.Dir = config.Git.ClonePath + project.Name
	out, err := cmd.Output()
	// log.Debug(out)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"dir":                  cmd.Dir,
			"project_key":          key,
			"sonar_source_project": source,
			"sonar_host":           url,
			"sonar_login,token":    login,
			"out":                  out,
			"error":                err.Error(),
		}).Debug("Jobs. Run Command failed!")
		return err
	} else {
		log.Log.Debug("Jobs. Run Command successfully :)")
		return nil
	}

}

func RemoveDir(project models.Project) error {

	path := configs.Cfg.Git.ClonePath + project.Name

	err := os.RemoveAll(path)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"path":  path,
			"error": err.Error(),
		}).Debug("Jobs. Remove dir failed!")
		return err
	}

	log.Log.Debug("Jobs. Remove dir successfully :)")
	return nil
}
