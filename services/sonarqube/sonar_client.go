package sonarqube

import (
	"encoding/base64"
	"errors"

	"github.com/imroc/req"
	"github.com/mhd7966/VulEQ/configs"
	"github.com/mhd7966/VulEQ/log"
	"github.com/mhd7966/VulEQ/models"
	"github.com/sirupsen/logrus"
)

func Ping() error {

	config := configs.Cfg.SonarQube

	// header := req.Header{
	// 	"Host":          config.Host,
	// 	"Accept":        "application/json",
	// 	"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Auth)),
	// 	// "Cookie":          "XSRF-TOKEN=" + config.XToken + "; PGADMIN_LANGUAGE=en; JWT-SESSION=" + config.JWTSession,
	// 	// "X-XSRF-TOKEN":    config.XToken,
	// }

	url := "http://" + config.Host + "/api/system/ping"
	r, err := req.Post(url)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"URL":    url,
			"Header": "header",
			"error":  err.Error(),
		}).Debug("SONAR. Ping request failed!")
		return err
	}

	if r.Response().StatusCode == 200 {
		log.Log.WithFields(logrus.Fields{
			"response_code": r.Response().StatusCode,
		}).Info("SONAR. PING PONG :))")
		return nil
	} else {

		log.Log.WithFields(logrus.Fields{
			"response_code": r.Response().StatusCode,
		}).Fatal("SONAR. NO PING  :||||")
		return nil
	}
}

func CreateProject(project *models.Project) error {

	config := configs.Cfg.SonarQube

	header := req.Header{
		"Host":            config.Host,
		"Accept":          "application/json",
		"Accept-Language": "en-US,en;q=0.5",
		"Connection":      "keep-alive",
		"Content-Length":  "8",
		"Authorization":   "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Auth)),
		// "Cookie":          "XSRF-TOKEN=" + config.XToken + "; PGADMIN_LANGUAGE=en; JWT-SESSION=" + config.JWTSession,
		"Origin": "http://" + config.Host,
		// "X-XSRF-TOKEN":    config.XToken,
	}
	param := req.Param{
		"project": project.Key,
		"name":    project.Name,
	}

	url := "http://" + config.Host + "/api/projects/create"
	r, err := req.Post(url, header, param)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"URL":       url,
			"Header":    header,
			"Parameter": param,
			"error":     err.Error(),
		}).Debug("SONAR. Create Project sonar request failed!")
		return err
	}

	log.Log.Debug("SONAR. Create Project sonar request ok :)")
	var resp map[string]interface{}
	err = r.ToJSON(&resp)

	if err != nil {
		log.Log.Debug("SONAR. convert Create Project response to map failed!")
		return err
	}

	_, ok := resp["project"]
	if !ok {
		log.Log.WithFields(logrus.Fields{
			"Response": r.String(),
		}).Debug("SONAR. Response of Create Project sonar request is incorrect!")
		return errors.New(r.String())
	}
	log.Log.Debug("SONAR. Create Project sonar request succesfull :)")

	return nil
}

func GenerateToken(project *models.Project) error {

	config := configs.Cfg.SonarQube

	header := req.Header{
		"Host":            config.Host,
		"Accept":          "application/json",
		"Accept-Language": "en-US,en;q=0.5",
		"Connection":      "keep-alive",
		"Content-Length":  "8",
		// "Cookie":          "PGADMIN_LANGUAGE=en; XSRF-TOKEN=" + config.XToken + "; JWT-SESSION=" + config.JWTSession,
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Auth)),
		"Origin":        "http://" + config.Host,
		// "X-XSRF-TOKEN":    config.XToken,
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      "http://" + config.Host + "/dashboard?id=" + project.Key,
	}

	param := req.Param{
		"name": project.Token,
	}

	url := "http://" + config.Host + "/api/user_tokens/generate"
	r, err := req.Post(url, header, param)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"URL":       url,
			"Header":    header,
			"Parameter": param,
			"error":     err.Error(),
		}).Debug("SONAR. Generate Token sonar request failed!")
		return err
	}
	log.Log.Debug("SONAR. Generate Token sonar request ok :)")

	var resp map[string]interface{}
	err = r.ToJSON(&resp)

	if err != nil {
		log.Log.Debug("SONAR. convert Generate Token response to map failed!")
		return err
	}

	_, ok := resp["token"]
	if !ok {
		log.Log.WithFields(logrus.Fields{
			"Response": r.String(),
		}).Debug("SONAR. Response of Generate Token sonar request is incorrect!")
		return errors.New("the content of the generate token response is incorrect")
	} else {
		log.Log.Debug("SONAR. Generate Token sonar request succesfull :)")
		project.SonarToken = resp["token"].(string)
	}

	return nil
}

