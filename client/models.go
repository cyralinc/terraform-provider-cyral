package client

type Datamap struct {
	// ID     string         `json:"id"`
	Labels []DatamapLabel `json:"labels,omitempty"`
}

type DatamapLabel struct {
	Name string      `json:"name,omitempty"`
	Info []LabelInfo `json:"info,omitempty"`
}

type LabelInfo struct {
	Repo       string   `json:"repo,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
}
