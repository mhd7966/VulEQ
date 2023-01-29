package models

type Scan struct {
	ID               int
	ProjectID        int
	GitCommitHash    string
	PipelineID       int
	NumBug           string
	NumVulnerability string
	NumDebt          string
	NumCodeSmell     string
	NumFile          string
	NumDuplicateLine string
	LineCode         string
	LineComment      string
}