func GetIssues(project *models.Project) ([]models.Issue, error) {

	config := configs.Cfg.SonarQube

	header := req.Header{
		"Host":            config.Host,
		"Accept":          "application/json",
		"Accept-Language": "en-US,en;q=0.5",
		"Connection":      "keep-alive",
		// "Cookie":          "XSRF-TOKEN=" + config.XToken + "; JWT-SESSION=" + config.JWTSession,
		"Origin": "http://" + config.Host,
		// "X-XSRF-TOKEN":    config.XToken,
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Auth)),
		"Referer":       "http://" + config.Host + "/project/issues?id=" + project.Key + "&resolved=false",
	}

	param := req.QueryParam{
		"componentKeys":    project.Key,
		"s":                "FILE_LINE",
		"resolved":         "false",
		"ps":               "100",
		"facets":           "owaspTop10,sansTop25,severities,sonarsourceSecurity,types",
		"additionalFields": "_all",
		"timeZone":         "Asia/Tehran",
	}
	// only url is required, others are optional.
	url := "http://" + config.Host + "/api/issues/search"
	r, err := req.Post(url, header, param)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"URL":       url,
			"Header":    header,
			"Parameter": param,
			"error":     err.Error(),
		}).Debug("SONAR. Get Issues sonar request Failed!")
		return nil, err
	}
	log.Log.Debug("SONAR. Get Issues sonar request ok :)")

	var resp models.IssuesResponse
	err = r.ToJSON(&resp)

	if err != nil {
		log.Log.Debug("SONAR. convert Get Issues response to IssuesResponse model failed!")
		return nil, err
	}

	log.Log.Info(resp.Issues)
	log.Log.Debug("SONAR. Get Issues sonar request succesfull :)")
	return resp.Issues, nil

}

func GetMeasures(project *models.Project) ([]models.Measure, error) {

	config := configs.Cfg.SonarQube

	header := req.Header{
		"Host":            config.Host,
		"Accept":          "application/json",
		"Accept-Language": "en-US,en;q=0.5",
		"Connection":      "keep-alive",
		// "Cookie":          "XSRF-TOKEN=" + config.XToken + "; JWT-SESSION=" + config.JWTSession,
		"Origin":        "http://" + config.Host,
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Auth)),
		// "X-XSRF-TOKEN":    config.XToken,
		"Referer": "http://" + config.Host + "/component_measures?id=" + project.Key,
	}

	param := req.QueryParam{
		"additionalFields": "period",
		"component":        project.Key,
		"metricKeys":       "bugs,vulnerabilities,sqale_index,code_smells,files,duplicated_lines,ncloc,comment_lines",
	}

	url := "http://" + config.Host + "/api/measures/component"
	r, err := req.Post(url, header, param)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"URL":       url,
			"Header":    header,
			"Parameter": param,
			"error":     err.Error(),
		}).Debug("SONAR. Get Measures sonar request Failed!")
		return nil, err
	}
	log.Log.Debug("SONAR. Get Measures sonar request ok :)")

	var resp models.MeasuresResponse
	err = r.ToJSON(&resp)

	if err != nil {
		log.Log.Debug("SONAR. convert Get Measures response to MeasuresResponse model failed!")
		return nil, err
	}
	log.Log.Info(resp.Component.Measures)
	log.Log.Debug("SONAR. Get Measues sonar request succesfull :)")
	return resp.Component.Measures, nil

}
