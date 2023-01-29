package models

type IssuesResponse struct {
	Total       int              `json:"total"`
	P           int              `json:"p"`
	Ps          int              `json:"ps"`
	Paging      PagingStruct     `json:"paging"`
	EffortTotal int              `json:"effortTotal"`
	Issues      []Issue          `json:"issues"`
	Components  []IssueComponent `json:"components"`
	Rules       []Rule           `json:"rules"`
	Users       []User           `json:"users"`
	Languages   []Language       `json:"languages"`
	Facets      []Facet          `json:"facets"`
}

type PagingStruct struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}

type Issue struct {
	Key          string          `json:"key"`
	Rule         string          `json:"rule"`
	Severity     string          `json:"severity"`
	Component    string          `json:"component"`
	Project      string          `json:"project"`
	Line         int             `json:"line"`
	Hash         string          `json:"hash"`
	TextRange    TextRangeStruct `json:"textRange"`
	Flows        []string        `json:"flows"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	Effort       string          `json:"effort"`
	Debt         string          `json:"debt"`
	Author       string          `json:"author"`
	Tags         []string        `json:"tags"`
	Transitions  []string        `json:"transitions"`
	Actions      []string        `json:"actions"`
	Comments     []string        `json:"comments"`
	CreationDate string          `json:"creationDate"`
	UpdateDate   string          `json:"updateDate"`
	Type         string          `json:"type"`
	Scope        string          `json:"scope"`
}

type TextRangeStruct struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartOffset int `json:"startOffset"`
	EndOffset   int `json:"endOffset"`
}

type IssueComponent struct {
	Key       string `json:"key"`
	Enabled   bool   `json:"enabled"`
	Qualifier string `json:"qualifier"`
	Name      string `json:"name"`
	LongName  string `json:"longName"`
	Path      string `json:"path"`
}

type Rule struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Lang     string `json:"lang"`
	Status   string `json:"status"`
	LangName string `json:"langName"`
}

type User struct {
	Login  string `json:"login"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type Language struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Facet struct {
	Property string  `json:"property"`
	Values   []Value `json:"values"`
}

type Value struct {
	Val   string `json:"val"`
	Count int    `json:"count"`
}

type ReturnIssueResponse struct {
	Component string `json:"component"`
	Line      int    `json:"line"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
	Type      string `json:"type"`
}
