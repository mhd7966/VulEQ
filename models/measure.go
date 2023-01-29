package models

type MeasuresResponse struct {
	Component MeasureComponent `json:"component"`
}

type MeasureComponent struct {
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	Qualifier string    `json:"qualifier"`
	Measures  []Measure `json:"measures"`
}

type Measure struct {
	Metric    string `json:"metric"`
	Value     string `json:"value"`
	BestValue bool   `json:"bestValue"`
}
