package client

type Datamap struct {
	// ID     string         `json:"id"`
	Labels []DatamapLabel `json:"labels,omitempty"`
}

type DatamapLabel struct {
	Repo       string   `json:"repo"`
	Attributes []string `json:"attributes"`
}
