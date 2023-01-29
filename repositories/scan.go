package repositories

import (
	"github.com/mhd7966/VulEQ/connections"
	"github.com/mhd7966/VulEQ/log"
	"github.com/mhd7966/VulEQ/models"
	"github.com/sirupsen/logrus"
)

func CreateScan(scan *models.Scan, projectID int) error {

	query := "INSERT INTO scan(project_id, git_commit_hash, pipeline_id) VALUES ($1, $2, $3) RETURNING id"

	result, err := connections.DB.Query(query, projectID, scan.GitCommitHash, scan.PipelineID)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan":            scan,
			"project_id":      projectID,
			"git_commit_hash": scan.GitCommitHash,
			"pipeline_id":     scan.PipelineID,
			"error":           err.Error(),
		}).Debug("Repo. Execution *CreateScan* query in DB have error!")
		return err
	}

	var id int

	for result.Next() {
		err = result.Scan(&id)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"id":    id,
				"error": err.Error(),
			}).Debug("Repo. Scan result of *CreateScan* to id query have error!")
			return err
		}

	}

	scan.ID = id

	log.Log.WithField("scan_id", scan.ID).Debug("Repo. Create Scan Finish :))")
	return nil
}

func UpdateScan(s *models.Scan) error {

	query := "UPDATE scan SET num_bug=$1, num_vulnerability=$2, num_debt=$3, num_code_smell=$4, num_file=$5, num_duplicate_line=$6, line_code=$7, line_comment=$8 Where id=$9"

	_, err := connections.DB.Query(query, s.NumBug, s.NumVulnerability, s.NumDebt, s.NumCodeSmell, s.NumFile, s.NumDuplicateLine, s.LineCode, s.LineComment, s.ID)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan_id": s,
			"error":   err.Error(),
		}).Debug("Repo. Execution *UpdateScan* query in DB have error!")
		return err
	}

	log.Log.Debug("Repo. Update Scan Finish :))")
	return nil
}

func GetScan(scanID int) (*models.Scan, error) {

	scan := new(models.Scan)

	query := "SELECT * FROM scan WHERE id=$1"
	result, err := connections.DB.Query(query, scanID)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
			"error":   err.Error(),
		}).Debug("Repo. Execution *GetScan* query in DB have error!")
		return nil, err
	}

	for result.Next() {
		err = result.Scan(&scan.ID, &scan.ProjectID, &scan.GitCommitHash, &scan.PipelineID, &scan.NumBug,
			&scan.NumVulnerability, &scan.NumDebt, &scan.NumCodeSmell, &scan.NumFile,
			&scan.NumDuplicateLine, &scan.LineCode, &scan.LineComment)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"scan":  scan,
				"error": err.Error(),
			}).Debug("Repo. Scan result of *GetScan* to scan model query have error!")
			return nil, err
		}
	}

	log.Log.WithField("scan", scan).Debug("Repo. Get Scan Finish :))")
	return scan, nil
}

func GetAllScan(delProjectID int) ([]models.Scan, error) {

	scan := new(models.Scan)

	query := "SELECT * FROM scan WHERE project_id = (SELECT id FROM project WHERE del_project_id = $1)"

	result, err := connections.DB.Query(query, delProjectID)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"del_project_id": delProjectID,
			"error":          err.Error(),
		}).Debug("Repo. Execution *GetAllScan* query in DB have error!")
		return nil, err
	}

	var scans []models.Scan

	for result.Next() {
		err = result.Scan(&scan.ID, &scan.ProjectID, &scan.GitCommitHash, &scan.PipelineID, &scan.NumBug,
			&scan.NumVulnerability, &scan.NumDebt, &scan.NumCodeSmell, &scan.NumFile,
			&scan.NumDuplicateLine, &scan.LineCode, &scan.LineComment)

		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"scan":  scan,
				"error": err.Error(),
			}).Debug("Repo. Scan result of *GetAllScan* to scan model query have error!")
			return nil, err
		}

		scans = append(scans, *scan)

	}
	log.Log.WithField("scans", scans).Debug("Repo. Get All Scan Finish :))")
	return scans, nil
}

func ScanExist(scanID int) (bool, error) {

	query := "SELECT * FROM scan WHERE id=$1"
	result, err := connections.DB.Query(query, scanID)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
			"error":   err.Error(),
		}).Debug("Repo. Execution *ExistScan* query in DB have error!")
		return false, err
	}
	exist := result.Next()

	log.Log.WithField("exist_scan", exist).Debug("Repo. Exist Scan Finish :))")
	return exist, nil

}

func RemoveScan(scanID int) error {
	query := "DELETE FROM scan WHERE id=$1"
	_, err := connections.DB.Query(query, scanID)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"scan_id": scanID,
			"error":   err.Error(),
		}).Debug("Repo. Execution *DeleteScan* query in DB have error!")
		return err
	}

	log.Log.Debug("Repo. Delete Scan Finish :))")
	return nil
}

// func VerifyUserScan(user_id string, del_project_id string, scan_id string) (bool, error){
// 	query := "SELECT * FROM scan WHERE project_id = (SELECT CAST(@(SELECT id FROM project WHERE user_id = $1 and del_project_id = $2) as varchar(10))) and id=$3"
// 	result, err := connections.DB.Query(query, user_id, del_project_id, scan_id)

// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"user_id": user_id,
// 			"scan_id": scan_id,
// 			"error": err.Error(),
// 		}).Debug("Execution *VerifyUserScan* query in DB have error!")
// 		return false, err
// 	}
// 	verify := result.Next()

// 	log.WithFields(log.Fields{
// 		"verify": verify,
// 		"user_id": user_id,
// 		"scan_id": scan_id,
// 	}).Debug("Repo.Verify User Scan Finish :))")
// 	return verify, nil
// }
