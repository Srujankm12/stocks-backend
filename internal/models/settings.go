package models

type Settings struct {
	Unit         string `json:"unit"`
	Category     string `json:"category"`
	Suppliers    string `json:"suppliers"`
	LeadTime     string `json:"leadtime"`
	StdNonStd    string `json:"stdnonstd"`
	Warranty     string `json:"warranty"`
	IssueAgainst string `json:"issueagainst"`
	Seller       string `json:"seller"`
	Buyer        string `json:"buyer"`
	Region       string `json:"region"`
}
