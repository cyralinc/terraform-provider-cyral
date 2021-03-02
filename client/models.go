package client

type Datamap struct {
	Labels []DatamapLabel `json:"items,omitempty"`
}

type DatamapLabel struct {
	Repo       string   `json:"repo"`
	Attributes []string `json:"attributes"`
}
