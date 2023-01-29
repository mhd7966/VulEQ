package models

type Project struct {
	ID           int
	DelProjectID int
	GitURL       string
	GitToken     string
	GitBranch    string
	UserID       int
	Date         string
	Name         string
	Key          string
	Token        string
	SonarToken   string
	ScanCounter  int
}

type ProjectBody struct {
	DelProjectID  int   `json:"del_project_id" binding:"required"`
	GitURL        string `json:"git_url" binding:"required"`
	GitToken      string `json:"git_token" binding:"required"`
	GitBranch     string `json:"git_branch" binding:"required"`
	GitCommitHash string `json:"git_commit_hash" binding:"required"`
	PipelineID    int   `json:"pipeline_id" binding:"required"`
}
